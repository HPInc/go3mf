package go3mf

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"sort"
	"strconv"
	"strings"

	xml3mf "github.com/qmuntal/go3mf/internal/xml"
)

const defaultFloatPrecision = 4

// XMLEncoder is based on the encoding/xml.Encoder implementation.
// It is modified to allow custom local namespaces and selfclosing nodes.
type XMLEncoder struct {
	floatPresicion int
	relationships  []Relationship
	p              xml3mf.Printer
}

// newXMLEncoder returns a new encoder that writes to w.
func newXMLEncoder(w io.Writer, floatPresicion int) *XMLEncoder {
	return &XMLEncoder{
		floatPresicion: floatPresicion,
		p:              xml3mf.Printer{Writer: bufio.NewWriter(w)},
	}
}

// AddRelationship adds a relationship to the encoded model.
// Duplicated relationships will be removed before encoding.
func (enc *XMLEncoder) AddRelationship(r Relationship) {
	enc.relationships = append(enc.relationships, r)
}

// FloatPresicion returns the float presicion to use
// when encoding floats.
func (enc *XMLEncoder) FloatPresicion() int {
	return enc.floatPresicion
}

// EncodeToken writes the given XML token to the stream.
func (enc *XMLEncoder) EncodeToken(t xml.Token) {
	p := &enc.p
	switch t := t.(type) {
	case xml.StartElement:
		p.WriteStart(&t)
	case xml.EndElement:
		p.WriteEnd(t.Name)
	case xml.CharData:
		xml.EscapeText(p, t)
	}
}

// Flush flushes any buffered XML to the underlying writer.
func (enc *XMLEncoder) Flush() error {
	return enc.p.Flush()
}

// SetAutoClose define if a start token will be self closed.
// Callers should not end the start token if the encode is in
// auto close mode.
func (enc *XMLEncoder) SetAutoClose(autoClose bool) {
	enc.p.AutoClose = autoClose
}

type packagePart interface {
	io.Writer
	AddRelationship(Relationship)
}

type packageWriter interface {
	Create(name, contentType string) (packagePart, error)
	AddRelationship(Relationship)
	Close() error
}

// Marshaler is the interface implemented by objects
// that can marshal themselves into valid XML elements.
type Marshaler interface {
	Marshal3MF(x *XMLEncoder) error
}

// AttrMarshaler is the interface implemented by objects that can marshal
// themselves into valid XML attributes.
type AttrMarshaler interface {
	Marshal3MFAttr(*XMLEncoder) ([]xml.Attr, error)
}

// MarshalModel returns the XML encoding of m.
func MarshalModel(m *Model) ([]byte, error) {
	var b bytes.Buffer
	if err := new(Encoder).writeModel(newXMLEncoder(&b, defaultFloatPrecision), m); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// An Encoder writes Model data to an output stream.
//
// See the documentation for strconv.FormatFloat for details about the FloatPrecision behaviour.
type Encoder struct {
	FloatPrecision int
	w              packageWriter
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		FloatPrecision: defaultFloatPrecision,
		w:              newOpcWriter(w),
	}
}

// Encode writes the XML encoding of m to the stream.
func (e *Encoder) Encode(m *Model) error {
	if err := e.writeAttachements(m.Attachments); err != nil {
		return err
	}
	rootName := m.PathOrDefault()
	e.w.AddRelationship(Relationship{Type: RelType3DModel, Path: rootName})
	for _, r := range m.RootRelationships {
		e.w.AddRelationship(r)
	}

	w, err := e.w.Create(rootName, ContentType3DModel)
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	enc := newXMLEncoder(w, e.FloatPrecision)
	enc.relationships = make([]Relationship, len(m.Relationships))
	copy(enc.relationships, m.Relationships)
	for path := range m.Childs {
		enc.AddRelationship(Relationship{Type: RelType3DModel, Path: path})
	}
	if err = e.writeModel(enc, m); err != nil {
		return err
	}
	for _, r := range enc.relationships {
		w.AddRelationship(r)
	}
	if err = e.writeChildModels(m); err != nil {
		return err
	}

	return e.w.Close()
}

