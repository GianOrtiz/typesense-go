package typesense

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// SearchResponse is the response for a search.
type SearchResponse struct {
	FacetCounts []FacetCount      `json:"facet_counts"`
	Found       int               `json:"found"`
	Hits        []SearchResultHit `json:"hits"`
}

// FacetCount is the representation of a typesense facet count.
type FacetCount struct {
	FieldName string `json:"field_name"`
	Counts    []struct {
		Count int    `json:"count"`
		Value string `json:"value"`
	} `json:"counts"`
}

// SearchResultHit represents a typesense search result hit.
type SearchResultHit struct {
	Highlights []SearchHighlight      `json:"highlights"`
	Document   map[string]interface{} `json:"document"`
}

// SearchHighlight represents the highlight of texts in the
// search result.
type SearchHighlight struct {
	Field    string   `json:"field"`
	Snippet  string   `json:"snippet"`
	Snippets []string `json:"snippets"`
	Indices  []int    `json:"indices"`
}

// SearchOptions is the options to make a search.
type SearchOptions struct {
	Query               string
	QueryBy             string
	FilterBy            string
	SortBy              string
	FacetBy             string
	MaxFacetValues      int
	NumTypos            int
	Prefix              bool
	Page                int
	PerPage             int
	IncludeFields       string
	ExcludeFields       string
	DropTokensThreshold int
}

func (opts *SearchOptions) encodeForm() (string, error) {
	data := url.Values{}
	if opts.Query == "" {
		return "", fmt.Errorf("Query is a required field")
	}
	data.Set("q", opts.Query)
	if opts.QueryBy == "" {
		return "", fmt.Errorf("QueryBy is a required field")
	}
	data.Set("query_by", opts.QueryBy)
	if opts.FilterBy != "" {
		data.Set("filter_by", opts.FilterBy)
	}
	if opts.SortBy != "" {
		data.Set("sort_by", opts.SortBy)
	}
	if opts.FacetBy != "" {
		data.Set("facet_by", opts.FacetBy)
	}
	if opts.MaxFacetValues > 0 {
		data.Set("max_facet_values", strconv.Itoa(opts.MaxFacetValues))
	}
	if opts.NumTypos >= 0 {
		data.Set("num_typos", strconv.Itoa(opts.NumTypos))
	}
	if opts.Page >= 0 {
		data.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		data.Set("per_page", strconv.Itoa(opts.PerPage))
	}
	if opts.IncludeFields != "" {
		data.Set("include_fields", opts.IncludeFields)
	}
	if opts.ExcludeFields != "" {
		data.Set("exclude_fields", opts.ExcludeFields)
	}
	if opts.DropTokensThreshold >= 0 {
		data.Set("drop_tokens_threshold", strconv.Itoa(opts.DropTokensThreshold))
	}
	data.Set("prefix", fmt.Sprintf("%v", opts.Prefix))
	return data.Encode(), nil
}

// IndexDocument index a new document in the collection.
func (c *Client) IndexDocument(collectionName string, document interface{}) *DocumentResponse {
	documentResponse := DocumentResponse{}
	method := http.MethodPost
	url := fmt.Sprintf(
		"%s://%s:%s/%s/%s/documents",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		collectionsEndpoint,
		collectionName,
	)
	body, _ := json.Marshal(document)
	req, _ := http.NewRequest(method, url, bytes.NewReader(body))
	req.Header.Add(defaultHeaderKey, c.masterNode.APIKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		documentResponse.Error = err
		return &documentResponse
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		documentResponse.Error = ErrCollectionNotFound
		return &documentResponse
	}
	documentResponse.Data, documentResponse.Error = ioutil.ReadAll(resp.Body)
	return &documentResponse
}

// RetrieveDocument retrieves a document in the collection by its id.
func (c *Client) RetrieveDocument(collectionName, documentID string) *DocumentResponse {
	documentResponse := DocumentResponse{}
	method := http.MethodGet
	url := fmt.Sprintf(
		"%s://%s:%s/%s/%s/documents/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		collectionsEndpoint,
		collectionName,
		documentID,
	)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add(defaultHeaderKey, c.masterNode.APIKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		documentResponse.Error = err
		return &documentResponse
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		documentResponse.Error = ErrCollectionNotFound
		return &documentResponse
	}
	documentResponse.Data, documentResponse.Error = ioutil.ReadAll(resp.Body)
	return &documentResponse
}

// DeleteDocument deletes a document in the collection by its id.
func (c *Client) DeleteDocument(collectionName, documentID string) *DocumentResponse {
	documentResponse := DocumentResponse{}
	method := http.MethodDelete
	url := fmt.Sprintf(
		"%s://%s:%s/%s/%s/documents/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		collectionsEndpoint,
		collectionName,
		documentID,
	)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add(defaultHeaderKey, c.masterNode.APIKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		documentResponse.Error = err
		return &documentResponse
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		documentResponse.Error = ErrCollectionNotFound
		return &documentResponse
	}
	documentResponse.Data, documentResponse.Error = ioutil.ReadAll(resp.Body)
	return &documentResponse
}

// Search searches for documents using the search options.
func (c *Client) Search(collectionName, query, queryBy string, searchOptions *SearchOptions) (*SearchResponse, error) {
	urlEncodedForm := fmt.Sprintf("q=%s&query_by=%s", query, queryBy)
	var err error
	if searchOptions != nil {
		searchOptions.Query = query
		searchOptions.QueryBy = queryBy
		urlEncodedForm, err = searchOptions.encodeForm()
		if err != nil {
			return nil, err
		}
	}

	method := http.MethodGet
	url := fmt.Sprintf(
		"%s://%s:%s/%s/%s/documents/search?%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		collectionsEndpoint,
		collectionName,
		urlEncodedForm,
	)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add(defaultHeaderKey, c.masterNode.APIKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, ErrCollectionNotFound
	}
	if resp.StatusCode == 400 {
		var apiResponse APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(apiResponse.Message)
	}
	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}
	return &searchResponse, nil
}
