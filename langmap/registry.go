// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package langmap

import (
	_ "embed"
	"encoding/json"
	"strings"
	"sync"
	"unicode"
)

// registryData holds generated alias mapping data.
type registryData struct {
	Aliases           map[string]string            `json:"aliases"`
	LanguageNames     map[string]string            `json:"language_names"`
	LocaleAliases     map[string]string            `json:"locale_aliases"`
	ProviderOverrides map[string]map[string]string `json:"provider_overrides"`
}

var (
	//go:embed registry.json
	languagesRegistryDB []byte

	registryOnce sync.Once
	registry     registryData
)

// Normalize resolves language alias/human name into canonical language tag.
func Normalize(input string) (string, bool) {
	data := loadRegistry()
	if data == nil {
		return "", false
	}

	code := normalizeLanguageTag(input)
	if code != "" {
		if _, ok := data.LanguageNames[strings.ToLower(code)]; ok {
			return code, true
		}
	}

	key := normalizeAliasKey(input)
	if key == "" {
		return "", false
	}

	if value, ok := data.Aliases[key]; ok {
		return value, true
	}
	if value, ok := data.LocaleAliases[key]; ok {
		return value, true
	}

	return "", false
}

// NormalizeWithDefaultRegion normalizes language input and adds default region
// when normalized tag does not already include region subtag.
func NormalizeWithDefaultRegion(input, defaultRegion string) (string, bool) {
	code, ok := Normalize(input)
	if !ok {
		return "", false
	}

	if hasRegionSubtag(code) {
		return code, true
	}

	region, ok := normalizeRegion(defaultRegion)
	if !ok {
		return code, true
	}

	return appendRegionSubtag(code, region), true
}

// ResolveForProvider resolves input language into provider-oriented value.
func ResolveForProvider(provider, input string) (string, bool) {
	data := loadRegistry()
	if data == nil {
		return "", false
	}

	canonical, ok := Normalize(input)
	if !ok {
		return "", false
	}

	providerKey := strings.ToLower(strings.TrimSpace(provider))
	codeKey := strings.ToLower(canonical)
	if providerMap, found := data.ProviderOverrides[providerKey]; found {
		if value, ok := providerMap[codeKey]; ok {
			return value, true
		}
	}

	if providerKey == "deepl" || providerKey == "deeplfree" {
		return deeplDefaultCode(canonical), true
	}

	return canonical, true
}

// DisplayName returns human-friendly language name for canonical or alias input.
func DisplayName(input string) (string, bool) {
	data := loadRegistry()
	if data == nil {
		return "", false
	}

	canonical, ok := Normalize(input)
	if !ok {
		return "", false
	}

	value, ok := data.LanguageNames[strings.ToLower(canonical)]
	return value, ok
}

// AliasesCount returns number of generated language aliases.
func AliasesCount() int {
	data := loadRegistry()
	if data == nil {
		return 0
	}

	return len(data.Aliases)
}

// loadRegistry lazily decodes embedded generated mapping data.
func loadRegistry() *registryData {
	registryOnce.Do(func() {
		registry = registryData{
			Aliases:           make(map[string]string),
			LanguageNames:     make(map[string]string),
			LocaleAliases:     make(map[string]string),
			ProviderOverrides: make(map[string]map[string]string),
		}
		if len(languagesRegistryDB) == 0 {
			return
		}
		if err := json.Unmarshal(languagesRegistryDB, &registry); err != nil {
			registry = registryData{}
		}
	})
	if registry.Aliases == nil && registry.LanguageNames == nil {
		return nil
	}

	return &registry
}

// normalizeAliasKey normalizes alias/human-name lookup key.
func normalizeAliasKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}
	value = strings.ReplaceAll(value, "_", " ")
	value = strings.ReplaceAll(value, "-", " ")

	var builder strings.Builder
	builder.Grow(len(value))
	space := false
	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsNumber(r):
			builder.WriteRune(r)
			space = false
		case unicode.IsSpace(r), r == '/', r == ',', r == '(', r == ')':
			if space {
				continue
			}
			builder.WriteByte(' ')
			space = true
		}
	}

	return strings.TrimSpace(builder.String())
}

// normalizeLanguageTag canonicalizes BCP-47-like language tags.
func normalizeLanguageTag(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.Trim(value, "-")
	if value == "" {
		return ""
	}

	parts := strings.Split(value, "-")
	for index := range parts {
		parts[index] = strings.TrimSpace(parts[index])
		if parts[index] == "" {
			return ""
		}
		switch {
		case index == 0:
			parts[index] = strings.ToLower(parts[index])
		case len(parts[index]) == 2 || len(parts[index]) == 3 && isNumeric(parts[index]):
			parts[index] = strings.ToUpper(parts[index])
		case len(parts[index]) == 4:
			parts[index] = strings.ToUpper(parts[index][:1]) + strings.ToLower(parts[index][1:])
		default:
			parts[index] = strings.ToLower(parts[index])
		}
	}

	return strings.Join(parts, "-")
}

// isNumeric reports whether input contains only decimal digits.
func isNumeric(value string) bool {
	for _, r := range value {
		if !unicode.IsNumber(r) {
			return false
		}
	}

	return value != ""
}

// deeplDefaultCode derives default DeepL code for canonical BCP-47 tag.
func deeplDefaultCode(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}

	lower := strings.ToLower(code)
	switch {
	case lower == "zh", strings.HasPrefix(lower, "zh-"):
		return "ZH"
	case lower == "pt-br":
		return "PT-BR"
	case lower == "pt", lower == "pt-pt":
		return "PT-PT"
	case lower == "en", lower == "en-us":
		return "EN-US"
	case lower == "en-gb":
		return "EN-GB"
	}

	tag := normalizeLanguageTag(code)
	if tag == "" {
		return strings.ToUpper(code)
	}

	parts := strings.Split(tag, "-")
	if len(parts) == 1 {
		return strings.ToUpper(parts[0])
	}
	if len(parts[1]) == 2 {
		return strings.ToUpper(parts[0] + "-" + parts[1])
	}

	return strings.ToUpper(parts[0])
}

// hasRegionSubtag reports whether canonical BCP-47-like tag has region part.
func hasRegionSubtag(tag string) bool {
	parts := strings.Split(tag, "-")
	if len(parts) < 2 {
		return false
	}

	for index := 1; index < len(parts); index++ {
		part := parts[index]
		if len(part) == 2 || (len(part) == 3 && isNumeric(part)) {
			return true
		}
	}

	return false
}

// normalizeRegion extracts canonical region subtag from input value.
func normalizeRegion(value string) (string, bool) {
	value = normalizeLanguageTag(value)
	if value == "" {
		return "", false
	}

	parts := strings.Split(value, "-")
	if len(parts) == 1 {
		if len(parts[0]) == 2 || (len(parts[0]) == 3 && isNumeric(parts[0])) {
			return strings.ToUpper(parts[0]), true
		}

		return "", false
	}

	for index := 1; index < len(parts); index++ {
		part := parts[index]
		if len(part) == 2 || (len(part) == 3 && isNumeric(part)) {
			return strings.ToUpper(part), true
		}
	}

	return "", false
}

// appendRegionSubtag appends region to canonical tag preserving script subtags.
func appendRegionSubtag(tag, region string) string {
	if tag == "" || region == "" {
		return tag
	}

	parts := strings.Split(tag, "-")
	if len(parts) >= 2 && len(parts[1]) == 4 {
		return parts[0] + "-" + parts[1] + "-" + region
	}

	return parts[0] + "-" + region
}
