// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkValidateRequest(b *testing.B) {
	request := benchRequest(200, 32)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if err := ValidateRequest(request); err != nil {
			b.Fatalf("ValidateRequest error: %v", err)
		}
	}
}

func BenchmarkSplitRequest(b *testing.B) {
	request := benchRequest(500, 48)
	options := BatchOptions{
		MaxItems:   25,
		MaxChars:   1500,
		OnOverflow: OverflowSplit,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		batches, err := SplitRequest(request, options)
		if err != nil {
			b.Fatalf("SplitRequest error: %v", err)
		}
		if len(batches) == 0 {
			b.Fatal("SplitRequest returned empty batches")
		}
	}
}

func BenchmarkCacheTranslatorMissThenHit(b *testing.B) {
	base := benchTranslator{}
	cached := NewCacheTranslator(base, CacheOptions{})
	request := benchRequest(100, 48)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if _, err := cached.Translate(context.Background(), request); err != nil {
			b.Fatalf("CacheTranslator.Translate error: %v", err)
		}
	}
}

func BenchmarkFallbackTranslatorFirstSuccess(b *testing.B) {
	translator := NewFallbackTranslator(
		benchTranslator{},
		benchFailTranslator{},
	)
	request := benchRequest(100, 32)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if _, err := translator.Translate(context.Background(), request); err != nil {
			b.Fatalf("FallbackTranslator.Translate error: %v", err)
		}
	}
}

type benchTranslator struct{}

func (benchTranslator) Translate(_ context.Context, request Request) (Result, error) {
	items := make([]TranslatedItem, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, TranslatedItem{
			ID:   item.ID,
			Text: item.Text,
		})
	}

	return Result{
		Provider: "bench",
		Items:    items,
	}, nil
}

func (benchTranslator) Capabilities() Capabilities {
	return Capabilities{Provider: "bench"}
}

type benchFailTranslator struct{}

func (benchFailTranslator) Translate(_ context.Context, _ Request) (Result, error) {
	return Result{}, ErrProviderTemporary
}

func (benchFailTranslator) Capabilities() Capabilities {
	return Capabilities{Provider: "bench-fail"}
}

func benchRequest(items int, textLen int) Request {
	out := Request{
		SourceLang: "en",
		TargetLang: "ru",
		Items:      make([]Item, 0, items),
	}

	for index := range items {
		out.Items = append(out.Items, Item{
			ID:   fmt.Sprintf("id-%d", index),
			Text: benchText(textLen),
		})
	}

	return out
}

func benchText(length int) string {
	if length <= 0 {
		return "x"
	}

	bytes := make([]byte, 0, length)
	for index := range length {
		bytes = append(bytes, byte('a'+index%26))
	}

	return string(bytes)
}
