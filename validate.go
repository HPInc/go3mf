package go3mf

import (
	"encoding/xml"
	"image/color"
	"sort"
	"strings"

	"github.com/qmuntal/go3mf/errors"
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
func (m *Model) Validate() error {
	errs := new(errors.ErrorList)
	errs.Append(validateRelationship(m, m.RootRelationships, ""))
	errs.Append(m.validateNamespaces())
	rootPath := m.PathOrDefault()
	sortedChilds := m.sortedChilds()
	for _, path := range sortedChilds {
		c := m.Childs[path]
		if path == rootPath {
			errs.Append(errors.ErrOPCDuplicatedModelName)
		} else {
			errs.Append(validateRelationship(m, c.Relationships, path))
		}
	}

	errs.Append(validateRelationship(m, m.Relationships, rootPath))
	errs.Append(checkMetadadata(m, m.Metadata))

	for _, ext := range m.Specs {
		if ext, ok := ext.(SpecValidator); ok {
			errs.Append(ext.ValidateModel(m))
		}
	}

	for _, path := range sortedChilds {
		c := m.Childs[path]
		err := c.Resources.validate(m, path)
		errs.Append(errors.NewPath(c.Resources, path, err))
	}
	err := m.Resources.validate(m, rootPath)
	errs.Append(errors.New(m.Resources, err))

	err = m.Build.validate(m, rootPath)
	errs.Append(errors.New(m.Build, err))
	return errs.ErrorOrNil()
}

func (item *Item) validate(m *Model, path string) error {
	errs := new(errors.ErrorList)
	opath := item.ObjectPath(path)
	if item.ObjectID == 0 {
		errs.Append(&errors.MissingFieldError{Name: attrObjectID})
	} else if obj, ok := m.FindObject(opath, item.ObjectID); ok {
		if obj.ObjectType == ObjectTypeOther {
			errs.Append(errors.ErrOtherItem)
		}
	} else {
		errs.Append(errors.ErrMissingResource)
	}
	errs.Append(checkMetadadata(m, item.Metadata))
	return errs.ErrorOrNil()
}

func (b *Build) validate(m *Model, path string) error {
	errs := new(errors.ErrorList)
	for i, item := range b.Items {
		err := item.validate(m, path)
		errs.Append(errors.NewIndexed(item, i, err))
	}
	return errs.ErrorOrNil()
}

var allowedMetadataNames = [...]string{ // sorted
	"application", "copyright", "creationdate", "description", "designer",
	"licenseterms", "modificationdate", "rating", "title",
}

func (m *Metadata) validate(model *Model) error {
	if m.Name.Local == "" {
		return &errors.MissingFieldError{Name: attrName}
	}
	errs := new(errors.ErrorList)
	if m.Name.Space == "" {
		nm := strings.ToLower(m.Name.Local)
		n := sort.SearchStrings(allowedMetadataNames[:], nm)
		if n >= len(allowedMetadataNames) || allowedMetadataNames[n] != nm {
			errs.Append(errors.ErrMetadataName)
		}
	} else {
		if _, ok := model.Specs[m.Name.Space]; !ok {
			errs.Append(errors.ErrMetadataNamespace)
		}
	}
	return errs.ErrorOrNil()
}

func checkMetadadata(model *Model, md []Metadata) error {
	errs := new(errors.ErrorList)
	names := make(map[xml.Name]struct{})
	for i, m := range md {
		err := m.validate(model)
		errs.Append(errors.NewIndexed(m, i, err))
		if _, ok := names[m.Name]; ok {
			errs.Append(errors.NewIndexed(m, i, errors.ErrMetadataDuplicated))
		}
		names[m.Name] = struct{}{}
	}
	return errs.ErrorOrNil()
}

func (r *BaseMaterials) Validate(m *Model, path string) error {
	errs := new(errors.ErrorList)
	if r.ID == 0 {
		errs.Append(errors.ErrMissingID)
	}
	if len(r.Materials) == 0 {
		errs.Append(errors.ErrEmptyResourceProps)
	}
	for j, b := range r.Materials {
		if b.Name == "" {
			errs.Append(errors.NewIndexed(b, j, &errors.MissingFieldError{Name: attrName}))
		}
		if b.Color == (color.RGBA{}) {
			errs.Append(errors.NewIndexed(b, j, &errors.MissingFieldError{Name: attrDisplayColor}))
		}
	}
	return errs.ErrorOrNil()
}

func (res *Resources) validate(m *Model, path string) error {
	errs := new(errors.ErrorList)
	assets := make(map[uint32]struct{})
	for i, r := range res.Assets {
		aErrs := new(errors.ErrorList)
		id := r.Identify()
		if id != 0 {
			if _, ok := assets[id]; ok {
				aErrs.Append(errors.ErrDuplicatedID)
			}
		}
		assets[id] = struct{}{}

		if r, ok := r.(*BaseMaterials); ok {
			aErrs.Append(r.Validate(m, path))
		}

		for _, ext := range m.Specs {
			if ext, ok := ext.(SpecValidator); ok {
				aErrs.Append(ext.ValidateAsset(m, path, r))
			}
		}
		errs.Append(errors.NewIndexed(r, i, aErrs))
	}
	for i, r := range res.Objects {
		if r.ID != 0 {
			if _, ok := assets[r.ID]; ok {
				errs.Append(errors.NewIndexed(r, i, errors.ErrDuplicatedID))
			}
		}
		assets[r.ID] = struct{}{}
		err := r.Validate(m, path)
		errs.Append(errors.NewIndexed(r, i, err))
	}
	return errs.ErrorOrNil()
}

func (r *Object) Validate(m *Model, path string) error {
	res, _ := m.FindResources(path)
	errs := new(errors.ErrorList)
	if r.ID == 0 {
		errs.Append(errors.ErrMissingID)
	}
	if r.DefaultPIndex != 0 && r.DefaultPID == 0 {
		errs.Append(&errors.MissingFieldError{Name: attrPID})
	}
	if (r.Mesh != nil && len(r.Components) > 0) || (r.Mesh == nil && len(r.Components) == 0) {
		errs.Append(errors.ErrInvalidObject)
	}
	if r.Mesh != nil {
		if r.DefaultPID != 0 {
			if a, ok := res.FindAsset(r.DefaultPID); ok {
				if a, ok := a.(PropertyGroup); ok {
					if int(r.DefaultPIndex) >= a.Len() {
						errs.Append(errors.ErrIndexOutOfBounds)
					}
				}
			} else {
				errs.Append(errors.ErrMissingResource)
			}
		}
		err := r.validateMesh(m, path)
		errs.Append(errors.New(r.Mesh, err))
	}
	if len(r.Components) > 0 {
		if r.DefaultPID != 0 {
			errs.Append(errors.ErrComponentsPID)
		}
		errs.Append(r.validateComponents(m, path))
	}
	for _, ext := range m.Specs {
		if ext, ok := ext.(SpecValidator); ok {
			errs.Append(ext.ValidateObject(m, path, r))
		}
	}
	return errs.ErrorOrNil()
}

func (r *Object) validateMesh(m *Model, path string) error {
	res, _ := m.FindResources(path)
	errs := new(errors.ErrorList)
	switch r.ObjectType {
	case ObjectTypeModel, ObjectTypeSolidSupport:
		if len(r.Mesh.Vertices) < 3 {
			errs.Append(errors.ErrInsufficientVertices)
		}
		if len(r.Mesh.Triangles) <= 3 && len(r.Mesh.Any) == 0 {
			errs.Append(errors.ErrInsufficientTriangles)
		}
	}

	nodeCount := uint32(len(r.Mesh.Vertices))
	for i, face := range r.Mesh.Triangles {
		i0, i1, i2 := face.Indices[0], face.Indices[1], face.Indices[2]
		if i0 == i1 || i0 == i2 || i1 == i2 {
			errs.Append(errors.NewIndexed(face, i, errors.ErrDuplicatedIndices))
		}
		if i0 >= nodeCount || i1 >= nodeCount || i2 >= nodeCount {
			errs.Append(errors.NewIndexed(face, i, errors.ErrIndexOutOfBounds))
		}
		if face.PID != 0 {
			if face.PID == r.DefaultPID && face.PIndices[0] == r.DefaultPIndex &&
				face.PIndices[1] == r.DefaultPIndex && face.PIndices[2] == r.DefaultPIndex {
				continue
			}
			if a, ok := res.FindAsset(face.PID); ok {
				if a, ok := a.(PropertyGroup); ok {
					l := a.Len()
					if int(face.PIndices[0]) >= l || int(face.PIndices[1]) >= l || int(face.PIndices[2]) >= l {
						errs.Append(errors.NewIndexed(face, i, errors.ErrIndexOutOfBounds))
					}
				}
			} else {
				errs.Append(errors.NewIndexed(face, i, errors.ErrMissingResource))
			}
		}
	}
	return errs.ErrorOrNil()
}

func (r *Object) validateComponents(m *Model, path string) error {
	errs := new(errors.ErrorList)
	for j, c := range r.Components {
		if c.ObjectID == 0 {
			errs.Append(errors.NewIndexed(c, j, &errors.MissingFieldError{Name: attrObjectID}))
		} else if ref, ok := m.FindObject(c.ObjectPath(path), c.ObjectID); ok {
			if ref.ID == r.ID && c.ObjectPath(path) == path {
				errs.Append(errors.NewIndexed(c, j, errors.ErrRecursion))
			}
		} else {
			errs.Append(errors.NewIndexed(c, j, errors.ErrMissingResource))
		}
	}
	return errs.ErrorOrNil()
}

func (m *Model) validateNamespaces() error {
	for _, ext := range m.Specs {
		if ext.Required() {
			if _, ok := ext.(*UnknownSpec); ok {
				return errors.ErrRequiredExt
			}
		}
	}
	return nil
}

func validateRelationship(m *Model, rels []Relationship, path string) error {
	errs := new(errors.ErrorList)
	type partrel struct{ path, rel string }
	visitedParts := make(map[partrel]struct{})
	var hasPrintTicket bool
	for i, r := range rels {
		if r.Path == "" || r.Path[0] != '/' || strings.Contains(r.Path, "/.") {
			errs.Append(errors.NewIndexed(r, i, errors.ErrOPCPartName))
		} else {
			if _, ok := findAttachment(m.Attachments, r.Path); !ok {
				errs.Append(errors.NewIndexed(r, i, errors.ErrOPCRelTarget))
			}
			if _, ok := visitedParts[partrel{r.Path, r.Type}]; ok {
				errs.Append(errors.NewIndexed(r, i, errors.ErrOPCDuplicatedRel))
			}
			visitedParts[partrel{r.Path, r.Type}] = struct{}{}
		}
		switch r.Type {
		case RelTypePrintTicket:
			if a, ok := findAttachment(m.Attachments, r.Path); ok {
				if a.ContentType != ContentTypePrintTicket {
					errs.Append(errors.NewIndexed(r, i, errors.ErrOPCContentType))
				}
				if hasPrintTicket {
					errs.Append(errors.NewIndexed(r, i, errors.ErrOPCDuplicatedTicket))
				}
				hasPrintTicket = true
			}
		}
	}
	if path != "" {
		for _, err := range errs.Errors {
			if err, ok := err.(*errors.Error); ok {
				err.Path = path
			}
		}
	}
	return errs.ErrorOrNil()
}

func findAttachment(att []Attachment, path string) (*Attachment, bool) {
	for _, a := range att {
		if strings.EqualFold(a.Path, path) {
			return &a, true
		}
	}
	return nil, false
}
