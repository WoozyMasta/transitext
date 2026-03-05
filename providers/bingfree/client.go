// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package bingfree provides translation via unofficial Bing endpoint.
package bingfree

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/woozymasta/transitext"
)

const (
	// defaultHostURL is Bing host URL.
	defaultHostURL = "https://www.bing.com"

	// defaultTranslatorPath is translator page path used for credentials extraction.
	defaultTranslatorPath = "/translator"

	// defaultTranslatePath is unofficial Bing translate endpoint.
	defaultTranslatePath = "/ttranslatev3"

	// defaultIid is Bing IID query value.
	defaultIid = "translator.5024.1"

	// defaultUserAgent is browser-like UA used for web endpoints.
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"

	// defaultTimeout is provider HTTP timeout.
	defaultTimeout = 20 * time.Second

	// defaultMaxItems is default request item limit per batch.
	defaultMaxItems = 10

	// defaultMaxChars is default request text-char limit per batch.
	defaultMaxChars = 1000

	// defaultMaxTextChars is per-item text length bound.
	defaultMaxTextChars = 1000
)

//nolint:gosec // Static parser pattern, not credential material.
const credentialsPrefix = "var params_AbusePreventionHelper = ["

// Options controls bingfree provider behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// Request controls low-level HTTP header/cookie/user-agent shaping.
	Request *transitext.HTTPRequestOptions `json:"request,omitempty" yaml:"request,omitempty"`

	// HostURL overrides Bing host URL.
	HostURL string `json:"host_url,omitempty" yaml:"host_url,omitempty"`

	// UserAgent overrides default request user agent.
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// MaxItems limits items per one transitext batch.
	MaxItems int `json:"max_items,omitempty" yaml:"max_items,omitempty"`

	// MaxChars limits total chars per one transitext batch.
	MaxChars int `json:"max_chars,omitempty" yaml:"max_chars,omitempty"`

	// MaxTextChars limits one input text length.
	MaxTextChars int `json:"max_text_chars,omitempty" yaml:"max_text_chars,omitempty"`
}

// Translator is unofficial Bing translation provider.
type Translator struct {
	// client performs HTTP calls to provider.
	client *http.Client

	// hostURL stores normalized Bing host URL.
	hostURL string

	// credentials stores cached provider credentials.
	credentials credentials

	// maxItems limits batch size by items.
	maxItems int

	// maxChars limits batch size by chars.
	maxChars int

	// maxTextChars limits one item text length.
	maxTextChars int

	// credentialsLock guards credentials refresh.
	credentialsLock sync.Mutex
}

// credentials stores Bing token material with expiration.
type credentials struct {
	// expiration is token expiration timestamp.
	expiration time.Time

	// token is Bing request token.
	token string

	// key is Bing numeric key.
	key int64
}

// New creates bingfree provider.
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

	hostURL := strings.TrimSpace(options.HostURL)
	if hostURL == "" {
		hostURL = defaultHostURL
	}
	hostURL = strings.TrimRight(hostURL, "/")

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

	return &Translator{
		client:       client,
		hostURL:      hostURL,
		maxItems:     maxItems,
		maxChars:     maxChars,
		maxTextChars: maxTextChars,
	}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.NewCapabilities(
		"bingfree",
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

// Translate translates request using unofficial Bing endpoint.
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
		Provider: "bingfree",
		Items:    items,
	}, nil
}

// translateBatch translates batch items one-by-one using shared credentials.
func (translator *Translator) translateBatch(
	ctx context.Context,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	creds, err := translator.getOrUpdateCredentials(ctx)
	if err != nil {
		return nil, err
	}

	source := strings.TrimSpace(request.SourceLang)
	if source == "" || strings.EqualFold(source, "auto") {
		source = "auto-detect"
	}
	target := strings.TrimSpace(request.TargetLang)

	items := make([]transitext.TranslatedItem, 0, len(request.Items))
	for _, item := range request.Items {
		text, detectedSource, err := translator.translateOne(
			ctx,
			creds,
			bingHotPatch(source),
			bingHotPatch(target),
			item.Text,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, transitext.TranslatedItem{
			ID:             item.ID,
			Text:           text,
			DetectedSource: detectedSource,
		})
	}

	return items, nil
}

