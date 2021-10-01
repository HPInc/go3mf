// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package materials

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/hpinc/go3mf/spec"
)

// Marshal3MF encodes the resource.
func (r *ColorGroup) Marshal3MF(x spec.Encoder, _ *xml.StartElement) error {
	xs := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrColorGroup}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrID}, Value: strconv.FormatUint(uint64(r.ID), 10)},
	}}
	x.EncodeToken(xs)
	x.SetAutoClose(true)
	x.SetSkipAttrEscape(true)
	start := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrColor}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrColor}},
	}}
	for _, c := range r.Colors {
		start.Attr[0].Value = spec.FormatRGBA(c)
		x.EncodeToken(start)
	}
	x.SetSkipAttrEscape(false)
	x.SetAutoClose(false)
	x.EncodeToken(xs.End())
	return nil
}

// Marshal3MF encodes the resource.
func (r *Texture2DGroup) Marshal3MF(x spec.Encoder, _ *xml.StartElement) error {
	xs := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrTexture2DGroup}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrID}, Value: strconv.FormatUint(uint64(r.ID), 10)},
		{Name: xml.Name{Local: attrTexID}, Value: strconv.FormatUint(uint64(r.TextureID), 10)},
	}}
	x.EncodeToken(xs)
	x.SetAutoClose(true)
	x.SetSkipAttrEscape(true)
	prec := x.FloatPresicion()
	start := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrTex2DCoord}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrU}},
		{Name: xml.Name{Local: attrV}},
	}}
	for _, c := range r.Coords {
		start.Attr[0].Value = strconv.FormatFloat(float64(c.U()), 'f', prec, 32)
		start.Attr[1].Value = strconv.FormatFloat(float64(c.V()), 'f', prec, 32)
		x.EncodeToken(start)
	}
	x.SetSkipAttrEscape(false)
	x.SetAutoClose(false)
	x.EncodeToken(xs.End())
	return nil
}

// Marshal3MF encodes the resource.
func (r *CompositeMaterials) Marshal3MF(x spec.Encoder, _ *xml.StartElement) error {
	indices := make([]string, len(r.Indices))
	for i, idx := range r.Indices {
		indices[i] = strconv.FormatUint(uint64(idx), 10)
	}
	xs := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrCompositematerials}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrID}, Value: strconv.FormatUint(uint64(r.ID), 10)},
		{Name: xml.Name{Local: attrMatID}, Value: strconv.FormatUint(uint64(r.MaterialID), 10)},
		{Name: xml.Name{Local: attrMatIndices}, Value: strings.Join(indices, " ")},
	}}
	x.EncodeToken(xs)
	x.SetAutoClose(true)
	x.SetSkipAttrEscape(true)
	for _, c := range r.Composites {
		values := make([]string, len(c.Values))
		for i, v := range c.Values {
			values[i] = strconv.FormatFloat(float64(v), 'f', x.FloatPresicion(), 32)
		}
		x.EncodeToken(xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrComposite}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrValues}, Value: strings.Join(values, " ")},
		}})
	}
	x.SetSkipAttrEscape(false)
	x.SetAutoClose(false)
	x.EncodeToken(xs.End())
	return nil
}

// Marshal3MF encodes the resource.
func (r *MultiProperties) Marshal3MF(x spec.Encoder, _ *xml.StartElement) error {
	pids := make([]string, len(r.PIDs))
	for i, idx := range r.PIDs {
		pids[i] = strconv.FormatUint(uint64(idx), 10)
	}
	methods := make([]string, len(r.BlendMethods))
	for i, method := range r.BlendMethods {
		methods[i] = method.String()
	}
	xs := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrMultiProps}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrID}, Value: strconv.FormatUint(uint64(r.ID), 10)},
		{Name: xml.Name{Local: attrPIDs}, Value: strings.Join(pids, " ")},
		{Name: xml.Name{Local: attrBlendMethods}, Value: strings.Join(methods, " ")},
	}}
	x.EncodeToken(xs)
	x.SetAutoClose(true)
	x.SetSkipAttrEscape(true)
	for _, mu := range r.Multis {
		indices := make([]string, len(mu.PIndices))
		for i, v := range mu.PIndices {
			indices[i] = strconv.FormatUint(uint64(v), 10)
		}
		x.EncodeToken(xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrMulti}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrPIndices}, Value: strings.Join(indices, " ")},
		}})
	}
	x.SetSkipAttrEscape(false)
	x.SetAutoClose(false)
	x.EncodeToken(xs.End())
	return nil
}

// Marshal3MF encodes the resource.
func (r *Texture2D) Marshal3MF(x spec.Encoder, _ *xml.StartElement) error {
	x.AddRelationship(spec.Relationship{Path: r.Path, Type: RelTypeTexture3D})
	xs := xml.StartElement{Name: xml.Name{Space: Namespace, Local: attrTexture2D}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrID}, Value: strconv.FormatUint(uint64(r.ID), 10)},
		{Name: xml.Name{Local: attrPath}, Value: r.Path},
		{Name: xml.Name{Local: attrContentType}, Value: r.ContentType.String()},
	}}
	if r.TileStyleU != TileWrap {
		xs.Attr = append(xs.Attr, xml.Attr{Name: xml.Name{Local: attrTileStyleU}, Value: r.TileStyleU.String()})
	}
	if r.TileStyleV != TileWrap {
		xs.Attr = append(xs.Attr, xml.Attr{Name: xml.Name{Local: attrTileStyleV}, Value: r.TileStyleV.String()})
	}
	if r.Filter != TextureFilterAuto {
		xs.Attr = append(xs.Attr, xml.Attr{Name: xml.Name{Local: attrFilter}, Value: r.Filter.String()})
	}
	x.SetAutoClose(true)
	x.EncodeToken(xs)
	x.SetAutoClose(false)
	return nil
}
