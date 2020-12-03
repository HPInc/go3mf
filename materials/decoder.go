package materials

import (
	"strconv"
	"strings"

	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
	"github.com/qmuntal/go3mf/spec/encoding"
)

func (e Spec) NewElementDecoder(el interface{}, nodeName string) (child encoding.ElementDecoder) {
	switch nodeName {
	case attrColorGroup:
		child = &colorGroupDecoder{resources: el.(*go3mf.Resources)}
	case attrTexture2DGroup:
		child = &tex2DGroupDecoder{resources: el.(*go3mf.Resources)}
	case attrTexture2D:
		child = &texture2DDecoder{resources: el.(*go3mf.Resources)}
	case attrCompositematerials:
		child = &compositeMaterialsDecoder{resources: el.(*go3mf.Resources)}
	case attrMultiProps:
		child = &multiPropertiesDecoder{resources: el.(*go3mf.Resources)}
	}
	return
}

func (e Spec) DecodeAttribute(_ interface{}, _ encoding.Attr) error { return nil }

type colorGroupDecoder struct {
	baseDecoder
	resources    *go3mf.Resources
	resource     ColorGroup
	colorDecoder colorDecoder
}

func (d *colorGroupDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *colorGroupDecoder) Child(name encoding.Name) (child encoding.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrColor {
		child = &d.colorDecoder
	}
	return
}

func (d *colorGroupDecoder) Start(attrs []encoding.Attr) (err error) {
	d.colorDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			id, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
			break
		}
	}
	return
}

type colorDecoder struct {
	baseDecoder
	resource *ColorGroup
}

func (d *colorDecoder) Start(attrs []encoding.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrColor {
			c, err1 := encoding.ParseRGBA(string(a.Value))
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.Colors = append(d.resource.Colors, c)
		}
	}
	return
}

type tex2DCoordDecoder struct {
	baseDecoder
	resource *Texture2DGroup
}

func (d *tex2DCoordDecoder) Start(attrs []encoding.Attr) (err error) {
	var u, v float32
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		val, err1 := strconv.ParseFloat(string(a.Value), 32)
		if err1 != nil {
			err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
		}
		switch a.Name.Local {
		case attrU:
			u = float32(val)
		case attrV:
			v = float32(val)
		}
	}
	d.resource.Coords = append(d.resource.Coords, TextureCoord{u, v})
	return
}

type tex2DGroupDecoder struct {
	baseDecoder
	resources         *go3mf.Resources
	resource          Texture2DGroup
	tex2DCoordDecoder tex2DCoordDecoder
}

func (d *tex2DGroupDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *tex2DGroupDecoder) Child(name encoding.Name) (child encoding.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrTex2DCoord {
		child = &d.tex2DCoordDecoder
	}
	return
}

func (d *tex2DGroupDecoder) Start(attrs []encoding.Attr) (err error) {
	d.tex2DCoordDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
		case attrTexID:
			val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.TextureID = uint32(val)
		}
	}
	return
}

type texture2DDecoder struct {
	baseDecoder
	resources *go3mf.Resources
	resource  Texture2D
}

func (d *texture2DDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *texture2DDecoder) Start(attrs []encoding.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
		case attrPath:
			d.resource.Path = string(a.Value)
		case attrContentType:
			d.resource.ContentType, _ = newTexture2DType(string(a.Value))
		case attrTileStyleU:
			d.resource.TileStyleU, _ = newTileStyle(string(a.Value))
		case attrTileStyleV:
			d.resource.TileStyleV, _ = newTileStyle(string(a.Value))
		case attrFilter:
			d.resource.Filter, _ = newTextureFilter(string(a.Value))
		}
	}
	return
}

type compositeMaterialsDecoder struct {
	baseDecoder
	resources        *go3mf.Resources
	resource         CompositeMaterials
	compositeDecoder compositeDecoder
}

func (d *compositeMaterialsDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *compositeMaterialsDecoder) Child(name encoding.Name) (child encoding.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrComposite {
		child = &d.compositeDecoder
	}
	return
}

func (d *compositeMaterialsDecoder) Start(attrs []encoding.Attr) (err error) {
	d.compositeDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
		case attrMatID:
			val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.MaterialID = uint32(val)
		case attrMatIndices:
			for _, f := range strings.Fields(string(a.Value)) {
				val, err1 := strconv.ParseUint(f, 10, 32)
				if err1 != nil {
					err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
				}
				d.resource.Indices = append(d.resource.Indices, uint32(val))
			}
		}
	}
	return
}

type compositeDecoder struct {
	baseDecoder
	resource *CompositeMaterials
}

func (d *compositeDecoder) Start(attrs []encoding.Attr) (err error) {
	var composite Composite
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrValues {
			for _, f := range strings.Fields(string(a.Value)) {
				val, err1 := strconv.ParseFloat(f, 32)
				if err1 != nil {
					err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
				}
				composite.Values = append(composite.Values, float32(val))
			}
		}
	}
	d.resource.Composites = append(d.resource.Composites, composite)
	return
}

type multiPropertiesDecoder struct {
	baseDecoder
	resources    *go3mf.Resources
	resource     MultiProperties
	multiDecoder multiDecoder
}

func (d *multiPropertiesDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *multiPropertiesDecoder) Child(name encoding.Name) (child encoding.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrMulti {
		child = &d.multiDecoder
	}
	return
}

func (d *multiPropertiesDecoder) Start(attrs []encoding.Attr) (err error) {
	d.multiDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
		case attrBlendMethods:
			for _, f := range strings.Fields(string(a.Value)) {
				val, _ := newBlendMethod(f)
				d.resource.BlendMethods = append(d.resource.BlendMethods, val)
			}
		case attrPIDs:
			for _, f := range strings.Fields(string(a.Value)) {
				val, err1 := strconv.ParseUint(f, 10, 32)
				if err1 != nil {
					err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
				}
				d.resource.PIDs = append(d.resource.PIDs, uint32(val))
			}
		}
	}
	return
}

type multiDecoder struct {
	baseDecoder
	resource *MultiProperties
}

func (d *multiDecoder) Start(attrs []encoding.Attr) (err error) {
	var multi Multi
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrPIndices {
			for _, f := range strings.Fields(string(a.Value)) {
				val, err1 := strconv.ParseUint(f, 10, 32)
				if err1 != nil {
					err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
				}
				multi.PIndices = append(multi.PIndices, uint32(val))
			}
		}
	}
	d.resource.Multis = append(d.resource.Multis, multi)
	return
}

type baseDecoder struct {
}

func (d *baseDecoder) End() {}
