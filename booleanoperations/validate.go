package booleanoperations

import (
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func (Spec) Validate(model interface{}, path string, e interface{}) error {
	switch e := e.(type) {
	case *go3mf.Model:
		return validateModel(e)
	}
	return nil
}
func validateModel(m *go3mf.Model) error {
	var errs error
	for i, obj := range m.Resources.Objects {
		var iErrs error
		iErrs = errors.Append(iErrs, validateObject(m, "", obj))
		if iErrs != nil {
			errs = errors.Append(errs, errors.Wrap(errors.WrapIndex(iErrs, obj, i), m.Resources))
		}
	}
	return errs
}
func validateObject(m *go3mf.Model, path string, obj *go3mf.Object) error {
	components := obj.Components
	if components != nil {
		return validateComponents(obj.Components)
	}
	return nil
}
func validateComponents(m *go3mf.Components) error {
	var errs error
	if p := GetOperationAttr(m); p == nil {
		errs = errors.Append(errs, errors.Wrap(errors.NewMissingFieldError(attrCompsBoolOperOperation), m))
	}
	if p := GetAssociationAttr(m); p == nil {
		errs = errors.Append(errs, errors.Wrap(errors.NewMissingFieldError(attrCompsBoolOperAssociation), m))
	}
	return errs
}
