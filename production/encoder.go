package production

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
)

func (u *UUID) Marshal3MFAttr(_ *go3mf.XMLEncoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: ExtensionName, Local: attrProdUUID}, Value: string(*u)},
	}, nil
}

func (p *PathUUID) Marshal3MFAttr(_ *go3mf.XMLEncoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: ExtensionName, Local: attrPath}, Value: string(p.Path)},
		{Name: xml.Name{Space: ExtensionName, Local: attrProdUUID}, Value: string(p.UUID)},
	}, nil
}
