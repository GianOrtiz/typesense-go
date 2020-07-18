package typesense

import (
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
	if collectionSchema.Name == "" {
		return nil, ErrCollectionNameRequired
	} else if len(collectionSchema.Fields) == 0 {
		return nil, ErrCollectionFieldsRequired
	}
	method := http.MethodPost
	url := fmt.Sprintf(
		"%s://%s:%s/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		collectionsEndpoint,
	)
	collectionJSON, _ := json.Marshal(collectionSchema)
	resp, err := c.apiCall(method, url, collectionJSON)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusConflict {
		return nil, ErrCollectionDuplicate
	} else if resp.StatusCode == http.StatusBadRequest {
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
	resp, err := c.apiCall(method, url, nil)
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

// RetrieveCollection retrieves a single collection by
// its name.
func (c *Client) RetrieveCollection(collectionName string) (*Collection, error) {
	method := http.MethodGet
	url := fmt.Sprintf(
		"%s://%s:%s/%s/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		collectionsEndpoint,
		collectionName,
	)
	resp, err := c.apiCall(method, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrCollectionNotFound
	}
	var collection Collection
	if err := json.NewDecoder(resp.Body).Decode(&collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

// DeleteCollection deletes a collection by its name.
func (c *Client) DeleteCollection(collectionName string) (*Collection, error) {
	method := http.MethodDelete
	url := fmt.Sprintf(
		"%s://%s:%s/%s/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		collectionsEndpoint,
		collectionName,
	)
	resp, err := c.apiCall(method, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrCollectionNotFound
	}
	var collection Collection
	if err := json.NewDecoder(resp.Body).Decode(&collection); err != nil {
		return nil, err
	}
	return &collection, nil
}
