package slices

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
)

// RegisterExtension registers this extension in the decoder instance.
func RegisterExtension(d *go3mf.Decoder) {
	d.RegisterNodeDecoderExtension(ExtensionName, nodeDecoder)
	d.RegisterDecodeAttributeExtension(ExtensionName, decodeAttribute)
}

func nodeDecoder(_ interface{}, nodeName string) go3mf.NodeDecoder {
	if nodeName == attrSliceStack {
		return &sliceStackDecoder{}
	}
	return nil
}

func decodeAttribute(s *go3mf.Scanner, parentNode interface{}, attr xml.Attr) {
	switch t := parentNode.(type) {
	case *go3mf.ObjectResource:
		objectAttrDecoder(s, t, attr)
	}
}

// objectAttrDecoder decodes the slice attributes of an ObjectReosurce.
func objectAttrDecoder(scanner *go3mf.Scanner, o *go3mf.ObjectResource, attr xml.Attr) {
	switch attr.Name.Local {
	case attrSliceRefID:
		ObjectSliceStackInfo(o).SliceStackID = scanner.ParseUint32(attrSliceRefID, attr.Value)
	case attrMeshRes:
		var ok bool
		ObjectSliceStackInfo(o).SliceResolution, ok = newSliceResolution(attr.Value)
		if !ok {
			scanner.InvalidAttr(attrMeshRes, attr.Value, false)
		}
	}
}

type sliceStackDecoder struct {
	baseDecoder
	resource SliceStackResource
}

func (d *sliceStackDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
}

func (d *sliceStackDecoder) Close() {
	if len(d.resource.Stack.Refs) > 0 && len(d.resource.Stack.Slices) > 0 {
		d.Scanner.GenericError(true, "slicestack contains slices and slicerefs")
	}
	d.Scanner.AddResource(&d.resource)
}

func (d *sliceStackDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName {
		if name.Local == attrSlice {
			child = &sliceDecoder{resource: &d.resource}
		} else if name.Local == attrSliceRef {
			child = &sliceRefDecoder{resource: &d.resource}
		}
	}
	return
}

func (d *sliceStackDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		switch a.Name.Local {
		case attrID:
			d.resource.ID = d.Scanner.ParseResourceID(a.Value)
		case attrZBottom:
			d.resource.Stack.BottomZ = d.Scanner.ParseFloat32Optional(attrZBottom, a.Value)
		}
	}
}

type sliceRefDecoder struct {
	baseDecoder
	resource *SliceStackResource
}

func (d *sliceRefDecoder) Attributes(attrs []xml.Attr) {
	var (
		sliceStackID uint32
		path         string
	)
	for _, a := range attrs {
		switch a.Name.Local {
		case attrSliceRefID:
			sliceStackID = d.Scanner.ParseUint32(attrSliceRefID, a.Value)
		case attrSlicePath:
			path = a.Value
		}
	}
	if sliceStackID == 0 {
		d.Scanner.MissingAttr(attrSliceRefID)
	}
	if path == d.resource.ModelPath {
		d.Scanner.GenericError(true, "a slicepath is invalid")
	}
	d.resource.Stack.Refs = append(d.resource.Stack.Refs, SliceRef{SliceStackID: sliceStackID, Path: path})
}

// TODO: validate coherency after decoding.
// func (d *sliceRefDecoder) addSliceRef(sliceStackID uint32, path string) {
// 	resource, exist := d.Scanner.FindResource(path, sliceStackID)
// 	if !exist {
// 		d.Scanner.GenericError(true, "non-existent referenced resource")
// 	} else if _, isSlice := resource.(*SliceStackResource); !isSlice {
// 		d.Scanner.GenericError(true, "non-slicestack referenced resource")
// 	}
// }

type sliceDecoder struct {
	baseDecoder
	resource               *SliceStackResource
	slice                  Slice
	polygonDecoder         polygonDecoder
	polygonVerticesDecoder polygonVerticesDecoder
}

