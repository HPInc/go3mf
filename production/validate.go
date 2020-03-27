package production

import (
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

func (e *Spec) ValidateAsset(_ *go3mf.Model, _ string, _ go3mf.Asset) error {
	return nil
}

func (e *Spec) ValidateModel(m *go3mf.Model) error {
	var (
		u    *UUID
		errs = new(specerr.ErrorList)
	)
	if !m.Build.AnyAttr.Get(&u) {
		errs.Append(specerr.New(m.Build, &specerr.MissingFieldError{Name: attrProdUUID}))
	} else if validateUUID(string(*u)) != nil {
		errs.Append(specerr.New(m.Build, specerr.ErrUUID))
	}
	for i, item := range m.Build.Items {
		iErrs := new(specerr.ErrorList)
		var p *PathUUID
		if !item.AnyAttr.Get(&p) {
			iErrs.Append(&specerr.MissingFieldError{Name: attrProdUUID})
		} else {
			iErrs.Append(e.validatePathUUID(m, "", p))
		}
		errs.Append(specerr.New(m.Build, specerr.NewIndexed(item, i, iErrs)))
	}
	return errs.ErrorOrNil()
}

func (e *Spec) ValidateObject(m *go3mf.Model, path string, obj *go3mf.Object) error {
	var (
		u    *UUID
		errs = new(specerr.ErrorList)
	)
	if !obj.AnyAttr.Get(&u) {
		errs.Append(&specerr.MissingFieldError{Name: attrProdUUID})
	} else if validateUUID(string(*u)) != nil {
		errs.Append(specerr.ErrUUID)
	}
	var p *PathUUID
	for i, c := range obj.Components {
		cErrs := new(specerr.ErrorList)
		if !c.AnyAttr.Get(&p) {
			cErrs.Append(&specerr.MissingFieldError{Name: attrProdUUID})
		} else {
			cErrs.Append(e.validatePathUUID(m, path, p))
		}
		errs.Append(specerr.NewIndexed(c, i, cErrs))
	}
	return errs.ErrorOrNil()
}

func (e *Spec) validatePathUUID(m *go3mf.Model, path string, p *PathUUID) error {
	errs := new(specerr.ErrorList)
	if p.UUID == "" {
		errs.Append(&specerr.MissingFieldError{Name: attrProdUUID})
	} else if validateUUID(string(p.UUID)) != nil {
		errs.Append(specerr.ErrUUID)
	}
	if p.Path != "" {
		if path == "" || path == m.PathOrDefault() { // root
			// Path is validated as part if the core validations
			if !e.Required() {
				errs.Append(specerr.ErrProdExtRequired)
			}
		} else {
			errs.Append(specerr.ErrProdRefInNonRoot)
		}
	}
	return errs.ErrorOrNil()
}
