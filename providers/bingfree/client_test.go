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
					body := `<html><script>var params_AbusePreventionHelper = [1700000000000,"token123"];</script></html>`
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

	creds, err := parseCredentials(`var params_AbusePreventionHelper = [1700000000000,"token123"];`)
	if err != nil {
		t.Fatalf("parseCredentials error: %v", err)
	}
	if creds.token != "token123" {
		t.Fatalf("token = %q, want token123", creds.token)
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
