package mock

import (
	"net/http"
)

// HTTPClient is the mocks for the typesense http client.
type HTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do implements the typesense http client interface.
func (c HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.DoFunc(req)
}
