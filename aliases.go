package typesense

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Alias is the representation of a collection alias.
type Alias struct {
	Name           string `json:"name,omitempty"`
	CollectionName string `json:"collection_name"`
}

// CreateAlias creates a new collection alias for a collection or updates
// the alias if it already exists.
func (c *Client) CreateAlias(aliasName string, alias *Alias) (*Alias, error) {
	method := http.MethodPut
	url := fmt.Sprintf(
		"%s://%s:%s/aliases/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		aliasName,
	)
	body, _ := json.Marshal(alias)
	res, err := c.apiCall(method, url, body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		responseBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, HTTPError{
			Status:       res.StatusCode,
			ResponseBody: responseBody,
		}
	}
	var upsertAlias Alias
	if err := json.NewDecoder(res.Body).Decode(&upsertAlias); err != nil {
		return nil, err
	}
	return &upsertAlias, nil
}

// RetrievesAlias retrieves an alias metadata by its name.
func (c *Client) RetrieveAlias(aliasName string) (*Alias, error) {
	method := http.MethodGet
	url := fmt.Sprintf(
		"%s://%s:%s/aliases/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		aliasName,
	)
	res, err := c.apiCall(method, url, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		responseBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, HTTPError{
			Status:       res.StatusCode,
			ResponseBody: responseBody,
		}
	}
	var alias Alias
	if err := json.NewDecoder(res.Body).Decode(&alias); err != nil {
		return nil, err
	}
	return &alias, nil
}

// RetrieveAliases retrieve all aliases in Typesense.
func (c *Client) RetrieveAliases() ([]*Alias, error) {
	method := http.MethodGet
	url := fmt.Sprintf(
		"%s://%s:%s/aliases",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
	)
	res, err := c.apiCall(method, url, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		responseBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, HTTPError{
			Status:       res.StatusCode,
			ResponseBody: responseBody,
		}
	}
	var body struct {
		Aliases []*Alias `json:"aliases"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	return body.Aliases, nil
}

// DeleteAlias deletes an alias by its name.
func (c *Client) DeleteAlias(aliasName string) (*Alias, error) {
	method := http.MethodDelete
	url := fmt.Sprintf(
		"%s://%s:%s/aliases/%s",
		c.masterNode.Protocol,
		c.masterNode.Host,
		c.masterNode.Port,
		aliasName,
	)
	res, err := c.apiCall(method, url, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		responseBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, HTTPError{
			Status:       res.StatusCode,
			ResponseBody: responseBody,
		}
	}
	var alias Alias
	if err := json.NewDecoder(res.Body).Decode(&alias); err != nil {
		return nil, err
	}
	return &alias, nil
}
