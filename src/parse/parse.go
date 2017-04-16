package parse

import (
	"bytes"
	"errors"

	"golang.org/x/net/html"
)

func BytesToHtmlNode(s []byte) (node *html.Node, err error) {
	node, err = html.Parse(bytes.NewReader(s))
	if err != nil {
		return nil, errors.New("Could not parse HTML bytes.")
	}

	return node, nil
}

func NodeIterHelper(node *html.Node, elements chan *html.Node, filter func(*html.Node) bool) {
	if filter(node) {
		elements <- node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		NodeIterHelper(child, elements, filter)
	}
}

func NodeIter(node *html.Node, filter func(*html.Node) bool) (result []*html.Node, err error) {
	elements := make(chan *html.Node)
	go func() {
		NodeIterHelper(node, elements, filter)
		close(elements)
	}()

	result = make([]*html.Node, len(elements))

	for element := range elements {
		result = append(result, element)
	}

	return result, nil
}

func Elements(node *html.Node) (result []*html.Node, err error) {
	return NodeIter(node, func(n *html.Node) bool {
		return n.Type == html.ElementNode
	})
}

func Anchors(node *html.Node) (result []*html.Node, err error) {
	return NodeIter(node, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "a"
	})
}

func Href(anchor *html.Node) (attr html.Attribute, err error) {
	for _, attr := range anchor.Attr {
		if attr.Key == "href" {
			return attr, nil
		}
	}

	return html.Attribute{}, errors.New("Could not find anchor href.")
}

func ValidHrefs(anchors []*html.Node) (hrefs []string, errs []error) {
	for _, anchor := range anchors {
		href, err := Href(anchor)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		hrefs = append(hrefs, href.Val)
	}

	return hrefs, errs
}
