package materials

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/qmuntal/go3mf"
)

// RegisterExtension registers this extension in the decoder instance.
func RegisterExtension(d *go3mf.Decoder) {
	d.RegisterNodeDecoderExtension(ExtensionName, nodeDecoder)
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
	baseDecoder
	resource     ColorGroupResource
	colorDecoder colorDecoder
}

func (d *colorGroupDecoder) End() {
	d.Scanner.AddAsset(&d.resource)
}

func (d *colorGroupDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrColor {
		child = &d.colorDecoder
	}
	return
}

func (d *colorGroupDecoder) Start(attrs []xml.Attr) {
	d.colorDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			id, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			d.resource.ID, d.Scanner.ResourceID = uint32(id), uint32(id)
			break
		}
	}
}

type colorDecoder struct {
	baseDecoder
	resource *ColorGroupResource
}

func (d *colorDecoder) Start(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrColor {
			c, err := go3mf.ParseRGBA(a.Value)
			if err != nil {
				d.Scanner.InvalidAttr(attrColor, a.Value, true)
			}
			d.resource.Colors = append(d.resource.Colors, c)
		}
	}
}

type tex2DCoordDecoder struct {
	baseDecoder
	resource *Texture2DGroupResource
}

func (d *tex2DCoordDecoder) Start(attrs []xml.Attr) {
	var u, v float32
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		val, err := strconv.ParseFloat(a.Value, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
		}
		switch a.Name.Local {
		case attrU:
			u = float32(val)
		case attrV:
			v = float32(val)
		}
	}
	d.resource.Coords = append(d.resource.Coords, TextureCoord{u, v})
}

type tex2DGroupDecoder struct {
	baseDecoder
	resource          Texture2DGroupResource
	tex2DCoordDecoder tex2DCoordDecoder
}

func (d *tex2DGroupDecoder) End() {
	d.Scanner.AddAsset(&d.resource)
}

func (d *tex2DGroupDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrTex2DCoord {
		child = &d.tex2DCoordDecoder
	}
	return
}

func (d *tex2DGroupDecoder) Start(attrs []xml.Attr) {
	d.tex2DCoordDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			d.resource.ID, d.Scanner.ResourceID = uint32(id), uint32(id)
		case attrTexID:
			val, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			d.resource.TextureID = uint32(val)
		}
	}
}

type texture2DDecoder struct {
	baseDecoder
	resource Texture2DResource
}

func (d *texture2DDecoder) End() {
	d.Scanner.AddAsset(&d.resource)
}

func (d *texture2DDecoder) Start(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			d.resource.ID, d.Scanner.ResourceID = uint32(id), uint32(id)
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
	baseDecoder
	resource         CompositeMaterialsResource
	compositeDecoder compositeDecoder
}

func (d *compositeMaterialsDecoder) End() {
	d.Scanner.AddAsset(&d.resource)
}

func (d *compositeMaterialsDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrComposite {
		child = &d.compositeDecoder
	}
	return
}

func (d *compositeMaterialsDecoder) Start(attrs []xml.Attr) {
	d.compositeDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			d.resource.ID, d.Scanner.ResourceID = uint32(id), uint32(id)
		case attrMatID:
			val, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			d.resource.MaterialID = uint32(val)
		case attrMatIndices:
			for _, f := range strings.Fields(a.Value) {
				val, err := strconv.ParseUint(f, 10, 32)
				if err != nil {
					d.Scanner.InvalidAttr(a.Name.Local, f, true)
				}
				d.resource.Indices = append(d.resource.Indices, uint32(val))
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
	baseDecoder
	resource *CompositeMaterialsResource
}

func (d *compositeDecoder) Start(attrs []xml.Attr) {
	composite := Composite{}
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrValues {
			for _, f := range strings.Fields(a.Value) {
				val, err := strconv.ParseFloat(f, 32)
				if err != nil {
					d.Scanner.InvalidAttr(a.Name.Local, f, true)
				}
				composite.Values = append(composite.Values, float32(val))
			}
		}
	}
	if len(composite.Values) == 0 {
		d.Scanner.MissingAttr(attrValues)
	}
	d.resource.Composites = append(d.resource.Composites, composite)
}

type multiPropertiesDecoder struct {
	baseDecoder
	resource     MultiPropertiesResource
	multiDecoder multiDecoder
}

func (d *multiPropertiesDecoder) End() {
	d.Scanner.AddAsset(&d.resource)
}

func (d *multiPropertiesDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrMulti {
		child = &d.multiDecoder
	}
	return
}

func (d *multiPropertiesDecoder) Start(attrs []xml.Attr) {
	d.multiDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			d.resource.ID, d.Scanner.ResourceID = uint32(id), uint32(id)
		case attrBlendMethods:
			for _, f := range strings.Fields(a.Value) {
				val, _ := newBlendMethod(f)
				d.resource.BlendMethods = append(d.resource.BlendMethods, val)
			}
		case attrPIDs:
			for _, f := range strings.Fields(a.Value) {
				val, err := strconv.ParseUint(f, 10, 32)
				if err != nil {
					d.Scanner.InvalidAttr(a.Name.Local, f, true)
				}
				d.resource.Resources = append(d.resource.Resources, uint32(val))
			}
		}
	}
	if len(d.resource.Resources) == 0 {
		d.Scanner.MissingAttr(attrPIDs)
	}
}

type multiDecoder struct {
	baseDecoder
	resource *MultiPropertiesResource
}

func (d *multiDecoder) Start(attrs []xml.Attr) {
	multi := Multi{}
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrPIndices {
			for _, f := range strings.Fields(a.Value) {
				val, err := strconv.ParseUint(f, 10, 32)
				if err != nil {
					d.Scanner.InvalidAttr(a.Name.Local, f, true)
				}
				multi.PIndex = append(multi.PIndex, uint32(val))
			}
		}
	}
	if len(multi.PIndex) == 0 {
		d.Scanner.MissingAttr(attrPIndices)
	}
	d.resource.Multis = append(d.resource.Multis, multi)
}

type baseDecoder struct {
	Scanner *go3mf.Scanner
}

func (d *baseDecoder) Text([]byte)                      {}
func (d *baseDecoder) Child(xml.Name) go3mf.NodeDecoder { return nil }
func (d *baseDecoder) End()                             {}
func (d *baseDecoder) SetScanner(s *go3mf.Scanner)      { d.Scanner = s }
