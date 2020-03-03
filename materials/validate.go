package materials

import (
	"errors"
	"image/color"
	"sort"
	"strings"

	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

var (
	ErrTextureReference       = errors.New("MUST reference to a texture resource")
	ErrCompositeBase          = errors.New("MUST reference to a basematerials group")
	ErrMaterialMulti          = errors.New("material, if included, MUST be positioned in the first layer")
	ErrMultiRefMulti          = errors.New("MUST NOT contain any references to a multiproperties")
	ErrMultiRefMultipleColors = errors.New("MUST NOT contain more than one reference to a colorgroup")
	ErrMissingMultiBlend      = errors.New("MUST NOT have more blendmethods than layers – 1")
	ErrMissingTexturePart     = errors.New("texture part MUST be added as an attachment")
)

// Validate checks that the model is conformant with the 3MF spec.
// Core spec related checks are not reported.
func Validate(model *go3mf.Model) []error {
	var hasExt bool
	for _, ext := range model.Namespaces {
		if ext.Space == ExtensionName {
			hasExt = true
			break
		}
	}
	if !hasExt {
		return nil
	}
	var err []error
	err = validate(model, model.PathOrDefault(), &model.Resources, err)

	s := make([]string, 0, len(model.Childs))
	for path := range model.Childs {
		s = append(s, path)
	}
	sort.Strings(s)
	for _, path := range s {
		c := model.Childs[path]
		err = validate(model, path, &c.Resources, err)
	}
	return err
}

func validate(model *go3mf.Model, path string, res *go3mf.Resources, err []error) []error {
	v := validator{m: model, path: path, res: res}
	v.Validate()
	return append(err, v.warnings...)
}

type validator struct {
	m        *go3mf.Model
	warnings []error
	path     string
	res      *go3mf.Resources
	ids      map[uint32]interface{}
	c        chan error
}

func (v *validator) AddWarning(err ...error) {
	v.warnings = append(v.warnings, err...)
}

func (v *validator) Validate() {
	v.ids = make(map[uint32]interface{})
	for i, r := range v.res.Assets {
		v.ids[r.Identify()] = r
		switch r := r.(type) {
		case *ColorGroupResource:
			v.validateColorGroup(i, r)
		case *Texture2DGroupResource:
			v.validateTextureGroup(i, r)
		case *CompositeMaterialsResource:
			v.validateComposite(i, r)
		case *MultiPropertiesResource:
			v.validateMulti(i, r)
		case *Texture2DResource:
			v.validateTexture(i, r)
		}
	}
}

func (v *validator) validateColorGroup(i int, r *ColorGroupResource) {
	if len(r.Colors) == 0 {
		v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrEmptyResourceProps))
	}
	var emptyColor color.RGBA
	for j, c := range r.Colors {
		if c == emptyColor {
			v.AddWarning(specerr.NewAsset(v.path, i, r, &specerr.ResourcePropertyError{
				Index: j,
				Err:   &specerr.MissingFieldError{Name: attrColor},
			}))
		}
	}
}

func (v *validator) validateTextureGroup(i int, r *Texture2DGroupResource) {
	if r.TextureID == 0 {
		v.AddWarning(specerr.NewAsset(v.path, i, r, &specerr.MissingFieldError{Name: attrTexID}))
	} else if text, ok := v.ids[r.TextureID]; ok {
		if _, ok := text.(*Texture2DResource); !ok {
			v.AddWarning(specerr.NewAsset(v.path, i, r, ErrTextureReference))
		}
	} else {
		v.AddWarning(specerr.NewAsset(v.path, i, r, ErrTextureReference))
	}
	if len(r.Coords) == 0 {
		v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrEmptyResourceProps))
	}
}

func (v *validator) validateTexture(i int, r *Texture2DResource) {
	if r.Path == "" {
		v.AddWarning(specerr.NewAsset(v.path, i, r, &specerr.MissingFieldError{Name: attrPath}))
	} else {
		var hasTexture bool
		for _, a := range v.m.Attachments {
			if strings.EqualFold(a.Path, r.Path) {
				hasTexture = true
				break
			}
		}
		if !hasTexture {
			v.AddWarning(specerr.NewAsset(v.path, i, r, ErrMissingTexturePart))
		}
	}
	if r.ContentType == 0 {
		v.AddWarning(specerr.NewAsset(v.path, i, r, &specerr.MissingFieldError{Name: attrContentType}))
	}
}

func (v *validator) validateMulti(i int, r *MultiPropertiesResource) {
	if len(r.PIDs) == 0 {
		v.AddWarning(specerr.NewAsset(v.path, i, r, &specerr.MissingFieldError{Name: attrPIDs}))
	}
	if len(r.BlendMethods) > len(r.PIDs)-1 {
		v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrMultiBlend))
	}
	if len(r.Multis) == 0 {
		v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrEmptyResourceProps))
	}
	var (
		colorCount        int
		resourceUndefined bool
		lengths           = make([]int, len(r.PIDs))
	)
	for j, pid := range r.PIDs {
		if pr, ok := v.ids[pid]; ok {
			switch pr := pr.(type) {
			case *go3mf.BaseMaterialsResource:
				if j != 0 {
					v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrMaterialMulti))
				}
				lengths[j] = len(pr.Materials)
			case *CompositeMaterialsResource:
				if j != 0 {
					v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrMaterialMulti))
				}
				lengths[j] = len(pr.Composites)
			case *MultiPropertiesResource:
				v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrMultiRefMulti))
			case *ColorGroupResource:
				if colorCount == 1 {
					v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrMultiColors))
				}
				colorCount++
				lengths[j] = len(pr.Colors)
			}
		} else if !resourceUndefined {
			resourceUndefined = true
			v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrMissingResource))
		}
	}
	for j, m := range r.Multis {
		for k, index := range m.PIndex {
			if k < len(r.PIDs) && lengths[k] < int(index) {
				v.AddWarning(specerr.NewAsset(v.path, i, r, &specerr.ResourcePropertyError{
					Index: j,
					Err:   specerr.ErrIndexOutOfBounds,
				}))
				break
			}
		}
	}
}
func (v *validator) validateComposite(i int, r *CompositeMaterialsResource) {
	if r.MaterialID == 0 {
		v.AddWarning(specerr.NewAsset(v.path, i, r, &specerr.MissingFieldError{Name: attrMatID}))
	} else if mat, ok := v.ids[r.MaterialID]; ok {
		if bm, ok := mat.(*go3mf.BaseMaterialsResource); ok {
			for _, index := range r.Indices {
				if int(index) > len(bm.Materials) {
					v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrIndexOutOfBounds))
					break
				}
			}
		} else {
			v.AddWarning(specerr.NewAsset(v.path, i, r, ErrCompositeBase))
		}
	} else {
		v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrMissingResource))
	}
	if len(r.Indices) == 0 {
		v.AddWarning(specerr.NewAsset(v.path, i, r, &specerr.MissingFieldError{Name: attrMatIndices}))
	}
	if len(r.Composites) == 0 {
		v.AddWarning(specerr.NewAsset(v.path, i, r, specerr.ErrEmptyResourceProps))
	}
}
