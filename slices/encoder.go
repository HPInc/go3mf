package slices

import (
	"encoding/xml"
	"strconv"

	"github.com/qmuntal/go3mf"
)

func (s *SliceStackResource) Marshal3MF(x *go3mf.XMLEncoder, _ xml.StartElement) error {
	xs := xml.StartElement{Name: xml.Name{Space: ExtensionName, Local: attrSliceStack}, Attr: []xml.Attr{
		{Name: xml.Name{Space: ExtensionName, Local: attrID}, Value: strconv.FormatUint(uint64(s.ID), 10)},
	}}
	x.EncodeToken(xs)

	x.EncodeToken(xs.End())
	return x.Flush()
}
