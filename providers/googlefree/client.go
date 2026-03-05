// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package googlefree provides translation via unofficial Google endpoint.
package googlefree

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/woozymasta/transitext"
)

const (
	// defaultTranslatePath is the py-googletrans translate endpoint path.
	defaultTranslatePath = "/translate_a/single"

	// defaultClientValue is py-googletrans client mode for no-token calls.
	defaultClientValue = "gtx"

	// defaultUserAgent is browser-like UA used by py-googletrans.
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"

	// defaultTimeout is provider HTTP timeout.
	defaultTimeout = 20 * time.Second

	// defaultMaxItems is default request item limit per batch.
	defaultMaxItems = 10

	// defaultMaxChars is default request text-char limit per batch.
	defaultMaxChars = 5000

	// defaultMaxTextChars is default single text-char limit.
	defaultMaxTextChars = 5000
)

var defaultServiceHosts = []string{"translate.googleapis.com"}

var defaultDTParams = []string{
	"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t",
}

// Options controls googlefree translator behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// Request controls low-level HTTP header/cookie/user-agent shaping.
	Request *transitext.HTTPRequestOptions `json:"request,omitempty" yaml:"request,omitempty"`

	// ClientValue overrides "client" query parameter.
	ClientValue string `json:"client_value,omitempty" yaml:"client_value,omitempty"`

	// UserAgent overrides default request user agent.
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`

	// ServiceHosts contains hosts for /translate_a/single endpoint.
	ServiceHosts []string `json:"service_hosts,omitempty" yaml:"service_hosts,omitempty"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// MaxItems limits items per one transitext batch.
	MaxItems int `json:"max_items,omitempty" yaml:"max_items,omitempty"`

	// MaxChars limits total chars per one transitext batch.
	MaxChars int `json:"max_chars,omitempty" yaml:"max_chars,omitempty"`

	// MaxTextChars limits one input text length.
	MaxTextChars int `json:"max_text_chars,omitempty" yaml:"max_text_chars,omitempty"`

	// Concurrency limits parallel per-item HTTP calls.
	Concurrency int `json:"concurrency,omitempty" yaml:"concurrency,omitempty"`
}

// Translator is unofficial Google endpoint translator.
type Translator struct {
	// client performs HTTP calls to provider.
	client *http.Client

	// clientValue stores endpoint "client" query value.
	clientValue string

	// serviceHosts stores translate hosts rotated per request.
	serviceHosts []string

	// maxItems is request item bound.
	maxItems int

	// maxChars is request char bound.
	maxChars int

	// maxTextChars is one text length bound.
	maxTextChars int

	// concurrency limits parallel per-item HTTP calls.
	concurrency int

	// hostCursor stores round-robin host index.
	hostCursor int

	// hostLock protects hostCursor.
	hostLock sync.Mutex
}

// New creates googlefree translator.
func New(options Options) *Translator {
	client := options.HTTPClient
	if client == nil {
		timeout := options.Timeout
		if timeout <= 0 {
			timeout = defaultTimeout
		}
		client = &http.Client{Timeout: timeout}
	}

	serviceHosts := normalizeHosts(options.ServiceHosts)
	clientValue := strings.TrimSpace(options.ClientValue)
	if clientValue == "" {
		clientValue = defaultClientValue
	}

	maxItems := options.MaxItems
	if maxItems <= 0 {
		maxItems = defaultMaxItems
	}

	maxChars := options.MaxChars
	if maxChars <= 0 {
		maxChars = defaultMaxChars
	}

	maxTextChars := options.MaxTextChars
	if maxTextChars <= 0 {
		maxTextChars = defaultMaxTextChars
	}

	concurrency := options.Concurrency
	if concurrency <= 0 {
		concurrency = 2
	}

	requestOptions := transitext.HTTPRequestOptions{}
	if options.Request != nil {
		requestOptions = *options.Request
	}
	if strings.TrimSpace(requestOptions.UserAgent) == "" {
		requestOptions.UserAgent = options.UserAgent
	}
	client = transitext.ConfigureHTTPClient(client, transitext.HTTPRequestDefaults{
		UserAgent: defaultUserAgent,
	}, requestOptions)

	return &Translator{
		client:       client,
		serviceHosts: serviceHosts,
		clientValue:  clientValue,
		maxItems:     maxItems,
		maxChars:     maxChars,
		maxTextChars: maxTextChars,
		concurrency:  concurrency,
	}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.NewCapabilities(
		"googlefree",
		transitext.ProviderUnstable,
		false,
		transitext.CapabilitiesOptions{
			SupportsBatch: true,
			MaxBatchItems: translator.maxItems,
			MaxBatchChars: translator.maxChars,
			MaxTextChars:  translator.maxTextChars,
		},
	)
}

// Translate translates input items with unofficial Google endpoint.
func (translator *Translator) Translate(
	ctx context.Context,
	request transitext.Request,
) (transitext.Result, error) {
	out, err := transitext.TranslateBatches(
		ctx,
		request,
		translator.Capabilities(),
		translator.translateBatch,
	)
	if err != nil {
		return transitext.Result{}, err
	}

	return transitext.Result{
		Provider: "googlefree",
		Items:    out,
	}, nil
}

