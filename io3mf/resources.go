package io3mf

import (
	"encoding/xml"
	"image/color"

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
			d.file.AddWarning(&ReadError{InvalidOptionalValue, "go3mf: composite materials extension not supported"})
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

func (d *baseMaterialsDecoder) Open() {
	d.resource.ModelPath = d.file.path
	d.baseMaterialDecoder.resource = &d.resource
}

func (d *baseMaterialsDecoder) Close() bool {
	if !d.file.parser.CloseResource() {
		return false
	}
	d.file.AddResource(&d.resource)
	return true
}

func (d *baseMaterialsDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrBase {
		child = &d.baseMaterialDecoder
	}
	return
}

func (d *baseMaterialsDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space != "" || a.Name.Local != attrID {
			continue
		}
		d.resource.ID, ok = d.file.parser.ParseResourceID(a.Value)
	}
	return ok
}

type baseMaterialDecoder struct {
	emptyDecoder
	resource *go3mf.BaseMaterialsResource
}

func (d *baseMaterialDecoder) Attributes(attrs []xml.Attr) bool {
	var name string
	var withColor bool
	baseColor := color.RGBA{}
	ok := true
	for _, a := range attrs {
		switch a.Name.Local {
		case attrName:
			name = a.Value
		case attrBaseMaterialColor:
			var err error
			baseColor, err = strToSRGB(a.Value)
			withColor = true
			if err != nil {
				ok = d.file.parser.InvalidRequiredAttr(attrBaseMaterialColor)
			}
		}
		if !ok {
			return false
		}
	}
	if name == "" {
		ok = d.file.parser.MissingAttr(attrName)
	}
	if !withColor {
		ok = d.file.parser.MissingAttr(attrBaseMaterialColor)
	}
	d.resource.Materials = append(d.resource.Materials, go3mf.BaseMaterial{Name: name, Color: baseColor})
	return ok
}
