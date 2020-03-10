package slices

import (
	"encoding/xml"
	"strconv"

	"github.com/qmuntal/go3mf"
)

// Marshal3MFAttr encodes the resource attributes.
func (s *SliceStackInfo) Marshal3MFAttr(_ *go3mf.XMLEncoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: ExtensionName, Local: attrSliceRefID}, Value: strconv.FormatUint(uint64(s.SliceStackID), 10)},
		{Name: xml.Name{Space: ExtensionName, Local: attrMeshRes}, Value: s.SliceResolution.String()},
	}, nil
}

// Marshal3MF encodes the resource.
func (s *SliceStack) Marshal3MF(x *go3mf.XMLEncoder) error {
	xs := xml.StartElement{Name: xml.Name{Space: ExtensionName, Local: attrSliceStack}, Attr: []xml.Attr{
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
		x.EncodeToken(xml.StartElement{Name: xml.Name{Space: ExtensionName, Local: attrSliceRef}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrSliceRefID}, Value: strconv.FormatUint(uint64(r.SliceStackID), 10)},
			{Name: xml.Name{Local: attrSlicePath}, Value: r.Path},
		}})
	}
	x.SetAutoClose(false)
	x.EncodeToken(xs.End())
	return nil
}

func (s *Slice) marshal3MF(x *go3mf.XMLEncoder) {
	xs := xml.StartElement{Name: xml.Name{Space: ExtensionName, Local: attrSlice}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrZTop}, Value: strconv.FormatFloat(float64(s.TopZ), 'f', x.FloatPresicion(), 32)},
	}}
	x.EncodeToken(xs)

	marshalVertices(x, s.Vertices)
	marshalPolygons(x, s.Polygons)

	x.EncodeToken(xs.End())
}

func marshalPolygons(x *go3mf.XMLEncoder, ply []Polygon) {
	for _, p := range ply {
		xp := xml.StartElement{Name: xml.Name{Space: ExtensionName, Local: attrPolygon}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrStartV}, Value: strconv.FormatUint(uint64(p.StartV), 10)},
		}}
		x.EncodeToken(xp)
		x.SetAutoClose(true)
		for _, s := range p.Segments {
			xs := xml.StartElement{Name: xml.Name{Space: ExtensionName, Local: attrSegment}, Attr: []xml.Attr{
				{Name: xml.Name{Local: attrV2}, Value: strconv.FormatUint(uint64(s.V2), 10)},
			}}
			if s.PID != 0 {
				if s.P1 != s.P2 {
					xs.Attr = append(xs.Attr,
						xml.Attr{Name: xml.Name{Local: attrPID}, Value: strconv.FormatUint(uint64(s.PID), 10)},
						xml.Attr{Name: xml.Name{Local: attrP1}, Value: strconv.FormatUint(uint64(s.P1), 10)},
						xml.Attr{Name: xml.Name{Local: attrP2}, Value: strconv.FormatUint(uint64(s.P2), 10)},
					)
				} else {
					xs.Attr = append(xs.Attr,
						xml.Attr{Name: xml.Name{Local: attrPID}, Value: strconv.FormatUint(uint64(s.PID), 10)},
						xml.Attr{Name: xml.Name{Local: attrP1}, Value: strconv.FormatUint(uint64(s.P1), 10)},
					)
				}
			}
			x.EncodeToken(xs)
		}
		x.SetAutoClose(false)
		x.EncodeToken(xp.End())
	}
}

func marshalVertices(x *go3mf.XMLEncoder, vs []go3mf.Point2D) {
	xv := xml.StartElement{Name: xml.Name{Space: ExtensionName, Local: attrVertices}}
	x.EncodeToken(xv)
	x.SetAutoClose(true)
	for _, v := range vs {
		x.EncodeToken(xml.StartElement{Name: xml.Name{Space: ExtensionName, Local: attrVertex}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrX}, Value: strconv.FormatFloat(float64(v.X()), 'f', x.FloatPresicion(), 32)},
			{Name: xml.Name{Local: attrY}, Value: strconv.FormatFloat(float64(v.Y()), 'f', x.FloatPresicion(), 32)},
		}})
	}
	x.SetAutoClose(false)
	x.EncodeToken(xv.End())
}
