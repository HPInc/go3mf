package go3mf

import (
	"encoding/xml"
	"image/color"
	"sort"
	"strings"

	specerr "github.com/qmuntal/go3mf/errors"
)

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
	rootPath := m.PathOrDefault()
	if err := m.validateNamespaces(); err != nil {
		errs = append(errs, err)
	}
	sortedChilds := m.sortedChilds()
	for _, path := range sortedChilds {
		c := m.Childs[path]
		if path == rootPath {
			errs = append(errs, specerr.ErrOPCDuplicatedModelName)
		} else {
			for _, err := range validateRelationship(m, c.Relationships, path) {
				if err, ok := err.(*specerr.Error); ok {
					err.Path = path
				}
				errs = append(errs, err)
			}
		}
	}

	errs = append(errs, validateRelationship(m, m.Relationships, rootPath)...)
	errs = append(errs, checkMetadadata(m, m.Metadata)...)

	for _, ext := range m.ExtensionSpecs {
		if ext, ok := ext.(interface {
			ValidateModel(*Model) []error
		}); ok {
			errs = append(errs, ext.ValidateModel(m)...)
		}
	}

	for _, path := range sortedChilds {
		c := m.Childs[path]
		for _, err := range c.Resources.Validate(m, path) {
			errs = append(errs, specerr.NewPath(c.Resources, path, err))
		}
	}
	for _, err := range m.Resources.Validate(m, rootPath) {
		errs = append(errs, specerr.New(m.Resources, err))
	}
	for _, err := range m.Build.Validate(m, rootPath) {
		errs = append(errs, specerr.New(m.Build, err))
	}
	return errs
}

func (item *Item) Validate(m *Model, path string) []error {
	var errs []error
	opath := item.ObjectPath(path)
	if item.ObjectID == 0 {
		errs = append(errs, &specerr.MissingFieldError{Name: attrObjectID})
	} else if obj, ok := m.FindObject(opath, item.ObjectID); ok {
		if obj.ObjectType == ObjectTypeOther {
			errs = append(errs, specerr.ErrOtherItem)
		}
	} else {
		errs = append(errs, specerr.ErrMissingResource)
	}
	errs = append(errs, checkMetadadata(m, item.Metadata)...)
	for _, ext := range m.ExtensionSpecs {
		if ext, ok := ext.(interface {
			ValidateItem(*Item, []error) []error
		}); ok {
			errs = ext.ValidateItem(item, errs)
		}
	}
	return errs
}

func (b *Build) Validate(m *Model, path string) []error {
	var errs []error
	for i, item := range b.Items {
		for _, err := range item.Validate(m, path) {
			errs = append(errs, specerr.NewIndexed(item, i, err))
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
			errs = append(errs, specerr.NewIndexed(m, i, err))
		}
		if _, ok := names[m.Name]; ok {
			errs = append(errs, specerr.NewIndexed(m, i, specerr.ErrMetadataDuplicated))
		}
		names[m.Name] = struct{}{}
	}
	return errs
}

func (r *BaseMaterials) Validate(m *Model, path string) []error {
	var errs []error
	if r.ID == 0 {
		errs = append(errs, specerr.ErrMissingID)
	}
	if len(r.Materials) == 0 {
		errs = append(errs, specerr.ErrEmptyResourceProps)
	}
	for j, b := range r.Materials {
		if b.Name == "" {
			errs = append(errs, specerr.NewIndexed(b, j, &specerr.MissingFieldError{Name: attrName}))
		}
		if b.Color == (color.RGBA{}) {
			errs = append(errs, specerr.NewIndexed(b, j, &specerr.MissingFieldError{Name: attrDisplayColor}))
		}
	}
	return errs
}

