package gointerfaceutils

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func MatchQuery(doc interface{}, query url.Values) bool {
	for k, v := range query {
		key := k
		if key == "search" {
			if fullTextSearch(doc, v) {
				return true
			}
		}
		condition := "="
		//First determine if the last character in the key is one of (!,>,<)
		if len(k) > 1 {
			lastChar := k[len(k)-1]
			switch lastChar {
			case '!':
				key = k[:len(k)-1]
				condition = "!="
			case '<':
				key = k[:len(k)-1]
				condition = "<="
			case '>':
				key = k[:len(k)-1]
				condition = ">="
			}
		}
		//Parse the key path
		path, err := parseObjPath(key)
		if err != nil {
			return false
		}
		//Get the value at the query key
		obj, err := getValueAtPath(doc, path)
		if err != nil {
			return false
		}
		//Now see if the value matches the obj
		if !match(obj, v, condition) {
			return false
		}
	}

	return true
}

func fullTextSearch(doc interface{}, search []string) bool {
	switch doc.(type) {
	case []interface{}:
		if matchArray(doc.([]interface{}), search, "=") {
			return true
		}
	case map[string]interface{}:
		//Iterate through all the attributes of the object
		for k, v := range doc.(map[string]interface{}) {
			for _, s := range search {
				if k == s {
					return true
				}
			}
			if fullTextSearch(v, search) {
				return true
			}
		}
	case string:
		if matchString(doc.(string), search, "=") {
			return true
		}
	case float64:
		if matchFloat64(doc.(float64), search, "=") {
			return true
		}
	case bool:
		if matchBool(doc.(bool), search, "=") {
			return true
		}
	case nil:
		if matchNull(search, "=") {
			return true
		}
	default:
		return false
	}

	return false
}

func match(obj interface{}, value interface{}, condition string) bool {
	switch obj.(type) {
	case map[string]interface{}:
		return matchMap(obj.(map[string]interface{}), value, condition)
	case []interface{}:
		return matchArray(obj.([]interface{}), value, condition)
	case string:
		return matchString(obj.(string), value, condition)
	case float64:
		return matchFloat64(obj.(float64), value, condition)
	case bool:
		return matchBool(obj.(bool), value, condition)
	case nil:
		return matchNull(value, condition)
	}
	return false
}

func matchMap(obj map[string]interface{}, value interface{}, condition string) bool {
	switch value.(type) {
	case map[string]interface{}:
		switch condition {
		case "=":
			return Equals(obj, value)
		case "!=":
			return !Equals(obj, value)
		}
		return false
	default:
		return false
	}
}

