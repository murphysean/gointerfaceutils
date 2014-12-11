package gointerfaceutils

import "testing"

var positiveDocPathTests = []interface{}{
	`/obj/something/here`,
	[]string{"obj", "something", "here"},
	`/`,
	[]string{},
	`/obj/0/twenty`,
	[]string{"obj", "0", "twenty"},
	`/some/really/long/path/0/10000/to/a/thing`,
	[]string{"some", "really", "long", "path", "0", "10000", "to", "a", "thing"},
}

var negativeDocPathTests = []string{
	`obj/something`,
	`/obj//something`,
}

var positiveObjPathTests = []interface{}{
	`doc.obj.something.here`,
	[]string{"obj", "something", "here"},
	`doc`,
	[]string{"doc"},
	`doc.obj[0].twenty`,
	[]string{"obj", "0", "twenty"},
	`doc.some.really.long.path[0][10000].to.a.thing`,
	[]string{"some", "really", "long", "path", "0", "10000", "to", "a", "thing"},
}

var negativeObjPathTests = []string{}

func checkStringArrayEquals(arr1, arr2 []string) bool {
	for i, v := range arr1 {
		if v != arr2[i] {
			return false
		}
	}
	return true
}

func TestDocPath(t *testing.T) {
	for i := 0; i < len(positiveDocPathTests); i += 2 {
		path, err := parseDocPath(positiveDocPathTests[i].(string))
		if err != nil {
			t.Error(err, positiveDocPathTests[i])
			return
		}

		result := checkStringArrayEquals(path, positiveDocPathTests[i+1].([]string))
		if !result {
			t.Errorf("The path %v should have ended up equal to %v, but instead is %v", positiveDocPathTests[i], positiveDocPathTests[i+1], path)
			return
		}
	}

	for _, v := range negativeDocPathTests {
		_, err := parseDocPath(v)
		if err == nil {
			t.Errorf("The invalid path %v should have resulted in an error", v)
			return
		}
	}
}

func TestObjPath(t *testing.T) {
	for i := 0; i < len(positiveObjPathTests); i += 2 {
		path, err := parseObjPath(positiveObjPathTests[i].(string))
		if err != nil {
			t.Error(err, positiveObjPathTests[i])
			return
		}

		result := checkStringArrayEquals(path, positiveObjPathTests[i+1].([]string))
		if !result {
			t.Errorf("The path %v should have ended up equal to %v, but instead is %v", positiveObjPathTests[i], positiveObjPathTests[i+1], path)
			return
		}
	}

	for _, v := range negativeObjPathTests {
		_, err := parseObjPath(v)
		if err == nil {
			t.Errorf("The invalid path %v should have resulted in an error", v)
			return
		}
	}
}
