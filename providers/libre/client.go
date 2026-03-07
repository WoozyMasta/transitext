// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package libre provides LibreTranslate provider.
package libre

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
	// defaultBaseURL is public LibreTranslate endpoint.
	defaultBaseURL = "https://libretranslate.com/translate"

	// defaultTimeout is provider HTTP timeout.
	defaultTimeout = 20 * time.Second
)

// Options controls LibreTranslate provider behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// BaseURL overrides LibreTranslate endpoint URL.
	BaseURL string `json:"base_url,omitempty" yaml:"base_url,omitempty" jsonschema:"format=uri,example=https://libretranslate.com/translate"`

	// API key for LibreTranslate instance (optional for some deployments).
	//nolint:gosec // Runtime credential from external config.
	APIKey string `json:"api_key,omitempty" yaml:"api_key,omitempty" jsonschema:"minLength=1"`

	// Format controls source text format: "text" or "html".
	Format string `json:"format,omitempty" yaml:"format,omitempty" jsonschema:"enum=text,enum=html,default=text"`

	// BatchMaxItems limits request batch size by item count.
	BatchMaxItems int `json:"batch_max_items,omitempty" yaml:"batch_max_items,omitempty" jsonschema:"minimum=1,maximum=1000"`

	// BatchMaxChars limits request batch size by total chars.
	BatchMaxChars int `json:"batch_max_chars,omitempty" yaml:"batch_max_chars,omitempty" jsonschema:"minimum=1,maximum=1000000"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" jsonschema:"minimum=0,default=20000000000"`
}

// Translator is LibreTranslate provider.
type Translator struct {
	// options stores provider configuration.
	options Options
}

// New creates LibreTranslate provider.
func New(options Options) *Translator {
	return &Translator{options: options}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.NewCapabilities(
		"libre",
		transitext.ProviderStable,
		false,
		transitext.CapabilitiesOptions{
			SupportsBatch: true,
			SupportsHTML:  true,
			MaxBatchItems: translator.options.BatchMaxItems,
			MaxBatchChars: translator.options.BatchMaxChars,
		},
	)
}

// Translate translates request using LibreTranslate API.
func (translator *Translator) Translate(
	ctx context.Context,
	request transitext.Request,
) (transitext.Result, error) {
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
		Provider: "libre",
		Items:    items,
	}, nil
}

// translateBatch translates batch items one-by-one with Libre endpoint.
func (translator *Translator) translateBatch(
	ctx context.Context,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	source := strings.TrimSpace(request.SourceLang)
	if source == "" {
		source = "auto"
	}

	items := make([]transitext.TranslatedItem, 0, len(request.Items))
	for _, item := range request.Items {
		translated, detectedSource, err := translator.translateOne(
			ctx,
			item.Text,
			source,
			request.TargetLang,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, transitext.TranslatedItem{
			ID:             item.ID,
			Text:           translated,
			DetectedSource: detectedSource,
		})
	}

	return items, nil
}

// translateOne translates one text item.
func (translator *Translator) translateOne(
	ctx context.Context,
	text string,
	source string,
	target string,
) (string, string, error) {
	payload := map[string]any{
		"q":      text,
		"source": source,
		"target": target,
		"format": translator.format(),
	}
	if strings.TrimSpace(translator.options.APIKey) != "" {
		payload["api_key"] = translator.options.APIKey
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("libre marshal request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		translator.baseURL(),
		bytes.NewReader(payloadBytes),
	)
	if err != nil {
		return "", "", fmt.Errorf("libre build request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.httpClient().Do(httpRequest)
	if err != nil {
		return "", "", fmt.Errorf("libre request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("libre read response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", "", fmt.Errorf(
			"libre response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		return "", "", fmt.Errorf("libre parse response: %w", err)
	}
	translated, _ := decoded["translatedText"].(string)
	if translated == "" {
		return "", "", fmt.Errorf("libre response missing translatedText: %w", transitext.ErrProviderPermanent)
	}

	detectedSource := ""
	if detectedValue, ok := decoded["detectedLanguage"]; ok {
		if detectedMap, ok := detectedValue.(map[string]any); ok {
			if language, ok := detectedMap["language"].(string); ok {
				detectedSource = language
			}
		}
	}
	if detectedSource == "" {
		if source != "auto" {
			detectedSource = source
		}
	}

	return translated, detectedSource, nil
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
