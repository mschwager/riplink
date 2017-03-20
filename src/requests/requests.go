package requests

import (
	"io/ioutil"
	"net/http"
	"sync"
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

func SendRequests(client Client, requests []*http.Request) (results []*Result, err error) {
	wg := sync.WaitGroup{}

	wg.Add(len(requests))

	for _, request := range requests {
		go func(innerRequest *http.Request) {
			defer wg.Done()

			_, code, err := SendRequest(client, innerRequest)

			result := &Result{
				Url:  innerRequest.URL.String(),
				Code: code,
				Err:  err,
			}

			results = append(results, result)
		}(request)
	}

	wg.Wait()

	return results, nil
}
