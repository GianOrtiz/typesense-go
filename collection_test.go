package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestRetrieveCollection(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		collectionJSON, _ := json.Marshal(&testCollection)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(bytes.NewReader(collectionJSON)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	collection, err := client.RetrieveCollection(testCollection.Name)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if !reflect.DeepEqual(*collection, testCollection) {
		t.Errorf("Expected to receive %v, received %v", testCollection, *collection)
	}
}

func TestRetrieveCollection_notFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body: ioutil.NopCloser(strings.NewReader(`{"json": "collection not found"}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.RetrieveCollection(testCollection.Name)
	if err != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, err)
	}
}

func TestDeleteCollection(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		collectionJSON, _ := json.Marshal(&testCollection)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(bytes.NewReader(collectionJSON)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	collection, err := client.DeleteCollection(testCollection.Name)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if !reflect.DeepEqual(*collection, testCollection) {
		t.Errorf("Expected to receive %v, received %v", testCollection, *collection)
	}
}

func TestDeleteCollection_notFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body: ioutil.NopCloser(strings.NewReader(`{"json": "collection not found"}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.DeleteCollection(testCollection.Name)
	if err != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, err)
	}
}