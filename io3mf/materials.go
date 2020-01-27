package io3mf

import (
	"encoding/xml"
	"strings"

	go3mf "github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/iohelper"
)

type colorGroupDecoder struct {
	iohelper.EmptyDecoder
	resource     go3mf.ColorGroupResource
	colorDecoder colorDecoder
}

func (d *colorGroupDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.colorDecoder.resource = &d.resource
}

func (d *colorGroupDecoder) Close() bool {
	d.Scanner.AddResource(&d.resource)
	return d.Scanner.CloseResource()
}

func (d *colorGroupDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsMaterialSpec && name.Local == attrColor {
		child = &d.colorDecoder
	}
	return
}

func (d *colorGroupDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			d.resource.ID, ok = d.Scanner.ParseResourceID(a.Value)
			break
		}
	}
	return ok
}

type colorDecoder struct {
	iohelper.EmptyDecoder
	resource *go3mf.ColorGroupResource
}

func (d *colorDecoder) Attributes(attrs []xml.Attr) bool {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrColor {
			c, err := iohelper.ReadRGB(a.Value)
			if err != nil {
				return d.Scanner.InvalidRequiredAttr(attrColor, a.Value)
			}
			d.resource.Colors = append(d.resource.Colors, c)
		}
	}
	return true
}

type tex2DCoordDecoder struct {
	iohelper.EmptyDecoder
	resource *go3mf.Texture2DGroupResource
}

func (d *tex2DCoordDecoder) Attributes(attrs []xml.Attr) bool {
	var u, v float32
	ok := true
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrU:
			u, ok = d.Scanner.ParseFloat32Required(attrU, a.Value)
		case attrV:
			v, ok = d.Scanner.ParseFloat32Required(attrV, a.Value)
		}
		if !ok {
			break
		}
	}
	d.resource.Coords = append(d.resource.Coords, go3mf.TextureCoord{float32(u), float32(v)})
	return ok
}

type tex2DGroupDecoder struct {
	iohelper.EmptyDecoder
	resource          go3mf.Texture2DGroupResource
	tex2DCoordDecoder tex2DCoordDecoder
}

func (d *tex2DGroupDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.tex2DCoordDecoder.resource = &d.resource
}

func (d *tex2DGroupDecoder) Close() bool {
	d.Scanner.AddResource(&d.resource)
	return d.Scanner.CloseResource()
}

func (d *tex2DGroupDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsMaterialSpec && name.Local == attrTex2DCoord {
		child = &d.tex2DCoordDecoder
	}
	return
}

func (d *tex2DGroupDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			d.resource.ID, ok = d.Scanner.ParseResourceID(a.Value)
		case attrTexID:
			d.resource.TextureID, ok = d.Scanner.ParseUint32Required(attrTexID, a.Value)
		}
		if !ok {
			break
		}
	}
	return ok
}

type texture2DDecoder struct {
	iohelper.EmptyDecoder
	resource go3mf.Texture2DResource
}

func (d *texture2DDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
}

func (d *texture2DDecoder) Close() bool {
	d.Scanner.AddResource(&d.resource)
	return d.Scanner.CloseResource()
}

func (d *texture2DDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			d.resource.ID, ok = d.Scanner.ParseResourceID(a.Value)
		case attrPath:
			d.resource.Path = a.Value
		case attrContentType:
			d.resource.ContentType, _ = newTexture2DType(a.Value)
		case attrTileStyleU:
			d.resource.TileStyleU, _ = newTileStyle(a.Value)
		case attrTileStyleV:
			d.resource.TileStyleV, _ = newTileStyle(a.Value)
		case attrFilter:
			d.resource.Filter, _ = newTextureFilter(a.Value)
		}
		if !ok {
			break
		}
	}
	if d.resource.Path == "" {
		return d.Scanner.MissingAttr(attrPath)
	}
	return ok
}

type compositeMaterialsDecoder struct {
	iohelper.EmptyDecoder
	resource         go3mf.CompositeMaterialsResource
	compositeDecoder compositeDecoder
}

