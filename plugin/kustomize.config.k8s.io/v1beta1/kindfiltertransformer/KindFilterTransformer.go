// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

//go:generate go run sigs.k8s.io/kustomize/v3/cmd/pluginator
package main

import (
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resid"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	// Filter contains the list of resource names to filter out
	Filter []string `json:"filter,omitempty" yaml:"filter,omitempty"`
}

//noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin

func (p *plugin) Config(
	_ ifc.Loader, _ *resmap.Factory, c []byte) (err error) {
	p.Filter = []string{}
	return yaml.Unmarshal(c, p)
}

func (p *plugin) Transform(m resmap.ResMap) (err error) {
	ids := m.AllIds()
	filter := newFilterSet(p.Filter)
	for _, id := range ids {
		if filter.In(id) {
			m.Remove(id)
		}
	}
	return nil
}

type FilterSet struct {
	resids map[string]struct{}
}

func newFilterSet(m []string) *FilterSet {
	fs := &FilterSet{
		resids: map[string]struct{}{},
	}
	for _, res := range m {
		fs.resids[res] = struct{}{}
	}
	return fs
}

func (fs *FilterSet) In(key resid.ResId) bool {
	_, ok := fs.resids[key.Kind]
	return ok
}
