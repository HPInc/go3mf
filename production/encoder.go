package production

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
	specxml "github.com/qmuntal/go3mf/spec/xml"
	"github.com/qmuntal/go3mf/uuid"
)

func (s *Spec) PreProcessEncode() {
	if s.DisableAutoUUID {
		return
	}
	var buildAttr *BuildAttr
	if !s.m.Build.AnyAttr.Get(&buildAttr) {
		s.m.Build.AnyAttr = append(s.m.Build.AnyAttr, &BuildAttr{UUID: uuid.New()})
	}
	for _, item := range s.m.Build.Items {
		var itemAttr *ItemAttr
		if !item.AnyAttr.Get(&itemAttr) {
			item.AnyAttr = append(item.AnyAttr, &ItemAttr{UUID: uuid.New()})
		}
	}
	s.m.WalkObjects(func(s string, o *go3mf.Object) error {
		var objAttr *ObjectAttr
		if !o.AnyAttr.Get(&objAttr) {
			o.AnyAttr = append(o.AnyAttr, &ObjectAttr{UUID: uuid.New()})
		}
		for _, c := range o.Components {
			var compAttr *ComponentAttr
			if !c.AnyAttr.Get(&compAttr) {
				c.AnyAttr = append(c.AnyAttr, &ComponentAttr{UUID: uuid.New()})
			}
		}
		return nil
	})
}

// Marshal3MFAttr encodes the resource attributes.
func (u *BuildAttr) Marshal3MFAttr(_ specxml.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrProdUUID}, Value: u.UUID},
	}, nil
}

// Marshal3MFAttr encodes the resource attributes.
func (u *ObjectAttr) Marshal3MFAttr(_ specxml.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrProdUUID}, Value: u.UUID},
	}, nil
}

// Marshal3MFAttr encodes the resource attributes.
func (p *ItemAttr) Marshal3MFAttr(_ specxml.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrPath}, Value: p.Path},
		{Name: xml.Name{Space: Namespace, Local: attrProdUUID}, Value: p.UUID},
	}, nil
}

// Marshal3MFAttr encodes the resource attributes.
func (p *ComponentAttr) Marshal3MFAttr(_ specxml.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrPath}, Value: p.Path},
		{Name: xml.Name{Space: Namespace, Local: attrProdUUID}, Value: p.UUID},
	}, nil
}
