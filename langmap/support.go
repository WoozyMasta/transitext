// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package langmap

import (
	_ "embed"
	"encoding/json"
	"slices"
	"strings"
	"sync"
)

// supportRaw stores serialized provider support payload.
type supportRaw struct {
	Providers map[string][]string `json:"providers"`
	Inherits  map[string]string   `json:"inherits"`
}

// supportData stores normalized provider support lookup tables.
type supportData struct {
	Providers map[string]map[string]string
	Inherits  map[string]string
	Provider  []string
}

var (
	//go:embed support.json
	providerSupportDB []byte

	supportOnce sync.Once
	support     supportData
)

// SupportedByProvider resolves and validates language for a provider.
//
// It returns provider-ready language code when support was confirmed.
func SupportedByProvider(provider, input string) (string, bool) {
	data := loadSupport()
	if data == nil {
		return "", false
	}

	providerKey := strings.ToLower(strings.TrimSpace(provider))
	if providerKey == "" {
		return "", false
	}

	set, ok := providerCodeSet(providerKey, data)
	if !ok {
		return "", false
	}

	code, ok := ResolveForProvider(providerKey, input)
	if !ok {
		return "", false
	}

	resolved, ok := findSupportedCode(set, code)
	if !ok {
		return "", false
	}

	return resolved, true
}

// SupportingProviders returns providers that support the given language.
func SupportingProviders(input string) []string {
	data := loadSupport()
	if data == nil {
		return nil
	}

	result := make([]string, 0, len(data.Provider))
	for _, provider := range data.Provider {
		if _, ok := SupportedByProvider(provider, input); !ok {
			continue
		}

		result = append(result, provider)
	}

	return result
}

// Providers returns known provider names from support registry.
func Providers() []string {
	data := loadSupport()
	if data == nil {
		return nil
	}

	result := make([]string, len(data.Provider))
	copy(result, data.Provider)

	return result
}

// SupportedLanguages returns provider-supported language codes.
//
// For inherited providers this function resolves the final provider set.
func SupportedLanguages(provider string) ([]string, bool) {
	data := loadSupport()
	if data == nil {
		return nil, false
	}

	providerKey := strings.ToLower(strings.TrimSpace(provider))
	if providerKey == "" {
		return nil, false
	}

	set, ok := providerCodeSet(providerKey, data)
	if !ok {
		return nil, false
	}

	result := make([]string, 0, len(set))
	for _, code := range set {
		result = append(result, code)
	}
	slices.Sort(result)

	return result, true
}

// loadSupport lazily decodes embedded provider support payload.
func loadSupport() *supportData {
	supportOnce.Do(func() {
		support = supportData{
			Providers: make(map[string]map[string]string),
			Inherits:  make(map[string]string),
			Provider:  make([]string, 0),
		}

		if len(providerSupportDB) == 0 {
			return
		}

		var raw supportRaw
		if err := json.Unmarshal(providerSupportDB, &raw); err != nil {
			support = supportData{}
			return
		}

		for provider, codes := range raw.Providers {
			key := strings.ToLower(strings.TrimSpace(provider))
			if key == "" {
				continue
			}

			codeMap := make(map[string]string, len(codes))
			for _, code := range codes {
				normalized := normalizeLanguageTag(code)
				if normalized == "" {
					continue
				}
				codeMap[strings.ToLower(normalized)] = normalized
			}
			if len(codeMap) == 0 {
				continue
			}
			support.Providers[key] = codeMap
		}

		for provider, base := range raw.Inherits {
			key := strings.ToLower(strings.TrimSpace(provider))
			parent := strings.ToLower(strings.TrimSpace(base))
			if key == "" || parent == "" {
				continue
			}
			support.Inherits[key] = parent
		}

		providerSet := make(map[string]struct{})
		for provider := range support.Providers {
			providerSet[provider] = struct{}{}
		}
		for provider := range support.Inherits {
			providerSet[provider] = struct{}{}
		}

		for provider := range providerSet {
			support.Provider = append(support.Provider, provider)
		}
		slices.Sort(support.Provider)
	})

	if support.Providers == nil && support.Inherits == nil {
		return nil
	}

	return &support
}

// providerCodeSet resolves provider support set following inheritance chain.
func providerCodeSet(
	provider string,
	data *supportData,
) (map[string]string, bool) {
	visited := make(map[string]struct{})
	current := provider
	for current != "" {
		if _, seen := visited[current]; seen {
			return nil, false
		}
		visited[current] = struct{}{}

		if set, ok := data.Providers[current]; ok {
			return set, true
		}

		next, ok := data.Inherits[current]
		if !ok {
			return nil, false
		}
		current = next
	}

	return nil, false
}

// findSupportedCode matches language code against provider support set.
func findSupportedCode(set map[string]string, code string) (string, bool) {
	candidate := normalizeLanguageTag(code)
	for candidate != "" {
		if value, ok := set[strings.ToLower(candidate)]; ok {
			return value, true
		}

		candidate = trimLastSubtag(candidate)
	}

	return "", false
}

// trimLastSubtag removes last BCP-47 subtag from language code.
func trimLastSubtag(code string) string {
	index := strings.LastIndex(code, "-")
	if index <= 0 {
		return ""
	}

	return code[:index]
}
