package gointerfaceutils

import (
	"errors"
	"time"
)

//OBJECT-----------------------------------------------------------------------------------------
func MustGetObjectAtDocPath(doc interface{}, docPath string) map[string]interface{} {
	o, _ := GetObjectAtDocPath(doc, docPath)
	return o
}

func GetObjectAtDocPath(doc interface{}, docPath string) (map[string]interface{}, error) {
	val, err := GetValueAtDocPath(doc, docPath)
	if err != nil {
		return nil, err
	}

	if o, ok := val.(map[string]interface{}); !ok {
		return nil, errors.New("Value at " + docPath + " path was not an object")
	} else {
		return o, nil
	}
}

func MustGetObjectAtObjPath(doc interface{}, objPath string) map[string]interface{} {
	o, _ := GetObjectAtObjPath(doc, objPath)
	return o
}

func GetObjectAtObjPath(doc interface{}, objPath string) (map[string]interface{}, error) {
	val, err := GetValueAtObjPath(doc, objPath)
	if err != nil {
		return nil, err
	}

	if o, ok := val.(map[string]interface{}); !ok {
		return nil, errors.New("Value at " + objPath + " path was not an object")
	} else {
		return o, nil
	}
}

//ARRAY------------------------------------------------------------------------------------------
func MustGetArrayAtDocPath(doc interface{}, docPath string) []interface{} {
	a, _ := GetArrayAtDocPath(doc, docPath)
	return a
}

func GetArrayAtDocPath(doc interface{}, docPath string) ([]interface{}, error) {
	val, err := GetValueAtDocPath(doc, docPath)
	if err != nil {
		return nil, err
	}

	if a, ok := val.([]interface{}); !ok {
		return nil, errors.New("Value at " + docPath + " path was not an array")
	} else {
		return a, nil
	}
}

func MustGetArrayAtObjPath(doc interface{}, objPath string) []interface{} {
	a, _ := GetArrayAtObjPath(doc, objPath)
	return a
}

func GetArrayAtObjPath(doc interface{}, objPath string) ([]interface{}, error) {
	val, err := GetValueAtObjPath(doc, objPath)
	if err != nil {
		return nil, err
	}

	if o, ok := val.([]interface{}); !ok {
		return nil, errors.New("Value at " + objPath + " path was not an array")
	} else {
		return o, nil
	}
}

//STRING-----------------------------------------------------------------------------------------
func MustGetStringAtDocPath(doc interface{}, docPath string) string {
	s, _ := GetStringAtDocPath(doc, docPath)
	return s
}

func GetStringAtDocPath(doc interface{}, docPath string) (string, error) {
	val, err := GetValueAtDocPath(doc, docPath)
	if err != nil {
		return "", err
	}

	if s, ok := val.(string); !ok {
		return "", errors.New("Value at " + docPath + " path was not a string")
	} else {
		return s, nil
	}
}

func MustGetStringAtObjPath(doc interface{}, objPath string) string {
	s, _ := GetStringAtObjPath(doc, objPath)
	return s
}

func GetStringAtObjPath(doc interface{}, objPath string) (string, error) {
	val, err := GetValueAtObjPath(doc, objPath)
	if err != nil {
		return "", err
	}

	if s, ok := val.(string); !ok {
		return "", errors.New("Value at " + objPath + " path was not a string")
	} else {
		return s, nil
	}
}

//Time-------------------------------------------------------------------------------------------
func MustGetTimeAtDocPath(doc interface{}, docPath string) time.Time {
	t, _ := GetTimeAtDocPath(doc, docPath)
	return t
}

func GetTimeAtDocPath(doc interface{}, docPath string) (time.Time, error) {
	val, err := GetValueAtDocPath(doc, docPath)
	if err != nil {
		return time.Time{}, err
	}

	switch val.(type) {
	case time.Time:
		return val.(time.Time), nil
	case string:
		return time.Parse(time.RFC3339Nano, val.(string))
	case float64:
		return time.Unix(0, int64(val.(float64)*1000000)), nil
	default:
		return time.Time{}, errors.New("Value at " + docPath + " path was not a valid time value")
	}
}

func MustGetTimeAtObjPath(doc interface{}, objPath string) time.Time {
	t, _ := GetTimeAtObjPath(doc, objPath)
	return t
}

func GetTimeAtObjPath(doc interface{}, objPath string) (time.Time, error) {
	val, err := GetValueAtObjPath(doc, objPath)
	if err != nil {
		return time.Time{}, err
	}

	switch val.(type) {
	case time.Time:
		return val.(time.Time), nil
	case string:
		return time.Parse(time.RFC3339Nano, val.(string))
	case float64:
		return time.Unix(0, int64(val.(float64)*1000000)), nil
	default:
		return time.Time{}, errors.New("Value at " + objPath + " path was not a valid time value")
	}
}

//FLOAT------------------------------------------------------------------------------------------
func MustGetFloatAtDocPath(doc interface{}, docPath string) float64 {
	i, _ := GetFloatAtDocPath(doc, docPath)
	return i
}

func GetFloatAtDocPath(doc interface{}, docPath string) (float64, error) {
	val, err := GetValueAtDocPath(doc, docPath)
	if err != nil {
		return 0.0, err
	}

	if s, ok := val.(float64); !ok {
		return 0.0, errors.New("Value at " + docPath + " path was not a float64")
	} else {
		return s, nil
	}
}

func MustGetFloatAtObjPath(doc interface{}, objPath string) float64 {
	i, _ := GetFloatAtObjPath(doc, objPath)
	return i
}

func GetFloatAtObjPath(doc interface{}, objPath string) (float64, error) {
	val, err := GetValueAtObjPath(doc, objPath)
	if err != nil {
		return 0.0, err
	}

	if i, ok := val.(float64); !ok {
		return 0.0, errors.New("Value at " + objPath + " path was not a float64")
	} else {
		return i, nil
	}
}

//BOOLEAN----------------------------------------------------------------------------------------
func MustGetBooleanAtDocPath(doc interface{}, docPath string) bool {
	b, _ := GetBooleanAtDocPath(doc, docPath)
	return b
}

func GetBooleanAtDocPath(doc interface{}, docPath string) (bool, error) {
	val, err := GetValueAtDocPath(doc, docPath)
	if err != nil {
		return false, err
	}

	if b, ok := val.(bool); !ok {
		return false, errors.New("Value at " + docPath + " path was not a boolean")
	} else {
		return b, nil
	}
}

func MustGetBooleanAtObjPath(doc interface{}, objPath string) bool {
	b, _ := GetBooleanAtObjPath(doc, objPath)
	return b
}

func GetBooleanAtObjPath(doc interface{}, objPath string) (bool, error) {
	val, err := GetValueAtObjPath(doc, objPath)
	if err != nil {
		return false, err
	}

	if b, ok := val.(bool); !ok {
		return false, errors.New("Value at " + objPath + " path was not a boolean")
	} else {
		return b, nil
	}
}
