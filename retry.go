// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	// DefaultRetryAttempts is default number of retry attempts.
	DefaultRetryAttempts = 3

	// DefaultRetryDelay is default delay before first retry.
	DefaultRetryDelay = 300 * time.Millisecond

	// DefaultRetryBackoff is default retry delay multiplier.
	DefaultRetryBackoff = 2.0
)

// RetryOptions controls retry wrapper behavior.
type RetryOptions struct {
	// Attempts defines maximum total attempts including the first try.
	Attempts int `json:"attempts,omitempty" yaml:"attempts,omitempty"`

	// Delay defines delay before first retry.
	Delay time.Duration `json:"delay,omitempty" yaml:"delay,omitempty"`

	// MaxDelay defines maximum delay cap between retries.
	MaxDelay time.Duration `json:"max_delay,omitempty" yaml:"max_delay,omitempty"`

	// Backoff multiplies delay after each failed attempt.
	Backoff float64 `json:"backoff,omitempty" yaml:"backoff,omitempty"`

	// RetryPermanent allows retry for permanent provider errors.
	RetryPermanent bool `json:"retry_permanent,omitempty" yaml:"retry_permanent,omitempty"`
}

// RetryTranslator retries wrapped translator on retryable errors.
type RetryTranslator struct {
	// base is wrapped translator.
	base Translator

	// options controls retry logic.
	options RetryOptions
}

// NewRetryTranslator creates retry wrapper around translator.
func NewRetryTranslator(base Translator, options RetryOptions) *RetryTranslator {
	return &RetryTranslator{
		base:    base,
		options: normalizeRetryOptions(options),
	}
}

// Capabilities returns wrapped translator capabilities.
func (translator *RetryTranslator) Capabilities() Capabilities {
	return translator.base.Capabilities()
}

// Translate calls wrapped translator with retry policy.
func (translator *RetryTranslator) Translate(
	ctx context.Context,
	request Request,
) (Result, error) {
	options := normalizeRetryOptions(translator.options)

	var lastErr error
	delay := options.Delay
	for attempt := 1; attempt <= options.Attempts; attempt++ {
		result, err := translator.base.Translate(ctx, request)
		if err == nil {
			return result, nil
		}

		lastErr = err
		if !shouldRetryError(err, options.RetryPermanent) {
			return Result{}, err
		}
		if attempt == options.Attempts {
			break
		}

		if err := waitRetryDelay(ctx, delay); err != nil {
			return Result{}, err
		}

		delay = nextRetryDelay(delay, options)
	}

	return Result{}, fmt.Errorf(
		"retry attempts exhausted (%d): %w",
		options.Attempts,
		lastErr,
	)
}

// normalizeRetryOptions fills retry defaults.
func normalizeRetryOptions(options RetryOptions) RetryOptions {
	if options.Attempts <= 0 {
		options.Attempts = DefaultRetryAttempts
	}
	if options.Delay <= 0 {
		options.Delay = DefaultRetryDelay
	}
	if options.Backoff <= 0 {
		options.Backoff = DefaultRetryBackoff
	}

	return options
}

// shouldRetryError reports whether error can be retried.
func shouldRetryError(err error, retryPermanent bool) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if errors.Is(err, ErrInvalidRequest) {
		return false
	}
	if errors.Is(err, ErrProviderTemporary) {
		return true
	}
	if retryPermanent && errors.Is(err, ErrProviderPermanent) {
		return true
	}

	return false
}

// waitRetryDelay blocks for retry delay or context cancellation.
func waitRetryDelay(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// nextRetryDelay computes next retry delay with backoff and cap.
func nextRetryDelay(delay time.Duration, options RetryOptions) time.Duration {
	next := time.Duration(float64(delay) * options.Backoff)
	if next <= 0 {
		next = delay
	}
	if options.MaxDelay > 0 && next > options.MaxDelay {
		return options.MaxDelay
	}

	return next
}
