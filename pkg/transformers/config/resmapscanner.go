package config

import (
	"fmt"
	"regexp"
	"strings"

	"sigs.k8s.io/kustomize/v3/pkg/expansion"
	"sigs.k8s.io/kustomize/v3/pkg/gvk"
	"sigs.k8s.io/kustomize/v3/pkg/resid"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/types"
)

type ResMapScanner struct {
	manualVars types.VarSet
	manualRefs fsSlice
	autoVars   types.VarSet
	autoRefs   fsSlice
}

type detectedRef struct {
	pathSlice     []string
	detectedNames []string
}

// KindRegistry contains a cache of the GVK used in the resources.
var KindRegistry = map[string]gvk.Gvk{}
var IndexRegex *regexp.Regexp = regexp.MustCompile(`<[0-9]+>`)

const (
	ParentInline string = "parent-inline"
	Dot          string = "."
	Slash        string = "/"
	Hash         string = "#"
)

// NewResMapScanner returns a new ResMapScanner
// that detects $(Kind.name.path) style variables with values.
func NewResMapScanner(userVars types.VarSet, userRefs fsSlice) *ResMapScanner {
	manualVars := userVars.Copy()
	manualRefs := fsSlice{}
	autoVars := types.NewVarSet()
	autoRefs := fsSlice{}

	return &ResMapScanner{
		manualVars: manualVars,
		manualRefs: manualRefs,
		autoVars:   autoVars,
		autoRefs:   autoRefs,
	}
}

// DiscoveredVars returns the list of Var to add to the
// consolidated var section of the kustomization.yaml(s)
// This allows the user to not have to do that manually.
func (rv *ResMapScanner) DiscoveredVars() types.VarSet {
	return rv.autoVars
}

// DiscoveredConfig returns a TransformerConfig containing
// a consolidated VarReference sections.
// This allows the user to not have to do that manually.
func (rv *ResMapScanner) DiscoveredConfig() *TransformerConfig {
	return &TransformerConfig{
		VarReference: rv.autoRefs,
	}
}

// Walk the path (curPath) of a resource (in) and collect
// the detected Var in the refMap object
func (rv *ResMapScanner) collectReferences(refMap map[string]detectedRef, curPath []string, id resid.ResId, in interface{}) {

	switch typedIn := in.(type) {
	case []interface{}:
		// Check each member of the slice for variable reference.
		// Note that create fake members _x_ in the path that will need to be removed
		// when creating the varReference.
		for idx, v := range typedIn {
			rv.collectReferences(refMap, append(curPath, fmt.Sprintf("<%v>", idx)), id, v)
		}
	case map[string]interface{}:
		// Check each member of the map for variable reference.
		for key, v := range typedIn {
			rv.collectReferences(refMap, append(curPath, key), id, v)
		}
	case string:
		// Look for potential variable references within the string
		detectedNames := expansion.Detect(typedIn)
		if len(detectedNames) == 0 {
			return
		}

		if len(detectedNames) == 1 && curPath[len(curPath)-1] == ParentInline {
			normalizedPath := make([]string, len(curPath)-1)
			copy(normalizedPath, curPath[:len(curPath)-1])
			refMap[strings.Join(normalizedPath, Hash)] = detectedRef{
				pathSlice:     normalizedPath,
				detectedNames: detectedNames,
			}
			return
		}

		pathSlice := make([]string, len(curPath))
		copy(pathSlice, curPath)
		refMap[strings.Join(pathSlice, Hash)] = detectedRef{
			pathSlice:     pathSlice,
			detectedNames: detectedNames,
		}
	}
}

// filterVar checks the syntax of the detected var and validate
// it matches the $(Kind.name.path) patterns.
func (rv *ResMapScanner) buildVar(detectedName string) (*types.Var, error) {
	// Note:
	// $(Ingress.name.metadata.annotations['ingress.auth-secretkubernetes.io/auth-secret'])
	// will be splitted in s in a strange way, but the fieldpath is rebuild
	// properly a few line further.
	s := strings.Split(detectedName, Dot)

	if len(s) < 3 {
		return nil, fmt.Errorf("var %s does not match expected "+
			"pattern $(Kind.name.path)", detectedName)
	}

	kind := s[0]
	name := s[1]
	fieldPath := strings.Join(s[2:], Dot)

	if _, ok := KindRegistry[kind]; !ok {
		// We don't have a entry for that kind.
		return nil, fmt.Errorf("var $(%s) referencing an unknown  "+
			"or conflicting Kind %s", detectedName, kind)
	}

	group := KindRegistry[kind].Group
	version := KindRegistry[kind].Version
	objref := types.Target{
		Gvk: gvk.Gvk{
			Group:   group,
			Version: version,
			Kind:    kind,
		},
		APIVersion: group + Slash + version,
		Name:       name,
	}
	fieldref := types.FieldSelector{
		FieldPath: fieldPath,
	}
	tVar := &types.Var{
		Name:     detectedName,
		ObjRef:   objref,
		FieldRef: fieldref,
	}

	return tVar, nil
}

