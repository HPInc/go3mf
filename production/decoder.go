package production

import (
	"encoding/xml"
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
	"github.com/qmuntal/go3mf/spec/encoding"
	"github.com/qmuntal/go3mf/uuid"
)

func (e Spec) NewElementDecoder(_ interface{}, _ xml.Name) encoding.ElementDecoder {
	return nil
}

func (e Spec) PostProcessDecode() {
	if GetBuildAttr(&e.m.Build) == nil {
		e.m.Build.AnyAttr = append(e.m.Build.AnyAttr, &BuildAttr{UUID: uuid.New()})
	}
	for _, item := range e.m.Build.Items {
		pu := GetItemAttr(item)
		if pu == nil {
			item.AnyAttr = append(item.AnyAttr, &ItemAttr{
				UUID: uuid.New(),
			})
		} else if pu.UUID == "" {
			pu.UUID = uuid.New()
		}
	}
	e.fillResourceUUID(&e.m.Resources)
	for _, c := range e.m.Childs {
		e.fillResourceUUID(&c.Resources)
	}
	return
}

func (e Spec) fillResourceUUID(res *go3mf.Resources) {
	for _, o := range res.Objects {
		if GetObjectAttr(o) == nil {
			o.AnyAttr = append(o.AnyAttr, &ObjectAttr{UUID: uuid.New()})
		}
		for _, c := range o.Components {
			pu := GetComponentAttr(c)
			if pu == nil {
				c.AnyAttr = append(c.AnyAttr, &ComponentAttr{
					UUID: uuid.New(),
				})
			} else if pu.UUID == "" {
				pu.UUID = uuid.New()
			}
		}
	}
}

func (e Spec) DecodeAttribute(parentNode interface{}, attr encoding.Attr) (err error) {
	switch t := parentNode.(type) {
	case *go3mf.Build:
		if attr.Name.Local == attrProdUUID {
			if err1 := uuid.Validate(string(attr.Value)); err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(attr.Name.Local, true))
			}
			t.AnyAttr = append(t.AnyAttr, &BuildAttr{UUID: string(attr.Value)})
		}
	case *go3mf.Item:
		switch attr.Name.Local {
		case attrProdUUID:
			if err1 := uuid.Validate(string(attr.Value)); err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(attr.Name.Local, true))
			}
			if ext := GetItemAttr(t); ext != nil {
				ext.UUID = string(attr.Value)
			} else {
				t.AnyAttr = append(t.AnyAttr, &ItemAttr{UUID: string(attr.Value)})
			}
		case attrPath:
			if ext := GetItemAttr(t); ext != nil {
				ext.Path = string(attr.Value)
			} else {
				t.AnyAttr = append(t.AnyAttr, &ItemAttr{Path: string(attr.Value)})
			}
		}
	case *go3mf.Object:
		if attr.Name.Local == attrProdUUID {
			if err1 := uuid.Validate(string(attr.Value)); err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(attr.Name.Local, true))
			}
			t.AnyAttr = append(t.AnyAttr, &ObjectAttr{UUID: string(attr.Value)})
		}
	case *go3mf.Component:
		switch attr.Name.Local {
		case attrProdUUID:
			if err1 := uuid.Validate(string(attr.Value)); err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(attr.Name.Local, true))
			}
			if ext := GetComponentAttr(t); ext != nil {
				ext.UUID = string(attr.Value)
			} else {
				t.AnyAttr = append(t.AnyAttr, &ComponentAttr{UUID: string(attr.Value)})
			}
		case attrPath:
			if ext := GetComponentAttr(t); ext != nil {
				ext.Path = string(attr.Value)
			} else {
				t.AnyAttr = append(t.AnyAttr, &ComponentAttr{Path: string(attr.Value)})
			}
		}
	}
	return
}
