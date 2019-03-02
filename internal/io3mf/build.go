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
	if err := d.parseAttr(se); err != nil {
		return err
	}
	for {
		t, err := d.x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if d.r.namespaceAttr(tp.Name.Space) == nsCoreSpec {
				if tp.Name.Local == "item" {
					bd := buildItemDecoder{x: d.x, r: d.r, model: d.model}
					if err := bd.Decode(se); err != nil {
						return err
					}
				}
			}
		}
	}
}

func (d *buildDecoder) parseAttr(se xml.StartElement) error {
	for _, a := range se.Attr {
		if d.r.namespaceAttr(a.Name.Space) == nsProductionSpec {
			if se.Name.Local == attrProdUUID {
				if d.model.UUID() != uuid.Nil {
					return errors.New("go3mf: duplicated build uuid attribute")
				}
				id := uuid.FromStringOrNil(a.Value)
				return d.model.SetUUID(id)
			}
		}
	}

	if d.model.UUID() == uuid.Nil && d.model.RootPath == d.model.Path && d.r.namespaceRegistered(nsCoreSpec) {
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
	if err := d.parseAttr(se); err != nil {
		return err
	}
	return nil
}

func (d *buildItemDecoder) parseAttr(se xml.StartElement) error {
	for _, a := range se.Attr {
		switch d.r.namespaceAttr(a.Name.Space) {
		case nsProductionSpec:
			if se.Name.Local == attrProdUUID {
				if d.item.UUID() != uuid.Nil {
					return errors.New("go3mf: duplicated build item uuid attribute")
				}
				id := uuid.FromStringOrNil(a.Value)
				if err := d.item.SetUUID(id); err != nil {
					return err
				}
			} else if a.Name.Local == attrProdPath {
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

	if d.model.UUID() == uuid.Nil && d.model.RootPath == d.model.Path && d.r.namespaceRegistered(nsCoreSpec) {
		d.r.Warnings = append(d.r.Warnings, &ReadError{MissingMandatoryValue, "go3mf: a UUID for a build is missing"})
	}
	return nil
}

func (d *buildItemDecoder) parseCoreAttr(a xml.Attr) error {
	switch a.Name.Local {
	case attrObjectID:
		if d.hasObjectID {
			return errors.New("go3mf: duplicated build item objectid attribute")
		}
		var err error
		if d.objectID, err = strconv.ParseUint(a.Value, 10, 64); err != nil {
			return errors.New("go3mf: build item id is not valid")
		}
		d.hasObjectID = true
	}
	return nil
}
