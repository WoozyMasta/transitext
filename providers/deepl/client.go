// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package deepl provides official DeepL API provider.
package deepl

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
	// defaultPaidURL is DeepL paid API endpoint.
	defaultPaidURL = "https://api.deepl.com/v2/translate"

	// defaultFreeURL is DeepL free API endpoint.
	defaultFreeURL = "https://api-free.deepl.com/v2/translate"

	// defaultTimeout is provider HTTP timeout.
	defaultTimeout = 20 * time.Second
)

// Options controls DeepL provider behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// URL overrides DeepL endpoint URL.
	URL string `json:"url,omitempty" yaml:"url,omitempty"`

	// Auth key for DeepL API.
	//nolint:gosec // Runtime credential from external config.
	AuthKey string `json:"auth_key,omitempty" yaml:"auth_key,omitempty"`

	// SourceLang sets default source language when request.SourceLang is empty.
	SourceLang string `json:"source_lang,omitempty" yaml:"source_lang,omitempty"`

	// Formality controls formality mode where supported.
	Formality string `json:"formality,omitempty" yaml:"formality,omitempty"`

	// SplitSentences controls DeepL sentence splitting behavior.
	SplitSentences string `json:"split_sentences,omitempty" yaml:"split_sentences,omitempty"`

	// BatchMaxItems limits request batch size by item count.
	BatchMaxItems int `json:"batch_max_items,omitempty" yaml:"batch_max_items,omitempty"`

	// BatchMaxChars limits request batch size by total chars.
	BatchMaxChars int `json:"batch_max_chars,omitempty" yaml:"batch_max_chars,omitempty"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// UseFreeAPI selects free endpoint when URL is empty.
	UseFreeAPI bool `json:"use_free_api,omitempty" yaml:"use_free_api,omitempty"`

	// PreserveFormatting preserves source formatting when true.
	PreserveFormatting bool `json:"preserve_formatting,omitempty" yaml:"preserve_formatting,omitempty"`
}

// Translator is official DeepL provider.
type Translator struct {
	// options stores provider configuration.
	options Options
}

// New creates DeepL provider.
func New(options Options) *Translator {
	return &Translator{options: options}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.Capabilities{
		Provider:             "deepl",
		Stability:            transitext.ProviderStable,
		OfficialAPI:          true,
		SupportsGlossary:     false,
		SupportsInstructions: false,
		SupportsBatch:        true,
		SupportsHTML:         false,
		MaxBatchItems:        translator.options.BatchMaxItems,
		MaxBatchChars:        translator.options.BatchMaxChars,
	}
}

// Translate translates request using DeepL API.
func (translator *Translator) Translate(
	ctx context.Context,
	request transitext.Request,
) (transitext.Result, error) {
	if err := transitext.ValidateRequest(request); err != nil {
		return transitext.Result{}, err
	}
	if strings.TrimSpace(translator.options.AuthKey) == "" {
		return transitext.Result{}, fmt.Errorf(
			"deepl auth_key is required: %w",
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
		Provider: "deepl",
		Items:    items,
	}, nil
}

// translateBatch sends one DeepL request for one transitext batch.
func (translator *Translator) translateBatch(
	ctx context.Context,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	texts := make([]string, 0, len(request.Items))
	for _, item := range request.Items {
		texts = append(texts, item.Text)
	}

	sourceLang := strings.TrimSpace(request.SourceLang)
	if sourceLang == "" || strings.EqualFold(sourceLang, "auto") {
		sourceLang = strings.TrimSpace(translator.options.SourceLang)
	}

	payload := map[string]any{
		"text":        texts,
		"target_lang": request.TargetLang,
	}
	if sourceLang != "" {
		payload["source_lang"] = sourceLang
	}
	if strings.TrimSpace(translator.options.Formality) != "" {
		payload["formality"] = strings.TrimSpace(translator.options.Formality)
	}
	if strings.TrimSpace(translator.options.SplitSentences) != "" {
		payload["split_sentences"] = strings.TrimSpace(translator.options.SplitSentences)
	}
	if translator.options.PreserveFormatting {
		payload["preserve_formatting"] = 1
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("deepl marshal request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		translator.endpointURL(),
		bytes.NewReader(payloadBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("deepl build request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "DeepL-Auth-Key "+translator.options.AuthKey)

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.httpClient().Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("deepl request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("deepl read response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf(
			"deepl response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	return parseResponse(request.Items, body)
}

// endpointURL returns configured DeepL endpoint URL.
func (translator *Translator) endpointURL() string {
	if strings.TrimSpace(translator.options.URL) != "" {
		return strings.TrimSpace(translator.options.URL)
	}
	if translator.options.UseFreeAPI {
		return defaultFreeURL
	}

	return defaultPaidURL
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

// parseResponse decodes DeepL response body into transitext items.
func parseResponse(
	input []transitext.Item,
	payload []byte,
) ([]transitext.TranslatedItem, error) {
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return nil, fmt.Errorf("deepl parse response: %w", err)
	}

	translationsValue, ok := decoded["translations"]
	if !ok {
		return nil, fmt.Errorf("deepl response missing translations: %w", transitext.ErrProviderPermanent)
	}
	translations, ok := translationsValue.([]any)
	if !ok {
		return nil, fmt.Errorf("deepl response malformed translations: %w", transitext.ErrProviderPermanent)
	}
	if len(translations) != len(input) {
		return nil, fmt.Errorf(
			"deepl response size mismatch: got %d, want %d: %w",
			len(translations),
			len(input),
			transitext.ErrProviderPermanent,
		)
	}

	out := make([]transitext.TranslatedItem, 0, len(input))
	for index := range input {
		itemMap, ok := translations[index].(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"deepl response malformed translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		text, _ := itemMap["text"].(string)
		if text == "" {
			return nil, fmt.Errorf(
				"deepl response missing text at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		detectedSource, _ := itemMap["detected_source_language"].(string)

		out = append(out, transitext.TranslatedItem{
			ID:             input[index].ID,
			Text:           text,
			DetectedSource: detectedSource,
		})
	}

	return out, nil
}
