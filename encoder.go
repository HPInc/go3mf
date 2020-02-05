package go3mf

import (
	"context"
	"encoding/xml"
	"io"
	"strconv"

	"github.com/qmuntal/opc"
)

type tokenEncoder interface {
	EncodeToken(t xml.Token)
	Flush() error
}

type packageWriter interface {
	Create(name, contentType string) (io.Writer, error)
	AddRelationship(*relationship)
	Close() error
}

type Encoder struct {
	w packageWriter
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: &opcWriter{opc.NewWriter(w)},
	}
}

func (e *Encoder) Encode(ctx context.Context, m *Model) error {
	rootName := m.Path
	if rootName == "" {
		rootName = uriDefault3DModel
	}
	w, err := e.w.Create(rootName, contentType3DModel)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	x := newXmlEncoder(w)
	if err = e.writeModel(x, m); err != nil {
		return err
	}
	e.w.AddRelationship(&relationship{
		ID: "1", Type: RelTypeModel3D, TargetURI: rootName,
	})
	return e.w.Close()
}

func (e *Encoder) writeModel(x tokenEncoder, m *Model) error {
	attrs := []xml.Attr{
		{Name: xml.Name{Local: attrXmlns}, Value: ExtensionName},
		{Name: xml.Name{Local: attrUnit}, Value: m.Units.String()},
		{Name: xml.Name{Space: nsXML, Local: attrLang}, Value: m.Language},
	}
	if m.Thumbnail != "" {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: attrThumbnail}, Value: m.Thumbnail})
	}
	for _, a := range m.Namespaces {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Space: attrXmlns, Local: a.Local}, Value: a.Space})
	}

	tm := xml.StartElement{Name: xml.Name{Local: attrModel}, Attr: attrs}
	x.EncodeToken(tm)

	e.writeMetadata(x, m.Metadata)

	if err := e.writeResources(x, m); err != nil {
		return err
	}

	e.writeBuild(x, m)
	x.EncodeToken(tm.End())
	return x.Flush()
}

func (e *Encoder) writeMetadataGroup(x tokenEncoder, m []Metadata) {
	xm := newXmlNodeEncoder(x, attrMetadataGroup, 0)
	xm.Close()
	e.writeMetadata(x, m)
	xm.End()
}

func (e *Encoder) writeBuild(x tokenEncoder, m *Model) {
	xb := newXmlNodeEncoder(x, attrBuild, 0)
	xb.Close()

	for _, i := range m.Build.Items {
		xi := newXmlNodeEncoder(x, attrItem, 3)
		xi.Attribute(attrObjectID, strconv.FormatUint(uint64(i.ObjectID), 10))
		if i.HasTransform() {
			xi.Attribute(attrTransform, FormatMatrix(i.Transform))
		}
		xi.OptionalAttribute(attrPartNumber, i.PartNumber)
		if len(i.Metadata) != 0 {
			xi.Close()
			e.writeMetadataGroup(x, i.Metadata)
		}

		xi.End()
	}

	xb.End()
}

func (e *Encoder) writeResources(x tokenEncoder, m *Model) error {
	xt := newXmlNodeEncoder(x, attrResources, 0)
	xt.Close()
	for _, r := range m.Resources {
		switch r := r.(type) {
		case *BaseMaterialsResource:
			e.writeBaseMaterial(x, r)
		case *ObjectResource:
			e.writeObject(x, r)
		}

		if err := x.Flush(); err != nil {
			return err
		}
	}
	xt.End()
	return nil
}

func (e *Encoder) writeMetadata(x tokenEncoder, metadata []Metadata) {
	for _, md := range metadata {
		xn := newXmlNodeEncoder(x, attrMetadata, 3)
		xn.Attribute(attrName, md.Name)
		if md.Preserve {
			xn.Attribute(attrPreserve, strconv.FormatBool(md.Preserve))
		}
		xn.OptionalAttribute(attrType, md.Type)
		xn.TextEnd(md.Value)
	}
}

func (e *Encoder) writeObject(x tokenEncoder, r *ObjectResource) {
	xo := newXmlNodeEncoder(x, attrObject, 7)
	xo.Attribute(attrID, strconv.FormatUint(uint64(r.ID), 10))
	xo.OptionalAttribute(attrType, r.ObjectType.String())
	xo.OptionalAttribute(attrThumbnail, r.Thumbnail)
	xo.OptionalAttribute(attrPartNumber, r.PartNumber)
	xo.OptionalAttribute(attrName, r.Name)
	if r.Mesh != nil {
		if r.DefaultPropertyID != 0 {
			xo.Attribute(attrPID, strconv.FormatUint(uint64(r.DefaultPropertyID), 10))
		}
		if r.DefaultPropertyIndex != 0 {
			xo.Attribute(attrPIndex, strconv.FormatUint(uint64(r.DefaultPropertyIndex), 10))
		}
	}
	xo.Close()

	if len(r.Metadata) != 0 {
		e.writeMetadataGroup(x, r.Metadata)
	}

	if r.Mesh != nil {
		e.writeMesh(x, r.Mesh)
	} else {
		e.writeComponents(x, r.Components)
	}
	xo.End()
}

