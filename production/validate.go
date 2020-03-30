package production

import (
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func (e *Spec) ValidateAsset(_ *go3mf.Model, _ string, _ go3mf.Asset) error {
	return nil
}

func (e *Spec) ValidateModel(m *go3mf.Model) error {
	var (
		u    *UUID
		errs error
	)
	if !m.Build.AnyAttr.Get(&u) {
		errs = errors.Append(errs, errors.Wrap(&errors.MissingFieldError{Name: attrProdUUID}, m.Build))
	} else if validateUUID(string(*u)) != nil {
		errs = errors.Append(errs, errors.Wrap(errors.ErrUUID, m.Build))
	}
	for i, item := range m.Build.Items {
		var iErrs error
		var p *PathUUID
		if !item.AnyAttr.Get(&p) {
			iErrs = errors.Append(iErrs, &errors.MissingFieldError{Name: attrProdUUID})
		} else {
			iErrs = errors.Append(iErrs, e.validatePathUUID(m, "", p))
		}
		if iErrs != nil {
			errs = errors.Append(errs, errors.Wrap(errors.WrapIndex(iErrs, item, i), m.Build))
		}
	}
	return errs
}

func (e *Spec) ValidateObject(m *go3mf.Model, path string, obj *go3mf.Object) error {
	var (
		u    *UUID
		errs error
	)
	if !obj.AnyAttr.Get(&u) {
		errs = errors.Append(errs, &errors.MissingFieldError{Name: attrProdUUID})
	} else if validateUUID(string(*u)) != nil {
		errs = errors.Append(errs, errors.ErrUUID)
	}
	var p *PathUUID
	for i, c := range obj.Components {
		var cErrs error
		if !c.AnyAttr.Get(&p) {
			cErrs = errors.Append(cErrs, &errors.MissingFieldError{Name: attrProdUUID})
		} else {
			cErrs = errors.Append(cErrs, e.validatePathUUID(m, path, p))
		}
		if cErrs != nil {
			errs = errors.Append(errs, errors.WrapIndex(cErrs, c, i))
		}
	}
	return errs
}

func (e *Spec) validatePathUUID(m *go3mf.Model, path string, p *PathUUID) error {
	var errs error
	if p.UUID == "" {
		errs = errors.Append(errs, &errors.MissingFieldError{Name: attrProdUUID})
	} else if validateUUID(string(p.UUID)) != nil {
		errs = errors.Append(errs, errors.ErrUUID)
	}
	if p.Path != "" {
		if path == "" || path == m.PathOrDefault() { // root
			// Path is validated as part if the core validations
			if !e.Required() {
				errs = errors.Append(errs, errors.ErrProdExtRequired)
			}
		} else {
			errs = errors.Append(errs, errors.ErrProdRefInNonRoot)
		}
	}
	return errs
}
