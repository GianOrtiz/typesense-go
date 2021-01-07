package typesense

import (
	"errors"
	"fmt"
)

// ErrConnNotReady is the error that alerts that the connection with the Typesense API
// could not be established, it can be because of a connection  timeout, a unauthorized
// response or a fail.
var ErrConnNotReady = errors.New("typesense connection is not ready")

// ErrCollectionNotFound returned when Typesense can't find the collection.
var ErrCollectionNotFound = errors.New("collection was not found")

// ErrCollectionNameRequired returned when the user tries to create a collection without a name.
var ErrCollectionNameRequired = errors.New("collection name is required")

// ErrCollectionFieldsRequired returned when the user tries to create a collection without
// its fields.
var ErrCollectionFieldsRequired = errors.New("collection fields is required")

// ErrCollectionDuplicate returned when the user tries to create a collection with a name that
// already exists.
var ErrCollectionDuplicate = errors.New("a collection with this name already exists")

// ErrNotFound returned when no resource was found for the request.
var ErrNotFound = errors.New("the resouce you are trying to fetch from Typesense does not exist")

// ErrQueryRequired returned when the user didn't specify a query to search for.
var ErrQueryRequired = errors.New("query field is required")

// ErrQueryByRequired returned when the user didn't specfify fields to query by, the url field
// `query_by` is a required field.
var ErrQueryByRequired = errors.New("query by field is required")

// ErrUnauthorized returned when the API key does not match the Typesense API key.
var ErrUnauthorized = errors.New("the api key does not match the Typesense api key")

// ErrDuplicateID returned when the document the user is trying to index has an id that is already
// in the collection.
var ErrDuplicateID = errors.New("the document you are trying to index has an id that already exists in the collection")

// HTTPError returns an error when an unexpected response status code is
// received.
type HTTPError struct {
	Status       int    `json:"status"`
	ResponseBody []byte `json:"body"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf(
		"unexpected response status code %v. Response body contents %s",
		e.Status,
		string(e.ResponseBody),
	)
}

// APIError is an error returned from the API.
type APIError struct {
	Message string `json:"message"`
}

// Error returns a string representation of the error.
func (e APIError) Error() string {
	return e.Message
}