// getOrUpdateCredentials returns cached credentials or refreshes when expired.
func (translator *Translator) getOrUpdateCredentials(
	ctx context.Context,
) (credentials, error) {
	translator.credentialsLock.Lock()
	defer translator.credentialsLock.Unlock()

	now := time.Now()
	if translator.credentials.token != "" && now.Before(translator.credentials.expiration) {
		return translator.credentials, nil
	}

	urlValue := translator.hostURL + defaultTranslatorPath
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, urlValue, nil)
	if err != nil {
		return credentials{}, fmt.Errorf("bingfree build credentials request: %w", err)
	}

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.client.Do(httpRequest)
	if err != nil {
		return credentials{}, fmt.Errorf("bingfree credentials request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return credentials{}, fmt.Errorf("bingfree read credentials response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return credentials{}, fmt.Errorf(
			"bingfree credentials response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	parsed, err := parseCredentials(string(body))
	if err != nil {
		return credentials{}, err
	}
	translator.credentials = parsed

	return parsed, nil
}

// translateOne translates one item via Bing translate endpoint.
func (translator *Translator) translateOne(
	ctx context.Context,
	creds credentials,
	source string,
	target string,
	text string,
) (string, string, error) {
	values := url.Values{
		"fromLang": {source},
		"text":     {text},
		"to":       {target},
		"token":    {creds.token},
		"key":      {strconv.FormatInt(creds.key, 10)},
	}

	ig := strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
	endpoint := fmt.Sprintf(
		"%s%s?isVertical=1&IG=%s&IID=%s",
		translator.hostURL,
		defaultTranslatePath,
		ig,
		defaultIid,
	)
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", "", fmt.Errorf("bingfree build translate request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.client.Do(httpRequest)
	if err != nil {
		return "", "", fmt.Errorf("bingfree translate request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("bingfree read response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", "", fmt.Errorf(
			"bingfree response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	return parseTranslateResponse(body)
}

// parseCredentials extracts token and key from translator page script.
func parseCredentials(html string) (credentials, error) {
	index := strings.Index(html, credentialsPrefix)
	if index < 0 {
		return credentials{}, fmt.Errorf(
			"bingfree credentials not found: %w",
			transitext.ErrProviderPermanent,
		)
	}

	start := index + len(credentialsPrefix)
	end := strings.Index(html[start:], "]")
	if end < 0 {
		return credentials{}, fmt.Errorf(
			"bingfree credentials malformed: %w",
			transitext.ErrProviderPermanent,
		)
	}
	payload := html[start : start+end]
	parts := strings.SplitN(payload, ",", 2)
	if len(parts) != 2 {
		return credentials{}, fmt.Errorf(
			"bingfree credentials malformed: %w",
			transitext.ErrProviderPermanent,
		)
	}

	key, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		key = time.Now().UnixMilli()
	}

	token := strings.TrimSpace(parts[1])
	token = strings.Trim(token, `"`)
	if token == "" {
		return credentials{}, fmt.Errorf(
			"bingfree token not found: %w",
			transitext.ErrProviderPermanent,
		)
	}

	expiration := time.UnixMilli(key + 3600000)
	return credentials{
		token:      token,
		key:        key,
		expiration: expiration,
	}, nil
}

// parseTranslateResponse parses Bing translation response payload.
func parseTranslateResponse(payload []byte) (string, string, error) {
	var raw any
	if err := json.Unmarshal(payload, &raw); err != nil {
		return "", "", fmt.Errorf("bingfree parse response: %w", err)
	}

	if rootMap, ok := raw.(map[string]any); ok {
		if _, hasStatus := rootMap["statusCode"]; hasStatus {
			return "", "", fmt.Errorf("bingfree api error: %w", transitext.ErrProviderTemporary)
		}
	}

	rootArray, ok := raw.([]any)
	if !ok || len(rootArray) == 0 {
		return "", "", fmt.Errorf("bingfree malformed response: %w", transitext.ErrProviderPermanent)
	}

	first, ok := rootArray[0].(map[string]any)
	if !ok {
		return "", "", fmt.Errorf("bingfree malformed response: %w", transitext.ErrProviderPermanent)
	}
	translationsValue, ok := first["translations"]
	if !ok {
		return "", "", fmt.Errorf("bingfree missing translations: %w", transitext.ErrProviderPermanent)
	}
	translations, ok := translationsValue.([]any)
	if !ok || len(translations) == 0 {
		return "", "", fmt.Errorf("bingfree missing translations: %w", transitext.ErrProviderPermanent)
	}
	translation, ok := translations[0].(map[string]any)
	if !ok {
		return "", "", fmt.Errorf("bingfree malformed translation: %w", transitext.ErrProviderPermanent)
	}
	text, _ := translation["text"].(string)
	if text == "" {
		return "", "", fmt.Errorf("bingfree empty translation: %w", transitext.ErrProviderPermanent)
	}

	detectedSource := ""
	if detectedValue, ok := first["detectedLanguage"]; ok {
		if detected, ok := detectedValue.(map[string]any); ok {
			if language, ok := detected["language"].(string); ok {
				detectedSource = language
			}
		}
	}

	return text, detectedSource, nil
}

// bingHotPatch converts generic codes to Bing-specific ones.
func bingHotPatch(language string) string {
	switch language {
	case "lg":
		return "lug"
	case "no":
		return "nb"
	case "ny":
		return "nya"
	case "rn":
		return "run"
	case "sr":
		return "sr-Cyrl"
	case "mn":
		return "mn-Cyrl"
	case "tlh":
		return "tlh-Latn"
	case "zh-CN":
		return "zh-Hans"
	case "zh-TW":
		return "zh-Hant"
	default:
		return language
	}
}
