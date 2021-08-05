package packngo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

/* deadcode, for now
func mockDoRequestWithHeader(doFn func(req *http.Request, v interface{}) (*Response, error), newRequestFn func(method, path string, body interface{}) (*http.Request, error)) func(string, map[string]string, string, interface{}, interface{}) (*Response, error) {
	return func(method string, headers map[string]string, path string, body, v interface{}) (*Response, error) {
		req, err := newRequestFn(method, path, body)
		for k, v := range headers {
			req.Header.Add(k, v)
		}

		if err != nil {
			return nil, err
		}
		return doFn(req, v)
	}
}
*/

func mockNewRequest() func(string, string, interface{}) (*http.Request, error) {
	baseURL := &url.URL{}
	apiKey, consumerToken, userAgent := "", "", ""
	return func(method, path string, body interface{}) (*http.Request, error) {
		// relative path to append to the endpoint url, no leading slash please
		if path[0] == '/' {
			path = path[1:]
		}
		rel, err := url.Parse(path)
		if err != nil {
			return nil, err
		}

		u := baseURL.ResolveReference(rel)

		// json encode the request body, if any
		buf := new(bytes.Buffer)
		if body != nil {
			err := json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err := http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}

		req.Close = true

		req.Header.Add("X-Auth-Token", apiKey)
		req.Header.Add("X-Consumer-Token", consumerToken)

		req.Header.Add("Content-Type", mediaType)
		req.Header.Add("Accept", mediaType)
		req.Header.Add("User-Agent", userAgent)
		return req, nil
	}
}

func mockDoRequest(newRequestFn func(string, string, interface{}) (*http.Request, error), doFn func(*http.Request, interface{}) (*Response, error)) func(method, path string, body, v interface{}) (*Response, error) {
	return func(method, path string, body, v interface{}) (*Response, error) {
		req, err := newRequestFn(method, path, body)
		if err != nil {
			return nil, err
		}
		return doFn(req, v)
	}
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
