package parse_test

import (
	"testing"

	"github.com/mschwager/riplink/src/parse"
)

func TestBytesToHtmlNodeBasic(t *testing.T) {
	html := []byte("")

	_, err := parse.BytesToHtmlNode(html)

	if err != nil {
		t.Error("Failed to parse HTML:", err)
	}
}

func TestElementsBasic(t *testing.T) {
	html := []byte(`
	<html>
	<head>
	</head>
	<body>
	</body>
	</html>
	`)

	node, err1 := parse.BytesToHtmlNode(html)
	elements, err2 := parse.Elements(node)

	expected_length := 3

	if err1 != nil || err2 != nil || len(elements) != expected_length {
		t.Error("Failed to parse HTML elements: ", elements)
	}
}

func TestAnchorsEmpty(t *testing.T) {
	html := []byte(`
	<html>
	<head>
	</head>
	<body>
	</body>
	</html>
	`)

	node, err1 := parse.BytesToHtmlNode(html)
	anchors, err2 := parse.Anchors(node)

	expected_length := 0

	if err1 != nil || err2 != nil || len(anchors) != expected_length {
		t.Error("Failed to parse HTML anchors: ", err2)
	}
}

func TestAnchorsBasic(t *testing.T) {
	html := []byte(`
	<html>
	<head>
	</head>
	<body>
		<h1>Header</h1>
		<a href="example.com">Test1</a>
		<a href="example.com">Test2</a>
	</body>
	</html>
	`)

	node, err1 := parse.BytesToHtmlNode(html)
	anchors, err2 := parse.Anchors(node)

	expected_length := 2

	if err1 != nil || err2 != nil || len(anchors) != expected_length {
		t.Error("Failed to parse HTML anchors: ", err2)
	}
}

func TestHrefBasic(t *testing.T) {
	html := []byte(`
	<html>
	<head>
	</head>
	<body>
		<a href="example.com">Test1</a>
	</body>
	</html>
	`)

	node, err1 := parse.BytesToHtmlNode(html)
	anchors, err2 := parse.Anchors(node)
	href, err3 := parse.Href(anchors[0])

	expected_href := "example.com"

	if err1 != nil || err2 != nil || err3 != nil || href.Val != expected_href {
		t.Error("Failed to parse HTML anchors: ", err3)
	}
}

func TestHrefNoAnchorHref(t *testing.T) {
	html := []byte(`
	<html>
	<head>
	</head>
	<body>
		<a>Test1</a>
	</body>
	</html>
	`)

	node, err1 := parse.BytesToHtmlNode(html)
	anchors, err2 := parse.Anchors(node)
	href, err3 := parse.Href(anchors[0])

	expected_href := ""

	if err1 != nil || err2 != nil || err3 == nil || href.Val != expected_href {
		t.Error("Failed to parse HTML anchors: ", err3)
	}
}

func TestValidHrefsWithValidHref(t *testing.T) {
	html := []byte(`
	<html>
	<head>
	</head>
	<body>
		<a href="example.com">Test1</a>
	</body>
	</html>
	`)

	node, err1 := parse.BytesToHtmlNode(html)
	anchors, err2 := parse.Anchors(node)
	hrefs, errs := parse.ValidHrefs(anchors)

	expected_href := "example.com"

	if err1 != nil || err2 != nil || len(errs) != 0 || hrefs[0].Val != expected_href {
		t.Error("Failed to parse HTML anchors: ", errs)
	}
}

func TestValidHrefsWithInvalidHref(t *testing.T) {
	html := []byte(`
	<html>
	<head>
	</head>
	<body>
		<a>Test1</a>
	</body>
	</html>
	`)

	node, err1 := parse.BytesToHtmlNode(html)
	anchors, err2 := parse.Anchors(node)
	hrefs, errs := parse.ValidHrefs(anchors)

	if err1 != nil || err2 != nil || len(errs) != 1 || len(hrefs) != 0 {
		t.Error("Failed to parse HTML anchors: ", hrefs)
	}
}
