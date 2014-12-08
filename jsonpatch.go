package gointerfaceutils

import (
	"errors"
	"fmt"
	"strconv"
)

func Patch(doc interface{}, patch interface{}) (ret interface{}, err error) {
	//First make a deep copy of the document, so that if any error occurs we can retain the original document
	ret, err = Copy(doc)
	if err != nil {
		return doc, err
	}

	//Ensure that the patch document is an array
	if patchDoc, ok := patch.([]interface{}); ok {
		for i, v := range patchDoc {
			if patchOp, ok := v.(map[string]interface{}); ok {
				pathInterface, pexists := patchOp["path"]
				if !pexists {
					return doc, errors.New(fmt.Sprintf("The Patch Operation @ index %v did not have a valid path attribute", i))
				}
				pathString, pis := pathInterface.(string)
				if !pis {
					return doc, errors.New(fmt.Sprintf("The Patch Operation @ index %v did not have a valid string path attribute", i))
				}
				path, err := parseDocPath(pathString)
				if err != nil {
					return doc, errors.New(fmt.Sprintf("The Patch Operation @ index %v did not have a valid path attribute. Resulted in error: ", i, err.Error()))
				}

				switch patchOp["op"] {
				case "add":
					value, vexists := patchOp["value"]
					if !vexists {
						return doc, errors.New(fmt.Sprintf("The Add Patch Operation @ index %v did not have a valid value attribute", i))
					}
					ret, err = patchAdd(ret, path, value)
					if err != nil {
						return doc, errors.New(fmt.Sprintf("The Add Patch Operation @ index %v failed with error: %v", i, err.Error()))
					}
				case "remove":
					ret, err = patchRemove(ret, path)
					if err != nil {
						return doc, errors.New(fmt.Sprintf("The Remove Patch Operation @ index %v failed with error: %v", i, err.Error()))
					}
				case "replace":
					value, vexists := patchOp["value"]
					if !vexists {
						return doc, errors.New(fmt.Sprintf("The Replace Patch Operation @ index %v did not have a valid value attribute", i))
					}
					ret, err = patchReplace(ret, path, value)
					if err != nil {
						return doc, errors.New(fmt.Sprintf("The Replace Patch Operation @ index %v failed with error: %v", i, err.Error()))
					}
				case "move":
					fromInterface, fexists := patchOp["from"]
					if !fexists {
						return doc, errors.New(fmt.Sprintf("The Move Patch Operation @ index %v did not have a valid from attribute", i))
					}
					fromString, fis := fromInterface.(string)
					if !fis {
						return doc, errors.New(fmt.Sprintf("The Move Patch Operation @ index %v did not have a valid string from attribute", i))
					}
					from, err := parseDocPath(fromString)
					if err != nil {
						return doc, errors.New(fmt.Sprintf("The Move Patch Operation @ index %v did not have a valid from attribute. Resulted in error: ", i, err.Error()))
					}
					ret, err = patchMove(ret, path, from)
					if err != nil {
						return doc, errors.New(fmt.Sprintf("The Move Patch Operation @ index %v failed with error: %v", i, err.Error()))
					}
				case "copy":
					fromInterface, fexists := patchOp["from"]
					if !fexists {
						return doc, errors.New(fmt.Sprintf("The Move Patch Operation @ index %v did not have a valid from attribute", i))
					}
					fromString, fis := fromInterface.(string)
					if !fis {
						return doc, errors.New(fmt.Sprintf("The Move Patch Operation @ index %v did not have a valid string from attribute", i))
					}
					from, err := parseDocPath(fromString)
					if err != nil {
						return doc, errors.New(fmt.Sprintf("The Move Patch Operation @ index %v did not have a valid from attribute. Resulted in error: ", i, err.Error()))
					}
					ret, err = patchCopy(ret, path, from)
					if err != nil {
						return doc, errors.New(fmt.Sprintf("The Copy Patch Operation @ index %v failed with error: %v", i, err.Error()))
					}
				case "test":
					value, vexists := patchOp["value"]
					if !vexists {
						return doc, errors.New(fmt.Sprintf("The Test Patch Operation @ index %v did not have a valid value attribute", i))
					}
					_, err = patchTest(ret, path, value)
					if err != nil {
						return doc, errors.New(fmt.Sprintf("The Test Patch Operation @ index %v failed with error: %v", i, err.Error()))
					}
				default:
					return doc, errors.New(fmt.Sprintf("The Patch Operation @ index %v did not have a valid op attribute", i))
				}
			} else {
				return doc, errors.New(fmt.Sprintf("The Patch Operation @ index %v was not a object", i))
			}
		}
	} else {
		return doc, errors.New("The Patch Document was not an array")
	}
	return ret, nil
}

