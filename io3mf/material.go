package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	go3mf "github.com/qmuntal/go3mf"
)

type colorGroupDecoder struct {
	emptyDecoder
	resource     go3mf.ColorGroupResource
	colorDecoder colorDecoder
}

func (d *colorGroupDecoder) Open() error {
	d.resource.ModelPath = d.ModelFile().Path()
	d.colorDecoder.resource = &d.resource
	return nil
}

func (d *colorGroupDecoder) Close() error {
	if d.resource.ID == 0 {
		return errors.New("go3mf: missing color group id attribute")
	}
	d.ModelFile().AddResource(&d.resource)
	return nil
}

func (d *colorGroupDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsMaterialSpec && name.Local == attrColor {
		child = &d.colorDecoder
	}
	return
}

func (d *colorGroupDecoder) Attributes(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			if d.resource.ID != 0 {
				return errors.New("go3mf: duplicated color group id attribute")
			}
			var id uint64
			if id, err = strconv.ParseUint(a.Value, 10, 32); err != nil {
				return errors.New("go3mf: color group id is not valid")
			}
			d.resource.ID = uint32(id)
		}
	}
	return nil
}

type colorDecoder struct {
	emptyDecoder
	resource *go3mf.ColorGroupResource
}

func (d *colorDecoder) Attributes(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Local == attrColor {
			c, err := strToSRGB(a.Value)
			if err != nil {
				return err
			}
			d.resource.Colors = append(d.resource.Colors, c)
		}
	}
	return nil
}

type tex2DCoordDecoder struct {
	emptyDecoder
	resource *go3mf.Texture2DGroupResource
}

func (d *tex2DCoordDecoder) Attributes(attrs []xml.Attr) error {
	var u, v float64
	var err error
	for _, a := range attrs {
		switch a.Name.Local {
		case attrU:
			u, err = strconv.ParseFloat(a.Value, 64)
		case attrV:
			v, err = strconv.ParseFloat(a.Value, 64)
		}
		if err != nil {
			return err
		}
	}
	d.resource.Coords = append(d.resource.Coords, go3mf.TextureCoord{float32(u), float32(v)})
	return nil
}

type tex2DGroupDecoder struct {
	emptyDecoder
	resource          go3mf.Texture2DGroupResource
	tex2DCoordDecoder tex2DCoordDecoder
}

func (d *tex2DGroupDecoder) Open() error {
	d.resource.ModelPath = d.ModelFile().Path()
	d.tex2DCoordDecoder.resource = &d.resource
	return nil
}

func (d *tex2DGroupDecoder) Close() error {
	if d.resource.ID == 0 {
		return errors.New("go3mf: missing color group id attribute")
	}
	d.ModelFile().AddResource(&d.resource)
	return nil
}

func (d *tex2DGroupDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsMaterialSpec && name.Local == attrTex2DCoord {
		child = &d.tex2DCoordDecoder
	}
	return
}

func (d *tex2DGroupDecoder) Attributes(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			if d.resource.ID != 0 {
				return errors.New("go3mf: duplicated tex2Coord group id attribute")
			}
			var id uint64
			if id, err = strconv.ParseUint(a.Value, 10, 32); err != nil {
				return errors.New("go3mf: tex2Coord group id is not valid")
			}
			d.resource.ID = uint32(id)
		case attrTexID:
			if d.resource.TextureID != 0 {
				return errors.New("go3mf: duplicated tex2Coord group texid attribute")
			}
			var id uint64
			if id, err = strconv.ParseUint(a.Value, 10, 32); err != nil {
				return errors.New("go3mf: tex2Coord group texid is not valid")
			}
			d.resource.TextureID = uint32(id)
		}
	}
	return nil
}

type texture2DDecoder struct {
	emptyDecoder
	resource go3mf.Texture2DResource
}

func (d *texture2DDecoder) Open() error {
	d.resource.ModelPath = d.ModelFile().Path()
	return nil
}

func (d *texture2DDecoder) Close() error {
	if d.resource.ID == 0 {
		return errors.New("go3mf: missing texture2d id attribute")
	}
	d.ModelFile().AddResource(&d.resource)
	return nil
}

func (d *texture2DDecoder) Attributes(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		var err error
		ok := true
		switch a.Name.Local {
		case attrID:
			if d.resource.ID != 0 {
				err = errors.New("go3mf: duplicated texture2d id attribute")
			} else {
				var id uint64
				id, err = strconv.ParseUint(a.Value, 10, 32)
				d.resource.ID = uint32(id)
			}
		case attrPath:
			d.resource.Path = a.Value
		case attrContentType:
			d.resource.ContentType, ok = newTexture2DType(a.Value)
		case attrTileStyleU:
			d.resource.TileStyleU, ok = newTileStyle(a.Value)
		case attrTileStyleV:
			d.resource.TileStyleV, ok = newTileStyle(a.Value)
		case attrFilter:
			d.resource.Filter, ok = newTextureFilter(a.Value)
		}
		if err != nil || !ok {
			return errors.New("go3mf: texture2d attribute not valid")
		}
	}
	return nil
}
