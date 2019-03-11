package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	go3mf "github.com/qmuntal/go3mf"
)

type colorGroupDecoder struct {
	x            *xml.Decoder
	r            *Reader
	colorMapping *colorMapping
	id           uint64
	colorIndex   uint64
}

func (d *colorGroupDecoder) Decode(se xml.StartElement) error {
	if err := d.parseAttr(se.Attr); err != nil {
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
		case xml.EndElement:
			if tp.Name.Space == nsMaterialSpec && tp.Name.Local == attrColorGroup {
				return nil
			}
		}
	}
}

func (d *colorGroupDecoder) parseAttr(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
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
	if err := d.parseAttr(se.Attr); err != nil {
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
		case xml.EndElement:
			if tp.Name.Space == nsMaterialSpec && tp.Name.Local == attrTexture2DGroup {
				return nil
			}
		}
	}
}

func (d *tex2DGroupDecoder) parseAttr(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
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
		case attrV:
			v, err = strconv.ParseFloat(a.Value, 64)
		}
		if err != nil {
			return err
		}
	}
	d.texCoordMapping.register(d.id, d.texCoordIndex, d.textureID, float32(u), float32(v))
	d.texCoordIndex++
	return nil
}

type texture2DDecoder struct {
	x              *xml.Decoder
	r              *Reader
	model          *go3mf.Model
	id             uint64
	path           string
	contentType    go3mf.Texture2DType
	styleU, styleV go3mf.TileStyle
	filter         go3mf.TextureFilter
}

func (d *texture2DDecoder) Decode(se xml.StartElement) error {
	if err := d.parseAttr(se.Attr); err != nil {
		return err
	}
	if d.id == 0 {
		return errors.New("go3mf: missing texture2d id attribute")
	}
	texture2d := go3mf.NewTexture2DResource(d.id)
	texture2d.Path = d.path
	texture2d.ContentType = d.contentType
	if d.styleU != "" {
		texture2d.TileStyleU = d.styleU
	}
	if d.styleV != "" {
		texture2d.TileStyleV = d.styleV
	}
	if d.filter != "" {
		texture2d.Filter = d.filter
	}
	d.model.Resources = append(d.model.Resources, texture2d)
	return nil
}

func (d *texture2DDecoder) parseAttr(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		var err error
		ok := true
		switch a.Name.Local {
		case attrID:
			if d.id != 0 {
				err = errors.New("go3mf: duplicated texture2d id attribute")
			} else {
				d.id, err = strconv.ParseUint(a.Value, 10, 64)
			}
		case attrPath:
			d.path = a.Value
		case attrContentType:
			d.contentType, ok = go3mf.NewTexture2DType(a.Value)
		case attrTileStyleU:
			d.styleU, ok = go3mf.NewTileStyle(a.Value)
		case attrTileStyleV:
			d.styleV, ok = go3mf.NewTileStyle(a.Value)
		case attrFilter:
			d.filter, ok = go3mf.NewTextureFilter(a.Value)
		}
		if err != nil || !ok {
			return errors.New("go3mf: texture2d attribute not valid")
		}
	}
	return nil
}