// normalizedPathSlice
func (rv *ResMapScanner) normalizePathSlice(pathSlice []string) ([]string, error) {
	normalizedPath := []string{}
	for idx, elt := range pathSlice {
		if strings.Contains(elt, Hash) {
			// Oops. The separator character we used to build our key
			// is also part of path itself.
			return nil, fmt.Errorf("Potential issue triggered by # in yaml field %s", pathSlice[idx])
		}

		// According to fieldspec.PathSlice, need to espace the slash
		normalizedElt := strings.Replace(elt, Slash, escapedForwardSlash, -1)

		// we also need to remove the <x> that we added when collecting
		// the references. according to MutateField getFirstPathElement we should
		// replace it by "[]".
		normalizedElt = IndexRegex.ReplaceAllString(normalizedElt, "")
		if normalizedElt != "" {
			normalizedPath = append(normalizedPath, normalizedElt)
		}
	}
	return normalizedPath, nil
}

// buildFieldSpec builds FieldSpec to add to the VarReference.
func (rv *ResMapScanner) buildFieldSpec(id resid.ResId, pathSlice []string) (*FieldSpec, error) {

	normalizedPath, err := rv.normalizePathSlice(pathSlice)
	if err != nil {
		return nil, err
	}

	// varReference path are using / as separator.
	path := strings.Join(normalizedPath, Slash)

	return &FieldSpec{
		Gvk:  gvk.FromKind(id.Kind),
		Path: path,
	}, nil
}

// BuildAutoConfig scans the ResMap and detects the $(Kind.name.path) pattern
func (rv *ResMapScanner) BuildAutoConfig(m resmap.ResMap) error {
	// Cache the GVK used by the project
	for _, res := range m.Resources() {
		// TODO(jeb): We need to check for conflicting information
		// in t.manualVars and in the resources.
		// TODO(jeb): We also need to check that for a specific
		// kind, the apiversion is used everywhere.
		groupversionkind := res.GetGvk()
		KindRegistry[groupversionkind.Kind] = groupversionkind
	}

	for _, res := range m.Resources() {
		id := res.OrgId()
		referenceMap := map[string]detectedRef{}
		startPath := []string{}

		// Let's walk the resource and collect the variable
		// and the location where the pattern has been detected
		rv.collectReferences(referenceMap, startPath, id, res.Map())

		for _, detectedRef := range referenceMap {
			varReference, err := rv.buildFieldSpec(id, detectedRef.pathSlice)
			if err != nil {
				// Algorithm can't deal with that referencePath,
				// probably because it is not actually used for a variable.
				continue
			}

			for _, detectedName := range detectedRef.detectedNames {
				tVar, err := rv.buildVar(detectedName)
				if err != nil {
					// Algorithm can't deal with the detected variable names,
					// probably because it is not actually a variable.
					continue
				}

				// Check that this could be a variable by looking for the corresponding resource.
				// First try in the same namespace as the current resource.
				targetId := resid.NewResIdWithNamespace(tVar.ObjRef.GVK(), tVar.ObjRef.Name, res.CurId().Namespace)
				idMatcher := targetId.Equals
				matched := m.GetMatchingResourcesByCurrentId(idMatcher)

				if len(matched) == 0 {
					// Look using original namespace and name
					targetId = resid.NewResIdWithNamespace(tVar.ObjRef.GVK(), tVar.ObjRef.Name, res.OrgId().Namespace)
					idMatcher = targetId.Equals
					matched = m.GetMatchingResourcesByOriginalId(idMatcher)
				}

				if len(matched) == 0 {
					// Look using current name but no namespace
					targetId = resid.NewResId(tVar.ObjRef.GVK(), tVar.ObjRef.Name)
					idMatcher = targetId.GvknEquals
					matched = m.GetMatchingResourcesByCurrentId(idMatcher)
				}

				if len(matched) == 0 {
					// Look using original name but no namespace
					targetId = resid.NewResId(tVar.ObjRef.GVK(), tVar.ObjRef.Name)
					idMatcher = targetId.GvknEquals
					matched = m.GetMatchingResourcesByOriginalId(idMatcher)
				}

				// If not, this is probably not a variable.
				if len(matched) == 1 {
					_, err := matched[0].GetFieldValue(tVar.FieldRef.FieldPath)
					if err != nil {
						// We detected $(validkind.validname.invalidfieldspec)
						// This is probably not a variable.
						continue
					}
					tVar.ObjRef.Name = matched[0].OrgId().Name
					tVar.ObjRef.Namespace = matched[0].OrgId().Namespace
					err = rv.manualVars.Absorb(*tVar)
					if err != nil {
						// Detected a GVK for that var which conflicts with the manual entered one.
						// Let's trust the user and ignore that potential variable.
						// TODO(jeb): Probably won't detect potential duplicate
						// between if fieldSpec have differnet format: spec.foo[bar] and spec.foo.bar
						continue
					}

					rv.autoVars.Merge(*tVar)
					matched[0].AbsorbRefVarName(*tVar)
					rv.autoRefs, _ = rv.autoRefs.mergeOne(*varReference)
				}
			}
		}
	}
	return nil
}
