// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"sync"
	"time"
)

// RateLimitOptions controls rate limit wrapper behavior.
type RateLimitOptions struct {
	// MinInterval is minimal delay between outbound requests.
	// Set `0` to disable throttling.
	MinInterval time.Duration `json:"min_interval,omitempty" yaml:"min_interval,omitempty" jsonschema:"minimum=0"`
}

// RateLimitTranslator enforces minimum interval between Translate calls.
type RateLimitTranslator struct {
	// nextAllowed is next allowed call timestamp.
	nextAllowed time.Time

	// base is wrapped translator.
	base Translator

	// options stores limiter settings.
	options RateLimitOptions

	// lock guards nextAllowed.
	lock sync.Mutex
}

// NewRateLimitTranslator creates rate-limit wrapper around translator.
func NewRateLimitTranslator(
	base Translator,
	options RateLimitOptions,
) *RateLimitTranslator {
	return &RateLimitTranslator{
		base:    base,
		options: options,
	}
}

// Capabilities returns wrapped translator capabilities.
func (translator *RateLimitTranslator) Capabilities() Capabilities {
	return translator.base.Capabilities()
}

// Translate executes wrapped Translate with rate limiting.
func (translator *RateLimitTranslator) Translate(
	ctx context.Context,
	request Request,
) (Result, error) {
	if err := translator.waitTurn(ctx); err != nil {
		return Result{}, err
	}

	return translator.base.Translate(ctx, request)
}

// waitTurn blocks until this call is allowed by rate limiter.
func (translator *RateLimitTranslator) waitTurn(ctx context.Context) error {
	minInterval := translator.options.MinInterval
	if minInterval <= 0 {
		return nil
	}

	translator.lock.Lock()
	now := time.Now()
	waitFor := time.Duration(0)
	if now.Before(translator.nextAllowed) {
		waitFor = translator.nextAllowed.Sub(now)
	}
	nextBase := now
	if waitFor > 0 {
		nextBase = translator.nextAllowed
	}
	translator.nextAllowed = nextBase.Add(minInterval)
	translator.lock.Unlock()

	if waitFor <= 0 {
		return nil
	}

	timer := time.NewTimer(waitFor)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
