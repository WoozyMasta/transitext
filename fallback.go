// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"errors"
	"fmt"
)

// FallbackTranslator runs translators in order until one succeeds.
type FallbackTranslator struct {
	// translators contains ordered fallback chain.
	translators []Translator
}

// NewFallbackTranslator creates fallback chain wrapper.
func NewFallbackTranslator(translators ...Translator) *FallbackTranslator {
	chain := make([]Translator, 0, len(translators))
	for _, translator := range translators {
		if translator == nil {
			continue
		}

		chain = append(chain, translator)
	}

	return &FallbackTranslator{translators: chain}
}

// Capabilities returns capabilities of primary translator.
func (translator *FallbackTranslator) Capabilities() Capabilities {
	if len(translator.translators) == 0 {
		return Capabilities{}
	}

	return translator.translators[0].Capabilities()
}

// Translate executes translators in order until first success.
func (translator *FallbackTranslator) Translate(
	ctx context.Context,
	request Request,
) (Result, error) {
	if len(translator.translators) == 0 {
		return Result{}, fmt.Errorf("fallback chain is empty: %w", ErrInvalidRequest)
	}

	errorsList := make([]error, 0, len(translator.translators))
	for _, current := range translator.translators {
		result, err := current.Translate(ctx, request)
		if err == nil {
			if result.Provider == "" {
				result.Provider = current.Capabilities().Provider
			}

			return result, nil
		}

		errorsList = append(errorsList, err)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return Result{}, err
		}
		if errors.Is(err, ErrInvalidRequest) {
			return Result{}, err
		}
	}

	return Result{}, errors.Join(errorsList...)
}
