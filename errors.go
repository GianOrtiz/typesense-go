package typesense

import "errors"

// What I need to do to show better errors is get the code from
// https://github.com/typesense/typesense/blob/master/src/core_api.cpp
// them check every return of errors, for example, at line 59 we
// have the code:
// res.set_400("Parameter `name` is required")
// and add a check to see if name is being send, creating a new err
// for this.

// ErrConnNotReady means the connection with the Typesense API
// could not be established, it can be because of a connection
// timeout, a unauthorized response or a fail.
var ErrConnNotReady = errors.New("typesense connection is not ready")

// ErrCollectionNotFound means Typesense does not have this collection
// registered in it.
var ErrCollectionNotFound = errors.New("collection was not found")

// ErrCollectionNameRequired returned when the user tries to create a
// new collection without a name.
var ErrCollectionNameRequired = errors.New("collection name is required")

// ErrCollectionFieldsRequired returned when the user tries to create a
// new collection without fields.
var ErrCollectionFieldsRequired = errors.New("collection fields is required")

// ErrCollectionDuplicate returned when the user tries to create a new
// collection with a name that already exists.
var ErrCollectionDuplicate = errors.New("a collection with this name already exists")

// ErrNotFound returned when no resource was found for the request.
var ErrNotFound = errors.New("not found")

// ErrQueryRequired means the user didn't specify a query to search for.
var ErrQueryRequired = errors.New("query field is required")

// ErrQueryByRequired means the user didn't specfify fields to query by,
// the url field `query_by` is a required field in typesense.
var ErrQueryByRequired = errors.New("query by field is required")
