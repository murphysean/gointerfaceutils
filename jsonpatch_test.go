package gointerfaceutils

import (
	"encoding/json"
	"testing"
)

var positivePatchTests = []string{
	`{ "foo": "bar"}`,
	`[{ "op": "add", "path": "/baz", "value": "qux" }]`,
	`{"baz": "qux","foo": "bar"}`,
	`{ "foo": [ "bar", "baz" ] }`,
	`[{ "op": "add", "path": "/foo/1", "value": "qux" }]`,
	`{ "foo": [ "bar", "qux", "baz" ] }`,
	`{"baz": "qux","foo": "bar"}`,
	`[{ "op": "remove", "path": "/baz" }]`,
	`{ "foo": "bar" }`,
	`{ "foo": [ "bar", "qux", "baz" ] }`,
	`[{ "op": "remove", "path": "/foo/1" }]`,
	`{ "foo": [ "bar", "baz" ] }`,
	`{"baz": "qux","foo": "bar"}`,
	`[{ "op": "replace", "path": "/baz", "value": "boo" }]`,
	`{"baz": "boo","foo": "bar"}`,
	`{"foo": {"bar": "baz","waldo": "fred"},"qux": {"corge": "grault"}}`,
	`[{ "op": "move", "from": "/foo/waldo", "path": "/qux/thud" }]`,
	`{"foo": {"bar": "baz"},"qux": {"corge": "grault","thud": "fred"}}`,
	`{ "foo": [ "all", "grass", "cows", "eat" ] }`,
	`[{ "op": "move", "from": "/foo/1", "path": "/foo/3" }]`,
	`{ "foo": [ "all", "cows", "eat", "grass" ] }`,
	`{"baz": "qux","foo": [ "a", 2, "c" ]}`,
	`[{ "op": "test", "path": "/baz", "value": "qux" },{ "op": "test", "path": "/foo/1", "value": 2 }]`,
	`{"baz": "qux","foo": [ "a", 2, "c" ]}`,
	`{ "foo": "bar" }`,
	`[{ "op": "add", "path": "/child", "value": { "grandchild": { } } }]`,
	`{"foo": "bar","child": {"grandchild": {}}}`,
	`{ "foo": "bar" }`,
	`[{ "op": "add", "path": "/baz", "value": "qux", "xyz": 123 }]`,
	`{"foo": "bar","baz": "qux"}`,
	`{ "foo": ["bar"] }`,
	`[{ "op": "add", "path": "/foo/-", "value": ["abc", "def"] }]`,
	`{ "foo": ["bar", ["abc", "def"]] }`,
}

var negativePatchTests = []string{
	`{ "baz": "qux" }`,
	`[{ "op": "test", "path": "/baz", "value": "bar" }]`,
	`{ "foo": "bar" }`,
	`[{ "op": "add", "path": "/baz/bat", "value": "qux" }]`,
}

func TestJsonPatch(t *testing.T) {
	var err error
	for i := 0; i < len(positivePatchTests); i += 3 {
		var orig interface{}
		err = json.Unmarshal([]byte(positivePatchTests[i]), &orig)
		if err != nil {
			t.Error(err, i, positivePatchTests[i])
			return
		}

		var patch interface{}
		err = json.Unmarshal([]byte(positivePatchTests[i+1]), &patch)
		if err != nil {
			t.Error(err, i, positivePatchTests[i+1])
			return
		}

		var res interface{}
		err = json.Unmarshal([]byte(positivePatchTests[i+2]), &res)
		if err != nil {
			t.Error(err, i, positivePatchTests[i+2])
			return
		}

		var doc interface{}
		doc, err = Patch(orig, patch)
		if err != nil {
			t.Error(err, i)
			return
		}

		equal := Equals(doc, res)
		if !equal {
			t.Errorf("The doc %v patched with %v resulted in %v when it should have resulted in %v", orig, patch, doc, res)
			return
		}
	}

	for i := 0; i < len(negativePatchTests); i += 2 {
		var orig interface{}
		err = json.Unmarshal([]byte(negativePatchTests[i]), &orig)
		if err != nil {
			continue
		}

		var patch interface{}
		err = json.Unmarshal([]byte(negativePatchTests[i+1]), &patch)
		if err != nil {
			continue
		}

		var doc interface{}
		doc, err = Patch(orig, patch)
		if err != nil {
			continue
		}

		t.Errorf("The doc %v patched with %v resulted in %v when it should have resulted in an error", orig, patch, doc)
		return
	}
}
