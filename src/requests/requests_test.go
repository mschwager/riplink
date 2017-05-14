package requests_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/mschwager/riplink/src/requests"
)

type MockClient struct {
	Body []byte
	Code int
	Err  error
}

func (client MockClient) Do(request *http.Request) (response *http.Response, err error) {
	response = &http.Response{
		StatusCode: client.Code,
		Body:       ioutil.NopCloser(bytes.NewReader(client.Body)),
	}

	return response, client.Err
}

func TestSendRequestBasic(t *testing.T) {
	body := []byte{}
	code := 200

	client := MockClient{
		Body: body,
		Code: code,
	}

	request, _ := http.NewRequest("UNUSED", "UNUSED", nil)

	result_body, result_code, result_err := requests.SendRequest(client, request)

	if len(result_body) != 0 || result_code != code || result_err != nil {
		t.Error("Failed to parse request: ", client)
	}
}

func TestSendRequestError(t *testing.T) {
	body := []byte{}
	code := 0
	err := errors.New("")

	client := MockClient{
		Body: body,
		Code: code,
		Err:  err,
	}

	request, _ := http.NewRequest("UNUSED", "UNUSED", nil)

	result_body, result_code, result_err := requests.SendRequest(client, request)

	if result_body != nil || result_code != code || result_err != err {
		t.Error("Failed to parse request: ", client)
	}
}

func TestRecursiveQueryToChanBasic(t *testing.T) {
	url := "https://example.com"
	body := []byte{}
	code := 200
	var depth uint = 0

	client := MockClient{
		Body: body,
		Code: code,
	}

	results := make(chan *requests.Result)

	go requests.RecursiveQueryToChan(client, url, depth, false, results)

	for result := range results {
		if result.Url != url || result.Code != code || result.Err != nil {
			t.Error("Failed to parse request: ", client)
		}
	}
}

func TestRecursiveQueryToChanError(t *testing.T) {
	url := "https://example.com"
	body := []byte{}
	code := 0
	var depth uint = 0
	err := errors.New("")

	client := MockClient{
		Body: body,
		Code: code,
		Err:  err,
	}

	results := make(chan *requests.Result)

	go requests.RecursiveQueryToChan(client, url, depth, false, results)

	for result := range results {
		if result.Url != url || result.Code != code || result.Err != err {
			t.Error("Failed to parse request: ", client)
		}
	}
}

func TestRecursiveQueryToChanRecurse(t *testing.T) {
	url := "https://example1.com"
	nestedUrl := "https://example2.com"
	body := []byte(fmt.Sprintf(`
	<html>
	<head>
	</head>
	<body>
		<a href="%s"></a>
	</body>
	</html>
	`, nestedUrl))
	code := 200
	var depth uint = 1

	client := MockClient{
		Body: body,
		Code: code,
	}

	results := make(chan *requests.Result)

	go requests.RecursiveQueryToChan(client, url, depth, false, results)

	for result := range results {
		if (result.Url != url && result.Url != nestedUrl) || result.Code != code || result.Err != nil {
			t.Error("Failed to parse request: ", client)
		}
	}
}

func TestRecursiveQueryToChanAttributeErrors(t *testing.T) {
	url := "https://example1.com"
	body := []byte(`
	<html>
	<head>
	</head>
	<body>
		<a href="mailto:test@example.com"></a>
		<a></a>
	</body>
	</html>
	`)
	code := 200
	var depth uint = 1

	client := MockClient{
		Body: body,
		Code: code,
	}

	results := make(chan *requests.Result)

	go requests.RecursiveQueryToChan(client, url, depth, false, results)

	for result := range results {
		if result.Url != url && result.Err != nil {
			t.Error("Failed to parse request: ", client)
		}
	}
}

func TestRecursiveQueryToChanAvoidDuplicateRequest(t *testing.T) {
	url := "https://example.com"
	body := []byte(`
	<html>
	<head>
	</head>
	<body>
		<a href="https://example.com"></a>
	</body>
	</html>
	`)
	code := 200
	var depth uint = 1

	client := MockClient{
		Body: body,
		Code: code,
	}

	results := make(chan *requests.Result)

	go requests.RecursiveQueryToChan(client, url, depth, false, results)

	count := 0
	for range results {
		count++
	}

	if count != 1 {
		t.Error("Failed to parse request: ", client)
	}
}

func TestRecursiveQueryToChanAvoidRelativeDuplicateRequest(t *testing.T) {
	url := "https://example.com"
	body := []byte(`
	<html>
	<head>
	</head>
	<body>
		<a href="https://example.com/relative/path"></a>
		<a href="/relative/path"></a>
	</body>
	</html>
	`)
	code := 200
	var depth uint = 1

	client := MockClient{
		Body: body,
		Code: code,
	}

	results := make(chan *requests.Result)

	go requests.RecursiveQueryToChan(client, url, depth, false, results)

	count := 0
	for range results {
		count++
	}

	if count != 2 {
		t.Error("Failed to parse request: ", client)
	}
}

func TestRecursiveQueryToChanAvoidDifferentDomain(t *testing.T) {
	url := "https://example.com"
	body := []byte(`
	<html>
	<head>
	</head>
	<body>
		<a href="https://different.com"></a>
	</body>
	</html>
	`)
	code := 200
	var depth uint = 1

	client := MockClient{
		Body: body,
		Code: code,
	}

	results := make(chan *requests.Result)

	go requests.RecursiveQueryToChan(client, url, depth, true, results)

	count := 0
	for range results {
		count++
	}

	if count != 1 {
		t.Error("Failed to parse request: ", client)
	}
}
