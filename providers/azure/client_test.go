// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package azure

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
		Key:    "test-key",
		Region: "westeurope",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				calls++
				if request.URL.Path != "/translate" {
					t.Fatalf("path = %q, want /translate", request.URL.Path)
				}
				if request.URL.Query().Get("api-version") != "3.0" {
					t.Fatalf("api-version query missing")
				}
				if got := request.Header.Get("Ocp-Apim-Subscription-Key"); got != "test-key" {
					t.Fatalf("subscription key = %q, want test-key", got)
				}
				if got := request.Header.Get("Ocp-Apim-Subscription-Region"); got != "westeurope" {
					t.Fatalf("region = %q, want westeurope", got)
				}

				body := `[{"detectedLanguage":{"language":"en"},"translations":[{"text":"Привет","to":"ru"}]},
{"detectedLanguage":{"language":"en"},"translations":[{"text":"Пока","to":"ru"}]}]`
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
			{ID: "2", Text: "Bye"},
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
	if result.Items[0].Text != "Привет" || result.Items[1].Text != "Пока" {
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
