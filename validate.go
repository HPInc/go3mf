package go3mf

import (
	"encoding/xml"
	"image/color"
	"sort"
	"strings"

	"github.com/qmuntal/go3mf/errors"
)

func (m *Model) sortedChilds() []string {
	if len(m.Childs) == 0 {
		return nil
	}
	s := make([]string, 0, len(m.Childs))
	for path := range m.Childs {
		s = append(s, path)
	}
	sort.Strings(s)
	return s
}

// Validate checks that the model is conformant with the 3MF spec.
func (m *Model) Validate() error {
	var errs error
	errs = errors.Append(errs, validateRelationship(m, m.RootRelationships, ""))
	errs = errors.Append(errs, m.validateNamespaces())
	rootPath := m.PathOrDefault()
	sortedChilds := m.sortedChilds()
	for _, path := range sortedChilds {
		c := m.Childs[path]
		if path == rootPath {
			errs = errors.Append(errs, errors.ErrOPCDuplicatedModelName)
		} else {
			errs = errors.Append(errs, validateRelationship(m, c.Relationships, path))
		}
	}

	errs = errors.Append(errs, validateRelationship(m, m.Relationships, rootPath))
	errs = errors.Append(errs, checkMetadadata(m, m.Metadata))

	for _, ext := range m.Specs {
		if ext, ok := ext.(SpecValidator); ok {
			errs = errors.Append(errs, ext.ValidateModel(m))
		}
	}

	for _, path := range sortedChilds {
		c := m.Childs[path]
		err := c.Resources.validate(m, path)
		if err != nil {
			errs = errors.Append(errs, errors.WrapPath(err, c.Resources, path))
		}
	}
	err := m.Resources.validate(m, rootPath)
	if err != nil {
		errs = errors.Append(errs, errors.Wrap(err, m.Resources))
	}
	err = m.Build.validate(m, rootPath)
	if err != nil {
		errs = errors.Append(errs, errors.Wrap(err, m.Build))
	}
	return errs
}

func (item *Item) validate(m *Model, path string) error {
	var errs error
	opath := item.ObjectPath(path)
	if item.ObjectID == 0 {
		errs = errors.Append(errs, &errors.MissingFieldError{Name: attrObjectID})
	} else if obj, ok := m.FindObject(opath, item.ObjectID); ok {
		if obj.ObjectType == ObjectTypeOther {
			errs = errors.Append(errs, errors.ErrOtherItem)
		}
	} else {
		errs = errors.Append(errs, errors.ErrMissingResource)
	}
	return errors.Append(errs, checkMetadadata(m, item.Metadata))
}

func (b *Build) validate(m *Model, path string) error {
	var errs error
	for i, item := range b.Items {
		err := item.validate(m, path)
		if err != nil {
			errs = errors.Append(errs, errors.WrapIndex(err, item, i))
		}
	}
	return errs
}

var allowedMetadataNames = [...]string{ // sorted
	"application", "copyright", "creationdate", "description", "designer",
	"licenseterms", "modificationdate", "rating", "title",
}

func (m *Metadata) validate(model *Model) error {
	if m.Name.Local == "" {
		return &errors.MissingFieldError{Name: attrName}
	}
	var errs error
	if m.Name.Space == "" {
		nm := strings.ToLower(m.Name.Local)
		n := sort.SearchStrings(allowedMetadataNames[:], nm)
		if n >= len(allowedMetadataNames) || allowedMetadataNames[n] != nm {
			errs = errors.Append(errs, errors.ErrMetadataName)
		}
	} else {
		if _, ok := model.Specs[m.Name.Space]; !ok {
			errs = errors.Append(errs, errors.ErrMetadataNamespace)
		}
	}
	return errs
}

func checkMetadadata(model *Model, md []Metadata) error {
	var errs error
	names := make(map[xml.Name]struct{})
	for i, m := range md {
		err := m.validate(model)
		errs = errors.Append(errs, errors.WrapIndex(err, m, i))
		if _, ok := names[m.Name]; ok {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrMetadataDuplicated, m, i))
		}
		names[m.Name] = struct{}{}
	}
	return errs
}

func (r *BaseMaterials) Validate(m *Model, path string) error {
	var errs error
	if r.ID == 0 {
		errs = errors.Append(errs, errors.ErrMissingID)
	}
	if len(r.Materials) == 0 {
		errs = errors.Append(errs, errors.ErrEmptyResourceProps)
	}
	for j, b := range r.Materials {
		if b.Name == "" {
			errs = errors.Append(errs, errors.WrapIndex(&errors.MissingFieldError{Name: attrName}, b, j))
		}
		if b.Color == (color.RGBA{}) {
			errs = errors.Append(errs, errors.WrapIndex(&errors.MissingFieldError{Name: attrDisplayColor}, b, j))
		}
	}
	return errs
}

func (res *Resources) validate(m *Model, path string) error {
	var errs error
	assets := make(map[uint32]struct{})
	for i, r := range res.Assets {
		var aErrs error
		id := r.Identify()
		if id != 0 {
			if _, ok := assets[id]; ok {
				aErrs = errors.Append(aErrs, errors.ErrDuplicatedID)
			}
		}
		assets[id] = struct{}{}

		if r, ok := r.(*BaseMaterials); ok {
			aErrs = errors.Append(aErrs, r.Validate(m, path))
		}

		for _, ext := range m.Specs {
			if ext, ok := ext.(SpecValidator); ok {
				aErrs = errors.Append(aErrs, ext.ValidateAsset(m, path, r))
			}
		}
		errs = errors.Append(errs, errors.WrapIndex(aErrs, r, i))
	}
	for i, r := range res.Objects {
		if r.ID != 0 {
			if _, ok := assets[r.ID]; ok {
				errs = errors.Append(errs, errors.WrapIndex(errors.ErrDuplicatedID, r, i))
			}
		}
		assets[r.ID] = struct{}{}
		err := r.Validate(m, path)
		errs = errors.Append(errs, errors.WrapIndex(err, r, i))
	}
	return errs
}

