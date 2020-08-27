package typesense

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// SearchResponse is the default Typesense response for a serch.
type SearchResponse struct {
	FacetCounts []FacetCount      `json:"facet_counts"`
	Found       int               `json:"found"`
	Hits        []SearchResultHit `json:"hits"`
}

// FacetCount is the representation of a Typesense facet count.
type FacetCount struct {
	FieldName string `json:"field_name"`
	Counts    []struct {
		Count int    `json:"count"`
		Value string `json:"value"`
	} `json:"counts"`
}

// SearchResultHit represents a Typesense search result hit. Every
// retrieved document from a search will have the type map[string]interface{}.
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

// SearchOptions is all options that will be used to create
// a form url encoded to search in Typesense.
type SearchOptions struct {
	Query               string
	QueryBy             string
	FilterBy            string
	SortBy              string
	FacetBy             string
	MaxFacetValues      *int
	NumTypos            *int
	Prefix              bool
	Page                *int
	PerPage             *int
	IncludeFields       string
	ExcludeFields       string
	DropTokensThreshold *int
}

func (opts *SearchOptions) encodeForm() (string, error) {
	data := url.Values{}
	if opts.Query == "" {
		return "", ErrQueryRequired
	}
	data.Set("q", opts.Query)
	if opts.QueryBy == "" {
		return "", ErrQueryByRequired
	}
	data.Set("query_by", opts.QueryBy)
	opts.setOptionalFields(&data)
	return data.Encode(), nil
}

func (opts *SearchOptions) setOptionalFields(data *url.Values) {
	if opts.FilterBy != "" {
		data.Set("filter_by", opts.FilterBy)
	}
	if opts.SortBy != "" {
		data.Set("sort_by", opts.SortBy)
	}
	if opts.FacetBy != "" {
		data.Set("facet_by", opts.FacetBy)
	}
	if opts.MaxFacetValues != nil {
		data.Set("max_facet_values", strconv.Itoa(*opts.MaxFacetValues))
	}
	if opts.NumTypos != nil {
		data.Set("num_typos", strconv.Itoa(*opts.NumTypos))
	}
	if opts.Page != nil {
		data.Set("page", strconv.Itoa(*opts.Page))
	}
	if opts.PerPage != nil {
		data.Set("per_page", strconv.Itoa(*opts.PerPage))
	}
	if opts.IncludeFields != "" {
		data.Set("include_fields", opts.IncludeFields)
	}
	if opts.ExcludeFields != "" {
		data.Set("exclude_fields", opts.ExcludeFields)
	}
	if opts.DropTokensThreshold != nil {
		data.Set("drop_tokens_threshold", strconv.Itoa(*opts.DropTokensThreshold))
	}
	if opts.Prefix {
		data.Set("prefix", fmt.Sprintf("%v", opts.Prefix))
	}
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
	resp, err := c.apiCall(method, url, body)
	if err != nil {
		documentResponse.Error = err
		return &documentResponse
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		documentResponse.Error = ErrCollectionNotFound
		return &documentResponse
	} else if resp.StatusCode == http.StatusUnauthorized {
		documentResponse.Error = ErrUnauthorized
		return &documentResponse
	} else if resp.StatusCode == http.StatusConflict {
		documentResponse.Error = ErrDuplicateID
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
	resp, err := c.apiCall(method, url, nil)
	if err != nil {
		documentResponse.Error = err
		return &documentResponse
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		documentResponse.Error = ErrCollectionNotFound
		return &documentResponse
	} else if resp.StatusCode == http.StatusUnauthorized {
		documentResponse.Error = ErrUnauthorized
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
	resp, err := c.apiCall(method, url, nil)
	if err != nil {
		documentResponse.Error = err
		return &documentResponse
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		documentResponse.Error = ErrCollectionNotFound
		return &documentResponse
	} else if resp.StatusCode == http.StatusUnauthorized {
		documentResponse.Error = ErrUnauthorized
		return &documentResponse
	}
	documentResponse.Data, documentResponse.Error = ioutil.ReadAll(resp.Body)
	return &documentResponse
}

// Search searches for the query using the queryBy argument
// and other options in searchOptions in the Typesense API.
func (c *Client) Search(collectionName, query, queryBy string, searchOptions *SearchOptions) (*SearchResponse, error) {
	if searchOptions == nil {
		searchOptions = &SearchOptions{
			Query:   query,
			QueryBy: queryBy,
		}
	}
	urlEncodedForm, err := searchOptions.encodeForm()
	if err != nil {
		return nil, err
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
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	} else if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	} else if resp.StatusCode == http.StatusBadRequest {
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
