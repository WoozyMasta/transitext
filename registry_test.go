// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestProviderRegistryBuild(t *testing.T) {
	t.Parallel()

	registry := NewProviderRegistry()
	err := registry.Register("googlefree", func(
		options map[string]any,
	) (Translator, error) {
		_ = options["x"]
		return &registryMockTranslator{
			capabilities: Capabilities{Provider: "googlefree"},
			result:       Result{Provider: "googlefree"},
		}, nil
	})
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}

	translator, err := registry.Build("googlefree", map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	if translator.Capabilities().Provider != "googlefree" {
		t.Fatalf("provider = %q, want googlefree", translator.Capabilities().Provider)
	}
}

func TestProviderRegistryUnknownProvider(t *testing.T) {
	t.Parallel()

	registry := NewProviderRegistry()
	_, err := registry.Build("missing", nil)
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("error = %v, want ErrInvalidRequest", err)
	}
}

func TestProviderRegistryIDs(t *testing.T) {
	t.Parallel()

	registry := NewProviderRegistry()
	_ = registry.Register("b", func(map[string]any) (Translator, error) {
		return &registryMockTranslator{}, nil
	})
	_ = registry.Register("a", func(map[string]any) (Translator, error) {
		return &registryMockTranslator{}, nil
	})
	ids := registry.IDs()
	if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Fatalf("ids = %#v, want [a b]", ids)
	}
}

func TestRateLimitTranslator(t *testing.T) {
	t.Parallel()

	base := &registryMockTranslator{
		capabilities: Capabilities{Provider: "mock"},
		result:       Result{Provider: "mock"},
	}
	limiter := NewRateLimitTranslator(base, RateLimitOptions{
		MinInterval: 30 * time.Millisecond,
	})

	start := time.Now()
	_, err := limiter.Translate(context.Background(), Request{
		TargetLang: "en",
		Items:      []Item{{ID: "1", Text: "hello"}},
	})
	if err != nil {
		t.Fatalf("first Translate error: %v", err)
	}
	_, err = limiter.Translate(context.Background(), Request{
		TargetLang: "en",
		Items:      []Item{{ID: "2", Text: "world"}},
	})
	if err != nil {
		t.Fatalf("second Translate error: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed < 30*time.Millisecond {
		t.Fatalf("elapsed = %v, want >= 30ms", elapsed)
	}
	if base.calls != 2 {
		t.Fatalf("calls = %d, want 2", base.calls)
	}
}

// registryMockTranslator is deterministic translator for registry tests.
type registryMockTranslator struct {
	// capabilities returns static capabilities.
	capabilities Capabilities

	// result is returned from Translate.
	result Result

	// calls tracks Translate invocation count.
	calls int
}

// Capabilities returns static mock capabilities.
func (translator *registryMockTranslator) Capabilities() Capabilities {
	return translator.capabilities
}

// Translate returns predefined result.
func (translator *registryMockTranslator) Translate(
	_ context.Context,
	_ Request,
) (Result, error) {
	translator.calls++
	return translator.result, nil
}
