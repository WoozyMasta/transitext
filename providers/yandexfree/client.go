// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package yandexfree provides translation via unofficial Yandex endpoint.
package yandexfree

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

	"github.com/google/uuid"
	"github.com/woozymasta/transitext"
)

const (
	// defaultBaseURL is unofficial Yandex API base URL.
	defaultBaseURL = "https://translate.yandex.net/api/v1/tr.json"

	// defaultUserAgent is Yandex mobile app-like UA.
	defaultUserAgent = "ru.yandex.translate/3.20.2024"

	// defaultTimeout is provider HTTP timeout.
	defaultTimeout = 20 * time.Second
)

// Options controls yandexfree provider behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// Request controls low-level HTTP header/cookie/user-agent shaping.
	Request *transitext.HTTPRequestOptions `json:"request,omitempty" yaml:"request,omitempty"`

	// BaseURL overrides Yandex free API base URL.
	BaseURL string `json:"base_url,omitempty" yaml:"base_url,omitempty"`

	// UserAgent overrides default user-agent.
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// MaxItems limits items per one transitext batch.
	MaxItems int `json:"max_items,omitempty" yaml:"max_items,omitempty"`

	// MaxChars limits total chars per one transitext batch.
	MaxChars int `json:"max_chars,omitempty" yaml:"max_chars,omitempty"`
}

// Translator is unofficial Yandex translation provider.
type Translator struct {
	// ucidExpiresAt stores UCID expiration timestamp.
	ucidExpiresAt time.Time

	// client performs HTTP calls.
	client *http.Client

	// baseURL stores normalized API base URL.
	baseURL string

	// maxItems limits batch size by items.
	maxItems int

	// maxChars limits batch size by chars.
	maxChars int

	// ucidLock guards UCID refresh.
	ucidLock sync.Mutex

	// ucid stores cached UCID.
	ucid uuid.UUID
}

// New creates yandexfree provider.
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

	baseURL := strings.TrimSpace(options.BaseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	maxItems := options.MaxItems
	if maxItems <= 0 {
		maxItems = 10
	}
	maxChars := options.MaxChars
	if maxChars <= 0 {
		maxChars = 2000
	}

	return &Translator{
		client:   client,
		baseURL:  baseURL,
		maxItems: maxItems,
		maxChars: maxChars,
	}
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.Capabilities{
		Provider:             "yandexfree",
		Stability:            transitext.ProviderUnstable,
		OfficialAPI:          false,
		SupportsGlossary:     false,
		SupportsInstructions: false,
		SupportsBatch:        true,
		SupportsHTML:         false,
		MaxBatchItems:        translator.maxItems,
		MaxBatchChars:        translator.maxChars,
	}
}

// Translate translates request using unofficial Yandex endpoint.
func (translator *Translator) Translate(
	ctx context.Context,
	request transitext.Request,
) (transitext.Result, error) {
	if err := transitext.ValidateRequest(request); err != nil {
		return transitext.Result{}, err
	}

	batchOptions := request.Batch
	if batchOptions.MaxItems <= 0 {
		batchOptions.MaxItems = translator.maxItems
	}
	if batchOptions.MaxChars <= 0 {
		batchOptions.MaxChars = translator.maxChars
	}
	if batchOptions.OnOverflow == "" {
		batchOptions.OnOverflow = transitext.OverflowSplit
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
		Provider: "yandexfree",
		Items:    items,
	}, nil
}

// translateBatch translates items one-by-one through /translate endpoint.
func (translator *Translator) translateBatch(
	ctx context.Context,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	ucid := translator.getOrUpdateUCID()
	items := make([]transitext.TranslatedItem, 0, len(request.Items))

	for _, item := range request.Items {
		text, detectedSource, err := translator.translateOne(
			ctx,
			ucid,
			item.Text,
			request.SourceLang,
			request.TargetLang,
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

// translateOne sends one request to Yandex free endpoint.
func (translator *Translator) translateOne(
	ctx context.Context,
	ucid uuid.UUID,
	text string,
	source string,
	target string,
) (string, string, error) {
	query := url.Values{
		"ucid":   {strings.ReplaceAll(ucid.String(), "-", "")},
		"srv":    {"android"},
		"format": {"text"},
	}
	endpoint := translator.baseURL + "/translate?" + query.Encode()

	source = strings.TrimSpace(source)
	target = strings.TrimSpace(target)
	lang := yandexHotPatch(target)
	if source != "" && !strings.EqualFold(source, "auto") {
		lang = yandexHotPatch(source) + "-" + yandexHotPatch(target)
	}

	form := url.Values{
		"text": {text},
		"lang": {lang},
	}
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", "", fmt.Errorf("yandexfree build request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.client.Do(httpRequest)
	if err != nil {
		return "", "", fmt.Errorf("yandexfree request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("yandexfree read response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", "", fmt.Errorf(
			"yandexfree response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	return parseTranslateResponse(body)
}

// parseTranslateResponse decodes translate endpoint response body.
func parseTranslateResponse(payload []byte) (string, string, error) {
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return "", "", fmt.Errorf("yandexfree parse response: %w", err)
	}

	codeFloat, hasCode := decoded["code"].(float64)
	if hasCode && int(codeFloat) != 200 {
		return "", "", fmt.Errorf("yandexfree api error code=%d: %w", int(codeFloat), transitext.ErrProviderTemporary)
	}

	textsValue, ok := decoded["text"]
	if !ok {
		return "", "", fmt.Errorf("yandexfree missing text: %w", transitext.ErrProviderPermanent)
	}
	texts, ok := textsValue.([]any)
	if !ok || len(texts) == 0 {
		return "", "", fmt.Errorf("yandexfree missing text: %w", transitext.ErrProviderPermanent)
	}
	text, _ := texts[0].(string)
	if text == "" {
		return "", "", fmt.Errorf("yandexfree empty text: %w", transitext.ErrProviderPermanent)
	}

	detectedSource := ""
	if lang, ok := decoded["lang"].(string); ok && lang != "" {
		if split := strings.SplitN(lang, "-", 2); len(split) == 2 {
			detectedSource = reversePatch(split[0])
		}
	}

	return text, detectedSource, nil
}

// getOrUpdateUCID returns cached UCID or rotates it every 360 seconds.
func (translator *Translator) getOrUpdateUCID() uuid.UUID {
	translator.ucidLock.Lock()
	defer translator.ucidLock.Unlock()

	now := time.Now()
	if translator.ucid == uuid.Nil || !now.Before(translator.ucidExpiresAt) {
		translator.ucid = uuid.New()
		translator.ucidExpiresAt = now.Add(360 * time.Second)
	}

	return translator.ucid
}

// yandexHotPatch converts generic codes to Yandex-specific ones.
func yandexHotPatch(language string) string {
	switch language {
	case "pt-PT":
		return "pt"
	case "pt":
		return "pt-BR"
	case "zh-CN":
		return "zh"
	default:
		return language
	}
}

// reversePatch converts Yandex codes back to generic ones.
func reversePatch(language string) string {
	switch language {
	case "pt":
		return "pt-PT"
	default:
		return language
	}
}
