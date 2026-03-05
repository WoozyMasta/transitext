// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package openai

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/woozymasta/transitext"
)

func TestTranslateSingleBatch(t *testing.T) {
	t.Parallel()

	var calls int
	translator := New(Options{
		AuthToken: "test-key",
		Model:     "test-model",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				calls++
				if request.URL.Path != "/v1/chat/completions" {
					t.Fatalf("path = %q, want /v1/chat/completions", request.URL.Path)
				}
				if got := request.Header.Get("Authorization"); got != "Bearer test-key" {
					t.Fatalf("authorization = %q, want Bearer test-key", got)
				}

				body := `{"choices":[{"message":{"content":"[\"Hello\",\"Exit\"]"}}]}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		},
	})

	result, err := translator.Translate(context.Background(), transitext.Request{
		SourceLang: "fr",
		TargetLang: "en",
		Items: []transitext.Item{
			{ID: "1", Text: "bonjour"},
			{ID: "2", Text: "sortie"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
	if len(result.Items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(result.Items))
	}
	if result.Items[0].Text != "Hello" || result.Items[1].Text != "Exit" {
		t.Fatalf("translations = %#v", result.Items)
	}
}

func TestTranslateBatchSplitByOptions(t *testing.T) {
	t.Parallel()

	var calls int
	translator := New(Options{
		AuthToken:     "test-key",
		BatchMaxItems: 1,
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				calls++
				body := `{"choices":[{"message":{"content":"[\"X\"]"}}]}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		},
	})

	_, err := translator.Translate(context.Background(), transitext.Request{
		TargetLang: "en",
		Items: []transitext.Item{
			{ID: "1", Text: "a"},
			{ID: "2", Text: "b"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
}

func TestParseJSONArrayStrict(t *testing.T) {
	t.Parallel()

	_, err := parseJSONArray(`prefix ["hello"] suffix`, true)
	if err == nil {
		t.Fatal("parseJSONArray strict error = nil, want error")
	}
}

func TestParseJSONArrayTolerant(t *testing.T) {
	t.Parallel()

	values, err := parseJSONArray(`prefix ["hello"] suffix`, false)
	if err != nil {
		t.Fatalf("parseJSONArray tolerant error: %v", err)
	}
	if len(values) != 1 || values[0] != "hello" {
		t.Fatalf("values = %#v, want [hello]", values)
	}
}

// roundTripFunc adapts function into http.RoundTripper.
type roundTripFunc func(request *http.Request) (*http.Response, error)

// RoundTrip executes adapter function.
func (function roundTripFunc) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	return function(request)
}
