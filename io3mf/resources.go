package io3mf

import (
	"encoding/xml"
	"errors"
	"image/color"
	"strconv"

	go3mf "github.com/qmuntal/go3mf"
)

type resourceDecoder struct {
	r             *Reader
	path          string
	progressCount int
}

func (d *resourceDecoder) Open() error {
	if !d.r.progress.progress(0.2, StageReadResources) {
		return ErrUserAborted
	}
	d.r.progress.pushLevel(0.2, 0.9)
	return nil
}

func (d *resourceDecoder) Close() error {
	d.r.progress.popLevel()
	return nil
}

func (d *resourceDecoder) Decode(x xml.TokenReader) error {
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
		od := objectDecoder{r: d.r, resource: go3mf.ObjectResource{ModelPath: d.path}}
		err = od.Decode(x, se.Attr)
		d.r.progress.popLevel()
	case attrBaseMaterials:
		md := baseMaterialsDecoder{r: d.r, resource: go3mf.BaseMaterialsResource{ModelPath: d.path}}
		err = md.Decode(x, se.Attr)
	}
	return
}

func (d *resourceDecoder) processMaterialContent(x xml.TokenReader, se xml.StartElement) error {
	switch se.Name.Local {
	case attrColorGroup:
		cd := colorGroupDecoder{r: d.r, resource: go3mf.ColorGroupResource{ModelPath: d.path}}
		return cd.Decode(x, se.Attr)
	case attrTexture2DGroup:
		td := tex2DGroupDecoder{r: d.r, resource: go3mf.Texture2DGroupResource{ModelPath: d.path}}
		return td.Decode(x, se.Attr)
	case attrTexture2D:
		td := texture2DDecoder{r: d.r, resource: go3mf.Texture2DResource{ModelPath: d.path}}
		return td.Decode(se.Attr)
	case attrComposite:
		d.r.addWarning(&ReadError{InvalidOptionalValue, "go3mf: composite materials extension not supported"})
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
	sd := sliceStackDecoder{r: d.r, resource: go3mf.SliceStackResource{ModelPath: d.path}}
	err := sd.Decode(x, se.Attr)
	d.r.progress.popLevel()
	return err
}

type baseMaterialsDecoder struct {
	r        *Reader
	resource go3mf.BaseMaterialsResource
}

func (d *baseMaterialsDecoder) parseAttr(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Space != "" || a.Name.Local != attrID {
			continue
		}
		if d.resource.ID == 0 {
			d.resource.ID, err = strconv.ParseUint(a.Value, 10, 64)
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

func (d *baseMaterialsDecoder) Decode(x xml.TokenReader, attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	if d.resource.ID == 0 {
		return errors.New("go3mf: missing base materials resource id attribute")
	}
	if err := d.parseContent(x); err != nil {
		return err
	}
	d.r.addResource(&d.resource)
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
	d.resource.Materials = append(d.resource.Materials, go3mf.BaseMaterial{Name: name, Color: baseColor})
	return nil
}
