package requests

import (
	"io"
	"io/ioutil"
	"net/http"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

func Request(client Client, method string, url string, body io.Reader) (responseBody string, responseCode int, err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", 0, err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", 0, err
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", 0, err
	}

	return string(bytes), response.StatusCode, nil
}
