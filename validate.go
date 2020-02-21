package go3mf

import (
	"errors"
	"sort"
	"strings"

	specerr "github.com/qmuntal/go3mf/errors"
)

type validatorResource struct {
	path string
	id   uint32
}

type validator struct {
	m        *Model
	strict   bool
	warnings []error
	ids      map[validatorResource]struct{}
}

func Validate(model *Model, strict bool) []error {
	v := validator{m: model, strict: strict}
	v.Validate()
	return v.warnings
}

func (v *validator) AddWarning(err error) {
	v.warnings = append(v.warnings, err)
}

// Validate checks that the model is conformant with the 3MF spec.
func (v *validator) Validate() {
	v.ids = make(map[validatorResource]struct{})
	v.validateRelationship(v.m.RootRelationships, "")
	rootPath := v.m.Path
	if rootPath == "" {
		rootPath = DefaultPartModelName
	}
	names := map[string]struct{}{
		rootPath: struct{}{},
	}
	for path, c := range v.m.Childs {
		if _, ok := names[path]; ok {
			v.AddWarning(specerr.ErrOPCDuplicatedModelName)
		}
		names[path] = struct{}{}
		v.validateRelationship(c.Relationships, path)
	}
	v.validateRelationship(v.m.Relationships, rootPath)

	v.validateNamespaces()

	for path, c := range v.m.Childs {
		v.validateResources(&c.Resources, path)
	}
	v.validateResources(&v.m.Resources, rootPath)
}

var allowedMetadataNames = [...]string{ // sorted
	"application", "copyright", "creationdate", "description", "designer",
	"licenseterms", "modificationdate", "rating", "title",
}

func (v *validator) checkMetdadata(md []Metadata) []error {
	var errs []error
	names := make(map[string]struct{})
	for i, m := range md {
		in := strings.Index(m.Name, ":")
		if in < 0 {
			if n := sort.SearchStrings(allowedMetadataNames[:], m.Name); n >= len(allowedMetadataNames) {
				errs = append(errs, &specerr.MetadataError{Index: i, Err: specerr.ErrMetadataName})
			}
		} else {
			var found bool
			space := m.Name[0:in]
			for _, ns := range v.m.Namespaces {
				if ns.Space == space {
					found = true
					break
				}
			}
			if !found {
				errs = append(errs, &specerr.MetadataError{Index: i, Err: specerr.ErrMetadataNamespace})
			}
		}
		if _, ok := names[m.Name]; ok {
			errs = append(errs, &specerr.MetadataError{Index: i, Err: specerr.ErrMetadataDuplicated})
		}
		names[m.Name] = struct{}{}
	}
	return errs
}

func (v *validator) validateResources(resources *Resources, path string) {
	assets := make(map[uint32]Asset)
	for i, r := range resources.Assets {
		id := r.Identify()
		if id == 0 {
			v.AddWarning(specerr.NewAsset(path, i, specerr.ErrMissingID))
		} else if _, ok := v.ids[validatorResource{path, id}]; ok {
			v.AddWarning(specerr.NewAsset(path, i, specerr.ErrDuplicatedID))
		}
		v.ids[validatorResource{path, id}] = struct{}{}
		assets[id] = r
		switch r := r.(type) {
		case *BaseMaterialsResource:
			for j, b := range r.Materials {
				if b.Name == "" {
					v.AddWarning(specerr.NewAsset(path, i, &specerr.BaseError{
						Index: j,
						Err:   &specerr.MissingFieldError{Name: attrName}},
					))
				}
			}
		}
	}
	for i, r := range resources.Objects {
		if r.ID == 0 {
			v.AddWarning(specerr.NewObject(path, i, specerr.ErrMissingID))
		} else if _, ok := v.ids[validatorResource{path, r.ID}]; ok {
			v.AddWarning(specerr.NewObject(path, i, specerr.ErrDuplicatedID))
		}
		v.ids[validatorResource{path, r.ID}] = struct{}{}
		if r.DefaultPIndex != 0 && r.DefaultPID == 0 {
			v.AddWarning(specerr.NewObject(path, i, &specerr.MissingFieldError{Name: attrPID}))
		}
		if r.Mesh != nil {
			if r.DefaultPID != 0 {
				if a, ok := assets[r.DefaultPID]; ok {
					if a, ok := a.(*BaseMaterialsResource); ok {
						if int(r.DefaultPIndex) > len(a.Materials) {
							v.AddWarning(specerr.NewObject(path, i, specerr.ErrIndexOutOfBounds))
						}
					}
				} else {
					v.AddWarning(specerr.NewObject(path, i, specerr.ErrMissingObject))
				}
			}
			v.validateMesh(r, path, i, assets)
		} else {
			if r.DefaultPID != 0 {
				v.AddWarning(specerr.NewObject(path, i, specerr.ErrComponentsPID))
			}
			v.validateComponents(r, path, i)
		}
	}
}

