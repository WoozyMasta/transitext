// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package deeplfree

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/woozymasta/transitext"
)

func TestTranslateSuccess(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		if request.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", request.Method)
		}
		if got := request.Header.Get("Origin"); got != "https://www.deepl.com" {
			t.Fatalf("Origin header = %q, want https://www.deepl.com", got)
		}
		if got := request.Header.Get("Referer"); got != "https://www.deepl.com/" {
			t.Fatalf("Referer header = %q, want https://www.deepl.com/", got)
		}
		if got := request.Header.Get("Accept-Language"); got != "en-US,en;q=0.9" {
			t.Fatalf("Accept-Language header = %q, want en-US,en;q=0.9", got)
		}
		if got := request.Header.Get("Cookie"); got != "dl_session=test-session" {
			t.Fatalf("Cookie header = %q, want dl_session=test-session", got)
		}

		rawBody, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		if !bytes.Contains(rawBody, []byte(`"method": "LMT_handle_texts"`)) &&
			!bytes.Contains(rawBody, []byte(`"method" : "LMT_handle_texts"`)) {
			t.Fatalf("request body missing patched method spacing: %s", string(rawBody))
		}

		var payload map[string]any
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if payload["method"] != "LMT_handle_texts" {
			t.Fatalf("method payload = %#v, want LMT_handle_texts", payload["method"])
		}

		params, ok := payload["params"].(map[string]any)
		if !ok {
			t.Fatalf("params is not object: %#v", payload["params"])
		}
		texts, ok := params["texts"].([]any)
		if !ok {
			t.Fatalf("texts is not array: %#v", params["texts"])
		}
		if len(texts) != 2 {
			t.Fatalf("texts len = %d, want 2", len(texts))
		}

		first, ok := texts[0].(map[string]any)
		if !ok {
			t.Fatalf("texts[0] is not object: %#v", texts[0])
		}
		if got := int(first["requestAlternatives"].(float64)); got != 5 {
			t.Fatalf("requestAlternatives = %d, want 5", got)
		}

		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"jsonrpc":"2.0",
			"result":{
				"texts":[{"text":"Privet"},{"text":"Mir"}],
				"lang":"ru"
			}
		}`))
	}))
	defer server.Close()

	translator := New(Options{
		URL:                 server.URL,
		RequestAlternatives: 5,
		DLSession:           "test-session",
	})
	result, err := translator.Translate(context.Background(), transitext.Request{
		SourceLang: "auto",
		TargetLang: "en",
		Items: []transitext.Item{
			{ID: "1", Text: "Привет"},
			{ID: "2", Text: "Мир"},
		},
	})
	if err != nil {
		t.Fatalf("Translate error: %v", err)
	}
	if got := result.Provider; got != "deeplfree" {
		t.Fatalf("provider = %q, want deeplfree", got)
	}
	if len(result.Items) != 2 {
		t.Fatalf("items len = %d, want 2", len(result.Items))
	}
	if result.Items[0].ID != "1" || result.Items[0].Text != "Privet" {
		t.Fatalf("item[0] = %#v", result.Items[0])
	}
	if result.Items[1].ID != "2" || result.Items[1].Text != "Mir" {
		t.Fatalf("item[1] = %#v", result.Items[1])
	}
	if result.Items[0].DetectedSource != "RU" {
		t.Fatalf("detected_source = %q, want RU", result.Items[0].DetectedSource)
	}
}

func TestTranslateProviderErrorTemporary(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		_ *http.Request,
	) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"jsonrpc":"2.0",
			"error":{"code":1042902,"message":"too many requests"}
		}`))
	}))
	defer server.Close()

	translator := New(Options{URL: server.URL})
	_, err := translator.Translate(context.Background(), transitext.Request{
		TargetLang: "en",
		Items:      []transitext.Item{{ID: "1", Text: "hello"}},
	})
	if err == nil {
		t.Fatal("Translate error = nil, want temporary provider error")
	}
	if !errors.Is(err, transitext.ErrProviderTemporary) {
		t.Fatalf("error = %v, want ErrProviderTemporary", err)
	}
}

func TestTranslateHTTPStatusError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		_ *http.Request,
	) {
		http.Error(writer, "bad gateway", http.StatusBadGateway)
	}))
	defer server.Close()

	translator := New(Options{URL: server.URL})
	_, err := translator.Translate(context.Background(), transitext.Request{
		TargetLang: "en",
		Items:      []transitext.Item{{ID: "1", Text: "hello"}},
	})
	if err == nil {
		t.Fatal("Translate error = nil, want temporary provider error")
	}
	if !errors.Is(err, transitext.ErrProviderTemporary) {
		t.Fatalf("error = %v, want ErrProviderTemporary", err)
	}
}

func TestPatchMethodSpacing(t *testing.T) {
	t.Parallel()

	base := []byte(`{"jsonrpc":"2.0","method":"LMT_handle_texts"}`)
	patchedA := patchMethodSpacing(24, base)
	if !bytes.Contains(patchedA, []byte(`"method" : "LMT_handle_texts"`)) {
		t.Fatalf("patchedA missing spaced key: %s", string(patchedA))
	}

	patchedB := patchMethodSpacing(1000, base)
	if !bytes.Contains(patchedB, []byte(`"method": "LMT_handle_texts"`)) {
		t.Fatalf("patchedB missing default spacing: %s", string(patchedB))
	}
}

func TestTimestampRounding(t *testing.T) {
	t.Parallel()

	timestamp := deeplTimestamp(3)
	if timestamp <= 0 {
		t.Fatalf("timestamp = %d, want positive", timestamp)
	}
	if timestamp%4 != 0 {
		t.Fatalf("timestamp = %d, want divisible by 4", timestamp)
	}
}
