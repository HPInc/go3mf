// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package slices

import (
	"encoding/xml"
	"strconv"

	"github.com/hpinc/go3mf"
	specerr "github.com/hpinc/go3mf/errors"
	"github.com/hpinc/go3mf/spec"
)

func (Spec) NewElementDecoder(parent interface{}, name string) spec.ElementDecoder {
	if name == attrSliceStack {
		return &sliceStackDecoder{resources: parent.(*go3mf.Resources)}
	}
	return nil
}

func (Spec) NewAttrGroup(parent xml.Name) spec.AttrGroup {
	if parent.Space == go3mf.Namespace {
		switch parent.Local {
		case "object":
			return new(ObjectAttr)
		}
	}
	return nil
}

func (u *ObjectAttr) Unmarshal3MFAttr(a spec.XMLAttr) error {
	switch a.Name.Local {
	case attrSliceRefID:
		val, err := strconv.ParseUint(string(a.Value), 10, 32)
		if err != nil {
			return specerr.NewParseAttrError(a.Name.Local, true)
		}
		u.SliceStackID = uint32(val)
	case attrMeshRes:
		res, ok := newMeshResolution(string(a.Value))
		if !ok {
			return specerr.NewParseAttrError(a.Name.Local, false)
		}
		u.MeshResolution = res
	}
	return nil
}

type sliceStackDecoder struct {
	baseDecoder
	resources *go3mf.Resources
	resource  SliceStack
}

func (d *sliceStackDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *sliceStackDecoder) Child(name xml.Name) (i int, child spec.ElementDecoder) {
	if name.Space == Namespace {
		if name.Local == attrSlice {
			child = &sliceDecoder{resource: &d.resource}
			i = len(d.resource.Slices)
		} else if name.Local == attrSliceRef {
			child = &sliceRefDecoder{resource: &d.resource}
			i = len(d.resource.Refs)
		}
	}
	return
}

func (d *sliceStackDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	for _, a := range attrs {
		switch a.Name.Local {
		case attrID:
			id, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
		case attrZBottom:
			val, err := strconv.ParseFloat(string(a.Value), 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
			}
			d.resource.BottomZ = float32(val)
		}
	}
	return errs
}

type sliceRefDecoder struct {
	baseDecoder
	resource *SliceStack
}

func (d *sliceRefDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		sliceStackID uint32
		path         string
		errs         error
	)
	for _, a := range attrs {
		switch a.Name.Local {
		case attrSliceRefID:
			val, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			sliceStackID = uint32(val)
		case attrSlicePath:
			path = string(a.Value)
		}
	}
	ref := SliceRef{SliceStackID: sliceStackID, Path: path}
	d.resource.Refs = append(d.resource.Refs, ref)
	return errs
}

type sliceDecoder struct {
	baseDecoder
	resource               *SliceStack
	slice                  Slice
	polygonDecoder         polygonDecoder
	polygonVerticesDecoder polygonVerticesDecoder
}

func (d *sliceDecoder) End() {
	d.resource.Slices = append(d.resource.Slices, d.slice)
}

func (d *sliceDecoder) Child(name xml.Name) (i int, child spec.ElementDecoder) {
	if name.Space == Namespace {
		if name.Local == attrVertices {
			child = &d.polygonVerticesDecoder
			i = -1
		} else if name.Local == attrPolygon {
			child = &d.polygonDecoder
			i = len(d.slice.Polygons)
		}
	}
	return
}

func (d *sliceDecoder) Start(attrs []spec.XMLAttr) error {
	d.polygonDecoder.slice = &d.slice
	d.polygonVerticesDecoder.slice = &d.slice
	var errs error
	for _, a := range attrs {
		if a.Name.Local == attrZTop {
			val, err := strconv.ParseFloat(string(a.Value), 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.slice.TopZ = float32(val)
			break
		}
	}
	return errs
}

type polygonVerticesDecoder struct {
	baseDecoder
	slice                *Slice
	polygonVertexDecoder polygonVertexDecoder
}

func (d *polygonVerticesDecoder) Start(_ []spec.XMLAttr) error {
	d.polygonVertexDecoder.slice = d.slice
	return nil
}

func (d *polygonVerticesDecoder) Child(name xml.Name) (i int, child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrVertex {
		child = &d.polygonVertexDecoder
		i = len(d.slice.Vertices.Vertex)
	}
	return
}

type polygonVertexDecoder struct {
	baseDecoder
	slice *Slice
}

func (d *polygonVertexDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		p    go3mf.Point2D
		errs error
	)
	for _, a := range attrs {
		val, err := strconv.ParseFloat(string(a.Value), 32)
		if err != nil {
			errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
		}
		switch a.Name.Local {
		case attrX:
			p[0] = float32(val)
		case attrY:
			p[1] = float32(val)
		}
	}
	d.slice.Vertices.Vertex = append(d.slice.Vertices.Vertex, p)
	return errs
}

type polygonDecoder struct {
	baseDecoder
	slice                 *Slice
	polygonSegmentDecoder polygonSegmentDecoder
}

func (d *polygonDecoder) Child(name xml.Name) (i int, child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrSegment {
		child = &d.polygonSegmentDecoder
		i = len(d.slice.Polygons)
	}
	return
}

func (d *polygonDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	polygonIndex := len(d.slice.Polygons)
	d.slice.Polygons = append(d.slice.Polygons, Polygon{})
	d.polygonSegmentDecoder.polygon = &d.slice.Polygons[polygonIndex]
	for _, a := range attrs {
		if a.Name.Local == attrStartV {
			val, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.slice.Polygons[polygonIndex].StartV = uint32(val)
			break
		}
	}
	return errs
}

type polygonSegmentDecoder struct {
	baseDecoder
	polygon *Polygon
}

func (d *polygonSegmentDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		segment      Segment
		hasP1, hasP2 bool
		errs         error
	)
	for _, a := range attrs {
		var required bool
		val, err := strconv.ParseUint(string(a.Value), 10, 32)
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
			errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, required))
		}
	}
	d.polygon.Segments = append(d.polygon.Segments, segment)
	return errs
}

type baseDecoder struct {
}

func (d *baseDecoder) End() {}