func (d *compositeMaterialsDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.compositeDecoder.resource = &d.resource
}

func (d *compositeMaterialsDecoder) Close() bool {
	d.Scanner.AddResource(&d.resource)
	return d.Scanner.CloseResource()
}

func (d *compositeMaterialsDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsMaterialSpec && name.Local == attrComposite {
		child = &d.compositeDecoder
	}
	return
}

func (d *compositeMaterialsDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			d.resource.ID, ok = d.Scanner.ParseResourceID(a.Value)
		case attrMatID:
			d.resource.MaterialID, ok = d.Scanner.ParseUint32Required(attrMatID, a.Value)
		case attrMatIndices:
			for _, f := range strings.Fields(a.Value) {
				var val uint32
				if val, ok = d.Scanner.ParseUint32Required(attrValues, f); ok {
					d.resource.Indices = append(d.resource.Indices, val)
				} else {
					break
				}
			}
		}
		if !ok {
			break
		}
	}
	if d.resource.MaterialID == 0 {
		ok = d.Scanner.MissingAttr(attrMatID)
	}
	if ok && len(d.resource.Indices) == 0 {
		ok = d.Scanner.MissingAttr(attrMatIndices)
	}
	return ok
}

type compositeDecoder struct {
	iohelper.EmptyDecoder
	resource *go3mf.CompositeMaterialsResource
}

func (d *compositeDecoder) Attributes(attrs []xml.Attr) (ok bool) {
	composite := go3mf.Composite{}
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrValues {
			for _, f := range strings.Fields(a.Value) {
				var val float64
				if val, ok = d.Scanner.ParseFloat64Required(attrValues, f); ok {
					composite.Values = append(composite.Values, val)
				} else {
					break
				}
			}
		}
	}
	if len(composite.Values) == 0 {
		ok = d.Scanner.MissingAttr(attrValues)
	}
	if ok {
		d.resource.Composites = append(d.resource.Composites, composite)
	}
	return ok
}

type multiPropertiesDecoder struct {
	iohelper.EmptyDecoder
	resource     go3mf.MultiPropertiesResource
	multiDecoder multiDecoder
}

func (d *multiPropertiesDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.multiDecoder.resource = &d.resource
}

func (d *multiPropertiesDecoder) Close() bool {
	d.Scanner.AddResource(&d.resource)
	return d.Scanner.CloseResource()
}

func (d *multiPropertiesDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsMaterialSpec && name.Local == attrMulti {
		child = &d.multiDecoder
	}
	return
}

func (d *multiPropertiesDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			d.resource.ID, ok = d.Scanner.ParseResourceID(a.Value)
		case attrBlendMethods:
			for _, f := range strings.Fields(a.Value) {
				val, _ := newBlendMethod(f)
				d.resource.BlendMethods = append(d.resource.BlendMethods, val)
			}
		case attrPIDs:
			for _, f := range strings.Fields(a.Value) {
				var val uint32
				if val, ok = d.Scanner.ParseUint32Required(attrPIDs, f); ok {
					d.resource.Resources = append(d.resource.Resources, val)
				} else {
					break
				}
			}
		}
		if !ok {
			break
		}
	}
	if ok && len(d.resource.Resources) == 0 {
		ok = d.Scanner.MissingAttr(attrPIDs)
	}
	return ok
}

type multiDecoder struct {
	iohelper.EmptyDecoder
	resource *go3mf.MultiPropertiesResource
}

func (d *multiDecoder) Attributes(attrs []xml.Attr) (ok bool) {
	multi := go3mf.Multi{}
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrPIndices {
			for _, f := range strings.Fields(a.Value) {
				var val uint32
				if val, ok = d.Scanner.ParseUint32Required(attrPIndices, f); ok {
					multi.ResourceIndices = append(multi.ResourceIndices, val)
				} else {
					break
				}
			}
		}
	}
	if len(multi.ResourceIndices) == 0 {
		ok = d.Scanner.MissingAttr(attrPIndices)
	}
	if ok {
		d.resource.Multis = append(d.resource.Multis, multi)
	}
	return ok
}
