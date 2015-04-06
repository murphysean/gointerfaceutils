package gointerfaceutils

import (
	"encoding/json"
	"testing"
)

var copyTests = []string{
	`64`,
	`64`,
	`"string"`,
	`"string"`,
	`["string1","string2","string3"]`,
	`["string1","string2","string3"]`,
	`null`,
	`null`,
	`{"key":"value"}`,
	`{"key":"value"}`,
	`{"string":"string","number":64,"null":null,"array":["string1","string2","string3"],"mapofmaps":{"submap":{"k":"v"}}}`,
	`{"string":"string","number":64,"null":null,"array":["string1","string2","string3"],"mapofmaps":{"submap":{"k":"v"}}}`,
}

func TestCopy(t *testing.T) {
	var doc1, doc2 interface{}
	var err error
	var result bool
	//Do the positive tests
	for i := 0; i < len(copyTests); i += 2 {
		err = json.Unmarshal([]byte(positiveEqualTests[i]), &doc1)
		if err != nil {
			t.Error(err)
			return
		}

		doc1copy, err := Copy(doc1)
		if err != nil {
			t.Error(err, doc1)
			return
		}

		err = json.Unmarshal([]byte(copyTests[i+1]), &doc2)
		if err != nil {
			t.Error(err, doc1)
			return
		}

		doc1copybytes, err := json.Marshal(&doc1copy)
		if err != nil {
			t.Error(err, doc1)
			return
		}

		result = Equals(doc1copy, doc2)
		if !result {
			t.Error("These two documents should be equal after copy:")
			t.Error("\t", string(doc1copybytes))
			t.Error("\t", copyTests[i+1])
			return
		}
	}

	//Do the negative tests
	for i := 0; i < len(copyTests); i += 2 {
		err = json.Unmarshal([]byte(copyTests[i]), &doc1)
		if err != nil {
			t.Error(err, doc1)
			return
		}

		doc1copy, err := Copy(doc1)
		if err != nil {
			t.Error(err, doc1)
			return
		}

		//Now modify the crap out of doc1, doc1copy should remain the same
		switch doc1copy.(type) {
		case map[string]interface{}:
			doc1.(map[string]interface{})["the world"] = "mine"
		case []interface{}:
			doc1.([]interface{})[0] = "ELEMENT 0 IS NOW THIS"
			doc1 = append(doc1.([]interface{}), float64(885))
		case []string:
			doc1.([]string)[0] = "ELEMENT 0 IS NOW THIS"
			doc1 = append(doc1.([]string), "885")
		case string:
			doc1 = "A Whole New String"
		case float64:
			doc1 = float64(885)
		case nil:
			doc1 = float64(885)
		}

		err = json.Unmarshal([]byte(copyTests[i+1]), &doc2)
		if err != nil {
			t.Error(err, doc1)
			return
		}

		doc1copybytes, err := json.Marshal(&doc1copy)
		if err != nil {
			t.Error(err, doc1)
			return
		}

		result = Equals(doc1copy, doc2)
		if !result {
			t.Error("These two documents should still be equal after copy and original doc modification:")
			t.Error("\t", string(doc1copybytes))
			t.Error("\t", copyTests[i+1])
			return
		}
	}
}