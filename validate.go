package go3mf

import (
	"encoding/xml"
	"image/color"
	"sort"
	"strings"

	specerr "github.com/qmuntal/go3mf/errors"
)

type validator interface {
	Validate(*Model, string) []error
}

func (m *Model) sortedChilds() []string {
	s := make([]string, 0, len(m.Childs))
	for path := range m.Childs {
		s = append(s, path)
	}
	sort.Strings(s)
	return s
}

// Validate checks that the model is conformant with the 3MF spec.
func (m *Model) Validate() []error {
	errs := []error{}
	errs = append(errs, validateRelationship(m, m.RootRelationships, "")...)
	if err := m.validateNamespaces(); err != nil {
		errs = append(errs, err)
	}
	rootPath := m.PathOrDefault()
	sortedChilds := m.sortedChilds()
	for _, path := range sortedChilds {
		c := m.Childs[path]
		if path == rootPath {
			errs = append(errs, specerr.ErrOPCDuplicatedModelName)
		} else {
			errs = append(errs, validateRelationship(m, c.Relationships, path)...)
		}
	}
	errs = append(errs, validateRelationship(m, m.Relationships, rootPath)...)
	errs = append(errs, checkMetadadata(m, m.Metadata)...)

	for _, path := range sortedChilds {
		c := m.Childs[path]
		errs = append(errs, c.Resources.validate(m, path)...)
	}
	errs = append(errs, m.Resources.validate(m, rootPath)...)
	return append(errs, m.Build.Validate(m, rootPath)...)
}

func (ext ExtensionAttr) validate(m *Model, path string) []error {
	var errs []error
	for _, a := range ext {
		if a, ok := a.(validator); ok {
			errs = append(errs, a.Validate(m, path)...)
		}
	}
	return errs
}

func (item *Item) Validate(m *Model, path string) []error {
	var errs []error
	opath := item.ObjectPath(path)
	if item.ObjectID == 0 {
		errs = append(errs, &specerr.MissingFieldError{attrObjectID})
	} else if obj, ok := m.FindObject(opath, item.ObjectID); ok {
		if obj.ObjectType == ObjectTypeOther {
			errs = append(errs, specerr.ErrOtherItem)
		}
	} else {
		errs = append(errs, specerr.ErrMissingResource)
	}
	errs = append(errs, checkMetadadata(m, item.Metadata)...)
	errs = append(errs, item.ExtensionAttr.validate(m, path)...)
	return errs
}

func (b *Build) Validate(m *Model, path string) []error {
	var errs []error
	for _, err := range b.ExtensionAttr.validate(m, path) {
		errs = append(errs, &specerr.BuildError{Err: err})
	}
	for i, item := range b.Items {
		for _, err := range item.Validate(m, path) {
			errs = append(errs, specerr.NewItem(i, err))
		}
	}
	return errs
}

var allowedMetadataNames = [...]string{ // sorted
	"application", "copyright", "creationdate", "description", "designer",
	"licenseterms", "modificationdate", "rating", "title",
}

func (m *Metadata) Validate(model *Model) []error {
	if m.Name.Local == "" {
		return []error{&specerr.MissingFieldError{Name: attrName}}
	}
	var errs []error
	if m.Name.Space == "" {
		nm := strings.ToLower(m.Name.Local)
		n := sort.SearchStrings(allowedMetadataNames[:], nm)
		if n >= len(allowedMetadataNames) || allowedMetadataNames[n] != nm {
			errs = append(errs, specerr.ErrMetadataName)
		}
	} else {
		var found bool
		for _, ns := range model.Namespaces {
			if ns.Space == m.Name.Space {
				found = true
				break
			}
		}
		if !found {
			errs = append(errs, specerr.ErrMetadataNamespace)
		}
	}
	return errs
}

func checkMetadadata(model *Model, md []Metadata) []error {
	var errs []error
	names := make(map[xml.Name]struct{})
	for i, m := range md {
		for _, err := range m.Validate(model) {
			errs = append(errs, &specerr.IndexedError{Name: attrMetadata, Index: i, Err: err})
		}
		if _, ok := names[m.Name]; ok {
			errs = append(errs, &specerr.IndexedError{Name: attrMetadata, Index: i, Err: specerr.ErrMetadataDuplicated})
		}
		names[m.Name] = struct{}{}
	}
	return errs
}

