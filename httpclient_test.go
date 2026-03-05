// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestConfigureHTTPClientAppliesDefaultsAndOverrides(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: roundTripFunc(func(
			request *http.Request,
		) (*http.Response, error) {
			if got := request.Header.Get("User-Agent"); got != "override-ua" {
				t.Fatalf("User-Agent = %q, want override-ua", got)
			}
			if got := request.Header.Get("X-Default"); got != "a" {
				t.Fatalf("X-Default = %q, want a", got)
			}
			if got := request.Header.Get("X-Override"); got != "b" {
				t.Fatalf("X-Override = %q, want b", got)
			}
			if got := request.Header.Get("Cookie"); got != "a=1; z=9" {
				t.Fatalf("Cookie = %q, want a=1; z=9", got)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("ok")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	configured := ConfigureHTTPClient(client, HTTPRequestDefaults{
		UserAgent: "default-ua",
		Headers: map[string]string{
			"X-Default": "a",
		},
		Cookies: map[string]string{
			"z": "9",
		},
	}, HTTPRequestOptions{
		UserAgent: "override-ua",
		Headers: map[string]string{
			"X-Override": "b",
		},
		Cookies: map[string]string{
			"a": "1",
		},
	})

	request, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}
	response, err := configured.Do(request)
	if err != nil {
		t.Fatalf("Do error: %v", err)
	}
	_ = response.Body.Close()
}

func TestConfigureHTTPClientDoesNotOverrideExplicitRequestHeaders(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: roundTripFunc(func(
			request *http.Request,
		) (*http.Response, error) {
			if got := request.Header.Get("User-Agent"); got != "explicit-ua" {
				t.Fatalf("User-Agent = %q, want explicit-ua", got)
			}
			if got := request.Header.Get("X-Test"); got != "explicit" {
				t.Fatalf("X-Test = %q, want explicit", got)
			}
			if got := request.Header.Get("Cookie"); got != "k=v" {
				t.Fatalf("Cookie = %q, want k=v", got)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("ok")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	configured := ConfigureHTTPClient(client, HTTPRequestDefaults{
		UserAgent: "default-ua",
		Headers: map[string]string{
			"X-Test": "default",
		},
		Cookies: map[string]string{
			"a": "1",
		},
	}, HTTPRequestOptions{})

	request, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}
	request.Header.Set("User-Agent", "explicit-ua")
	request.Header.Set("X-Test", "explicit")
	request.Header.Set("Cookie", "k=v")

	response, err := configured.Do(request)
	if err != nil {
		t.Fatalf("Do error: %v", err)
	}
	_ = response.Body.Close()
}

type roundTripFunc func(request *http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	return function(request)
}
