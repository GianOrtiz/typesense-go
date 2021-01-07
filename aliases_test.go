package typesense

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	testAlias = Alias{
		Name:           "alias",
		CollectionName: "collection",
	}
)

func TestCreateAlias(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(testAlias)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.CreateAlias(testAlias.Name, &testAlias)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
}

func TestRetrieveAlias(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(testAlias)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	alias, err := client.RetrieveAlias(testAlias.Name)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if alias.Name != testAlias.Name || alias.CollectionName != testAlias.CollectionName {
		t.Errorf(
			"Expected to receive alias %+v, received %+v",
			testAlias,
			alias,
		)
	}
}

func TestRetrieveAliases(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(map[string]interface{}{
			"aliases": []*Alias{&testAlias},
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
	aliases, err := client.RetrieveAliases()
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if len(aliases) != 1 {
		t.Errorf("Expected to receive exact one alias, received: %d", len(aliases))
	}
	alias := aliases[0]
	if alias.Name != testAlias.Name || alias.CollectionName != testAlias.CollectionName {
		t.Errorf(
			"Expected to receive first alias %+v, received %+v",
			testAlias,
			alias,
		)
	}
}

func TestDeleteAlias(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(testAlias)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.DeleteAlias(testAlias.Name)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
}
