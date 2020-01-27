package io3mf

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/iohelper"
)

type objectDecoder struct {
	iohelper.EmptyDecoder
	resource go3mf.ObjectResource
}

func (d *objectDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
}

func (d *objectDecoder) Close() bool {
	return d.Scanner.CloseResource()
}

func (d *objectDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		switch a.Name.Space {
		case nsProductionSpec:
			if a.Name.Local == attrProdUUID {
				if err := validateUUID(a.Value); err != nil {
					ok = d.Scanner.InvalidRequiredAttr(attrProdUUID, a.Value)
				} else {
					d.resource.UUID = a.Value
				}
			}
		case "":
			ok = d.parseCoreAttr(a)
		default:
			if ext, ok := extensionDecoder[a.Name.Space]; ok {
				ok = ext.DecodeAttribute(d.Scanner, &d.resource, a)
			}
		}
		if !ok {
			break
		}
	}
	return ok
}

func (d *objectDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec {
		if name.Local == attrMesh {
			child = &meshDecoder{resource: go3mf.MeshResource{ObjectResource: d.resource}}
		} else if name.Local == attrComponents {
			if d.resource.DefaultPropertyID != 0 {
				d.Scanner.GenericError(true, "default PID is not supported for component objects")
			}
			child = &componentsDecoder{resource: go3mf.ComponentsResource{ObjectResource: d.resource}}
		} else if name.Local == attrMetadataGroup {
			child = &metadataGroupDecoder{metadatas: &d.resource.Metadata}
		}
	}
	return
}

func (d *objectDecoder) parseCoreAttr(a xml.Attr) bool {
	ok := true
	switch a.Name.Local {
	case attrID:
		d.resource.ID, ok = d.Scanner.ParseResourceID(a.Value)
	case attrType:
		d.resource.ObjectType, ok = newObjectType(a.Value)
		if !ok {
			ok = true
			d.Scanner.InvalidOptionalAttr(attrType, a.Value)
		}
	case attrThumbnail:
		d.resource.Thumbnail = a.Value
	case attrName:
		d.resource.Name = a.Value
	case attrPartNumber:
		d.resource.PartNumber = a.Value
	case attrPID:
		d.resource.DefaultPropertyID = d.Scanner.ParseUint32Optional(attrPID, a.Value)
	case attrPIndex:
		d.resource.DefaultPropertyIndex = d.Scanner.ParseUint32Optional(attrPIndex, a.Value)
	}
	return ok
}

type componentsDecoder struct {
	iohelper.EmptyDecoder
	resource         go3mf.ComponentsResource
	componentDecoder componentDecoder
}

func (d *componentsDecoder) Open() {
	d.componentDecoder.resource = &d.resource
}
func (d *componentsDecoder) Close() bool {
	d.Scanner.AddResource(&d.resource)
	return true
}

func (d *componentsDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrComponent {
		child = &d.componentDecoder
	}
	return
}

type componentDecoder struct {
	iohelper.EmptyDecoder
	resource *go3mf.ComponentsResource
}

func (d *componentDecoder) Attributes(attrs []xml.Attr) bool {
	var (
		component go3mf.Component
		path      string
		objectID  uint32
	)
	ok := true
	for _, a := range attrs {
		switch a.Name.Space {
		case nsProductionSpec:
			if a.Name.Local == attrProdUUID {
				if err := validateUUID(a.Value); err != nil {
					ok = d.Scanner.InvalidRequiredAttr(attrProdUUID, a.Value)
				} else {
					component.UUID = a.Value
				}
			} else if a.Name.Local == attrPath {
				path = a.Value
			}
		case "":
			if a.Name.Local == attrObjectID {
				objectID, ok = d.Scanner.ParseUint32Required(attrObjectID, a.Value)
			} else if a.Name.Local == attrTransform {
				component.Transform = d.Scanner.ParseToMatrixOptional(attrTransform, a.Value)
			}
		}
		if !ok {
			return false
		}
	}
	return ok && d.addComponent(&component, path, objectID)
}

func (d *componentDecoder) addComponent(component *go3mf.Component, path string, objectID uint32) bool {
	ok := true
	if component.UUID == "" && d.Scanner.NamespaceRegistered(nsProductionSpec) {
		ok = d.Scanner.MissingAttr(attrProdUUID)
	}
	if ok && path != "" && !d.Scanner.IsRoot {
		ok = d.Scanner.GenericError(true, "path attribute in a non-root file is not supported")
	}
	if !ok {
		return false
	}

	resource, ok := d.Scanner.FindResource(path, uint32(objectID))
	if !ok {
		ok = d.Scanner.GenericError(true, "non-existent referenced object")
	} else if component.Object, ok = resource.(go3mf.Object); !ok {
		ok = d.Scanner.GenericError(true, "non-object referenced resource")
	}
	if ok {
		d.resource.Components = append(d.resource.Components, component)
	}
	return ok
}
