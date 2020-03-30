package production

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
)

func (e Spec) NewNodeDecoder(_ interface{}, _ string) go3mf.NodeDecoder {
	return nil
}

func (e Spec) OnDecoded(m *go3mf.Model) error {
	var (
		uuid *UUID
		pu   *PathUUID
	)
	if !m.Build.AnyAttr.Get(&uuid) {
		m.Build.AnyAttr = append(m.Build.AnyAttr, NewUUID())
	}
	for _, item := range m.Build.Items {
		if !item.AnyAttr.Get(&pu) {
			item.AnyAttr = append(item.AnyAttr, &PathUUID{
				UUID: *NewUUID(),
			})
		} else if pu.UUID == "" {
			pu.UUID = *NewUUID()
		}
	}
	e.fillResourceUUID(&m.Resources)
	for _, c := range m.Childs {
		e.fillResourceUUID(&c.Resources)
	}
	return nil
}

func (e Spec) fillResourceUUID(res *go3mf.Resources) {
	var (
		pu   *PathUUID
		uuid *UUID
	)
	for _, o := range res.Objects {
		if !o.AnyAttr.Get(&uuid) {
			o.AnyAttr = append(o.AnyAttr, NewUUID())
		}
		for _, c := range o.Components {
			if !c.AnyAttr.Get(&pu) {
				c.AnyAttr = append(c.AnyAttr, &PathUUID{
					UUID: *NewUUID(),
				})
			} else if pu.UUID == "" {
				pu.UUID = *NewUUID()
			}
		}
	}
}

func (e Spec) DecodeAttribute(s *go3mf.Scanner, parentNode interface{}, attr xml.Attr) {
	var (
		uuid UUID
		err  error
	)
	switch t := parentNode.(type) {
	case *go3mf.Build:
		if attr.Name.Local == attrProdUUID {
			if uuid, err = ParseUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, true)
			}
			t.AnyAttr = append(t.AnyAttr, &uuid)
		}
	case *go3mf.Item:
		switch attr.Name.Local {
		case attrProdUUID:
			if uuid, err = ParseUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, true)
			}
			var ext *PathUUID
			if t.AnyAttr.Get(&ext) {
				ext.UUID = uuid
			} else {
				t.AnyAttr = append(t.AnyAttr, &PathUUID{UUID: uuid})
			}
		case attrPath:
			var ext *PathUUID
			if t.AnyAttr.Get(&ext) {
				ext.Path = attr.Value
			} else {
				t.AnyAttr = append(t.AnyAttr, &PathUUID{Path: attr.Value})
			}
		}
	case *go3mf.Object:
		if attr.Name.Local == attrProdUUID {
			if uuid, err = ParseUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, true)
			}
			t.AnyAttr = append(t.AnyAttr, &uuid)
		}
	case *go3mf.Component:
		switch attr.Name.Local {
		case attrProdUUID:
			if uuid, err = ParseUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, true)
			}
			var ext *PathUUID
			if t.AnyAttr.Get(&ext) {
				ext.UUID = uuid
			} else {
				t.AnyAttr = append(t.AnyAttr, &PathUUID{UUID: uuid})
			}
		case attrPath:
			var ext *PathUUID
			if t.AnyAttr.Get(&ext) {
				ext.Path = attr.Value
			} else {
				t.AnyAttr = append(t.AnyAttr, &PathUUID{Path: attr.Value})
			}
		}
	}
}
