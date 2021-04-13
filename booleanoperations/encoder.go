package booleanoperations

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf/spec"
)

// Marshal3MFAttr encodes the resource attributes.
func (u *BooleanOperationAttr) Marshal3MFAttr(_ spec.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrCompsBoolOperAssociation}, Value: u.association.String()},
		{Name: xml.Name{Space: Namespace, Local: attrCompsBoolOperOperation}, Value: u.operation.String()},
	}, nil
}
