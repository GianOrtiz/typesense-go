package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const collectionsEndpoint = "collections"

// CollectionConfig is the collection model of typesense.
type CollectionConfig struct {
	Name                string            `json:"name"`
	Fields              []CollectionField `json:"fields"`
	DefaultSortingField string            `json:"default_sorting_field"`
}

// Collection is the model of a collection in typesense.
type Collection struct {
	CollectionConfig
	NumDocuments int   `json:"num_documents"`
	CreatedAt    int64 `json:"created_at"`
}

// CollectionField is the representation of a typesense
// collection field.
type CollectionField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Facet bool   `json:"facet"`
}

// CreateCollection creates a new collection.
func (c *Client) CreateCollection(collectionCfg CollectionConfig) (*Collection, error) {
	method := http.MethodPost
	url := fmt.Sprintf(
		"%s://%s:%s/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		collectionsEndpoint,
	)
	collectionJSON, _ := json.Marshal(collectionCfg)
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

// RetrieveCollections get all collections in typesense.
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
