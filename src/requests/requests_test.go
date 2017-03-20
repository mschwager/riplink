package requests_test

import (
	"bytes"
	"errors"
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

func TestSendRequestsBasic(t *testing.T) {
	body := []byte{}
	code := 200
	url := "URL"

	client := MockClient{
		Body: body,
		Code: code,
	}

	var preparedRequests []*http.Request

	request1, _ := http.NewRequest("UNUSED", url, nil)
	preparedRequests = append(preparedRequests, request1)

	request2, _ := http.NewRequest("UNUSED", url, nil)
	preparedRequests = append(preparedRequests, request2)

	results, err := requests.SendRequests(client, preparedRequests)

	for _, result := range results {
		if result.Url != url || result.Code != code || err != nil {
			t.Error("Failed to parse request: ", client)
		}
	}
}

func TestSendRequestsError(t *testing.T) {
	body := []byte{}
	code := 0
	url := "URL"
	err := errors.New("")

	client := MockClient{
		Body: body,
		Code: code,
		Err:  err,
	}

	var preparedRequests []*http.Request

	request1, _ := http.NewRequest("UNUSED", url, nil)
	preparedRequests = append(preparedRequests, request1)

	request2, _ := http.NewRequest("UNUSED", url, nil)
	preparedRequests = append(preparedRequests, request2)

	results, _ := requests.SendRequests(client, preparedRequests)

	for _, result := range results {
		if result.Url != url || result.Code != code || result.Err != err {
			t.Error("Failed to parse request: ", client)
		}
	}
}
