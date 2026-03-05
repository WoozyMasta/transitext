// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package google

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
		Key: "test-key",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				calls++
				if request.URL.Path != "/language/translate/v2" {
					t.Fatalf("path = %q, want /language/translate/v2", request.URL.Path)
				}
				if got := request.URL.Query().Get("key"); got != "test-key" {
					t.Fatalf("key query = %q, want test-key", got)
				}
				body := `{"data":{"translations":[{"translatedText":"Hello","detectedSourceLanguage":"fr"},{"translatedText":"Exit","detectedSourceLanguage":"fr"}]}}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		},
	})

	result, err := translator.Translate(context.Background(), transitext.Request{
		SourceLang: "auto",
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
	if result.Items[0].DetectedSource != "fr" {
		t.Fatalf("detected_source = %q, want fr", result.Items[0].DetectedSource)
	}
}

func TestTranslateSplitByBatchMaxItems(t *testing.T) {
	t.Parallel()

	var calls int
	translator := New(Options{
		Key:           "test-key",
		BatchMaxItems: 1,
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				calls++
				body := `{"data":{"translations":[{"translatedText":"X","detectedSourceLanguage":"en"}]}}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		},
	})

	_, err := translator.Translate(context.Background(), transitext.Request{
		TargetLang: "ru",
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

// roundTripFunc adapts function into http.RoundTripper.
type roundTripFunc func(request *http.Request) (*http.Response, error)

// RoundTrip executes adapter function.
func (function roundTripFunc) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	return function(request)
}
