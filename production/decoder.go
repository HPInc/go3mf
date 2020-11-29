package production

import (
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
	"github.com/qmuntal/go3mf/uuid"
)

func (e Spec) PostProcessDecode(m *go3mf.Model) {
	var (
		buildAttr *BuildAttr
		pu        *ItemAttr
	)
	if !m.Build.AnyAttr.Get(&buildAttr) {
		m.Build.AnyAttr = append(m.Build.AnyAttr, &BuildAttr{UUID: uuid.New()})
	}
	for _, item := range m.Build.Items {
		if !item.AnyAttr.Get(&pu) {
			item.AnyAttr = append(item.AnyAttr, &ItemAttr{
				UUID: uuid.New(),
			})
		} else if pu.UUID == "" {
			pu.UUID = uuid.New()
		}
	}
	e.fillResourceUUID(&m.Resources)
	for _, c := range m.Childs {
		e.fillResourceUUID(&c.Resources)
	}
	return
}

func (e Spec) fillResourceUUID(res *go3mf.Resources) {
	var (
		pu         *ComponentAttr
		objectAttr *ObjectAttr
	)
	for _, o := range res.Objects {
		if !o.AnyAttr.Get(&objectAttr) {
			o.AnyAttr = append(o.AnyAttr, &ObjectAttr{UUID: uuid.New()})
		}
		for _, c := range o.Components {
			if !c.AnyAttr.Get(&pu) {
				c.AnyAttr = append(c.AnyAttr, &ComponentAttr{
					UUID: uuid.New(),
				})
			} else if pu.UUID == "" {
				pu.UUID = uuid.New()
			}
		}
	}
}

func (e Spec) DecodeAttribute(parentNode interface{}, attr go3mf.XMLAttr) (err error) {
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
			var ext *ItemAttr
			if t.AnyAttr.Get(&ext) {
				ext.UUID = string(attr.Value)
			} else {
				t.AnyAttr = append(t.AnyAttr, &ItemAttr{UUID: string(attr.Value)})
			}
		case attrPath:
			var ext *ItemAttr
			if t.AnyAttr.Get(&ext) {
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
			var ext *ComponentAttr
			if t.AnyAttr.Get(&ext) {
				ext.UUID = string(attr.Value)
			} else {
				t.AnyAttr = append(t.AnyAttr, &ComponentAttr{UUID: string(attr.Value)})
			}
		case attrPath:
			var ext *ComponentAttr
			if t.AnyAttr.Get(&ext) {
				ext.Path = string(attr.Value)
			} else {
				t.AnyAttr = append(t.AnyAttr, &ComponentAttr{Path: string(attr.Value)})
			}
		}
	}
	return
}