func (r *Object) Validate(m *Model, path string) error {
	res, _ := m.FindResources(path)
	var errs error
	if r.ID == 0 {
		errs = errors.Append(errs, errors.ErrMissingID)
	}
	if r.DefaultPIndex != 0 && r.DefaultPID == 0 {
		errs = errors.Append(errs, &errors.MissingFieldError{Name: attrPID})
	}
	if (r.Mesh != nil && len(r.Components) > 0) || (r.Mesh == nil && len(r.Components) == 0) {
		errs = errors.Append(errs, errors.ErrInvalidObject)
	}
	if r.Mesh != nil {
		if r.DefaultPID != 0 {
			if a, ok := res.FindAsset(r.DefaultPID); ok {
				if a, ok := a.(PropertyGroup); ok {
					if int(r.DefaultPIndex) >= a.Len() {
						errs = errors.Append(errs, errors.ErrIndexOutOfBounds)
					}
				}
			} else {
				errs = errors.Append(errs, errors.ErrMissingResource)
			}
		}
		err := r.validateMesh(m, path)
		if err != nil {
			errs = errors.Append(errs, errors.Wrap(err, r.Mesh))
		}
	}
	if len(r.Components) > 0 {
		if r.DefaultPID != 0 {
			errs = errors.Append(errs, errors.ErrComponentsPID)
		}
		errs = errors.Append(errs, r.validateComponents(m, path))
	}
	for _, ext := range m.Specs {
		if ext, ok := ext.(SpecValidator); ok {
			errs = errors.Append(errs, ext.ValidateObject(m, path, r))
		}
	}
	return errs
}

func (r *Object) validateMesh(m *Model, path string) error {
	res, _ := m.FindResources(path)
	var errs error
	switch r.ObjectType {
	case ObjectTypeModel, ObjectTypeSolidSupport:
		if len(r.Mesh.Vertices) < 3 {
			errs = errors.Append(errs, errors.ErrInsufficientVertices)
		}
		if len(r.Mesh.Triangles) <= 3 && len(r.Mesh.Any) == 0 {
			errs = errors.Append(errs, errors.ErrInsufficientTriangles)
		}
	}

	nodeCount := uint32(len(r.Mesh.Vertices))
	for i, face := range r.Mesh.Triangles {
		i0, i1, i2 := face.Indices[0], face.Indices[1], face.Indices[2]
		if i0 == i1 || i0 == i2 || i1 == i2 {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrDuplicatedIndices, face, i))
		}
		if i0 >= nodeCount || i1 >= nodeCount || i2 >= nodeCount {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrIndexOutOfBounds, face, i))
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
						errs = errors.Append(errs, errors.WrapIndex(errors.ErrIndexOutOfBounds, face, i))
					}
				}
			} else {
				errs = errors.Append(errs, errors.WrapIndex(errors.ErrMissingResource, face, i))
			}
		}
	}
	return errs
}

func (r *Object) validateComponents(m *Model, path string) error {
	var errs error
	for j, c := range r.Components {
		if c.ObjectID == 0 {
			errs = errors.Append(errs, errors.WrapIndex(&errors.MissingFieldError{Name: attrObjectID}, c, j))
		} else if ref, ok := m.FindObject(c.ObjectPath(path), c.ObjectID); ok {
			if ref.ID == r.ID && c.ObjectPath(path) == path {
				errs = errors.Append(errs, errors.WrapIndex(errors.ErrRecursion, c, j))
			}
		} else {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrMissingResource, c, j))
		}
	}
	return errs
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
	var errs error
	type partrel struct{ path, rel string }
	visitedParts := make(map[partrel]struct{})
	var hasPrintTicket bool
	for i, r := range rels {
		if r.Path == "" || r.Path[0] != '/' || strings.Contains(r.Path, "/.") {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrOPCPartName, r, i))
		} else {
			if _, ok := findAttachment(m.Attachments, r.Path); !ok {
				errs = errors.Append(errs, errors.WrapIndex(errors.ErrOPCRelTarget, r, i))
			}
			if _, ok := visitedParts[partrel{r.Path, r.Type}]; ok {
				errs = errors.Append(errs, errors.WrapIndex(errors.ErrOPCDuplicatedRel, r, i))
			}
			visitedParts[partrel{r.Path, r.Type}] = struct{}{}
		}
		switch r.Type {
		case RelTypePrintTicket:
			if a, ok := findAttachment(m.Attachments, r.Path); ok {
				if a.ContentType != ContentTypePrintTicket {
					errs = errors.Append(errs, errors.WrapIndex(errors.ErrOPCContentType, r, i))
				}
				if hasPrintTicket {
					errs = errors.Append(errs, errors.WrapIndex(errors.ErrOPCDuplicatedTicket, r, i))
				}
				hasPrintTicket = true
			}
		}
	}
	if path != "" && errs != nil {
		for _, err := range errs.(*errors.List).Errors {
			if err, ok := err.(*errors.Error); ok {
				err.Path = path
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
