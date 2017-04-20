package requests

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/mschwager/riplink/src/parse"
	"github.com/mschwager/riplink/src/rpurl"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type Result struct {
	Url  string
	Code int
	Err  error
}

func SendRequest(client Client, request *http.Request) (responseBody []byte, responseCode int, err error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}

	return bytes, response.StatusCode, nil
}

func RecursiveQueryToChanHelper(client Client, queryUrl string, depth uint, results chan *Result, wg *sync.WaitGroup) {
	defer wg.Done()

	sendResult := func(url string, code int, err error) {
		result := &Result{
			Url:  url,
			Code: code,
			Err:  err,
		}
		results <- result
	}

	request, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		sendResult(queryUrl, 0, err)
		return
	}

	response, code, err := SendRequest(client, request)
	if err != nil {
		sendResult(queryUrl, code, err)
		return
	}

	sendResult(queryUrl, code, nil)

	if depth <= 0 {
		return
	}

	node, err := parse.BytesToHtmlNode(response)
	if err != nil {
		sendResult(queryUrl, code, err)
		return
	}

	anchors := parse.Anchors(node)

	hrefs, errs := parse.ValidHrefs(anchors)
	for _, err := range errs {
		sendResult(queryUrl, code, err)
	}

	urls, errs := rpurl.AbsoluteHttpUrls(queryUrl, hrefs)
	for _, err := range errs {
		sendResult(queryUrl, code, err)
	}

	for _, url := range urls {
		wg.Add(1)
		go RecursiveQueryToChanHelper(client, url, depth-1, results, wg)
	}

}

func RecursiveQueryToChan(client Client, queryUrl string, depth uint, results chan *Result) {
	defer close(results)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go RecursiveQueryToChanHelper(client, queryUrl, depth, results, &wg)

	wg.Wait()
}
