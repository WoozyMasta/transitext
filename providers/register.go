// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

// Package providers registers built-in transitext providers.
package providers

import (
	"encoding/json"
	"fmt"

	"github.com/woozymasta/transitext"
	"github.com/woozymasta/transitext/providers/azure"
	"github.com/woozymasta/transitext/providers/bingfree"
	"github.com/woozymasta/transitext/providers/deepl"
	"github.com/woozymasta/transitext/providers/deeplfree"
	"github.com/woozymasta/transitext/providers/google"
	"github.com/woozymasta/transitext/providers/googlefree"
	"github.com/woozymasta/transitext/providers/libre"
	"github.com/woozymasta/transitext/providers/microsoftfree"
	"github.com/woozymasta/transitext/providers/openai"
	"github.com/woozymasta/transitext/providers/yandex"
	"github.com/woozymasta/transitext/providers/yandexfree"
)

// NewDefaultRegistry creates registry with built-in provider factories.
func NewDefaultRegistry() (*transitext.ProviderRegistry, error) {
	registry := transitext.NewProviderRegistry()
	if err := RegisterDefaults(registry); err != nil {
		return nil, err
	}

	return registry, nil
}

// RegisterDefaults registers all built-in provider factories.
func RegisterDefaults(registry *transitext.ProviderRegistry) error {
	if registry == nil {
		return fmt.Errorf("registry is nil: %w", transitext.ErrInvalidRequest)
	}

	if err := registry.Register("googlefree", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[googlefree.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode googlefree options: %w", err)
		}

		return googlefree.New(decoded), nil
	}); err != nil {
		return err
	}

	if err := registry.Register("microsoftfree", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[microsoftfree.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode microsoftfree options: %w", err)
		}

		return microsoftfree.New(decoded), nil
	}); err != nil {
		return err
	}

	if err := registry.Register("bingfree", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[bingfree.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode bingfree options: %w", err)
		}

		return bingfree.New(decoded), nil
	}); err != nil {
		return err
	}

	if err := registry.Register("deepl", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[deepl.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode deepl options: %w", err)
		}

		return deepl.New(decoded), nil
	}); err != nil {
		return err
	}

	if err := registry.Register("deeplfree", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[deeplfree.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode deeplfree options: %w", err)
		}

		return deeplfree.New(decoded), nil
	}); err != nil {
		return err
	}

	if err := registry.Register("google", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[google.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode google options: %w", err)
		}

		return google.New(decoded), nil
	}); err != nil {
		return err
	}

	if err := registry.Register("azure", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[azure.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode azure options: %w", err)
		}

		return azure.New(decoded), nil
	}); err != nil {
		return err
	}

	if err := registry.Register("yandex", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[yandex.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode yandex options: %w", err)
		}

		return yandex.New(decoded), nil
	}); err != nil {
		return err
	}

	if err := registry.Register("yandexfree", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[yandexfree.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode yandexfree options: %w", err)
		}

		return yandexfree.New(decoded), nil
	}); err != nil {
		return err
	}

	if err := registry.Register("openai", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[openai.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode openai options: %w", err)
		}

		return openai.New(decoded), nil
	}); err != nil {
		return err
	}

	if err := registry.Register("libre", func(
		options map[string]any,
	) (transitext.Translator, error) {
		decoded, err := decodeOptions[libre.Options](options)
		if err != nil {
			return nil, fmt.Errorf("decode libre options: %w", err)
		}

		return libre.New(decoded), nil
	}); err != nil {
		return err
	}

	return nil
}

// decodeOptions decodes generic map options into typed provider options.
func decodeOptions[T any](options map[string]any) (T, error) {
	var out T
	if options == nil {
		return out, nil
	}

	payload, err := json.Marshal(options)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(payload, &out); err != nil {
		return out, err
	}

	return out, nil
}
