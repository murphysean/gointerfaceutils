package gointerfaceutils

func mapEquals(doc1 map[string]interface{}, doc2 interface{}) (ret bool) {
	var doc2map map[string]interface{}
	switch doc2.(type) {
	case map[string]interface{}:
		doc2map = doc2.(map[string]interface{})
	default:
		return false
	}

	for k, _ := range doc1 {
		//Does the key exist in part 2?
		if _, ok := doc2map[k]; !ok {
			return false
		}

		//Do the values equal?
		equal := Equals(doc1[k], doc2map[k])
		if !equal {
			return false
		}
	}

	return true
}

func arrayEquals(doc1 []interface{}, doc2 interface{}) (ret bool) {
	var doc2arr []interface{}
	switch doc2.(type) {
	case []interface{}:
		doc2arr = doc2.([]interface{})
	default:
		return false
	}

	//Do the lengths match
	if len(doc1) != len(doc2arr) {
		return false
	}

	for i, v := range doc1 {
		//Do the values equal?
		equal := Equals(v, doc2arr[i])
		if !equal {
			return false
		}
	}

	return true
}

func Equals(doc1 interface{}, doc2 interface{}) (ret bool) {
	switch doc1.(type) {
	case map[string]interface{}:
		return mapEquals(doc1.(map[string]interface{}), doc2)
	case []interface{}:
		return arrayEquals(doc1.([]interface{}), doc2)
	case string:
		if val, ok := doc2.(string); ok && doc1.(string) == val {
			return true
		}
	case float64:
		if val, ok := doc2.(float64); ok && doc1.(float64) == val {
			return true
		}
	case bool:
		if val, ok := doc2.(bool); ok && doc1.(bool) == val {
			return true
		}
	case nil:
		if doc2 == nil {
			return true
		}
	default:
		return false
	}

	return false
}
