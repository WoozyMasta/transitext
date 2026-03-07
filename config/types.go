// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package config

import (
	"github.com/woozymasta/transitext"
	"github.com/woozymasta/transitext/providers/azure"
	"github.com/woozymasta/transitext/providers/bingfree"
	"github.com/woozymasta/transitext/providers/deepl"
	"github.com/woozymasta/transitext/providers/deeplfree"
	"github.com/woozymasta/transitext/providers/google"
	"github.com/woozymasta/transitext/providers/googlefree"
	"github.com/woozymasta/transitext/providers/libre"
	"github.com/woozymasta/transitext/providers/microsoft"
	"github.com/woozymasta/transitext/providers/openai"
	"github.com/woozymasta/transitext/providers/yandex"
	"github.com/woozymasta/transitext/providers/yandexfree"
)

// Config is a full transitext configuration document.
type Config struct {
	// Providers contains per-provider transport and API options.
	Providers Providers `json:"providers" yaml:"providers"`

	// Wrappers contains pipeline wrapper options applied around provider.
	Wrappers Wrappers `json:"wrappers" yaml:"wrappers"`
}

// Providers groups built-in provider option blocks.
type Providers struct {

	// Microsoft configures unofficial Microsoft Edge backend.
	Microsoft microsoft.Options `json:"microsoft" yaml:"microsoft"`

	// Google configures official Google Translate API backend.
	Google google.Options `json:"google" yaml:"google"`

	// Yandex configures official Yandex Cloud Translate backend.
	Yandex yandex.Options `json:"yandex" yaml:"yandex"`

	// Azure configures official Azure Translator backend.
	Azure azure.Options `json:"azure" yaml:"azure"`

	// Libre configures LibreTranslate backend.
	Libre libre.Options `json:"libre" yaml:"libre"`

	// DeepL configures official DeepL API backend.
	DeepL deepl.Options `json:"deepl" yaml:"deepl"`

	// BingFree configures unofficial Bing web backend.
	BingFree bingfree.Options `json:"bingfree" yaml:"bingfree"`

	// YandexFree configures unofficial Yandex web backend.
	YandexFree yandexfree.Options `json:"yandexfree" yaml:"yandexfree"`

	// DeepLFree configures unofficial DeepL web backend.
	DeepLFree deeplfree.Options `json:"deeplfree" yaml:"deeplfree"`

	// GoogleFree configures unofficial Google web backend.
	GoogleFree googlefree.Options `json:"googlefree" yaml:"googlefree"`

	// OpenAI configures OpenAI-compatible chat translation backend.
	OpenAI openai.Options `json:"openai" yaml:"openai"`
}

// Wrappers groups wrapper options used by transitext.Wrap(...).
type Wrappers struct {

	// Context configures context marker injection mode.
	Context transitext.ContextOptions `json:"context" yaml:"context"`

	// Retry configures automatic retries for retryable errors.
	Retry transitext.RetryOptions `json:"retry" yaml:"retry"`

	// Partial configures per-item fallback behavior on batch failures.
	Partial transitext.PartialOptions `json:"partial" yaml:"partial"`

	// LongText configures auto-splitting for oversized text items.
	LongText transitext.LongTextOptions `json:"long_text" yaml:"long_text"`

	// RateLimit configures minimum interval between outbound requests.
	RateLimit transitext.RateLimitOptions `json:"rate_limit" yaml:"rate_limit"`

	// Cache configures in-memory response caching.
	Cache transitext.CacheOptions `json:"cache" yaml:"cache"`
}
