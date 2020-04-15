package typesense

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
