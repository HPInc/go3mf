package production

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
)

// RegisterExtension registers this extension in the decoder instance.
func RegisterExtension(d *go3mf.Decoder) {
	d.RegisterDecodeAttributeExtension(ExtensionName, decodeAttribute)
	d.RegisterFileFilterExtension(ExtensionName, fileFilter)
}

func fileFilter(relType string) bool {
	return relType == go3mf.RelTypeModel3D
}

func decodeAttribute(s *go3mf.Scanner, parentNode interface{}, attr xml.Attr) {
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
			*BuildAttr(t) = uuid
		}
	case *go3mf.Item:
		switch attr.Name.Local {
		case attrProdUUID:
			if uuid, err = NewUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
			}
			ItemAttr(t).UUID = uuid
		case attrPath:
			ItemAttr(t).Path = attr.Value
		}
	case *go3mf.ObjectResource:
		if attr.Name.Local == attrProdUUID {
			if uuid, err = NewUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
			}
			*ObjectAttr(t) = uuid
		}
	case *go3mf.Component:
		switch attr.Name.Local {
		case attrProdUUID:
			if uuid, err = NewUUID(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
			}
			ComponentAttr(t).UUID = uuid
		case attrPath:
			ComponentAttr(t).Path = attr.Value
		}
	}
}
