package go3mf

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"

	"github.com/qmuntal/opc"
)

type tokenEncoder interface {
	EncodeToken(t xml.Token)
	Flush() error
	SetAutoClose(bool)
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
	if err = e.writeModel(newXmlEncoder(w), m); err != nil {
		return err
	}
	e.w.AddRelationship(&relationship{
		ID: "rel0", Type: RelTypeModel3D, TargetURI: rootName,
	})

	if err = e.writeAttachements(m); err != nil {
		return err
	}

	return e.w.Close()
}

func (e *Encoder) writeAttachements(m *Model) error {
	for i, a := range m.Attachments {
		w, err := e.w.Create(a.Path, a.ContentType)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, a.Stream)
		if err != nil {
			return err
		}
		e.w.AddRelationship(&relationship{
			ID: fmt.Sprintf("rel%d", i+1), Type: a.RelationshipType, TargetURI: a.Path,
		})
	}
	return nil
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
	xm := xml.StartElement{Name: xml.Name{Local: attrMetadataGroup}}
	x.EncodeToken(xm)
	e.writeMetadata(x, m)
	x.EncodeToken(xm.End())
}

func (e *Encoder) writeBuild(x tokenEncoder, m *Model) {
	xb := xml.StartElement{Name: xml.Name{Local: attrBuild}}
	x.EncodeToken(xb)
	x.SetAutoClose(true)
	for _, item := range m.Build.Items {
		xi := xml.StartElement{Name: xml.Name{Local: attrItem}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrObjectID}, Value: strconv.FormatUint(uint64(item.ObjectID), 10)},
		}}
		if item.HasTransform() {
			xi.Attr = append(xi.Attr, xml.Attr{
				Name: xml.Name{Local: attrTransform}, Value: FormatMatrix(item.Transform),
			})
		}
		if item.PartNumber != "" {
			xi.Attr = append(xi.Attr, xml.Attr{
				Name: xml.Name{Local: attrPartNumber}, Value: item.PartNumber,
			})
		}
		if len(item.Metadata) != 0 {
			x.SetAutoClose(false)
			x.EncodeToken(xi)
			e.writeMetadataGroup(x, item.Metadata)
			x.EncodeToken(xi.End())
			x.SetAutoClose(true)
		}
	}
	x.SetAutoClose(false)
	x.EncodeToken(xb.End())
}

func (e *Encoder) writeResources(x tokenEncoder, m *Model) error {
	xt := xml.StartElement{Name: xml.Name{Local: attrResources}}
	x.EncodeToken(xt)
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
	x.EncodeToken(xt.End())
	return nil
}

