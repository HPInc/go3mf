package materials

import (
	"image/color"
	"strings"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func (e *Spec) ValidateModel(_ *go3mf.Model) error {
	return nil
}

func (e *Spec) ValidateObject(_ *go3mf.Model, _ string, _ *go3mf.Object) error {
	return nil
}

func (e *Spec) ValidateAsset(m *go3mf.Model, path string, r go3mf.Asset) error {
	errs := new(errors.ErrorList)
	switch r := r.(type) {
	case *ColorGroup:
		e.validateColorGroup(path, r, errs)
	case *Texture2DGroup:
		e.validateTexture2DGroup(m, path, r, errs)
	case *Texture2D:
		e.validateTexture2D(m, path, r, errs)
	case *MultiProperties:
		e.validateMultiProps(m, path, r, errs)
	case *CompositeMaterials:
		e.validateCompositeMat(m, path, r, errs)
	}
	return errs.ErrorOrNil()
}

func (e *Spec) validateColorGroup(path string, r *ColorGroup, errs *errors.ErrorList) {
	if r.ID == 0 {
		errs.Append(errors.ErrMissingID)
	}
	if len(r.Colors) == 0 {
		errs.Append(errors.ErrEmptyResourceProps)
	}
	for j, c := range r.Colors {
		if c == (color.RGBA{}) {
			errs.Append(errors.NewIndexed(c, j, &errors.MissingFieldError{Name: attrColor}))
		}
	}
}

func (e *Spec) validateTexture2DGroup(m *go3mf.Model, path string, r *Texture2DGroup, errs *errors.ErrorList) {
	if r.ID == 0 {
		errs.Append(errors.ErrMissingID)
	}
	if r.TextureID == 0 {
		errs.Append(&errors.MissingFieldError{Name: attrTexID})
	} else if text, ok := m.FindAsset(path, r.TextureID); ok {
		if _, ok := text.(*Texture2D); !ok {
			errs.Append(errors.ErrTextureReference)
		}
	} else {
		errs.Append(errors.ErrTextureReference)
	}
	if len(r.Coords) == 0 {
		errs.Append(errors.ErrEmptyResourceProps)
	}
}

func (e *Spec) validateTexture2D(m *go3mf.Model, path string, r *Texture2D, errs *errors.ErrorList) {
	if r.ID == 0 {
		errs.Append(errors.ErrMissingID)
	}
	if r.Path == "" {
		errs.Append(&errors.MissingFieldError{Name: attrPath})
	} else {
		var hasTexture bool
		for _, a := range m.Attachments {
			if strings.EqualFold(a.Path, r.Path) {
				hasTexture = true
				break
			}
		}
		if !hasTexture {
			errs.Append(errors.ErrMissingTexturePart)
		}
	}
	if r.ContentType == 0 {
		errs.Append(&errors.MissingFieldError{Name: attrContentType})
	}
}

func (e *Spec) validateMultiProps(m *go3mf.Model, path string, r *MultiProperties, errs *errors.ErrorList) {
	if r.ID == 0 {
		errs.Append(errors.ErrMissingID)
	}
	if len(r.PIDs) == 0 {
		errs.Append(&errors.MissingFieldError{Name: attrPIDs})
	}
	if len(r.BlendMethods) > len(r.PIDs)-1 {
		errs.Append(errors.ErrMultiBlend)
	}
	if len(r.Multis) == 0 {
		errs.Append(errors.ErrEmptyResourceProps)
	}
	var (
		colorCount        int
		resourceUndefined bool
		lengths           = make([]int, len(r.PIDs))
	)
	for j, pid := range r.PIDs {
		if pr, ok := m.FindAsset(path, pid); ok {
			switch pr := pr.(type) {
			case *go3mf.BaseMaterials:
				if j != 0 {
					errs.Append(errors.ErrMaterialMulti)
				}
				lengths[j] = len(pr.Materials)
			case *CompositeMaterials:
				if j != 0 {
					errs.Append(errors.ErrMaterialMulti)
				}
				lengths[j] = len(pr.Composites)
			case *MultiProperties:
				errs.Append(errors.ErrMultiRefMulti)
			case *ColorGroup:
				if colorCount == 1 {
					errs.Append(errors.ErrMultiColors)
				}
				colorCount++
				lengths[j] = len(pr.Colors)
			}
		} else if !resourceUndefined {
			resourceUndefined = true
			errs.Append(errors.ErrMissingResource)
		}
	}
	for j, m := range r.Multis {
		for k, index := range m.PIndices {
			if k < len(r.PIDs) && lengths[k] < int(index) {
				errs.Append(errors.NewIndexed(m, j, errors.ErrIndexOutOfBounds))
				break
			}
		}
	}
}

func (e *Spec) validateCompositeMat(m *go3mf.Model, path string, r *CompositeMaterials, errs *errors.ErrorList) {
	if r.ID == 0 {
		errs.Append(errors.ErrMissingID)
	}
	if r.MaterialID == 0 {
		errs.Append(&errors.MissingFieldError{Name: attrMatID})
	} else if mat, ok := m.FindAsset(path, r.MaterialID); ok {
		if bm, ok := mat.(*go3mf.BaseMaterials); ok {
			for _, index := range r.Indices {
				if int(index) > len(bm.Materials) {
					errs.Append(errors.ErrIndexOutOfBounds)
					break
				}
			}
		} else {
			errs.Append(errors.ErrCompositeBase)
		}
	} else {
		errs.Append(errors.ErrMissingResource)
	}
	if len(r.Indices) == 0 {
		errs.Append(&errors.MissingFieldError{Name: attrMatIndices})
	}
	if len(r.Composites) == 0 {
		errs.Append(errors.ErrEmptyResourceProps)
	}
}
