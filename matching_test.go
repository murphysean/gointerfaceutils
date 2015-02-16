package gointerfaceutils

import (
	"encoding/json"
	"net/url"
	"testing"
)

var positiveMatchingTests = []string{
	`firstName=Sean`,
	`{"firstName":"Sean","lastName":"Bob"}`,
	`obj.name.first=Sean`,
	`{"name":{"first":"Sean"}}`,
	`obj.phones[2].digits=987654321`,
	`{"name":{"first":"Sean"},"phones":[{"digits":12345},{"digits":56789},{"digits":987654321}]}`,
	`firstName!=Sean`,
	`{"firstName":"Shiz","lastName":"Bob"}`,
	`rating%3E=3`,
	`{"rating":5}`,
	`rating%3C=7`,
	`{"rating":5}`,
	`lastName=%5E%5Ba-z%5D%2B%5C%5B%5B0-9%5D%2B%5C%5D%24`,
	`{"firstName":"Sean","lastName":"adam[23]"}`,
	`rating%3E=3`,
	`{"rating":"5"}`,
	`createDate%3E=2006-01-02T15%3A04%3A05.999999999-07%3A00`,
	`{"createDate":"2014-01-02T15:04:05-07:00"}`,
}

func TestMatchingPath(t *testing.T) {
	for i := 0; i < len(positiveMatchingTests); i += 2 {
		q, _ := url.ParseQuery(positiveMatchingTests[i])
		var d interface{}
		json.Unmarshal([]byte(positiveMatchingTests[i+1]), &d)
		if !MatchQuery(d, q) {
			t.Log(d, q)
			t.Error("Obj Should Match Query")
			return
		}
	}
}

func TestMatchMap(t *testing.T) {
	map1 := map[string]interface{}{"key": "value"}
	map2 := map[string]interface{}{"key": "value"}
	map3 := map[string]interface{}{"k": "v"}

	if !matchMap(map1, map2, "=") {
		t.Error("Map Should Match Map")
		return
	}

	if matchMap(map2, map3, "=") {
		t.Error("Map Shouldn't Match Map")
		return
	}

	if !matchMap(map2, map3, "!=") {
		t.Error("Map Shouldn't Match Map")
		return
	}
}

func TestMatchArray(t *testing.T) {
	arr := []interface{}{"cool", "beans", "bro", "waz", "up", "doc", nil, "la", "dee", "da"}
	idarr := []interface{}{"cool", "beans", "bro", "waz", "up", "doc", nil, "la", "dee", "da"}
	idarr2 := []interface{}{"cool", "beans", "bro", "waz", "up", "doc", nil, "la", "dee", "da", 99, map[string]interface{}{"key": "value"}}
	if !matchArray(arr, []string{"cool", "beans"}, "=") {
		t.Error("Strings Should Match Array")
		return
	}

	if !matchArray(arr, "cool", "=") {
		t.Error("String Should Match Array")
		return
	}

	if !matchArray(arr, idarr, "=") {
		t.Error("Array Should Match Array")
		return
	}

	if matchArray(arr, []string{"cool", "beaner"}, "=") {
		t.Error("Strings Shouldn't Match Array")
		return
	}

	if matchArray(arr, "coolness", "=") {
		t.Error("String Shouldn't Match Array")
		return
	}

	if matchArray(arr, idarr2, "=") {
		t.Error("Array Shouldn't Match Array")
		return
	}
}

func TestMatchString(t *testing.T) {
	if !matchString("hello", "hello", "=") {
		t.Error("Strings Should Match")
		return
	}

	if !matchString("hello", "world", "!=") {
		t.Error("Strings Shouldn't Match")
		return
	}

	if !matchString("hello", []string{"hello"}, "=") {
		t.Error("Strings was contained in array")
		return
	}

	if matchString("hello", "world", "=") {
		t.Error("Strings Shouldn't Match")
		return
	}

	if matchString("hello", []string{"world"}, "=") {
		t.Error("String was not contained in array")
		return
	}

	if matchString("hello", []string{"hello", "world"}, "!=") {
		t.Error("String was contained in array")
		return
	}
}

func TestMatchFloat64(t *testing.T) {
	if !matchFloat64(float64(100.123), float64(100.123), "=") {
		t.Error("Numbers Should Match")
		return
	}

	if !matchFloat64(float64(100.123), 100.123, "=") {
		t.Error("Numbers Should Match")
		return
	}

	if !matchFloat64(float64(100), 100, "=") {
		t.Error("Numbers Should Match")
		return
	}

	if !matchFloat64(float64(100.123), "100.123", "=") {
		t.Error("Numbers and String Should Match")
		return
	}

	if matchFloat64(float64(100.123), float64(10000000.12345), "=") {
		t.Error("Numbers Shouldn't Match")
		return
	}

	if matchFloat64(float64(6.5), float64(2.3), "<=") {
		t.Error("6.5 >= 2.3")
		return
	}

	if !matchFloat64(float64(6.5), float64(2.3), ">=") {
		t.Error("6.5 >= 2.3")
		return
	}
}

func TestMatchBool(t *testing.T) {
	if !matchBool(true, "true", "=") || !matchBool(true, "TRUE", "=") || !matchBool(true, "t", "=") || !matchBool(true, "T", "=") {
		t.Error("The string true should match true")
		return
	}

	if !matchBool(false, "false", "=") || !matchBool(false, "FALSE", "=") || !matchBool(false, "f", "=") || !matchBool(false, "F", "=") {
		t.Error("The string false should match false")
		return
	}

	if !matchBool(true, true, "=") {
		t.Error("True should match true")
		return
	}

	if !matchBool(false, false, "=") {
		t.Error("False should match false")
		return
	}

	if matchBool(true, 100, "=") {
		t.Error("True should match number")
		return
	}

	if matchBool(true, false, "=") {
		t.Error("True should match False")
		return
	}

	if matchBool(true, true, "!=") {
		t.Error("True should not match true with a '!=' condition")
		return
	}

	if !matchBool(true, true, ">=") {
		t.Error("True should match true with a '>=' condition")
		return
	}

	if !matchBool(true, true, "<=") {
		t.Error("True should match true with a '<=' condition")
		return
	}
}

func TestMatchNull(t *testing.T) {
	if !matchNull("nil", "=") || !matchNull("NIL", "=") {
		t.Error("The string nil should match null")
		return
	}

	if !matchNull("null", "=") || !matchNull("NULL", "=") {
		t.Error("The string null should match null")
		return
	}

	if !matchNull(nil, "=") {
		t.Error("The actual nil value should match null")
		return
	}

	if matchNull("helloworld", "=") {
		t.Error("A string shouldn't match null")
		return
	}

	if matchNull(1234, "=") {
		t.Error("A number shouldn't match null")
		return
	}

	if !matchNull("nil", ">=") || !matchNull("NIL", "<=") {
		t.Error("Greater and Less conditions should match")
		return
	}

	if matchNull("nil", "!=") {
		t.Error("Not Equal should not match nil")
		return
	}
}
