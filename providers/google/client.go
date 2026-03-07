// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package google provides official Google Translate API provider.
package google

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
	// defaultBaseURL is official Google Translate v2 endpoint.
	defaultBaseURL = "https://translation.googleapis.com/language/translate/v2"

	// defaultTimeout is provider HTTP timeout.
	defaultTimeout = 20 * time.Second

	// defaultBatchMaxChars is conservative official sync-request char limit.
	defaultBatchMaxChars = 30000
)

// Options controls Google provider behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// BaseURL overrides official Google endpoint URL.
	BaseURL string `json:"base_url,omitempty" yaml:"base_url,omitempty" jsonschema:"format=uri,example=https://translation.googleapis.com/language/translate/v2"`

	// API key for Google Translate API.
	//nolint:gosec // Runtime credential from external config.
	Key string `json:"key,omitempty" yaml:"key,omitempty" jsonschema:"minLength=1"`

	// Format controls source text format: "text" or "html".
	Format string `json:"format,omitempty" yaml:"format,omitempty" jsonschema:"enum=text,enum=html,default=text"`

	// Model selects Google translation model when supported.
	Model string `json:"model,omitempty" yaml:"model,omitempty" jsonschema:"maxLength=64"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" jsonschema:"minimum=0,default=20000000000"`

	// BatchMaxItems limits request batch size by item count.
	BatchMaxItems int `json:"batch_max_items,omitempty" yaml:"batch_max_items,omitempty" jsonschema:"minimum=1,maximum=1000"`

	// BatchMaxChars limits request batch size by total chars.
	BatchMaxChars int `json:"batch_max_chars,omitempty" yaml:"batch_max_chars,omitempty" jsonschema:"minimum=1,maximum=30000,default=30000"`
}

// Translator is official Google Translate provider.
type Translator struct {
	// options stores provider configuration.
	options Options
}

// New creates Google provider.
func New(options Options) *Translator {
	if options.BatchMaxChars <= 0 {
		options.BatchMaxChars = defaultBatchMaxChars
	}

	return &Translator{options: options}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.NewCapabilities(
		"google",
		transitext.ProviderStable,
		true,
		transitext.CapabilitiesOptions{
			SupportsBatch: true,
			SupportsHTML:  true,
			MaxBatchItems: translator.options.BatchMaxItems,
			MaxBatchChars: translator.options.BatchMaxChars,
		},
	)
}

// Translate translates request using official Google API.
func (translator *Translator) Translate(
	ctx context.Context,
	request transitext.Request,
) (transitext.Result, error) {
	if strings.TrimSpace(translator.options.Key) == "" {
		return transitext.Result{}, fmt.Errorf(
			"google api key is required: %w",
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
		Provider: "google",
		Model:    translator.options.Model,
		Items:    items,
	}, nil
}

// translateBatch sends one Google API request for one transitext batch.
func (translator *Translator) translateBatch(
	ctx context.Context,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	input := make([]string, 0, len(request.Items))
	for _, item := range request.Items {
		input = append(input, item.Text)
	}

	payload := map[string]any{
		"q":      input,
		"target": request.TargetLang,
		"format": translator.format(),
	}
	if strings.TrimSpace(request.SourceLang) != "" &&
		!strings.EqualFold(strings.TrimSpace(request.SourceLang), "auto") {
		payload["source"] = request.SourceLang
	}
	if strings.TrimSpace(translator.options.Model) != "" {
		payload["model"] = strings.TrimSpace(translator.options.Model)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("google marshal request: %w", err)
	}

	endpoint, err := url.Parse(translator.baseURL())
	if err != nil {
		return nil, fmt.Errorf("google parse base url: %w", err)
	}
	query := endpoint.Query()
	query.Set("key", translator.options.Key)
	endpoint.RawQuery = query.Encode()

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint.String(),
		bytes.NewReader(payloadBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("google build request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.httpClient().Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("google request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("google read response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf(
			"google response %d: %w",
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

// format returns normalized text format value.
func (translator *Translator) format() string {
	format := strings.ToLower(strings.TrimSpace(translator.options.Format))
	if format != "html" {
		return "text"
	}

	return format
}

// parseResponse decodes Google response body into transitext items.
func parseResponse(
	input []transitext.Item,
	payload []byte,
) ([]transitext.TranslatedItem, error) {
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return nil, fmt.Errorf("google parse response: %w", err)
	}

	dataValue, ok := decoded["data"]
	if !ok {
		return nil, fmt.Errorf("google response missing data: %w", transitext.ErrProviderPermanent)
	}
	dataMap, ok := dataValue.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("google response malformed data: %w", transitext.ErrProviderPermanent)
	}
	translationsValue, ok := dataMap["translations"]
	if !ok {
		return nil, fmt.Errorf("google response missing translations: %w", transitext.ErrProviderPermanent)
	}
	translations, ok := translationsValue.([]any)
	if !ok {
		return nil, fmt.Errorf("google response malformed translations: %w", transitext.ErrProviderPermanent)
	}

	if len(translations) != len(input) {
		return nil, fmt.Errorf(
			"google response size mismatch: got %d, want %d: %w",
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
				"google response malformed translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		translatedText, ok := itemMap["translatedText"].(string)
		if !ok {
			if value, ok := itemMap["translated_text"].(string); ok {
				translatedText = value
			}
		}
		if translatedText == "" {
			return nil, fmt.Errorf(
				"google response missing translated text at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}

		source, _ := itemMap["detectedSourceLanguage"].(string)
		if source == "" {
			source, _ = itemMap["detected_source_language"].(string)
		}

		out = append(out, transitext.TranslatedItem{
			ID:             input[index].ID,
			Text:           translatedText,
			DetectedSource: source,
		})
	}

	return out, nil
}
