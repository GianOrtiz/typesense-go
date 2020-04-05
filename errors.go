package main

import "errors"

// ErrConnNotReady is returned when we can't establish a connection
// with the typesense node.
var ErrConnNotReady = errors.New("typesense connection is not ready")

// ErrCollectionNotFound is used when the typesense can't find the collection.
var ErrCollectionNotFound = errors.New("collection was not found")
