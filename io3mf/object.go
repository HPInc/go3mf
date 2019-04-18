package io3mf

import (
	"encoding/xml"

	"github.com/gofrs/uuid"
	go3mf "github.com/qmuntal/go3mf"
)

type objectDecoder struct {
	emptyDecoder
	progressCount int
	resource      go3mf.ObjectResource
}

func (d *objectDecoder) Open() {
	d.resource.ModelPath = d.file.path
}

func (d *objectDecoder) Close() bool {
	return d.file.parser.CloseResource()
}

func (d *objectDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		switch a.Name.Space {
		case nsProductionSpec:
			if a.Name.Local == attrProdUUID {
				if _, err := uuid.FromString(a.Value); err != nil {
					ok = d.file.parser.InvalidRequiredAttr(attrProdUUID, a.Value)
				} else {
					d.resource.UUID = a.Value
				}
			}
		case nsSliceSpec:
			ok = d.parseSliceAttr(a)
		case "":
			ok = d.parseCoreAttr(a)
		}
		if !ok {
			break
		}
	}
	return ok
}

func (d *objectDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec {
		if name.Local == attrMesh {
			child = &meshDecoder{resource: go3mf.MeshResource{ObjectResource: d.resource}}
		} else if name.Local == attrComponents {
			if d.resource.DefaultPropertyID != 0 {
				d.file.parser.GenericError(true, "default PID is not supported for component objects")
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
		d.resource.ID, ok = d.file.parser.ParseResourceID(a.Value)
	case attrType:
		d.resource.ObjectType, ok = newObjectType(a.Value)
		if !ok {
			ok = true
			d.file.parser.InvalidOptionalAttr(attrType, a.Value)
		}
	case attrThumbnail:
		d.resource.Thumbnail = a.Value
	case attrName:
		d.resource.Name = a.Value
	case attrPartNumber:
		d.resource.PartNumber = a.Value
	case attrPID:
		d.resource.DefaultPropertyID = d.file.parser.ParseUint32Optional(attrPID, a.Value)
	case attrPIndex:
		d.resource.DefaultPropertyIndex = d.file.parser.ParseUint32Optional(attrPIndex, a.Value)
	}
	return ok
}

func (d *objectDecoder) parseSliceAttr(a xml.Attr) bool {
	ok := true
	switch a.Name.Local {
	case attrSliceRefID:
		d.resource.SliceStackID, ok = d.file.parser.ParseUint32Required(attrSliceRefID, a.Value)
	case attrMeshRes:
		d.resource.SliceResoultion, ok = newSliceResolution(a.Value)
		if !ok {
			ok = true
			d.file.parser.InvalidOptionalAttr(attrMeshRes, a.Value)
		}
	}
	return ok
}

type componentsDecoder struct {
	emptyDecoder
	resource         go3mf.ComponentsResource
	componentDecoder componentDecoder
}

func (d *componentsDecoder) Open() {
	d.componentDecoder.resource = &d.resource
}
func (d *componentsDecoder) Close() bool {
	d.file.AddResource(&d.resource)
	return true
}

func (d *componentsDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrComponent {
		child = &d.componentDecoder
	}
	return
}

type componentDecoder struct {
	emptyDecoder
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
				if _, err := uuid.FromString(a.Value); err != nil {
					ok = d.file.parser.InvalidRequiredAttr(attrProdUUID, a.Value)
				} else {
					component.UUID = a.Value
				}
			} else if a.Name.Local == attrPath {
				path = a.Value
			}
		case "":
			if a.Name.Local == attrObjectID {
				objectID, ok = d.file.parser.ParseUint32Required(attrObjectID, a.Value)
			} else if a.Name.Local == attrTransform {
				var err error
				component.Transform, err = strToMatrix(a.Value)
				if err != nil {
					d.file.parser.InvalidOptionalAttr(attrTransform, a.Value)
				}
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
	if component.UUID == "" && d.file.NamespaceRegistered(nsProductionSpec) {
		ok = d.file.parser.MissingAttr(attrProdUUID)
	}
	if ok && path != "" && !d.file.isRoot {
		ok = d.file.parser.GenericError(true, "path attribute in a non-root file is not supported")
	}
	if !ok {
		return false
	}

	resource, ok := d.file.FindResource(path, uint32(objectID))
	if !ok {
		ok = d.file.parser.GenericError(true, "non-existent referenced object")
	} else if component.Object, ok = resource.(go3mf.Object); !ok {
		ok = d.file.parser.GenericError(true, "non-object referenced resource")
	}
	if ok {
		d.resource.Components = append(d.resource.Components, component)
	}
	return ok
}
