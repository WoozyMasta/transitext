// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package libre

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/woozymasta/transitext"
)

func TestTranslateBatch(t *testing.T) {
	t.Parallel()

	var calls int
	translator := New(Options{
		BaseURL: "https://libretranslate.com/translate",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				calls++
				if request.URL.Path != "/translate" {
					t.Fatalf("path = %q, want /translate", request.URL.Path)
				}
				body := `{"translatedText":"Привет","detectedLanguage":{"language":"en"}}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		},
	})

	result, err := translator.Translate(context.Background(), transitext.Request{
		SourceLang: "en",
		TargetLang: "ru",
		Items: []transitext.Item{
			{ID: "1", Text: "Hello"},
			{ID: "2", Text: "Hello 2"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2 (per-item libre requests)", calls)
	}
	if len(result.Items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(result.Items))
	}
	if result.Items[0].Text != "Привет" || result.Items[1].Text != "Привет" {
		t.Fatalf("translations = %#v", result.Items)
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