func (res *Resources) Validate(m *Model, path string) []error {
	var errs []error
	assets := make(map[uint32]struct{})
	for i, r := range res.Assets {
		var aErrs []error
		id := r.Identify()
		if id != 0 {
			if _, ok := assets[id]; ok {
				aErrs = append(aErrs, specerr.ErrDuplicatedID)
			}
		}
		assets[id] = struct{}{}

		if r, ok := r.(*BaseMaterials); ok {
			aErrs = append(aErrs, r.Validate(m, path)...)
		}

		for _, ext := range m.ExtensionSpecs {
			if ext, ok := ext.(interface {
				ValidateAsset(*Model, string, Asset) []error
			}); ok {
				aErrs = append(aErrs, ext.ValidateAsset(m, path, r)...)
			}
		}
		for _, err := range aErrs {
			errs = append(errs, specerr.NewIndexed(r, i, err))
		}
	}
	for i, r := range res.Objects {
		if r.ID != 0 {
			if _, ok := assets[r.ID]; ok {
				errs = append(errs, specerr.NewIndexed(r, i, specerr.ErrDuplicatedID))
			}
		}
		assets[r.ID] = struct{}{}
		for _, err := range r.Validate(m, path) {
			errs = append(errs, specerr.NewIndexed(r, i, err))
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
		for _, err := range r.validateMesh(m, path) {
			errs = append(errs, specerr.New(r.Mesh, err))
		}
	}
	if len(r.Components) > 0 {
		if r.DefaultPID != 0 {
			errs = append(errs, specerr.ErrComponentsPID)
		}
		errs = append(errs, r.validateComponents(m, path)...)
	}
	for _, ext := range m.ExtensionSpecs {
		if ext, ok := ext.(interface {
			ValidateObject(*Model, string, *Object) []error
		}); ok {
			errs = append(errs, ext.ValidateObject(m, path, r)...)
		}
	}
	return errs
}

func (r *Object) validateMesh(m *Model, path string) []error {
	res, _ := m.FindResources(path)
	var errs []error
	switch r.ObjectType {
	case ObjectTypeModel, ObjectTypeSolidSupport:
		if len(r.Mesh.Nodes) < 3 {
			errs = append(errs, specerr.ErrInsufficientVertices)
		}
		if len(r.Mesh.Faces) <= 3 && len(r.Mesh.Extension) == 0 {
			errs = append(errs, specerr.ErrInsufficientTriangles)
		}
	}

	nodeCount := uint32(len(r.Mesh.Nodes))
	for i, face := range r.Mesh.Faces {
		i0, i1, i2 := face.NodeIndices[0], face.NodeIndices[1], face.NodeIndices[2]
		if i0 == i1 || i0 == i2 || i1 == i2 {
			errs = append(errs, specerr.NewIndexed(face, i, specerr.ErrDuplicatedIndices))
		}
		if i0 >= nodeCount || i1 >= nodeCount || i2 >= nodeCount {
			errs = append(errs, specerr.NewIndexed(face, i, specerr.ErrIndexOutOfBounds))
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
						errs = append(errs, specerr.NewIndexed(face, i, specerr.ErrIndexOutOfBounds))
					}
				}
			} else {
				errs = append(errs, specerr.NewIndexed(face, i, specerr.ErrMissingResource))
			}
		}
	}
	return errs
}

func (r *Object) validateComponents(m *Model, path string) []error {
	var errs []error
	for j, c := range r.Components {
		if c.ObjectID == 0 {
			errs = append(errs, specerr.NewIndexed(c, j, &specerr.MissingFieldError{Name: attrObjectID}))
		} else if ref, ok := m.FindObject(c.ObjectPath(path), c.ObjectID); ok {
			if ref.ID == r.ID && c.ObjectPath(path) == path {
				errs = append(errs, specerr.NewIndexed(c, j, specerr.ErrRecursiveComponent))
			}
		} else {
			errs = append(errs, specerr.NewIndexed(c, j, specerr.ErrMissingResource))
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
			errs = append(errs, specerr.NewIndexed(r, i, specerr.ErrOPCPartName))
		} else {
			if _, ok := findAttachment(m.Attachments, r.Path); !ok {
				errs = append(errs, specerr.NewIndexed(r, i, specerr.ErrOPCRelTarget))
			}
			if _, ok := visitedParts[partrel{r.Path, r.Type}]; ok {
				errs = append(errs, specerr.NewIndexed(r, i, specerr.ErrOPCDuplicatedRel))
			}
			visitedParts[partrel{r.Path, r.Type}] = struct{}{}
		}
		switch r.Type {
		case RelTypePrintTicket:
			if a, ok := findAttachment(m.Attachments, r.Path); ok {
				if a.ContentType != ContentTypePrintTicket {
					errs = append(errs, specerr.NewIndexed(r, i, specerr.ErrOPCContentType))
				}
				if hasPrintTicket {
					errs = append(errs, specerr.NewIndexed(r, i, specerr.ErrOPCDuplicatedTicket))
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
