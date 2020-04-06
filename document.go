package main

import (
	"bytes"
	"encoding/json"
)

// DocumentResponse is the representation of a response
// that contains the data of a document.
type DocumentResponse struct {
	Data  []byte
	Error error
}

// UnmarshalDocument will unmarshal the document data into
// the given interface.
func (ds *DocumentResponse) UnmarshalDocument(document interface{}) error {
	if ds.Error != nil {
		return ds.Error
	}
	err := json.NewDecoder(bytes.NewReader(ds.Data)).Decode(&document)
	return err
}
