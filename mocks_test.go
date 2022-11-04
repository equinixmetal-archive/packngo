package packngo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// MockClient makes it simpler to test the Client
type MockClient struct {
	fnNewRequest          func(method, path string, body interface{}) (*http.Request, error)
	fnDo                  func(req *http.Request, v interface{}) (*Response, error)
	fnDoRequest           func(method, path string, body, v interface{}) (*Response, error)
	fnDoRequestWithHeader func(method string, headers map[string]string, path string, body, v interface{}) (*Response, error)
}

var _ requestDoer = &MockClient{}

// NewRequest uses the mock NewRequest function
func (mc *MockClient) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	return mc.fnNewRequest(method, path, body)
}

// Do uses the mock Do function
func (mc *MockClient) Do(req *http.Request, v interface{}) (*Response, error) {
	return mc.fnDo(req, v)
}

// DoRequest uses the mock DoRequest function
func (mc *MockClient) DoRequest(method, path string, body, v interface{}) (*Response, error) {
	return mc.fnDoRequest(method, path, body, v)
}

// DoRequestWithHeader uses the mock DoRequestWithHeader function
func (mc *MockClient) DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*Response, error) {
	return mc.fnDoRequestWithHeader(method, headers, path, body, v)
}

func mockResponse(code int, body string, req *http.Request) *Response {
	return &Response{Response: &http.Response{
		Status:        fmt.Sprintf("%d Ignored", code),
		StatusCode:    code,
		Body:          ioutil.NopCloser(bytes.NewReader([]byte(body))), // TODO: io.NopCloser requires go 1.16+
		ContentLength: int64(len(body)),
		Request:       req,
	}}
}
