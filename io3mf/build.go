package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	"github.com/gofrs/uuid"
	go3mf "github.com/qmuntal/go3mf"
)

type buildDecoder struct {
	emptyDecoder
	model *go3mf.Model
}

func (d *buildDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrItem {
		child = &buildItemDecoder{model: d.model}
	}
	return
}

func (d *buildDecoder) Attributes(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Space == nsProductionSpec && a.Name.Local == attrProdUUID {
			if d.model.UUID != "" {
				return errors.New("go3mf: duplicated build uuid attribute")
			}
			if _, err := uuid.FromString(a.Value); err != nil {
				return errors.New("go3mf: build uuid is not valid")
			}
			d.model.UUID = a.Value
		}
	}

	if d.model.UUID == "" && d.ModelFile().NamespaceRegistered(nsProductionSpec) {
		d.ModelFile().AddWarning(&ReadError{MissingMandatoryValue, "go3mf: a UUID for a build is missing"})
	}
	return nil
}

type buildItemDecoder struct {
	emptyDecoder
	model      *go3mf.Model
	item       go3mf.BuildItem
	objectID   uint64
	objectPath string
}

func (d *buildItemDecoder) Close() error {
	if d.objectID == 0 {
		return errors.New("go3mf: build item does not have objectid attribute")
	}

	return d.processItem()
}

func (d *buildItemDecoder) processItem() error {
	resource, ok := d.ModelFile().FindResource(d.objectPath, uint32(d.objectID))
	if !ok {
		return errors.New("go3mf: could not find build item object")
	}
	d.item.Object, ok = resource.(go3mf.Object)
	if !ok {
		return errors.New("go3mf: could not find build item object")
	}
	if d.item.Object.Type() == go3mf.ObjectTypeOther {
		d.ModelFile().AddWarning(&ReadError{InvalidMandatoryValue, "go3mf: build item must not reference object of type OTHER"})
	}
	if !d.item.IsValidForSlices() {
		d.ModelFile().AddWarning(&ReadError{InvalidMandatoryValue, "go3mf: A slicestack posesses a nonplanar transformation"})
	}
	d.model.BuildItems = append(d.model.BuildItems, &d.item)
	return nil
}

func (d *buildItemDecoder) Attributes(attrs []xml.Attr) error {
	for _, a := range attrs {
		switch a.Name.Space {
		case nsProductionSpec:
			if a.Name.Local == attrProdUUID {
				if d.item.UUID != "" {
					return errors.New("go3mf: duplicated build item uuid attribute")
				}
				if _, err := uuid.FromString(a.Value); err != nil {
					return errors.New("go3mf: build item uuid is not valid")
				}
				d.item.UUID = a.Value
			} else if a.Name.Local == attrPath {
				if d.objectPath != "" {
					return errors.New("go3mf: duplicated build item path attribute")
				}
				d.objectPath = a.Value
			}
		case "":
			if err := d.parseCoreAttr(a); err != nil {
				return err
			}
		}
	}

	if d.item.UUID == "" && d.ModelFile().NamespaceRegistered(nsProductionSpec) {
		d.ModelFile().AddWarning(&ReadError{MissingMandatoryValue, "go3mf: a UUID for a build item is missing"})
	}
	return nil
}

func (d *buildItemDecoder) parseCoreAttr(a xml.Attr) (err error) {
	switch a.Name.Local {
	case attrObjectID:
		if d.objectID != 0 {
			return errors.New("go3mf: duplicated build item objectid attribute")
		}
		if d.objectID, err = strconv.ParseUint(a.Value, 10, 32); err != nil {
			return errors.New("go3mf: build item id is not valid")
		}
	case attrPartNumber:
		d.item.PartNumber = a.Value
	case attrTransform:
		d.item.Transform, err = strToMatrix(a.Value)
	}
	return
}
