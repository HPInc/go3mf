package io3mf

import (
	"encoding/xml"

	go3mf "github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/mesh"
)

type sliceStackDecoder struct {
	emptyDecoder
	progressCount int
	resource      go3mf.SliceStackResource
}

func (d *sliceStackDecoder) Open() {
	d.resource.ModelPath = d.file.path
}

func (d *sliceStackDecoder) Close() bool {
	ok := true
	if len(d.resource.Stack.Refs) > 0 && len(d.resource.Stack.Slices) > 0 {
		ok = d.file.parser.GenericError(true, "slicestack contains slices and slicerefs")
	}
	d.file.AddResource(&d.resource)
	return d.file.parser.CloseResource() && ok
}
func (d *sliceStackDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsSliceSpec {
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
			d.resource.ID, ok = d.file.parser.ParseResourceID(a.Value)
		case attrZBottom:
			d.resource.Stack.BottomZ = d.file.parser.ParseFloat32Optional(a.Name.Local, a.Value)
		}
	}
	return ok
}

type sliceRefDecoder struct {
	emptyDecoder
	resource *go3mf.SliceStackResource
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
			sliceStackID, ok = d.file.parser.ParseUint32Required(a.Name.Local, a.Value)
		case attrSlicePath:
			path = a.Value
		}
	}

	return ok && d.addSliceRef(sliceStackID, path)
}

func (d *sliceRefDecoder) addSliceRef(sliceStackID uint32, path string) bool {
	ok := sliceStackID != 0 || d.file.parser.MissingAttr(attrSliceRefID)
	ok = ok && path != d.resource.ModelPath || d.file.parser.GenericError(true, "a slicepath is invalid")
	resource, ok := d.file.FindResource(path, sliceStackID)
	if !ok {
		ok = d.file.parser.GenericError(true, "non-existent referenced resource")
	} else if _, ok = resource.(*go3mf.SliceStackResource); !ok {
		ok = d.file.parser.GenericError(true, "non-slicestack referenced resource")
	}
	if ok {
		d.resource.Stack.Refs = append(d.resource.Stack.Refs, go3mf.SliceRef{SliceStackID: sliceStackID, Path: path})
	}
	return ok
}

type sliceDecoder struct {
	emptyDecoder
	resource               *go3mf.SliceStackResource
	slice                  mesh.Slice
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
func (d *sliceDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsSliceSpec {
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
			d.slice.TopZ, ok = d.file.parser.ParseFloat32Required(a.Name.Local, a.Value)
			break
		}
	}
	if !hasTopZ {
		ok = d.file.parser.MissingAttr(attrZTop)
	}
	return ok
}

type polygonVerticesDecoder struct {
	emptyDecoder
	slice                *mesh.Slice
	polygonVertexDecoder polygonVertexDecoder
}

func (d *polygonVerticesDecoder) Open() {
	d.polygonVertexDecoder.slice = d.slice
}

func (d *polygonVerticesDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsSliceSpec && name.Local == attrVertex {
		child = &d.polygonVertexDecoder
	}
	return
}

type polygonVertexDecoder struct {
	emptyDecoder
	slice *mesh.Slice
}

func (d *polygonVertexDecoder) Attributes(attrs []xml.Attr) bool {
	var x, y float32
	ok := true
	for _, a := range attrs {
		switch a.Name.Local {
		case attrX:
			x, ok = d.file.parser.ParseFloat32Required(a.Name.Local, a.Value)
		case attrY:
			y, ok = d.file.parser.ParseFloat32Required(a.Name.Local, a.Value)
		}
		if !ok {
			break
		}
	}
	d.slice.AddVertex(x, y)
	return ok
}

type polygonDecoder struct {
	emptyDecoder
	slice                 *mesh.Slice
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
		return d.file.parser.GenericError(true, "a closed slice polygon is actually a line")
	}
	return true
}

func (d *polygonDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsSliceSpec && name.Local == attrSegment {
		child = &d.polygonSegmentDecoder
	}
	return
}

func (d *polygonDecoder) Attributes(attrs []xml.Attr) bool {
	var start uint32
	ok := true
	for _, a := range attrs {
		if a.Name.Local == attrStartV {
			start, ok = d.file.parser.ParseUint32Required(a.Name.Local, a.Value)
			break
		}
	}
	if ok {
		err := d.slice.AddPolygonIndex(d.polygonIndex, int(start))
		if err != nil {
			ok = d.file.parser.GenericError(true, err.Error())
		}
	}
	return ok
}

type polygonSegmentDecoder struct {
	emptyDecoder
	slice        *mesh.Slice
	polygonIndex int
}

func (d *polygonSegmentDecoder) Attributes(attrs []xml.Attr) bool {
	var v2 uint32
	ok := true
	for _, a := range attrs {
		if a.Name.Local == attrV2 {
			v2, ok = d.file.parser.ParseUint32Required(a.Name.Local, a.Value)
			break
		}
	}
	if ok {
		err := d.slice.AddPolygonIndex(d.polygonIndex, int(v2))
		if err != nil {
			ok = d.file.parser.GenericError(true, err.Error())
		}
	}
	return ok
}
