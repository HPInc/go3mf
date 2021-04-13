package booleanoperations

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf/spec"
)

// Marshal3MFAttr encodes the resource attributes.
func (u *OperationAttr) Marshal3MFAttr(_ spec.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrCompsBoolOperOperation}, Value: u.operation.String()},
	}, nil
}

func (u *AssociationAttr) Marshal3MFAttr(_ spec.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: Namespace, Local: attrCompsBoolOperAssociation}, Value: u.association.String()},
	}, nil
}
