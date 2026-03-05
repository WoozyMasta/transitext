// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext_test

import (
	"context"
	"fmt"
	"time"

	"github.com/woozymasta/transitext"
)

func Example_pipeline() {
	primary := &exampleStaticTranslator{
		capabilities: transitext.Capabilities{
			Provider: "googlefree",
		},
		result: transitext.Result{
			Provider: "googlefree",
			Items: []transitext.TranslatedItem{
				{ID: "1", Text: "Привет"},
			},
		},
	}
	fallback := &exampleStaticTranslator{
		capabilities: transitext.Capabilities{
			Provider: "microsoft",
		},
		result: transitext.Result{
			Provider: "microsoft",
			Items: []transitext.TranslatedItem{
				{ID: "1", Text: "Привет"},
			},
		},
	}

	pipeline := transitext.NewCacheTranslator(
		transitext.NewRateLimitTranslator(
			transitext.NewRetryTranslator(
				transitext.NewFallbackTranslator(primary, fallback),
				transitext.RetryOptions{
					Attempts: 3,
					Delay:    5 * time.Millisecond,
				},
			),
			transitext.RateLimitOptions{
				MinInterval: 10 * time.Millisecond,
			},
		),
		transitext.CacheOptions{},
	)

	result, err := pipeline.Translate(context.Background(), transitext.Request{
		SourceLang: "en",
		TargetLang: "ru",
		Items: []transitext.Item{
			{ID: "1", Text: "Hello"},
		},
	})
	if err != nil {
		fmt.Println("error")
		return
	}

	fmt.Printf("%s %s\n", result.Provider, result.Items[0].Text)
	// Output: googlefree Привет
}

// exampleStaticTranslator is deterministic example translator.
type exampleStaticTranslator struct {
	// capabilities returns static capability payload.
	capabilities transitext.Capabilities

	// result returns static successful response.
	result transitext.Result
}

// Capabilities returns static capabilities.
func (translator *exampleStaticTranslator) Capabilities() transitext.Capabilities {
	return translator.capabilities
}

// Translate returns static result.
func (translator *exampleStaticTranslator) Translate(
	_ context.Context,
	_ transitext.Request,
) (transitext.Result, error) {
	return translator.result, nil
}
