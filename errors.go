package main

import "errors"

// ErrConnNotReady is returned when we can't establish a connection
// with the typesense node.
var ErrConnNotReady = errors.New("typesense connection is not ready")
