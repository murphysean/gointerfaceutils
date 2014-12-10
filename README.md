Go Interface Utils
===

Introduction
---
Go Interface Utils is a set of utilities for interfaces. It was primarily intented to interact with decoded json documents. It has methods to merge documents using both the JSON Merge Patch and JSON Patch methods.

I created it when dealing with unstructured json documents in go. It provides a nice wrapper on top of the already awesome go basics.

It also has methods to deep copy and deep compare documents, match documents based on url query parameters,  and get values within a document based on either path or object based selectors.

This library operates upon the interface{} representations of unmarshaled json documents. The idea is that you have probably already unmarshaled the data to do your own processing before you want to merge, or patch. In this case you should just be able to pass in the raw structures.

In order to get the interface{} pointer from a json string:

	var doc interface{}
	err := json.Unmarshal([]byte(jsonstring), &doc)

### Copy
If you'd like to deep copy an object it's as simple as:

	theCopy, err := gointerfaceutils.Copy(theOriginalDocument)

### Compare
Comparing is also fairly trivial:

	booleanResult := gointerfaceutils.Equals(document1, document2)
	

### Select
Before you start getting and setting values, it's helpful to understand paths. There are two types:
* Document Path: This is much like a path in your filesystem `/` is the root and then you can have folders within the root. Some examples might include `/name/first` and `/phoneNumbers/0/digits`
* Object Path: This is much like using json in javascript. The json object can be selected by using .attribute and [arrayIndex]. An example might look like `user.name.first` or `user.phoneNumbers[0].digits`. In this case the root name of the document is ignored and only things following the first . are considered part of the path.

You can get and set objects at certain paths like so:

	value, err := gointerfaceutils.GetValueAtDocPath(document, "/phoneNumbers/0/digits")
	value, err := gointerfaceutils.GetValueAtObjPath(document, "user.phoneNumbers[0].digits")
	
	doc, err := gointerfaceutils.SetValueAtDocPath(document, "/phoneNumbers/0/digits", value)
	doc, err := gointerfaceutils.SetValueAtObjPath(document, "user.phoneNumbers[0].digits", value)

### JSON Match with Query Parameters
Sometimes you want to filter, or select documents that match some query criteria. Using some of the functionality that the library offers there is another function that will tell you if a document matches query parameters.

You can use it thusly:

	booleanResult := gointerfaceutils.MatchQuery(document, query url.Values)
	
It attempts to be smart about how it does matching, for instance if you query something like `user.identifiers=secretAgent` against this json document `{"identifiers":["secretAgent","instaKill"]}` you would get a match.

### JSON Merge Patch
Visit [RFC 7386](https://tools.ietf.org/html/rfc7386) to learn more about JSON Merge Patch.

To utilize the library, simply provide the documents like so:

	document, err := gointerfaceutils.MergePatch(target, patch)
	
The target is the original, the patch is the proposed changes. If anything goes wrong, the original doc is returned along with the error.

### JSON Patch
Visit [RFC 6902](https://tools.ietf.org/html/rfc6902) to learn more about JavaScript Object Notation (JSON) Patch

To utilize this function, simply provide the documents like so:

	doc, err := gointerfaceutils.Patch(document, patch)
	
If anything goes wrong, the original unchanged document is returned along with the error.