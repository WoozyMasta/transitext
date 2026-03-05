// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
)

// TranslatorFactory creates translator from generic provider options map.
type TranslatorFactory func(options map[string]any) (Translator, error)

// ProviderRegistry stores provider factories by id.
type ProviderRegistry struct {
	// factories holds registered provider factories.
	factories map[string]TranslatorFactory

	// lock guards factories map.
	lock sync.RWMutex
}

// NewProviderRegistry creates empty provider registry.
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		factories: make(map[string]TranslatorFactory),
	}
}

// Register adds or replaces provider factory by id.
func (registry *ProviderRegistry) Register(
	id string,
	factory TranslatorFactory,
) error {
	normalized := normalizeProviderID(id)
	if normalized == "" {
		return fmt.Errorf("provider id is required: %w", ErrInvalidRequest)
	}
	if factory == nil {
		return fmt.Errorf("provider factory is required: %w", ErrInvalidRequest)
	}

	registry.lock.Lock()
	defer registry.lock.Unlock()
	registry.factories[normalized] = factory

	return nil
}

// Build creates translator instance by provider id.
func (registry *ProviderRegistry) Build(
	id string,
	options map[string]any,
) (Translator, error) {
	normalized := normalizeProviderID(id)
	if normalized == "" {
		return nil, fmt.Errorf("provider id is required: %w", ErrInvalidRequest)
	}

	registry.lock.RLock()
	factory, ok := registry.factories[normalized]
	registry.lock.RUnlock()
	if !ok {
		return nil, fmt.Errorf(
			"provider %q is not registered: %w",
			normalized,
			ErrInvalidRequest,
		)
	}

	normalizedOptions := maps.Clone(options)
	if normalizedOptions == nil {
		normalizedOptions = make(map[string]any)
	}

	translator, err := factory(normalizedOptions)
	if err != nil {
		return nil, fmt.Errorf("build provider %q: %w", normalized, err)
	}

	return translator, nil
}

// IDs returns sorted list of registered provider IDs.
func (registry *ProviderRegistry) IDs() []string {
	registry.lock.RLock()
	ids := make([]string, 0, len(registry.factories))
	for id := range registry.factories {
		ids = append(ids, id)
	}
	registry.lock.RUnlock()
	slices.Sort(ids)

	return ids
}

// normalizeProviderID normalizes provider id format.
func normalizeProviderID(id string) string {
	return strings.ToLower(strings.TrimSpace(id))
}
