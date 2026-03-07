// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package deeplfree provides unofficial DeepL web JSON-RPC translation.
package deeplfree

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/woozymasta/transitext"
)

const (
	// defaultURL is unofficial DeepL web JSON-RPC endpoint.
	defaultURL = "https://www2.deepl.com/jsonrpc"

	// defaultUserAgent is browser-like user agent for web endpoint.
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"

	// defaultSplitMode is sentence splitting mode for JSON-RPC flow.
	defaultSplitMode = "newlines"

	// defaultAcceptLanguage is preferred request locale for web endpoint.
	defaultAcceptLanguage = "en-US,en;q=0.9"

	// defaultTimeout is provider HTTP timeout.
	defaultTimeout = 20 * time.Second

	// defaultMaxItems is default request item limit per batch.
	defaultMaxItems = 20

	// defaultMaxChars is default request text-char limit per batch.
	defaultMaxChars = 5000

	// defaultMaxTextChars is default single text-char limit.
	defaultMaxTextChars = 5000

	// defaultRequestAlternatives is requested alternatives count per text.
	defaultRequestAlternatives = 3
)

// Options controls deeplfree provider behavior.
type Options struct {
	// HTTPClient is optional custom HTTP client.
	HTTPClient *http.Client `json:"-" yaml:"-"`

	// URL overrides DeepL web endpoint URL.
	URL string `json:"url,omitempty" yaml:"url,omitempty" jsonschema:"format=uri,example=https://www2.deepl.com/jsonrpc"`

	// UserAgent overrides default request user agent.
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty" jsonschema:"maxLength=512"`

	// Request controls low-level HTTP header/cookie/user-agent shaping.
	Request *transitext.HTTPRequestOptions `json:"request,omitempty" yaml:"request,omitempty"`

	// AcceptLanguage overrides Accept-Language header value.
	AcceptLanguage string `json:"accept_language,omitempty" yaml:"accept_language,omitempty" jsonschema:"maxLength=128,default=en-US,en;q=0.9"`

	// DLSession optionally sends dl_session cookie for authenticated web mode.
	DLSession string `json:"dl_session,omitempty" yaml:"dl_session,omitempty" jsonschema:"maxLength=512"`

	// SplitMode controls DeepL splitting mode (for example "newlines").
	SplitMode string `json:"split_mode,omitempty" yaml:"split_mode,omitempty" jsonschema:"maxLength=32,default=newlines"`

	// RequestAlternatives sets requestAlternatives for each text.
	RequestAlternatives int `json:"request_alternatives,omitempty" yaml:"request_alternatives,omitempty" jsonschema:"minimum=1,maximum=10,default=3"`

	// MaxItems limits items per one transitext batch.
	MaxItems int `json:"max_items,omitempty" yaml:"max_items,omitempty" jsonschema:"minimum=1,maximum=100,default=20"`

	// MaxChars limits total chars per one transitext batch.
	MaxChars int `json:"max_chars,omitempty" yaml:"max_chars,omitempty" jsonschema:"minimum=1,maximum=50000,default=5000"`

	// MaxTextChars limits one input text length.
	MaxTextChars int `json:"max_text_chars,omitempty" yaml:"max_text_chars,omitempty" jsonschema:"minimum=1,maximum=5000,default=5000"`

	// Timeout is request timeout when HTTPClient is not provided.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" jsonschema:"minimum=0,default=20000000000"`
}

// Translator is unofficial DeepL web translator.
type Translator struct {
	// client performs HTTP calls to provider.
	client *http.Client

	// url stores JSON-RPC endpoint URL.
	url string

	// splitMode stores request split mode.
	splitMode string

	// requestAlternatives controls alternatives count per text item.
	requestAlternatives int

	// maxItems is request item bound.
	maxItems int

	// maxChars is request char bound.
	maxChars int

	// maxTextChars is one text length bound.
	maxTextChars int

	// idCounter stores JSON-RPC request id sequence.
	idCounter atomic.Int64
}

