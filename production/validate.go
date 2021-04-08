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

func (Spec) Validate(model interface{}, path string, e interface{}) error {
	switch e := e.(type) {
	case *go3mf.Model:
		return validateModel(e)
	case *go3mf.Object:
		return validateObject(model.(*go3mf.Model), path, e)
	}
	return nil
}

func validateModel(m *go3mf.Model) error {
	var errs error
	u := GetBuildAttr(&m.Build)
	if u == nil {
		errs = errors.Append(errs, errors.Wrap(errors.NewMissingFieldError(attrProdUUID), m.Build))
	} else if uuid.Validate(u.UUID) != nil {
		errs = errors.Append(errs, errors.Wrap(ErrUUID, m.Build))
	}
	for i, item := range m.Build.Items {
		var iErrs error

		if p := GetItemAttr(item); p != nil {
			iErrs = errors.Append(iErrs, validatePathUUID(m, "", p))
		} else {
			iErrs = errors.Append(iErrs, errors.NewMissingFieldError(attrProdUUID))
		}
		if iErrs != nil {
			errs = errors.Append(errs, errors.Wrap(errors.WrapIndex(iErrs, item, i), m.Build))
		}
	}
	return errs
}

func validateObject(m *go3mf.Model, path string, obj *go3mf.Object) error {
	var errs error
	u := GetObjectAttr(obj)
	if u == nil {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrProdUUID))
	} else if uuid.Validate(u.UUID) != nil {
		errs = errors.Append(errs, ErrUUID)
	}
	if obj.Components != nil {
		var cErrs error
		for i, c := range obj.Components.Component {
			var err error
			if p := GetComponentAttr(c); p != nil {
				err = errors.Append(err, validatePathUUID(m, path, p))
			} else {
				err = errors.Append(err, errors.NewMissingFieldError(attrProdUUID))
			}
			if err != nil {
				cErrs = errors.Append(cErrs, errors.WrapIndex(err, c, i))
			}
		}
		if cErrs != nil {
			errs = errors.Append(errs, errors.Wrap(cErrs, obj.Components))
		}
	}
	return errs
}

func validatePathUUID(m *go3mf.Model, path string, p uuidPath) error {
	var errs error
	if p.getUUID() == "" {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrProdUUID))
	} else if uuid.Validate(string(p.getUUID())) != nil {
		errs = errors.Append(errs, ErrUUID)
	}
	if p.ObjectPath() != "" {
		if path == "" || path == m.PathOrDefault() { // root
			// Path is validated as part if the core validations
		} else {
			errs = errors.Append(errs, ErrProdRefInNonRoot)
		}
	}
	return errs
}
