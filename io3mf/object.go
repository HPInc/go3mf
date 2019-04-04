package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	"github.com/gofrs/uuid"
	go3mf "github.com/qmuntal/go3mf"
)

type objectDecoder struct {
	emptyDecoder
	r             *Reader
	progressCount int
	resource      go3mf.ObjectResource
}

func (d *objectDecoder) Open() error {
	if !d.r.progress.progress(1.0-2.0/float64(d.progressCount+2), StageReadResources) {
		return ErrUserAborted
	}
	d.r.progress.pushLevel(1.0-2.0/float64(d.progressCount+2), 1.0-2.0/float64(d.progressCount+1+2))
	return nil
}

func (d *objectDecoder) Close() error {
	d.r.progress.popLevel()
	return nil
}

func (d *objectDecoder) Attributes(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		switch a.Name.Space {
		case nsProductionSpec:
			if a.Name.Local == attrProdUUID {
				if d.resource.UUID != "" {
					d.r.addWarning(&ReadError{InvalidMandatoryValue, "go3mf: duplicated object resource uuid attribute"})
				}
				if _, err = uuid.FromString(a.Value); err != nil {
					err = errors.New("go3mf: object resource uuid is not valid")
				} else {
					d.resource.UUID = a.Value
				}
			}
		case nsSliceSpec:
			err = d.parseSliceAttr(a)
		case "":
			err = d.parseCoreAttr(a)
		}
		if err != nil {
			break
		}
	}
	return
}

func (d *objectDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec {
		if name.Local == attrMesh {
			child = &meshDecoder{r: d.r, resource: go3mf.MeshResource{ObjectResource: d.resource}}
		} else if name.Local == attrComponents {
			if d.resource.DefaultPropertyID != 0 {
				d.r.addWarning(&ReadError{InvalidOptionalValue, "go3mf: a components object must not have a default PID"})
			}
			child = &componentsDecoder{r: d.r, resource: go3mf.ComponentsResource{ObjectResource: d.resource}}
		}
	}
	return
}

func (d *objectDecoder) parseCoreAttr(a xml.Attr) (err error) {
	switch a.Name.Local {
	case attrID:
		if d.resource.ID != 0 {
			err = errors.New("go3mf: duplicated object resource id attribute")
		} else {
			d.resource.ID, err = strconv.ParseUint(a.Value, 10, 64)
			if err != nil {
				err = errors.New("go3mf: object resource id is not valid")
			}
		}
	case attrType:
		var ok bool
		d.resource.ObjectType, ok = newObjectType(a.Value)
		if !ok {
			d.r.addWarning(&ReadError{InvalidOptionalValue, "go3mf: object resource type is not valid"})
		}
	case attrThumbnail:
		d.resource.Thumbnail = a.Value
	case attrName:
		d.resource.Name = a.Value
	case attrPartNumber:
		d.resource.PartNumber = a.Value
	case attrPID:
		d.resource.DefaultPropertyID, err = strconv.ParseUint(a.Value, 10, 64)
		if err != nil {
			err = errors.New("go3mf: object resource pid is not valid")
		}
	case attrPIndex:
		d.resource.DefaultPropertyIndex, err = strconv.ParseUint(a.Value, 10, 64)
		if err != nil {
			err = errors.New("go3mf: object resource ºpindex is not valid")
		}
	}
	return
}

func (d *objectDecoder) parseSliceAttr(a xml.Attr) (err error) {
	switch a.Name.Local {
	case attrSliceRefID:
		if d.resource.SliceStackID != 0 {
			d.r.addWarning(&ReadError{InvalidOptionalValue, "go3mf: duplicated object resource slicestackid attribute"})
		}
		d.resource.SliceStackID, err = strconv.ParseUint(a.Value, 10, 64)
		if err != nil {
			err = errors.New("go3mf: object resource slicestackid is not valid")
		}
	case attrMeshRes:
		var ok bool
		d.resource.SliceResoultion, ok = newSliceResolution(a.Value)
		if !ok {
			err = errors.New("go3mf: object resource sliceresolution is not valid")
		}
	}
	return
}

type componentsDecoder struct {
	emptyDecoder
	r                *Reader
	resource         go3mf.ComponentsResource
	componentDecoder componentDecoder
}

func (d *componentsDecoder) Open() error {
	d.componentDecoder.r = d.r
	d.componentDecoder.resource = &d.resource
	return nil
}
func (d *componentsDecoder) Close() error {
	d.r.addResource(&d.resource)
	return nil
}

func (d *componentsDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrComponent {
		child = &d.componentDecoder
	}
	return
}

type componentDecoder struct {
	emptyDecoder
	r        *Reader
	resource *go3mf.ComponentsResource
}

func (d *componentDecoder) Attributes(attrs []xml.Attr) (err error) {
	var component go3mf.Component
	var path string
	var objectID uint64
	for _, a := range attrs {
		switch a.Name.Space {
		case nsProductionSpec:
			if a.Name.Local == attrProdUUID {
				if component.UUID != "" {
					d.r.addWarning(&ReadError{InvalidMandatoryValue, "go3mf: duplicated component uuid attribute"})
				}
				if _, err = uuid.FromString(a.Value); err != nil {
					err = errors.New("go3mf: component uuid is not valid")
				} else {
					component.UUID = a.Value
				}
			} else if a.Name.Local == attrPath {
				if path != "" {
					d.r.addWarning(&ReadError{InvalidMandatoryValue, "go3mf: duplicated component path attribute"})
				}
				path = a.Value
			}
		case "":
			if a.Name.Local == attrObjectID {
				if objectID != 0 {
					err = errors.New("go3mf: duplicated component objectid attribute")
				}
				objectID, err = strconv.ParseUint(a.Value, 10, 64)
				if err != nil {
					err = errors.New("go3mf: component id is not valid")
				}
			} else if a.Name.Local == attrTransform {
				component.Transform, err = strToMatrix(a.Value)
			}
		}
		if err != nil {
			break
		}
	}
	if component.UUID == "" && d.r.namespaceRegistered(nsProductionSpec) {
		d.r.addWarning(&ReadError{MissingMandatoryValue, "go3mf: a UUID for a component is missing"})
	}
	resource, ok := d.r.Model.FindResource(path, objectID)
	if !ok {
		err = errors.New("go3mf: could not find component object")
	}
	component.Object, ok = resource.(go3mf.Object)
	if !ok {
		return errors.New("go3mf: could not find component object")
	}
	d.resource.Components = append(d.resource.Components, &component)
	return
}
