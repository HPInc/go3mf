package production

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
)

func (e Spec) NewNodeDecoder(_ interface{}, _ string) go3mf.NodeDecoder {
	return nil
}

func (e Spec) DecodeAttribute(s *go3mf.Scanner, parentNode interface{}, attr xml.Attr) {
	var (
		uuid UUID
		err  error
	)
	switch t := parentNode.(type) {
	case *go3mf.Build:
		if attr.Name.Local == attrProdUUID {
			if uuid, err = NewUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
			}
			t.AnyAttr = append(t.AnyAttr, &uuid)
		}
	case *go3mf.Item:
		switch attr.Name.Local {
		case attrProdUUID:
			if uuid, err = NewUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
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
			if uuid, err = NewUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
			}
			t.AnyAttr = append(t.AnyAttr, &uuid)
		}
	case *go3mf.Component:
		switch attr.Name.Local {
		case attrProdUUID:
			if uuid, err = NewUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
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
