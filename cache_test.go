// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"testing"
)

func TestCacheTranslatorCachesByDefaultRequestShape(t *testing.T) {
	t.Parallel()

	base := &cacheMockTranslator{
		capabilities: Capabilities{Provider: "mock"},
		result: Result{
			Provider: "mock",
			Items: []TranslatedItem{
				{ID: "1", Text: "hello"},
			},
		},
	}
	cache := NewCacheTranslator(base, CacheOptions{})

	request := Request{
		SourceLang: "auto",
		TargetLang: "en",
		Items: []Item{
			{ID: "1", Text: "bonjour"},
		},
		Metadata: map[string]string{"job": "a"},
	}

	first, err := cache.Translate(context.Background(), request)
	if err != nil {
		t.Fatalf("first translate error: %v", err)
	}
	second, err := cache.Translate(context.Background(), request)
	if err != nil {
		t.Fatalf("second translate error: %v", err)
	}
	if base.calls != 1 {
		t.Fatalf("base calls = %d, want 1", base.calls)
	}
	if first.Items[0].Text != "hello" || second.Items[0].Text != "hello" {
		t.Fatalf("unexpected result: %#v / %#v", first, second)
	}
	if cache.Len() != 1 {
		t.Fatalf("cache len = %d, want 1", cache.Len())
	}
}

func TestCacheTranslatorIncludeMetadata(t *testing.T) {
	t.Parallel()

	base := &cacheMockTranslator{
		capabilities: Capabilities{Provider: "mock"},
		result: Result{
			Provider: "mock",
			Items: []TranslatedItem{
				{ID: "1", Text: "hello"},
			},
		},
	}
	cache := NewCacheTranslator(base, CacheOptions{
		IncludeMetadata: true,
	})

	request := Request{
		SourceLang: "auto",
		TargetLang: "en",
		Items: []Item{
			{ID: "1", Text: "bonjour"},
		},
	}

	request.Metadata = map[string]string{"job": "a"}
	_, err := cache.Translate(context.Background(), request)
	if err != nil {
		t.Fatalf("first translate error: %v", err)
	}
	request.Metadata = map[string]string{"job": "b"}
	_, err = cache.Translate(context.Background(), request)
	if err != nil {
		t.Fatalf("second translate error: %v", err)
	}

	if base.calls != 2 {
		t.Fatalf("base calls = %d, want 2", base.calls)
	}
}

func TestCacheTranslatorClear(t *testing.T) {
	t.Parallel()

	base := &cacheMockTranslator{
		capabilities: Capabilities{Provider: "mock"},
		result:       Result{Provider: "mock"},
	}
	cache := NewCacheTranslator(base, CacheOptions{})

	request := Request{
		TargetLang: "en",
		Items:      []Item{{ID: "1", Text: "x"}},
	}
	_, _ = cache.Translate(context.Background(), request)
	if cache.Len() != 1 {
		t.Fatalf("cache len = %d, want 1", cache.Len())
	}
	cache.Clear()
	if cache.Len() != 0 {
		t.Fatalf("cache len = %d, want 0", cache.Len())
	}
}

// cacheMockTranslator is deterministic translator for cache tests.
type cacheMockTranslator struct {
	// capabilities returns static capabilities.
	capabilities Capabilities

	// result is returned for each successful translate call.
	result Result

	// calls tracks number of calls.
	calls int
}

// Capabilities returns static capabilities.
func (translator *cacheMockTranslator) Capabilities() Capabilities {
	return translator.capabilities
}

// Translate returns configured result and increments call counter.
func (translator *cacheMockTranslator) Translate(
	_ context.Context,
	_ Request,
) (Result, error) {
	translator.calls++
	return translator.result, nil
}
