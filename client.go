package typesense

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultHeaderKey = "X-TYPESENSE-API-KEY"

type httpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

// Client is the client to communicate with the Typesense API.
type Client struct {
	httpClient       httpClient
	masterNode       *Node
	readReplicaNodes []*Node
}

// Node is a Typesense node, either the master or a read replica.
type Node struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	APIKey   string `json:"apiKey"`
}

// APIResponse is the default API message response.
type APIResponse struct {
	Message string `json:"message"`
}

// NewClient configures a client using the master node and timeout
// seconds.
func NewClient(masterNode *Node, timeoutSeconds int, replicaNodes ...*Node) (*Client, error) {
	var readReplicas []*Node
	for _, replica := range replicaNodes {
		readReplicas = append(readReplicas, replica)
	}

	client := Client{
		httpClient: &http.Client{
			Timeout: time.Duration(time.Second * time.Duration(timeoutSeconds)),
		},
		masterNode:       masterNode,
		readReplicaNodes: readReplicas,
	}

	if err := client.Ping(); err != nil {
		return nil, err
	}

	return &client, nil
}

// Ping checks if the client has a connection with the Typesense API.
func (c *Client) Ping() error {
	if ok := c.Health(); !ok {
		return ErrConnNotReady
	}
	return nil
}

// DebugInfo retrieves the debug information from the Typesense API.
func (c *Client) DebugInfo() (string, error) {
	method := http.MethodGet
	url := fmt.Sprintf("%s://%s:%s/debug", c.masterNode.Protocol, c.masterNode.Host, c.masterNode.Port)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add(defaultHeaderKey, c.masterNode.APIKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	type debugResponse struct {
		Version string `json:"version"`
	}
	var debug debugResponse
	if err := json.NewDecoder(resp.Body).Decode(&debug); err != nil {
		return "", err
	}
	return debug.Version, nil
}

// Health checks the health information from the Typesense API.
func (c *Client) Health() bool {
	method := http.MethodGet
	url := fmt.Sprintf("%s://%s:%s/health", c.masterNode.Protocol, c.masterNode.Host, c.masterNode.Port)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add(defaultHeaderKey, c.masterNode.APIKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	type healthResponse struct {
		OK bool `json:"ok"`
	}
	var health healthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return false
	}
	return health.OK
}
