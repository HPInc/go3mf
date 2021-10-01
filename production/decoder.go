// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package production

import (
	specerr "github.com/hpinc/go3mf/errors"
	"github.com/hpinc/go3mf/spec"
	"github.com/hpinc/go3mf/uuid"
)

func (Spec) NewElementDecoder(_ interface{}, _ string) spec.ElementDecoder {
	return nil
}

func (Spec) NewAttr3MF(parent string) spec.AttrGroup {
	switch parent {
	case "build":
		return new(BuildAttr)
	case "item":
		return new(ItemAttr)
	case "object":
		return new(ObjectAttr)
	case "component":
		return new(ComponentAttr)
	}
	return nil
}

func (u *BuildAttr) Unmarshal3MFAttr(a spec.XMLAttr) error {
	if a.Name.Local == attrProdUUID {
		if err := uuid.Validate(string(a.Value)); err != nil {
			return specerr.NewParseAttrError(a.Name.Local, true)
		}
		u.UUID = string(a.Value)
	}
	return nil
}

func (u *ItemAttr) Unmarshal3MFAttr(a spec.XMLAttr) error {
	if a.Name.Local == attrProdUUID {
		if err := uuid.Validate(string(a.Value)); err != nil {
			return specerr.NewParseAttrError(a.Name.Local, true)
		}
		u.UUID = string(a.Value)
	} else if a.Name.Local == attrPath {
		u.Path = string(a.Value)
	}
	return nil
}

func (u *ObjectAttr) Unmarshal3MFAttr(a spec.XMLAttr) error {
	if a.Name.Local == attrProdUUID {
		if err := uuid.Validate(string(a.Value)); err != nil {
			return specerr.NewParseAttrError(a.Name.Local, true)
		}
		u.UUID = string(a.Value)
	}
	return nil
}

func (u *ComponentAttr) Unmarshal3MFAttr(a spec.XMLAttr) error {
	if a.Name.Local == attrProdUUID {
		if err := uuid.Validate(string(a.Value)); err != nil {
			return specerr.NewParseAttrError(a.Name.Local, true)
		}
		u.UUID = string(a.Value)
	} else if a.Name.Local == attrPath {
		u.Path = string(a.Value)
	}
	return nil
}
