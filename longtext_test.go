// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestLongTextTranslatorSplitAndMerge(t *testing.T) {
	t.Parallel()

	base := &longTextMockTranslator{
		capabilities: Capabilities{
			Provider:     "mock",
			MaxTextChars: 32,
		},
	}
	translator := NewLongTextTranslator(base, LongTextOptions{})

	request := Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "a", Text: "alpha.\n\nbeta sentence. gamma sentence."},
			{ID: "b", Text: "short"},
		},
	}

	result, err := translator.Translate(context.Background(), request)
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if len(base.lastRequest.Items) <= len(request.Items) {
		t.Fatalf(
			"expanded request items = %d, want > %d",
			len(base.lastRequest.Items),
			len(request.Items),
		)
	}
	if len(result.Items) != len(request.Items) {
		t.Fatalf("result items len = %d, want %d", len(result.Items), len(request.Items))
	}
	if got, want := result.Items[0].ID, "a"; got != want {
		t.Fatalf("item[0].id = %q, want %q", got, want)
	}
	if got, want := result.Items[0].Text, strings.ToUpper(request.Items[0].Text); got != want {
		t.Fatalf("item[0].text = %q, want %q", got, want)
	}
	if got, want := result.Items[1].Text, strings.ToUpper(request.Items[1].Text); got != want {
		t.Fatalf("item[1].text = %q, want %q", got, want)
	}
}

func TestLongTextTranslatorErrorOnOverflow(t *testing.T) {
	t.Parallel()

	base := &longTextMockTranslator{
		capabilities: Capabilities{
			Provider:     "mock",
			MaxTextChars: 4,
		},
	}
	translator := NewLongTextTranslator(base, LongTextOptions{
		ErrorOnOverflow: true,
	})

	_, err := translator.Translate(context.Background(), Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "a", Text: "too long"},
		},
	})
	if err == nil {
		t.Fatal("Translate error = nil, want ErrTextTooLong")
	}
	if !errors.Is(err, ErrTextTooLong) {
		t.Fatalf("error = %v, want ErrTextTooLong", err)
	}
}

func TestSplitLongTextParagraphPriority(t *testing.T) {
	t.Parallel()

	text := "line 1.\n\nline 2.\n\nline 3."
	parts := splitLongText(text, 12)

	if len(parts) < 3 {
		t.Fatalf("parts len = %d, want >= 3", len(parts))
	}
	for index, part := range parts[:2] {
		if !strings.Contains(part, "\n\n") {
			t.Fatalf("parts[%d] = %q, want paragraph separator preserved", index, part)
		}
	}
}

// longTextMockTranslator is deterministic translator for long text tests.
type longTextMockTranslator struct {
	// capabilities stores static provider capabilities.
	capabilities Capabilities

	// lastRequest keeps request for assertions.
	lastRequest Request
}

// Capabilities returns static mocked capabilities.
func (translator *longTextMockTranslator) Capabilities() Capabilities {
	return translator.capabilities
}

// Translate stores request and uppercases every input text.
func (translator *longTextMockTranslator) Translate(
	_ context.Context,
	request Request,
) (Result, error) {
	translator.lastRequest = request

	items := make([]TranslatedItem, len(request.Items))
	for index, item := range request.Items {
		items[index] = TranslatedItem{
			ID:   item.ID,
			Text: strings.ToUpper(item.Text),
		}
	}

	return Result{
		Provider: "mock",
		Items:    items,
	}, nil
}
