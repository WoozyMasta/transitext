// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package microsoft provides translation via Edge unofficial endpoint.
package microsoft

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

	// modeCustomHeaders uses provided auth headers directly.
	modeCustomHeaders = "custom_headers"
)

// Options controls microsoft translator behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// Request controls low-level HTTP header/cookie/user-agent shaping.
	Request *transitext.HTTPRequestOptions `json:"request,omitempty" yaml:"request,omitempty"`

	// AuthURL overrides auth endpoint URL.
	AuthURL string `json:"auth_url,omitempty" yaml:"auth_url,omitempty" jsonschema:"format=uri,example=https://edge.microsoft.com/translate/auth"`

	// TranslateURL overrides translate endpoint URL.
	TranslateURL string `json:"translate_url,omitempty" yaml:"translate_url,omitempty" jsonschema:"format=uri,example=https://api-edge.cognitive.microsofttranslator.com/translate"`

	// Mode selects auth mode: "edge_free" or "custom_headers".
	Mode string `json:"mode,omitempty" yaml:"mode,omitempty" jsonschema:"enum=edge_free,enum=custom_headers,default=edge_free"`

	// AuthenticationHeaders are applied directly in custom_headers mode.
	AuthenticationHeaders map[string]string `json:"authentication_headers,omitempty" yaml:"authentication_headers,omitempty" jsonschema:"maxProperties=32"`

	// TranslateOptions adds optional query params to translate endpoint.
	TranslateOptions map[string]string `json:"translate_options,omitempty" yaml:"translate_options,omitempty" jsonschema:"maxProperties=32"`

	// UserAgent overrides default request user agent.
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty" jsonschema:"maxLength=512"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" jsonschema:"minimum=0,default=20000000000"`

	// MaxItems limits items per one provider HTTP request.
	MaxItems int `json:"max_items,omitempty" yaml:"max_items,omitempty" jsonschema:"minimum=1,maximum=1000,default=20"`

	// MaxChars limits total chars per one provider HTTP request.
	MaxChars int `json:"max_chars,omitempty" yaml:"max_chars,omitempty" jsonschema:"minimum=1,maximum=50000,default=4000"`
}

// Translator is unofficial edge translate provider.
type Translator struct {
	// client performs HTTP calls.
	client *http.Client

	// authURL stores auth endpoint URL.
	authURL *url.URL

	// translateURL stores translate endpoint URL.
	translateURL *url.URL

	// authenticationHeaders stores custom auth headers.
	authenticationHeaders map[string]string

	// translateOptions stores extra query options.
	translateOptions map[string]string

	// useCustomHeadersMode selects custom_headers auth mode.
	useCustomHeadersMode bool

	// maxItems limits batch item count.
	maxItems int

	// maxChars limits batch total chars.
	maxChars int
}

// New creates microsoft translator.
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
		client:               client,
		authURL:              parseURLOrDefault(authURL, defaultAuthURL),
		translateURL:         parseURLOrDefault(translateURL, defaultTranslateURL),
		useCustomHeadersMode: isCustomHeadersMode(options.Mode),
		authenticationHeaders: cloneHeaders(
			options.AuthenticationHeaders,
		),
		translateOptions: cloneHeaders(options.TranslateOptions),
		maxItems:         maxItems,
		maxChars:         maxChars,
	}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.NewCapabilities(
		"microsoft",
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
	authHeaders, err := translator.resolveAuthHeaders(ctx)
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
			return translator.translateBatch(batchCtx, authHeaders, batch)
		},
	)
	if err != nil {
		return transitext.Result{}, err
	}

	return transitext.Result{
		Provider: "microsoft",
		Items:    items,
	}, nil
}

// resolveAuthHeaders resolves headers for selected mode.
func (translator *Translator) resolveAuthHeaders(
	ctx context.Context,
) (map[string]string, error) {
	if translator.useCustomHeadersMode {
		if len(translator.authenticationHeaders) == 0 {
			return nil, fmt.Errorf(
				"microsoft authentication_headers are required in custom_headers mode: %w",
				transitext.ErrInvalidRequest,
			)
		}

		return cloneHeaders(translator.authenticationHeaders), nil
	}

	token, err := translator.fetchToken(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Authorization": normalizeBearer(token)}, nil
}

// fetchToken obtains edge auth token.
func (translator *Translator) fetchToken(ctx context.Context) (string, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		translator.authURL.String(),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("build microsoft auth request: %w", err)
	}

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("microsoft auth request: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read microsoft auth response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", fmt.Errorf(
			"microsoft auth response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	token := strings.TrimSpace(string(body))
	if token == "" {
		return "", fmt.Errorf("microsoft empty token: %w", transitext.ErrProviderPermanent)
	}

	return token, nil
}

// translateBatch translates one batch with a single API call.
func (translator *Translator) translateBatch(
	ctx context.Context,
	authHeaders map[string]string,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	source := strings.TrimSpace(request.SourceLang)
	isAutoSource := source == "" || strings.EqualFold(source, "auto")

	endpoint := *translator.translateURL

	query := endpoint.Query()
	query.Set("api-version", "3.0")
	query.Set("includeSentenceLength", "false")
	query.Set("to", request.TargetLang)
	if !isAutoSource {
		query.Set("from", source)
	}
	for key, value := range translator.translateOptions {
		normalizedKey := strings.TrimSpace(key)
		if normalizedKey == "" {
			continue
		}
		if strings.TrimSpace(value) == "" {
			continue
		}

		query.Set(normalizedKey, strings.TrimSpace(value))
	}
	endpoint.RawQuery = query.Encode()

	payload := make([]map[string]string, 0, len(request.Items))
	for _, item := range request.Items {
		payload = append(payload, map[string]string{"text": item.Text})
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal microsoft payload: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint.String(),
		bytes.NewReader(payloadBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("build microsoft translate request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	for key, value := range authHeaders {
		httpRequest.Header.Set(key, value)
	}

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("microsoft translate request: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read microsoft translate response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf(
			"microsoft translate response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	return parseBatchResponse(request.Items, body)
}

// normalizeBearer normalizes token into Bearer header value.
func normalizeBearer(token string) string {
	trimmed := strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(trimmed), "bearer ") {
		return trimmed
	}

	return "Bearer " + trimmed
}

// normalizeMode returns supported mode or default.
func isCustomHeadersMode(mode string) bool {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case modeCustomHeaders:
		return true
	default:
		return false
	}
}

// cloneHeaders copies key-value map.
func cloneHeaders(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}

	out := make(map[string]string, len(in))
	for key, value := range in {
		normalizedKey := strings.TrimSpace(key)
		if normalizedKey == "" {
			continue
		}

		out[normalizedKey] = strings.TrimSpace(value)
	}
	if len(out) == 0 {
		return nil
	}

	return out
}

// parseURLOrDefault parses candidate URL or returns parsed fallback.
func parseURLOrDefault(candidate string, fallback string) *url.URL {
	parsed, err := url.Parse(strings.TrimSpace(candidate))
	if err == nil && parsed != nil && parsed.Scheme != "" && parsed.Host != "" {
		return parsed
	}

	parsed, _ = url.Parse(fallback)
	return parsed
}

// parseBatchResponse decodes edge batch translation response.
func parseBatchResponse(
	input []transitext.Item,
	payload []byte,
) ([]transitext.TranslatedItem, error) {
	var responseData []map[string]any
	if err := json.Unmarshal(payload, &responseData); err != nil {
		return nil, fmt.Errorf("parse microsoft response: %w", err)
	}
	if len(responseData) != len(input) {
		return nil, fmt.Errorf(
			"microsoft response size mismatch: got %d, want %d: %w",
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
				"microsoft response missing translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		translations, ok := translationsValue.([]any)
		if !ok || len(translations) == 0 {
			return nil, fmt.Errorf(
				"microsoft response missing translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		first, ok := translations[0].(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"microsoft response malformed translation at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}
		text, ok := first["text"].(string)
		if !ok || text == "" {
			return nil, fmt.Errorf(
				"microsoft response missing text at index %d: %w",
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
