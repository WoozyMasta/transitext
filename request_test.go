// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

package transitext

import (
	"context"
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

func TestValidateRequestMaxTextChars(t *testing.T) {
	t.Parallel()

	err := ValidateRequest(Request{
		TargetLang: "en",
		Batch: BatchOptions{
			MaxTextChars: 3,
		},
		Items: []Item{
			{ID: "1", Text: "abcd"},
		},
	})
	if !errors.Is(err, ErrTextTooLong) {
		t.Fatalf("error = %v, want ErrTextTooLong", err)
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

func TestResolveBatchOptionsFromCapabilities(t *testing.T) {
	t.Parallel()

	options := ResolveBatchOptions(BatchOptions{}, Capabilities{
		MaxBatchItems: 7,
		MaxBatchChars: 777,
		MaxTextChars:  77,
	})

	if options.MaxItems != 7 {
		t.Fatalf("max items = %d, want 7", options.MaxItems)
	}
	if options.MaxChars != 777 {
		t.Fatalf("max chars = %d, want 777", options.MaxChars)
	}
	if options.MaxTextChars != 77 {
		t.Fatalf("max text chars = %d, want 77", options.MaxTextChars)
	}
	if options.OnOverflow != OverflowSplit {
		t.Fatalf("on overflow = %q, want %q", options.OnOverflow, OverflowSplit)
	}
}

func TestSplitRequestMaxTextCharsError(t *testing.T) {
	t.Parallel()

	request := Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "1", Text: "abcdef"},
		},
	}
	_, err := SplitRequest(request, BatchOptions{
		MaxItems:     5,
		MaxChars:     100,
		MaxTextChars: 3,
	})
	if !errors.Is(err, ErrTextTooLong) {
		t.Fatalf("error = %v, want ErrTextTooLong", err)
	}
}

func TestTranslateBatches(t *testing.T) {
	t.Parallel()

	request := Request{
		TargetLang: "en",
		Items: []Item{
			{ID: "1", Text: "aaa"},
			{ID: "2", Text: "bbb"},
			{ID: "3", Text: "ccc"},
		},
		Batch: BatchOptions{
			MaxItems: 2,
		},
	}

	var calls int
	items, err := TranslateBatches(
		context.Background(),
		request,
		Capabilities{},
		func(_ context.Context, batch Request) ([]TranslatedItem, error) {
			calls++
			out := make([]TranslatedItem, 0, len(batch.Items))
			for _, item := range batch.Items {
				out = append(out, TranslatedItem{
					ID:   item.ID,
					Text: item.Text + "_ok",
				})
			}

			return out, nil
		},
	)
	if err != nil {
		t.Fatalf("TranslateBatches error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("batch calls = %d, want 2", calls)
	}
	if len(items) != 3 {
		t.Fatalf("items len = %d, want 3", len(items))
	}
	if items[0].Text != "aaa_ok" || items[2].Text != "ccc_ok" {
		t.Fatalf("unexpected items: %#v", items)
	}
}
