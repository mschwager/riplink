package rpurl_test

import (
	"testing"

	"github.com/mschwager/riplink/src/rpurl"
)

func TestHasHostEmpty(t *testing.T) {
	urlIn := ""

	result, err := rpurl.HasHost(urlIn)

	if result != false || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestHasHostValid(t *testing.T) {
	urlIn := "https://example.com"

	result, err := rpurl.HasHost(urlIn)

	if result != true || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestHasHostOnlyPath(t *testing.T) {
	urlIn := "/test"

	result, err := rpurl.HasHost(urlIn)

	if result != false || err != nil {
		t.Error("Failed to parse URL:", urlIn)
	}
}

func TestHasHostOnlyFragment(t *testing.T) {
	urlIn := "#fragment"

	result, err := rpurl.HasHost(urlIn)

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
