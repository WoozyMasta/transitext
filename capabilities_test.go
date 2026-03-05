// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import "testing"

func TestNewCapabilities(t *testing.T) {
	t.Parallel()

	capabilities := NewCapabilities(
		"demo",
		ProviderStable,
		true,
		CapabilitiesOptions{
			SupportsBatch: true,
			SupportsHTML:  true,
			MaxBatchItems: 12,
			MaxBatchChars: 3456,
			MaxTextChars:  789,
		},
	)
	if capabilities.Provider != "demo" {
		t.Fatalf("provider = %q, want demo", capabilities.Provider)
	}
	if capabilities.Stability != ProviderStable {
		t.Fatalf("stability = %q, want %q", capabilities.Stability, ProviderStable)
	}
	if !capabilities.OfficialAPI {
		t.Fatal("official api = false, want true")
	}
	if !capabilities.SupportsBatch || !capabilities.SupportsHTML {
		t.Fatalf("supports flags = %#v", capabilities)
	}
	if capabilities.MaxBatchItems != 12 ||
		capabilities.MaxBatchChars != 3456 ||
		capabilities.MaxTextChars != 789 {
		t.Fatalf("limits = %#v", capabilities)
	}
}