/*
The "add" operation performs one of the following functions,
depending upon what the target location references:

* If the target location specifies an array index, a new value is
inserted into the array at the specified index.
*If the target location specifies an object member that does not
already exist, a new member is added to the object.
*If the target location specifies an object member that does exist,
that member's value is replaced.

The operation object MUST contain a "value" member whose content
specifies the value to be added.

For example:
{ "op": "add", "path": "/a/b/c", "value": [ "foo", "bar" ] }

When the operation is applied, the target location MUST reference one of:
*The root of the target document - whereupon the specified value
becomes the entire content of the target document.
*A member to add to an existing object - whereupon the supplied
value is added to that object at the indicated location.  If the
member already exists, it is replaced by the specified value.
*An element to add to an existing array - whereupon the supplied
value is added to the array at the indicated location.  Any
elements at or above the specified index are shifted one position
to the right.  The specified index MUST NOT be greater than the
number of elements in the array.  If the "-" character is used to
index the end of the array (see [RFC6901]), this has the effect of
appending the value to the array.

Because this operation is designed to add to existing objects and
arrays, its target location will often not exist.  Although the
pointer's error handling algorithm will thus be invoked, this
specification defines the error handling behavior for "add" pointers
to ignore that error and add the value as specified.

However, the object itself or an array containing it does need to
exist, and it remains an error for that not to be the case.  For
example, an "add" with a target location of "/a/b" starting with this
document:

{ "a": { "foo": 1 } }

is not an error, because "a" exists, and "b" will be added to its
value.  It is an error in this document:

{ "q": { "bar": 2 } }

because "a" does not exist.
*/
func patchAdd(doc interface{}, path []string, value interface{}) (ret interface{}, err error) {
	if len(path) == 0 {
		return value, nil
	}
	selector := path[len(path)-1]
	parentValue, err := getValueAtPath(doc, path[:len(path)-1])
	if err != nil {
		return doc, err
	}

	switch parentValue.(type) {
	case map[string]interface{}:
		parentValue.(map[string]interface{})[selector] = value
	case []interface{}:
		if selector == "-" {
			selector = strconv.Itoa(len(parentValue.([]interface{})))
		}
		i, err := strconv.Atoi(selector)
		if err != nil {
			return doc, errors.New("The selector " + selector + " wasn't a valid index into the array")
		}
		if 0 > i || i > len(parentValue.([]interface{})) {
			return doc, errors.New("The selector " + selector + " was an out of bound index into the array")
		}

		newArr := append(parentValue.([]interface{}), nil)
		copy(newArr[i+1:], newArr[i:])
		newArr[i] = value

		doc, err = setValueAtPath(doc, path[:len(path)-1], newArr)
	}

	return doc, nil
}

/*
The "remove" operation removes the value at the target location.

The target location MUST exist for the operation to be successful.

For example:

{ "op": "remove", "path": "/a/b/c" }

If removing an element from an array, any elements above the
specified index are shifted one position to the left.
*/
func patchRemove(doc interface{}, path []string) (ret interface{}, err error) {
	if len(path) == 0 {
		return nil, nil
	}

	selector := path[len(path)-1]

	parentValue, err := getValueAtPath(doc, path[:len(path)-1])
	if err != nil {
		return doc, err
	}

	switch parentValue.(type) {
	case map[string]interface{}:
		if _, exists := parentValue.(map[string]interface{})[selector]; !exists {
			return doc, errors.New("The selector " + selector + " didn't exist")
		}
		delete(parentValue.(map[string]interface{}), selector)
	case []interface{}:
		parentArr := parentValue.([]interface{})
		if selector == "-" {
			selector = strconv.Itoa(len(parentArr))
		}
		i, err := strconv.Atoi(selector)
		if err != nil {
			return doc, errors.New("The selector " + selector + " wasn't a valid index into the array")
		}
		if 0 > i || i >= len(parentArr) {
			return doc, errors.New("The selector " + selector + " was an out of bound index into the array")
		}

		copy(parentArr[i:], parentArr[i+1:])
		parentArr[len(parentArr)-1] = nil // or the zero value of T
		parentArr = parentArr[:len(parentArr)-1]

		doc, err = setValueAtPath(doc, path[:len(path)-1], parentArr)
	}

	return doc, nil
}

