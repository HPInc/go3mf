package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	go3mf "github.com/qmuntal/go3mf"
)

type sliceStackDecoder struct {
	r             *Reader
	sliceStack    go3mf.SliceStackResource
	progressCount uint64
}

func (d *sliceStackDecoder) Decode(x xml.TokenReader, se xml.StartElement) error {
	d.sliceStack.SliceStack = new(go3mf.SliceStack)
	if err := d.parseAttr(se.Attr); err != nil {
		return err
	}
	if d.sliceStack.ID == 0 {
		return errors.New("go3mf: missing slice stack id attribute")
	}
	if err := d.parseContent(x); err != nil {
		return err
	}
	d.r.addResource(&d.sliceStack)
	return nil
}

func (d *sliceStackDecoder) parseAttr(attrs []xml.Attr) error {
	for _, a := range attrs {
		var err error
		switch a.Name.Local {
		case attrID:
			if d.sliceStack.ID != 0 {
				err = errors.New("go3mf: duplicated slicestack id attribute")
			} else {
				d.sliceStack.ID, err = strconv.ParseUint(a.Value, 10, 64)
			}
		case attrZBottom:
			var bottomZ float64
			bottomZ, err = strconv.ParseFloat(a.Value, 32)
			d.sliceStack.SliceStack.BottomZ = float32(bottomZ)
		}
		if err != nil {
			return errors.New("go3mf: texture2d attribute not valid")
		}
	}
	return nil
}

func (d *sliceStackDecoder) parseContent(x xml.TokenReader) error {
	var hasSliceRef, hasSlice bool
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space != nsSliceSpec {
				continue
			}
			if tp.Name.Local == attrSlice {
				if hasSliceRef {
					err = errors.New("go3mf: slicestack contains slices and slicerefs")
				} else {
					hasSlice = true
					err = d.parseSlice(x, tp)
				}
			} else if tp.Name.Local == attrSliceRef {
				if hasSlice {
					err = errors.New("go3mf: slicestack contains slices and slicerefs")
				} else {
					hasSliceRef = true
					err = d.parseSliceRef(tp)
				}
			}
		case xml.EndElement:
			if tp.Name.Space == nsSliceSpec && tp.Name.Local == attrSliceStack {
				return nil
			}
		}
		if err != nil {
			return err
		}
	}
}

func (d *sliceStackDecoder) parseSliceRef(se xml.StartElement) error {
	var sliceStackID uint64
	var path string
	var err error
	for _, a := range se.Attr {
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

func (d *sliceStackDecoder) addSliceRef(sliceStackID uint64, path string) error {
	if path == d.sliceStack.ModelPath {
		return errors.New("go3mf: a slicepath is invalid")
	}
	resource, ok := d.r.Model.FindResource(sliceStackID, path)
	if !ok {
		return errors.New("go3mf: a sliceref points to a unexisting resource")
	}
	sliceStackResource, ok := resource.(*go3mf.SliceStackResource)
	if !ok {
		return errors.New("go3mf: a sliceref points to a resource that is not an slicestack")
	}
	sliceStackResource.TimesRefered++
	for _, s := range sliceStackResource.Slices {
		if _, err := d.sliceStack.AddSlice(s); err != nil {
			return err
		}
	}
	d.sliceStack.UsesSliceRef = true
	return nil
}

func (d *sliceStackDecoder) parseSlice(x xml.TokenReader, se xml.StartElement) (err error) {
	if len(d.sliceStack.Slices)%readSliceUpdate == readSliceUpdate-1 {
		d.progressCount++
		if !d.r.progress.progress(1.0-2.0/float64(d.progressCount+2), StageReadSlices) {
			return ErrUserAborted
		}
	}
	sd := sliceDecoder{r: d.r, sliceStack: d.sliceStack.SliceStack}
	return sd.Decode(x, se)
}

type sliceDecoder struct {
	r          *Reader
	sliceStack *go3mf.SliceStack
	slice      go3mf.Slice
	hasTopZ    bool
}

func (d *sliceDecoder) Decode(x xml.TokenReader, se xml.StartElement) error {
	if err := d.parseAttr(se.Attr); err != nil {
		return err
	}
	if !d.hasTopZ {
		return errors.New("go3mf: missing slice topz attribute")
	}
	if err := d.parseContent(x); err != nil {
		return err
	}
	d.sliceStack.Slices = append(d.sliceStack.Slices, &d.slice)
	return nil
}

func (d *sliceDecoder) parseAttr(attrs []xml.Attr) (err error) {
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

func (d *sliceDecoder) parseContent(x xml.TokenReader) error {
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space != nsSliceSpec {
				continue
			}
			if tp.Name.Local == attrVertices {
				err = d.parseVertices(x, tp)
			} else if tp.Name.Local == attrPolygon {
				err = d.parsePolygons(x, tp)
			}
		case xml.EndElement:
			if tp.Name.Space == nsSliceSpec && tp.Name.Local == attrSlice {
				return nil
			}
		}
		if err != nil {
			return err
		}
	}
}

func (d *sliceDecoder) parseVertices(x xml.TokenReader, se xml.StartElement) error {
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsSliceSpec && tp.Name.Local == attrVertex {
				if err = d.parseVertex(tp.Attr); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if tp.Name.Space == nsSliceSpec && tp.Name.Local == attrVertices {
				return nil
			}
		}
		if err != nil {
			return err
		}
	}
}

func (d *sliceDecoder) parseVertex(attrs []xml.Attr) error {
	var x, y float64
	for _, a := range attrs {
		var err error
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
	return nil
}

func (d *sliceDecoder) parsePolygons(x xml.TokenReader, se xml.StartElement) error {
	polygonIndex := d.slice.BeginPolygon()
	if err := d.parsePolygonAttr(polygonIndex, se.Attr); err != nil {
		return err
	}
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsSliceSpec && tp.Name.Local == attrSegment {
				if err = d.addSegment(polygonIndex, tp.Attr); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if tp.Name.Space == nsSliceSpec && tp.Name.Local == attrPolygon {
				if !d.slice.IsPolygonValid(polygonIndex) {
					return errors.New("go3mf: a closed slice polygon is actually a line")
				}
				return nil
			}
		}
		if err != nil {
			return err
		}
	}
}

func (d *sliceDecoder) addSegment(polygonIndex int, attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Local == attrV2 {
			var v264 uint64
			v264, err = strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				err = errors.New("go3mf: a polygon has an invalid v2 attribute")
			} else {
				d.slice.AddPolygonIndex(polygonIndex, int(v264))
			}
			break
		}
	}
	return
}

func (d *sliceDecoder) parsePolygonAttr(polygonIndex int, attrs []xml.Attr) (err error) {
	var start64 uint64
	for _, a := range attrs {
		if a.Name.Local == attrStartV {
			start64, err = strconv.ParseUint(a.Value, 10, 32)
			break
		}
	}
	if err == nil {
		err = d.slice.AddPolygonIndex(polygonIndex, int(start64))
	}
	return
}