func (r *BaseMaterialsResource) Validate(m *Model, path string) []error {
	var errs []error
	if r.ID == 0 {
		errs = append(errs, specerr.ErrMissingID)
	}
	if len(r.Materials) == 0 {
		errs = append(errs, specerr.ErrEmptyResourceProps)
	}
	for j, b := range r.Materials {
		if b.Name == "" {
			errs = append(errs, &specerr.IndexedError{Name: attrBase, Index: j, Err: &specerr.MissingFieldError{Name: attrName}})
		}
		if b.Color == (color.RGBA{}) {
			errs = append(errs, &specerr.IndexedError{Name: attrBase, Index: j, Err: &specerr.MissingFieldError{Name: attrDisplayColor}})
		}
	}
	return errs
}

func (res *Resources) validate(m *Model, path string) []error {
	var errs []error
	assets := make(map[uint32]struct{})
	for i, r := range res.Assets {
		id := r.Identify()
		if id != 0 {
			if _, ok := assets[id]; ok {
				errs = append(errs, specerr.NewAsset(path, i, r, specerr.ErrDuplicatedID))
			}
		}
		assets[id] = struct{}{}
		if r, ok := r.(validator); ok {
			for _, err := range r.Validate(m, path) {
				errs = append(errs, specerr.NewAsset(path, i, r, err))
			}
		}
	}
	for i, r := range res.Objects {
		if r.ID != 0 {
			if _, ok := assets[r.ID]; ok {
				errs = append(errs, specerr.NewObject(path, i, specerr.ErrDuplicatedID))
			}
		}
		assets[r.ID] = struct{}{}
		for _, err := range r.Validate(m, path) {
			errs = append(errs, specerr.NewObject(path, i, err))
		}
	}
	return errs
}

func (r *Object) Validate(m *Model, path string) []error {
	res, _ := m.FindResources(path)
	var errs []error
	if r.ID == 0 {
		errs = append(errs, specerr.ErrMissingID)
	}
	if r.DefaultPIndex != 0 && r.DefaultPID == 0 {
		errs = append(errs, &specerr.MissingFieldError{Name: attrPID})
	}
	if (r.Mesh != nil && len(r.Components) > 0) || (r.Mesh == nil && len(r.Components) == 0) {
		errs = append(errs, specerr.ErrInvalidObject)
	}
	errs = append(errs, r.ExtensionAttr.validate(m, path)...)
	if r.Mesh != nil {
		if r.DefaultPID != 0 {
			if a, ok := res.FindAsset(r.DefaultPID); ok {
				if a, ok := a.(propertyGroup); ok {
					if int(r.DefaultPIndex) >= a.Len() {
						errs = append(errs, specerr.ErrIndexOutOfBounds)
					}
				}
			} else {
				errs = append(errs, specerr.ErrMissingResource)
			}
		}
		errs = append(errs, r.validateMesh(m, res)...)
	}
	if len(r.Components) > 0 {
		if r.DefaultPID != 0 {
			errs = append(errs, specerr.ErrComponentsPID)
		}
		errs = append(errs, r.validateComponents(m, path)...)
	}
	return errs
}

