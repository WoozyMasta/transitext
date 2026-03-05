// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package yandex

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/woozymasta/transitext"
)

func TestTranslateBatchWithAPIKey(t *testing.T) {
	t.Parallel()

	var calls int
	translator := New(Options{
		APIKey:   "test-key",
		FolderID: "folder-id",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				calls++
				if request.URL.Path != "/translate/v2/translate" {
					t.Fatalf("path = %q, want /translate/v2/translate", request.URL.Path)
				}
				if got := request.Header.Get("Authorization"); got != "Api-Key test-key" {
					t.Fatalf("authorization = %q, want Api-Key test-key", got)
				}

				body := `{"translations":[{"text":"Привет","detectedLanguageCode":"en"},{"text":"Пока","detectedLanguageCode":"en"}]}`
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

func TestTranslateRequiresFolderWithAPIKey(t *testing.T) {
	t.Parallel()

	translator := New(Options{
		APIKey: "test-key",
	})
	_, err := translator.Translate(context.Background(), transitext.Request{
		TargetLang: "ru",
		Items:      []transitext.Item{{ID: "1", Text: "Hello"}},
	})
	if err == nil {
		t.Fatal("Translate error = nil, want folder_id validation error")
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
