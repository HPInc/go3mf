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
					ok = d.file.parser.InvalidRequiredAttr(attrProdUUID)
				}
				d.resource.UUID = a.Value
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
				d.file.parser.GenericError(true, "a components object must not have a default PID")
			}
			child = &componentsDecoder{resource: go3mf.ComponentsResource{ObjectResource: d.resource}}
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
			d.file.parser.InvalidOptionalAttr(attrType)
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
			d.file.parser.InvalidOptionalAttr(attrMeshRes)
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
					ok = d.file.parser.InvalidRequiredAttr(attrProdUUID)
				}
				component.UUID = a.Value
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
					d.file.parser.InvalidOptionalAttr(attrTransform)
				}
			}
		}
		if !ok {
			return false
		}
	}

	if component.UUID == "" && d.file.NamespaceRegistered(nsProductionSpec) {
		if !d.file.parser.MissingAttr(attrProdUUID) {
			return false
		}
	}

	if path != "" && !d.file.isRoot {
		if !d.file.parser.GenericError(true, "a component in a non-root model has a path attribute") {
			return false
		}
	}

	resource, ok := d.file.FindResource(path, uint32(objectID))
	if !ok {
		return d.file.parser.GenericError(true, "could not find component object")
	}
	component.Object, ok = resource.(go3mf.Object)
	if !ok {
		return d.file.parser.GenericError(true, "a component points to a non-object resource")
	}
	d.resource.Components = append(d.resource.Components, &component)
	return true
}