func (d *sliceDecoder) Open() {
	d.polygonDecoder.slice = &d.slice
	d.polygonVerticesDecoder.slice = &d.slice
}
func (d *sliceDecoder) Close() {
	d.resource.Stack.Slices = append(d.resource.Stack.Slices, &d.slice)
}
func (d *sliceDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName {
		if name.Local == attrVertices {
			child = &d.polygonVerticesDecoder
		} else if name.Local == attrPolygon {
			child = &d.polygonDecoder
		}
	}
	return
}

func (d *sliceDecoder) Attributes(attrs []xml.Attr) {
	var hasTopZ bool
	for _, a := range attrs {
		if a.Name.Local == attrZTop {
			hasTopZ = true
			d.slice.TopZ = d.Scanner.ParseFloat32(attrZTop, a.Value)
			break
		}
	}
	if !hasTopZ {
		d.Scanner.MissingAttr(attrZTop)
	}
}

type polygonVerticesDecoder struct {
	baseDecoder
	slice                *Slice
	polygonVertexDecoder polygonVertexDecoder
}

func (d *polygonVerticesDecoder) Open() {
	d.polygonVertexDecoder.slice = d.slice
}

func (d *polygonVerticesDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrVertex {
		child = &d.polygonVertexDecoder
	}
	return
}

type polygonVertexDecoder struct {
	baseDecoder
	slice *Slice
}

func (d *polygonVertexDecoder) Attributes(attrs []xml.Attr) {
	var x, y float32
	for _, a := range attrs {
		switch a.Name.Local {
		case attrX:
			x = d.Scanner.ParseFloat32(attrX, a.Value)
		case attrY:
			y = d.Scanner.ParseFloat32(attrY, a.Value)
		}
	}
	d.slice.AddVertex(x, y)
}

type polygonDecoder struct {
	baseDecoder
	slice                 *Slice
	polygonIndex          int
	polygonSegmentDecoder polygonSegmentDecoder
}

func (d *polygonDecoder) Open() {
	d.polygonIndex = d.slice.BeginPolygon()
	d.polygonSegmentDecoder.slice = d.slice
	d.polygonSegmentDecoder.polygonIndex = d.polygonIndex
}

func (d *polygonDecoder) Close() {
	if !d.slice.IsPolygonValid(d.polygonIndex) {
		d.Scanner.GenericError(true, "a closed slice polygon is actually a line")
	}
}

func (d *polygonDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrSegment {
		child = &d.polygonSegmentDecoder
	}
	return
}

func (d *polygonDecoder) Attributes(attrs []xml.Attr) {
	var start uint32
	for _, a := range attrs {
		if a.Name.Local == attrStartV {
			start = d.Scanner.ParseUint32(attrStartV, a.Value)
			break
		}
	}
	err := d.slice.AddPolygonIndex(d.polygonIndex, int(start))
	if err != nil {
		d.Scanner.GenericError(true, err.Error())
	}
}

type polygonSegmentDecoder struct {
	baseDecoder
	slice        *Slice
	polygonIndex int
}

func (d *polygonSegmentDecoder) Attributes(attrs []xml.Attr) {
	var v2 uint32
	for _, a := range attrs {
		if a.Name.Local == attrV2 {
			v2 = d.Scanner.ParseUint32(attrV2, a.Value)
			break
		}
	}
	err := d.slice.AddPolygonIndex(d.polygonIndex, int(v2))
	if err != nil {
		d.Scanner.GenericError(true, err.Error())
	}
}

type baseDecoder struct {
	Scanner *go3mf.Scanner
}

func (d *baseDecoder) Open()                            {}
func (d *baseDecoder) Attributes([]xml.Attr)            {}
func (d *baseDecoder) Text([]byte)                      {}
func (d *baseDecoder) Child(xml.Name) go3mf.NodeDecoder { return nil }
func (d *baseDecoder) Close()                           {}
func (d *baseDecoder) SetScanner(s *go3mf.Scanner)      { d.Scanner = s }