/*
The "replace" operation replaces the value at the target location
with a new value.  The operation object MUST contain a "value" member
whose content specifies the replacement value.

The target location MUST exist for the operation to be successful.

For example:

{ "op": "replace", "path": "/a/b/c", "value": 42 }

This operation is functionally identical to a "remove" operation for
a value, followed immediately by an "add" operation at the same
location with the replacement value.
*/
func patchReplace(doc interface{}, path []string, value interface{}) (ret interface{}, err error) {
	if len(path) == 0 {
		return value, nil
	}

	selector := path[len(path)-1]

	parentValue, err := getValueAtPath(doc, path[:len(path)-1])
	if err != nil {
		return doc, err
	}

	switch parentValue.(type) {
	case map[string]interface{}:
		if _, exists := parentValue.(map[string]interface{})[selector]; !exists {
			return doc, errors.New("The selector " + selector + " didn't exist")
		}
		parentValue.(map[string]interface{})[selector] = value
	case []interface{}:
		parentArr := parentValue.([]interface{})
		if selector == "-" {
			selector = strconv.Itoa(len(parentArr))
		}
		i, err := strconv.Atoi(selector)
		if err != nil {
			return doc, errors.New("The selector " + selector + " wasn't a valid index into the array")
		}
		if 0 > i || i >= len(parentArr) {
			return doc, errors.New("The selector " + selector + " was an out of bound index into the array")
		}

		parentArr[i] = value
	}

	return doc, nil
}

/*
The "move" operation removes the value at a specified location and
adds it to the target location.

The operation object MUST contain a "from" member, which is a string
containing a JSON Pointer value that references the location in the
target document to move the value from.

The "from" location MUST exist for the operation to be successful.

For example:

{ "op": "move", "from": "/a/b/c", "path": "/a/b/d" }

This operation is functionally identical to a "remove" operation on
the "from" location, followed immediately by an "add" operation at
the target location with the value that was just removed.

The "from" location MUST NOT be a proper prefix of the "path"
location; i.e., a location cannot be moved into one of its children.
*/
func patchMove(doc interface{}, path []string, from []string) (ret interface{}, err error) {
	value, err := getValueAtPath(doc, from)
	if err != nil {
		return doc, err
	}

	doc, err = patchRemove(doc, from)
	if err != nil {
		return doc, err
	}
	doc, err = patchAdd(doc, path, value)
	if err != nil {
		return doc, err
	}

	return doc, nil
}

/*
The "copy" operation copies the value at a specified location to the
target location.

The operation object MUST contain a "from" member, which is a string
containing a JSON Pointer value that references the location in the
target document to copy the value from.

The "from" location MUST exist for the operation to be successful.

For example:

{ "op": "copy", "from": "/a/b/c", "path": "/a/b/e" }

This operation is functionally identical to an "add" operation at the
target location using the value specified in the "from" member.
*/
func patchCopy(doc interface{}, path []string, from []string) (ret interface{}, err error) {
	value, err := getValueAtPath(doc, from)
	if err != nil {
		return doc, err
	}

	doc, err = patchAdd(doc, path, value)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

/*
The "test" operation tests that a value at the target location is
equal to a specified value.

The operation object MUST contain a "value" member that conveys the
value to be compared to the target location's value.

The target location MUST be equal to the "value" value for the
operation to be considered successful.

Here, "equal" means that the value at the target location and the
value conveyed by "value" are of the same JSON type, and that they
are considered equal by the following rules for that type:

strings: are considered equal if they contain the same number of
Unicode characters and their code points are byte-by-byte equal.

numbers: are considered equal if their values are numerically
equal.

arrays: are considered equal if they contain the same number of
values, and if each value can be considered equal to the value at
the corresponding position in the other array, using this list of
type-specific rules.

objects: are considered equal if they contain the same number of
members, and if each member can be considered equal to a member in
the other object, by comparing their keys (as strings) and their
values (using this list of type-specific rules).

literals (false, true, and null): are considered equal if they are
the same.

Note that the comparison that is done is a logical comparison; e.g.,
whitespace between the member values of an array is not significant.

Also, note that ordering of the serialization of object members is
not significant.

For example:

{ "op": "test", "path": "/a/b/c", "value": "foo" }
*/
func patchTest(doc interface{}, path []string, value interface{}) (ret interface{}, err error) {
	val, err := getValueAtPath(doc, path)
	if err != nil {
		return doc, err
	}

	eq := Equals(val, value)
	if !eq {
		return doc, errors.New("The Patch Test Failed")
	}

	return doc, nil
}