func matchArray(obj []interface{}, value interface{}, condition string) bool {
	switch value.(type) {
	case []interface{}:
		switch condition {
		case "=":
			return Equals(obj, value)
		case "!=":
			return !Equals(obj, value)
		}
	case []string:
		//Do something of a contains
		allMatch := true
		//For each interface in the obj
		for _, s := range value.([]string) {
			oneMatch := false
			//See if it matches one of the string values
		INLOOP:
			for _, v := range obj {
				switch v.(type) {
				case string:
					if matchString(v.(string), s, condition) {
						oneMatch = true
						break INLOOP
					}
				case float64:
					if matchFloat64(v.(float64), s, condition) {
						oneMatch = true
						break INLOOP
					}
				case bool:
					if matchBool(v.(bool), s, condition) {
						oneMatch = true
						break INLOOP
					}
				case nil:
					if matchNull(s, condition) {
						oneMatch = true
						break INLOOP
					}
				default:
					continue INLOOP
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
				if matchString(v.(string), s, condition) {
					return true
				}
			case float64:
				if matchFloat64(v.(float64), s, condition) {
					return true
				}
			case bool:
				if matchBool(v.(bool), s, condition) {
					return true
				}
			case nil:
				if matchNull(s, condition) {
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

func matchString(obj string, value interface{}, condition string) bool {
	switch value.(type) {
	case []interface{}:
		switch condition {
		case "=":
			for _, v := range value.([]interface{}) {
				if matchString(obj, v, "=") {
					return true
				}
			}
			return false
		case "!=":
			for _, v := range value.([]interface{}) {
				if matchString(obj, v, "=") {
					return false
				}
			}
			return true
		}
	case []string:
		switch condition {
		case "!=":
			for _, v := range value.([]string) {
				if matchString(obj, v, "=") {
					return false
				}
			}
			return true
		default:
			for _, v := range value.([]string) {
				if matchString(obj, v, condition) {
					return true
				}
			}
		}
	case string:
		if obj == value.(string) {
			if condition == "=" {
				return true
			} else if condition == "!=" {
				return false
			}
		}
		//The query value could be a Regular Expression
		if sv, ok := value.(string); ok && (condition == "=" || condition == "!=") {
			if re, err := regexp.Compile(sv); err == nil {
				if re.MatchString(obj) {
					if condition == "=" {
						return true
					} else if condition == "!=" {
						return false
					}
				}
				if condition == "!=" {
					return true
				}
			}
		}
		//The object could be a string encoded number
		if num, err := strconv.ParseFloat(obj, 64); err == nil {
			return matchFloat64(num, value.(string), condition)
		}
		//The object could be a string encoded date
		if t1, err := time.Parse(time.RFC3339Nano, obj); err == nil {
			if t2, err := time.Parse(time.RFC3339Nano, value.(string)); err == nil {
				switch condition {
				case "=":
					return t1.Equal(t2)
				case "!=":
					return !t1.Equal(t2)
				case ">=":
					return t1.After(t2)
				case "<=":
					return t1.Before(t2)
				}
			}
		}
	default:
		return false
	}
	return false
}

func matchFloat64(obj float64, value interface{}, condition string) bool {
	switch value.(type) {
	case []interface{}:
		for _, v := range value.([]interface{}) {
			if matchFloat64(obj, v, condition) {
				return true
			}
		}
	case []string:
		for _, v := range value.([]string) {
			if matchFloat64(obj, v, condition) {
				return true
			}
		}
	case string:
		v, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return false
		}
		return matchFloat64(obj, v, condition)
	case float64:
		switch condition {
		case "=":
			return obj == value.(float64)
		case "!=":
			return obj != value.(float64)
		case ">=":
			return obj >= value.(float64)
		case "<=":
			return obj <= value.(float64)
		}
		if obj == value.(float64) {
			return true
		}
	case float32:
		return matchFloat64(obj, float64(value.(float32)), condition)
	case int64:
		return matchFloat64(obj, float64(value.(int64)), condition)
	case int32:
		return matchFloat64(obj, float64(value.(int32)), condition)
	case int16:
		return matchFloat64(obj, float64(value.(int16)), condition)
	case int8:
		return matchFloat64(obj, float64(value.(int8)), condition)
	case int:
		return matchFloat64(obj, float64(value.(int)), condition)
	case byte:
		return matchFloat64(obj, float64(value.(byte)), condition)
	default:
		return false
	}
	return false
}

func matchBool(obj bool, value interface{}, condition string) bool {
	switch value.(type) {
	case []interface{}:
		for _, v := range value.([]interface{}) {
			if matchBool(obj, v, condition) {
				return true
			}
		}
	case []string:
		for _, v := range value.([]string) {
			if matchBool(obj, v, condition) {
				return true
			}
		}
	case string:
		s := value.(string)
		if (obj && strings.ToLower(s) == "true") || (obj && strings.ToLower(s) == "t") {
			if condition == "!=" {
				return false
			}
			return true
		}
		if (!obj && strings.ToLower(s) == "false") || (!obj && strings.ToLower(s) == "f") {
			if condition == "!=" {
				return false
			}
			return true
		}
	case bool:
		if obj == value.(bool) {
			if condition == "!=" {
				return false
			}
			return true
		}
		return false
	default:
		return false
	}
	return false
}

func matchNull(value interface{}, condition string) bool {
	switch value.(type) {
	case []interface{}:
		for _, v := range value.([]interface{}) {
			if matchNull(v, condition) {
				return true
			}
		}
	case []string:
		for _, v := range value.([]string) {
			if matchNull(v, condition) {
				return true
			}
		}
	case string:
		s := value.(string)
		if condition == "!=" && (strings.ToLower(s) == "nil" || strings.ToLower(s) == "null") {
			return false
		}
		if strings.ToLower(s) == "nil" || strings.ToLower(s) == "null" {
			return true
		}
	case nil:
		if condition == "!=" {
			return false
		}
		return true
	default:
		return false
	}
	return false
}
