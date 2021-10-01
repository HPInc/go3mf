// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package slices

import (
	"encoding/xml"
	"strconv"

	"github.com/hpinc/go3mf"
	"github.com/hpinc/go3mf/spec"
)

// Marshal3MF encodes the resource attributes.
func (s *ObjectAttr) Marshal3MF(_ spec.Encoder, start *xml.StartElement) error {
	start.Attr = append(start.Attr,
		xml.Attr{Name: xml.Name{Space: Namespace, Local: attrSliceRefID}, Value: strconv.FormatUint(uint64(s.SliceStackID), 10)},
		xml.Attr{Name: xml.Name{Space: Namespace, Local: attrMeshRes}, Value: s.MeshResolution.String()},
	)
	return nil
}

// Marshal3MF encodes the resource.
func (s *SliceStack) Marshal3MF(x spec.Encoder, _ *xml.StartElement) error {
	xs := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrSliceStack}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrID}, Value: strconv.FormatUint(uint64(s.ID), 10)},
	}}
	if s.BottomZ != 0 {
		xs.Attr = append(xs.Attr, xml.Attr{
			Name:  xml.Name{Local: attrZBottom},
			Value: strconv.FormatFloat(float64(s.BottomZ), 'f', x.FloatPresicion(), 32),
		})
	}
	x.EncodeToken(xs)
	for _, sl := range s.Slices {
		sl.marshal3MF(x)
	}
	x.SetAutoClose(true)
	for _, r := range s.Refs {
		x.EncodeToken(xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrSliceRef}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrSliceRefID}, Value: strconv.FormatUint(uint64(r.SliceStackID), 10)},
			{Name: xml.Name{Local: attrSlicePath}, Value: r.Path},
		}})
	}
	x.SetAutoClose(false)
	x.EncodeToken(xs.End())
	return nil
}

func (s *Slice) marshal3MF(x spec.Encoder) {
	xs := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrSlice}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrZTop}, Value: strconv.FormatFloat(float64(s.TopZ), 'f', x.FloatPresicion(), 32)},
	}}
	x.EncodeToken(xs)

	marshalVertices(x, s.Vertices)
	marshalPolygons(x, s.Polygons)

	x.EncodeToken(xs.End())
}

func marshalPolygons(x spec.Encoder, ply []Polygon) {
	xs := xml.StartElement{
		Name: xml.Name{Space: Namespace, Local: attrSegment},
	}
	xsattrs := []xml.Attr{
		{Name: xml.Name{Local: attrV2}},
		{Name: xml.Name{Local: attrPID}},
		{Name: xml.Name{Local: attrP1}},
		{Name: xml.Name{Local: attrP2}},
	}
	x.SetSkipAttrEscape(true)
	for _, p := range ply {
		xp := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrPolygon}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrStartV}, Value: strconv.FormatUint(uint64(p.StartV), 10)},
		}}
		x.EncodeToken(xp)
		x.SetAutoClose(true)
		for _, s := range p.Segments {
			xsattrs[0].Value = strconv.FormatUint(uint64(s.V2), 10)
			xs.Attr = xsattrs[:1]
			if s.PID != 0 {
				xsattrs[1].Value = strconv.FormatUint(uint64(s.PID), 10)
				xsattrs[2].Value = strconv.FormatUint(uint64(s.P1), 10)
				if s.P1 != s.P2 {
					xsattrs[3].Value = strconv.FormatUint(uint64(s.P2), 10)
					xs.Attr = xsattrs[:4]
				} else {
					xs.Attr = xsattrs[:3]
				}
			}
			x.EncodeToken(xs)
		}
		x.SetAutoClose(false)
		x.EncodeToken(xp.End())
	}
	x.SetSkipAttrEscape(false)
}

func marshalVertices(x spec.Encoder, vs []go3mf.Point2D) {
	xv := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrVertices}}
	x.EncodeToken(xv)
	x.SetAutoClose(true)
	x.SetSkipAttrEscape(true)
	prec := x.FloatPresicion()
	start := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrVertex}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrX}},
		{Name: xml.Name{Local: attrY}},
	}}
	for _, v := range vs {
		start.Attr[0].Value = strconv.FormatFloat(float64(v.X()), 'f', prec, 32)
		start.Attr[1].Value = strconv.FormatFloat(float64(v.Y()), 'f', prec, 32)
		x.EncodeToken(start)
	}
	x.SetSkipAttrEscape(false)
	x.SetAutoClose(false)
	x.EncodeToken(xv.End())
}