func (v *validator) validateMesh(r *Object, path string, index int, assets map[uint32]Asset) {
	nodeCount := uint32(len(r.Mesh.Nodes))
	for i, face := range r.Mesh.Faces {
		i0, i1, i2 := face.NodeIndices[0], face.NodeIndices[1], face.NodeIndices[2]
		if i0 == i1 || i0 == i2 || i1 == i2 {
			v.AddWarning(specerr.NewObject(path, index, &specerr.TriangleError{Index: i, Err: specerr.ErrDuplicatedIndices}))
		}
		if i0 >= nodeCount || i1 >= nodeCount || i2 >= nodeCount {
			v.AddWarning(specerr.NewObject(path, index, &specerr.TriangleError{Index: i, Err: specerr.ErrIndexOutOfBounds}))
		}
		if face.PID != 0 {
			if a, ok := assets[face.PID]; ok {
				if a, ok := a.(*BaseMaterialsResource); ok {
					if (face.PIndex[1] != face.PIndex[0] && face.PIndex[1] != 0) ||
						(face.PIndex[2] != face.PIndex[0] && face.PIndex[2] != 0) {
						v.AddWarning(specerr.NewObject(path, index, &specerr.TriangleError{Index: i, Err: specerr.ErrBaseMaterialGradient}))
					}
					if int(face.PIndex[0]) > len(a.Materials) {
						v.AddWarning(specerr.NewObject(path, index, &specerr.TriangleError{Index: i, Err: specerr.ErrIndexOutOfBounds}))
					}
				}
			} else {
				v.AddWarning(specerr.NewObject(path, index, &specerr.TriangleError{Index: i, Err: specerr.ErrMissingObject}))
			}
		}
	}
	switch r.ObjectType {
	case ObjectTypeModel, ObjectTypeSolidSupport:
		v.validateMeshCoherency(r, path, index)
	}
}

func (v *validator) validateMeshCoherency(r *Object, path string, index int) {
	if len(r.Mesh.Nodes) < 3 {
		v.AddWarning(specerr.NewObject(path, index, specerr.ErrEmptyTriangles))
	}
	if len(r.Mesh.Faces) <= 3 {
		v.AddWarning(specerr.NewObject(path, index, specerr.ErrEmptyTriangles))
	}

	var edgeCounter uint32
	pairMatching := newPairMatch()
	for _, face := range r.Mesh.Faces {
		for j := uint32(0); j < 3; j++ {
			n1, n2 := face.NodeIndices[j], face.NodeIndices[(j+1)%3]
			if _, ok := pairMatching.CheckMatch(n1, n2); !ok {
				pairMatching.AddMatch(n1, n2, edgeCounter)
				edgeCounter++
			}
		}
	}

	positive, negative := make([]uint32, edgeCounter), make([]uint32, edgeCounter)
	for _, face := range r.Mesh.Faces {
		for j := uint32(0); j < 3; j++ {
			n1, n2 := face.NodeIndices[j], face.NodeIndices[(j+1)%3]
			edgeIndex, _ := pairMatching.CheckMatch(n1, n2)
			if n1 <= n2 {
				positive[edgeIndex]++
			} else {
				negative[edgeIndex]++
			}
		}
	}

	for i := uint32(0); i < edgeCounter; i++ {
		if positive[i] != 1 || negative[i] != 1 {
			v.AddWarning(specerr.NewObject(path, index, specerr.ErrMeshConsistency))
		}
	}
}

func (v *validator) validateComponents(r *Object, path string, index int) {
	for j, c := range r.Components {
		if c.ObjectID == 0 {
			v.AddWarning(specerr.NewObject(path, index, &specerr.ComponentError{
				Index: j,
				Err:   &specerr.MissingFieldError{Name: attrObjectID}},
			))
		} else if _, ok := v.ids[validatorResource{path, c.ObjectID}]; !ok {
			v.AddWarning(specerr.NewObject(path, index, &specerr.ComponentError{
				Index: j,
				Err:   specerr.ErrMissingObject},
			))
		}
	}
}

func (v *validator) validateNamespaces() {
	for _, r := range v.m.RequiredExtensions {
		var found bool
		for _, ns := range v.m.Namespaces {
			if ns.Space == r {
				found = true
				break
			}
		}
		if !found {
			v.AddWarning(errors.New("go3mf: unsupported required extension"))
		}
	}
}

func (v *validator) validateRelationship(rels []Relationship, path string) {
	type partrel struct{ path, rel string }
	visitedParts := make(map[partrel]struct{})
	var hasPrintTicket bool
	for i, r := range rels {
		if r.Path == "" || r.Path[0] != '/' || strings.Contains(r.Path, "/.") {
			v.AddWarning(&specerr.RelationshipError{Path: path, Index: i, Err: specerr.ErrOPCPartName})
		} else {
			if _, ok := findAttachment(v.m.Attachments, r.Path); !ok {
				v.AddWarning(&specerr.RelationshipError{Path: path, Index: i, Err: specerr.ErrOPCRelTarget})
			}
			if _, ok := visitedParts[partrel{r.Path, r.Type}]; ok {
				v.AddWarning(&specerr.RelationshipError{Path: path, Index: i, Err: specerr.ErrOPCDuplicatedRel})
			}
			visitedParts[partrel{r.Path, r.Type}] = struct{}{}
		}
		switch r.Type {
		case RelTypePrintTicket:
			if a, ok := findAttachment(v.m.Attachments, r.Path); ok {
				if a.ContentType != ContentTypePrintTicket {
					v.AddWarning(&specerr.RelationshipError{Path: path, Index: i, Err: specerr.ErrOPCContentType})
				}
				if hasPrintTicket {
					v.AddWarning(&specerr.RelationshipError{Path: path, Index: i, Err: specerr.ErrOPCDuplicatedTicket})
				}
				hasPrintTicket = true
			}
		}
	}
}

func findAttachment(att []Attachment, path string) (*Attachment, bool) {
	for _, a := range att {
		if strings.EqualFold(a.Path, path) {
			return &a, true
		}
	}
	return nil, false
}
