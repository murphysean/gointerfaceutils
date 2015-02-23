package gointerfaceutils

import (
	"encoding/json"
	"testing"
)

var getterTestObj = `{
	"object":{
		"object":{
			"string":"string",
			"bool":true
		},
		"array":["a","b","c","d"],
		"string":"string",
		"integer":45.66,
		"boolean":true
	},
	"array":["a","b","c","d"],
	"string":"string",
	"timestring":"2015-02-22T21:27:55-07:00",
	"timemillis":1424665675000,
	"boolean":true
}`

func TestGetters(t *testing.T) {
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(getterTestObj), &obj)
	if err != nil {
		t.Error(err)
		return
	}

	topObj, err := GetObjectAtObjPath(obj, "")
	if err != nil {
		t.Error(err)
		return
	}

	if !Equals(obj, topObj) {
		t.Error("Not Equal")
		return
	}

	arr, err := GetArrayAtObjPath(obj, "obj.array")
	if err != nil {
		t.Error(err)
		return
	}

	if !Equals(obj["array"], arr) {
		t.Error("Not Equal")
		return
	}

	str, err := GetStringAtObjPath(obj, "obj.object.object.string")
	if err != nil {
		t.Error(err)
		return
	}

	if !Equals("string", str) {
		t.Error("Not Equal")
		return
	}

	i, err := GetFloatAtObjPath(obj, "obj.object.integer")
	if err != nil {
		t.Error(err)
		return
	}

	if !Equals(45.66, i) {
		t.Error("Not Equal")
		return
	}

	if !Equals(true, MustGetBooleanAtObjPath(obj, "obj.boolean")) {
		t.Error("Not Equal")
		return
	}

	if MustGetTimeAtObjPath(obj, "obj.timestring").UnixNano() != 1424665675000000000 {
		t.Error("Not Equal")
		return
	}

	if MustGetTimeAtObjPath(obj, "obj.timemillis").UnixNano() != 1424665675000000000 {
		t.Error("Not Equal")
		return
	}
}
