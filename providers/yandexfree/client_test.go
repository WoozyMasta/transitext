// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package yandexfree

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/woozymasta/transitext"
)

func TestTranslateFlow(t *testing.T) {
	t.Parallel()

	var calls int
	translator := New(Options{
		BaseURL: "https://translate.yandex.net/api/v1/tr.json",
		HTTPClient: &http.Client{
			Transport: roundTripperFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				calls++
				if request.URL.Path != "/api/v1/tr.json/translate" {
					t.Fatalf("path = %q, want /api/v1/tr.json/translate", request.URL.Path)
				}
				body := `{"code":200,"lang":"en-ru","text":["Привет"]}`
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
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
	if len(result.Items) != 1 || result.Items[0].Text != "Привет" {
		t.Fatalf("result = %#v", result)
	}
	if result.Items[0].DetectedSource != "en" {
		t.Fatalf("detected source = %q, want en", result.Items[0].DetectedSource)
	}
}

type roundTripperFunc func(request *http.Request) (*http.Response, error)

func (function roundTripperFunc) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	return function(request)
}
