package io3mf

import (
	"encoding/xml"
	"errors"
	"image/color"
	"strconv"

	go3mf "github.com/qmuntal/go3mf"
)

type resourceDecoder struct {
	emptyDecoder
	progressCount int
}

func (d *resourceDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec {
		switch name.Local {
		case attrObject:
			d.progressCount++
			child = &objectDecoder{progressCount: d.progressCount}
		case attrBaseMaterials:
			child = new(baseMaterialsDecoder)
		}
	} else if name.Space == nsMaterialSpec {
		switch name.Local {
		case attrColorGroup:
			child = new(colorGroupDecoder)
		case attrTexture2DGroup:
			child = new(tex2DGroupDecoder)
		case attrTexture2D:
			child = new(texture2DDecoder)
		case attrComposite:
			d.ModelFile().AddWarning(&ReadError{InvalidOptionalValue, "go3mf: composite materials extension not supported"})
		}
	} else if name.Space == nsSliceSpec && name.Local == attrSliceStack {
		d.progressCount++
		child = &sliceStackDecoder{progressCount: d.progressCount}
	}
	return
}

type baseMaterialsDecoder struct {
	emptyDecoder
	resource            go3mf.BaseMaterialsResource
	baseMaterialDecoder baseMaterialDecoder
}

func (d *baseMaterialsDecoder) Open() error {
	d.resource.ModelPath = d.ModelFile().Path()
	d.baseMaterialDecoder.resource = &d.resource
	return nil
}

func (d *baseMaterialsDecoder) Close() error {
	if d.resource.ID == 0 {
		return errors.New("go3mf: missing base materials resource id attribute")
	}
	d.ModelFile().AddResource(&d.resource)
	return nil
}

func (d *baseMaterialsDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrBase {
		child = &d.baseMaterialDecoder
	}
	return
}

func (d *baseMaterialsDecoder) Attributes(attrs []xml.Attr) (err error) {
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

type baseMaterialDecoder struct {
	emptyDecoder
	resource *go3mf.BaseMaterialsResource
}

func (d *baseMaterialDecoder) Attributes(attrs []xml.Attr) (err error) {
	var name string
	var withColor bool
	baseColor := color.RGBA{}

	for _, a := range attrs {
		switch a.Name.Local {
		case attrName:
			name = a.Value
		case attrBaseMaterialColor:
			baseColor, err = strToSRGB(a.Value)
			if err != nil {
				return
			}
			withColor = true
		}
	}
	if name == "" || !withColor {
		return errors.New("go3mf: missing base material attributes")
	}
	d.resource.Materials = append(d.resource.Materials, go3mf.BaseMaterial{Name: name, Color: baseColor})
	return
}
