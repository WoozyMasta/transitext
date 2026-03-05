// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/transitext

/*
Package langmap provides language alias normalization and provider-oriented
language code resolution.

The package uses generated mapping data and is intended for cases where your
input language value is not stable: human names, legacy aliases, locale-like
forms, or project-specific labels.

It exposes small helpers to:
  - normalize names and aliases to canonical BCP-47 tags;
  - resolve provider-specific target/source codes;
  - validate provider language support from embedded support matrix;
  - list providers that support a given language;
  - get a display name from canonical or alias input.

Minimal flow:

	code, ok := langmap.Normalize("chinesesimp")
	// code == "zh-Hans", ok == true

	deeplCode, ok := langmap.ResolveForProvider("deepl", "chinesesimp")
	// deeplCode == "ZH", ok == true

	checked, ok := langmap.SupportedByProvider("deepl", "russian")
	// checked == "ru", ok == true

	providers := langmap.SupportingProviders("ukrainian")
	// providers contains providers where language is supported now.

Use this package before provider calls when user-facing language input should
be accepted in multiple forms but translation backends still require strict
codes.
*/
package langmap
