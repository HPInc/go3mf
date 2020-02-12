package production

import "encoding/xml"

func (u *UUID) Marshal3MFAttr() ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: ExtensionName, Local: attrProdUUID}, Value: string(*u)},
	}, nil
}

func (p *PathUUID) Marshal3MFAttr() ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: ExtensionName, Local: attrPath}, Value: string(p.Path)},
		{Name: xml.Name{Space: ExtensionName, Local: attrProdUUID}, Value: string(p.UUID)},
	}, nil
}
