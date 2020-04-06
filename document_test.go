package typesense

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestUnmarshalDocument(t *testing.T) {
	type someStruct struct {
		Name string `json:"name"`
	}
	aStruct := someStruct{Name: "name"}
	jsonData, _ := json.Marshal(aStruct)
	resp := DocumentResponse{
		Data:  jsonData,
		Error: nil,
	}
	var s someStruct
	if err := resp.UnmarshalDocument(&s); err != nil {
		t.Errorf("Expected to receive no errors, received %v", err)
	}
	if !reflect.DeepEqual(s, aStruct) {
		t.Errorf("Expected to receive %v, received %v", aStruct, s)
	}
}

func TestUnmarshalDocument_withError(t *testing.T) {
	type someStruct struct {
		Name string `json:"name"`
	}
	aStruct := someStruct{Name: "name"}
	jsonData, _ := json.Marshal(aStruct)
	errDocumentNotFound := errors.New("document not found")
	resp := DocumentResponse{
		Data:  jsonData,
		Error: errDocumentNotFound,
	}
	var s someStruct
	if err := resp.UnmarshalDocument(&s); err != errDocumentNotFound {
		t.Errorf("Expected to receive error %v, received %v", errDocumentNotFound, err)
	}
}
