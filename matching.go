package gointerfaceutils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func MatchQuery(doc interface{}, query url.Values) bool {
	for k, v := range query {
		//Parse the key path
		path, err := parseObjPath(k)
		if err != nil {
			return false
		}
		//Get the value at the query key
		obj, err := getValueAtPath(doc, path)
		if err != nil {
			fmt.Println(err)
			return false
		}
		//Now see if the value matches the obj
		if !match(obj, v) {
			return false
		}
	}

	return true
}

func match(obj interface{}, value interface{}) bool {
	switch obj.(type) {
	case map[string]interface{}:
		return matchMap(obj.(map[string]interface{}), value)
	case []interface{}:
		return matchArray(obj.([]interface{}), value)
	case string:
		return matchString(obj.(string), value)
	case float64:
		return matchFloat64(obj.(float64), value)
	case bool:
		return matchBool(obj.(bool), value)
	case nil:
		return matchNull(value)
	}
	return false
}

func matchMap(obj map[string]interface{}, value interface{}) bool {
	switch value.(type) {
	case map[string]interface{}:
		return Equals(obj, value)
	default:
		return false
	}
}

func matchArray(obj []interface{}, value interface{}) bool {
	switch value.(type) {
	case []interface{}:
		return Equals(obj, value)
	case []string:
		//Do something of a contains
		allMatch := true
		//For each interface in the obj
		for _, s := range value.([]string) {
			oneMatch := false
			//See if it matches one of the string values
			for _, v := range obj {
				switch v.(type) {
				case string:
					if matchString(v.(string), s) {
						oneMatch = true
						break
					}
				case float64:
					if matchFloat64(v.(float64), s) {
						oneMatch = true
						break
					}
				case bool:
					if matchBool(v.(bool), s) {
						oneMatch = true
						break
					}
				case nil:
					if matchNull(s) {
						oneMatch = true
						break
					}
				default:
					continue
				}
			}
			if !oneMatch {
				allMatch = false
				return false
			}
		}
		if allMatch {
			return true
		}
	case string:
		s := value.(string)
		for _, v := range obj {
			switch v.(type) {
			case string:
				if matchString(v.(string), s) {
					return true
				}
			case float64:
				if matchFloat64(v.(float64), s) {
					return true
				}
			case bool:
				if matchBool(v.(bool), s) {
					return true
				}
			case nil:
				if matchNull(s) {
					return true
				}
			default:
				continue
			}
		}
	default:
		return false
	}
	return false
}

func matchString(obj string, value interface{}) bool {
	switch value.(type) {
	case []interface{}:
		for _, v := range value.([]interface{}) {
			if matchString(obj, v) {
				return true
			}
		}
	case []string:
		for _, v := range value.([]string) {
			if matchString(obj, v) {
				return true
			}
		}
	case string:
		if obj == value.(string) {
			return true
		}
	default:
		return false
	}
	return false
}

func matchFloat64(obj float64, value interface{}) bool {
	switch value.(type) {
	case []interface{}:
		for _, v := range value.([]interface{}) {
			if matchFloat64(obj, v) {
				return true
			}
		}
	case []string:
		for _, v := range value.([]string) {
			if matchFloat64(obj, v) {
				return true
			}
		}
	case string:
		v, err := strconv.ParseFloat(value.(string), 64)
		if err == nil && obj == v {
			return true
		}
	case float64:
		if obj == value.(float64) {
			return true
		}
	case float32:
		if obj == float64(value.(float32)) {
			return true
		}
	case int64:
		if obj == float64(value.(int64)) {
			return true
		}
	case int32:
		if obj == float64(value.(int32)) {
			return true
		}
	case int16:
		if obj == float64(value.(int16)) {
			return true
		}
	case int8:
		if obj == float64(value.(int8)) {
			return true
		}
	case int:
		if obj == float64(value.(int)) {
			return true
		}
	case byte:
		if obj == float64(value.(byte)) {
			return true
		}
	default:
		return false
	}
	return false
}

func matchBool(obj bool, value interface{}) bool {
	switch value.(type) {
	case []interface{}:
		for _, v := range value.([]interface{}) {
			if matchBool(obj, v) {
				return true
			}
		}
	case []string:
		for _, v := range value.([]string) {
			if matchBool(obj, v) {
				return true
			}
		}
	case string:
		s := value.(string)
		if (obj && strings.ToLower(s) == "true") || (obj && strings.ToLower(s) == "t") {
			return true
		}
		if (!obj && strings.ToLower(s) == "false") || (!obj && strings.ToLower(s) == "f") {
			return true
		}
	case bool:
		if obj == value.(bool) {
			return true
		}
		return false
	default:
		return false
	}
	return false
}

func matchNull(value interface{}) bool {
	switch value.(type) {
	case []interface{}:
		for _, v := range value.([]interface{}) {
			if matchNull(v) {
				return true
			}
		}
	case []string:
		for _, v := range value.([]string) {
			if matchNull(v) {
				return true
			}
		}
	case string:
		s := value.(string)
		if strings.ToLower(s) == "nil" || strings.ToLower(s) == "null" {
			return true
		}
	case nil:
		return true
	default:
		return false
	}
	return false
}
