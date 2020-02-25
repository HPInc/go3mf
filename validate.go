package go3mf

import (
	"encoding/xml"
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
	warnings []error
	ids      map[validatorResource]interface{}
}

// Validate checks that the model is conformant with the 3MF spec.
func Validate(model *Model) []error {
	v := validator{m: model}
	v.Validate()
	return v.warnings
}

func (v *validator) AddWarning(err ...error) {
	v.warnings = append(v.warnings, err...)
}

func (v *validator) sortedChilds() []string {
	s := make([]string, 0, len(v.m.Childs))
	for path := range v.m.Childs {
		s = append(s, path)
	}
	sort.Strings(s)
	return s
}

func (v *validator) Validate() {
	v.ids = make(map[validatorResource]interface{})
	v.validateRelationship(v.m.RootRelationships, "")

	v.validateNamespaces()

	rootPath := v.m.Path
	if rootPath == "" {
		rootPath = DefaultPartModelName
	}
	sortedChilds := v.sortedChilds()
	for _, path := range sortedChilds {
		c := v.m.Childs[path]
		if path == rootPath {
			v.AddWarning(specerr.ErrOPCDuplicatedModelName)
		} else {
			v.validateRelationship(c.Relationships, path)
		}
	}
	v.validateRelationship(v.m.Relationships, rootPath)
	v.AddWarning(v.checkMetadadata(v.m.Metadata)...)

	for _, path := range sortedChilds {
		c := v.m.Childs[path]
		v.validateResources(&c.Resources, path)
	}
	v.validateResources(&v.m.Resources, rootPath)
	v.validateBuild(rootPath)
}

func (v *validator) validateBuild(rootPath string) {
	for i, item := range v.m.Build.Items {
		opath := item.ObjectPath(rootPath)
		if item.ObjectID == 0 {
			v.AddWarning(specerr.NewItem(i, &specerr.MissingFieldError{attrObjectID}))
		} else if r, ok := v.ids[validatorResource{opath, item.ObjectID}]; ok {
			if obj, ok := r.(*Object); ok {
				if obj.ObjectType == ObjectTypeOther {
					v.AddWarning(specerr.NewItem(i, specerr.ErrOtherItem))
				}
			} else {
				v.AddWarning(specerr.NewItem(i, specerr.ErrNonObject))
			}
		} else {
			v.AddWarning(specerr.NewItem(i, specerr.ErrMissingResource))
		}
		for _, err := range v.checkMetadadata(item.Metadata) {
			v.AddWarning(specerr.NewItem(i, err))
		}
	}
}

func (v *validator) checkMetadadata(md []Metadata) []error {
	var allowedMetadataNames = [...]string{ // sorted
		"application", "copyright", "creationdate", "description", "designer",
		"licenseterms", "modificationdate", "rating", "title",
	}
	var errs []error
	names := make(map[xml.Name]struct{})
	for i, m := range md {
		if m.Name.Local == "" {
			errs = append(errs, &specerr.MetadataError{Index: i, Err: &specerr.MissingFieldError{Name: attrName}})
			continue
		}
		if m.Name.Space == "" {
			nm := strings.ToLower(m.Name.Local)
			n := sort.SearchStrings(allowedMetadataNames[:], nm)
			if n >= len(allowedMetadataNames) || allowedMetadataNames[n] != nm {
				errs = append(errs, &specerr.MetadataError{Index: i, Err: specerr.ErrMetadataName})
			}
		} else {
			var found bool
			for _, ns := range v.m.Namespaces {
				if ns.Space == m.Name.Space {
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
		v.ids[validatorResource{path, id}] = r
		assets[id] = r
		switch r := r.(type) {
		case *BaseMaterialsResource:
			if len(r.Materials) == 0 {
				v.AddWarning(specerr.NewAsset(path, i, specerr.ErrEmptySlice))
			} else {
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
	}
	for i, r := range resources.Objects {
		if r.ID == 0 {
			v.AddWarning(specerr.NewObject(path, i, specerr.ErrMissingID))
		} else if _, ok := v.ids[validatorResource{path, r.ID}]; ok {
			v.AddWarning(specerr.NewObject(path, i, specerr.ErrDuplicatedID))
		}
		v.ids[validatorResource{path, r.ID}] = r
		if r.DefaultPIndex != 0 && r.DefaultPID == 0 {
			v.AddWarning(specerr.NewObject(path, i, &specerr.MissingFieldError{Name: attrPID}))
		}
		if (r.Mesh != nil && len(r.Components) > 0) || (r.Mesh == nil && len(r.Components) == 0) {
			v.AddWarning(specerr.NewObject(path, i, specerr.ErrInvalidObject))
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
					v.AddWarning(specerr.NewObject(path, i, specerr.ErrMissingResource))
				}
			}
			v.validateMesh(r, path, i, assets)
		}
		if len(r.Components) > 0 {
			if r.DefaultPID != 0 {
				v.AddWarning(specerr.NewObject(path, i, specerr.ErrComponentsPID))
			}
			v.validateComponents(r, path, i)
		}
	}
}

func (v *validator) validateMesh(r *Object, path string, index int, assets map[uint32]Asset) {
	switch r.ObjectType {
	case ObjectTypeModel, ObjectTypeSolidSupport:
		if len(r.Mesh.Nodes) < 3 {
			v.AddWarning(specerr.NewObject(path, index, specerr.ErrInsufficientVertices))
		}
		if len(r.Mesh.Faces) <= 3 {
			v.AddWarning(specerr.NewObject(path, index, specerr.ErrInsufficientTriangles))
		}
	}

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
				v.AddWarning(specerr.NewObject(path, index, &specerr.TriangleError{Index: i, Err: specerr.ErrMissingResource}))
			}
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
		} else if ref, ok := v.ids[validatorResource{c.ObjectPath(path), c.ObjectID}]; ok {
			if ref == r {
				v.AddWarning(specerr.NewObject(path, index, &specerr.ComponentError{
					Index: j,
					Err:   specerr.ErrRecursiveComponent},
				))
			} else if _, ok := ref.(*Object); !ok {
				v.AddWarning(specerr.NewObject(path, index, &specerr.ComponentError{
					Index: j,
					Err:   specerr.ErrNonObject},
				))
			}
		} else {
			v.AddWarning(specerr.NewObject(path, index, &specerr.ComponentError{
				Index: j,
				Err:   specerr.ErrMissingResource},
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
			v.AddWarning(specerr.ErrRequiredExt)
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
