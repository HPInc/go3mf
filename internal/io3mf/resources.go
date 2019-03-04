package io3mf

import (
	"encoding/xml"
	"errors"
	"image/color"
	"strconv"

	mdl "github.com/qmuntal/go3mf/internal/model"
	"github.com/qmuntal/go3mf/internal/progress"
)

type resourceDecoder struct {
	x             *xml.Decoder
	r             *Reader
	model         *mdl.Model
	progressCount int
}

func (d *resourceDecoder) Decode(se xml.StartElement) error {
	for {
		t, err := d.x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsCoreSpec {
				if err := d.processCoreContent(tp); err != nil {
					return err
				}
			} else if tp.Name.Space == nsMaterialSpec {
				if err := d.processMaterialContent(tp); err != nil {
					return err
				}
			}
		}
	}
}

func (d *resourceDecoder) processCoreContent(se xml.StartElement) error {
	if se.Name.Local == attrObject {
		d.progressCount++
		if !d.r.progress.Progress(1.0-2.0/float64(d.progressCount+2), progress.StageReadResources) {
			return ErrUserAborted
		}
		d.r.progress.PushLevel(1.0-2.0/float64(d.progressCount+2), 1.0-2.0/float64(d.progressCount+1+2))

		d.r.progress.PopLevel()
	} else if se.Name.Local == attrBaseMaterials {
		md := baseMaterialsDecoder{x: d.x, r: d.r, model: d.model}
		if err := md.Decode(se); err != nil {
			return err
		}
	}
	return nil
}

func (d *resourceDecoder) processMaterialContent(se xml.StartElement) error {
	return nil
}

type baseMaterialsDecoder struct {
	x             *xml.Decoder
	r             *Reader
	model         *mdl.Model
	baseMaterials *mdl.BaseMaterialsResource
}

func (d *baseMaterialsDecoder) parseAttr(se xml.StartElement) error {
	for _, a := range se.Attr {
		if a.Name.Space == "" && a.Name.Local == attrID {
			if d.baseMaterials != nil {
				return errors.New("go3mf: duplicated base materials id attribute")
			}
			id, err := strconv.ParseUint(a.Value, 10, 64)
			if err != nil {
				return errors.New("go3mf: base materials id is not valid")
			}
			d.baseMaterials, err = mdl.NewBaseMaterialsResource(id, d.model)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *baseMaterialsDecoder) Decode(se xml.StartElement) error {
	if err := d.parseAttr(se); err != nil {
		return err
	}
	if d.baseMaterials == nil {
		return errors.New("go3mf: missing base materials resource id attribute")
	}
	return d.model.AddResource(d.baseMaterials)
}

func (d *baseMaterialsDecoder) parseContent() error {
	for {
		t, err := d.x.Token()
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
		}
	}
}

func (d *baseMaterialsDecoder) addBaseMaterial(attrs []xml.Attr) error {
	baseMaterial := mdl.BaseMaterial{
		Color: color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	}
	for _, a := range attrs {
		switch a.Name.Local {
		case attrBaseMaterialName:
			baseMaterial.Name = a.Value
		case attrBaseMaterialColor:
			c, err := strToSRGB(a.Value)
			if err != nil {
				return err
			}
			baseMaterial.Color = c
		}
	}
	d.baseMaterials.Materials = append(d.baseMaterials.Materials, &baseMaterial)
	return nil
}
