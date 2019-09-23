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

package config

import (
	"reflect"
	"sort"
	"testing"

	"sigs.k8s.io/kustomize/v3/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/v3/pkg/gvk"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/resmaptest"
	"sigs.k8s.io/kustomize/v3/pkg/resource"
	"sigs.k8s.io/kustomize/v3/pkg/types"
)

var cmap = gvk.Gvk{Kind: "ConfigMap"}
var catalog = types.Target{
	Gvk: gvk.Gvk{
		Group:   "my.org",
		Version: "v1",
		Kind:    "SomeCatalog",
	},
	APIVersion: "my.org/v1",
	Name:       "catalog-name",
}

// newVarRefSlice sorts the fsSlice according to the path
// instead of the GVK (see fielspec.go)
func newVarRefSlice(fs []FieldSpec) fsSlice {
	va := make([]FieldSpec, len(fs))
	copy(va, fs)
	sort.Slice(va, func(i, j int) bool {
		return va[i].Path < va[j].Path
	})
	return va
}

// TestResMapScanner tests the detection of VarRef in ResMap
func TestResMapScanner(t *testing.T) {
	rf := resource.NewFactory(
		kunstruct.NewKunstructuredFactoryImpl())
	type given struct {
		toScan resmap.ResMap
	}
	type expected struct {
		vars         []types.Var
		varReference fsSlice
	}
	testCases := []struct {
		description string
		given       given
		expected    expected
	}{
		{
			description: "auto-detect-var",
			given: given{
				toScan: resmaptest_test.NewRmBuilder(t, rf).
					Add(map[string]interface{}{
						"apiVersion": catalog.APIVersion,
						"kind":       catalog.Gvk.Kind,
						"metadata": map[string]interface{}{
							"name": "catalog-name",
							"annotations": map[string]interface{}{
								"my.org/after": "simple",
							},
						},
						"spec": map[string]interface{}{
							"key_name": "key_value",
						}}).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
							"annotations": map[string]interface{}{
								"my.org": `$(SomeCatalog.catalog-name.metadata.annotations["my.org/after"])`,
								"my/org": `$(SomeCatalog.catalog-name.metadata.annotations["my.org/after"])`,
							},
						},
						"data": map[string]interface{}{
							"item1": "$(SomeCatalog.catalog-name.spec.key_name)",
							"item2": "bla",
						}}).ResMap(),
			},
			expected: expected{
				vars: []types.Var{
					{Name: `SomeCatalog.catalog-name.metadata.annotations["my.org/after"]`,
						ObjRef:   catalog,
						FieldRef: types.FieldSelector{FieldPath: `metadata.annotations["my.org/after"]`}},
					{Name: "SomeCatalog.catalog-name.spec.key_name",
						ObjRef:   catalog,
						FieldRef: types.FieldSelector{FieldPath: "spec.key_name"}},
				},
				varReference: fsSlice{
					{Gvk: cmap, Path: `metadata/annotations/my.org`},
					{Gvk: cmap, Path: `metadata/annotations/my\/org`},
					{Gvk: cmap, Path: `data/item1`},
				},
			},
		},
		{
			description: "auto-detect-parent-inline",
			given: given{
				toScan: resmaptest_test.NewRmBuilder(t, rf).
					Add(map[string]interface{}{
						"apiVersion": catalog.APIVersion,
						"kind":       catalog.Gvk.Kind,
						"metadata": map[string]interface{}{
							"name": "catalog-name",
							"annotations": map[string]interface{}{
								"my.org/after": "simple",
							},
						},
						"spec": map[string]interface{}{
							"key_name": "key_value",
						}}).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							ParentInline: "$(SomeCatalog.catalog-name.spec.key_name)",
							"foofield3":  "bla",
						}}).ResMap(),
			},
			expected: expected{
				vars: []types.Var{
					{Name: "SomeCatalog.catalog-name.spec.key_name",
						ObjRef:   catalog,
						FieldRef: types.FieldSelector{FieldPath: "spec.key_name"}},
				},
				varReference: fsSlice{
					{Gvk: cmap, Path: "data"},
				},
			},
		},
		{
			description: "auto-detect-arrays",
			given: given{
				toScan: resmaptest_test.NewRmBuilder(t, rf).
					Add(map[string]interface{}{
						"apiVersion": catalog.APIVersion,
						"kind":       catalog.Gvk.Kind,
						"metadata": map[string]interface{}{
							"name": "catalog-name",
							"annotations": map[string]interface{}{
								"my.org/after": "simple",
							},
						},
						"spec": map[string]interface{}{
							"key_name": "key_value",
						}}).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"someslice": []interface{}{
								map[string]interface{}{
									"item1": "foo",
									"item2": "bla",
								},
								map[string]interface{}{
									"item1": "foo",
									"item2": "$(SomeCatalog.catalog-name.spec.key_name)",
								},
							},
						}}).ResMap(),
			},
			expected: expected{
				vars: []types.Var{
					{Name: "SomeCatalog.catalog-name.spec.key_name",
						ObjRef:   catalog,
						FieldRef: types.FieldSelector{FieldPath: "spec.key_name"}},
				},
				varReference: fsSlice{
					{Gvk: cmap, Path: "data/someslice/item2"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// arrange
			// TODO(jeb): Create a second serie of set where
			// the user manually entered data is provided.
			tr := NewResMapScanner(types.VarSet{}, fsSlice{})

			// act
			err := tr.BuildAutoConfig(tc.given.toScan)

			// assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			va, ve := newVarRefSlice(tr.DiscoveredConfig().VarReference), newVarRefSlice(tc.expected.varReference)
			if !reflect.DeepEqual(va, ve) {
				t.Fatalf("VarReference actual doesn't match expected: \nACTUAL:\n%v\nEXPECTED:\n%v", va, ve)
			}

			varaset := tr.DiscoveredVars()
			varsa, varse := varaset.AsSlice(), tc.expected.vars
			if len(varsa) != len(varse) {
				t.Fatalf("unexpected size %d", len(varsa))
			}
			for i := range varsa {
				if !varsa[i].DeepEqual(varse[i]) {
					t.Fatalf("unexpected varsa[%d]:\n  %v\n  %v", i, varsa[i], varse[i])
				}
			}
		})
	}
}
