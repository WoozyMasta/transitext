// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package providers

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/woozymasta/transitext"
)

func TestFreeProvidersIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test is skipped in short mode")
	}

	registry, err := NewDefaultRegistry()
	if err != nil {
		t.Fatalf("NewDefaultRegistry error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	providers := []string{
		"googlefree",
		"bingfree",
		"deeplfree",
		"microsoftfree",
		"yandexfree",
	}

	for _, providerID := range providers {
		providerID := providerID
		t.Run(providerID, func(t *testing.T) {
			translator, err := registry.Build(providerID, nil)
			if err != nil {
				t.Fatalf("Build(%q) error: %v", providerID, err)
			}

			result, err := translator.Translate(ctx, transitext.Request{
				SourceLang: "en",
				TargetLang: "ru",
				Items: []transitext.Item{
					{ID: "1", Text: "Hello world"},
				},
			})
			if err != nil {
				if errors.Is(err, transitext.ErrProviderTemporary) {
					t.Skipf("Translate(%q) temporary provider failure: %v", providerID, err)
				}

				t.Fatalf("Translate(%q) error: %v", providerID, err)
			}
			if got := result.Provider; got == "" {
				t.Fatalf("Translate(%q) provider is empty", providerID)
			}
			if len(result.Items) != 1 {
				t.Fatalf("Translate(%q) items len = %d, want 1", providerID, len(result.Items))
			}
			if strings.TrimSpace(result.Items[0].Text) == "" {
				t.Fatalf("Translate(%q) empty translated text", providerID)
			}
		})
	}
}
