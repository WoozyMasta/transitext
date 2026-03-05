// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import "context"

// ProviderStability describes provider stability level.
type ProviderStability string

const (
	// ProviderStable marks official stable provider APIs.
	ProviderStable ProviderStability = "stable"

	// ProviderUnstable marks unofficial or unstable provider APIs.
	ProviderUnstable ProviderStability = "unstable"
)

// OverflowMode describes what to do when batch limits are exceeded.
type OverflowMode string

const (
	// OverflowSplit splits request into multiple batches.
	OverflowSplit OverflowMode = "split"

	// OverflowError fails request when limits are exceeded.
	OverflowError OverflowMode = "error"
)

// Translator translates text items.
type Translator interface {
	// Translate translates request and returns ordered per-item result.
	Translate(ctx context.Context, request Request) (Result, error)

	// Capabilities reports provider runtime capabilities.
	Capabilities() Capabilities
}

// Request describes translation input.
type Request struct {
	// Metadata stores opaque client-side metadata.
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// SourceLang is source language tag (or "auto"/empty for auto detect).
	SourceLang string `json:"source_lang,omitempty" yaml:"source_lang,omitempty"`

	// TargetLang is output language tag.
	TargetLang string `json:"target_lang" yaml:"target_lang"`

	// Hints carries optional provider-independent translation hints.
	Hints Hints `json:"hints" yaml:"hints"`

	// Items is ordered translation input list.
	Items []Item `json:"items" yaml:"items"`

	// Batch controls chunking behavior for this request.
	Batch BatchOptions `json:"batch" yaml:"batch"`
}

// Item is one translation input unit.
type Item struct {
	// ID is stable item identifier preserved in output.
	ID string `json:"id,omitempty" yaml:"id,omitempty"`

	// Text is input text to translate.
	Text string `json:"text" yaml:"text"`
}

// Hints carries optional translation guidance.
type Hints struct {
	// Glossary maps terms that should be preferred by provider.
	Glossary map[string]string `json:"glossary,omitempty" yaml:"glossary,omitempty"`

	// Domain describes business/subject domain.
	Domain string `json:"domain,omitempty" yaml:"domain,omitempty"`

	// Instructions is user guidance added for translation behavior.
	Instructions string `json:"instructions,omitempty" yaml:"instructions,omitempty"`

	// SystemPrompt is provider-level system prompt when supported.
	SystemPrompt string `json:"system_prompt,omitempty" yaml:"system_prompt,omitempty"`

	// Preserve lists placeholders/tokens that must stay unchanged.
	Preserve []string `json:"preserve,omitempty" yaml:"preserve,omitempty"`
}

// BatchOptions controls provider chunking.
type BatchOptions struct {
	// OnOverflow controls behavior when limits are exceeded.
	OnOverflow OverflowMode `json:"on_overflow,omitempty" yaml:"on_overflow,omitempty"`

	// MaxItems limits number of items per request batch.
	MaxItems int `json:"max_items,omitempty" yaml:"max_items,omitempty"`

	// MaxChars limits text characters per request batch.
	MaxChars int `json:"max_chars,omitempty" yaml:"max_chars,omitempty"`

	// MaxTextChars limits one text item length when configured.
	MaxTextChars int `json:"max_text_chars,omitempty" yaml:"max_text_chars,omitempty"`

	// Parallel limits concurrent batch calls (wrapper-level concern).
	Parallel int `json:"parallel,omitempty" yaml:"parallel,omitempty"`
}

// Result is translation output.
type Result struct {
	// Provider identifies provider implementation.
	Provider string `json:"provider,omitempty" yaml:"provider,omitempty"`

	// Model identifies provider model/engine when available.
	Model string `json:"model,omitempty" yaml:"model,omitempty"`

	// Items keeps output in input order.
	Items []TranslatedItem `json:"items,omitempty" yaml:"items,omitempty"`

	// Usage carries provider usage stats when available.
	Usage Usage `json:"usage" yaml:"usage"`
}

// TranslatedItem is one translated output item.
type TranslatedItem struct {
	// Error contains per-item error details if translation failed.
	Error *ItemError `json:"error,omitempty" yaml:"error,omitempty"`

	// ID mirrors input item ID.
	ID string `json:"id,omitempty" yaml:"id,omitempty"`

	// Text is translated output.
	Text string `json:"text,omitempty" yaml:"text,omitempty"`

	// DetectedSource is provider-reported source language.
	DetectedSource string `json:"detected_source,omitempty" yaml:"detected_source,omitempty"`
}

// ItemError describes per-item translation failure.
type ItemError struct {
	// Code is provider or domain error code.
	Code string `json:"code,omitempty" yaml:"code,omitempty"`

	// Message is human-readable error text.
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	// Temporary indicates retry may succeed later.
	Temporary bool `json:"temporary,omitempty" yaml:"temporary,omitempty"`
}

// Usage contains provider token/character accounting.
type Usage struct {
	// CharsIn is counted input characters.
	CharsIn int `json:"chars_in,omitempty" yaml:"chars_in,omitempty"`

	// CharsOut is counted output characters.
	CharsOut int `json:"chars_out,omitempty" yaml:"chars_out,omitempty"`

	// TokensIn is counted input tokens when provider reports it.
	TokensIn int `json:"tokens_in,omitempty" yaml:"tokens_in,omitempty"`

	// TokensOut is counted output tokens when provider reports it.
	TokensOut int `json:"tokens_out,omitempty" yaml:"tokens_out,omitempty"`
}

// Capabilities describes translator features and constraints.
type Capabilities struct {
	// Provider is short provider id.
	Provider string `json:"provider" yaml:"provider"`

	// Stability marks provider API stability.
	Stability ProviderStability `json:"stability" yaml:"stability"`

	// OfficialAPI reports whether provider uses official API.
	OfficialAPI bool `json:"official_api" yaml:"official_api"`

	// SupportsGlossary reports glossary hint support.
	SupportsGlossary bool `json:"supports_glossary,omitempty" yaml:"supports_glossary,omitempty"`

	// SupportsInstructions reports prompt/instruction support.
	SupportsInstructions bool `json:"supports_instructions,omitempty" yaml:"supports_instructions,omitempty"`

	// SupportsBatch reports batch request support.
	SupportsBatch bool `json:"supports_batch,omitempty" yaml:"supports_batch,omitempty"`

	// SupportsHTML reports HTML-aware translation support.
	SupportsHTML bool `json:"supports_html,omitempty" yaml:"supports_html,omitempty"`

	// MaxBatchItems reports hard provider batch item limit if known.
	MaxBatchItems int `json:"max_batch_items,omitempty" yaml:"max_batch_items,omitempty"`

	// MaxBatchChars reports hard provider batch character limit if known.
	MaxBatchChars int `json:"max_batch_chars,omitempty" yaml:"max_batch_chars,omitempty"`

	// MaxTextChars reports hard single text limit if known.
	MaxTextChars int `json:"max_text_chars,omitempty" yaml:"max_text_chars,omitempty"`
}