func (e *Encoder) writeMetadata(x tokenEncoder, metadata []Metadata) {
	for _, md := range metadata {
		xn := xml.StartElement{Name: xml.Name{Local: attrMetadata}, Attr: []xml.Attr{
			{Name: xml.Name{Local: attrName}, Value: md.Name},
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

func (e *Encoder) writeObject(x tokenEncoder, r *ObjectResource) {
	xo := xml.StartElement{Name: xml.Name{Local: attrObject}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrID}, Value: strconv.FormatUint(uint64(r.ID), 10)},
	}}
	if r.ObjectType != ObjectTypeModel {
		xo.Attr = append(xo.Attr, xml.Attr{Name: xml.Name{Local: attrType}, Value: r.ObjectType.String()})
	}
	if r.Thumbnail != "" {
		xo.Attr = append(xo.Attr, xml.Attr{Name: xml.Name{Local: attrThumbnail}, Value: r.Thumbnail})
	}
	if r.PartNumber != "" {
		xo.Attr = append(xo.Attr, xml.Attr{Name: xml.Name{Local: attrPartNumber}, Value: r.PartNumber})
	}
	if r.Name != "" {
		xo.Attr = append(xo.Attr, xml.Attr{Name: xml.Name{Local: attrName}, Value: r.Name})
	}
	if r.Mesh != nil {
		if r.DefaultPropertyID != 0 {
			xo.Attr = append(xo.Attr, xml.Attr{
				Name: xml.Name{Local: attrPID}, Value: strconv.FormatUint(uint64(r.DefaultPropertyID), 10),
			})
		}
		if r.DefaultPropertyIndex != 0 {
			xo.Attr = append(xo.Attr, xml.Attr{
				Name: xml.Name{Local: attrPIndex}, Value: strconv.FormatUint(uint64(r.DefaultPropertyIndex), 10),
			})
		}
	}
	x.EncodeToken(xo)

	if len(r.Metadata) != 0 {
		e.writeMetadataGroup(x, r.Metadata)
	}

	if r.Mesh != nil {
		e.writeMesh(x, r.Mesh)
	} else {
		e.writeComponents(x, r.Components)
	}
	x.EncodeToken(xo.End())
}

func (e *Encoder) writeComponents(x tokenEncoder, comps []*Component) {
	xcs := xml.StartElement{Name: xml.Name{Local: attrComponents}}
	x.EncodeToken(xcs)
	x.SetAutoClose(true)
	for _, c := range comps {
		t := xml.StartElement{
			Name: xml.Name{Local: attrComponent}, Attr: []xml.Attr{
				{Name: xml.Name{Local: attrObjectID}, Value: strconv.FormatUint(uint64(c.ObjectID), 10)},
			},
		}
		if c.HasTransform() {
			t.Attr = append(t.Attr, xml.Attr{Name: xml.Name{Local: attrTransform}, Value: FormatMatrix(c.Transform)})
		}
		x.EncodeToken(t)
	}
	x.SetAutoClose(false)
	x.EncodeToken(xcs.End())
}

func (e *Encoder) writeMesh(x tokenEncoder, m *Mesh) {
	xm := xml.StartElement{Name: xml.Name{Local: attrMesh}}
	x.EncodeToken(xm)
	xvs := xml.StartElement{Name: xml.Name{Local: attrVertices}}
	x.EncodeToken(xvs)
	x.SetAutoClose(true)
	for _, v := range m.Nodes {
		x.EncodeToken(xml.StartElement{
			Name: xml.Name{Local: attrVertex},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: attrX}, Value: strconv.FormatFloat(float64(v.X()), 'f', 3, 32)},
				{Name: xml.Name{Local: attrY}, Value: strconv.FormatFloat(float64(v.Y()), 'f', 3, 32)},
				{Name: xml.Name{Local: attrZ}, Value: strconv.FormatFloat(float64(v.Z()), 'f', 3, 32)},
			},
		})
	}
	x.SetAutoClose(false)
	x.EncodeToken(xvs.End())

	xvt := xml.StartElement{Name: xml.Name{Local: attrTriangles}}
	x.EncodeToken(xvt)
	x.SetAutoClose(true)
	for _, v := range m.Faces {
		t := xml.StartElement{
			Name: xml.Name{Local: attrTriangle},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: attrV1}, Value: strconv.FormatUint(uint64(v.NodeIndices[0]), 10)},
				{Name: xml.Name{Local: attrV2}, Value: strconv.FormatUint(uint64(v.NodeIndices[1]), 10)},
				{Name: xml.Name{Local: attrV3}, Value: strconv.FormatUint(uint64(v.NodeIndices[2]), 10)},
			},
		}
		if v.Resource != 0 {
			t.Attr = append(t.Attr, xml.Attr{
				Name: xml.Name{Local: attrPID}, Value: strconv.FormatUint(uint64(v.Resource), 10),
			})
			if v.ResourceIndices[0] != 0 {
				t.Attr = append(t.Attr, xml.Attr{
					Name: xml.Name{Local: attrP1}, Value: strconv.FormatUint(uint64(v.ResourceIndices[0]), 10),
				})
				if v.ResourceIndices[1] != 0 {
					t.Attr = append(t.Attr, xml.Attr{
						Name: xml.Name{Local: attrP2}, Value: strconv.FormatUint(uint64(v.ResourceIndices[1]), 10),
					})
				}
				if v.ResourceIndices[2] != 0 {
					t.Attr = append(t.Attr, xml.Attr{
						Name: xml.Name{Local: attrP3}, Value: strconv.FormatUint(uint64(v.ResourceIndices[2]), 10),
					})
				}
			}
		}
		x.EncodeToken(t)
	}
	x.SetAutoClose(false)
	x.EncodeToken(xvt.End())
	x.EncodeToken(xm.End())
}

func (e *Encoder) writeBaseMaterial(x tokenEncoder, r *BaseMaterialsResource) {
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