func (e *Encoder) writeComponents(x tokenEncoder, comps []*Component) {
	xcs := newXmlNodeEncoder(x, attrComponents, 0)
	xcs.Close()

	for _, c := range comps {
		xc := newXmlNodeEncoder(x, attrComponent, 2)
		xc.Attribute(attrObjectID, strconv.FormatUint(uint64(c.ObjectID), 10))
		if c.HasTransform() {
			xc.Attribute(attrTransform, FormatMatrix(c.Transform))
		}
		xc.End()
	}
	xcs.End()
}

func (e *Encoder) writeMesh(x tokenEncoder, m *Mesh) {
	xm := newXmlNodeEncoder(x, attrMesh, 0)
	xm.Close()
	xvs := newXmlNodeEncoder(x, attrVertices, 0)
	xvs.Close()
	for _, v := range m.Nodes {
		xv := newXmlNodeEncoder(x, attrVertex, 3)
		xv.Attribute(attrX, strconv.FormatFloat(float64(v.X()), 'f', 3, 32))
		xv.Attribute(attrY, strconv.FormatFloat(float64(v.Y()), 'f', 3, 32))
		xv.Attribute(attrZ, strconv.FormatFloat(float64(v.Z()), 'f', 3, 32))
		xv.End()
	}
	xvs.End()
	xvt := newXmlNodeEncoder(x, attrTriangles, 0)
	xvt.Close()
	for _, v := range m.Faces {
		xv := newXmlNodeEncoder(x, attrTriangle, 3)
		xv.Attribute(attrV1, strconv.FormatUint(uint64(v.NodeIndices[0]), 10))
		xv.Attribute(attrV2, strconv.FormatUint(uint64(v.NodeIndices[1]), 10))
		xv.Attribute(attrV3, strconv.FormatUint(uint64(v.NodeIndices[2]), 10))
		if v.Resource != 0 {
			xv.Attribute(attrPID, strconv.FormatUint(uint64(v.Resource), 10))
			if v.ResourceIndices[0] != 0 {
				xv.Attribute(attrP1, strconv.FormatUint(uint64(v.ResourceIndices[0]), 10))
				if v.ResourceIndices[1] != 0 {
					xv.Attribute(attrP2, strconv.FormatUint(uint64(v.ResourceIndices[1]), 10))
				}
				if v.ResourceIndices[2] != 0 {
					xv.Attribute(attrP3, strconv.FormatUint(uint64(v.ResourceIndices[2]), 10))
				}
			}
		}
		xv.End()
	}
	xvt.End()
	xm.End()
}

func (e *Encoder) writeBaseMaterial(x tokenEncoder, r *BaseMaterialsResource) {
	xt := newXmlNodeEncoder(x, attrBaseMaterials, 1)
	xt.Attribute(attrID, strconv.FormatUint(uint64(r.ID), 10))
	xt.Close()
	for _, ma := range r.Materials {
		xn := newXmlNodeEncoder(x, attrBase, 2)
		xn.Attribute(attrName, ma.Name)
		xn.Attribute(attrDisplayColor, ma.ColorString())
		xn.End()
	}
	xt.End()
}

type xmlNodeEncoder struct {
	x      tokenEncoder
	start  xml.StartElement
	closed bool
}

func newXmlNodeEncoder(x tokenEncoder, name string, cap int) *xmlNodeEncoder {
	return &xmlNodeEncoder{
		x: x,
		start: xml.StartElement{
			Name: xml.Name{Local: name},
			Attr: make([]xml.Attr, 0, cap),
		},
	}
}

func (e *xmlNodeEncoder) Attribute(name string, value string) {
	e.start.Attr = append(e.start.Attr, xml.Attr{
		Name:  xml.Name{Local: name},
		Value: value,
	})
}

func (e *xmlNodeEncoder) OptionalAttribute(name string, value string) {
	if value != "" {
		e.Attribute(name, value)
	}
}

func (e *xmlNodeEncoder) TextEnd(txt string) {
	e.Close()
	e.x.EncodeToken(xml.CharData(txt))
	e.x.EncodeToken(e.start.End())
}

func (e *xmlNodeEncoder) Close() {
	if !e.closed {
		e.closed = true
		e.x.EncodeToken(e.start)
	}
}

func (e *xmlNodeEncoder) End() {
	e.Close()
	e.x.EncodeToken(e.start.End())
}
