// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
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

// SplitRequest splits request into batches according to options.
// Returned batches reuse input item storage and should be treated as immutable.
func SplitRequest(request Request, options BatchOptions) ([]Request, error) {
	options = NormalizeBatchOptions(options)

	batches := make([]Request, 0, max(1, len(request.Items)/options.MaxItems+1))
	start := 0
	currentChars := 0

	for index, item := range request.Items {
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
