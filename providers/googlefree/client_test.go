// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package googlefree

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/woozymasta/transitext"
)

func TestTranslateSingleEndpoint(t *testing.T) {
	t.Parallel()

	var requests int32
	translator := New(Options{
		HTTPClient: &http.Client{
			Transport: roundTripperFunc(func(
				request *http.Request,
			) (*http.Response, error) {
				atomic.AddInt32(&requests, 1)
				if request.URL.Path != "/translate_a/single" {
					t.Fatalf("path = %q, want /translate_a/single", request.URL.Path)
				}
				if got := request.URL.Query().Get("client"); got != "gtx" {
					t.Fatalf("client = %q, want gtx", got)
				}
				if request.URL.Query().Get("tk") == "" {
					t.Fatal("tk query parameter should exist")
				}
				if got := request.Header.Get("User-Agent"); got == "" {
					t.Fatal("user-agent should be set")
				}

				q := request.URL.Query().Get("q")
				body := `[[["` + strings.ToUpper(q) + `","` + q + `",null,null,10]],null,"fr"]`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		},
		MaxItems:    10,
		MaxChars:    1000,
		Concurrency: 2,
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
	if len(result.Items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(result.Items))
	}
	if got := atomic.LoadInt32(&requests); got != 2 {
		t.Fatalf("requests = %d, want 2", got)
	}
	if result.Items[0].Text != "BONJOUR" || result.Items[1].Text != "SORTIE" {
		t.Fatalf("translations = %#v", result.Items)
	}
}

func TestParseSingleResponse(t *testing.T) {
	t.Parallel()

	payload := []byte(`[[["Hello","bonjour",null,null,10],[" world"," monde",null,null,10]],null,"fr"]`)
	translated, source, err := parseSingleResponse(payload)
	if err != nil {
		t.Fatalf("parseSingleResponse error: %v", err)
	}
	if translated != "Hello world" {
		t.Fatalf("translated = %q, want %q", translated, "Hello world")
	}
	if source != "fr" {
		t.Fatalf("source = %q, want %q", source, "fr")
	}
}

type roundTripperFunc func(request *http.Request) (*http.Response, error)

func (function roundTripperFunc) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	return function(request)
}
