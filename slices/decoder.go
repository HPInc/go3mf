package slices

import (
	"encoding/xml"
	"strconv"

	"github.com/qmuntal/go3mf"
)

func (e Spec) OnDecoded(_ *go3mf.Model) error {
	return nil
}

func (e Spec) NewNodeDecoder(_ interface{}, nodeName string) go3mf.NodeDecoder {
	if nodeName == attrSliceStack {
		return &sliceStackDecoder{}
	}
	return nil
}

func (e Spec) DecodeAttribute(s *go3mf.Scanner, parentNode interface{}, attr xml.Attr) {
	switch t := parentNode.(type) {
	case *go3mf.Object:
		objectAttrDecoder(s, t, attr)
	}
}

// objectAttrDecoder decodes the slice attributes of an ObjectReosurce.
func objectAttrDecoder(scanner *go3mf.Scanner, o *go3mf.Object, a xml.Attr) {
	switch a.Name.Local {
	case attrSliceRefID:
		val, err := strconv.ParseUint(a.Value, 10, 32)
		if err != nil {
			scanner.InvalidAttr(a.Name.Local, true)
		}

		var ext *SliceStackInfo
		if o.AnyAttr.Get(&ext) {
			ext.SliceStackID = uint32(val)
		} else {
			o.AnyAttr = append(o.AnyAttr, &SliceStackInfo{SliceStackID: uint32(val)})
		}
	case attrMeshRes:
		res, ok := newMeshResolution(a.Value)
		if !ok {
			scanner.InvalidAttr(attrMeshRes, false)
		}
		var ext *SliceStackInfo
		if o.AnyAttr.Get(&ext) {
			ext.MeshResolution = res
		} else {
			o.AnyAttr = append(o.AnyAttr, &SliceStackInfo{MeshResolution: res})
		}
	}
}

type sliceStackDecoder struct {
	baseDecoder
	resource SliceStack
}

func (d *sliceStackDecoder) End() {
	d.Scanner.AddAsset(&d.resource)
}

func (d *sliceStackDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == Namespace {
		if name.Local == attrSlice {
			child = &sliceDecoder{resource: &d.resource}
		} else if name.Local == attrSliceRef {
			child = &sliceRefDecoder{resource: &d.resource}
		}
	}
	return
}

func (d *sliceStackDecoder) Start(attrs []xml.Attr) {
	for _, a := range attrs {
		switch a.Name.Local {
		case attrID:
			id, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, true)
			}
			d.resource.ID, d.Scanner.ResourceID = uint32(id), uint32(id)
		case attrZBottom:
			val, err := strconv.ParseFloat(a.Value, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, false)
			}
			d.resource.BottomZ = float32(val)
		}
	}
}

type sliceRefDecoder struct {
	baseDecoder
	resource *SliceStack
}

func (d *sliceRefDecoder) Start(attrs []xml.Attr) {
	var (
		sliceStackID uint32
		path         string
	)
	for _, a := range attrs {
		switch a.Name.Local {
		case attrSliceRefID:
			val, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, true)
			}
			sliceStackID = uint32(val)
		case attrSlicePath:
			path = a.Value
		}
	}
	d.resource.Refs = append(d.resource.Refs, SliceRef{SliceStackID: sliceStackID, Path: path})
}

type sliceDecoder struct {
	baseDecoder
	resource               *SliceStack
	slice                  Slice
	polygonDecoder         polygonDecoder
	polygonVerticesDecoder polygonVerticesDecoder
}

func (d *sliceDecoder) End() {
	d.resource.Slices = append(d.resource.Slices, &d.slice)
}
func (d *sliceDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == Namespace {
		if name.Local == attrVertices {
			child = &d.polygonVerticesDecoder
		} else if name.Local == attrPolygon {
			child = &d.polygonDecoder
		}
	}
	return
}

func (d *sliceDecoder) Start(attrs []xml.Attr) {
	d.polygonDecoder.slice = &d.slice
	d.polygonVerticesDecoder.slice = &d.slice
	for _, a := range attrs {
		if a.Name.Local == attrZTop {
			val, err := strconv.ParseFloat(a.Value, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, true)
			}
			d.slice.TopZ = float32(val)
			break
		}
	}
}

type polygonVerticesDecoder struct {
	baseDecoder
	slice                *Slice
	polygonVertexDecoder polygonVertexDecoder
}

func (d *polygonVerticesDecoder) Start(_ []xml.Attr) {
	d.polygonVertexDecoder.slice = d.slice
}

func (d *polygonVerticesDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == Namespace && name.Local == attrVertex {
		child = &d.polygonVertexDecoder
	}
	return
}

type polygonVertexDecoder struct {
	baseDecoder
	slice *Slice
}

func (d *polygonVertexDecoder) Start(attrs []xml.Attr) {
	var p go3mf.Point2D
	for _, a := range attrs {
		val, err := strconv.ParseFloat(a.Value, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, true)
		}
		switch a.Name.Local {
		case attrX:
			p[0] = float32(val)
		case attrY:
			p[1] = float32(val)
		}
	}

	d.slice.Vertices = append(d.slice.Vertices, p)
}

type polygonDecoder struct {
	baseDecoder
	slice                 *Slice
	polygonSegmentDecoder polygonSegmentDecoder
}

func (d *polygonDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == Namespace && name.Local == attrSegment {
		child = &d.polygonSegmentDecoder
	}
	return
}

func (d *polygonDecoder) Start(attrs []xml.Attr) {
	polygonIndex := len(d.slice.Polygons)
	d.slice.Polygons = append(d.slice.Polygons, Polygon{})
	d.polygonSegmentDecoder.polygon = &d.slice.Polygons[polygonIndex]
	for _, a := range attrs {
		if a.Name.Local == attrStartV {
			val, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, true)
			}
			d.slice.Polygons[polygonIndex].StartV = uint32(val)
			break
		}
	}
}

type polygonSegmentDecoder struct {
	baseDecoder
	polygon *Polygon
}

func (d *polygonSegmentDecoder) Start(attrs []xml.Attr) {
	var (
		segment      Segment
		hasP1, hasP2 bool
	)
	for _, a := range attrs {
		var required bool
		val, err := strconv.ParseUint(a.Value, 10, 32)
		switch a.Name.Local {
		case attrV2:
			segment.V2 = uint32(val)
			required = true
		case attrPID:
			segment.PID = uint32(val)
		case attrP1:
			segment.P1 = uint32(val)
			hasP1 = true
		case attrP2:
			segment.P2 = uint32(val)
			hasP2 = true
		}
		if hasP1 && !hasP2 {
			segment.P2 = segment.P1
		}
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, required)
		}
	}
	d.polygon.Segments = append(d.polygon.Segments, segment)
}

type baseDecoder struct {
	Scanner *go3mf.Scanner
}

func (d *baseDecoder) Text([]byte)                      {}
func (d *baseDecoder) Child(xml.Name) go3mf.NodeDecoder { return nil }
func (d *baseDecoder) End()                             {}
func (d *baseDecoder) SetScanner(s *go3mf.Scanner)      { d.Scanner = s }
