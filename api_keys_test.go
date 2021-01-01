package typesense

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var testAPIKey = APIKey{
	ID:          1,
	Value:       "some-value",
	Description: "description",
	Actions:     []string{ActionAll},
	Collections: []string{"companies"},
	ExpiresAt:   time.Now().Add(time.Duration(time.Hour * 5)),
}

func TestCreateAPIKey(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(testAPIKey)
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.CreateAPIKey(testAPIKey)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
}

func TestGetAPIKey(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(testAPIKey)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.GetAPIKey(testAPIKey.ID)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
}

func TestGetAPIKey_notFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(testAPIKey)
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.GetAPIKey(testAPIKey.ID)
	if err != ErrNotFound {
		t.Errorf("Expected to receive error: %v, received %v", ErrNotFound, err)
	}
}

func TestGetAPIKeys(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(map[string]interface{}{
			"keys": []*APIKey{&testAPIKey},
		})
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	keys, err := client.GetAPIKeys()
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if len(keys) == 0 {
		t.Errorf("Expected to receive at least one API key, received zero")
	}
}

func TestDeleteAPIKey(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	err := client.DeleteAPIKey(1)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
}
