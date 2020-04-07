package typesense

import (
	"bytes"
	"encoding/json"
)

// DocumentResponse is the response returned with a
// document. Because the document can't be retrieved
// immediatyl we wrap the data into this struct so it
// can be unmarshaled after.
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
