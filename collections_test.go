package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var (
	collectionCfg = CollectionConfig{
		Name: "companies",
		Fields: []CollectionField{
			CollectionField{
				Name: "name",
				Type: "string",
			},
		},
	}
	collection = Collection{
		collectionCfg,
		0,
		0,
	}
)

func TestCreateCollection(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		collectionData, _ := json.Marshal(collection)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(collectionData)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	collectionResp, err := client.CreateCollection(collectionCfg)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if collectionResp == nil {
		t.Errorf("Expected to receive a collection as the response, received %v", collectionResp)
	}
}

func TestCreateCollection_conflict(t *testing.T) {
	errorMessage := "collection companies already exist"
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusConflict,
			Body:       ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"message": %q}`, errorMessage))),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.CreateCollection(collectionCfg)
	if err == nil || err.Error() != errorMessage {
		t.Errorf("Expected to receive error message %q, received error %v", errorMessage, err)
	}
}

func TestRetrieveCollections(t *testing.T) {
	jsonBody := `[{"name": "companies", "num_documents": 0, "fields": [{"name": "name", "type": "string", "facet": false}]}]`
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(jsonBody)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	collections, err := client.RetrieveCollections()
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if collections == nil {
		t.Errorf("Expected to receive collections, received nil")
	}
}
