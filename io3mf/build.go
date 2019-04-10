package io3mf

import (
	"encoding/xml"

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

func (d *buildDecoder) Attributes(attrs []xml.Attr) bool {
	for _, a := range attrs {
		if a.Name.Space == nsProductionSpec && a.Name.Local == attrProdUUID {
			if _, err := uuid.FromString(a.Value); err != nil {
				return d.file.parser.InvalidRequiredAttr(attrProdUUID)
			}
			d.model.UUID = a.Value
		}
	}

	if d.model.UUID == "" && d.file.NamespaceRegistered(nsProductionSpec) {
		return d.file.parser.MissingAttr(attrProdUUID)
	}
	return true
}

type buildItemDecoder struct {
	emptyDecoder
	model      *go3mf.Model
	item       go3mf.BuildItem
	objectID   uint32
	objectPath string
}

func (d *buildItemDecoder) Close() bool {
	if d.objectID == 0 {
		return d.file.parser.CloseResource()
	}

	if d.processItem() {
		if d.item.Object.Type() == go3mf.ObjectTypeOther {
			if !d.file.parser.GenericError(true, "build item must not reference object of type OTHER") {
				return false
			}
		}

	}
	return true
}

func (d *buildItemDecoder) processItem() bool {
	resource, ok := d.file.FindResource(d.objectPath, uint32(d.objectID))
	if !ok {
		return d.file.parser.GenericError(true, "could not find build item object")
	}
	d.item.Object, ok = resource.(go3mf.Object)
	if !ok {
		return d.file.parser.GenericError(true, "a build item points to a non-object resource")
	}
	d.model.BuildItems = append(d.model.BuildItems, &d.item)
	return true
}

func (d *buildItemDecoder) Attributes(attrs []xml.Attr) bool {
	for _, a := range attrs {
		switch a.Name.Space {
		case nsProductionSpec:
			if a.Name.Local == attrProdUUID {
				if _, err := uuid.FromString(a.Value); err != nil {
					return d.file.parser.InvalidRequiredAttr(attrProdUUID)
				}
				d.item.UUID = a.Value
			} else if a.Name.Local == attrPath {
				d.objectPath = a.Value
			}
		case "":
			if !d.parseCoreAttr(a) {
				return false
			}
		}
	}

	if d.item.UUID == "" && d.file.NamespaceRegistered(nsProductionSpec) {
		return d.file.parser.MissingAttr(attrProdUUID)
	}
	return true
}

func (d *buildItemDecoder) parseCoreAttr(a xml.Attr) bool {
	ok := true
	switch a.Name.Local {
	case attrObjectID:
		d.objectID, ok = d.file.parser.ParseResourceID(a.Value)
	case attrPartNumber:
		d.item.PartNumber = a.Value
	case attrTransform:
		var err error
		d.item.Transform, err = strToMatrix(a.Value)
		if err != nil {
			d.file.parser.InvalidOptionalAttr(attrTransform)
		}
	}
	return ok
}
