package typesense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type testDocumentStruct struct {
	Field1 string `json:"field1"`
	Field2 int32  `json:"field2"`
}

const (
	collectionNameTest = "test"

	searchResultTest = `
		{
			"facet_counts": [],
			"found": 62,
			"hits": [
				{
					"highlights": [
						{
							"field": "title",
							"snippet": "<mark>Harry</mark> <mark>Potter</mark> and the Philosopher's Stone"
						}
					],
					"document": {
						"authors": [
							"J.K. Rowling", "Mary GrandPré"
						],
						"authors_facet": [
							"J.K. Rowling", "Mary GrandPré"
						],
						"average_rating": 4.44,
						"id": "2",
						"image_url": "https://images.gr-assets.com/books/1474154022m/3.jpg",
						"publication_year": 1997,
						"publication_year_facet": "1997",
						"ratings_count": 4602479,
						"title": "Harry Potter and the Philosopher's Stone"
					}
				}
			]
		}
	`
)

var (
	testDocument = testDocumentStruct{
		Field1: "test",
		Field2: 10,
	}
)

func TestEncodeForm(t *testing.T) {
	opts := SearchOptions{
		Query:               "query",
		QueryBy:             "name",
		FilterBy:            "age>3",
		SortBy:              "age",
		FacetBy:             "tags",
		MaxFacetValues:      2,
		NumTypos:            2,
		Prefix:              true,
		Page:                1,
		PerPage:             5,
		IncludeFields:       "name",
		ExcludeFields:       "full_name",
		DropTokensThreshold: 5,
	}
	if _, err := opts.encodeForm(); err != nil {
		t.Errorf("Expected no errors, received %v", err)
	}
}

func TestIndexDocument(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		documentJSON, _ := json.Marshal(testDocument)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(documentJSON)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	documentResp := client.IndexDocument(collectionNameTest, testDocument)
	if documentResp.Error != nil {
		t.Errorf("Expected to receive no errors, received %v", documentResp.Error)
	}
}

func TestIndexDocument_collectionNotFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "collection not found"}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	documentResp := client.IndexDocument(collectionNameTest, testDocument)
	if documentResp.Error != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, documentResp.Error)
	}
}

func TestRetrieveDocument(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		documentJSON, _ := json.Marshal(testDocument)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(documentJSON)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	documentResp := client.RetrieveDocument(collectionNameTest, testDocument.Field1)
	if documentResp.Error != nil {
		t.Errorf("Expected to receive no errors, received %v", documentResp.Error)
	}
}

func TestRetrieveDocument_collectionNotFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "collection not found"}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	documentResp := client.RetrieveDocument(collectionNameTest, testDocument.Field1)
	if documentResp.Error != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, documentResp.Error)
	}
}

func TestDeleteDocument(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		documentJSON, _ := json.Marshal(testDocument)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(documentJSON)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	documentResp := client.DeleteDocument(collectionNameTest, testDocument.Field1)
	if documentResp.Error != nil {
		t.Errorf("Expected to receive no errors, received %v", documentResp.Error)
	}
}

func TestDeleteDocument_collectionNotFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "collection not found"}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	documentResp := client.DeleteDocument(collectionNameTest, testDocument.Field1)
	if documentResp.Error != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, documentResp.Error)
	}
}

func TestSearch(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(searchResultTest)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	searchResp, err := client.Search("books", "harry potter", "title", nil)
	if err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if len(searchResp.Hits) == 0 {
		t.Errorf("Expected to get at least one hit, got %d", len(searchResp.Hits))
	}
}

func TestSearch_collectionNotFound(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "collection not found"}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	_, err := client.Search("books", "harry potter", "title", nil)
	if err != ErrCollectionNotFound {
		t.Errorf("Expected to receive error %v, received %v", ErrCollectionNotFound, err)
	}
}

func TestSearch_missingRequiredField(t *testing.T) {
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "query_by is required"}`)),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	if _, err := client.Search("books", "harry potter", "", nil); err != ErrQueryByRequired {
		t.Errorf("Expected to receive error %v, received %v", ErrQueryByRequired, err)
	}

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "query_by is required"}`)),
		}, nil
	}
	if _, err := client.Search("books", "", "title", nil); err != ErrQueryRequired {
		t.Errorf("Expected to receive error %v, received %v", ErrQueryRequired, err)
	}
}

func TestSearch_badRequest(t *testing.T) {
	errorMessage := "sort_by must be of type number"
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"message": %q}`, errorMessage))),
		}, nil
	}
	client := Client{
		httpClient: mockClient,
		masterNode: testMasterNode,
	}
	if _, err := client.Search("books", "harry porter", "title", nil); err.Error() != errorMessage {
		t.Errorf("Expected to receive error %q, received %q", errorMessage, err.Error())
	}
}
