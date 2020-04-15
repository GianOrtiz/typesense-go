package typesense

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/GianOrtiz/typesense-go/mock"
)

var (
	testMasterNode = &Node{
		Host:     "localhost",
		Port:     "8108",
		Protocol: "http",
		APIKey:   "secret",
	}

	mockClient = mock.HTTPClient{}
)

func TestNewClient(t *testing.T) {
	timeout := 2
	client := NewClient(testMasterNode, timeout)
	if client == nil {
		t.Errorf("Expected to receive a configured client, received <nil>")
	}
}

func TestPing_ready(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"ok": true}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	if err := client.Ping(); err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
}

func TestPing_notReady(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"ok": false}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	if err := client.Ping(); err != ErrConnNotReady {
		t.Errorf("Expected error %v, received %v", ErrConnNotReady, err)
	}
}

func TestDebugInfo(t *testing.T) {
	defaultVersion := "0.11"
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"version": %q}`, defaultVersion))),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	version, err := client.DebugInfo()
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if version != defaultVersion {
		t.Errorf("Expected to receive value %v, received %v", defaultVersion, version)
	}
}
