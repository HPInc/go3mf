package slices

import (
	"strconv"

	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
	"github.com/qmuntal/go3mf/spec/xml"
)

func (e Spec) NewResourcesElementDecoder(resources *go3mf.Resources, nodeName string) xml.NodeDecoder {
	if nodeName == attrSliceStack {
		return &sliceStackDecoder{resources: resources}
	}
	return nil
}

func (e Spec) DecodeAttribute(parentNode interface{}, attr xml.Attr) error {
	switch t := parentNode.(type) {
	case *go3mf.Object:
		return objectAttrDecoder(t, attr)
	}
	return nil
}

// objectAttrDecoder decodes the slice attributes of an ObjectReosurce.
func objectAttrDecoder(o *go3mf.Object, a xml.Attr) (err error) {
	switch a.Name.Local {
	case attrSliceRefID:
		val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
		if err1 != nil {
			err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
		}

		var ext *SliceStackInfo
		if o.AnyAttr.Get(&ext) {
			ext.SliceStackID = uint32(val)
		} else {
			o.AnyAttr = append(o.AnyAttr, &SliceStackInfo{SliceStackID: uint32(val)})
		}
	case attrMeshRes:
		res, ok := newMeshResolution(string(a.Value))
		if !ok {
			err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, false))
		}
		var ext *SliceStackInfo
		if o.AnyAttr.Get(&ext) {
			ext.MeshResolution = res
		} else {
			o.AnyAttr = append(o.AnyAttr, &SliceStackInfo{MeshResolution: res})
		}
	}
	return
}

type sliceStackDecoder struct {
	baseDecoder
	resources *go3mf.Resources
	resource  SliceStack
}

func (d *sliceStackDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *sliceStackDecoder) Child(name xml.Name) (child xml.NodeDecoder) {
	if name.Space == Namespace {
		if name.Local == attrSlice {
			child = &sliceDecoder{resource: &d.resource}
		} else if name.Local == attrSliceRef {
			child = &sliceRefDecoder{resource: &d.resource}
		}
	}
	return
}

func (d *sliceStackDecoder) Start(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		switch a.Name.Local {
		case attrID:
			id, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
		case attrZBottom:
			val, err1 := strconv.ParseFloat(string(a.Value), 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, false))
			}
			d.resource.BottomZ = float32(val)
		}
	}
	return
}

type sliceRefDecoder struct {
	baseDecoder
	resource *SliceStack
}

func (d *sliceRefDecoder) Start(attrs []xml.Attr) (err error) {
	var (
		sliceStackID uint32
		path         string
	)
	for _, a := range attrs {
		switch a.Name.Local {
		case attrSliceRefID:
			val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			sliceStackID = uint32(val)
		case attrSlicePath:
			path = string(a.Value)
		}
	}
	d.resource.Refs = append(d.resource.Refs, SliceRef{SliceStackID: sliceStackID, Path: path})
	return
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
func (d *sliceDecoder) Child(name xml.Name) (child xml.NodeDecoder) {
	if name.Space == Namespace {
		if name.Local == attrVertices {
			child = &d.polygonVerticesDecoder
		} else if name.Local == attrPolygon {
			child = &d.polygonDecoder
		}
	}
	return
}

func (d *sliceDecoder) Start(attrs []xml.Attr) (err error) {
	d.polygonDecoder.slice = &d.slice
	d.polygonVerticesDecoder.slice = &d.slice
	for _, a := range attrs {
		if a.Name.Local == attrZTop {
			val, err1 := strconv.ParseFloat(string(a.Value), 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.slice.TopZ = float32(val)
			break
		}
	}
	return
}

type polygonVerticesDecoder struct {
	baseDecoder
	slice                *Slice
	polygonVertexDecoder polygonVertexDecoder
}

func (d *polygonVerticesDecoder) Start(_ []xml.Attr) error {
	d.polygonVertexDecoder.slice = d.slice
	return nil
}

func (d *polygonVerticesDecoder) Child(name xml.Name) (child xml.NodeDecoder) {
	if name.Space == Namespace && name.Local == attrVertex {
		child = &d.polygonVertexDecoder
	}
	return
}

type polygonVertexDecoder struct {
	baseDecoder
	slice *Slice
}

func (d *polygonVertexDecoder) Start(attrs []xml.Attr) (err error) {
	var p go3mf.Point2D
	for _, a := range attrs {
		val, err1 := strconv.ParseFloat(string(a.Value), 32)
		if err1 != nil {
			err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
		}
		switch a.Name.Local {
		case attrX:
			p[0] = float32(val)
		case attrY:
			p[1] = float32(val)
		}
	}
	d.slice.Vertices = append(d.slice.Vertices, p)
	return
}

type polygonDecoder struct {
	baseDecoder
	slice                 *Slice
	polygonSegmentDecoder polygonSegmentDecoder
}

func (d *polygonDecoder) Child(name xml.Name) (child xml.NodeDecoder) {
	if name.Space == Namespace && name.Local == attrSegment {
		child = &d.polygonSegmentDecoder
	}
	return
}

func (d *polygonDecoder) Start(attrs []xml.Attr) (err error) {
	polygonIndex := len(d.slice.Polygons)
	d.slice.Polygons = append(d.slice.Polygons, Polygon{})
	d.polygonSegmentDecoder.polygon = &d.slice.Polygons[polygonIndex]
	for _, a := range attrs {
		if a.Name.Local == attrStartV {
			val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.slice.Polygons[polygonIndex].StartV = uint32(val)
			break
		}
	}
	return
}

type polygonSegmentDecoder struct {
	baseDecoder
	polygon *Polygon
}

func (d *polygonSegmentDecoder) Start(attrs []xml.Attr) (err error) {
	var (
		segment      Segment
		hasP1, hasP2 bool
	)
	for _, a := range attrs {
		var required bool
		val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
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
		if err1 != nil {
			err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, required))
		}
	}
	d.polygon.Segments = append(d.polygon.Segments, segment)
	return
}

type baseDecoder struct {
}

func (d *baseDecoder) End() {}
