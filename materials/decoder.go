package materials

import (
	"encoding/xml"
	"strings"

	"github.com/qmuntal/go3mf"
)

func init() {
	go3mf.RegisterNewNodeDecoder(ExtensionName, nodeDecoder)
}

func nodeDecoder(_ interface{}, nodeName string) (child go3mf.NodeDecoder) {
	switch nodeName {
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
	return
}

type colorGroupDecoder struct {
	go3mf.BaseDecoder
	resource     ColorGroupResource
	colorDecoder colorDecoder
}

func (d *colorGroupDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.colorDecoder.resource = &d.resource
}

func (d *colorGroupDecoder) Close() {
	d.Scanner.AddResource(&d.resource)
	d.Scanner.CloseResource()
}

func (d *colorGroupDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrColor {
		child = &d.colorDecoder
	}
	return
}

func (d *colorGroupDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			d.resource.ID = d.Scanner.ParseResourceID(a.Value)
			break
		}
	}
}

type colorDecoder struct {
	go3mf.BaseDecoder
	resource *ColorGroupResource
}

func (d *colorDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrColor {
			c, err := go3mf.ParseRGB(a.Value)
			if err != nil {
				d.Scanner.InvalidRequiredAttr(attrColor, a.Value)
			}
			d.resource.Colors = append(d.resource.Colors, c)
		}
	}
}

type tex2DCoordDecoder struct {
	go3mf.BaseDecoder
	resource *Texture2DGroupResource
}

func (d *tex2DCoordDecoder) Attributes(attrs []xml.Attr) {
	var u, v float32
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrU:
			u = d.Scanner.ParseFloat32Required(attrU, a.Value)
		case attrV:
			v = d.Scanner.ParseFloat32Required(attrV, a.Value)
		}
	}
	d.resource.Coords = append(d.resource.Coords, TextureCoord{float32(u), float32(v)})
}

type tex2DGroupDecoder struct {
	go3mf.BaseDecoder
	resource          Texture2DGroupResource
	tex2DCoordDecoder tex2DCoordDecoder
}

func (d *tex2DGroupDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.tex2DCoordDecoder.resource = &d.resource
}

func (d *tex2DGroupDecoder) Close() {
	d.Scanner.AddResource(&d.resource)
	d.Scanner.CloseResource()
}

func (d *tex2DGroupDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrTex2DCoord {
		child = &d.tex2DCoordDecoder
	}
	return
}

func (d *tex2DGroupDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			d.resource.ID = d.Scanner.ParseResourceID(a.Value)
		case attrTexID:
			d.resource.TextureID = d.Scanner.ParseUint32Required(attrTexID, a.Value)
		}
	}
}

type texture2DDecoder struct {
	go3mf.BaseDecoder
	resource Texture2DResource
}

func (d *texture2DDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
}

func (d *texture2DDecoder) Close() {
	d.Scanner.AddResource(&d.resource)
	d.Scanner.CloseResource()
}

func (d *texture2DDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			d.resource.ID = d.Scanner.ParseResourceID(a.Value)
		case attrPath:
			d.resource.Path = a.Value
		case attrContentType:
			d.resource.ContentType, _ = newTexture2DType(a.Value)
		case attrTileStyleU:
			d.resource.TileStyleU, _ = newTileStyle(a.Value)
		case attrTileStyleV:
			d.resource.TileStyleV, _ = newTileStyle(a.Value)
		case attrFilter:
			d.resource.Filter, _ = newTextureFilter(a.Value)
		}
	}
	if d.resource.Path == "" {
		d.Scanner.MissingAttr(attrPath)
	}
}

type compositeMaterialsDecoder struct {
	go3mf.BaseDecoder
	resource         CompositeMaterialsResource
	compositeDecoder compositeDecoder
}

func (d *compositeMaterialsDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.compositeDecoder.resource = &d.resource
}

func (d *compositeMaterialsDecoder) Close() {
	d.Scanner.AddResource(&d.resource)
	d.Scanner.CloseResource()
}

func (d *compositeMaterialsDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrComposite {
		child = &d.compositeDecoder
	}
	return
}

func (d *compositeMaterialsDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			d.resource.ID = d.Scanner.ParseResourceID(a.Value)
		case attrMatID:
			d.resource.MaterialID = d.Scanner.ParseUint32Required(attrMatID, a.Value)
		case attrMatIndices:
			for _, f := range strings.Fields(a.Value) {
				val := d.Scanner.ParseUint32Required(attrValues, f)
				d.resource.Indices = append(d.resource.Indices, val)
			}
		}
	}
	if d.resource.MaterialID == 0 {
		d.Scanner.MissingAttr(attrMatID)
	}
	if len(d.resource.Indices) == 0 {
		d.Scanner.MissingAttr(attrMatIndices)
	}
}

type compositeDecoder struct {
	go3mf.BaseDecoder
	resource *CompositeMaterialsResource
}

func (d *compositeDecoder) Attributes(attrs []xml.Attr) {
	composite := Composite{}
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrValues {
			for _, f := range strings.Fields(a.Value) {
				val := d.Scanner.ParseFloat32Required(attrValues, f)
				composite.Values = append(composite.Values, val)
			}
		}
	}
	if len(composite.Values) == 0 {
		d.Scanner.MissingAttr(attrValues)
	}
	d.resource.Composites = append(d.resource.Composites, composite)
}

type multiPropertiesDecoder struct {
	go3mf.BaseDecoder
	resource     MultiPropertiesResource
	multiDecoder multiDecoder
}

func (d *multiPropertiesDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.multiDecoder.resource = &d.resource
}

func (d *multiPropertiesDecoder) Close() {
	d.Scanner.AddResource(&d.resource)
	d.Scanner.CloseResource()
}

func (d *multiPropertiesDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrMulti {
		child = &d.multiDecoder
	}
	return
}

func (d *multiPropertiesDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			d.resource.ID = d.Scanner.ParseResourceID(a.Value)
		case attrBlendMethods:
			for _, f := range strings.Fields(a.Value) {
				val, _ := newBlendMethod(f)
				d.resource.BlendMethods = append(d.resource.BlendMethods, val)
			}
		case attrPIDs:
			for _, f := range strings.Fields(a.Value) {
				val := d.Scanner.ParseUint32Required(attrPIDs, f)
				d.resource.Resources = append(d.resource.Resources, val)
			}
		}
	}
	if len(d.resource.Resources) == 0 {
		d.Scanner.MissingAttr(attrPIDs)
	}
}

type multiDecoder struct {
	go3mf.BaseDecoder
	resource *MultiPropertiesResource
}

func (d *multiDecoder) Attributes(attrs []xml.Attr) {
	multi := Multi{}
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrPIndices {
			for _, f := range strings.Fields(a.Value) {
				val := d.Scanner.ParseUint32Required(attrPIndices, f)
				multi.ResourceIndices = append(multi.ResourceIndices, val)
			}
		}
	}
	if len(multi.ResourceIndices) == 0 {
		d.Scanner.MissingAttr(attrPIndices)
	}
	d.resource.Multis = append(d.resource.Multis, multi)
}
