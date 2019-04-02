package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	"github.com/gofrs/uuid"
	go3mf "github.com/qmuntal/go3mf"
)

type buildDecoder struct {
	r *Reader
}

func (d *buildDecoder) Decode(x xml.TokenReader, attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrItem {
				bd := buildItemDecoder{r: d.r}
				if err := bd.Decode(tp.Attr); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrBuild {
				return nil
			}
		}
	}
}

func (d *buildDecoder) parseAttr(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Space == nsProductionSpec && a.Name.Local == attrProdUUID {
			if d.r.Model.UUID != "" {
				return errors.New("go3mf: duplicated build uuid attribute")
			}
			if _, err := uuid.FromString(a.Value); err != nil {
				return errors.New("go3mf: build uuid is not valid")
			}
			d.r.Model.UUID = a.Value
		}
	}

	if d.r.Model.UUID == "" && d.r.namespaceRegistered(nsProductionSpec) {
		d.r.addWarning(&ReadError{MissingMandatoryValue, "go3mf: a UUID for a build is missing"})
	}
	return nil
}

type buildItemDecoder struct {
	r          *Reader
	item       go3mf.BuildItem
	objectID   uint64
	objectPath string
}

func (d *buildItemDecoder) Decode(attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	if d.objectID == 0 {
		return errors.New("go3mf: build item does not have objectid attribute")
	}

	return d.processItem()
}

func (d *buildItemDecoder) processItem() error {
	resource, ok := d.r.Model.FindResource(d.objectPath, d.objectID)
	if !ok {
		return errors.New("go3mf: could not find build item object")
	}
	d.item.Object, ok = resource.(go3mf.Object)
	if !ok {
		return errors.New("go3mf: could not find build item object")
	}
	if d.item.Object.Type() == go3mf.ObjectTypeOther {
		d.r.addWarning(&ReadError{InvalidMandatoryValue, "go3mf: build item must not reference object of type OTHER"})
	}
	if !d.item.IsValidForSlices() {
		d.r.addWarning(&ReadError{InvalidMandatoryValue, "go3mf: A slicestack posesses a nonplanar transformation"})
	}
	d.r.Model.BuildItems = append(d.r.Model.BuildItems, &d.item)
	return nil
}

func (d *buildItemDecoder) parseAttr(attrs []xml.Attr) error {
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

	if d.item.UUID == "" && d.r.namespaceRegistered(nsProductionSpec) {
		d.r.addWarning(&ReadError{MissingMandatoryValue, "go3mf: a UUID for a build item is missing"})
	}
	return nil
}

func (d *buildItemDecoder) parseCoreAttr(a xml.Attr) (err error) {
	switch a.Name.Local {
	case attrObjectID:
		if d.objectID != 0 {
			return errors.New("go3mf: duplicated build item objectid attribute")
		}
		if d.objectID, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
			return errors.New("go3mf: build item id is not valid")
		}
	case attrPartNumber:
		d.item.PartNumber = a.Value
	case attrTransform:
		d.item.Transform, err = strToMatrix(a.Value)
	}
	return
}
