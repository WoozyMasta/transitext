// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"errors"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	t.Parallel()

	err := ValidateRequest(Request{
		SourceLang: "auto",
		TargetLang: "en",
		Items: []Item{
			{ID: "1", Text: "bonjour"},
		},
	})
	if err != nil {
		t.Fatalf("ValidateRequest error: %v", err)
	}
}

func TestValidateRequestMissingTargetLang(t *testing.T) {
	t.Parallel()

	err := ValidateRequest(Request{
		Items: []Item{{ID: "1", Text: "bonjour"}},
	})
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("error = %v, want ErrInvalidRequest", err)
	}
}

func TestSplitRequest(t *testing.T) {
	t.Parallel()

	request := Request{
		SourceLang: "auto",
		TargetLang: "en",
		Items: []Item{
			{ID: "1", Text: "a"},
			{ID: "2", Text: "b"},
			{ID: "3", Text: "c"},
		},
	}
	batches, err := SplitRequest(request, BatchOptions{
		MaxItems:   2,
		MaxChars:   100,
		OnOverflow: OverflowSplit,
	})
	if err != nil {
		t.Fatalf("SplitRequest error: %v", err)
	}
	if len(batches) != 2 {
		t.Fatalf("len(batches) = %d, want 2", len(batches))
	}
	if len(batches[0].Items) != 2 || len(batches[1].Items) != 1 {
		t.Fatalf("batch sizes = [%d, %d], want [2,1]", len(batches[0].Items), len(batches[1].Items))
	}
}

func TestSplitRequestOverflowError(t *testing.T) {
	t.Parallel()

	request := Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "1", Text: "a"},
			{ID: "2", Text: "b"},
		},
	}
	_, err := SplitRequest(request, BatchOptions{
		MaxItems:   1,
		MaxChars:   100,
		OnOverflow: OverflowError,
	})
	if !errors.Is(err, ErrBatchTooLarge) {
		t.Fatalf("error = %v, want ErrBatchTooLarge", err)
	}
}
