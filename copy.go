package gointerfaceutils

import "errors"

func Copy(doc interface{}) (thecopy interface{}, err error) {
	switch doc.(type) {
	case map[string]interface{}:
		ret := make(map[string]interface{})
		for k, _ := range doc.(map[string]interface{}) {
			val, err := Copy(doc.(map[string]interface{})[k])
			if err != nil {
				return doc, err
			}
			ret[k] = val
		}
		return ret, nil
	case []interface{}:
		ret := make([]interface{}, len(doc.([]interface{})))
		for i, _ := range doc.([]interface{}) {
			val, err := Copy(doc.([]interface{})[i])
			if err != nil {
				return doc, err
			}
			ret[i] = val
		}
		return ret, nil
	case string:
		return doc.(string), nil
	case float64:
		return doc.(float64), nil
	case nil:
		return nil, nil
	default:
		return nil, errors.New("Encountered an unknown Type")
	}
}
