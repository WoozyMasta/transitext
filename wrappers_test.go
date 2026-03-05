// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryTranslatorSucceedsAfterTemporaryErrors(t *testing.T) {
	t.Parallel()

	mock := &mockTranslator{
		capabilities: Capabilities{Provider: "mock"},
		sequence: []mockResponse{
			{err: ErrProviderTemporary},
			{err: ErrProviderTemporary},
			{result: Result{Provider: "mock", Items: []TranslatedItem{{ID: "1", Text: "ok"}}}},
		},
	}
	retry := NewRetryTranslator(mock, RetryOptions{
		Attempts: 3,
		Delay:    time.Millisecond,
	})

	result, err := retry.Translate(context.Background(), Request{
		TargetLang: "en",
		Items:      []Item{{ID: "1", Text: "bonjour"}},
	})
	if err != nil {
		t.Fatalf("RetryTranslator error: %v", err)
	}
	if mock.calls != 3 {
		t.Fatalf("calls = %d, want 3", mock.calls)
	}
	if len(result.Items) != 1 || result.Items[0].Text != "ok" {
		t.Fatalf("result = %#v", result)
	}
}

func TestRetryTranslatorStopsOnInvalidRequest(t *testing.T) {
	t.Parallel()

	mock := &mockTranslator{
		capabilities: Capabilities{Provider: "mock"},
		sequence: []mockResponse{
			{err: ErrInvalidRequest},
			{result: Result{Provider: "mock"}},
		},
	}
	retry := NewRetryTranslator(mock, RetryOptions{
		Attempts: 3,
		Delay:    time.Millisecond,
	})

	_, err := retry.Translate(context.Background(), Request{
		TargetLang: "en",
		Items:      []Item{{ID: "1", Text: "bonjour"}},
	})
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("error = %v, want ErrInvalidRequest", err)
	}
	if mock.calls != 1 {
		t.Fatalf("calls = %d, want 1", mock.calls)
	}
}

func TestFallbackTranslatorUsesSecondProvider(t *testing.T) {
	t.Parallel()

	first := &mockTranslator{
		capabilities: Capabilities{Provider: "first"},
		sequence: []mockResponse{
			{err: ErrProviderTemporary},
		},
	}
	second := &mockTranslator{
		capabilities: Capabilities{Provider: "second"},
		sequence: []mockResponse{
			{result: Result{Items: []TranslatedItem{{ID: "1", Text: "hello"}}}},
		},
	}

	fallback := NewFallbackTranslator(first, second)
	result, err := fallback.Translate(context.Background(), Request{
		TargetLang: "en",
		Items:      []Item{{ID: "1", Text: "bonjour"}},
	})
	if err != nil {
		t.Fatalf("FallbackTranslator error: %v", err)
	}
	if result.Provider != "second" {
		t.Fatalf("provider = %q, want second", result.Provider)
	}
	if result.Items[0].Text != "hello" {
		t.Fatalf("text = %q, want hello", result.Items[0].Text)
	}
}

func TestFallbackTranslatorReturnsJoinedError(t *testing.T) {
	t.Parallel()

	firstErr := errors.New("first failed")
	secondErr := errors.New("second failed")
	fallback := NewFallbackTranslator(
		&mockTranslator{
			capabilities: Capabilities{Provider: "first"},
			sequence:     []mockResponse{{err: firstErr}},
		},
		&mockTranslator{
			capabilities: Capabilities{Provider: "second"},
			sequence:     []mockResponse{{err: secondErr}},
		},
	)

	_, err := fallback.Translate(context.Background(), Request{
		TargetLang: "en",
		Items:      []Item{{ID: "1", Text: "bonjour"}},
	})
	if err == nil {
		t.Fatal("FallbackTranslator error = nil, want joined error")
	}
	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("joined error does not contain source errors: %v", err)
	}
}

// mockResponse stores one mocked translate call outcome.
type mockResponse struct {
	// err returns translation error when non-nil.
	err error

	// result returns translation result on success.
	result Result
}

// mockTranslator is deterministic translator for wrapper tests.
type mockTranslator struct {
	// sequence defines ordered call outcomes.
	sequence []mockResponse

	// capabilities returns static capabilities.
	capabilities Capabilities

	// calls tracks Translate invocation count.
	calls int
}

// Capabilities returns static mock capabilities.
func (translator *mockTranslator) Capabilities() Capabilities {
	return translator.capabilities
}

// Translate returns predefined response for current call.
func (translator *mockTranslator) Translate(
	_ context.Context,
	_ Request,
) (Result, error) {
	translator.calls++
	if len(translator.sequence) == 0 {
		return Result{}, errors.New("no mock responses left")
	}

	response := translator.sequence[0]
	translator.sequence = translator.sequence[1:]
	if response.err != nil {
		return Result{}, response.err
	}

	return response.result, nil
}
