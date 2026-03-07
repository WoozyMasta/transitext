// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"strings"
)

//nolint:gosec // Marker tokens are not credentials.
const (
	// defaultContextToken is marker before context payload.
	defaultContextToken = "[[|:|]]"

	// defaultTextToken is marker before source text payload.
	defaultTextToken = "[[|>|]]"

	// defaultEndToken is marker after source text payload.
	defaultEndToken = "[[|#|]]"
)

// ContextOptions controls experimental context-injection wrapper behavior.
type ContextOptions struct {
	// Context is default translation context injected for each item.
	Context string `json:"context,omitempty" yaml:"context,omitempty" jsonschema:"maxLength=2000"`

	// ContextByID overrides context for specific item IDs.
	ContextByID map[string]string `json:"context_by_id,omitempty" yaml:"context_by_id,omitempty" jsonschema:"maxProperties=2048"`

	// ContextToken is marker inserted before context payload.
	ContextToken string `json:"context_token,omitempty" yaml:"context_token,omitempty" jsonschema:"minLength=1,maxLength=64"`

	// TextToken is marker inserted before source text payload.
	TextToken string `json:"text_token,omitempty" yaml:"text_token,omitempty" jsonschema:"minLength=1,maxLength=64"`

	// EndToken is marker appended after source text payload.
	EndToken string `json:"end_token,omitempty" yaml:"end_token,omitempty" jsonschema:"minLength=1,maxLength=64"`
}

// ContextTranslator injects context into item text and strips marker envelope.
type ContextTranslator struct {
	// base is wrapped translator.
	base Translator

	// options controls context wrapping behavior.
	options ContextOptions
}

// NewContextTranslator creates context-injection wrapper around translator.
func NewContextTranslator(base Translator, options ContextOptions) *ContextTranslator {
	return &ContextTranslator{
		base:    base,
		options: normalizeContextOptions(options),
	}
}

// Capabilities returns wrapped translator capabilities.
func (translator *ContextTranslator) Capabilities() Capabilities {
	return translator.base.Capabilities()
}

// Translate injects context markers, delegates translation, then strips them.
func (translator *ContextTranslator) Translate(
	ctx context.Context,
	request Request,
) (Result, error) {
	modified := request
	modified.Items = append([]Item(nil), request.Items...)

	for index := range modified.Items {
		itemContext := translator.contextForItem(modified.Items[index].ID)
		if itemContext == "" {
			continue
		}

		modified.Items[index].Text = buildContextEnvelope(
			modified.Items[index].Text,
			itemContext,
			translator.options,
		)
	}

	result, err := translator.base.Translate(ctx, modified)
	if err != nil {
		return Result{}, err
	}

	for index := range result.Items {
		if stripped, ok := stripContextEnvelope(result.Items[index].Text, translator.options); ok {
			result.Items[index].Text = stripped
		}
	}

	return result, nil
}

// normalizeContextOptions fills missing marker defaults.
func normalizeContextOptions(options ContextOptions) ContextOptions {
	if strings.TrimSpace(options.ContextToken) == "" {
		options.ContextToken = defaultContextToken
	}
	if strings.TrimSpace(options.TextToken) == "" {
		options.TextToken = defaultTextToken
	}
	if strings.TrimSpace(options.EndToken) == "" {
		options.EndToken = defaultEndToken
	}

	return options
}

// contextForItem returns effective context for item by ID.
func (translator *ContextTranslator) contextForItem(id string) string {
	if translator.options.ContextByID != nil {
		if value := strings.TrimSpace(translator.options.ContextByID[id]); value != "" {
			return value
		}
	}

	return strings.TrimSpace(translator.options.Context)
}

// buildContextEnvelope creates marker envelope with context + source text.
func buildContextEnvelope(text string, itemContext string, options ContextOptions) string {
	var builder strings.Builder
	builder.Grow(
		len(options.ContextToken) +
			len(options.TextToken) +
			len(options.EndToken) +
			len(itemContext) + len(text) + 8,
	)
	_, _ = builder.WriteString(options.ContextToken)
	builder.WriteByte(' ')
	_, _ = builder.WriteString(itemContext)
	builder.WriteByte('\n')
	_, _ = builder.WriteString(options.TextToken)
	builder.WriteByte(' ')
	_, _ = builder.WriteString(text)
	builder.WriteByte('\n')
	_, _ = builder.WriteString(options.EndToken)

	return builder.String()
}

// stripContextEnvelope extracts translated payload between text and end tokens.
func stripContextEnvelope(text string, options ContextOptions) (string, bool) {
	textIndex := strings.Index(text, options.TextToken)
	if textIndex < 0 {
		return "", false
	}

	start := textIndex + len(options.TextToken)
	for start < len(text) && (text[start] == ' ' || text[start] == '\n' || text[start] == '\r' || text[start] == '\t') {
		start++
	}

	end := strings.Index(text[start:], options.EndToken)
	if end < 0 {
		return sanitizeExtractedText(text[start:]), true
	}

	return sanitizeExtractedText(text[start : start+end]), true
}

// sanitizeExtractedText removes marker-like trailing artifacts.
func sanitizeExtractedText(text string) string {
	value := strings.TrimSpace(text)
	for {
		index := strings.LastIndex(value, "\n[[")
		if index < 0 {
			break
		}

		tail := strings.TrimSpace(value[index+1:])
		if !isBracketMarker(tail) {
			break
		}

		value = strings.TrimSpace(value[:index])
	}

	if isBracketMarker(value) {
		return ""
	}

	return value
}

// isBracketMarker reports whether value looks like synthetic marker token.
func isBracketMarker(value string) bool {
	if !strings.HasPrefix(value, "[[") || !strings.HasSuffix(value, "]]") {
		return false
	}
	if strings.ContainsAny(value, " \t") {
		return false
	}
	if len(value) < 4 || len(value) > 64 {
		return false
	}

	return true
}
