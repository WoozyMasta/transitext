// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"errors"
	"fmt"
)

// PartialOptions controls partial-result wrapper behavior.
type PartialOptions struct {
	// ItemRetries is extra retry count for one failed item.
	ItemRetries int `json:"item_retries,omitempty" yaml:"item_retries,omitempty" jsonschema:"minimum=0,maximum=16"`

	// ContinueOnTemporary keeps processing after temporary provider errors.
	ContinueOnTemporary bool `json:"continue_on_temporary,omitempty" yaml:"continue_on_temporary,omitempty"`

	// ContinueOnContext keeps processing after context cancellation/deadline.
	ContinueOnContext bool `json:"continue_on_context,omitempty" yaml:"continue_on_context,omitempty"`
}

// PartialTranslator processes items one-by-one and returns partial results.
type PartialTranslator struct {
	// base is wrapped translator.
	base Translator

	// options controls partial handling behavior.
	options PartialOptions
}

// NewPartialTranslator creates partial-result wrapper around translator.
func NewPartialTranslator(base Translator, options PartialOptions) *PartialTranslator {
	return &PartialTranslator{
		base:    base,
		options: normalizePartialOptions(options),
	}
}

// Capabilities returns wrapped translator capabilities.
func (translator *PartialTranslator) Capabilities() Capabilities {
	return translator.base.Capabilities()
}

// Translate processes each item independently and returns partial progress.
func (translator *PartialTranslator) Translate(
	ctx context.Context,
	request Request,
) (Result, error) {
	if err := ValidateRequest(request); err != nil {
		return Result{}, err
	}

	result := Result{
		Provider: translator.base.Capabilities().Provider,
		Items:    make([]TranslatedItem, 0, len(request.Items)),
	}
	for index, item := range request.Items {
		translated, err := translator.translateItemWithRetry(ctx, request, item)
		if err == nil {
			result.Items = append(result.Items, translated)
			continue
		}

		if translator.shouldStop(err) {
			result.Items = append(result.Items, failedItem(item, err))
			for _, remaining := range request.Items[index+1:] {
				result.Items = append(result.Items, skippedItem(remaining))
			}

			return result, fmt.Errorf(
				"partial translation stopped at item %d/%d (%q): %w",
				index+1,
				len(request.Items),
				item.ID,
				err,
			)
		}

		result.Items = append(result.Items, failedItem(item, err))
	}

	return result, nil
}

// normalizePartialOptions fills wrapper defaults.
func normalizePartialOptions(options PartialOptions) PartialOptions {
	if options.ItemRetries < 0 {
		options.ItemRetries = 0
	}

	return options
}

// translateItemWithRetry translates one item with configured retries.
func (translator *PartialTranslator) translateItemWithRetry(
	ctx context.Context,
	base Request,
	item Item,
) (TranslatedItem, error) {
	subrequest := base
	subrequest.Items = []Item{item}

	var lastErr error
	attempts := translator.options.ItemRetries + 1
	for attempt := 1; attempt <= attempts; attempt++ {
		result, err := translator.base.Translate(ctx, subrequest)
		if err == nil {
			if len(result.Items) == 0 {
				return TranslatedItem{}, fmt.Errorf(
					"provider returned empty items for %q: %w",
					item.ID,
					ErrProviderPermanent,
				)
			}

			output := result.Items[0]
			if output.ID == "" {
				output.ID = item.ID
			}

			return output, nil
		}

		lastErr = err
		if translator.shouldStop(err) {
			break
		}
	}

	return TranslatedItem{}, lastErr
}

// shouldStop reports whether processing must stop after this error.
func (translator *PartialTranslator) shouldStop(err error) bool {
	if err == nil {
		return false
	}
	if !translator.options.ContinueOnContext &&
		(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
		return true
	}
	if !translator.options.ContinueOnTemporary && errors.Is(err, ErrProviderTemporary) {
		return true
	}

	return false
}

// failedItem builds translated item with failure payload.
func failedItem(item Item, err error) TranslatedItem {
	return TranslatedItem{
		ID: item.ID,
		Error: &ItemError{
			Code:      classifyErrorCode(err),
			Message:   err.Error(),
			Temporary: errors.Is(err, ErrProviderTemporary),
		},
	}
}

// skippedItem builds translated item skipped after stop condition.
func skippedItem(item Item) TranslatedItem {
	return TranslatedItem{
		ID: item.ID,
		Error: &ItemError{
			Code:      "not_processed",
			Message:   "not processed after previous stop condition",
			Temporary: true,
		},
	}
}

// classifyErrorCode maps error to stable per-item code.
func classifyErrorCode(err error) string {
	switch {
	case errors.Is(err, context.Canceled):
		return "context_canceled"
	case errors.Is(err, context.DeadlineExceeded):
		return "context_deadline"
	case errors.Is(err, ErrInvalidRequest):
		return "invalid_request"
	case errors.Is(err, ErrBatchTooLarge):
		return "batch_too_large"
	case errors.Is(err, ErrProviderTemporary):
		return "provider_temporary"
	case errors.Is(err, ErrProviderPermanent):
		return "provider_permanent"
	default:
		return "provider_error"
	}
}