func (e *Encoder) writeChildModels(m *Model) error {
	for path, child := range m.Childs {
		var (
			w   packagePart
			err error
		)
		if w, err = e.w.Create(path, ContentType3DModel); err != nil {
			return err
		}
		if _, err = w.Write([]byte(xml.Header)); err != nil {
			return err
		}
		enc := newXMLEncoder(w, e.FloatPrecision)
		enc.relationships = child.Relationships
		if err = e.writeChildModel(enc, m, path); err != nil {
			return err
		}
		for _, r := range enc.relationships {
			w.AddRelationship(r)
		}
	}
	return nil
}

func (e *Encoder) writeAttachements(att []Attachment) error {
	for _, a := range att {
		w, err := e.w.Create(a.Path, a.ContentType)
		if err == nil {
			_, err = io.Copy(w, a.Stream)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) modelToken(x *XMLEncoder, m *Model, isRoot bool) (xml.StartElement, error) {
	attrs := []xml.Attr{
		{Name: xml.Name{Local: attrXmlns}, Value: Namespace},
		{Name: xml.Name{Local: attrUnit}, Value: m.Units.String()},
		{Name: xml.Name{Space: nsXML, Local: attrLang}, Value: m.Language},
	}
	if isRoot && m.Thumbnail != "" {
		x.AddRelationship(Relationship{Path: m.Thumbnail, Type: RelTypeThumbnail})
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: attrThumbnail}, Value: m.Thumbnail})
	}
	sortedSpecs := m.sortedSpecs()
	for _, ns := range sortedSpecs {
		a := m.Specs[ns]
		attrs = append(attrs, xml.Attr{Name: xml.Name{Space: attrXmlns, Local: a.Local()}, Value: a.Namespace()})
	}
	var exts []string
	for _, ns := range sortedSpecs {
		a := m.Specs[ns]
		if a.Required() {
			exts = append(exts, a.Local())
		}
	}
	sort.Strings(exts)
	if len(exts) != 0 {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: attrReqExt}, Value: strings.Join(exts, " ")})
	}
	tm := xml.StartElement{Name: xml.Name{Local: attrModel}, Attr: attrs}
	m.AnyAttr.encode(x, &tm)
	return tm, nil
}

func (e *Encoder) writeChildModel(x *XMLEncoder, m *Model, path string) error {
	tm, _ := e.modelToken(x, m, false) // error already checked before
	x.EncodeToken(tm)

	child := m.Childs[path]
	if err := e.writeResources(x, &child.Resources); err != nil {
		return err
	}

	xb := xml.StartElement{Name: xml.Name{Local: attrBuild}}
	x.EncodeToken(xb)
	x.EncodeToken(xb.End())
	child.Any.encode(x)
	x.EncodeToken(tm.End())
	return x.Flush()
}

func (e *Encoder) writeModel(x *XMLEncoder, m *Model) error {
	tm, err := e.modelToken(x, m, true)
	if err != nil {
		return err
	}
	x.EncodeToken(tm)

	e.writeMetadata(x, m.Metadata)
	if err := e.writeResources(x, &m.Resources); err != nil {
		return err
	}
	e.writeBuild(x, m)
	m.Any.encode(x)
	x.EncodeToken(tm.End())
	return x.Flush()
}

func (e *Encoder) writeMetadataGroup(x *XMLEncoder, m []Metadata) {
	xm := xml.StartElement{Name: xml.Name{Local: attrMetadataGroup}}
	x.EncodeToken(xm)
	e.writeMetadata(x, m)
	x.EncodeToken(xm.End())
}

func (e *Encoder) writeBuild(x *XMLEncoder, m *Model) {
	xb := xml.StartElement{Name: xml.Name{Local: attrBuild}}
	m.Build.AnyAttr.encode(x, &xb)
	x.EncodeToken(xb)
	x.SetAutoClose(true)
	for _, item := range m.Build.Items {
		xi := xml.StartElement{Name: xml.Name{Local: attrItem}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrObjectID}, Value: strconv.FormatUint(uint64(item.ObjectID), 10)},
		}}
		if item.HasTransform() {
			xi.Attr = append(xi.Attr, xml.Attr{
				Name: xml.Name{Local: attrTransform}, Value: item.Transform.String(),
			})
		}
		if item.PartNumber != "" {
			xi.Attr = append(xi.Attr, xml.Attr{
				Name: xml.Name{Local: attrPartNumber}, Value: item.PartNumber,
			})
		}
		item.AnyAttr.encode(x, &xi)
		if len(item.Metadata) != 0 {
			x.SetAutoClose(false)
			x.EncodeToken(xi)
			e.writeMetadataGroup(x, item.Metadata)
			x.EncodeToken(xi.End())
			x.SetAutoClose(true)
		} else {
			x.EncodeToken(xi)
		}
	}
	x.SetAutoClose(false)
	x.EncodeToken(xb.End())
}

