package booleanoperations

import (
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"

	"github.com/qmuntal/go3mf/spec"
)

func (Spec) CreateElementDecoder(_ interface{}, _ string) spec.ElementDecoder {
	return nil
}

func (Spec) DecodeAttribute(parentNode interface{}, attr spec.Attr) (errs error) {
	t, ok := parentNode.(*go3mf.Components)
	if ok {
		switch attr.Name.Local {
		case attrCompsBoolOperAssociation:
			if ext := GetAssociationAttr(t); ext != nil {
				association, ok := newAssociation(string(attr.Value))
				if ok {
					t.AnyAttr = append(t.AnyAttr, &AssociationAttr{association: association})
				} else {
					errs = errors.Append(errs, errors.NewParseAttrError(attr.Name.Local, true))
				}
			} else {
				association, ok := newAssociation(string(attr.Value))
				if ok {
					t.AnyAttr = append(t.AnyAttr, &AssociationAttr{association: association})
				}
			}
		case attrCompsBoolOperOperation:
			if ext := GetOperationAttr(t); ext != nil {
				operation, ok := newOperation(string(attr.Value))
				if ok {
					t.AnyAttr = append(t.AnyAttr, &OperationAttr{operation: operation})
				} else {
					errs = errors.Append(errs, errors.NewParseAttrError(attr.Name.Local, true))
				}
			} else {
				operation, ok := newOperation(string(attr.Value))
				if ok {
					t.AnyAttr = append(t.AnyAttr, &OperationAttr{operation: operation})
				}
			}

		}
	}

	return
}
