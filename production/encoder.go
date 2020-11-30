package production

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/spec/encoding"
	"github.com/qmuntal/go3mf/uuid"
)

func (s *Spec) PreProcessEncode() {
	if s.DisableAutoUUID {
		return
	}
	buildAttr := GetBuildAttr(&s.m.Build)
	if buildAttr == nil {
		s.m.Build.AnyAttr = append(s.m.Build.AnyAttr, &BuildAttr{UUID: uuid.New()})
	}
	for _, item := range s.m.Build.Items {
		itemAttr := GetItemAttr(item)
		if itemAttr == nil {
			item.AnyAttr = append(item.AnyAttr, &ItemAttr{UUID: uuid.New()})
		}
	}
	s.m.WalkObjects(func(s string, o *go3mf.Object) error {
		objAttr := GetObjectAttr(o)
		if objAttr == nil {
			o.AnyAttr = append(o.AnyAttr, &ObjectAttr{UUID: uuid.New()})
		}
		for _, c := range o.Components {
			compAttr := GetComponentAttr(c)
			if compAttr == nil {
				c.AnyAttr = append(c.AnyAttr, &ComponentAttr{UUID: uuid.New()})
			}
		}
		return nil
	})
}

// Marshal3MFAttr encodes the resource attributes.
func (u *BuildAttr) Marshal3MFAttr(_ encoding.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrProdUUID}, Value: u.UUID},
	}, nil
}

// Marshal3MFAttr encodes the resource attributes.
func (u *ObjectAttr) Marshal3MFAttr(_ encoding.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrProdUUID}, Value: u.UUID},
	}, nil
}

// Marshal3MFAttr encodes the resource attributes.
func (p *ItemAttr) Marshal3MFAttr(_ encoding.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrPath}, Value: p.Path},
		{Name: xml.Name{Space: Namespace, Local: attrProdUUID}, Value: p.UUID},
	}, nil
}

// Marshal3MFAttr encodes the resource attributes.
func (p *ComponentAttr) Marshal3MFAttr(_ encoding.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrPath}, Value: p.Path},
		{Name: xml.Name{Space: Namespace, Local: attrProdUUID}, Value: p.UUID},
	}, nil
}
