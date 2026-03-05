// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package azure provides official Azure Translator provider.
package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/woozymasta/transitext"
)

const (
	// defaultBaseURL is official Azure Translator endpoint.
	defaultBaseURL = "https://api.cognitive.microsofttranslator.com/translate"

	// defaultTimeout is provider HTTP timeout.
	defaultTimeout = 20 * time.Second
)

// Options controls Azure provider behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// BaseURL overrides Azure translate endpoint URL.
	BaseURL string `json:"base_url,omitempty" yaml:"base_url,omitempty"`

	// Key is Azure Translator subscription key.
	//nolint:gosec // Runtime credential from external config.
	Key string `json:"key,omitempty" yaml:"key,omitempty"`

	// Region is Azure Translator resource region.
	Region string `json:"region,omitempty" yaml:"region,omitempty"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// BatchMaxItems limits request batch size by item count.
	BatchMaxItems int `json:"batch_max_items,omitempty" yaml:"batch_max_items,omitempty"`

	// BatchMaxChars limits request batch size by total chars.
	BatchMaxChars int `json:"batch_max_chars,omitempty" yaml:"batch_max_chars,omitempty"`
}

// Translator is official Azure Translator provider.
type Translator struct {
	// options stores provider configuration.
	options Options
}

// New creates Azure provider.
func New(options Options) *Translator {
	return &Translator{options: options}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.Capabilities{
		Provider:             "azure",
		Stability:            transitext.ProviderStable,
		OfficialAPI:          true,
		SupportsGlossary:     false,
		SupportsInstructions: false,
		SupportsBatch:        true,
		SupportsHTML:         true,
		MaxBatchItems:        translator.options.BatchMaxItems,
		MaxBatchChars:        translator.options.BatchMaxChars,
	}
}

// Translate translates request using Azure Translator API.
func (translator *Translator) Translate(
	ctx context.Context,
	request transitext.Request,
) (transitext.Result, error) {
	if err := transitext.ValidateRequest(request); err != nil {
		return transitext.Result{}, err
	}
	if strings.TrimSpace(translator.options.Key) == "" {
		return transitext.Result{}, fmt.Errorf(
			"azure key is required: %w",
			transitext.ErrInvalidRequest,
		)
	}

	batchOptions := request.Batch
	if batchOptions.MaxItems <= 0 && translator.options.BatchMaxItems > 0 {
		batchOptions.MaxItems = translator.options.BatchMaxItems
	}
	if batchOptions.MaxChars <= 0 && translator.options.BatchMaxChars > 0 {
		batchOptions.MaxChars = translator.options.BatchMaxChars
	}
	batches, err := transitext.SplitRequest(request, batchOptions)
	if err != nil {
		return transitext.Result{}, err
	}

	items := make([]transitext.TranslatedItem, 0, len(request.Items))
	for _, batch := range batches {
		batchItems, err := translator.translateBatch(ctx, batch)
		if err != nil {
			return transitext.Result{}, err
		}

		items = append(items, batchItems...)
	}

	return transitext.Result{
		Provider: "azure",
		Items:    items,
	}, nil
}

// translateBatch sends one Azure API request for one transitext batch.
func (translator *Translator) translateBatch(
	ctx context.Context,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	payload := make([]map[string]string, 0, len(request.Items))
	for _, item := range request.Items {
		payload = append(payload, map[string]string{"Text": item.Text})
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("azure marshal request: %w", err)
	}

	endpoint, err := url.Parse(translator.baseURL())
	if err != nil {
		return nil, fmt.Errorf("azure parse base url: %w", err)
	}
	query := endpoint.Query()
	query.Set("api-version", "3.0")
	query.Set("to", request.TargetLang)
	if strings.TrimSpace(request.SourceLang) != "" &&
		!strings.EqualFold(strings.TrimSpace(request.SourceLang), "auto") {
		query.Set("from", request.SourceLang)
	}
	endpoint.RawQuery = query.Encode()

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint.String(),
		bytes.NewReader(payloadBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("azure build request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Ocp-Apim-Subscription-Key", translator.options.Key)
	if strings.TrimSpace(translator.options.Region) != "" {
		httpRequest.Header.Set("Ocp-Apim-Subscription-Region", translator.options.Region)
	}

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.httpClient().Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("azure request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("azure read response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf(
			"azure response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	return parseResponse(request.Items, responseBody)
}

// httpClient returns configured HTTP client.
func (translator *Translator) httpClient() *http.Client {
	if translator.options.HTTPClient != nil {
		return translator.options.HTTPClient
	}

	timeout := translator.options.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	return &http.Client{Timeout: timeout}
}

// baseURL returns normalized API URL.
func (translator *Translator) baseURL() string {
	baseURL := strings.TrimSpace(translator.options.BaseURL)
	if baseURL == "" {
		return defaultBaseURL
	}

	return baseURL
}

// parseResponse decodes Azure response body into transitext items.
func parseResponse(
	input []transitext.Item,
	payload []byte,
) ([]transitext.TranslatedItem, error) {
	var decoded []map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return nil, fmt.Errorf("azure parse response: %w", err)
	}
	if len(decoded) != len(input) {
		return nil, fmt.Errorf(
			"azure response size mismatch: got %d, want %d: %w",
			len(decoded),
			len(input),
			transitext.ErrProviderPermanent,
		)
	}

	out := make([]transitext.TranslatedItem, 0, len(input))
	for index := range input {
		item := decoded[index]
		translationsValue, ok := item["translations"]
		if !ok {
			return nil, fmt.Errorf(
				"azure response missing translations at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		translations, ok := translationsValue.([]any)
		if !ok || len(translations) == 0 {
			return nil, fmt.Errorf(
				"azure response missing translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		first, ok := translations[0].(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"azure response malformed translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		text, _ := first["text"].(string)
		if text == "" {
			return nil, fmt.Errorf(
				"azure response empty text at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}

		detectedSource := ""
		if detectedValue, ok := item["detectedLanguage"]; ok {
			if detectedMap, ok := detectedValue.(map[string]any); ok {
				if language, ok := detectedMap["language"].(string); ok {
					detectedSource = language
				}
			}
		}

		out = append(out, transitext.TranslatedItem{
			ID:             input[index].ID,
			Text:           text,
			DetectedSource: detectedSource,
		})
	}

	return out, nil
}
