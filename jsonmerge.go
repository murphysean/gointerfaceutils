package gointerfaceutils

import "errors"

func MergePatch(target interface{}, patch interface{}) (ret interface{}, err error) {
	ret, err = Copy(target)
	if err != nil {
		return target, err
	}

	if _, ok := patch.(map[string]interface{}); ok {
		if _, ok := ret.(map[string]interface{}); !ok {
			ret = make(map[string]interface{})
		}
		for k, v := range patch.(map[string]interface{}) {
			if v == nil {
				delete(ret.(map[string]interface{}), k)
			} else {
				mergedval, err := MergePatch(ret.(map[string]interface{})[k], v)
				if err != nil {
					return target, nil
				}
				ret.(map[string]interface{})[k] = mergedval
			}
		}
		return ret, nil
	} else {
		return patch, nil
	}
	return target, errors.New("Not really sure what happened here")
}
