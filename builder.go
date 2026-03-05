// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

// Builder constructs translator pipelines with fluent API.
type Builder struct {
	// translator stores current pipeline head.
	translator Translator
}

// Wrap starts pipeline builder from base translator.
func Wrap(translator Translator) *Builder {
	return &Builder{translator: translator}
}

// Retry wraps current pipeline with RetryTranslator.
func (builder *Builder) Retry(options RetryOptions) *Builder {
	if builder == nil || builder.translator == nil {
		return builder
	}

	builder.translator = NewRetryTranslator(builder.translator, options)
	return builder
}

// RateLimit wraps current pipeline with RateLimitTranslator.
func (builder *Builder) RateLimit(options RateLimitOptions) *Builder {
	if builder == nil || builder.translator == nil {
		return builder
	}

	builder.translator = NewRateLimitTranslator(builder.translator, options)
	return builder
}

// Cache wraps current pipeline with CacheTranslator.
func (builder *Builder) Cache(options CacheOptions) *Builder {
	if builder == nil || builder.translator == nil {
		return builder
	}

	builder.translator = NewCacheTranslator(builder.translator, options)
	return builder
}

// Context wraps current pipeline with ContextTranslator.
func (builder *Builder) Context(options ContextOptions) *Builder {
	if builder == nil || builder.translator == nil {
		return builder
	}

	builder.translator = NewContextTranslator(builder.translator, options)
	return builder
}

// Fallback wraps current pipeline as first provider with fallbacks.
func (builder *Builder) Fallback(translators ...Translator) *Builder {
	if builder == nil || builder.translator == nil {
		return builder
	}

	chain := make([]Translator, 0, len(translators)+1)
	chain = append(chain, builder.translator)
	chain = append(chain, translators...)
	builder.translator = NewFallbackTranslator(chain...)

	return builder
}

// Build returns configured translator pipeline.
func (builder *Builder) Build() Translator {
	if builder == nil {
		return nil
	}

	return builder.translator
}
