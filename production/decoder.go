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

func fileFilter(relType string, isRootModel bool) bool {
	return isRootModel && relType == go3mf.RelTypeModel3D
}

func decodeAttribute(s *go3mf.Scanner, parentNode interface{}, attr xml.Attr) {
	switch t := parentNode.(type) {
	case *go3mf.Build:
		if attr.Name.Local == attrProdUUID {
			if err := BuildAttr(t).Set(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
			}
		}
	case *go3mf.Item:
		switch attr.Name.Local {
		case attrProdUUID:
			if err := ItemAttr(t).UUID.Set(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
			}
		case attrPath:
			ItemAttr(t).Path = attr.Value
		}
	case *go3mf.ObjectResource:
		if attr.Name.Local == attrProdUUID {
			if err := ObjectAttr(t).Set(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
			}
		}
	case *go3mf.Component:
		switch attr.Name.Local {
		case attrProdUUID:
			if err := ComponentAttr(t).UUID.Set(attr.Value); err != nil {
				s.InvalidAttr(attr.Name.Local, attr.Value, true)
			}
		case attrPath:
			ComponentAttr(t).Path = attr.Value
		}
	}
}
