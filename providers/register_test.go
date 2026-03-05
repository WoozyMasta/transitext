// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package providers

import (
	"slices"
	"testing"
)

func TestNewDefaultRegistry(t *testing.T) {
	t.Parallel()

	registry, err := NewDefaultRegistry()
	if err != nil {
		t.Fatalf("NewDefaultRegistry error: %v", err)
	}

	ids := registry.IDs()
	for _, provider := range []string{
		"azure",
		"bingfree",
		"deepl",
		"deeplfree",
		"google",
		"googlefree",
		"libre",
		"microsoft",
		"openai",
		"yandex",
		"yandexfree",
	} {
		if !slices.Contains(ids, provider) {
			t.Fatalf("registry ids %#v missing provider %q", ids, provider)
		}
	}
}

func TestBuildOpenAIFromMapOptions(t *testing.T) {
	t.Parallel()

	registry, err := NewDefaultRegistry()
	if err != nil {
		t.Fatalf("NewDefaultRegistry error: %v", err)
	}

	translator, err := registry.Build("openai", map[string]any{
		"auth_token":        "token",
		"model":             "gpt-4o-mini",
		"strict_json_array": true,
		"batch_max_items":   5,
	})
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	if got := translator.Capabilities().Provider; got != "openai" {
		t.Fatalf("provider = %q, want openai", got)
	}
}