// translateBatch translates one request batch with limited concurrency.
func (translator *Translator) translateBatch(
	ctx context.Context,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	source := request.SourceLang
	if strings.TrimSpace(source) == "" {
		source = "auto"
	}

	type jobResult struct {
		err   error
		item  transitext.TranslatedItem
		index int
	}
	sem := make(chan struct{}, translator.concurrency)
	results := make(chan jobResult, len(request.Items))

	var group sync.WaitGroup
	for index := range request.Items {
		group.Add(1)

		go func(index int) {
			defer group.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			text, detectedSource, err := translator.translateOne(
				ctx,
				source,
				request.TargetLang,
				request.Items[index].Text,
			)
			if err != nil {
				results <- jobResult{index: index, err: err}
				return
			}

			results <- jobResult{
				index: index,
				item: transitext.TranslatedItem{
					ID:             request.Items[index].ID,
					Text:           text,
					DetectedSource: detectedSource,
				},
			}
		}(index)
	}

	group.Wait()
	close(results)

	out := make([]transitext.TranslatedItem, len(request.Items))
	for result := range results {
		if result.err != nil {
			return nil, result.err
		}

		out[result.index] = result.item
	}

	return out, nil
}

// translateOne translates single text via /translate_a/single.
func (translator *Translator) translateOne(
	ctx context.Context,
	source string,
	target string,
	text string,
) (string, string, error) {
	endpoint := translator.buildEndpoint(source, target, text)

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return "", "", fmt.Errorf("build googlefree request: %w", err)
	}

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.client.Do(httpRequest)
	if err != nil {
		return "", "", fmt.Errorf("googlefree request: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("read googlefree response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", "", fmt.Errorf(
			"googlefree response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	translated, detectedSource, err := parseSingleResponse(body)
	if err != nil {
		return "", "", err
	}

	return translated, detectedSource, nil
}

// buildEndpoint builds query URL for one translation request.
func (translator *Translator) buildEndpoint(
	source string,
	target string,
	text string,
) string {
	host := translator.pickServiceHost()
	urlValue := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   defaultTranslatePath,
	}

	query := urlValue.Query()
	query.Set("client", translator.clientValue)
	query.Set("sl", source)
	query.Set("tl", target)
	query.Set("hl", target)
	query.Set("ie", "UTF-8")
	query.Set("oe", "UTF-8")
	query.Set("otf", "1")
	query.Set("ssel", "0")
	query.Set("tsel", "0")
	query.Set("tk", "xxxx")
	query.Set("q", text)
	for _, dt := range defaultDTParams {
		query.Add("dt", dt)
	}
	urlValue.RawQuery = query.Encode()

	return urlValue.String()
}

// pickServiceHost returns next host in round-robin order.
func (translator *Translator) pickServiceHost() string {
	translator.hostLock.Lock()
	defer translator.hostLock.Unlock()

	host := translator.serviceHosts[translator.hostCursor%len(translator.serviceHosts)]
	translator.hostCursor++

	return host
}

// parseSingleResponse extracts translated text and detected source language.
func parseSingleResponse(payload []byte) (string, string, error) {
	var data []any
	if err := json.Unmarshal(payload, &data); err != nil {
		return "", "", fmt.Errorf("parse googlefree response: %w", err)
	}
	if len(data) == 0 {
		return "", "", fmt.Errorf(
			"googlefree empty response: %w",
			transitext.ErrProviderPermanent,
		)
	}

	parts, ok := data[0].([]any)
	if !ok || len(parts) == 0 {
		return "", "", fmt.Errorf(
			"googlefree malformed translation response: %w",
			transitext.ErrProviderPermanent,
		)
	}

	var builder strings.Builder
	for _, part := range parts {
		segment, ok := part.([]any)
		if !ok || len(segment) == 0 {
			continue
		}

		text, ok := segment[0].(string)
		if !ok {
			continue
		}

		builder.WriteString(text)
	}
	if builder.Len() == 0 {
		return "", "", fmt.Errorf(
			"googlefree empty translation text: %w",
			transitext.ErrProviderPermanent,
		)
	}

	detectedSource := ""
	if len(data) > 2 {
		if source, ok := data[2].(string); ok {
			detectedSource = source
		}
	}

	return builder.String(), detectedSource, nil
}

// normalizeHosts validates and normalizes endpoint hosts list.
func normalizeHosts(hosts []string) []string {
	if len(hosts) == 0 {
		return append([]string(nil), defaultServiceHosts...)
	}

	out := make([]string, 0, len(hosts))
	seen := make(map[string]struct{}, len(hosts))
	for _, host := range hosts {
		normalized := strings.TrimSpace(host)
		normalized = strings.TrimPrefix(normalized, "https://")
		normalized = strings.TrimPrefix(normalized, "http://")
		normalized = strings.Trim(normalized, "/")
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}

		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}

	if len(out) == 0 {
		return append([]string(nil), defaultServiceHosts...)
	}

	return out
}
