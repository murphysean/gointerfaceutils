package gointerfaceutils

import (
	"errors"
	"regexp"
	"strings"
)

//A Doc path looks like a unix file path. It must begin with a /. Numbers will
//select array elements, and words will select map attributes.
//Examples
//'/name/first' -> ["name","first"]
//'/phonenumbers/0/number -> ["phonenumbers","0","number"]
func parseDocPath(docPath string) ([]string, error) {
	if !strings.HasPrefix(docPath, "/") {
		return nil, errors.New("The Path must begin at the root ('/')")
	}
	if strings.Contains(docPath, "//") {
		return nil, errors.New("The Path can't have null directories('//')")
	}
	if docPath == "/" {
		return []string{}, nil
	}
	return strings.Split(docPath[1:], "/"), nil
}

//A Obj path assumes that the doc is an object, and it uses attribute selectors.
//Array elements are enclosed in brackets [0] and words will select map
//attributes.
//An example
//'user.name.first' -> ["name","first"]
//'user.phonenumbers[0].number' -> ["phonenumbers","0","number"]
func parseObjPath(objPath string) ([]string, error) {
	//First turn all the array accessors into .<num> instead of [num]
	re := regexp.MustCompile("(\\[([0-9]+)\\])")
	objPath = re.ReplaceAllString(objPath, ".${2}")

	if !strings.ContainsAny(objPath, ".") {
		return []string{}, nil
	}

	//Now split it up
	if i := strings.Index(objPath, "."); i != -1 {
		return strings.Split(objPath[i+1:], "."), nil
	}
	return strings.Split(objPath[0:], "."), nil
}
