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
	var errs error
	u := GetBuildAttr(&e.m.Build)
	if u == nil {
		errs = errors.Append(errs, errors.Wrap(errors.NewMissingFieldError(attrProdUUID), e.m.Build))
	} else if uuid.Validate(u.UUID) != nil {
		errs = errors.Append(errs, errors.Wrap(ErrUUID, e.m.Build))
	}
	for i, item := range e.m.Build.Items {
		var iErrs error

		if p := GetItemAttr(item); p != nil {
			iErrs = errors.Append(iErrs, e.validatePathUUID("", p))
		} else {
			iErrs = errors.Append(iErrs, errors.NewMissingFieldError(attrProdUUID))
		}
		if iErrs != nil {
			errs = errors.Append(errs, errors.Wrap(errors.WrapIndex(iErrs, item, i), e.m.Build))
		}
	}
	return errs
}

func (e *Spec) ValidateObject(path string, obj *go3mf.Object) error {
	var errs error
	u := GetObjectAttr(obj)
	if u == nil {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrProdUUID))
	} else if uuid.Validate(u.UUID) != nil {
		errs = errors.Append(errs, ErrUUID)
	}
	for i, c := range obj.Components {
		var cErrs error
		if p := GetComponentAttr(c); p != nil {
			cErrs = errors.Append(cErrs, e.validatePathUUID(path, p))
		} else {
			cErrs = errors.Append(cErrs, errors.NewMissingFieldError(attrProdUUID))
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
		errs = errors.Append(errs, ErrUUID)
	}
	if p.ObjectPath() != "" {
		if path == "" || path == e.m.PathOrDefault() { // root
			// Path is validated as part if the core validations
		} else {
			errs = errors.Append(errs, ErrProdRefInNonRoot)
		}
	}
	return errs
}
