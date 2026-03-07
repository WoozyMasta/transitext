// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

var longTextSeparators = []string{
	"\r\n\r\n",
	"\n\n",
	"\r\n",
	"\n",
	". ",
	"! ",
	"? ",
	"。",
	"！",
	"？",
	"; ",
	" ",
}

// LongTextOptions controls single-item overflow handling.
type LongTextOptions struct {
	// ErrorOnOverflow fails when item exceeds MaxTextChars.
	// By default long items are split and merged automatically.
	ErrorOnOverflow bool `json:"error_on_overflow,omitempty" yaml:"error_on_overflow,omitempty"`

	// MaxTextChars overrides provider single-item limit.
	MaxTextChars int `json:"max_text_chars,omitempty" yaml:"max_text_chars,omitempty" jsonschema:"minimum=1,maximum=1000000"`
}

// LongTextTranslator splits oversized items and merges translated parts back.
type LongTextTranslator struct {
	// base handles actual translation transport.
	base Translator

	// options configures long item processing behavior.
	options LongTextOptions
}

// NewLongTextTranslator creates wrapper for oversized item handling.
func NewLongTextTranslator(
	base Translator,
	options LongTextOptions,
) *LongTextTranslator {
	return &LongTextTranslator{
		base:    base,
		options: options,
	}
}

// Capabilities reports wrapped translator capabilities.
func (translator *LongTextTranslator) Capabilities() Capabilities {
	capabilities := translator.base.Capabilities()
	if translator.options.MaxTextChars > 0 {
		capabilities.MaxTextChars = translator.options.MaxTextChars
	}

	return capabilities
}

// Translate handles long text items before forwarding request to base.
func (translator *LongTextTranslator) Translate(
	ctx context.Context,
	request Request,
) (Result, error) {
	if err := ValidateRequest(request); err != nil {
		return Result{}, err
	}

	maxTextChars := translator.resolveMaxTextChars()
	if maxTextChars <= 0 {
		return translator.base.Translate(ctx, request)
	}

	expanded, mapping, err := expandLongItems(
		request.Items,
		maxTextChars,
		translator.options.ErrorOnOverflow,
	)
	if err != nil {
		return Result{}, err
	}
	if len(expanded) == len(request.Items) {
		return translator.base.Translate(ctx, request)
	}

	patched := request
	patched.Items = expanded

	translated, err := translator.base.Translate(ctx, patched)
	if err != nil {
		return Result{}, err
	}
	if len(translated.Items) != len(expanded) {
		return Result{}, fmt.Errorf(
			"provider returned %d items for %d expanded items: %w",
			len(translated.Items),
			len(expanded),
			ErrProviderPermanent,
		)
	}

	return collapseLongResult(translated, request.Items, mapping), nil
}

// resolveMaxTextChars returns active single-item limit.
func (translator *LongTextTranslator) resolveMaxTextChars() int {
	if translator.options.MaxTextChars > 0 {
		return translator.options.MaxTextChars
	}

	return translator.base.Capabilities().MaxTextChars
}

// longTextGroup maps original item to expanded part index range.
type longTextGroup struct {
	// start is first expanded index.
	start int

	// end is exclusive expanded index.
	end int
}

// expandLongItems expands oversized source items into smaller parts.
func expandLongItems(
	items []Item,
	maxTextChars int,
	errorOnOverflow bool,
) ([]Item, []longTextGroup, error) {
	expanded := make([]Item, 0, len(items))
	mapping := make([]longTextGroup, len(items))

	for index, item := range items {
		parts := splitLongText(item.Text, maxTextChars)
		if len(parts) > 1 && errorOnOverflow {
			return nil, nil, fmt.Errorf(
				"item %q exceeds max_text_chars %d: %w",
				item.ID,
				maxTextChars,
				ErrTextTooLong,
			)
		}

		if len(parts) == 0 {
			parts = []string{item.Text}
		}

		start := len(expanded)

		for partIndex, part := range parts {
			expanded = append(expanded, Item{
				ID:   buildLongTextPartID(item.ID, partIndex),
				Text: part,
			})
		}

		mapping[index] = longTextGroup{
			start: start,
			end:   len(expanded),
		}
	}

	return expanded, mapping, nil
}

