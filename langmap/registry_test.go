package langmap

import "testing"

func TestNormalizeByName(t *testing.T) {
	code, ok := Normalize("English")
	if !ok {
		t.Fatalf("Normalize(English) not resolved")
	}
	if code != "en" {
		t.Fatalf("Normalize(English)=%q, want en", code)
	}
}

func TestDeprecatedAlias(t *testing.T) {
	code, ok := Normalize("iw")
	if !ok {
		t.Fatalf("Normalize(iw) not resolved")
	}
	if code != "he" {
		t.Fatalf("Normalize(iw)=%q, want he", code)
	}
}

func TestProviderResolveDeepL(t *testing.T) {
	code, ok := ResolveForProvider("deepl", "chinesesimp")
	if !ok {
		t.Fatalf("ResolveForProvider(deepl,chinesesimp) not resolved")
	}
	if code != "ZH" {
		t.Fatalf("ResolveForProvider(deepl,chinesesimp)=%q, want ZH", code)
	}
}

func TestNormalizeWithDefaultRegion_BaseLanguage(t *testing.T) {
	code, ok := NormalizeWithDefaultRegion("russian", "RU")
	if !ok {
		t.Fatalf("NormalizeWithDefaultRegion(russian,RU) not resolved")
	}
	if code != "ru-RU" {
		t.Fatalf("NormalizeWithDefaultRegion(russian,RU)=%q, want ru-RU", code)
	}
}

func TestNormalizeWithDefaultRegion_AlreadyRegional(t *testing.T) {
	code, ok := NormalizeWithDefaultRegion("ru-RU", "KZ")
	if !ok {
		t.Fatalf("NormalizeWithDefaultRegion(ru-RU,KZ) not resolved")
	}
	if code != "ru-RU" {
		t.Fatalf("NormalizeWithDefaultRegion(ru-RU,KZ)=%q, want ru-RU", code)
	}
}

func TestNormalizeWithDefaultRegion_DefaultLocaleInput(t *testing.T) {
	code, ok := NormalizeWithDefaultRegion("russian", "ru-RU")
	if !ok {
		t.Fatalf("NormalizeWithDefaultRegion(russian,ru-RU) not resolved")
	}
	if code != "ru-RU" {
		t.Fatalf(
			"NormalizeWithDefaultRegion(russian,ru-RU)=%q, want ru-RU",
			code,
		)
	}
}
