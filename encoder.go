package go3mf

import (
	"context"
	"encoding/xml"
	"github.com/qmuntal/opc"
	"io"
	"strconv"
)

type Encoder struct {
	w *opc.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: opc.NewWriter(w),
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
	_, err =w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	x := xml.NewEncoder(w)
	err = x.EncodeToken(xml.StartElement{Name: xml.Name{Local: attrModel}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrXmlns}, Value: ExtensionName},
		{Name: xml.Name{Local: attrUnit}, Value: m.Units.String()},
		{Name: xml.Name{Space: nsXML, Local: attrLang}, Value: m.Language},
	}})
	if err != nil {
		return err
	}

	if err = e.writeMetadata(x, m.Metadata); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) writeResources(x *xml.Encoder, m *Model) error {
	xt := newXmlNodeEncoder(x, attrResources, 0)
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

func (e *Encoder) writeMetadata(x *xml.Encoder, metadata []Metadata) error {
	for _, md := range metadata {
		xn := newXmlNodeEncoder(x, attrMetadata, 3)
		xn.Attribute(attrName, md.Name)
		xn.Attribute(attrPreserve, strconv.FormatBool(md.Preserve))
		xn.OptionalAttribute(attrType, md.Type)
		if err := xn.TextEnd(md.Value); err != nil {
			return err
		}
	}
	return nil
}


func (e *Encoder) writeObject(x *xml.Encoder, r *ObjectResource) error {
	xn := newXmlNodeEncoder(x, attrObject, 7)
	xn.Attribute(attrID, strconv.Itoa(int(r.ID)))
	xn.Attribute(attrType, r.ObjectType.String())
	xn.Attribute(attrThumbnail, r.Thumbnail)
	xn.Attribute(attrPartNumber, r.PartNumber)
	xn.Attribute(attrName, r.Name)
	xn.Attribute(attrPID, strconv.Itoa(int(r.DefaultPropertyID)))
	xn.Attribute(attrPID, strconv.Itoa(int(r.DefaultPropertyIndex)))
	xn.Close()

	xm := newXmlNodeEncoder(x, attrMetadataGroup, 0)
	if err := xm.Close(); err != nil {
		return err
	}
	if err := e.writeMetadata(x, r.Metadata); err != nil {
		return err
	}
	if err := xm.End(); err != nil {
		return err
	}
	  

	return xn.End()
}

func (e *Encoder) writeBaseMaterial(x *xml.Encoder, r *BaseMaterialsResource) error {
	xt := newXmlNodeEncoder(x, attrBaseMaterials, 1)
	xt.Attribute(attrID, strconv.Itoa(int(r.ID)))
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
	x *xml.Encoder
	start xml.StartElement
	closed bool
}

func newXmlNodeEncoder(x *xml.Encoder, name string, cap int) *xmlNodeEncoder {
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
		Name: xml.Name{Local: name},
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
	return e.x.EncodeToken(e.End())
}

func (e *xmlNodeEncoder) Close() error {
	if !e.closed {
		return e.x.EncodeToken(e.start)
	}
	return nil
}

func (e *xmlNodeEncoder) End() error {
	if err := e.Close(); err != nil {
		return err
	}
	return e.x.EncodeToken(e.End())
}