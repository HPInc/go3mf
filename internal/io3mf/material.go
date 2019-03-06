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
	model        *mdl.Model
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
	model           *mdl.Model
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
		return errors.New("go3mf: missing color group id attribute")
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
		if a.Name.Space == "" && se.Name.Local == attrID {
			switch se.Name.Local {
			case attrID:
				if d.id != 0 {
					return errors.New("go3mf: duplicated texture2d group id attribute")
				}
				var err error
				if d.id, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
					return errors.New("go3mf: texture2d group id is not valid")
				}
			case attrTexID:
				if d.textureID != 0 {
					return errors.New("go3mf: duplicated texture2d group texid attribute")
				}
				var err error
				if d.textureID, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
					return errors.New("go3mf: texture2d group texid is not valid")
				}
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
