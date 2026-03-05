// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"slices"
	"strings"
	"sync"
)

// CacheOptions controls cache wrapper behavior.
type CacheOptions struct {
	// IncludeHints includes request hints in cache key.
	IncludeHints bool `json:"include_hints,omitempty" yaml:"include_hints,omitempty"`

	// IncludeMetadata includes request metadata in cache key.
	IncludeMetadata bool `json:"include_metadata,omitempty" yaml:"include_metadata,omitempty"`
}

// CacheTranslator memoizes successful translation requests in memory.
type CacheTranslator struct {
	// base is wrapped translator.
	base Translator

	// entries stores cached translation results by key.
	entries map[string]Result

	// lock guards map operations.
	lock sync.RWMutex

	// options controls cache key behavior.
	options CacheOptions
}

// NewCacheTranslator creates cache wrapper around translator.
func NewCacheTranslator(base Translator, options CacheOptions) *CacheTranslator {
	return &CacheTranslator{
		base:    base,
		options: options,
		entries: make(map[string]Result),
	}
}

// Capabilities returns wrapped translator capabilities.
func (translator *CacheTranslator) Capabilities() Capabilities {
	return translator.base.Capabilities()
}

// Translate returns cached result when available, otherwise calls base.
func (translator *CacheTranslator) Translate(
	ctx context.Context,
	request Request,
) (Result, error) {
	key := buildCacheKey(translator.base.Capabilities().Provider, request, translator.options)

	translator.lock.RLock()
	cached, ok := translator.entries[key]
	translator.lock.RUnlock()
	if ok {
		return cloneResult(cached), nil
	}

	result, err := translator.base.Translate(ctx, request)
	if err != nil {
		return Result{}, err
	}

	translator.lock.Lock()
	translator.entries[key] = cloneResult(result)
	translator.lock.Unlock()

	return result, nil
}

// Len returns number of cached entries.
func (translator *CacheTranslator) Len() int {
	translator.lock.RLock()
	defer translator.lock.RUnlock()

	return len(translator.entries)
}

// Clear removes all cached entries.
func (translator *CacheTranslator) Clear() {
	translator.lock.Lock()
	defer translator.lock.Unlock()

	translator.entries = make(map[string]Result)
}

// buildCacheKey builds deterministic cache key for translation request.
func buildCacheKey(
	provider string,
	request Request,
	options CacheOptions,
) string {
	var builder strings.Builder
	builder.Grow(estimateCacheKeySize(request, options))

	_, _ = builder.WriteString("p=")
	_, _ = builder.WriteString(provider)
	_, _ = builder.WriteString("|s=")
	_, _ = builder.WriteString(request.SourceLang)
	_, _ = builder.WriteString("|t=")
	_, _ = builder.WriteString(request.TargetLang)
	_, _ = builder.WriteString("|i=")

	for _, item := range request.Items {
		_, _ = builder.WriteString(item.ID)
		builder.WriteByte('=')
		_, _ = builder.WriteString(item.Text)
		builder.WriteByte(';')
	}

	if options.IncludeHints {
		_, _ = builder.WriteString("|h:")
		_, _ = builder.WriteString(request.Hints.Domain)
		builder.WriteByte('|')
		_, _ = builder.WriteString(request.Hints.Instructions)
		builder.WriteByte('|')
		_, _ = builder.WriteString(request.Hints.SystemPrompt)
		builder.WriteByte('|')
		writeStringSlice(&builder, request.Hints.Preserve)
		builder.WriteByte('|')
		writeStringMap(&builder, request.Hints.Glossary)
	}

	if options.IncludeMetadata {
		_, _ = builder.WriteString("|m:")
		writeStringMap(&builder, request.Metadata)
	}

	return builder.String()
}

// cloneResult clones result to avoid external mutation of cached values.
func cloneResult(result Result) Result {
	cloned := result
	cloned.Items = append([]TranslatedItem(nil), result.Items...)
	for index := range cloned.Items {
		if cloned.Items[index].Error != nil {
			copied := *cloned.Items[index].Error
			cloned.Items[index].Error = &copied
		}
	}

	return cloned
}

// estimateCacheKeySize estimates key builder capacity.
func estimateCacheKeySize(request Request, options CacheOptions) int {
	size := 32 + len(request.SourceLang) + len(request.TargetLang)
	for _, item := range request.Items {
		size += len(item.ID) + len(item.Text) + 2
	}
	if options.IncludeHints {
		size += len(request.Hints.Domain) +
			len(request.Hints.Instructions) +
			len(request.Hints.SystemPrompt) +
			len(request.Hints.Preserve)*4 +
			len(request.Hints.Glossary)*16
	}
	if options.IncludeMetadata {
		size += len(request.Metadata) * 16
	}

	return size
}

// writeStringSlice writes string slice into key builder.
func writeStringSlice(builder *strings.Builder, values []string) {
	for _, value := range values {
		_, _ = builder.WriteString(value)
		builder.WriteByte(',')
	}
}

// writeStringMap writes deterministic key/value map into key builder.
func writeStringMap(builder *strings.Builder, values map[string]string) {
	if len(values) == 0 {
		return
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for _, key := range keys {
		_, _ = builder.WriteString(key)
		builder.WriteByte('=')
		_, _ = builder.WriteString(values[key])
		builder.WriteByte(';')
	}
}
