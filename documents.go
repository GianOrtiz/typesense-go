package typesense

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
// a form url encoded to search in Typesense. More information
// about the values can be found at https://typesense.org/docs/0.14.0/api/#search-collection.
type SearchOptions struct {
	// Query text to search for.
	Query string `url:"q"`

	// QueryBy represents fields to query_by.
	QueryBy []string

	// MaxHits is the max number of hits for the query search, value
	// increase may increase latency. Default value is 500.
	MaxHits *int

	// Prefix whether the query should be treated as a prefix or not.
	Prefix *bool

	// FilterBy represents filter conditions for refining your search
	// results.
	FilterBy []string

	// SortBy list of numerical values and their corresponding sort order
	// to sort results by.
	SortBy []string

	// FacetBy list of fields that will be used for faceting your results on.
	FacetBy []string

	// MaxFacetValues maximum number of facet values to be returned.
	MaxFacetValues *int

	// FacetQuery filter facet values by this paremeter, only values that
	// match the facet value will be matched.
	FacetQuery *string

	// NumTypos number of typographical errors (1 or 2) that would be
	// tolerated. Default value is 2.
	NumTypos *int

	// Page results from this specific page number would be fetched.
	Page *int

	// PerPage number of results to fetch per page. Default value is 10.
	PerPage *int

	// GroupBy aggregate search results by groups, groups must be a
	// facet field.
	GroupBy []string

	// GroupLimit maximum number of hits to return for every group
	// Default value is 3.
	GroupLimit *int

	// IncludeFields list of fields from the document to include in the search result.
	IncludeFields []string

	// ExcludeFields list of fields from the document to exclude in the search result.
	ExcludeFields []string

	// HighlightFullFields list of fields which should be highlighted fully without snippeting.
	// Default is all fields will be snipped.
	HighlightFullFields []string

	// SnippetThreshold Field values under this length will be fully highlighted, instead
	// of showing a snippet of relevant portion.
	// Default value is 30.
	SnippetThreshold *int

	// DropTokensThreshold if the number of hits is less than this value, Typesense
	// will try to drop tokens until the number of hits get to this value.
	// Default value is 10.
	DropTokensThreshold *int

	// TypoTokensThreshold if the number of results found for a specific query is less
	// than this number, Typesense will attempt to look for tokens with more typos
	// until enough results are found.
	// Default value is 100.
	TypoTokensThreshold *int

	// PinnedHits list of records to unconditionally include in the search results at
	// specific positions.
	PinnedHits []string

	// HiddenHits list of records to unconditionally hide from search results.
	Hiddenhits []string
}

func (opts *SearchOptions) encodeForm() (string, error) {
	data := url.Values{}
	if opts.Query == "" {
		return "", ErrQueryRequired
	}
	data.Set("q", opts.Query)
	if opts.QueryBy == nil || len(opts.QueryBy) == 0 {
		return "", ErrQueryByRequired
	}
	queryBy := strings.Join(opts.QueryBy, ",")
	data.Set("query_by", queryBy)
	opts.setOptionalFields(&data)
	return data.Encode(), nil
}

func (opts *SearchOptions) setOptionalFields(data *url.Values) {
	if opts.MaxHits != nil {
		data.Set("max_hits", strconv.Itoa(*opts.MaxHits))
	}
	if opts.Prefix != nil {
		data.Set("prefix", strconv.FormatBool(*opts.Prefix))
	}
	if opts.FilterBy != nil && len(opts.FilterBy) > 0 {
		filterBy := strings.Join(opts.FilterBy, " && ")
		data.Set("filter_by", filterBy)
	}
	if opts.SortBy != nil && len(opts.SortBy) > 0 {
		sortBy := strings.Join(opts.SortBy, ",")
		data.Set("sort_by", sortBy)
	}
	if opts.FacetBy != nil && len(opts.FacetBy) > 0 {
		facetBy := strings.Join(opts.FacetBy, ",")
		data.Set("facet_by", facetBy)
	}
	if opts.MaxFacetValues != nil {
		data.Set("max_facet_values", strconv.Itoa(*opts.MaxFacetValues))
	}
	if opts.FacetQuery != nil {
		data.Set("facet_query", *opts.FacetQuery)
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
	if opts.GroupBy != nil && len(opts.GroupBy) > 0 {
		groupBy := strings.Join(opts.GroupBy, ",")
		data.Set("group_by", groupBy)
	}
	if opts.GroupLimit != nil {
		data.Set("group_limit", strconv.Itoa(*opts.GroupLimit))
	}
	if opts.IncludeFields != nil && len(opts.IncludeFields) > 0 {
		includeFields := strings.Join(opts.IncludeFields, ",")
		data.Set("include_fields", includeFields)
	}
	if opts.ExcludeFields != nil && len(opts.ExcludeFields) > 0 {
		excludeFields := strings.Join(opts.ExcludeFields, ",")
		data.Set("exclude_fields", excludeFields)
	}
	if opts.HighlightFullFields != nil && len(opts.HighlightFullFields) > 0 {
		highlightFullFields := strings.Join(opts.HighlightFullFields, ",")
		data.Set("highlight_full_fields", highlightFullFields)
	}
	if opts.SnippetThreshold != nil {
		data.Set("snippet_threshold", strconv.Itoa(*opts.SnippetThreshold))
	}
	if opts.DropTokensThreshold != nil {
		data.Set("drop_tokens_threshold", strconv.Itoa(*opts.DropTokensThreshold))
	}
	if opts.TypoTokensThreshold != nil {
		data.Set("typo_tokens_threshold", strconv.Itoa(*opts.TypoTokensThreshold))
	}
	if opts.PinnedHits != nil && len(opts.PinnedHits) > 0 {
		pinnedHits := strings.Join(opts.PinnedHits, ",")
		data.Set("pinned_hits", pinnedHits)
	}
	if opts.Hiddenhits != nil && len(opts.Hiddenhits) > 0 {
		hiddenhits := strings.Join(opts.Hiddenhits, ",")
		data.Set("hidden_hits", hiddenhits)
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
	} else if resp.StatusCode == http.StatusBadRequest {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			apiErr.Message = "status bad request"
		}
		documentResponse.Error = apiErr
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
func (c *Client) Search(collectionName, query string, queryBy []string, searchOptions *SearchOptions) (*SearchResponse, error) {
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
