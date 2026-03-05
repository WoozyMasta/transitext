// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"strings"
	"testing"
)

func TestContextTranslatorPassThroughWithoutContext(t *testing.T) {
	t.Parallel()

	base := &contextMockTranslator{
		translateFunc: func(_ context.Context, request Request) (Result, error) {
			return Result{
				Provider: "mock",
				Items: []TranslatedItem{
					{ID: request.Items[0].ID, Text: request.Items[0].Text},
				},
			}, nil
		},
	}
	translator := NewContextTranslator(base, ContextOptions{})

	result, err := translator.Translate(context.Background(), Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "1", Text: "замок"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if result.Items[0].Text != "замок" {
		t.Fatalf("text = %q, want замок", result.Items[0].Text)
	}
}

func TestContextTranslatorInjectsAndStrips(t *testing.T) {
	t.Parallel()

	var seenInput string
	base := &contextMockTranslator{
		translateFunc: func(_ context.Context, request Request) (Result, error) {
			seenInput = request.Items[0].Text

			return Result{
				Provider: "mock",
				Items: []TranslatedItem{
					{ID: request.Items[0].ID, Text: "[[|:|]] in the context of a building\n[[|>|]] castle\n[[|#|]]"},
				},
			}, nil
		},
	}
	translator := NewContextTranslator(base, ContextOptions{
		Context: "в контексте здания",
	})

	result, err := translator.Translate(context.Background(), Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "1", Text: "замок"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if !strings.Contains(seenInput, "[[|:|]]") || !strings.Contains(seenInput, "[[|>|]]") {
		t.Fatalf("injected input missing markers: %q", seenInput)
	}
	if !strings.Contains(seenInput, "в контексте здания") {
		t.Fatalf("injected input missing context: %q", seenInput)
	}
	if result.Items[0].Text != "castle" {
		t.Fatalf("text = %q, want castle", result.Items[0].Text)
	}
}

func TestContextTranslatorContextByID(t *testing.T) {
	t.Parallel()

	base := &contextMockTranslator{
		translateFunc: func(_ context.Context, request Request) (Result, error) {
			if !strings.Contains(request.Items[0].Text, "в контексте двери") {
				t.Fatalf("item[0] context mismatch: %q", request.Items[0].Text)
			}
			if !strings.Contains(request.Items[1].Text, "в контексте здания") {
				t.Fatalf("item[1] context mismatch: %q", request.Items[1].Text)
			}

			return Result{
				Provider: "mock",
				Items: []TranslatedItem{
					{ID: "door", Text: "[[|>|]] lock [[|#|]]"},
					{ID: "building", Text: "[[|>|]] castle [[|#|]]"},
				},
			}, nil
		},
	}
	translator := NewContextTranslator(base, ContextOptions{
		Context: "default",
		ContextByID: map[string]string{
			"door":     "в контексте двери",
			"building": "в контексте здания",
		},
	})

	result, err := translator.Translate(context.Background(), Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "door", Text: "замок"},
			{ID: "building", Text: "замок"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if result.Items[0].Text != "lock" || result.Items[1].Text != "castle" {
		t.Fatalf("texts = %#v", result.Items)
	}
}

func TestSanitizeExtractedTextRemovesTrailingMarkerLine(t *testing.T) {
	t.Parallel()

	got := sanitizeExtractedText("Bankverbindung\n[[TRX_ENDE]]")
	if got != "Bankverbindung" {
		t.Fatalf("sanitizeExtractedText = %q, want Bankverbindung", got)
	}
}

type contextMockTranslator struct {
	translateFunc func(ctx context.Context, request Request) (Result, error)
}

func (translator *contextMockTranslator) Translate(
	ctx context.Context,
	request Request,
) (Result, error) {
	return translator.translateFunc(ctx, request)
}

func (*contextMockTranslator) Capabilities() Capabilities {
	return Capabilities{Provider: "mock"}
}
