// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"errors"
	"testing"
)

func TestPartialTranslatorRetriesOneItemAndContinues(t *testing.T) {
	t.Parallel()

	base := &partialMockTranslator{
		responses: map[string][]partialMockResponse{
			"a": {
				{err: ErrProviderPermanent},
				{result: Result{Items: []TranslatedItem{{Text: "A"}}}},
			},
			"b": {
				{result: Result{Items: []TranslatedItem{{Text: "B"}}}},
			},
		},
	}
	translator := NewPartialTranslator(base, PartialOptions{
		ItemRetries: 1,
	})

	result, err := translator.Translate(context.Background(), Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "a", Text: "a"},
			{ID: "b", Text: "b"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if len(result.Items) != 2 {
		t.Fatalf("items len = %d, want 2", len(result.Items))
	}
	if result.Items[0].Text != "A" || result.Items[1].Text != "B" {
		t.Fatalf("unexpected output: %#v", result.Items)
	}
}

func TestPartialTranslatorStopsOnTemporaryAndMarksRest(t *testing.T) {
	t.Parallel()

	base := &partialMockTranslator{
		responses: map[string][]partialMockResponse{
			"a": {
				{result: Result{Items: []TranslatedItem{{Text: "A"}}}},
			},
			"b": {
				{err: ErrProviderTemporary},
			},
		},
	}
	translator := NewPartialTranslator(base, PartialOptions{})

	result, err := translator.Translate(context.Background(), Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "a", Text: "a"},
			{ID: "b", Text: "b"},
			{ID: "c", Text: "c"},
		},
	})
	if err == nil {
		t.Fatal("Translate error = nil, want stop error")
	}
	if !errors.Is(err, ErrProviderTemporary) {
		t.Fatalf("error = %v, want ErrProviderTemporary", err)
	}
	if len(result.Items) != 3 {
		t.Fatalf("items len = %d, want 3", len(result.Items))
	}
	if result.Items[0].Text != "A" || result.Items[0].Error != nil {
		t.Fatalf("item[0] = %#v", result.Items[0])
	}
	if result.Items[1].Error == nil || result.Items[1].Error.Code != "provider_temporary" {
		t.Fatalf("item[1] error = %#v", result.Items[1].Error)
	}
	if result.Items[2].Error == nil || result.Items[2].Error.Code != "not_processed" {
		t.Fatalf("item[2] error = %#v", result.Items[2].Error)
	}
}

func TestPartialTranslatorContinuesOnTemporaryWhenEnabled(t *testing.T) {
	t.Parallel()

	base := &partialMockTranslator{
		responses: map[string][]partialMockResponse{
			"a": {{err: ErrProviderTemporary}},
			"b": {{result: Result{Items: []TranslatedItem{{Text: "B"}}}}},
		},
	}
	translator := NewPartialTranslator(base, PartialOptions{
		ContinueOnTemporary: true,
	})

	result, err := translator.Translate(context.Background(), Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "a", Text: "a"},
			{ID: "b", Text: "b"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if result.Items[0].Error == nil || result.Items[0].Error.Code != "provider_temporary" {
		t.Fatalf("item[0] error = %#v", result.Items[0].Error)
	}
	if result.Items[1].Text != "B" {
		t.Fatalf("item[1] = %#v", result.Items[1])
	}
}

// partialMockResponse stores one fake translation call outcome.
type partialMockResponse struct {
	// result is successful response payload.
	result Result

	// err is failed response payload.
	err error
}

// partialMockTranslator is deterministic translator for partial tests.
type partialMockTranslator struct {
	// responses stores queued outcomes per input text.
	responses map[string][]partialMockResponse
}

// Capabilities returns static provider id.
func (*partialMockTranslator) Capabilities() Capabilities {
	return Capabilities{Provider: "mock"}
}

// Translate consumes one queued outcome based on first item text.
func (translator *partialMockTranslator) Translate(
	_ context.Context,
	request Request,
) (Result, error) {
	key := request.Items[0].Text
	queue := translator.responses[key]
	if len(queue) == 0 {
		return Result{}, errors.New("missing mock response")
	}

	response := queue[0]
	translator.responses[key] = queue[1:]
	if response.err != nil {
		return Result{}, response.err
	}

	item := response.result.Items[0]
	if item.ID == "" {
		item.ID = request.Items[0].ID
	}

	return Result{
		Provider: "mock",
		Items:    []TranslatedItem{item},
	}, nil
}
