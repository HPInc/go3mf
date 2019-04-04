package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	go3mf "github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/mesh"
)

type sliceStackDecoder struct {
	emptyDecoder
	r             *Reader
	progressCount int
	resource      go3mf.SliceStackResource
	hasSlice      bool
}

func (d *sliceStackDecoder) Open() error {
	if !d.r.progress.progress(1.0-2.0/float64(d.progressCount+2), StageReadResources) {
		return ErrUserAborted
	}
	d.r.progress.pushLevel(1.0-2.0/float64(d.progressCount+2), 1.0-2.0/float64(d.progressCount+1+2))

	d.resource.SliceStack = new(go3mf.SliceStack)
	return nil
}
func (d *sliceStackDecoder) Close() error {
	if d.resource.ID == 0 {
		return errors.New("go3mf: missing slice stack id attribute")
	}
	if d.resource.UsesSliceRef && d.hasSlice {
		return errors.New("go3mf: slicestack contains slices and slicerefs")
	}
	d.r.addResource(&d.resource)
	d.r.progress.popLevel()
	return nil
}
func (d *sliceStackDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsSliceSpec {
		if name.Local == attrSlice {
			d.hasSlice = true
			child = &sliceDecoder{r: d.r, resource: &d.resource}
		} else if name.Local == attrSliceRef {
			child = &sliceRefDecoder{r: d.r, resource: &d.resource}
		}
	}
	return
}

func (d *sliceStackDecoder) Attributes(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		switch a.Name.Local {
		case attrID:
			if d.resource.ID != 0 {
				err = errors.New("go3mf: duplicated slicestack id attribute")
			} else {
				d.resource.ID, err = strconv.ParseUint(a.Value, 10, 64)
			}
		case attrZBottom:
			var bottomZ float64
			bottomZ, err = strconv.ParseFloat(a.Value, 32)
			d.resource.SliceStack.BottomZ = float32(bottomZ)
		}
		if err != nil {
			return errors.New("go3mf: texture2d attribute not valid")
		}
	}
	return
}

type sliceRefDecoder struct {
	emptyDecoder
	r        *Reader
	resource *go3mf.SliceStackResource
}

func (d *sliceRefDecoder) Attributes(attrs []xml.Attr) (err error) {
	var sliceStackID uint64
	var path string
	for _, a := range attrs {
		switch a.Name.Local {
		case attrSliceRefID:
			sliceStackID, err = strconv.ParseUint(a.Value, 10, 64)
		case attrSlicePath:
			path = a.Value
		}
	}
	if err != nil {
		return errors.New("go3mf: a sliceref has an invalid slicestackid attribute")
	}

	return d.addSliceRef(sliceStackID, path)
}

func (d *sliceRefDecoder) addSliceRef(sliceStackID uint64, path string) error {
	if path == d.resource.ModelPath {
		return errors.New("go3mf: a slicepath is invalid")
	}
	resource, ok := d.r.Model.FindResource(path, sliceStackID)
	if !ok {
		return errors.New("go3mf: a sliceref points to a unexisting resource")
	}
	sliceStackResource, ok := resource.(*go3mf.SliceStackResource)
	if !ok {
		return errors.New("go3mf: a sliceref points to a resource that is not an slicestack")
	}
	sliceStackResource.TimesRefered++
	for _, s := range sliceStackResource.SliceStack.Slices {
		if _, err := d.resource.AddSlice(s); err != nil {
			return err
		}
	}
	d.resource.UsesSliceRef = true
	return nil
}

type sliceDecoder struct {
	emptyDecoder
	r                      *Reader
	resource               *go3mf.SliceStackResource
	slice                  mesh.Slice
	hasTopZ                bool
	polygonDecoder         polygonDecoder
	polygonVerticesDecoder polygonVerticesDecoder
}

