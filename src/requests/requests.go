package requests

import (
	"io"
	"io/ioutil"
	"net/http"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type Result struct {
	Url  string
	Code int
	Err  error
}

func Request(client Client, method string, url string, body io.Reader) (responseBody []byte, responseCode int, err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return []byte{}, 0, err
	}

	response, err := client.Do(request)
	if err != nil {
		return []byte{}, 0, err
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, 0, err
	}

	return bytes, response.StatusCode, nil
}

func RequestToChan(client Client, method string, url string, body io.Reader, ch chan *Result) (err error) {
	_, code, err := Request(client, method, url, body)

	result := &Result{
		Url:  url,
		Code: code,
		Err:  err,
	}

	ch <- result

	return nil
}