// rpcRequest describes DeepL JSON-RPC request payload.
type rpcRequest struct {
	// JSONRPC stores JSON-RPC protocol version.
	JSONRPC string `json:"jsonrpc"`

	// Method stores rpc method name.
	Method string `json:"method"`

	// Params stores method parameters.
	Params rpcParams `json:"params"`

	// ID stores JSON-RPC request id.
	ID int64 `json:"id"`
}

// rpcParams describes LMT_handle_texts params.
type rpcParams struct {
	// Splitting configures sentence splitting mode.
	Splitting string `json:"splitting"`

	// Lang stores source and target languages.
	Lang rpcLang `json:"lang"`

	// Texts stores input text array.
	Texts []rpcText `json:"texts"`

	// Timestamp stores request timestamp expected by endpoint.
	Timestamp int64 `json:"timestamp"`
}

// rpcLang stores language params for request.
type rpcLang struct {
	// SourceLangUserSelected stores chosen source lang or "auto".
	SourceLangUserSelected string `json:"source_lang_user_selected"`

	// TargetLang stores target language tag.
	TargetLang string `json:"target_lang"`
}

// rpcText stores one text item in request.
type rpcText struct {
	// Text stores source text.
	Text string `json:"text"`

	// RequestAlternatives stores requested alternatives count.
	RequestAlternatives int `json:"requestAlternatives"` //nolint:tagliatelle // DeepL JSON-RPC requires camelCase key.
}

// rpcResponse describes DeepL JSON-RPC response payload.
type rpcResponse struct {
	// Result stores untyped success payload fields.
	Result map[string]json.RawMessage `json:"result"`

	// Error stores JSON-RPC error payload.
	Error *rpcError `json:"error"`
}

// rpcResultText stores one translated text item.
type rpcResultText struct {
	// Text stores translated value.
	Text string `json:"text"`
}

// rpcError stores JSON-RPC error payload.
type rpcError struct {
	// Message stores provider error message.
	Message string `json:"message"`

	// Code stores provider error code.
	Code int `json:"code"`
}

// New creates deeplfree translator.
func New(options Options) *Translator {
	client := options.HTTPClient
	if client == nil {
		timeout := options.Timeout
		if timeout <= 0 {
			timeout = defaultTimeout
		}

		client = &http.Client{Timeout: timeout}
	}

	url := strings.TrimSpace(options.URL)
	if url == "" {
		url = defaultURL
	}

	acceptLanguage := strings.TrimSpace(options.AcceptLanguage)
	if acceptLanguage == "" {
		acceptLanguage = defaultAcceptLanguage
	}

	splitMode := strings.TrimSpace(options.SplitMode)
	if splitMode == "" {
		splitMode = defaultSplitMode
	}

	requestAlternatives := options.RequestAlternatives
	if requestAlternatives <= 0 {
		requestAlternatives = defaultRequestAlternatives
	}

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

	requestOptions := transitext.HTTPRequestOptions{}
	if options.Request != nil {
		requestOptions = *options.Request
	}
	if strings.TrimSpace(requestOptions.UserAgent) == "" {
		requestOptions.UserAgent = options.UserAgent
	}
	if requestOptions.Headers == nil {
		requestOptions.Headers = make(map[string]string, 1)
	}
	if _, ok := requestOptions.Headers["Accept-Language"]; !ok {
		requestOptions.Headers["Accept-Language"] = acceptLanguage
	}
	if requestOptions.Cookies == nil {
		requestOptions.Cookies = make(map[string]string, 1)
	}
	if strings.TrimSpace(options.DLSession) != "" {
		if _, ok := requestOptions.Cookies["dl_session"]; !ok {
			requestOptions.Cookies["dl_session"] = strings.TrimSpace(options.DLSession)
		}
	}
	client = transitext.ConfigureHTTPClient(client, transitext.HTTPRequestDefaults{
		UserAgent: defaultUserAgent,
		Headers: map[string]string{
			"Accept":  "*/*",
			"Origin":  "https://www.deepl.com",
			"Referer": "https://www.deepl.com/",
		},
	}, requestOptions)

	translator := &Translator{
		client:              client,
		url:                 url,
		splitMode:           splitMode,
		requestAlternatives: requestAlternatives,
		maxItems:            maxItems,
		maxChars:            maxChars,
		maxTextChars:        maxTextChars,
	}
	translator.idCounter.Store((time.Now().UnixMilli()%900000 + 100000) * 1000)

	return translator
}