func (r *Object) validateMesh(m *Model, res *Resources) []error {
	var errs []error
	switch r.ObjectType {
	case ObjectTypeModel, ObjectTypeSolidSupport:
		if len(r.Mesh.Nodes) < 3 {
			errs = append(errs, specerr.ErrInsufficientVertices)
		}
		if len(r.Mesh.Faces) <= 3 {
			errs = append(errs, specerr.ErrInsufficientTriangles)
		}
	}

	nodeCount := uint32(len(r.Mesh.Nodes))
	for i, face := range r.Mesh.Faces {
		i0, i1, i2 := face.NodeIndices[0], face.NodeIndices[1], face.NodeIndices[2]
		if i0 == i1 || i0 == i2 || i1 == i2 {
			errs = append(errs, &specerr.IndexedError{Name: attrTriangle, Index: i, Err: specerr.ErrDuplicatedIndices})
		}
		if i0 >= nodeCount || i1 >= nodeCount || i2 >= nodeCount {
			errs = append(errs, &specerr.IndexedError{Name: attrTriangle, Index: i, Err: specerr.ErrIndexOutOfBounds})
		}
		if face.PID != 0 {
			if face.PID == r.DefaultPID && face.PIndex[0] == r.DefaultPIndex &&
				face.PIndex[1] == r.DefaultPIndex && face.PIndex[2] == r.DefaultPIndex {
				continue
			}
			if a, ok := res.FindAsset(face.PID); ok {
				if a, ok := a.(propertyGroup); ok {
					l := a.Len()
					if int(face.PIndex[0]) >= l || int(face.PIndex[1]) >= l || int(face.PIndex[2]) >= l {
						errs = append(errs, &specerr.IndexedError{Name: attrTriangle, Index: i, Err: specerr.ErrIndexOutOfBounds})
					}
				}
			} else {
				errs = append(errs, &specerr.IndexedError{Name: attrTriangle, Index: i, Err: specerr.ErrMissingResource})
			}
		}
	}
	return errs
}

func (r *Object) validateComponents(m *Model, path string) []error {
	var errs []error
	for j, c := range r.Components {
		if c.ObjectID == 0 {
			errs = append(errs, &specerr.IndexedError{Name: attrComponent, Index: j, Err: &specerr.MissingFieldError{Name: attrObjectID}})
		} else if ref, ok := m.FindObject(c.ObjectPath(path), c.ObjectID); ok {
			if ref.ID == r.ID && c.ObjectPath(path) == path {
				errs = append(errs, &specerr.IndexedError{Name: attrComponent, Index: j, Err: specerr.ErrRecursiveComponent})
			}
		} else {
			errs = append(errs, &specerr.IndexedError{Name: attrComponent, Index: j, Err: specerr.ErrMissingResource})
		}
		for _, err := range c.ExtensionAttr.validate(m, path) {
			errs = append(errs, &specerr.IndexedError{Name: attrComponent, Index: j, Err: err})
		}
	}
	return errs
}

func (m *Model) validateNamespaces() error {
	for _, r := range m.RequiredExtensions {
		var found bool
		for _, ns := range m.Namespaces {
			if ns.Space == r {
				found = true
				break
			}
		}
		if !found {
			return specerr.ErrRequiredExt
		}
	}
	return nil
}

func validateRelationship(m *Model, rels []Relationship, path string) []error {
	var errs []error
	type partrel struct{ path, rel string }
	visitedParts := make(map[partrel]struct{})
	var hasPrintTicket bool
	for i, r := range rels {
		if r.Path == "" || r.Path[0] != '/' || strings.Contains(r.Path, "/.") {
			errs = append(errs, &specerr.RelationshipError{Path: path, Index: i, Err: specerr.ErrOPCPartName})
		} else {
			if _, ok := findAttachment(m.Attachments, r.Path); !ok {
				errs = append(errs, &specerr.RelationshipError{Path: path, Index: i, Err: specerr.ErrOPCRelTarget})
			}
			if _, ok := visitedParts[partrel{r.Path, r.Type}]; ok {
				errs = append(errs, &specerr.RelationshipError{Path: path, Index: i, Err: specerr.ErrOPCDuplicatedRel})
			}
			visitedParts[partrel{r.Path, r.Type}] = struct{}{}
		}
		switch r.Type {
		case RelTypePrintTicket:
			if a, ok := findAttachment(m.Attachments, r.Path); ok {
				if a.ContentType != ContentTypePrintTicket {
					errs = append(errs, &specerr.RelationshipError{Path: path, Index: i, Err: specerr.ErrOPCContentType})
				}
				if hasPrintTicket {
					errs = append(errs, &specerr.RelationshipError{Path: path, Index: i, Err: specerr.ErrOPCDuplicatedTicket})
				}
				hasPrintTicket = true
			}
		}
	}
	return errs
}

func findAttachment(att []Attachment, path string) (*Attachment, bool) {
	for _, a := range att {
		if strings.EqualFold(a.Path, path) {
			return &a, true
		}
	}
	return nil, false
}
