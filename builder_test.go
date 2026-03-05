// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"testing"
	"time"
)

func TestBuilderPipeline(t *testing.T) {
	t.Parallel()

	primary := &builderMockTranslator{
		capabilities: Capabilities{Provider: "primary"},
		result: Result{
			Provider: "primary",
			Items:    []TranslatedItem{{ID: "1", Text: "ok"}},
		},
	}
	fallback := &builderMockTranslator{
		capabilities: Capabilities{Provider: "fallback"},
		result: Result{
			Provider: "fallback",
			Items:    []TranslatedItem{{ID: "1", Text: "fallback"}},
		},
	}

	pipeline := Wrap(primary).
		Fallback(fallback).
		Retry(RetryOptions{Attempts: 2, Delay: time.Millisecond}).
		RateLimit(RateLimitOptions{MinInterval: time.Millisecond}).
		Cache(CacheOptions{}).
		Build()

	if pipeline == nil {
		t.Fatal("pipeline is nil")
	}

	request := Request{
		TargetLang: "en",
		Items:      []Item{{ID: "1", Text: "bonjour"}},
	}
	first, err := pipeline.Translate(context.Background(), request)
	if err != nil {
		t.Fatalf("first Translate error: %v", err)
	}
	second, err := pipeline.Translate(context.Background(), request)
	if err != nil {
		t.Fatalf("second Translate error: %v", err)
	}

	if first.Provider != "primary" || second.Provider != "primary" {
		t.Fatalf("unexpected provider %q/%q", first.Provider, second.Provider)
	}
	if primary.calls != 1 {
		t.Fatalf("primary calls = %d, want 1 due to cache", primary.calls)
	}
	if fallback.calls != 0 {
		t.Fatalf("fallback calls = %d, want 0", fallback.calls)
	}
}

// builderMockTranslator is deterministic translator for builder tests.
type builderMockTranslator struct {
	// capabilities returns static capabilities.
	capabilities Capabilities

	// result returns static result.
	result Result

	// calls tracks translate call count.
	calls int
}

// Capabilities returns static capabilities.
func (translator *builderMockTranslator) Capabilities() Capabilities {
	return translator.capabilities
}

// Translate returns static result and increments call counter.
func (translator *builderMockTranslator) Translate(
	_ context.Context,
	_ Request,
) (Result, error) {
	translator.calls++
	return translator.result, nil
}
