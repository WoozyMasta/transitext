// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

/*
Package transitext provides provider-agnostic translation contracts and
pipeline wrappers for batch localization tooling.

Core package responsibilities:

  - Request/Result data model shared by all providers.
  - Validation and batch splitting utilities.
  - Runtime wrappers: retry, fallback, rate limit, cache, context passing,
    partial item retry, and long-text splitting.

Provider-specific implementations are in subpackages under
github.com/woozymasta/transitext/providers.

Typical flow:

 1. Build a provider translator (for example via providers registry).
 2. Compose wrappers with Wrap(...).
 3. Send Request with Items and optional Batch limits.
 4. Read ordered Result items with per-item errors when supported.

Minimal pipeline example:

	pipeline := transitext.Wrap(base).
		Retry(transitext.RetryOptions{Attempts: 3}).
		RateLimit(transitext.RateLimitOptions{MinInterval: 500 * time.Millisecond}).
		LongText(transitext.LongTextOptions{MaxTextChars: 4000}).
		Build()

	result, err := pipeline.Translate(ctx, transitext.Request{
		SourceLang: "en",
		TargetLang: "ru",
		Items: []transitext.Item{
			{ID: "hello", Text: "Hello world"},
		},
	})
*/
package transitext
