package rpurl_test

import (
	"testing"

	"github.com/mschwager/riplink/src/rpurl"
)

func TestIsRelativeEmpty(t *testing.T) {
	urlIn := ""

	result, err := rpurl.IsRelative(urlIn)

	if result != true || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestIsRelativeValid(t *testing.T) {
	urlIn := "https://example.com"

	result, err := rpurl.IsRelative(urlIn)

	if result != false || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestIsRelativeOnlyPath(t *testing.T) {
	urlIn := "/test"

	result, err := rpurl.IsRelative(urlIn)

	if result != true || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestIsRelativeOnlyFragment(t *testing.T) {
	urlIn := "#fragment"

	result, err := rpurl.IsRelative(urlIn)

	if result != true || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestIsHttpSchemeEmpty(t *testing.T) {
	urlIn := ""

	result, err := rpurl.IsHttpScheme(urlIn)

	if result != true || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestIsHttpSchemeValid(t *testing.T) {
	urlIn := "www.example.com"

	result, err := rpurl.IsHttpScheme(urlIn)

	if result != true || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestIsHttpSchemeValidHttp(t *testing.T) {
	urlIn := "http://example.com"

	result, err := rpurl.IsHttpScheme(urlIn)

	if result != true || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestIsHttpSchemeValidHttps(t *testing.T) {
	urlIn := "https://example.com"

	result, err := rpurl.IsHttpScheme(urlIn)

	if result != true || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestIsHttpSchemeInvalidMailto(t *testing.T) {
	urlIn := "mailto:test@example.com"

	result, err := rpurl.IsHttpScheme(urlIn)

	if result != false || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestIsHttpSchemeInvalidTel(t *testing.T) {
	urlIn := "tel:867-5309"

	result, err := rpurl.IsHttpScheme(urlIn)

	if result != false || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestAddBaseHostToPath(t *testing.T) {
	base := "https://example.com"
	urlIn := "/test"

	result, err := rpurl.AddBaseHost(base, urlIn)
	expected := "https://example.com/test"

	if result != expected || err != nil {
		t.Error("Failed to parse URL:", result)
	}
}

func TestAddBaseHostToFragment(t *testing.T) {
	base := "https://example.com"
	urlIn := "#fragment"

	result, err := rpurl.AddBaseHost(base, urlIn)
	expected := "https://example.com#fragment"

	if result != expected || err != nil {
		t.Error("Failed to parse URL:", result)
	}
}

func TestAddBaseHostJustFragment(t *testing.T) {
	base := "https://example.com"
	urlIn := "#"

	result, err := rpurl.AddBaseHost(base, urlIn)
	expected := "https://example.com"

	if result != expected || err != nil {
		t.Error("Failed to parse URL:", result)
	}
}