// Capabilities reports provider capabilities.
func (translator *Translator) Capabilities() transitext.Capabilities {
	return transitext.NewCapabilities(
		"deeplfree",
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

// Translate translates input items with unofficial DeepL web endpoint.
func (translator *Translator) Translate(
	ctx context.Context,
	request transitext.Request,
) (transitext.Result, error) {
	out, err := transitext.TranslateBatches(
		ctx,
		request,
		translator.Capabilities(),
		translator.translateBatch,
	)
	if err != nil {
		return transitext.Result{}, err
	}

	return transitext.Result{
		Provider: "deeplfree",
		Items:    out,
	}, nil
}

// translateBatch translates one request batch with one JSON-RPC call.
func (translator *Translator) translateBatch(
	ctx context.Context,
	request transitext.Request,
) ([]transitext.TranslatedItem, error) {
	source := normalizeSourceLang(request.SourceLang)
	target := normalizeTargetLang(request.TargetLang)

	texts := make([]rpcText, 0, len(request.Items))
	iCount := 0
	for _, item := range request.Items {
		texts = append(texts, rpcText{
			Text:                item.Text,
			RequestAlternatives: translator.requestAlternatives,
		})
		iCount += strings.Count(item.Text, "i")
	}

	payload := rpcRequest{
		JSONRPC: "2.0",
		Method:  "LMT_handle_texts",
		Params: rpcParams{
			Splitting: translator.splitMode,
			Lang: rpcLang{
				SourceLangUserSelected: source,
				TargetLang:             target,
			},
			Texts:     texts,
			Timestamp: deeplTimestamp(iCount),
		},
		ID: translator.nextID(),
	}

	body, err := translator.doRPC(ctx, payload)
	if err != nil {
		return nil, err
	}

	response, err := parseRPCResponse(body)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, providerError(response.Error)
	}
	if len(response.Result) == 0 {
		return nil, fmt.Errorf(
			"deeplfree response missing result: %w",
			transitext.ErrProviderPermanent,
		)
	}

	translatedTexts, detectedSource, err := decodeResultFields(response.Result)
	if err != nil {
		return nil, err
	}
	if len(translatedTexts) != len(request.Items) {
		return nil, fmt.Errorf(
			"deeplfree response size mismatch: got %d, want %d: %w",
			len(translatedTexts),
			len(request.Items),
			transitext.ErrProviderPermanent,
		)
	}

	out := make([]transitext.TranslatedItem, 0, len(request.Items))
	for index := range request.Items {
		translated := strings.TrimSpace(translatedTexts[index].Text)
		if translated == "" {
			return nil, fmt.Errorf(
				"deeplfree response missing text at index %d: %w",
				index,
				transitext.ErrProviderPermanent,
			)
		}

		itemSource := detectedSource
		if itemSource == "" && source != "auto" {
			itemSource = source
		}

		out = append(out, transitext.TranslatedItem{
			ID:             request.Items[index].ID,
			Text:           translated,
			DetectedSource: itemSource,
		})
	}

	return out, nil
}

// doRPC sends one JSON-RPC request and returns raw response body.
func (translator *Translator) doRPC(
	ctx context.Context,
	payload rpcRequest,
) ([]byte, error) {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("deeplfree marshal request: %w", err)
	}
	requestBody = patchMethodSpacing(payload.ID, requestBody)

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		translator.url,
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("deeplfree build request: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	//nolint:gosec // Provider intentionally performs outbound HTTP requests.
	response, err := translator.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("deeplfree request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("deeplfree read response: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf(
			"deeplfree response %d: %w",
			response.StatusCode,
			transitext.ErrProviderTemporary,
		)
	}

	return responseBody, nil
}

// parseRPCResponse decodes raw JSON-RPC response.
func parseRPCResponse(payload []byte) (rpcResponse, error) {
	var response rpcResponse
	if err := json.Unmarshal(payload, &response); err != nil {
		return rpcResponse{}, fmt.Errorf("deeplfree parse response: %w", err)
	}

	return response, nil
}

// decodeResultFields extracts translated texts and detected language.
func decodeResultFields(
	result map[string]json.RawMessage,
) ([]rpcResultText, string, error) {
	textsRaw, ok := result["texts"]
	if !ok {
		return nil, "", fmt.Errorf(
			"deeplfree response missing result.texts: %w",
			transitext.ErrProviderPermanent,
		)
	}

	var texts []rpcResultText
	if err := json.Unmarshal(textsRaw, &texts); err != nil {
		return nil, "", fmt.Errorf("deeplfree parse result.texts: %w", err)
	}

	detectedSource := ""
	if langRaw, ok := result["lang"]; ok {
		if err := json.Unmarshal(langRaw, &detectedSource); err != nil {
			return nil, "", fmt.Errorf("deeplfree parse result.lang: %w", err)
		}
		detectedSource = normalizeTargetLang(detectedSource)
	}

	return texts, detectedSource, nil
}

// patchMethodSpacing applies DeepL body-spacing quirk used by web clients.
func patchMethodSpacing(id int64, body []byte) []byte {
	if (id+5)%29 == 0 || (id+3)%13 == 0 {
		return bytes.Replace(body, []byte(`"method":"`), []byte(`"method" : "`), 1)
	}

	return bytes.Replace(body, []byte(`"method":"`), []byte(`"method": "`), 1)
}

// providerError maps rpc error payload into transitext category.
func providerError(value *rpcError) error {
	if value == nil {
		return fmt.Errorf(
			"deeplfree empty provider error: %w",
			transitext.ErrProviderPermanent,
		)
	}

	base := transitext.ErrProviderPermanent
	if isTemporaryErrorCode(value.Code) {
		base = transitext.ErrProviderTemporary
	}

	message := strings.TrimSpace(value.Message)
	if message == "" {
		message = "unknown provider error"
	}

	return fmt.Errorf(
		"deeplfree rpc error code=%d msg=%q: %w",
		value.Code,
		message,
		base,
	)
}

// isTemporaryErrorCode reports whether rpc code is likely retryable.
func isTemporaryErrorCode(code int) bool {
	if code == 429 || code == 503 || code == 504 {
		return true
	}
	if code >= 1042900 && code <= 1043999 {
		return true
	}

	return false
}

// normalizeSourceLang normalizes source language for deeplfree.
func normalizeSourceLang(source string) string {
	value := canonicalLang(source)
	if value == "" {
		return "auto"
	}
	if strings.EqualFold(value, "auto") {
		return "auto"
	}

	return value
}

// normalizeTargetLang normalizes target language for deeplfree.
func normalizeTargetLang(target string) string {
	value := canonicalLang(target)
	if value == "" {
		return target
	}

	return value
}

// canonicalLang uppercases and normalizes separators.
func canonicalLang(language string) string {
	value := strings.TrimSpace(language)
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.ToUpper(value)

	return value
}

// deeplTimestamp computes rounded timestamp expected by web endpoint.
func deeplTimestamp(iCount int) int64 {
	now := time.Now().UnixMilli()
	if iCount <= 0 {
		return now
	}

	divider := int64(iCount + 1)
	return now - (now % divider) + divider
}

// nextID returns next JSON-RPC id.
func (translator *Translator) nextID() int64 {
	return translator.idCounter.Add(1000)
}
