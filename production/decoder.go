package production

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
)

func init() {
	go3mf.RegisterExtensionDecoder(ExtensionName, &extensionDecoder{})
}

type extensionDecoder struct{}

func (d *extensionDecoder) NodeDecoder(_ interface{}, nodeName string) go3mf.NodeDecoder {
	return nil
}

func (d *extensionDecoder) DecodeAttribute(s *go3mf.Scanner, parentNode interface{}, attr xml.Attr) {
	switch t := parentNode.(type) {
	case *go3mf.Build:
		if attr.Name.Local == attrProdUUID {
			if err := validateUUID(attr.Value); err != nil {
				s.InvalidRequiredAttr(attr.Name.Local, attr.Value)
			}
			ExtensionBuild(t).UUID = attr.Value
		}
	case *go3mf.Item:
		switch attr.Name.Local {
		case attrProdUUID:
			if err := validateUUID(attr.Value); err != nil {
				s.InvalidRequiredAttr(attr.Name.Local, attr.Value)
			}
			ExtensionItem(t).UUID = attr.Value
		case attrPath:
			ExtensionItem(t).Path = attr.Value
		}
	case *go3mf.ObjectResource:
		if attr.Name.Local == attrProdUUID {
			if err := validateUUID(attr.Value); err != nil {
				s.InvalidRequiredAttr(attr.Name.Local, attr.Value)
			}
			ExtensionObject(t).UUID = attr.Value
		}
	case *go3mf.Component:
		switch attr.Name.Local {
		case attrProdUUID:
			if err := validateUUID(attr.Value); err != nil {
				s.InvalidRequiredAttr(attr.Name.Local, attr.Value)
			}
			ExtensionComponent(t).UUID = attr.Value
		case attrPath:
			ExtensionComponent(t).Path = attr.Value
		}
	}
}