package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	mdl "github.com/qmuntal/go3mf/internal/model"
)

type colorGroupDecoder struct {
	x            *xml.Decoder
	r            *Reader
	colorMapping *colorMapping
	id           uint64
	colorIndex   uint64
}

func (d *colorGroupDecoder) Decode(se xml.StartElement) error {
	if err := d.parseAttr(se); err != nil {
		return err
	}
	if d.id == 0 {
		return errors.New("go3mf: missing color group id attribute")
	}
	for {
		t, err := d.x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsMaterialSpec && tp.Name.Local == attrColor {
				if err := d.addColor(tp.Attr); err != nil {
					return err
				}
			}
		}
	}
}

func (d *colorGroupDecoder) parseAttr(se xml.StartElement) error {
	for _, a := range se.Attr {
		if a.Name.Space == "" && se.Name.Local == attrID {
			if d.id != 0 {
				return errors.New("go3mf: duplicated color group id attribute")
			}
			var err error
			if d.id, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
				return errors.New("go3mf: color group id is not valid")
			}
		}
	}
	return nil
}

func (d *colorGroupDecoder) addColor(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Local == attrColor {
			c, err := strToSRGB(a.Value)
			if err != nil {
				return err
			}
			d.colorMapping.register(d.id, d.colorIndex, c)
			d.colorIndex++
		}
	}
	return nil
}

type tex2DGroupDecoder struct {
	x               *xml.Decoder
	r               *Reader
	texCoordMapping *texCoordMapping
	id              uint64
	textureID       uint64
	texCoordIndex   uint64
}

func (d *tex2DGroupDecoder) Decode(se xml.StartElement) error {
	if err := d.parseAttr(se); err != nil {
		return err
	}
	if d.id == 0 {
		return errors.New("go3mf: missing tex2Coord group id attribute")
	}
	for {
		t, err := d.x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsMaterialSpec && tp.Name.Local == attrTex2DCoord {
				if err := d.addTextureCoord(tp.Attr); err != nil {
					return err
				}
			}
		}
	}
}

func (d *tex2DGroupDecoder) parseAttr(se xml.StartElement) error {
	for _, a := range se.Attr {
		if a.Name.Space == "" {
			continue
		}
		switch se.Name.Local {
		case attrID:
			if d.id != 0 {
				return errors.New("go3mf: duplicated tex2Coord group id attribute")
			}
			var err error
			if d.id, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
				return errors.New("go3mf: tex2Coord group id is not valid")
			}
		case attrTexID:
			if d.textureID != 0 {
				return errors.New("go3mf: duplicated tex2Coord group texid attribute")
			}
			var err error
			if d.textureID, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
				return errors.New("go3mf: tex2Coord group texid is not valid")
			}
		}
	}
	return nil
}

func (d *tex2DGroupDecoder) addTextureCoord(attrs []xml.Attr) error {
	var u, v float64
	var err error
	for _, a := range attrs {
		switch a.Name.Local {
		case attrU:
			u, err = strconv.ParseFloat(a.Value, 64)
			if err != nil {
				return err
			}
		case attrV:
			v, err = strconv.ParseFloat(a.Value, 64)
			if err != nil {
				return err
			}
		}
	}
	d.texCoordMapping.register(d.id, d.texCoordIndex, d.textureID, float32(u), float32(v))
	d.texCoordIndex++
	return nil
}

type texture2DDecoder struct {
	x              *xml.Decoder
	r              *Reader
	model          *mdl.Model
	id             uint64
	path           string
	contentType    mdl.Texture2DType
	styleU, styleV mdl.TileStyle
	filter         mdl.TextureFilter
}

func (d *texture2DDecoder) Decode(se xml.StartElement) error {
	if err := d.parseAttr(se); err != nil {
		return err
	}
	if d.id == 0 {
		return errors.New("go3mf: missing texture2d id attribute")
	}
	texture2d, err := mdl.NewTexture2DResource(d.id, d.model)
	if err != nil {
		return err
	}
	texture2d.Path = d.path
	texture2d.ContentType = d.contentType
	texture2d.TileStyleU = d.styleU
	texture2d.TileStyleV = d.styleV
	texture2d.Filter = d.filter
	return d.model.AddResource(texture2d)
}

func (d *texture2DDecoder) parseAttr(se xml.StartElement) error {
	for _, a := range se.Attr {
		if a.Name.Space == "" {
			continue
		}
		var err error
		var ok bool
		switch se.Name.Local {
		case attrID:
			if d.id != 0 {
				err = errors.New("go3mf: duplicated texture2d id attribute")
			} else {
				d.id, err = strconv.ParseUint(a.Value, 10, 64)
			}
		case attrPath:
			d.path = a.Value
		case attrContentType:
			d.contentType, ok = mdl.NewTexture2DType(a.Value)
		case attrTileStyleU:
			d.styleU, ok = mdl.NewTileStyle(a.Value)
		case attrTileStyleV:
			d.styleV, ok = mdl.NewTileStyle(a.Value)
		case attrFilter:
			d.filter, ok = mdl.NewTextureFilter(a.Value)
		}
		if err != nil || !ok {
			return errors.New("go3mf: texture2d attribute not valid")
		}
	}
	return nil
}
