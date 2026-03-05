// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package bingfree

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

	var pageCalls int
	var translateCalls int
	translator := New(Options{
		HostURL: "https://www.bing.com",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				switch request.URL.Path {
				case "/translator":
					pageCalls++
					body := `<html data-iid="translator.7777.1"><script>var params_AbusePreventionHelper = [1700000000000,"token123",3600000];</script>IG:"ABC123"</html>`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header:     make(http.Header),
					}, nil
				case "/ttranslatev3":
					translateCalls++
					body := `[{"detectedLanguage":{"language":"en"},"translations":[{"text":"Привет","to":"ru"}]}]`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header:     make(http.Header),
					}, nil
				default:
					t.Fatalf("unexpected path: %s", request.URL.Path)
					return nil, nil
				}
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
	if pageCalls != 1 || translateCalls != 1 {
		t.Fatalf("calls page=%d translate=%d, want 1/1", pageCalls, translateCalls)
	}
	if len(result.Items) != 1 || result.Items[0].Text != "Привет" {
		t.Fatalf("result = %#v", result)
	}
}

func TestParseCredentials(t *testing.T) {
	t.Parallel()

	key, token, ig, iid, expiresAtUnixMilli, err := parseCredentials(
		`<div data-iid="translator.7777.1"></div>IG:"ABC123";var params_AbusePreventionHelper = [1700000000000,"token123",3600000];`,
	)
	if err != nil {
		t.Fatalf("parseCredentials error: %v", err)
	}
	if token != "token123" {
		t.Fatalf("token = %q, want token123", token)
	}
	if key <= 0 {
		t.Fatalf("key = %d, want > 0", key)
	}
	if expiresAtUnixMilli <= key {
		t.Fatalf("expires_at = %d, want > key %d", expiresAtUnixMilli, key)
	}
	if ig != "ABC123" {
		t.Fatalf("ig = %q, want ABC123", ig)
	}
	if iid != "translator.7777.1" {
		t.Fatalf("iid = %q, want translator.7777.1", iid)
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
