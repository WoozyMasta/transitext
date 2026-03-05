// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package microsoftfree provides translation via Edge unofficial endpoint.
package microsoftfree

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
	// defaultAuthURL returns bearer token for edge translate endpoint.
	defaultAuthURL = "https://edge.microsoft.com/translate/auth"

	// defaultTranslateURL is edge translate endpoint.
	defaultTranslateURL = "https://api-edge.cognitive.microsofttranslator.com/translate"

	// defaultTimeout is provider HTTP timeout.
	defaultTimeout = 20 * time.Second

	// defaultUserAgent is browser-like UA for edge endpoint calls.
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"

	// defaultMaxItems is default request item limit per batch.
	defaultMaxItems = 20

	// defaultMaxChars is default request text-char limit per batch.
	defaultMaxChars = 4000
)

// Options controls microsoftfree translator behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// Request controls low-level HTTP header/cookie/user-agent shaping.
	Request *transitext.HTTPRequestOptions `json:"request,omitempty" yaml:"request,omitempty"`

	// AuthURL overrides auth endpoint URL.
	AuthURL string `json:"auth_url,omitempty" yaml:"auth_url,omitempty"`

	// TranslateURL overrides translate endpoint URL.
	TranslateURL string `json:"translate_url,omitempty" yaml:"translate_url,omitempty"`

	// UserAgent overrides default request user agent.
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// MaxItems limits items per one provider HTTP request.
	MaxItems int `json:"max_items,omitempty" yaml:"max_items,omitempty"`

	// MaxChars limits total chars per one provider HTTP request.
	MaxChars int `json:"max_chars,omitempty" yaml:"max_chars,omitempty"`
}

// Translator is unofficial edge translate provider.
type Translator struct {
	// client performs HTTP calls.
	client *http.Client

	// authURL stores auth endpoint URL.
	authURL string

	// translateURL stores translate endpoint URL.
	translateURL string

	// maxItems limits batch item count.
	maxItems int

	// maxChars limits batch total chars.
	maxChars int
}

// New creates microsoftfree translator.
func New(options Options) *Translator {
	client := options.HTTPClient
	if client == nil {
		timeout := options.Timeout
		if timeout <= 0 {
			timeout = defaultTimeout
		}
		client = &http.Client{Timeout: timeout}
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

	authURL := strings.TrimSpace(options.AuthURL)
	if authURL == "" {
		authURL = defaultAuthURL
	}

	translateURL := strings.TrimSpace(options.TranslateURL)
	if translateURL == "" {
		translateURL = defaultTranslateURL
	}

	maxItems := options.MaxItems
	if maxItems <= 0 {
		maxItems = defaultMaxItems
	}

	maxChars := options.MaxChars
	if maxChars <= 0 {
		maxChars = defaultMaxChars
	}

	return &Translator{
		client:       client,
		authURL:      authURL,
		translateURL: translateURL,
		maxItems:     maxItems,
		maxChars:     maxChars,
	}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.NewCapabilities(
		"microsoftfree",
		transitext.ProviderUnstable,
		false,
		transitext.CapabilitiesOptions{
			SupportsBatch: true,
			MaxBatchItems: translator.maxItems,
			MaxBatchChars: translator.maxChars,
		},
	)
}

// Translate translates request via edge translate endpoint.
func (translator *Translator) Translate(
	ctx context.Context,
	request transitext.Request,
) (transitext.Result, error) {
	token, err := translator.fetchToken(ctx)
	if err != nil {
		return transitext.Result{}, err
	}

	items, err := transitext.TranslateBatches(
		ctx,
		request,
		translator.Capabilities(),
		func(batchCtx context.Context, batch transitext.Request) (
			[]transitext.TranslatedItem,
			error,
		) {
			return translator.translateBatch(batchCtx, token, batch)
		},
	)
	if err != nil {
		return transitext.Result{}, err
	}

	return transitext.Result{
		Provider: "microsoftfree",
		Items:    items,
	}, nil
}

// fetchToken obtains edge auth token.
func (translator *Translator) fetchToken(ctx context.Context) (string, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		translator.authURL,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("build microsoftfree auth request: %w", err)
	}

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("microsoftfree auth request: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read microsoftfree auth response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", fmt.Errorf(
			"microsoftfree auth response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	token := strings.TrimSpace(string(body))
	if token == "" {
		return "", fmt.Errorf("microsoftfree empty token: %w", transitext.ErrProviderPermanent)
	}

	return token, nil
}

// translateBatch translates one batch with a single API call.
func (translator *Translator) translateBatch(
	ctx context.Context,
	token string,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	source := strings.TrimSpace(request.SourceLang)
	if source == "" || strings.EqualFold(source, "auto") {
		// No reliable auto-detect in free edge mode; keep predictable default.
		source = "en"
	}

	endpoint, err := url.Parse(translator.translateURL)
	if err != nil {
		return nil, fmt.Errorf("parse microsoftfree translate url: %w", err)
	}

	query := endpoint.Query()
	query.Set("api-version", "3.0")
	query.Set("includeSentenceLength", "false")
	query.Set("from", source)
	query.Set("to", request.TargetLang)
	endpoint.RawQuery = query.Encode()

	payload := make([]map[string]string, 0, len(request.Items))
	for _, item := range request.Items {
		payload = append(payload, map[string]string{"text": item.Text})
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal microsoftfree payload: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint.String(),
		bytes.NewReader(payloadBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("build microsoftfree translate request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", token)

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("microsoftfree translate request: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read microsoftfree translate response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf(
			"microsoftfree translate response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	return parseBatchResponse(request.Items, body)
}

// parseBatchResponse decodes edge batch translation response.
func parseBatchResponse(
	input []transitext.Item,
	payload []byte,
) ([]transitext.TranslatedItem, error) {
	var responseData []map[string]any
	if err := json.Unmarshal(payload, &responseData); err != nil {
		return nil, fmt.Errorf("parse microsoftfree response: %w", err)
	}
	if len(responseData) != len(input) {
		return nil, fmt.Errorf(
			"microsoftfree response size mismatch: got %d, want %d: %w",
			len(responseData),
			len(input),
			transitext.ErrProviderPermanent,
		)
	}

	out := make([]transitext.TranslatedItem, 0, len(input))
	for index := range responseData {
		item := responseData[index]
		translationsValue, ok := item["translations"]
		if !ok {
			return nil, fmt.Errorf(
				"microsoftfree response missing translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		translations, ok := translationsValue.([]any)
		if !ok || len(translations) == 0 {
			return nil, fmt.Errorf(
				"microsoftfree response missing translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		first, ok := translations[0].(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"microsoftfree response malformed translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		text, ok := first["text"].(string)
		if !ok || text == "" {
			return nil, fmt.Errorf(
				"microsoftfree response missing text at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}

		detectedSource := ""
		if detectedValue, ok := item["detectedLanguage"]; ok {
			if detectedMap, ok := detectedValue.(map[string]any); ok {
				if detectedLang, ok := detectedMap["language"].(string); ok {
					detectedSource = detectedLang
				}
			}
		}
		if detectedSource == "" {
			if detectedValue, ok := item["detected_language"]; ok {
				if detectedMap, ok := detectedValue.(map[string]any); ok {
					if detectedLang, ok := detectedMap["language"].(string); ok {
						detectedSource = detectedLang
					}
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
