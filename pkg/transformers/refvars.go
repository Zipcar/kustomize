/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package transformers

import (
	"log"
	"sigs.k8s.io/kustomize/v3/pkg/expansion"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/transformers/config"
)

type RefVarTransformer struct {
	varMap            map[string]interface{}
	replacementCounts map[string]int
	fieldSpecs        []config.FieldSpec
	mappingFunc       func(string) interface{}
}

const parentInline = "parent-inline"

// NewRefVarTransformer returns a new RefVarTransformer
// that replaces $(VAR) style variables with values.
// The fieldSpecs are the places to look for occurrences of $(VAR).
func NewRefVarTransformer(
	varMap map[string]interface{}, fs []config.FieldSpec) *RefVarTransformer {
	return &RefVarTransformer{
		varMap:     varMap,
		fieldSpecs: fs,
	}
}

// replacePrimitiveType checks if the field is a string. If not it will return
// the field. In case it is a string, it will expand the field, and perform
// a deepCopy of the expanded value in case it is not of primitiv type.
func (rv *RefVarTransformer) replacePrimitiveType(a interface{}) interface{} {
	s, ok := a.(string)
	if !ok {
		// This field is not of string type.
		// It can not contain a $(VAR)
		return a
	}

	// This field can potientially contain a $(VAR)
	expandedValue := expansion.Expand(s, rv.mappingFunc)

	// Let's perform a deep copy if we didn't inline
	// a primitive type
	return deepCopy(expandedValue)
}

// replaceParentInline allows to inline the complex tree of a variable
// at the same time it allows to replace and patch individual member of
// that inlined tree.
func (rv *RefVarTransformer) replaceParentInline(inMap map[string]interface{}) (interface{}, error) {
	s, _ := inMap[parentInline].(string)

	inlineValue := expansion.Expand(s, rv.mappingFunc)
	newMap, ok := inlineValue.(map[string]interface{})
	if !ok {
		log.Printf("inlining issue with %s", inlineValue)
		return inMap, nil
	}

	newMapCopy := deepCopyMap(newMap)
	mergedMap, err := deepMergeMap(newMapCopy, inMap)
	if err != nil {
		log.Printf("deepMerging issue with %s %v", newMap, err)
		return inMap, nil
	}

	delete(mergedMap, parentInline)
	return mergedMap, nil
}

// replaceVars accepts as 'in' a string, or string array, which can have
// embedded instances of $VAR style variables, e.g. a container command string.
// The function returns the string with the variables expanded to their final
// values.
func (rv *RefVarTransformer) replaceVars(in interface{}) (interface{}, error) {
	switch in.(type) {
	case []interface{}:
		var xs []interface{}
		for _, a := range in.([]interface{}) {
			// Attempt to expand item by item
			xs = append(xs, rv.replacePrimitiveType(a))
		}
		return xs, nil
	case map[string]interface{}:
		inMap := in.(map[string]interface{})

		// Deal with "parent-inline" special expansion
		if _, ok := inMap[parentInline]; ok {
			return rv.replaceParentInline(inMap)
		}

		// Attempt to expand field by field
		xs := make(map[string]interface{}, len(inMap))
		for k, v := range inMap {
			xs[k] = rv.replacePrimitiveType(v)
		}
		return xs, nil
	case string:
		// Attempt to expand this simple field
		return rv.replacePrimitiveType(in), nil
	case nil:
		return nil, nil
	default:
		// This field not contain a $(VAR) since it is not of string type.
		return in, nil
	}
}

// UnusedVars returns slice of Var names that were unused
// after a Transform run.
func (rv *RefVarTransformer) UnusedVars() []string {
	var unused []string
	for k := range rv.varMap {
		_, ok := rv.replacementCounts[k]
		if !ok {
			unused = append(unused, k)
		}
	}
	return unused
}

// Transform replaces $(VAR) style variables with values.
func (rv *RefVarTransformer) Transform(m resmap.ResMap) error {
	rv.replacementCounts = make(map[string]int)
	rv.mappingFunc = expansion.MappingFuncFor(
		rv.replacementCounts, rv.varMap)

	// Then replace the variables. The first pass may inline
	// complex subtree, when the second can replace variables
	// reference inlined during the first pass
	for i := 0; i < 2; i++ {
		for _, res := range m.Resources() {
			for _, fieldSpec := range rv.fieldSpecs {
				if res.OrgId().IsSelected(&fieldSpec.Gvk) {
					if err := MutateField(
						res.Map(), fieldSpec.PathSlice(),
						false, rv.replaceVars); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
