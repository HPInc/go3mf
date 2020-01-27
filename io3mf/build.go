package io3mf

import (
	"encoding/xml"

	go3mf "github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/iohelper"
)

type buildDecoder struct {
	iohelper.EmptyDecoder
}

func (d *buildDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrItem {
		child = &buildItemDecoder{}
	}
	return
}

func (d *buildDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space == nsProductionSpec && a.Name.Local == attrProdUUID {
			if err := validateUUID(a.Value); err != nil {
				ok = d.Scanner.InvalidRequiredAttr(attrProdUUID, a.Value)
			}
			d.Scanner.UUID = a.Value
			break
		}
	}

	if ok && d.Scanner.UUID == "" && d.Scanner.NamespaceRegistered(nsProductionSpec) {
		ok = d.Scanner.MissingAttr(attrProdUUID)
	}
	return ok
}

type buildItemDecoder struct {
	iohelper.EmptyDecoder
	item       go3mf.BuildItem
	objectID   uint32
	objectPath string
}

func (d *buildItemDecoder) Close() bool {
	return d.processItem() && d.Scanner.CloseResource()
}

func (d *buildItemDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrMetadataGroup {
		child = &metadataGroupDecoder{metadatas: &d.item.Metadata}
	}
	return
}

func (d *buildItemDecoder) processItem() bool {
	resource, ok := d.Scanner.FindResource(d.objectPath, uint32(d.objectID))
	if !ok {
		ok = d.Scanner.GenericError(true, "non-existent referenced object")
	} else if d.item.Object, ok = resource.(go3mf.Object); !ok {
		ok = d.Scanner.GenericError(true, "non-object referenced resource")
	}
	if ok {
		if d.item.Object != nil && d.item.Object.Type() == go3mf.ObjectTypeOther {
			ok = d.Scanner.GenericError(true, "referenced object cannot be have OTHER type")
		}
	}
	if ok {
		d.Scanner.BuildItems = append(d.Scanner.BuildItems, &d.item)
	}
	return ok
}

func (d *buildItemDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		switch a.Name.Space {
		case nsProductionSpec:
			if a.Name.Local == attrProdUUID {
				if err := validateUUID(a.Value); err != nil {
					ok = d.Scanner.InvalidRequiredAttr(attrProdUUID, a.Value)
				}
				d.item.UUID = a.Value
			} else if a.Name.Local == attrPath {
				d.objectPath = a.Value
			}
		case "":
			ok = d.parseCoreAttr(a)
		}
		if !ok {
			return false
		}
	}

	if d.item.UUID == "" && d.Scanner.NamespaceRegistered(nsProductionSpec) {
		ok = d.Scanner.MissingAttr(attrProdUUID)
	}
	return ok
}

func (d *buildItemDecoder) parseCoreAttr(a xml.Attr) bool {
	ok := true
	switch a.Name.Local {
	case attrObjectID:
		d.objectID, ok = d.Scanner.ParseResourceID(a.Value)
	case attrPartNumber:
		d.item.PartNumber = a.Value
	case attrTransform:
		d.item.Transform = d.Scanner.ParseToMatrixOptional(attrTransform, a.Value)
	}
	return ok
}
