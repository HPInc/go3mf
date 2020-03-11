package production

import (
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

func (u *UUID) Validate(m *go3mf.Model, path string, _ interface{}) []error {
	if validateUUID(string(*u)) != nil {
		return []error{specerr.ErrUUID}
	}
	return nil
}

func (p *PathUUID) Validate(m *go3mf.Model, path string, _ interface{}) []error {
	var errs []error
	if p.UUID == "" {
		errs = append(errs, &specerr.MissingFieldError{Name: attrProdUUID})
	} else if validateUUID(string(p.UUID)) != nil {
		errs = append(errs, specerr.ErrUUID)
	}
	if p.Path != "" {
		if m.PathOrDefault() == path { // root
			// Path is validated as part if the core validations
			var extRequired bool
			for _, r := range m.RequiredExtensions {
				if r == ExtensionName {
					extRequired = true
					break
				}
			}
			if !extRequired {
				errs = append(errs, specerr.ErrProdExtRequired)
			}
		} else {
			errs = append(errs, specerr.ErrProdRefInNonRoot)
		}
	}
	return errs
}