// collapseLongResult merges translated expanded items back to original shape.
func collapseLongResult(
	translated Result,
	original []Item,
	mapping []longTextGroup,
) Result {
	merged := translated
	merged.Items = make([]TranslatedItem, len(original))

	for index := range original {
		group := mapping[index]
		first := translated.Items[group.start]
		builder := strings.Builder{}
		builder.Grow(len(first.Text))
		builder.WriteString(first.Text)

		itemError := first.Error
		detectedSource := first.DetectedSource

		for partIndex := group.start + 1; partIndex < group.end; partIndex++ {
			part := translated.Items[partIndex]
			builder.WriteString(part.Text)
			if itemError == nil && part.Error != nil {
				itemError = part.Error
			}
			if detectedSource == "" {
				detectedSource = part.DetectedSource
			}
		}

		merged.Items[index] = TranslatedItem{
			ID:             original[index].ID,
			Text:           builder.String(),
			DetectedSource: detectedSource,
			Error:          itemError,
		}
	}

	return merged
}

// splitLongText splits text by paragraph/sentence separators and hard fallback.
func splitLongText(text string, maxTextChars int) []string {
	if maxTextChars <= 0 || countTextChars(text) <= maxTextChars {
		return []string{text}
	}

	parts := splitLongTextByRule(text, maxTextChars, 0)
	return mergeShortParts(parts, maxTextChars)
}

// splitLongTextByRule recursively applies separators by priority.
func splitLongTextByRule(text string, maxTextChars int, ruleIndex int) []string {
	if countTextChars(text) <= maxTextChars {
		return []string{text}
	}
	if ruleIndex >= len(longTextSeparators) {
		return splitTextHard(text, maxTextChars)
	}

	parts := splitKeepSeparator(text, longTextSeparators[ruleIndex])
	if len(parts) <= 1 {
		return splitLongTextByRule(text, maxTextChars, ruleIndex+1)
	}

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}

		if countTextChars(part) <= maxTextChars {
			result = append(result, part)
			continue
		}

		result = append(result, splitLongTextByRule(part, maxTextChars, ruleIndex+1)...)
	}

	return result
}

// splitKeepSeparator splits by separator and keeps it on the left fragment.
func splitKeepSeparator(text string, separator string) []string {
	if separator == "" {
		return []string{text}
	}

	result := make([]string, 0, 4)
	rest := text

	for {
		index := strings.Index(rest, separator)
		if index < 0 {
			break
		}

		cut := index + len(separator)
		result = append(result, rest[:cut])
		rest = rest[cut:]
	}

	if rest != "" {
		result = append(result, rest)
	}
	if len(result) == 0 {
		return []string{text}
	}

	return result
}

// splitTextHard performs rune-safe hard split by length.
func splitTextHard(text string, maxTextChars int) []string {
	runes := []rune(text)
	if len(runes) == 0 {
		return []string{text}
	}

	result := make([]string, 0, len(runes)/maxTextChars+1)
	for start := 0; start < len(runes); start += maxTextChars {
		end := min(start+maxTextChars, len(runes))

		result = append(result, string(runes[start:end]))
	}

	return result
}

// mergeShortParts merges adjacent fragments while they fit in limit.
func mergeShortParts(parts []string, maxTextChars int) []string {
	if len(parts) < 2 {
		return parts
	}

	merged := make([]string, 0, len(parts))
	current := parts[0]
	currentChars := countTextChars(current)

	for index := 1; index < len(parts); index++ {
		next := parts[index]
		nextChars := countTextChars(next)

		if currentChars+nextChars <= maxTextChars {
			current += next
			currentChars += nextChars
			continue
		}

		merged = append(merged, current)
		current = next
		currentChars = nextChars
	}

	merged = append(merged, current)
	return merged
}

// buildLongTextPartID builds deterministic id for expanded item.
func buildLongTextPartID(id string, partIndex int) string {
	if id == "" {
		return ""
	}

	return fmt.Sprintf("%s#%d", id, partIndex+1)
}

// countTextChars counts runes in text.
func countTextChars(text string) int {
	return utf8.RuneCountInString(text)
}