func (e *Encoder) writeResources(x *XMLEncoder, rs *Resources) error {
	xt := xml.StartElement{Name: xml.Name{Local: attrResources}}
	x.EncodeToken(xt)
	for _, r := range rs.Assets {
		switch r := r.(type) {
		case *BaseMaterials:
			e.writeBaseMaterial(x, r)
		case Marshaler:
			if err := r.Marshal3MF(x); err != nil {
				return err
			}
		}
		if err := x.Flush(); err != nil {
			return err
		}
	}

	for _, o := range rs.Objects {
		e.writeObject(x, o)
		if err := x.Flush(); err != nil {
			return err
		}
	}
	x.EncodeToken(xt.End())
	return nil
}

func (e *Encoder) writeMetadata(x *XMLEncoder, metadata []Metadata) {
	for _, md := range metadata {
		name := md.Name.Local
		if md.Name.Space != "" {
			name = md.Name.Space + ":" + name
		}
		xn := xml.StartElement{Name: xml.Name{Local: attrMetadata}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrName}, Value: name},
		}}
		if md.Preserve {
			xn.Attr = append(xn.Attr, xml.Attr{
				Name: xml.Name{Local: attrPreserve}, Value: strconv.FormatBool(md.Preserve),
			})
		}
		if md.Type != "" {
			xn.Attr = append(xn.Attr, xml.Attr{
				Name: xml.Name{Local: attrType}, Value: md.Type,
			})
		}
		x.EncodeToken(xn)
		x.EncodeToken(xml.CharData(md.Value))
		x.EncodeToken(xn.End())
	}
}

func (e *Encoder) writeObject(x *XMLEncoder, r *Object) {
	xo := xml.StartElement{Name: xml.Name{Local: attrObject}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrID}, Value: strconv.FormatUint(uint64(r.ID), 10)},
	}}
	if r.Type != ObjectTypeModel {
		xo.Attr = append(xo.Attr, xml.Attr{Name: xml.Name{Local: attrType}, Value: r.Type.String()})
	}
	if r.Thumbnail != "" {
		x.AddRelationship(Relationship{Path: r.Thumbnail, Type: RelTypeThumbnail})
		xo.Attr = append(xo.Attr, xml.Attr{Name: xml.Name{Local: attrThumbnail}, Value: r.Thumbnail})
	}
	if r.PartNumber != "" {
		xo.Attr = append(xo.Attr, xml.Attr{Name: xml.Name{Local: attrPartNumber}, Value: r.PartNumber})
	}
	if r.Name != "" {
		xo.Attr = append(xo.Attr, xml.Attr{Name: xml.Name{Local: attrName}, Value: r.Name})
	}
	if r.Mesh != nil {
		if r.PID != 0 {
			xo.Attr = append(xo.Attr, xml.Attr{
				Name: xml.Name{Local: attrPID}, Value: strconv.FormatUint(uint64(r.PID), 10),
			})
		}
		if r.PIndex != 0 {
			xo.Attr = append(xo.Attr, xml.Attr{
				Name: xml.Name{Local: attrPIndex}, Value: strconv.FormatUint(uint64(r.PIndex), 10),
			})
		}
	}
	r.AnyAttr.encode(x, &xo)
	x.EncodeToken(xo)

	if len(r.Metadata) != 0 {
		e.writeMetadataGroup(x, r.Metadata)
	}

	if r.Mesh != nil {
		e.writeMesh(x, r, r.Mesh)
	} else {
		e.writeComponents(x, r.Components)
	}
	x.EncodeToken(xo.End())
}

