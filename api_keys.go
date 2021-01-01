package typesense

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	// ActionDocumentsSearch allows only search requests.
	ActionDocumentsSearch = "documents:search"

	// ActionDocumentsGet allows fetching a single document.
	ActionDocumentsGet = "documents:get"

	// ActionCollectionsDelete allows a collection to be deleted.
	ActionCollectionsDelete = "collections:delete"

	// ActionCollectionsCreate allows a collection to be created.
	ActionCollectionsCreate = "collections:create"

	// ActionCollectionsAll allow all kinds of collection related operations.
	ActionCollectionsAll = "collections:*"

	// ActionAll Allows all operations.
	ActionAll = "*"
)

// APIKey is the model for a Typesense API key.
type APIKey struct {
	ID          int       `json:"id"`
	Value       string    `json:"value"`
	Description string    `json:"description,omitempty"`
	Actions     []string  `json:"actions"`
	Collections []string  `json:"collections"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// CreateAPIKey creates a new API key for the Typesense API.
func (c *Client) CreateAPIKey(key APIKey) (*APIKey, error) {
	method := http.MethodPost
	url := fmt.Sprintf(
		"%s://%s:%s/keys",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
	)
	body, err := json.Marshal(key)
	if err != nil {
		return nil, err
	}
	resp, err := c.apiCall(method, url, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("creation of api key returned status code: %d", resp.StatusCode)
	}
	var apiKey APIKey
	if err := json.NewDecoder(resp.Body).Decode(&apiKey); err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// GetAPIKey retrieves a Typesense API key by its id.
func (c *Client) GetAPIKey(id int) (*APIKey, error) {
	method := http.MethodGet
	url := fmt.Sprintf(
		"%s://%s:%s/keys/%d",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		id,
	)
	resp, err := c.apiCall(method, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(
			"retrieve of api key %d failed with status code: %d",
			id,
			resp.StatusCode,
		)
	}
	var key APIKey
	if err := json.NewDecoder(resp.Body).Decode(&key); err != nil {
		return nil, err
	}
	return &key, nil
}

// GetAPIKeys retrieve all API keys.
func (c *Client) GetAPIKeys() ([]*APIKey, error) {
	method := http.MethodGet
	url := fmt.Sprintf(
		"%s://%s:%s/keys",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
	)
	resp, err := c.apiCall(method, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"retrieve of api keys failed with status code: %d",
			resp.StatusCode,
		)
	}
	var keysWrapper struct {
		Keys []*APIKey `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&keysWrapper); err != nil {
		return nil, err
	}
	return keysWrapper.Keys, nil
}

// DeleteAPIKey deletes an API key by its id.
func (c *Client) DeleteAPIKey(id int) error {
	method := http.MethodDelete
	url := fmt.Sprintf(
		"%s://%s:%s/keys/%d",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		id,
	)
	resp, err := c.apiCall(method, url, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"delete api key %d failed with status code: %d",
			id,
			resp.StatusCode,
		)
	}
	return nil
}
