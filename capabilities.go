// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

// CapabilitiesOptions configures optional capability fields.
type CapabilitiesOptions struct {
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

// NewCapabilities builds Capabilities with a compact provider-side API.
func NewCapabilities(
	provider string,
	stability ProviderStability,
	officialAPI bool,
	options CapabilitiesOptions,
) Capabilities {
	return Capabilities{
		Provider:             provider,
		Stability:            stability,
		OfficialAPI:          officialAPI,
		SupportsGlossary:     options.SupportsGlossary,
		SupportsInstructions: options.SupportsInstructions,
		SupportsBatch:        options.SupportsBatch,
		SupportsHTML:         options.SupportsHTML,
		MaxBatchItems:        options.MaxBatchItems,
		MaxBatchChars:        options.MaxBatchChars,
		MaxTextChars:         options.MaxTextChars,
	}
}