func (e *Encoder) writeComponents(x *XMLEncoder, comps []*Component) {
	xcs := xml.StartElement{Name: xml.Name{Local: attrComponents}}
	x.EncodeToken(xcs)
	x.SetAutoClose(true)
	for _, c := range comps {
		xt := xml.StartElement{
			Name: xml.Name{Local: attrComponent}, Attr: []xml.Attr{
				{Name: xml.Name{Local: attrObjectID}, Value: strconv.FormatUint(uint64(c.ObjectID), 10)},
			},
		}
		if c.HasTransform() {
			xt.Attr = append(xt.Attr, xml.Attr{Name: xml.Name{Local: attrTransform}, Value: c.Transform.String()})
		}
		c.AnyAttr.encode(x, &xt)
		x.EncodeToken(xt)
	}
	x.SetAutoClose(false)
	x.EncodeToken(xcs.End())
}

func (e *Encoder) writeMesh(x *XMLEncoder, r *Object, m *Mesh) {
	xm := xml.StartElement{Name: xml.Name{Local: attrMesh}}
	m.AnyAttr.encode(x, &xm)
	x.EncodeToken(xm)
	xvs := xml.StartElement{Name: xml.Name{Local: attrVertices}}
	x.EncodeToken(xvs)
	x.SetAutoClose(true)
	for _, v := range m.Vertices {
		x.EncodeToken(xml.StartElement{
			Name: xml.Name{Local: attrVertex},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: attrX}, Value: strconv.FormatFloat(float64(v.X()), 'f', x.FloatPresicion(), 32)},
				{Name: xml.Name{Local: attrY}, Value: strconv.FormatFloat(float64(v.Y()), 'f', x.FloatPresicion(), 32)},
				{Name: xml.Name{Local: attrZ}, Value: strconv.FormatFloat(float64(v.Z()), 'f', x.FloatPresicion(), 32)},
			},
		})
	}
	x.SetAutoClose(false)
	x.EncodeToken(xvs.End())

	xvt := xml.StartElement{Name: xml.Name{Local: attrTriangles}}
	x.EncodeToken(xvt)
	x.SetAutoClose(true)
	for _, v := range m.Triangles {
		v1, v2, v3 := v.Indices()
		t := xml.StartElement{
			Name: xml.Name{Local: attrTriangle},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: attrV1}, Value: strconv.FormatUint(uint64(v1), 10)},
				{Name: xml.Name{Local: attrV2}, Value: strconv.FormatUint(uint64(v2), 10)},
				{Name: xml.Name{Local: attrV3}, Value: strconv.FormatUint(uint64(v3), 10)},
			},
		}
		pid := v.PID()
		if pid != 0 {
			p1, p2, p3 := v.PIndices()
			if (p1 != p2) || (p1 != p3) {
				t.Attr = append(t.Attr,
					xml.Attr{Name: xml.Name{Local: attrPID}, Value: strconv.FormatUint(uint64(pid), 10)},
					xml.Attr{Name: xml.Name{Local: attrP1}, Value: strconv.FormatUint(uint64(p1), 10)},
					xml.Attr{Name: xml.Name{Local: attrP2}, Value: strconv.FormatUint(uint64(p2), 10)},
					xml.Attr{Name: xml.Name{Local: attrP3}, Value: strconv.FormatUint(uint64(p3), 10)},
				)
			} else if (pid != r.PID) || (p1 != r.PIndex) {
				t.Attr = append(t.Attr,
					xml.Attr{Name: xml.Name{Local: attrPID}, Value: strconv.FormatUint(uint64(pid), 10)},
					xml.Attr{Name: xml.Name{Local: attrP1}, Value: strconv.FormatUint(uint64(p1), 10)},
				)
			}
		}
		x.EncodeToken(t)
	}
	x.SetAutoClose(false)
	x.EncodeToken(xvt.End())
	m.Any.encode(x)
	x.EncodeToken(xm.End())
}

func (e *Encoder) writeBaseMaterial(x *XMLEncoder, r *BaseMaterials) {
	xt := xml.StartElement{Name: xml.Name{Local: attrBaseMaterials}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrID}, Value: strconv.FormatUint(uint64(r.ID), 10)},
	}}
	x.EncodeToken(xt)
	x.SetAutoClose(true)
	for _, ma := range r.Materials {
		x.EncodeToken(xml.StartElement{
			Name: xml.Name{Local: attrBase},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: attrName}, Value: ma.Name},
				{Name: xml.Name{Local: attrDisplayColor}, Value: FormatRGBA(ma.Color)},
			},
		})
	}
	x.SetAutoClose(false)
	x.EncodeToken(xt.End())
}
