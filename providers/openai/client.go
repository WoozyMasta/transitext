// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package openai provides OpenAI-compatible translation provider.
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/woozymasta/transitext"
)

const (
	// defaultBaseURL is default OpenAI API base URL.
	defaultBaseURL = "https://api.openai.com/v1"

	// defaultModel is default model when not configured.
	defaultModel = "gpt-4o-mini"

	// defaultTimeout is default HTTP request timeout.
	defaultTimeout = 60 * time.Second
)

const (
	promptPreserve = "Preserve punctuation, spacing, and placeholders like {name}, %%s, {0}, <tag>."
	promptJSONOnly = "Return ONLY a JSON array of strings in the same order."
)

// Options controls OpenAI-compatible provider behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// BaseURL is OpenAI-compatible API base URL.
	BaseURL string `json:"base_url,omitempty" yaml:"base_url,omitempty"`

	// AuthToken is bearer token for API authentication.
	//nolint:gosec // This is runtime config field, not embedded secret.
	AuthToken string `json:"auth_token,omitempty" yaml:"auth_token,omitempty"`

	// Model is model identifier.
	Model string `json:"model,omitempty" yaml:"model,omitempty"`

	// SystemPrompt overrides default system prompt.
	SystemPrompt string `json:"system_prompt,omitempty" yaml:"system_prompt,omitempty"`

	// InstructionPrefix is prepended to request instructions.
	InstructionPrefix string `json:"instruction_prefix,omitempty" yaml:"instruction_prefix,omitempty"`

	// InstructionSuffix is appended to request instructions.
	InstructionSuffix string `json:"instruction_suffix,omitempty" yaml:"instruction_suffix,omitempty"`

	// Temperature sets sampling temperature.
	Temperature float64 `json:"temperature,omitempty" yaml:"temperature,omitempty"`

	// TopP sets nucleus sampling parameter.
	TopP float64 `json:"top_p,omitempty" yaml:"top_p,omitempty"`

	// MaxTokens sets response token cap when supported.
	MaxTokens int `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`

	// StrictJSONArray requires strict JSON array response parsing.
	StrictJSONArray bool `json:"strict_json_array,omitempty" yaml:"strict_json_array,omitempty"`

	// BatchMaxItems limits request batch size by items.
	BatchMaxItems int `json:"batch_max_items,omitempty" yaml:"batch_max_items,omitempty"`

	// BatchMaxChars limits request batch size by chars.
	BatchMaxChars int `json:"batch_max_chars,omitempty" yaml:"batch_max_chars,omitempty"`

	// Timeout is HTTP timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// Translator is OpenAI-compatible provider implementation.
type Translator struct {
	// options stores provider behavior config.
	options Options
}

// New creates OpenAI-compatible translator.
func New(options Options) *Translator {
	return &Translator{options: options}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.Capabilities{
		Provider:             "openai",
		Stability:            transitext.ProviderStable,
		OfficialAPI:          false,
		SupportsGlossary:     false,
		SupportsInstructions: true,
		SupportsBatch:        true,
		SupportsHTML:         false,
		MaxBatchItems:        translator.options.BatchMaxItems,
		MaxBatchChars:        translator.options.BatchMaxChars,
	}
}

