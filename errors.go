// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import "errors"

var (
	// ErrInvalidRequest indicates invalid translation request input.
	ErrInvalidRequest = errors.New("invalid translation request")

	// ErrUnsupportedLanguage indicates provider does not support language.
	ErrUnsupportedLanguage = errors.New("unsupported language")

	// ErrBatchTooLarge indicates request exceeds batch constraints.
	ErrBatchTooLarge = errors.New("batch too large")

	// ErrTextTooLong indicates single text item exceeds configured limit.
	ErrTextTooLong = errors.New("text too long")

	// ErrProviderTemporary indicates retryable upstream provider failure.
	ErrProviderTemporary = errors.New("temporary provider failure")

	// ErrProviderPermanent indicates non-retryable upstream provider failure.
	ErrProviderPermanent = errors.New("permanent provider failure")
)
