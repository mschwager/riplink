package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/mschwager/riplink/src/parse"
	"github.com/mschwager/riplink/src/requests"
	"github.com/mschwager/riplink/src/rpurl"

	"golang.org/x/net/html"
)

func main() {
	var queryUrl string
	flag.StringVar(&queryUrl, "url", "https://google.com", "URL to query")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")

	flag.Parse()

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	response, _, err := requests.Request(client, "GET", queryUrl, nil)
	if err != nil {
		panic(err)
	}

	node, err := parse.BytesToHtmlNode(response)
	if err != nil {
		panic(err)
	}

	anchors, err := parse.Anchors(node)
	if err != nil {
		panic(err)
	}

	// Filter invalid HREFs
	var hrefs []html.Attribute
	for _, anchor := range anchors {
		href, err := parse.Href(anchor)
		if err != nil {
			fmt.Println(err)
			continue
		}
		hrefs = append(hrefs, href)
	}

	// Attempt to include hostname in relative paths
	// E.g. Querying https://example.com yields /relative/path
	// gives https://example.com/relative/path
	var urls []string
	for _, href := range hrefs {
		url := href.Val
		hasHost, err := rpurl.HasHost(url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !hasHost {
			url, err = rpurl.AddBaseHost(queryUrl, url)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		urls = append(urls, url)
	}

	ch := make(chan *requests.Result)

	requestCount := 0

	for _, url := range urls {
		requestCount += 1
		go requests.RequestToChan(client, "GET", url, nil, ch)
	}

	for i := 0; i < requestCount; i++ {
		result := <-ch

		if result.Err != nil {
			fmt.Println(result.Err)
			continue
		}

		if verbose || result.Code < 200 || result.Code > 299 {
			fmt.Println(result.Url, result.Code)
		}
	}
}
