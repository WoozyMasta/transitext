// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"net/http"
	"slices"
	"strings"
)

// HTTPRequestOptions controls generic outbound request shaping for providers.
type HTTPRequestOptions struct {
	// Headers sets default headers applied when missing in request.
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" jsonschema:"maxProperties=64"`

	// Cookies sets default cookies applied when request cookie header is empty.
	Cookies map[string]string `json:"cookies,omitempty" yaml:"cookies,omitempty" jsonschema:"maxProperties=64"`

	// UserAgent sets default User-Agent when request does not set it explicitly.
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty" jsonschema:"maxLength=512"`
}

// HTTPRequestDefaults stores provider-owned defaults for request shaping.
type HTTPRequestDefaults struct {
	// Headers are provider default request headers.
	Headers map[string]string

	// Cookies are provider default request cookies.
	Cookies map[string]string

	// UserAgent is provider default User-Agent value.
	UserAgent string
}

// ConfigureHTTPClient applies default and override request options to client.
func ConfigureHTTPClient(
	client *http.Client,
	defaults HTTPRequestDefaults,
	options HTTPRequestOptions,
) *http.Client {
	if client == nil {
		return nil
	}
	if client.Transport == nil {
		client.Transport = http.DefaultTransport
	}

	userAgent := strings.TrimSpace(defaults.UserAgent)
	if value := strings.TrimSpace(options.UserAgent); value != "" {
		userAgent = value
	}

	headers := mergeStringMap(defaults.Headers, options.Headers)
	cookies := mergeStringMap(defaults.Cookies, options.Cookies)
	cookieHeader := buildCookieHeader(cookies)

	if userAgent == "" && len(headers) == 0 && cookieHeader == "" {
		return client
	}

	client.Transport = requestProfileRoundTripper{
		base:         client.Transport,
		userAgent:    userAgent,
		headers:      headers,
		cookieHeader: cookieHeader,
	}

	return client
}

// requestProfileRoundTripper applies default headers to outgoing requests.
type requestProfileRoundTripper struct {
	// base executes the real network request.
	base http.RoundTripper

	// userAgent stores default User-Agent.
	userAgent string

	// headers stores default headers.
	headers map[string]string

	// cookieHeader stores serialized default cookie header.
	cookieHeader string
}

// RoundTrip clones request and injects default headers when missing.
func (roundTripper requestProfileRoundTripper) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	cloned := request.Clone(request.Context())
	cloned.Header = request.Header.Clone()

	if roundTripper.userAgent != "" && cloned.Header.Get("User-Agent") == "" {
		cloned.Header.Set("User-Agent", roundTripper.userAgent)
	}

	for key, value := range roundTripper.headers {
		if cloned.Header.Get(key) == "" {
			cloned.Header.Set(key, value)
		}
	}

	if roundTripper.cookieHeader != "" && cloned.Header.Get("Cookie") == "" {
		cloned.Header.Set("Cookie", roundTripper.cookieHeader)
	}

	return roundTripper.base.RoundTrip(cloned)
}

// mergeStringMap merges defaults with overrides and drops empty values.
func mergeStringMap(
	defaults map[string]string,
	overrides map[string]string,
) map[string]string {
	if len(defaults) == 0 && len(overrides) == 0 {
		return nil
	}

	out := make(map[string]string, len(defaults)+len(overrides))
	for key, value := range defaults {
		trimmedKey := strings.TrimSpace(key)
		trimmedValue := strings.TrimSpace(value)
		if trimmedKey == "" || trimmedValue == "" {
			continue
		}

		out[trimmedKey] = trimmedValue
	}

	for key, value := range overrides {
		trimmedKey := strings.TrimSpace(key)
		trimmedValue := strings.TrimSpace(value)
		if trimmedKey == "" || trimmedValue == "" {
			continue
		}

		out[trimmedKey] = trimmedValue
	}

	if len(out) == 0 {
		return nil
	}

	return out
}

// buildCookieHeader serializes cookies map into deterministic header string.
func buildCookieHeader(cookies map[string]string) string {
	if len(cookies) == 0 {
		return ""
	}

	keys := make([]string, 0, len(cookies))
	for key := range cookies {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	values := make([]string, 0, len(keys))
	for _, key := range keys {
		trimmedKey := strings.TrimSpace(key)
		trimmedValue := strings.TrimSpace(cookies[key])
		if trimmedKey == "" || trimmedValue == "" {
			continue
		}

		values = append(values, trimmedKey+"="+trimmedValue)
	}

	if len(values) == 0 {
		return ""
	}

	return strings.Join(values, "; ")
}
