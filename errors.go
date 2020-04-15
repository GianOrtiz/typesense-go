package typesense

import "errors"

// ErrConnNotReady means the connection with the Typesense API
// could not be established, it can be because of a connection
// timeout, a unauthorized response or a fail.
var ErrConnNotReady = errors.New("typesense connection is not ready")

// ErrCollectionNotFound means Typesense does not have this collection
// registered in it.
var ErrCollectionNotFound = errors.New("collection was not found")

// ErrQueryRequired means the user didn't specify a query to search for.
var ErrQueryRequired = errors.New("query field is required")

// ErrQueryByRequired means the user didn't specfify fields to query by,
// the url field `query_by` is a required field in typesense.
var ErrQueryByRequired = errors.New("query by field is required")
