// +build integration

package integration

import (
	"log"
	"os"
	"testing"

	"github.com/GianOrtiz/typesense-go"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
)

const (
	testTypesenseHost     = "localhost"
	testTypesensePort     = "8108"
	testTypesenseProtocol = "http"
	testTypesenseAPIKey   = "api-key"

	strFieldName = "strField"
	strTestValue = "value"

	floatFieldName = "intField"
	floatTestValue = 3.5
)

var (
	testClient *typesense.Client

	testCollection = typesense.CollectionSchema{
		Name: "test",
		Fields: []typesense.CollectionField{
			{
				Name: strFieldName,
				Type: "string",
			},
			{
				Name: floatFieldName,
				Type: "float",
			},
		},
		DefaultSortingField: floatFieldName,
	}

	testDocument = map[string]interface{}{
		"id":           "0",
		strFieldName:   strTestValue,
		floatFieldName: floatTestValue,
	}

	testAlias = typesense.Alias{
		Name:           "alias",
		CollectionName: testCollection.Name,
	}
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker daemon: %v", err)
	}

	optsTypesense := dockertest.RunOptions{
		Repository: "typesense/typesense",
		Tag:        "0.14.0",
		Env: []string{
			"TYPESENSE_API_KEY=" + testTypesenseAPIKey,
			"TYPESENSE_DATA_DIR=/data",
		},
		Mounts: []string{
			"tmp:/data",
		},
		ExposedPorts: []string{"8108"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8108": {
				{
					HostIP:   "0.0.0.0",
					HostPort: "8108",
				},
			},
		},
	}
	typesenseRes, err := pool.RunWithOptions(&optsTypesense)
	if err != nil {
		log.Fatalf("Could not start MySQL: %v", err)
	}
	typesenseRes.Expire(100)
	defer func() {
		if err := pool.Purge(typesenseRes); err != nil {
			log.Fatalf("Could not purge Typesense: %v", err)
		}
	}()

	masterNode := &typesense.Node{
		Host:     testTypesenseHost,
		Port:     testTypesensePort,
		Protocol: testTypesenseProtocol,
		APIKey:   testTypesenseAPIKey,
	}
	err = pool.Retry(func() error {
		testClient = typesense.NewClient(masterNode, 40)
		return testClient.Ping()
	})
	if err != nil {
		log.Fatalf("Could not connect to the Typesense test instance: %v", err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

// The following tests must succeed sequentially in order to validate
// the client.

func TestCreateCollection(t *testing.T) {
	_, err := testClient.CreateCollection(testCollection)
	assert.Equal(t, nil, err)
	if err != nil {
		log.Fatal(err)
	}
}

func TestRetrieveCollections(t *testing.T) {
	cs, err := testClient.RetrieveCollections()
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(cs))
	if err == nil && len(cs) == 1 {
		assert.Equal(t, testCollection.Name, cs[0].Name)
	}
}

func TestRetrieveCollection(t *testing.T) {
	c, err := testClient.RetrieveCollection(testCollection.Name)
	assert.Equal(t, nil, err)
	if err != nil {
		assert.Equal(t, testCollection.Name, c.Name)
	}
}

func TestIndexDocument(t *testing.T) {
	docRes := testClient.IndexDocument(testCollection.Name, testDocument)
	assert.Equal(t, nil, docRes.Error)
	if docRes.Error != nil {
		log.Fatal(docRes.Error)
	}
	var returnDocument map[string]interface{}
	err := docRes.UnmarshalDocument(&returnDocument)
	assert.Equal(t, nil, err)
	assert.Equal(t, testDocument, returnDocument)
}

func TestRetrieveDocument(t *testing.T) {
	docID, _ := testDocument["id"].(string)
	docRes := testClient.RetrieveDocument(testCollection.Name, docID)
	assert.Equal(t, nil, docRes.Error)
	var returnDocument map[string]interface{}
	err := docRes.UnmarshalDocument(&returnDocument)
	assert.Equal(t, nil, err)
	assert.Equal(t, testDocument, returnDocument)
}

func TestSearchDocument(t *testing.T) {
	searchRes, err := testClient.Search(testCollection.Name, strTestValue, []string{strFieldName}, nil)
	assert.Equal(t, nil, err)
	if err == nil {
		assert.Equal(t, 1, searchRes.Found)
		if searchRes.Found > 1 {
			assert.Equal(t, testDocument, searchRes.Hits[0].Document)
		}
	}
}

func TestDeleteDocument(t *testing.T) {
	docID, _ := testDocument["id"].(string)
	docRes := testClient.DeleteDocument(testCollection.Name, docID)
	assert.Equal(t, nil, docRes.Error)
	var returnDocument map[string]interface{}
	err := docRes.UnmarshalDocument(&returnDocument)
	assert.Equal(t, nil, err)
	assert.Equal(t, testDocument, returnDocument)
}

func TestAPIKey(t *testing.T) {
	createKey := typesense.APIKey{
		Actions:     []string{typesense.ActionAll},
		Collections: []string{testCollection.Name},
		Description: "A new API key",
	}
	apiKey, err := testClient.CreateAPIKey(createKey)
	assert.Equal(t, nil, err)
	if err != nil {
		t.Fatal(err)
	}
	key, err := testClient.GetAPIKey(apiKey.ID)
	assert.Equal(t, nil, err)
	assert.Equal(t, apiKey.ID, key.ID)
	keys, err := testClient.GetAPIKeys()
	assert.Equal(t, nil, err)
	assert.Greater(t, len(keys), 0)
	err = testClient.DeleteAPIKey(apiKey.ID)
	assert.Equal(t, nil, err)
}

func TestCreateAlias(t *testing.T) {
	_, err := testClient.CreateAlias(
		testAlias.Name,
		&testAlias,
	)
	assert.Equal(t, nil, err)
}

func TestRetrieveAlias(t *testing.T) {
	a, err := testClient.RetrieveAlias(testAlias.Name)
	assert.Equal(t, nil, err)
	if a != nil {
		assert.Equal(t, *a, testAlias)
	}
}

func TestRetrieveAliases(t *testing.T) {
	as, err := testClient.RetrieveAliases()
	assert.Equal(t, nil, err)
	assert.GreaterOrEqual(t, len(as), 0)
}

func TestDeleteAlias(t *testing.T) {
	_, err := testClient.DeleteAlias(testAlias.Name)
	assert.Equal(t, nil, err)
}

func TestDeleteCollection(t *testing.T) {
	c, err := testClient.DeleteCollection(testCollection.Name)
	assert.Equal(t, nil, err)
	if err != nil {
		assert.Equal(t, testCollection.Name, c.Name)
	}
}
