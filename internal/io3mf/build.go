package io3mf

import (
	"encoding/xml"
	"errors"
	"github.com/go-gl/mathgl/mgl32"
	"strconv"

	"github.com/gofrs/uuid"
	mdl "github.com/qmuntal/go3mf/internal/model"
)

type buildDecoder struct {
	x     *xml.Decoder
	r     *Reader
	model *mdl.Model
}

func (d *buildDecoder) Decode(se xml.StartElement) error {
	if err := d.parseAttr(se.Attr); err != nil {
		return err
	}
	for {
		t, err := d.x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrItem {
				bd := buildItemDecoder{x: d.x, r: d.r, model: d.model}
				if err := bd.Decode(se); err != nil {
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
			if d.model.UUID != "" {
				return errors.New("go3mf: duplicated build uuid attribute")
			}
			if _, err := uuid.FromString(a.Value); err != nil {
				return errors.New("go3mf: build uuid is not valid")
			}
			d.model.UUID = a.Value
		}
	}

	if d.model.UUID == "" && d.model.RootPath == d.model.Path && d.r.namespaceRegistered(nsProductionSpec) {
		d.r.Warnings = append(d.r.Warnings, &ReadError{MissingMandatoryValue, "go3mf: a UUID for a build is missing"})
	}
	return nil
}

type buildItemDecoder struct {
	x           *xml.Decoder
	r           *Reader
	model       *mdl.Model
	item        *mdl.BuildItem
	objectID    uint64
	hasObjectID bool
	transform   mgl32.Mat4
}

func (d *buildItemDecoder) Decode(se xml.StartElement) error {
	d.item = new(mdl.BuildItem)
	if err := d.parseAttr(se.Attr); err != nil {
		return err
	}
	if !d.hasObjectID {
		return errors.New("go3mf: build item does not have objectid attribute")
	}

	return d.processItem()
}

func (d *buildItemDecoder) processItem() error {
	var path string
	if d.item.Path != "" {
		if d.model.Path != d.model.RootPath {
			return errors.New("go3mf: references in production extension go deeper than one level")
		}
		path = d.item.Path
	} else {
		path = d.model.Path
	}
	resource, ok := d.model.FindResource(d.objectID, path)
	if !ok {
		return errors.New("go3mf: could not find build item object")
	}
	obj, ok := resource.(*mdl.ObjectResource)
	if !ok {
		return errors.New("go3mf: could not find build item object")
	}
	if obj.Type() == mdl.ObjectTypeOther {
		d.r.Warnings = append(d.r.Warnings, &ReadError{InvalidMandatoryValue, "go3mf: build item must not reference object of type OTHER"})
	}
	d.item.Object = obj
	if !d.item.IsValidForSlices() {
		d.r.Warnings = append(d.r.Warnings, &ReadError{InvalidMandatoryValue, "go3mf: A slicestack posesses a nonplanar transformation"})
	}
	d.model.BuildItems = append(d.model.BuildItems, d.item)
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
				if d.item.Path != "" {
					return errors.New("go3mf: duplicated build item path attribute")
				}
				d.item.Path = a.Value
			}
		case "":
			if err := d.parseCoreAttr(a); err != nil {
				return err
			}
		}
	}

	if d.item.UUID == "" && d.model.RootPath == d.model.Path && d.r.namespaceRegistered(nsProductionSpec) {
		d.r.Warnings = append(d.r.Warnings, &ReadError{MissingMandatoryValue, "go3mf: a UUID for a build item is missing"})
	}
	return nil
}

func (d *buildItemDecoder) parseCoreAttr(a xml.Attr) (err error) {
	switch a.Name.Local {
	case attrObjectID:
		if d.hasObjectID {
			return errors.New("go3mf: duplicated build item objectid attribute")
		}
		if d.objectID, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
			return errors.New("go3mf: build item id is not valid")
		}
		d.hasObjectID = true
	case attrPartNumber:
		d.item.PartNumber = a.Value
	case attrTransform:
		d.item.Transform, err = strToMatrix(a.Value)
	}
	return
}
