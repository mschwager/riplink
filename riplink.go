package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mschwager/riplink/src/parse"
	"github.com/mschwager/riplink/src/requests"
)

func main() {
	var url string
	flag.StringVar(&url, "url", "https://google.com", "URL to query")

	flag.Parse()

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	response, _, err := requests.Request(client, "GET", url, nil)
	if err != nil {
		panic(err)
	}

	node, err := parse.StringToHtmlNode(response)
	if err != nil {
		panic(err)
	}

	anchors, err := parse.Anchors(node)
	if err != nil {
		panic(err)
	}

	for _, anchor := range anchors {
		href, err := parse.Href(anchor)
		if err != nil {
			fmt.Println(err)
			continue
		}

		url = href.Val
		_, statusCode, err := requests.Request(client, "HEAD", url, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("HREF: " + href.Val)
		fmt.Println("STATUS CODE: " + strconv.Itoa(statusCode))
	}
}
