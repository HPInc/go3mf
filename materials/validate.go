package materials

import (
	"image/color"
	"strings"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func (e *Spec) ValidateAsset(m *go3mf.Model, path string, r go3mf.Asset) (errs error) {
	switch r := r.(type) {
	case *ColorGroup:
		errs = e.validateColorGroup(path, r)
	case *Texture2DGroup:
		errs = e.validateTexture2DGroup(m, path, r)
	case *Texture2D:
		errs = e.validateTexture2D(m, path, r)
	case *MultiProperties:
		errs = e.validateMultiProps(m, path, r)
	case *CompositeMaterials:
		errs = e.validateCompositeMat(m, path, r)
	}
	return
}

func (e *Spec) validateColorGroup(path string, r *ColorGroup) (errs error) {
	if r.ID == 0 {
		errs = errors.Append(errs, errors.ErrMissingID)
	}
	if len(r.Colors) == 0 {
		errs = errors.Append(errs, errors.ErrEmptyResourceProps)
	}
	for j, c := range r.Colors {
		if c == (color.RGBA{}) {
			errs = errors.Append(errs, errors.WrapIndex(errors.NewMissingFieldError(attrColor), c, j))
		}
	}
	return
}

func (e *Spec) validateTexture2DGroup(m *go3mf.Model, path string, r *Texture2DGroup) (errs error) {
	if r.ID == 0 {
		errs = errors.Append(errs, errors.ErrMissingID)
	}
	if r.TextureID == 0 {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrTexID))
	} else if text, ok := m.FindAsset(path, r.TextureID); ok {
		if _, ok := text.(*Texture2D); !ok {
			errs = errors.Append(errs, errors.ErrTextureReference)
		}
	} else {
		errs = errors.Append(errs, errors.ErrTextureReference)
	}
	if len(r.Coords) == 0 {
		errs = errors.Append(errs, errors.ErrEmptyResourceProps)
	}
	return
}

func (e *Spec) validateTexture2D(m *go3mf.Model, path string, r *Texture2D) (errs error) {
	if r.ID == 0 {
		errs = errors.Append(errs, errors.ErrMissingID)
	}
	if r.Path == "" {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrPath))
	} else {
		var hasTexture bool
		for _, a := range m.Attachments {
			if strings.EqualFold(a.Path, r.Path) {
				hasTexture = true
				break
			}
		}
		if !hasTexture {
			errs = errors.Append(errs, errors.ErrMissingTexturePart)
		}
	}
	if r.ContentType == 0 {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrContentType))
	}
	return
}

func (e *Spec) validateMultiProps(m *go3mf.Model, path string, r *MultiProperties) (errs error) {
	if r.ID == 0 {
		errs = errors.Append(errs, errors.ErrMissingID)
	}
	if len(r.PIDs) == 0 {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrPIDs))
	}
	if len(r.BlendMethods) > len(r.PIDs)-1 {
		errs = errors.Append(errs, errors.ErrMultiBlend)
	}
	if len(r.Multis) == 0 {
		errs = errors.Append(errs, errors.ErrEmptyResourceProps)
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
					errs = errors.Append(errs, errors.ErrMaterialMulti)
				}
				lengths[j] = len(pr.Materials)
			case *CompositeMaterials:
				if j != 0 {
					errs = errors.Append(errs, errors.ErrMaterialMulti)
				}
				lengths[j] = len(pr.Composites)
			case *MultiProperties:
				errs = errors.Append(errs, errors.ErrMultiRefMulti)
			case *ColorGroup:
				if colorCount == 1 {
					errs = errors.Append(errs, errors.ErrMultiColors)
				}
				colorCount++
				lengths[j] = len(pr.Colors)
			}
		} else if !resourceUndefined {
			resourceUndefined = true
			errs = errors.Append(errs, errors.ErrMissingResource)
		}
	}
	for j, m := range r.Multis {
		for k, index := range m.PIndices {
			if k < len(r.PIDs) && lengths[k] < int(index) {
				errs = errors.Append(errs, errors.WrapIndex(errors.ErrIndexOutOfBounds, m, j))
				break
			}
		}
	}
	return
}

func (e *Spec) validateCompositeMat(m *go3mf.Model, path string, r *CompositeMaterials) (errs error) {
	if r.ID == 0 {
		errs = errors.Append(errs, errors.ErrMissingID)
	}
	if r.MaterialID == 0 {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrMatID))
	} else if mat, ok := m.FindAsset(path, r.MaterialID); ok {
		if bm, ok := mat.(*go3mf.BaseMaterials); ok {
			for _, index := range r.Indices {
				if int(index) > len(bm.Materials) {
					errs = errors.Append(errs, errors.ErrIndexOutOfBounds)
					break
				}
			}
		} else {
			errs = errors.Append(errs, errors.ErrCompositeBase)
		}
	} else {
		errs = errors.Append(errs, errors.ErrMissingResource)
	}
	if len(r.Indices) == 0 {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrMatIndices))
	}
	if len(r.Composites) == 0 {
		errs = errors.Append(errs, errors.ErrEmptyResourceProps)
	}
	return
}
