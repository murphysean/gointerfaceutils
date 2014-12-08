package gointerfaceutils

import (
	"encoding/json"
	"testing"
)

var positiveMergePatchTests = []string{
	`{"a":"b","c":{"d":"e","f":"g"}}`,
	`{"a":"z","c":{"f":null}}`,
	`{"a":"z","c":{"d":"e"}}`,
	`{"title": "Goodbye!","author" : {"givenName" : "John","familyName" : "Doe"},"tags":[ "example", "sample" ],"content": "This will be unchanged"}`,
	`{"title": "Hello!","phoneNumber": "+01-123-456-7890","author": {"familyName": null},"tags": [ "example" ]}`,
	`{"title": "Hello!","author" : {"givenName" : "John"},"tags": [ "example" ],"content": "This will be unchanged","phoneNumber": "+01-123-456-7890"}`,
	`{"a":"b"}`,
	`{"a":"c"}`,
	`{"a":"c"}`,
	`{"a":"b"}`,
	`{"b":"c"}`,
	`{"a":"b","b":"c"}`,
	`{"a":"b"}`,
	`{"a":null}`,
	`{}`,
	`{"a":"b","b":"c"}`,
	`{"a":null}`,
	`{"b":"c"}`,
	`{"a":["b"]}`,
	`{"a":"c"}`,
	`{"a":"c"}`,
	`{"a":"c"}`,
	`{"a":["b"]}`,
	`{"a":["b"]}`,
	`{"a": {"b": "c"}}`,
	`{"a": { "b": "d", "c": null}}`,
	`{"a": {"b": "d"}}`,
	`{"a": [{"b":"c"}]}`,
	`{"a": [1]}`,
	`{"a": [1]}`,
	`["a","b"]`,
	`["c","d"]`,
	`["c","d"]`,
	`{"a":"b"}`,
	`["c"]`,
	`["c"]`,
	`{"a":"foo"}`,
	`null`,
	`null`,
	`{"a":"foo"}`,
	`"bar"`,
	`"bar"`,
	`{"e":null}`,
	`{"a":1}`,
	`{"e":null,"a":1}`,
	`[1,2]`,
	`{"a":"b","c":null}`,
	`{"a":"b"}`,
	`{}`,
	`{"a":{"bb":{"ccc":null}}}`,
	`{"a":{"bb":{}}}`,
}

func TestMergePatch(t *testing.T) {
	var err error
	for i := 0; i < len(positiveMergePatchTests); i += 3 {
		var orig interface{}
		err = json.Unmarshal([]byte(positiveMergePatchTests[i]), &orig)
		if err != nil {
			t.Error(err, i, positiveMergePatchTests[i])
			return
		}

		var patch interface{}
		err = json.Unmarshal([]byte(positiveMergePatchTests[i+1]), &patch)
		if err != nil {
			t.Error(err, i, positiveMergePatchTests[i+1])
			return
		}

		var res interface{}
		err = json.Unmarshal([]byte(positiveMergePatchTests[i+2]), &res)
		if err != nil {
			t.Error(err, i, positiveMergePatchTests[i+2])
			return
		}

		var doc interface{}
		doc, err = MergePatch(orig, patch)

		equal := Equals(doc, res)
		if !equal {
			t.Errorf("The doc %v merged with %v resulted in %v when it should have resulted in %v", orig, patch, doc, res)
			return
		}
	}
}
