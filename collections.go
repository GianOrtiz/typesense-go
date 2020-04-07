package typesense

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const collectionsEndpoint = "collections"

// CollectionSchema is the definition of a collection schema
// to create in Typesense.
type CollectionSchema struct {
	Name                string            `json:"name"`
	Fields              []CollectionField `json:"fields"`
	DefaultSortingField string            `json:"default_sorting_field"`
}

// Collection is the model of a collection created in the
// Typesense API.
type Collection struct {
	CollectionSchema
	NumDocuments int   `json:"num_documents"`
	CreatedAt    int64 `json:"created_at"`
}

// CollectionField is a Typesense collection field.
type CollectionField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Facet bool   `json:"facet"`
}

// CreateCollection creates a new collection using the
// given collection schema.
func (c *Client) CreateCollection(collectionSchema CollectionSchema) (*Collection, error) {
	method := http.MethodPost
	url := fmt.Sprintf(
		"%s://%s:%s/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		collectionsEndpoint,
	)
	collectionJSON, _ := json.Marshal(collectionSchema)
	req, _ := http.NewRequest(method, url, bytes.NewReader(collectionJSON))
	req.Header.Add(defaultHeaderKey, c.masterNode.APIKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusConflict || resp.StatusCode == http.StatusBadRequest {
		var apiResponse APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(apiResponse.Message)
	}
	var collectionResponse Collection
	if err := json.NewDecoder(resp.Body).Decode(&collectionResponse); err != nil {
		return nil, err
	}
	return &collectionResponse, nil
}

// RetrieveCollections get all collections from Typesense.
func (c *Client) RetrieveCollections() ([]*Collection, error) {
	method := http.MethodGet
	url := fmt.Sprintf(
		"%s://%s:%s/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		collectionsEndpoint,
	)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add(defaultHeaderKey, c.masterNode.APIKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var collections []*Collection
	if err := json.NewDecoder(resp.Body).Decode(&collections); err != nil {
		return nil, err
	}
	return collections, nil
}
