package gointerfaceutils

import (
	"encoding/json"
	"testing"
)

var positiveEqualTests = []string{
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

var negativeEqualTests = []string{
	`64`,
	`65`,
	`"string"`,
	`"anotherstring"`,
	`["string5","string6","string9"]`,
	`["string5","string6","string7"]`,
	`null`,
	`1234`,
	`{"key":"value"}`,
	`{"key":"value2"}`,
	`{"string":"string","number":64,"null":null,"array":["string1","string2","string3"],"mapofmaps":{"submap":{"k":"v"}}}`,
	`{"string":"string","number":64,"null":null,"array":["string1","string2","string3"],"mapofmaps":{"submap":{"k":"v2"}}}`,
}

func TestEquals(t *testing.T) {
	var doc1, doc2 interface{}
	var err error
	var result bool
	//Do the positive tests
	for i := 0; i < len(positiveEqualTests); i += 2 {
		err = json.Unmarshal([]byte(positiveEqualTests[i]), &doc1)
		if err != nil {
			t.Error(err)
			return
		}

		err = json.Unmarshal([]byte(positiveEqualTests[i+1]), &doc2)
		if err != nil {
			t.Error(err)
			return
		}

		result = Equals(doc1, doc2)
		if !result {
			t.Error("These two documents should be equal:")
			t.Error("\t", positiveEqualTests[i])
			t.Error("\t", positiveEqualTests[i+1])
			return
		}
	}

	//Do the negative tests
	for i := 0; i < len(negativeEqualTests); i += 2 {
		err = json.Unmarshal([]byte(negativeEqualTests[i]), &doc1)
		if err != nil {
			t.Error(err)
			return
		}

		err = json.Unmarshal([]byte(negativeEqualTests[i+1]), &doc2)
		if err != nil {
			t.Error(err)
			return
		}

		result = Equals(doc1, doc2)
		if result {
			t.Error("These two documents should not be equal:")
			t.Error("\t", negativeEqualTests[i])
			t.Error("\t", negativeEqualTests[i+1])
			return
		}
	}
}
