package slices

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/io3mf"
	"github.com/qmuntal/go3mf/iohelper"
)

func init() {
	io3mf.RegisterExtensionDecoder(ExtensionName, &extensionDecoder{})
}

type extensionDecoder struct{}

func (d *extensionDecoder) NodeDecoder(_ interface{}, nodeName string) iohelper.NodeDecoder {
	if nodeName == attrSliceStack {
		return &sliceStackDecoder{}
	}
	return nil
}

func (d *extensionDecoder) DecodeAttribute(s *iohelper.Scanner, parentNode interface{}, attr xml.Attr) bool {
	switch t := parentNode.(type) {
	case *go3mf.ObjectResource:
		return objectAttrDecoder(s, t, attr)
	}
	return true
}

// objectAttrDecoder decodes the slice attributes of an ObjectReosurce.
func objectAttrDecoder(scanner *iohelper.Scanner, o *go3mf.ObjectResource, attr xml.Attr) bool {
	ok := true
	switch attr.Name.Local {
	case attrSliceRefID:
		ExtensionObjectResource(o).SliceStackID, ok = scanner.ParseUint32Required(attrSliceRefID, attr.Value)
	case attrMeshRes:
		ExtensionObjectResource(o).SliceResolution, ok = newSliceResolution(attr.Value)
		if !ok {
			ok = true
			scanner.InvalidOptionalAttr(attrMeshRes, attr.Value)
		}
	}
	return ok
}

type sliceStackDecoder struct {
	iohelper.EmptyDecoder
	resource SliceStackResource
}

func (d *sliceStackDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
}

func (d *sliceStackDecoder) Close() bool {
	ok := true
	if len(d.resource.Stack.Refs) > 0 && len(d.resource.Stack.Slices) > 0 {
		ok = d.Scanner.GenericError(true, "slicestack contains slices and slicerefs")
	}
	d.Scanner.AddResource(&d.resource)
	return d.Scanner.CloseResource() && ok
}

func (d *sliceStackDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == ExtensionName {
		if name.Local == attrSlice {
			child = &sliceDecoder{resource: &d.resource}
		} else if name.Local == attrSliceRef {
			child = &sliceRefDecoder{resource: &d.resource}
		}
	}
	return
}

func (d *sliceStackDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		switch a.Name.Local {
		case attrID:
			d.resource.ID, ok = d.Scanner.ParseResourceID(a.Value)
		case attrZBottom:
			d.resource.Stack.BottomZ = d.Scanner.ParseFloat32Optional(attrZBottom, a.Value)
		}
	}
	return ok
}

type sliceRefDecoder struct {
	iohelper.EmptyDecoder
	resource *SliceStackResource
}

func (d *sliceRefDecoder) Attributes(attrs []xml.Attr) bool {
	var (
		sliceStackID uint32
		path         string
	)
	ok := true
	for _, a := range attrs {
		switch a.Name.Local {
		case attrSliceRefID:
			sliceStackID, ok = d.Scanner.ParseUint32Required(attrSliceRefID, a.Value)
		case attrSlicePath:
			path = a.Value
		}
	}

	return ok && d.addSliceRef(sliceStackID, path)
}

func (d *sliceRefDecoder) addSliceRef(sliceStackID uint32, path string) bool {
	ok := true
	if sliceStackID == 0 {
		ok = d.Scanner.MissingAttr(attrSliceRefID)
	}
	if ok && path == d.resource.ModelPath {
		ok = d.Scanner.GenericError(true, "a slicepath is invalid")
	}
	if ok {
		resource, exist := d.Scanner.FindResource(path, sliceStackID)
		if !exist {
			ok = d.Scanner.GenericError(true, "non-existent referenced resource")
		} else if _, isSlice := resource.(*SliceStackResource); !isSlice {
			ok = d.Scanner.GenericError(true, "non-slicestack referenced resource")
		}
		if ok {
			d.resource.Stack.Refs = append(d.resource.Stack.Refs, SliceRef{SliceStackID: sliceStackID, Path: path})
		}
	}
	return ok
}

