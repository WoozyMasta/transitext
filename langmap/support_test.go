package langmap

import (
	"slices"
	"testing"
)

func TestSupportedByProvider_Direct(t *testing.T) {
	code, ok := SupportedByProvider("deepl", "russian")
	if !ok {
		t.Fatalf("SupportedByProvider(deepl,russian) not supported")
	}
	if code != "ru" {
		t.Fatalf("SupportedByProvider(deepl,russian)=%q, want ru", code)
	}
}

func TestSupportedByProvider_Inherited(t *testing.T) {
	code, ok := SupportedByProvider("googlefree", "english")
	if !ok {
		t.Fatalf("SupportedByProvider(googlefree,english) not supported")
	}
	if code != "en" {
		t.Fatalf("SupportedByProvider(googlefree,english)=%q, want en", code)
	}
}

func TestSupportedByProvider_Unsupported(t *testing.T) {
	if code, ok := SupportedByProvider("deepl", "udm"); ok {
		t.Fatalf("SupportedByProvider(deepl,udm)=%q, want unsupported", code)
	}
}

func TestSupportingProviders(t *testing.T) {
	providers := SupportingProviders("russian")
	if len(providers) == 0 {
		t.Fatal("SupportingProviders(russian) empty")
	}

	expected := []string{"azure", "deepl", "google", "libre", "yandex"}
	for _, provider := range expected {
		if !slices.Contains(providers, provider) {
			t.Fatalf("SupportingProviders(russian) missing %q", provider)
		}
	}
}

func TestProviders(t *testing.T) {
	providers := Providers()
	if len(providers) == 0 {
		t.Fatal("Providers() empty")
	}
	if !slices.Contains(providers, "deepl") {
		t.Fatal("Providers() missing deepl")
	}
	if !slices.Contains(providers, "googlefree") {
		t.Fatal("Providers() missing googlefree")
	}
}

func TestSupportedLanguages(t *testing.T) {
	languages, ok := SupportedLanguages("googlefree")
	if !ok {
		t.Fatal("SupportedLanguages(googlefree) not resolved")
	}
	if len(languages) == 0 {
		t.Fatal("SupportedLanguages(googlefree) empty")
	}
	if !slices.Contains(languages, "en") {
		t.Fatal("SupportedLanguages(googlefree) missing en")
	}
}
