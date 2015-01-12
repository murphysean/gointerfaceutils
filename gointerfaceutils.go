package gointerfaceutils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
)

func GetMD5HashForString(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetMD5HashForJSONDoc(doc interface{}) string {
	bytes, err := json.Marshal(&doc)
	if err != nil {
		return ""
	}
	text := string(bytes)

	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func SetValueAtDocPath(doc interface{}, docPath string, value interface{}) (interface{}, error) {
	path, err := parseDocPath(docPath)
	if err != nil {
		return nil, err
	}

	return setValueAtPath(doc, path, value)
}

func SetValueAtObjPath(doc interface{}, objPath string, value interface{}) (interface{}, error) {
	path, err := parseObjPath(objPath)
	if err != nil {
		return nil, err
	}

	return setValueAtPath(doc, path, value)
}

func GetValueAtDocPath(doc interface{}, docPath string) (interface{}, error) {
	path, err := parseDocPath(docPath)
	if err != nil {
		return nil, err
	}

	return getValueAtPath(doc, path)
}

func GetValueAtObjPath(doc interface{}, objPath string) (interface{}, error) {
	path, err := parseObjPath(objPath)
	if err != nil {
		return nil, err
	}

	return getValueAtPath(doc, path)
}

func getValueAtSelector(doc interface{}, selector string) (ret interface{}, err error) {
	switch doc.(type) {
	case map[string]interface{}:
		val, ok := doc.(map[string]interface{})[selector]
		if !ok {
			return doc, errors.New("The map didn't have any value at key " + selector)
		}
		return val, nil
	case []interface{}:
		i, err := strconv.Atoi(selector)
		if err != nil {
			return doc, errors.New("The selector " + selector + " wasn't a valid index into the array")
		}
		if 0 > i || i >= len(doc.([]interface{})) {
			return doc, errors.New("The selector " + selector + " was an out of bound index into the array")
		}
		return doc.([]interface{})[i], nil
	default:
		return doc, errors.New("The selector " + selector + " wasn't valid for the document")
	}
}

func getValueAtPath(doc interface{}, path []string) (ret interface{}, err error) {
	ret = doc
	for _, selector := range path {
		ret, err = getValueAtSelector(ret, selector)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func setValueAtPath(doc interface{}, path []string, value interface{}) (ret interface{}, err error) {
	switch doc.(type) {
	case nil:
		if len(path) == 0 {
			return value, nil
		}
		doc = make(map[string]interface{})
	case map[string]interface{}:
	case []interface{}:
	default:
		if len(path) == 0 {
			return value, nil
		}
		return doc, errors.New("Couldn't non destructivly create given path in doc")
	}

	ret = doc
	for i, selector := range path {
		//If this is the last node in the path, do the set
		if i == len(path)-1 {
			switch ret.(type) {
			case map[string]interface{}:
				ret.(map[string]interface{})[selector] = value
				break
			case []interface{}:

			default:
				return doc, errors.New("Couldn't non destructivly create given path in doc")
			}
			_, ok := ret.(map[string]interface{})
			if ok {

			}

		}
		//Get the value at selector
		ret, err = getValueAtSelector(ret, selector)
		if err != nil {
			switch ret.(type) {
			case map[string]interface{}:
				//TODO If the next value is 0, perhaps I should set this to an array instead?
				ret.(map[string]interface{})[selector] = make(map[string]interface{})
				ret = ret.(map[string]interface{})[selector]
			case []interface{}:
				i, err := strconv.Atoi(selector)
				if err != nil {
					return doc, errors.New("The selector " + selector + " wasn't a valid array index")
				}
				if 0 > i || i > len(ret.([]interface{})) {
					return doc, errors.New("The selector " + selector + " was an out of bound index into the array")
				}
				if i == len(ret.([]interface{})) {
					//TODO If the next value is 0, perhaps I should set this to an array instead?
					ret = append(ret.([]interface{}), make(map[string]interface{}))
					ret = ret.([]interface{})[len(ret.([]interface{}))-1]
				}
				return doc.([]interface{})[i], nil
			default:
				return doc, errors.New("Couldn't non destructivly create given path in doc")
			}
		}
	}
	return doc, nil
}
