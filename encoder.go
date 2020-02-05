package go3mf

import (
	"context"
	"encoding/xml"
	"io"
	"strconv"

	"github.com/qmuntal/opc"
)

type tokenEncoder interface {
	EncodeToken(t xml.Token) error
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
	if err = x.Flush(); err != nil {
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
	err := x.EncodeToken(tm)
	if err != nil {
		return err
	}

	if err = e.writeMetadata(x, m.Metadata); err != nil {
		return err
	}

	if err = e.writeResources(x, m); err != nil {
		return err
	}

	if err = e.writeBuild(x, m); err != nil {
		return err
	}

	return x.EncodeToken(tm.End())
}

func (e *Encoder) writeMetadataGroup(x tokenEncoder, m []Metadata) error {
	xm := newXmlNodeEncoder(x, attrMetadataGroup, 0)
	if err := xm.Close(); err != nil {
		return err
	}
	if err := e.writeMetadata(x, m); err != nil {
		return err
	}
	return xm.End()
}

func (e *Encoder) writeBuild(x tokenEncoder, m *Model) error {
	xb := newXmlNodeEncoder(x, attrBuild, 0)
	if err := xb.Close(); err != nil {
		return err
	}

	for _, i := range m.Build.Items {
		xi := newXmlNodeEncoder(x, attrItem, 3)
		xi.Attribute(attrObjectID, strconv.FormatUint(uint64(i.ObjectID), 10))
		if i.HasTransform() {
			xi.Attribute(attrTransform, FormatMatrix(i.Transform))
		}
		xi.OptionalAttribute(attrPartNumber, i.PartNumber)
		if len(i.Metadata) != 0 {
			if err := xi.Close(); err != nil {
				return err
			}
			if err := e.writeMetadataGroup(x, i.Metadata); err != nil {
				return err
			}
		}

		xi.End()
	}

	return xb.End()
}

func (e *Encoder) writeResources(x tokenEncoder, m *Model) error {
	xt := newXmlNodeEncoder(x, attrResources, 0)
	if err := xt.Close(); err != nil {
		return err
	}
	for _, r := range m.Resources {
		var err error
		switch r := r.(type) {
		case *BaseMaterialsResource:
			err = e.writeBaseMaterial(x, r)
		case *ObjectResource:
			err = e.writeObject(x, r)
		}
		if err != nil {
			return err
		}
	}
	return xt.End()
}

func (e *Encoder) writeMetadata(x tokenEncoder, metadata []Metadata) error {
	for _, md := range metadata {
		xn := newXmlNodeEncoder(x, attrMetadata, 3)
		xn.Attribute(attrName, md.Name)
		if md.Preserve {
			xn.Attribute(attrPreserve, strconv.FormatBool(md.Preserve))
		}
		xn.OptionalAttribute(attrType, md.Type)
		if err := xn.TextEnd(md.Value); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) writeObject(x tokenEncoder, r *ObjectResource) error {
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
	if err := xo.Close(); err != nil {
		return err
	}

	if len(r.Metadata) != 0 {
		if err := e.writeMetadataGroup(x, r.Metadata); err != nil {
			return err
		}
	}

	var err error
	if r.Mesh != nil {
		err = e.writeMesh(x, r.Mesh)
	} else {
		err = e.writeComponents(x, r.Components)
	}

	if err != nil {
		return err
	}
	return xo.End()
}

func (e *Encoder) writeComponents(x tokenEncoder, comps []*Component) error {
	xcs := newXmlNodeEncoder(x, attrComponents, 0)
	if err := xcs.Close(); err != nil {
		return err
	}

	for _, c := range comps {
		xc := newXmlNodeEncoder(x, attrComponent, 2)
		xc.Attribute(attrObjectID, strconv.FormatUint(uint64(c.ObjectID), 10))
		if c.HasTransform() {
			xc.Attribute(attrTransform, FormatMatrix(c.Transform))
		}
		xc.End()
	}

	return xcs.End()
}

func (e *Encoder) writeMesh(x tokenEncoder, m *Mesh) error {
	xm := newXmlNodeEncoder(x, attrMesh, 0)
	if err := xm.Close(); err != nil {
		return err
	}
	xvs := newXmlNodeEncoder(x, attrVertices, 0)
	if err := xvs.Close(); err != nil {
		return err
	}
	for _, v := range m.Nodes {
		xv := newXmlNodeEncoder(x, attrVertex, 3)
		xv.Attribute(attrX, strconv.FormatFloat(float64(v.X()), 'f', 3, 32))
		xv.Attribute(attrY, strconv.FormatFloat(float64(v.Y()), 'f', 3, 32))
		xv.Attribute(attrZ, strconv.FormatFloat(float64(v.Z()), 'f', 3, 32))
		if err := xv.End(); err != nil {
			return err
		}
	}
	if err := xvs.End(); err != nil {
		return err
	}
	xvt := newXmlNodeEncoder(x, attrTriangles, 0)
	if err := xvt.Close(); err != nil {
		return err
	}
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
		if err := xv.End(); err != nil {
			return err
		}
	}
	if err := xvt.End(); err != nil {
		return err
	}
	return xm.End()
}

func (e *Encoder) writeBaseMaterial(x tokenEncoder, r *BaseMaterialsResource) error {
	xt := newXmlNodeEncoder(x, attrBaseMaterials, 1)
	xt.Attribute(attrID, strconv.FormatUint(uint64(r.ID), 10))
	if err := xt.Close(); err != nil {
		return err
	}
	for _, ma := range r.Materials {
		xn := newXmlNodeEncoder(x, attrBase, 2)
		xn.Attribute(attrName, ma.Name)
		xn.Attribute(attrDisplayColor, ma.ColorString())
		if err := xn.End(); err != nil {
			return err
		}
	}
	return xt.End()
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

func (e *xmlNodeEncoder) TextEnd(txt string) error {
	if err := e.Close(); err != nil {
		return err
	}
	if err := e.x.EncodeToken(xml.CharData(txt)); err != nil {
		return err
	}
	return e.x.EncodeToken(e.start.End())
}

func (e *xmlNodeEncoder) Close() error {
	if !e.closed {
		e.closed = true
		return e.x.EncodeToken(e.start)
	}
	return nil
}

func (e *xmlNodeEncoder) End() error {
	if err := e.Close(); err != nil {
		return err
	}
	return e.x.EncodeToken(e.start.End())
}
