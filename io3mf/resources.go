package io3mf

import (
	"encoding/xml"
	"image/color"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/iohelper"
)

type resourceDecoder struct {
	iohelper.EmptyDecoder
}

func (d *resourceDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec {
		switch name.Local {
		case attrObject:
			child = &objectDecoder{}
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
		case attrCompositematerials:
			child = new(compositeMaterialsDecoder)
		case attrMultiProps:
			child = new(multiPropertiesDecoder)
		}
	}
	if ext, ok := extensionDecoder[name.Space]; ok {
		child = ext.NodeDecoder(name.Local)
	}
	return
}

type baseMaterialsDecoder struct {
	iohelper.EmptyDecoder
	resource            go3mf.BaseMaterialsResource
	baseMaterialDecoder baseMaterialDecoder
}

func (d *baseMaterialsDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.baseMaterialDecoder.resource = &d.resource
}

func (d *baseMaterialsDecoder) Close() bool {
	ok := d.Scanner.CloseResource()
	if ok {
		d.Scanner.AddResource(&d.resource)
	}
	return ok
}

func (d *baseMaterialsDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrBase {
		child = &d.baseMaterialDecoder
	}
	return
}

func (d *baseMaterialsDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			d.resource.ID, ok = d.Scanner.ParseResourceID(a.Value)
			break
		}
	}
	return ok
}

type baseMaterialDecoder struct {
	iohelper.EmptyDecoder
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
			baseColor, err = iohelper.ReadRGB(a.Value)
			withColor = true
			if err != nil {
				ok = d.Scanner.InvalidRequiredAttr(attrBaseMaterialColor, a.Value)
			}
		}
	}
	if ok {
		if name == "" {
			ok = d.Scanner.MissingAttr(attrName)
		}
		if !withColor {
			ok = d.Scanner.MissingAttr(attrBaseMaterialColor)
		}
		d.resource.Materials = append(d.resource.Materials, go3mf.BaseMaterial{Name: name, Color: baseColor})
	}
	return ok
}
