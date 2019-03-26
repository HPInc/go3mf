package io3mf

import (
	"encoding/xml"
	"errors"
	"image/color"
	"strconv"

	go3mf "github.com/qmuntal/go3mf"
)

var emptyEntry struct{}

type resourceEntry struct {
	ID, Index uint64
}

type colorMapping struct {
	entries   map[resourceEntry]color.RGBA
	resources map[uint64]struct{}
}

func (m *colorMapping) register(id, index uint64, c color.RGBA) {
	m.entries[resourceEntry{id, index}] = c
	m.resources[id] = emptyEntry
}

func (m *colorMapping) find(id, index uint64) (color.RGBA, bool) {
	if c, ok := m.entries[resourceEntry{id, index}]; ok {
		return c, true
	}
	return defaultColor, false
}

func (m *colorMapping) hasResource(id uint64) bool {
	_, ok := m.resources[id]
	return ok
}

type texCoord struct {
	id   uint64
	u, v float32
}

type texCoordMapping struct {
	entries   map[resourceEntry]texCoord
	resources map[uint64]struct{}
}

func (m *texCoordMapping) register(id, index, textureID uint64, u, v float32) {
	m.entries[resourceEntry{id, index}] = texCoord{textureID, u, v}
	m.resources[id] = emptyEntry
}

func (m *texCoordMapping) find(id, index uint64) (texCoord, bool) {
	if c, ok := m.entries[resourceEntry{id, index}]; ok {
		return c, true
	}
	return texCoord{}, false
}

func (m *texCoordMapping) hasResource(id uint64) bool {
	_, ok := m.resources[id]
	return ok
}

type resourceDecoder struct {
	r               *Reader
	path            string
	colorMapping    colorMapping
	texCoordMapping texCoordMapping
	progressCount   int
}

func (d *resourceDecoder) init() {
	d.colorMapping.entries = make(map[resourceEntry]color.RGBA)
	d.colorMapping.resources = make(map[uint64]struct{})
	d.texCoordMapping.entries = make(map[resourceEntry]texCoord)
	d.texCoordMapping.resources = make(map[uint64]struct{})
}

func (d *resourceDecoder) Decode(x xml.TokenReader, se xml.StartElement) error {
	d.init()
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsCoreSpec {
				err = d.processCoreContent(x, tp)
			} else if tp.Name.Space == nsMaterialSpec {
				err = d.processMaterialContent(x, tp)
			} else if tp.Name.Space == nsSliceSpec {
				err = d.processSliceContent(x, tp)
			}
		case xml.EndElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrResources {
				return nil
			}
		}
		if err != nil {
			return err
		}
	}
}

func (d *resourceDecoder) processCoreContent(x xml.TokenReader, se xml.StartElement) (err error) {
	switch se.Name.Local {
	case attrObject:
		d.progressCount++
		if !d.r.progress.progress(1.0-2.0/float64(d.progressCount+2), StageReadResources) {
			return ErrUserAborted
		}
		d.r.progress.pushLevel(1.0-2.0/float64(d.progressCount+2), 1.0-2.0/float64(d.progressCount+1+2))
		od := objectDecoder{r: d.r, texCoordMapping: &d.texCoordMapping, colorMapping: &d.colorMapping}
		od.obj.ModelPath = d.path
		err = od.Decode(x, se)
		d.r.progress.popLevel()
	case attrBaseMaterials:
		md := baseMaterialsDecoder{r: d.r}
		md.baseMaterials.ModelPath = d.path
		err = md.Decode(x, se)
	}
	return
}

func (d *resourceDecoder) processMaterialContent(x xml.TokenReader, se xml.StartElement) error {
	switch se.Name.Local {
	case attrColorGroup:
		cd := colorGroupDecoder{r: d.r, colorMapping: &d.colorMapping}
		return cd.Decode(x, se)
	case attrTexture2DGroup:
		td := tex2DGroupDecoder{r: d.r, texCoordMapping: &d.texCoordMapping}
		return td.Decode(x, se)
	case attrTexture2D:
		td := texture2DDecoder{r: d.r}
		td.texture.ModelPath = d.path
		return td.Decode(se)
	case attrComposite:
		d.r.Warnings = append(d.r.Warnings, &ReadError{InvalidOptionalValue, "go3mf: composite materials extension not supported"})
	}
	return nil
}

func (d *resourceDecoder) processSliceContent(x xml.TokenReader, se xml.StartElement) error {
	if se.Name.Local != attrSliceStack {
		return nil
	}
	d.progressCount++
	if !d.r.progress.progress(1.0-2.0/float64(d.progressCount+2), StageReadResources) {
		return ErrUserAborted
	}
	d.r.progress.pushLevel(1.0-2.0/float64(d.progressCount+2), 1.0-2.0/float64(d.progressCount+1+2))
	sd := sliceStackDecoder{r: d.r}
	sd.sliceStack.ModelPath = d.path
	err := sd.Decode(x, se)
	d.r.progress.popLevel()
	return err
}

type baseMaterialsDecoder struct {
	r             *Reader
	baseMaterials go3mf.BaseMaterialsResource
}

func (d *baseMaterialsDecoder) parseAttr(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Space != "" || a.Name.Local != attrID {
			continue
		}
		if d.baseMaterials.ID == 0 {
			d.baseMaterials.ID, err = strconv.ParseUint(a.Value, 10, 64)
			if err != nil {
				err = errors.New("go3mf: base materials id is not valid")
			}
		} else {
			err = errors.New("go3mf: duplicated base materials id attribute")
		}
		if err != nil {
			break
		}
	}
	return
}

func (d *baseMaterialsDecoder) Decode(x xml.TokenReader, se xml.StartElement) error {
	if err := d.parseAttr(se.Attr); err != nil {
		return err
	}
	if d.baseMaterials.ID == 0 {
		return errors.New("go3mf: missing base materials resource id attribute")
	}
	if err := d.parseContent(x); err != nil {
		return err
	}
	d.r.addResource(&d.baseMaterials)
	return nil
}

func (d *baseMaterialsDecoder) parseContent(x xml.TokenReader) error {
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrBase {
				if err := d.addBaseMaterial(tp.Attr); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrBaseMaterials {
				return nil
			}
		}
	}
}

func (d *baseMaterialsDecoder) addBaseMaterial(attrs []xml.Attr) error {
	var name string
	var withColor bool
	baseColor := color.RGBA{}
	for _, a := range attrs {
		switch a.Name.Local {
		case attrName:
			name = a.Value
		case attrBaseMaterialColor:
			var err error
			baseColor, err = strToSRGB(a.Value)
			if err != nil {
				return err
			}
			withColor = true
		}
	}
	if name == "" || !withColor {
		return errors.New("go3mf: missing base material attributes")
	}
	d.baseMaterials.Materials = append(d.baseMaterials.Materials, go3mf.BaseMaterial{Name: name, Color: baseColor})
	return nil
}