// Translate translates request with chat completion API.
func (translator *Translator) Translate(
	ctx context.Context,
	request transitext.Request,
) (transitext.Result, error) {
	if err := transitext.ValidateRequest(request); err != nil {
		return transitext.Result{}, err
	}
	if strings.TrimSpace(translator.options.AuthToken) == "" {
		return transitext.Result{}, fmt.Errorf(
			"openai auth token is required: %w",
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
		translated, err := translator.translateBatch(ctx, batch)
		if err != nil {
			return transitext.Result{}, err
		}

		items = append(items, translated...)
	}

	return transitext.Result{
		Provider: "openai",
		Model:    translator.model(),
		Items:    items,
	}, nil
}

// translateBatch translates one request batch in one chat completion call.
func (translator *Translator) translateBatch(
	ctx context.Context,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	input := make([]string, 0, len(request.Items))
	for _, item := range request.Items {
		input = append(input, item.Text)
	}
	userInput, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("openai marshal input: %w", err)
	}

	payload := map[string]any{
		"model": translator.model(),
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": buildSystemPrompt(request, translator.options),
			},
			{
				"role":    "user",
				"content": string(userInput),
			},
		},
	}
	if translator.options.Temperature != 0 {
		payload["temperature"] = translator.options.Temperature
	}
	if translator.options.TopP != 0 {
		payload["top_p"] = translator.options.TopP
	}
	if translator.options.MaxTokens > 0 {
		payload["max_tokens"] = translator.options.MaxTokens
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("openai marshal request: %w", err)
	}

	endpoint := strings.TrimRight(translator.baseURL(), "/") + "/chat/completions"
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("openai build request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+translator.options.AuthToken)

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.httpClient().Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("openai request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("openai read response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf(
			"openai response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	content, err := parseResponseContent(responseBody)
	if err != nil {
		return nil, err
	}
	translations, err := parseJSONArray(content, translator.options.StrictJSONArray)
	if err != nil {
		return nil, fmt.Errorf("openai parse translations: %w", err)
	}
	if len(translations) != len(request.Items) {
		return nil, fmt.Errorf(
			"openai response size mismatch: got %d, want %d: %w",
			len(translations),
			len(request.Items),
			transitext.ErrProviderPermanent,
		)
	}

	items := make([]transitext.TranslatedItem, 0, len(request.Items))
	for index := range request.Items {
		items = append(items, transitext.TranslatedItem{
			ID:   request.Items[index].ID,
			Text: translations[index],
		})
	}

	return items, nil
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

// baseURL returns normalized provider base URL.
func (translator *Translator) baseURL() string {
	baseURL := strings.TrimSpace(translator.options.BaseURL)
	if baseURL == "" {
		return defaultBaseURL
	}

	return baseURL
}

// model returns configured model.
func (translator *Translator) model() string {
	model := strings.TrimSpace(translator.options.Model)
	if model == "" {
		return defaultModel
	}

	return model
}

// buildSystemPrompt builds translation system prompt from options and request.
func buildSystemPrompt(request transitext.Request, options Options) string {
	custom := strings.TrimSpace(options.SystemPrompt)
	if strings.TrimSpace(request.Hints.SystemPrompt) != "" {
		custom = strings.TrimSpace(request.Hints.SystemPrompt)
	}
	if custom != "" {
		return appendInstructions(custom, request, options)
	}

	base := ""
	if strings.TrimSpace(request.SourceLang) != "" &&
		!strings.EqualFold(strings.TrimSpace(request.SourceLang), "auto") {
		base = fmt.Sprintf(
			"Translate the following strings from %s to %s. %s %s",
			request.SourceLang,
			request.TargetLang,
			promptPreserve,
			promptJSONOnly,
		)
	} else {
		base = fmt.Sprintf(
			"Translate the following strings into %s. %s %s",
			request.TargetLang,
			promptPreserve,
			promptJSONOnly,
		)
	}

	return appendInstructions(base, request, options)
}

// appendInstructions appends configured instruction blocks to system prompt.
func appendInstructions(
	base string,
	request transitext.Request,
	options Options,
) string {
	parts := make([]string, 0, 4)
	parts = append(parts, strings.TrimSpace(base))
	if strings.TrimSpace(options.InstructionPrefix) != "" {
		parts = append(parts, strings.TrimSpace(options.InstructionPrefix))
	}
	if strings.TrimSpace(request.Hints.Instructions) != "" {
		parts = append(parts, strings.TrimSpace(request.Hints.Instructions))
	}
	if strings.TrimSpace(options.InstructionSuffix) != "" {
		parts = append(parts, strings.TrimSpace(options.InstructionSuffix))
	}

	return strings.Join(parts, "\n\n")
}

// parseResponseContent extracts assistant message content from completion payload.
func parseResponseContent(payload []byte) (string, error) {
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return "", fmt.Errorf("openai parse response: %w", err)
	}

	choicesValue, ok := decoded["choices"]
	if !ok {
		return "", fmt.Errorf("openai response missing choices: %w", transitext.ErrProviderPermanent)
	}
	choices, ok := choicesValue.([]any)
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("openai response missing choices: %w", transitext.ErrProviderPermanent)
	}

	choice, ok := choices[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("openai response malformed choice: %w", transitext.ErrProviderPermanent)
	}
	messageValue, ok := choice["message"]
	if !ok {
		return "", fmt.Errorf("openai response missing message: %w", transitext.ErrProviderPermanent)
	}
	message, ok := messageValue.(map[string]any)
	if !ok {
		return "", fmt.Errorf("openai response malformed message: %w", transitext.ErrProviderPermanent)
	}
	contentValue, ok := message["content"]
	if !ok {
		return "", fmt.Errorf("openai response missing content: %w", transitext.ErrProviderPermanent)
	}

	switch content := contentValue.(type) {
	case string:
		return strings.TrimSpace(content), nil
	case []any:
		var builder strings.Builder
		for _, item := range content {
			part, ok := item.(map[string]any)
			if !ok {
				continue
			}
			textValue, ok := part["text"]
			if !ok {
				continue
			}
			text, ok := textValue.(string)
			if !ok {
				continue
			}

			builder.WriteString(text)
		}
		out := strings.TrimSpace(builder.String())
		if out == "" {
			return "", fmt.Errorf("openai response empty content: %w", transitext.ErrProviderPermanent)
		}

		return out, nil
	default:
		return "", fmt.Errorf("openai response unknown content type: %w", transitext.ErrProviderPermanent)
	}
}

// parseJSONArray parses response as JSON array of strings.
func parseJSONArray(content string, strict bool) ([]string, error) {
	var out []string
	if err := json.Unmarshal([]byte(content), &out); err == nil {
		return out, nil
	}
	if strict {
		return nil, errors.New("expected strict JSON array")
	}

	start := strings.Index(content, "[")
	end := strings.LastIndex(content, "]")
	if start >= 0 && end > start {
		if err := json.Unmarshal([]byte(content[start:end+1]), &out); err == nil {
			return out, nil
		}
	}

	return nil, errors.New("expected JSON array")
}
