package booleanoperations

import (
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"

	"github.com/qmuntal/go3mf/spec"
)

func (Spec) CreateElementDecoder(_ interface{}, _ string) spec.ElementDecoder {
	return nil
}

func (Spec) DecodeAttribute(parentNode interface{}, attr spec.Attr) (errs error) {
	if t, ok := parentNode.(*go3mf.Components); ok {
		switch attr.Name.Local {
		case attrCompsBoolOperAssociation:
			if ext := GetBooleanOperationAttr(t); ext != nil {
				if association, ok := newAssociation(string(attr.Value)); ok {
					ext.Association = association
				} else {
					errs = specerr.Append(errs, specerr.NewParseAttrError(attr.Name.Local, true))
				}
			} else {
				if association, ok := newAssociation(string(attr.Value)); ok {
					t.AnyAttr = append(t.AnyAttr, &BooleanOperationAttr{Association: association})
				} else {
					errs = specerr.Append(errs, specerr.NewParseAttrError(attr.Name.Local, true))
				}
			}
		case attrCompsBoolOperOperation:
			if ext := GetBooleanOperationAttr(t); ext != nil {
				if operation, ok := newOperation(string(attr.Value)); ok {
					ext.Operation = operation
				} else {
					errs = specerr.Append(errs, specerr.NewParseAttrError(attr.Name.Local, true))
				}
			} else {
				if operation, ok := newOperation(string(attr.Value)); ok {
					t.AnyAttr = append(t.AnyAttr, &BooleanOperationAttr{Operation: operation})
				} else {
					errs = specerr.Append(errs, specerr.NewParseAttrError(attr.Name.Local, true))
				}
			}

		}
	}

	return
}
