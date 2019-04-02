package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	go3mf "github.com/qmuntal/go3mf"
)

type colorGroupDecoder struct {
	r        *Reader
	resource go3mf.ColorGroupResource
}

func (d *colorGroupDecoder) Decode(x xml.TokenReader, attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	if d.resource.ID == 0 {
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
				d.r.addResource(&d.resource)
				return nil
			}
		}
	}
}

func (d *colorGroupDecoder) parseAttr(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			if d.resource.ID != 0 {
				return errors.New("go3mf: duplicated color group id attribute")
			}
			var err error
			if d.resource.ID, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
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
			d.resource.Colors = append(d.resource.Colors, c)
		}
	}
	return nil
}

type tex2DGroupDecoder struct {
	r        *Reader
	resource go3mf.Texture2DGroupResource
}

func (d *tex2DGroupDecoder) Decode(x xml.TokenReader, attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	if d.resource.ID == 0 {
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
				d.r.addResource(&d.resource)
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
			if d.resource.ID != 0 {
				return errors.New("go3mf: duplicated tex2Coord group id attribute")
			}
			var err error
			if d.resource.ID, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
				return errors.New("go3mf: tex2Coord group id is not valid")
			}
		case attrTexID:
			if d.resource.TextureID != 0 {
				return errors.New("go3mf: duplicated tex2Coord group texid attribute")
			}
			var err error
			if d.resource.TextureID, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
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
	d.resource.Coords = append(d.resource.Coords, go3mf.TextureCoord{float32(u), float32(v)})
	return nil
}

type texture2DDecoder struct {
	r        *Reader
	resource go3mf.Texture2DResource
}

func (d *texture2DDecoder) Decode(attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	if d.resource.ID == 0 {
		return errors.New("go3mf: missing texture2d id attribute")
	}
	d.r.addResource(&d.resource)
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
			if d.resource.ID != 0 {
				err = errors.New("go3mf: duplicated texture2d id attribute")
			} else {
				d.resource.ID, err = strconv.ParseUint(a.Value, 10, 64)
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
