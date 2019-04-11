package io3mf

import (
	"encoding/xml"

	go3mf "github.com/qmuntal/go3mf"
)

type colorGroupDecoder struct {
	emptyDecoder
	resource     go3mf.ColorGroupResource
	colorDecoder colorDecoder
}

func (d *colorGroupDecoder) Open() {
	d.resource.ModelPath = d.file.path
	d.colorDecoder.resource = &d.resource
}

func (d *colorGroupDecoder) Close() bool {
	d.file.parser.CloseResource()
	d.file.AddResource(&d.resource)
	return true
}

func (d *colorGroupDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsMaterialSpec && name.Local == attrColor {
		child = &d.colorDecoder
	}
	return
}

func (d *colorGroupDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			d.resource.ID, ok = d.file.parser.ParseResourceID(a.Value)
			break
		}
	}
	return ok
}

type colorDecoder struct {
	emptyDecoder
	resource *go3mf.ColorGroupResource
}

func (d *colorDecoder) Attributes(attrs []xml.Attr) bool {
	for _, a := range attrs {
		if a.Name.Local == attrColor {
			c, err := strToSRGB(a.Value)
			if err != nil {
				return d.file.parser.InvalidRequiredAttr(attrColor, a.Value)
			}
			d.resource.Colors = append(d.resource.Colors, c)
		}
	}
	return true
}

type tex2DCoordDecoder struct {
	emptyDecoder
	resource *go3mf.Texture2DGroupResource
}

func (d *tex2DCoordDecoder) Attributes(attrs []xml.Attr) bool {
	var u, v float32
	ok := true
	for _, a := range attrs {
		switch a.Name.Local {
		case attrU:
			u, ok = d.file.parser.ParseFloat32Required(a.Name.Local, a.Value)
		case attrV:
			v, ok = d.file.parser.ParseFloat32Required(a.Name.Local, a.Value)
		}
		if !ok {
			break
		}
	}
	d.resource.Coords = append(d.resource.Coords, go3mf.TextureCoord{float32(u), float32(v)})
	return ok
}

type tex2DGroupDecoder struct {
	emptyDecoder
	resource          go3mf.Texture2DGroupResource
	tex2DCoordDecoder tex2DCoordDecoder
}

func (d *tex2DGroupDecoder) Open() {
	d.resource.ModelPath = d.file.path
	d.tex2DCoordDecoder.resource = &d.resource
}

func (d *tex2DGroupDecoder) Close() bool {
	if !d.file.parser.CloseResource() {
		return false
	}
	d.file.AddResource(&d.resource)
	return true
}

func (d *tex2DGroupDecoder) Child(name xml.Name) (child nodeDecoder) {
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
			d.resource.ID, ok = d.file.parser.ParseResourceID(a.Value)
		case attrTexID:
			d.resource.TextureID, ok = d.file.parser.ParseUint32Required(a.Name.Local, a.Value)
		}
		if !ok {
			break
		}
	}
	return ok
}

type texture2DDecoder struct {
	emptyDecoder
	resource go3mf.Texture2DResource
}

func (d *texture2DDecoder) Open() {
	d.resource.ModelPath = d.file.path
}

func (d *texture2DDecoder) Close() bool {
	if !d.file.parser.CloseResource() {
		return false
	}
	d.file.AddResource(&d.resource)
	return true
}

func (d *texture2DDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			d.resource.ID, ok = d.file.parser.ParseResourceID(a.Value)
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
		return d.file.parser.MissingAttr(attrPath)
	}
	return ok
}
