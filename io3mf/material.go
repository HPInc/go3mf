package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	go3mf "github.com/qmuntal/go3mf"
)

type colorGroupDecoder struct {
	r            *Reader
	colorMapping *colorMapping
	id           uint64
	colorIndex   uint64
}

func (d *colorGroupDecoder) Decode(x xml.TokenReader, attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	if d.id == 0 {
		return errors.New("go3mf: missing color group id attribute")
	}
	for {
		t, err := x.Token()
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
	r               *Reader
	texCoordMapping *texCoordMapping
	id              uint64
	textureID       uint64
	texCoordIndex   uint64
}

func (d *tex2DGroupDecoder) Decode(x xml.TokenReader, attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	if d.id == 0 {
		return errors.New("go3mf: missing tex2Coord group id attribute")
	}
	for {
		t, err := x.Token()
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
	r       *Reader
	texture go3mf.Texture2DResource
}

func (d *texture2DDecoder) Decode(attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	if d.texture.ID == 0 {
		return errors.New("go3mf: missing texture2d id attribute")
	}
	d.r.addResource(&d.texture)
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
			if d.texture.ID != 0 {
				err = errors.New("go3mf: duplicated texture2d id attribute")
			} else {
				d.texture.ID, err = strconv.ParseUint(a.Value, 10, 64)
			}
		case attrPath:
			d.texture.Path = a.Value
		case attrContentType:
			d.texture.ContentType, ok = newTexture2DType(a.Value)
		case attrTileStyleU:
			d.texture.TileStyleU, ok = newTileStyle(a.Value)
		case attrTileStyleV:
			d.texture.TileStyleV, ok = newTileStyle(a.Value)
		case attrFilter:
			d.texture.Filter, ok = newTextureFilter(a.Value)
		}
		if err != nil || !ok {
			return errors.New("go3mf: texture2d attribute not valid")
		}
	}
	return nil
}
