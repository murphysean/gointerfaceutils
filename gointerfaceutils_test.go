package gointerfaceutils

import (
	"reflect"
	"testing"
)

func TestGetValueAtPath(t *testing.T) {
	user := make(map[string]interface{})
	username := make(map[string]interface{})
	username["first"] = "Sean"
	username["last"] = "Murphy"
	phonenumbers := make([]interface{}, 2)
	phonenumbers[0] = map[string]interface{}{"number": "123456789"}
	phonenumbers[1] = map[string]interface{}{"number": "987654321"}
	user["phones"] = phonenumbers
	user["name"] = username

	path, err := parseObjPath("user.phones[0].number")
	if err != nil {
		t.Error(err)
		return
	}

	phonenum, err := getValueAtPath(user, path)
	if err != nil {
		t.Error(err)
		return
	}

	if phonenum != "123456789" {
		t.Error("Phones don't match")
	}
}

func TestGetValueAtPathArray(t *testing.T) {
	user := make(map[string]interface{})
	username := make(map[string]interface{})
	username["first"] = "Sean"
	username["last"] = "Murphy"
	phonenumbers := make([]interface{}, 2)
	phonenumbers[0] = map[string]interface{}{"number": "123456789"}
	phonenumbers[1] = map[string]interface{}{"number": "987654321"}
	user["phones"] = phonenumbers
	user["name"] = username

	path, err := parseObjPath("user.phones")
	if err != nil {
		t.Error(err)
		return
	}

	phones, err := getValueAtPath(user, path)
	if err != nil {
		t.Error(err)
		return
	}

	if reflect.TypeOf(phones) != reflect.TypeOf(phonenumbers) {
		t.Error("Didn't get an array")
	}
}

func TestSetValueAtPath(t *testing.T) {
	path, err := parseObjPath("user.name.first")
	if err != nil {
		t.Error(err)
		return
	}
	doc, err := setValueAtPath(nil, path, "Sean")
	if err != nil {
		t.Error(err)
		return
	}

	firstname, err := getValueAtPath(doc, path)
	if err != nil {
		t.Error(err)
		return
	}

	if firstname != "Sean" {
		t.Error("Firstname is not what it was set to be")
	}
}

func TestSetValueAtRootPath(t *testing.T) {
	path, err := parseObjPath("")
	if err != nil {
		t.Error(err)
		return
	}
	doc, err := setValueAtPath(nil, path, "Sean")
	if err != nil {
		t.Error(err)
		return
	}

	if doc != "Sean" {
		t.Error("The document is not what it was set to be")
	}
}