type sliceDecoder struct {
	iohelper.EmptyDecoder
	resource               *SliceStackResource
	slice                  Slice
	polygonDecoder         polygonDecoder
	polygonVerticesDecoder polygonVerticesDecoder
}

func (d *sliceDecoder) Open() {
	d.polygonDecoder.slice = &d.slice
	d.polygonVerticesDecoder.slice = &d.slice
}
func (d *sliceDecoder) Close() bool {
	d.resource.Stack.Slices = append(d.resource.Stack.Slices, &d.slice)
	return true
}
func (d *sliceDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == ExtensionName {
		if name.Local == attrVertices {
			child = &d.polygonVerticesDecoder
		} else if name.Local == attrPolygon {
			child = &d.polygonDecoder
		}
	}
	return
}

func (d *sliceDecoder) Attributes(attrs []xml.Attr) bool {
	var hasTopZ bool
	ok := true
	for _, a := range attrs {
		if a.Name.Local == attrZTop {
			hasTopZ = true
			d.slice.TopZ, ok = d.Scanner.ParseFloat32Required(attrZTop, a.Value)
			break
		}
	}
	if !hasTopZ {
		ok = d.Scanner.MissingAttr(attrZTop)
	}
	return ok
}

type polygonVerticesDecoder struct {
	iohelper.EmptyDecoder
	slice                *Slice
	polygonVertexDecoder polygonVertexDecoder
}

func (d *polygonVerticesDecoder) Open() {
	d.polygonVertexDecoder.slice = d.slice
}

func (d *polygonVerticesDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrVertex {
		child = &d.polygonVertexDecoder
	}
	return
}

type polygonVertexDecoder struct {
	iohelper.EmptyDecoder
	slice *Slice
}

func (d *polygonVertexDecoder) Attributes(attrs []xml.Attr) bool {
	var x, y float32
	ok := true
	for _, a := range attrs {
		switch a.Name.Local {
		case attrX:
			x, ok = d.Scanner.ParseFloat32Required(attrX, a.Value)
		case attrY:
			y, ok = d.Scanner.ParseFloat32Required(attrY, a.Value)
		}
		if !ok {
			break
		}
	}
	d.slice.AddVertex(x, y)
	return ok
}

type polygonDecoder struct {
	iohelper.EmptyDecoder
	slice                 *Slice
	polygonIndex          int
	polygonSegmentDecoder polygonSegmentDecoder
}

func (d *polygonDecoder) Open() {
	d.polygonIndex = d.slice.BeginPolygon()
	d.polygonSegmentDecoder.slice = d.slice
	d.polygonSegmentDecoder.polygonIndex = d.polygonIndex
}

func (d *polygonDecoder) Close() bool {
	if !d.slice.IsPolygonValid(d.polygonIndex) {
		return d.Scanner.GenericError(true, "a closed slice polygon is actually a line")
	}
	return true
}

func (d *polygonDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrSegment {
		child = &d.polygonSegmentDecoder
	}
	return
}

func (d *polygonDecoder) Attributes(attrs []xml.Attr) bool {
	var start uint32
	ok := true
	for _, a := range attrs {
		if a.Name.Local == attrStartV {
			start, ok = d.Scanner.ParseUint32Required(attrStartV, a.Value)
			break
		}
	}
	if ok {
		err := d.slice.AddPolygonIndex(d.polygonIndex, int(start))
		if err != nil {
			ok = d.Scanner.GenericError(true, err.Error())
		}
	}
	return ok
}

type polygonSegmentDecoder struct {
	iohelper.EmptyDecoder
	slice        *Slice
	polygonIndex int
}

func (d *polygonSegmentDecoder) Attributes(attrs []xml.Attr) bool {
	var v2 uint32
	ok := true
	for _, a := range attrs {
		if a.Name.Local == attrV2 {
			v2, ok = d.Scanner.ParseUint32Required(attrV2, a.Value)
			break
		}
	}
	if ok {
		err := d.slice.AddPolygonIndex(d.polygonIndex, int(v2))
		if err != nil {
			ok = d.Scanner.GenericError(true, err.Error())
		}
	}
	return ok
}