func (d *sliceDecoder) Open() error {
	if len(d.resource.Slices)%readSliceUpdate == readSliceUpdate-1 {
		if !d.r.progress.progress(1.0-2.0/float64(len(d.resource.Slices)+2), StageReadSlices) {
			return ErrUserAborted
		}
	}
	d.polygonDecoder.r = d.r
	d.polygonDecoder.slice = &d.slice
	d.polygonVerticesDecoder.r = d.r
	d.polygonVerticesDecoder.slice = &d.slice
	return nil
}
func (d *sliceDecoder) Close() error {
	if !d.hasTopZ {
		return errors.New("go3mf: missing slice topz attribute")
	}
	d.resource.SliceStack.Slices = append(d.resource.SliceStack.Slices, &d.slice)
	return nil
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

func (d *sliceDecoder) Attributes(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Local == attrZTop {
			if d.hasTopZ {
				err = errors.New("go3mf: duplicated slice topz attribute")
			} else {
				d.hasTopZ = true
				var topZ float64
				topZ, err = strconv.ParseFloat(a.Value, 32)
				d.slice.TopZ = float32(topZ)
			}
		}
	}
	return
}

type polygonVerticesDecoder struct {
	emptyDecoder
	r                    *Reader
	slice                *mesh.Slice
	polygonVertexDecoder polygonVertexDecoder
}

func (d *polygonVerticesDecoder) Open() error {
	d.polygonVertexDecoder.r = d.r
	d.polygonVertexDecoder.slice = d.slice
	return nil
}

func (d *polygonVerticesDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsSliceSpec && name.Local == attrVertex {
		child = &d.polygonVertexDecoder
	}
	return
}

type polygonVertexDecoder struct {
	emptyDecoder
	r     *Reader
	slice *mesh.Slice
}

func (d *polygonVertexDecoder) Attributes(attrs []xml.Attr) (err error) {
	var x, y float64
	for _, a := range attrs {
		switch a.Name.Local {
		case attrX:
			x, err = strconv.ParseFloat(a.Value, 32)
		case attrY:
			y, err = strconv.ParseFloat(a.Value, 32)
		}
		if err != nil {
			return errors.New("go3mf: slice vertex has an invalid coordinate attribute")
		}
	}
	d.slice.AddVertex(float32(x), float32(y))
	return
}

type polygonDecoder struct {
	emptyDecoder
	r                     *Reader
	slice                 *mesh.Slice
	polygonIndex          int
	polygonSegmentDecoder polygonSegmentDecoder
}

func (d *polygonDecoder) Open() error {
	d.polygonIndex = d.slice.BeginPolygon()
	d.polygonSegmentDecoder.r = d.r
	d.polygonSegmentDecoder.slice = d.slice
	d.polygonSegmentDecoder.polygonIndex = d.polygonIndex
	return nil
}
func (d *polygonDecoder) Close() error {
	if !d.slice.IsPolygonValid(d.polygonIndex) {
		return errors.New("go3mf: a closed slice polygon is actually a line")
	}
	return nil
}
func (d *polygonDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsSliceSpec && name.Local == attrSegment {
		child = &d.polygonSegmentDecoder
	}
	return
}

func (d *polygonDecoder) Attributes(attrs []xml.Attr) (err error) {
	var start64 uint64
	for _, a := range attrs {
		if a.Name.Local == attrStartV {
			start64, err = strconv.ParseUint(a.Value, 10, 32)
			break
		}
	}
	if err == nil {
		err = d.slice.AddPolygonIndex(d.polygonIndex, int(start64))
	}
	return
}

type polygonSegmentDecoder struct {
	emptyDecoder
	r            *Reader
	slice        *mesh.Slice
	polygonIndex int
}

func (d *polygonSegmentDecoder) Attributes(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Local == attrV2 {
			var v264 uint64
			v264, err = strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				err = errors.New("go3mf: a polygon has an invalid v2 attribute")
			} else {
				d.slice.AddPolygonIndex(d.polygonIndex, int(v264))
			}
			break
		}
	}
	return
}
