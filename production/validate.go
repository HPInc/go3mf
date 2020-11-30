package production

import (
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
	"github.com/qmuntal/go3mf/uuid"
)

type uuidPath interface {
	getUUID() string
	ObjectPath() string
}

func (e *Spec) ValidateModel() error {
	var (
		u    *BuildAttr
		errs error
	)
	if !e.m.Build.AnyAttr.Get(&u) {
		errs = errors.Append(errs, errors.Wrap(errors.NewMissingFieldError(attrProdUUID), e.m.Build))
	} else if uuid.Validate(u.UUID) != nil {
		errs = errors.Append(errs, errors.Wrap(errors.ErrUUID, e.m.Build))
	}
	for i, item := range e.m.Build.Items {
		var iErrs error
		var p *ItemAttr
		if !item.AnyAttr.Get(&p) {
			iErrs = errors.Append(iErrs, errors.NewMissingFieldError(attrProdUUID))
		} else {
			iErrs = errors.Append(iErrs, e.validatePathUUID("", p))
		}
		if iErrs != nil {
			errs = errors.Append(errs, errors.Wrap(errors.WrapIndex(iErrs, item, i), e.m.Build))
		}
	}
	return errs
}

func (e *Spec) ValidateObject(path string, obj *go3mf.Object) error {
	var (
		u    *ObjectAttr
		errs error
	)
	if !obj.AnyAttr.Get(&u) {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrProdUUID))
	} else if uuid.Validate(u.UUID) != nil {
		errs = errors.Append(errs, errors.ErrUUID)
	}
	var p *ComponentAttr
	for i, c := range obj.Components {
		var cErrs error
		if !c.AnyAttr.Get(&p) {
			cErrs = errors.Append(cErrs, errors.NewMissingFieldError(attrProdUUID))
		} else {
			cErrs = errors.Append(cErrs, e.validatePathUUID(path, p))
		}
		if cErrs != nil {
			errs = errors.Append(errs, errors.WrapIndex(cErrs, c, i))
		}
	}
	return errs
}

func (e *Spec) validatePathUUID(path string, p uuidPath) error {
	var errs error
	if p.getUUID() == "" {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrProdUUID))
	} else if uuid.Validate(string(p.getUUID())) != nil {
		errs = errors.Append(errs, errors.ErrUUID)
	}
	if p.ObjectPath() != "" {
		if path == "" || path == e.m.PathOrDefault() { // root
			// Path is validated as part if the core validations
		} else {
			errs = errors.Append(errs, errors.ErrProdRefInNonRoot)
		}
	}
	return errs
}
