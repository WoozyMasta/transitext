// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package microsoft

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/woozymasta/transitext"
)

func TestTranslateEdgeFlow(t *testing.T) {
	t.Parallel()

	var authCalls int
	var translateCalls int

	translator := New(Options{
		AuthURL:      "https://edge.microsoft.com/translate/auth",
		TranslateURL: "https://api-edge.cognitive.microsofttranslator.com/translate",
		HTTPClient: &http.Client{
			Transport: roundTripperFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				switch {
				case strings.Contains(request.URL.Host, "edge.microsoft.com"):
					authCalls++
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("Bearer edge-token")),
						Header:     make(http.Header),
					}, nil
				case strings.Contains(request.URL.Host, "api-edge.cognitive.microsofttranslator.com"):
					translateCalls++
					if got := request.Header.Get("Authorization"); got != "Bearer edge-token" {
						t.Fatalf("authorization = %q, want Bearer edge-token", got)
					}
					body := `[{"detectedLanguage":{"language":"en"},"translations":[{"text":"Привет","to":"ru"}]},
{"detectedLanguage":{"language":"en"},"translations":[{"text":"Пока","to":"ru"}]}]`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header:     make(http.Header),
					}, nil
				default:
					t.Fatalf("unexpected host: %s", request.URL.Host)
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
			{ID: "2", Text: "Bye"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}

	if authCalls != 1 {
		t.Fatalf("authCalls = %d, want 1", authCalls)
	}
	if translateCalls != 1 {
		t.Fatalf("translateCalls = %d, want 1", translateCalls)
	}
	if len(result.Items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(result.Items))
	}
	if result.Items[0].Text != "Привет" || result.Items[1].Text != "Пока" {
		t.Fatalf("translations = %#v", result.Items)
	}
}

func TestTranslateCustomHeadersMode(t *testing.T) {
	t.Parallel()

	var authCalls int
	var translateCalls int

	translator := New(Options{
		Mode:         "custom_headers",
		AuthURL:      "https://edge.microsoft.com/translate/auth",
		TranslateURL: "https://api-edge.cognitive.microsofttranslator.com/translate",
		AuthenticationHeaders: map[string]string{
			"Ocp-Apim-Subscription-Key":    "k1",
			"Ocp-Apim-Subscription-Region": "westeurope",
		},
		TranslateOptions: map[string]string{
			"textType": "html",
		},
		HTTPClient: &http.Client{
			Transport: roundTripperFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				switch {
				case strings.Contains(request.URL.Host, "edge.microsoft.com"):
					authCalls++
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("token")),
						Header:     make(http.Header),
					}, nil
				case strings.Contains(request.URL.Host, "api-edge.cognitive.microsofttranslator.com"):
					translateCalls++
					if got := request.Header.Get("Ocp-Apim-Subscription-Key"); got != "k1" {
						t.Fatalf("subscription key = %q, want k1", got)
					}
					if got := request.Header.Get("Ocp-Apim-Subscription-Region"); got != "westeurope" {
						t.Fatalf("subscription region = %q, want westeurope", got)
					}
					if got := request.Header.Get("Authorization"); got != "" {
						t.Fatalf("authorization = %q, want empty", got)
					}

					query, err := url.ParseQuery(request.URL.RawQuery)
					if err != nil {
						t.Fatalf("parse query: %v", err)
					}
					if got := query.Get("from"); got != "" {
						t.Fatalf("from = %q, want empty for auto source", got)
					}
					if got := query.Get("textType"); got != "html" {
						t.Fatalf("textType = %q, want html", got)
					}

					body := `[{"detectedLanguage":{"language":"en"},"translations":[{"text":"Привет","to":"ru"}]}]`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header:     make(http.Header),
					}, nil
				default:
					t.Fatalf("unexpected host: %s", request.URL.Host)
					return nil, nil
				}
			}),
		},
	})

	result, err := translator.Translate(context.Background(), transitext.Request{
		TargetLang: "ru",
		Items: []transitext.Item{
			{ID: "1", Text: "Hello"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if authCalls != 0 {
		t.Fatalf("authCalls = %d, want 0", authCalls)
	}
	if translateCalls != 1 {
		t.Fatalf("translateCalls = %d, want 1", translateCalls)
	}
	if len(result.Items) != 1 || result.Items[0].Text != "Привет" {
		t.Fatalf("result = %#v", result)
	}
}

type roundTripperFunc func(request *http.Request) (*http.Response, error)

func (function roundTripperFunc) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	return function(request)
}
