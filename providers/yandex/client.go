// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package yandex provides official Yandex Cloud Translate provider.
package yandex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/woozymasta/transitext"
)

const (
	// defaultBaseURL is official Yandex Cloud Translate endpoint.
	defaultBaseURL = "https://translate.api.cloud.yandex.net/translate/v2/translate"

	// defaultTimeout is provider HTTP timeout.
	defaultTimeout = 20 * time.Second

	// defaultBatchMaxChars is conservative Yandex request-char limit.
	defaultBatchMaxChars = 10000
)

// Options controls Yandex provider behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// BaseURL overrides Yandex translate endpoint URL.
	BaseURL string `json:"base_url,omitempty" yaml:"base_url,omitempty" jsonschema:"format=uri,example=https://translate.api.cloud.yandex.net/translate/v2/translate"`

	// API key for Yandex Cloud Translate.
	//nolint:gosec // Runtime credential from external config.
	APIKey string `json:"api_key,omitempty" yaml:"api_key,omitempty" jsonschema:"minLength=1"`

	// IAMToken can be used instead of APIKey.
	//nolint:gosec // Runtime credential from external config.
	IAMToken string `json:"iam_token,omitempty" yaml:"iam_token,omitempty" jsonschema:"minLength=1"`

	// FolderID is Yandex Cloud folder id for API key auth mode.
	FolderID string `json:"folder_id,omitempty" yaml:"folder_id,omitempty" jsonschema:"maxLength=128"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" jsonschema:"minimum=0,default=20000000000"`

	// BatchMaxItems limits request batch size by item count.
	BatchMaxItems int `json:"batch_max_items,omitempty" yaml:"batch_max_items,omitempty" jsonschema:"minimum=1,maximum=1000"`

	// BatchMaxChars limits request batch size by total chars.
	BatchMaxChars int `json:"batch_max_chars,omitempty" yaml:"batch_max_chars,omitempty" jsonschema:"minimum=1,maximum=10000,default=10000"`
}

// Translator is official Yandex Translate provider.
type Translator struct {
	// options stores provider configuration.
	options Options
}

// New creates Yandex provider.
func New(options Options) *Translator {
	if options.BatchMaxChars <= 0 {
		options.BatchMaxChars = defaultBatchMaxChars
	}

	return &Translator{options: options}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.NewCapabilities(
		"yandex",
		transitext.ProviderStable,
		true,
		transitext.CapabilitiesOptions{
			SupportsBatch: true,
			MaxBatchItems: translator.options.BatchMaxItems,
			MaxBatchChars: translator.options.BatchMaxChars,
		},
	)
}

// Translate translates request using Yandex Cloud Translate API.
func (translator *Translator) Translate(
	ctx context.Context,
	request transitext.Request,
) (transitext.Result, error) {
	if strings.TrimSpace(translator.options.APIKey) == "" &&
		strings.TrimSpace(translator.options.IAMToken) == "" {
		return transitext.Result{}, fmt.Errorf(
			"yandex api_key or iam_token is required: %w",
			transitext.ErrInvalidRequest,
		)
	}
	if strings.TrimSpace(translator.options.APIKey) != "" &&
		strings.TrimSpace(translator.options.FolderID) == "" {
		return transitext.Result{}, fmt.Errorf(
			"yandex folder_id is required with api_key: %w",
			transitext.ErrInvalidRequest,
		)
	}

	items, err := transitext.TranslateBatches(
		ctx,
		request,
		translator.Capabilities(),
		translator.translateBatch,
	)
	if err != nil {
		return transitext.Result{}, err
	}

	return transitext.Result{
		Provider: "yandex",
		Items:    items,
	}, nil
}

// translateBatch sends one Yandex API request for one transitext batch.
func (translator *Translator) translateBatch(
	ctx context.Context,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	texts := make([]string, 0, len(request.Items))
	for _, item := range request.Items {
		texts = append(texts, item.Text)
	}

	payload := map[string]any{
		"targetLanguageCode": request.TargetLang,
		"texts":              texts,
	}
	if strings.TrimSpace(request.SourceLang) != "" &&
		!strings.EqualFold(strings.TrimSpace(request.SourceLang), "auto") {
		payload["sourceLanguageCode"] = request.SourceLang
	}
	if strings.TrimSpace(translator.options.APIKey) != "" {
		payload["folderId"] = translator.options.FolderID
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("yandex marshal request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		translator.baseURL(),
		bytes.NewReader(payloadBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("yandex build request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(translator.options.IAMToken) != "" {
		httpRequest.Header.Set("Authorization", "Bearer "+translator.options.IAMToken)
	} else {
		httpRequest.Header.Set("Authorization", "Api-Key "+translator.options.APIKey)
	}

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.httpClient().Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("yandex request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("yandex read response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf(
			"yandex response %d: %w",
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

// parseResponse decodes Yandex response body into transitext items.
func parseResponse(
	input []transitext.Item,
	payload []byte,
) ([]transitext.TranslatedItem, error) {
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return nil, fmt.Errorf("yandex parse response: %w", err)
	}

	translationsValue, ok := decoded["translations"]
	if !ok {
		return nil, fmt.Errorf("yandex response missing translations: %w", transitext.ErrProviderPermanent)
	}
	translations, ok := translationsValue.([]any)
	if !ok {
		return nil, fmt.Errorf("yandex response malformed translations: %w", transitext.ErrProviderPermanent)
	}
	if len(translations) != len(input) {
		return nil, fmt.Errorf(
			"yandex response size mismatch: got %d, want %d: %w",
			len(translations),
			len(input),
			transitext.ErrProviderPermanent,
		)
	}

	out := make([]transitext.TranslatedItem, 0, len(input))
	for index := range input {
		item, ok := translations[index].(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"yandex response malformed translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		text, _ := item["text"].(string)
		if text == "" {
			return nil, fmt.Errorf(
				"yandex response missing text at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		detectedSource, _ := item["detectedLanguageCode"].(string)
		if detectedSource == "" {
			detectedSource, _ = item["detected_language_code"].(string)
		}

		out = append(out, transitext.TranslatedItem{
			ID:             input[index].ID,
			Text:           text,
			DetectedSource: detectedSource,
		})
	}

	return out, nil
}
