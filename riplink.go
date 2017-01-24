package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mschwager/riplink/src/parse"
	"github.com/mschwager/riplink/src/requests"
	"github.com/mschwager/riplink/src/rpurl"
)

func main() {
	var url string
	flag.StringVar(&url, "url", "https://google.com", "URL to query")

	flag.Parse()

	client := &http.Client{
		Timeout: time.Second * 5,
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

	ch := make(chan *requests.Result)

	for _, anchor := range anchors {
		href, err := parse.Href(anchor)
		if err != nil {
			fmt.Println(err)
			continue
		}

		pageUrl := href.Val
		hasHost, err := rpurl.HasHost(pageUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !hasHost {
			pageUrl, err = rpurl.AddBaseHost(url, pageUrl)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		go requests.RequestToChan(client, "GET", pageUrl, nil, ch)
	}

	for response := range ch {
		if response.Err != nil {
			fmt.Println(response.Err)
			continue
		}

		fmt.Println("HREF: " + response.Url)
		fmt.Println("STATUS CODE: " + strconv.Itoa(response.Code))
	}
}
