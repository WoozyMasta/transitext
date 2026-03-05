// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"fmt"
	"strings"
)

const (
	// DefaultMaxBatchItems is the default upper bound for batch size by item.
	DefaultMaxBatchItems = 50

	// DefaultMaxBatchChars is the default upper bound for batch size by chars.
	DefaultMaxBatchChars = 5000
)

// ValidateRequest validates request structure and basic constraints.
func ValidateRequest(request Request) error {
	if strings.TrimSpace(request.TargetLang) == "" {
		return fmt.Errorf("target_lang is required: %w", ErrInvalidRequest)
	}
	if len(request.Items) == 0 {
		return fmt.Errorf("items are required: %w", ErrInvalidRequest)
	}

	for index, item := range request.Items {
		if item.Text == "" {
			return fmt.Errorf(
				"items[%d].text is required: %w",
				index,
				ErrInvalidRequest,
			)
		}

		if request.Batch.MaxTextChars > 0 && len(item.Text) > request.Batch.MaxTextChars {
			return fmt.Errorf(
				"items[%d] exceeds max_text_chars %d: %w",
				index,
				request.Batch.MaxTextChars,
				ErrTextTooLong,
			)
		}
	}

	return nil
}

// NormalizeBatchOptions returns safe and complete batch options.
func NormalizeBatchOptions(options BatchOptions) BatchOptions {
	if options.MaxItems <= 0 {
		options.MaxItems = DefaultMaxBatchItems
	}
	if options.MaxChars <= 0 {
		options.MaxChars = DefaultMaxBatchChars
	}
	if options.OnOverflow == "" {
		options.OnOverflow = OverflowSplit
	}

	return options
}

// ResolveBatchOptions merges request batch options with provider capabilities.
func ResolveBatchOptions(
	options BatchOptions,
	capabilities Capabilities,
) BatchOptions {
	if options.MaxItems <= 0 && capabilities.MaxBatchItems > 0 {
		options.MaxItems = capabilities.MaxBatchItems
	}
	if options.MaxChars <= 0 && capabilities.MaxBatchChars > 0 {
		options.MaxChars = capabilities.MaxBatchChars
	}
	if options.MaxTextChars <= 0 && capabilities.MaxTextChars > 0 {
		options.MaxTextChars = capabilities.MaxTextChars
	}

	return NormalizeBatchOptions(options)
}

// SplitRequest splits request into batches according to options.
// Returned batches reuse input item storage and should be treated as immutable.
func SplitRequest(request Request, options BatchOptions) ([]Request, error) {
	options = NormalizeBatchOptions(options)

	batches := make([]Request, 0, max(1, len(request.Items)/options.MaxItems+1))
	start := 0
	currentChars := 0

	for index, item := range request.Items {
		if options.MaxTextChars > 0 && len(item.Text) > options.MaxTextChars {
			return nil, fmt.Errorf(
				"item %q exceeds max text chars %d: %w",
				item.ID,
				options.MaxTextChars,
				ErrTextTooLong,
			)
		}

		itemChars := len(item.Text)
		if itemChars > options.MaxChars {
			return nil, fmt.Errorf(
				"item %q exceeds max chars %d: %w",
				item.ID,
				options.MaxChars,
				ErrBatchTooLarge,
			)
		}

		wouldOverflowItems := index-start >= options.MaxItems
		wouldOverflowChars := currentChars+itemChars > options.MaxChars

		if wouldOverflowItems || wouldOverflowChars {
			if options.OnOverflow == OverflowError {
				return nil, fmt.Errorf(
					"request exceeds batch limit: %w",
					ErrBatchTooLarge,
				)
			}

			if start < index {
				batch := request
				batch.Items = request.Items[start:index]
				batches = append(batches, batch)
			}

			start = index
			currentChars = 0
		}

		currentChars += itemChars
	}

	if start < len(request.Items) {
		batch := request
		batch.Items = request.Items[start:]
		batches = append(batches, batch)
	}

	return batches, nil
}

// TranslateBatchFunc translates one prepared request batch.
type TranslateBatchFunc func(
	ctx context.Context,
	request Request,
) ([]TranslatedItem, error)

// TranslateBatches validates request, splits by limits, and translates batches.
func TranslateBatches(
	ctx context.Context,
	request Request,
	capabilities Capabilities,
	translateBatch TranslateBatchFunc,
) ([]TranslatedItem, error) {
	if translateBatch == nil {
		return nil, fmt.Errorf("translate batch function is required: %w", ErrInvalidRequest)
	}
	if err := ValidateRequest(request); err != nil {
		return nil, err
	}

	batchOptions := ResolveBatchOptions(request.Batch, capabilities)
	batches, err := SplitRequest(request, batchOptions)
	if err != nil {
		return nil, err
	}

	items := make([]TranslatedItem, 0, len(request.Items))
	for _, batch := range batches {
		batchItems, err := translateBatch(ctx, batch)
		if err != nil {
			return nil, err
		}

		items = append(items, batchItems...)
	}

	return items, nil
}
