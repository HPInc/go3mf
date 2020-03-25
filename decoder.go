package go3mf

import (
	"encoding/xml"
	"image/color"
	"strconv"
	"strings"
)

type modelDecoder struct {
	baseDecoder
	model *Model
}

func (d *modelDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace {
		switch name.Local {
		case attrResources:
			child = &resourceDecoder{}
		case attrBuild:
			if d.Scanner.IsRoot {
				child = &buildDecoder{build: &d.model.Build}
			}
		case attrMetadata:
			if d.Scanner.IsRoot {
				child = &metadataDecoder{metadatas: &d.model.Metadata}
			}
		}
	} else if ext, ok := d.Scanner.extensionDecoder[name.Space]; ok {
		child = ext.NewNodeDecoder(d.model, name.Local)
	}
	return
}

func (d *modelDecoder) Start(attrs []xml.Attr) {
	if !d.Scanner.IsRoot {
		return
	}
	var requiredExts []string
	for _, a := range attrs {
		if a.Name.Space == "" {
			switch a.Name.Local {
			case attrUnit:
				var ok bool
				if d.model.Units, ok = newUnits(a.Value); !ok {
					d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
				}
			case attrThumbnail:
				d.model.Thumbnail = a.Value
			case attrReqExt:
				requiredExts = strings.Fields(a.Value)
			}
		} else {
			d.noCoreAttribute(a)
		}
	}

	for _, ext := range requiredExts {
		for _, x := range d.model.Specs {
			if x.Local() == ext {
				x.SetRequired(true)
				break
			}
		}
	}
}

func (d *modelDecoder) noCoreAttribute(a xml.Attr) {
	switch a.Name.Space {
	case nsXML:
		if a.Name.Local == attrLang {
			d.model.Language = a.Value
		}
	case attrXmlns:
		if ext, ok := d.model.Specs[a.Value]; ok {
			ext.SetLocal(a.Name.Local)
		} else {
			d.model.WithExtension(&UnknownSpec{SpaceName: a.Value, LocalName: a.Name.Local})
		}
	default:
		if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, d.model, a)
		}
	}
}

type metadataGroupDecoder struct {
	baseDecoder
	metadatas *[]Metadata
}

func (d *metadataGroupDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace && name.Local == attrMetadata {
		child = &metadataDecoder{metadatas: d.metadatas}
	}
	return
}

type metadataDecoder struct {
	baseDecoder
	metadatas *[]Metadata
	metadata  Metadata
}

func (d *metadataDecoder) Start(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrName:
			d.metadata.Name = a.Name
			i := strings.IndexByte(a.Value, ':')
			if i < 0 {
				d.metadata.Name.Local = a.Value
			} else if _, ok := d.Scanner.namespace(a.Value[0:i]); ok {
				d.metadata.Name.Space = a.Value[0:i]
				d.metadata.Name.Local = a.Value[i+1:]
			} else {
				d.metadata.Name.Local = a.Value
			}
		case attrType:
			d.metadata.Type = a.Value
		case attrPreserve:
			d.metadata.Preserve, _ = strconv.ParseBool(a.Value)
		}
	}
}

func (d *metadataDecoder) Text(txt []byte) {
	d.metadata.Value = string(txt)
}

func (d *metadataDecoder) End() {
	*d.metadatas = append(*d.metadatas, d.metadata)
}

type buildDecoder struct {
	baseDecoder
	build *Build
}

func (d *buildDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace && name.Local == attrItem {
		child = &buildItemDecoder{}
	}
	return
}

func (d *buildDecoder) Start(attrs []xml.Attr) {
	for _, a := range attrs {
		if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, d.build, a)
		}
	}
}

type buildItemDecoder struct {
	baseDecoder
	item Item
}

func (d *buildItemDecoder) End() {
	d.Scanner.BuildItems = append(d.Scanner.BuildItems, &d.item)
	d.Scanner.ResourceID = 0
}

func (d *buildItemDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace && name.Local == attrMetadataGroup {
		child = &metadataGroupDecoder{metadatas: &d.item.Metadata}
	}
	return
}

func (d *buildItemDecoder) Start(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space == "" {
			d.parseCoreAttr(a)
		} else if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, &d.item, a)
		}
	}
	return
}

func (d *buildItemDecoder) parseCoreAttr(a xml.Attr) {
	switch a.Name.Local {
	case attrObjectID:
		val, err := strconv.ParseUint(a.Value, 10, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
		}
		d.item.ObjectID = uint32(val)
		d.Scanner.ResourceID = d.item.ObjectID
	case attrPartNumber:
		d.item.PartNumber = a.Value
	case attrTransform:
		var ok bool
		d.item.Transform, ok = ParseMatrix(a.Value)
		if !ok {
			d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
		}
	}
}

type resourceDecoder struct {
	baseDecoder
}

func (d *resourceDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace {
		switch name.Local {
		case attrObject:
			child = &objectDecoder{}
		case attrBaseMaterials:
			child = new(baseMaterialsDecoder)
		}
	} else if ext, ok := d.Scanner.extensionDecoder[name.Space]; ok {
		child = ext.NewNodeDecoder(nil, name.Local)
	}
	return
}

type baseMaterialsDecoder struct {
	baseDecoder
	resource            BaseMaterials
	baseMaterialDecoder baseMaterialDecoder
}

func (d *baseMaterialsDecoder) End() {
	d.Scanner.AddAsset(&d.resource)
}

func (d *baseMaterialsDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace && name.Local == attrBase {
		child = &d.baseMaterialDecoder
	}
	return
}

func (d *baseMaterialsDecoder) Start(attrs []xml.Attr) {
	d.baseMaterialDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			id, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			d.resource.ID, d.Scanner.ResourceID = uint32(id), uint32(id)
			break
		}
	}
}

type baseMaterialDecoder struct {
	baseDecoder
	resource *BaseMaterials
}

func (d *baseMaterialDecoder) Start(attrs []xml.Attr) {
	var name string
	var baseColor color.RGBA
	for _, a := range attrs {
		switch a.Name.Local {
		case attrName:
			name = a.Value
		case attrDisplayColor:
			var err error
			baseColor, err = ParseRGBA(a.Value)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
		}
	}
	d.resource.Materials = append(d.resource.Materials, Base{Name: name, Color: baseColor})
	return
}

type meshDecoder struct {
	baseDecoder
	resource *Object
}

func (d *meshDecoder) Start(_ []xml.Attr) {
	d.resource.Mesh = new(Mesh)
}

func (d *meshDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace {
		if name.Local == attrVertices {
			child = &verticesDecoder{mesh: d.resource.Mesh}
		} else if name.Local == attrTriangles {
			child = &trianglesDecoder{resource: d.resource}
		}
	} else if ext, ok := d.Scanner.extensionDecoder[name.Space]; ok {
		child = ext.NewNodeDecoder(d.resource.Mesh, name.Local)
	}
	return
}

type verticesDecoder struct {
	baseDecoder
	mesh          *Mesh
	vertexDecoder vertexDecoder
}

func (d *verticesDecoder) Start(_ []xml.Attr) {
	d.vertexDecoder.mesh = d.mesh
}

func (d *verticesDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace && name.Local == attrVertex {
		child = &d.vertexDecoder
	}
	return
}

type vertexDecoder struct {
	baseDecoder
	mesh *Mesh
}

func (d *vertexDecoder) Start(attrs []xml.Attr) {
	var x, y, z float32
	for _, a := range attrs {
		val, err := strconv.ParseFloat(a.Value, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
		}
		switch a.Name.Local {
		case attrX:
			x = float32(val)
		case attrY:
			y = float32(val)
		case attrZ:
			z = float32(val)
		}
	}
	d.mesh.Vertices = append(d.mesh.Vertices, Point3D{x, y, z})
}

type trianglesDecoder struct {
	baseDecoder
	resource        *Object
	triangleDecoder triangleDecoder
}

func (d *trianglesDecoder) Start(_ []xml.Attr) {
	d.triangleDecoder.mesh = d.resource.Mesh
	d.triangleDecoder.defaultPropertyID = d.resource.DefaultPID
	d.triangleDecoder.defaultPropertyIndex = d.resource.DefaultPIndex

	if len(d.resource.Mesh.Triangles) == 0 && len(d.resource.Mesh.Vertices) > 0 {
		d.resource.Mesh.Triangles = make([]Triangle, 0, len(d.resource.Mesh.Vertices)-1)
	}
}

func (d *trianglesDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace && name.Local == attrTriangle {
		child = &d.triangleDecoder
	}
	return
}

type triangleDecoder struct {
	baseDecoder
	mesh                                    *Mesh
	defaultPropertyIndex, defaultPropertyID uint32
}

func (d *triangleDecoder) Start(attrs []xml.Attr) {
	var v1, v2, v3, pid, p1, p2, p3 uint32
	var hasPID, hasP1, hasP2, hasP3 bool
	for _, a := range attrs {
		required := true
		val, err := strconv.ParseUint(a.Value, 10, 32)
		switch a.Name.Local {
		case attrV1:
			v1 = uint32(val)
		case attrV2:
			v2 = uint32(val)
		case attrV3:
			v3 = uint32(val)
		case attrPID:
			pid = uint32(val)
			hasPID = true
			required = false
		case attrP1:
			p1 = uint32(val)
			hasP1 = true
			required = false
		case attrP2:
			p2 = uint32(val)
			hasP2 = true
			required = false
		case attrP3:
			p3 = uint32(val)
			hasP3 = true
			required = false
		}
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, a.Value, required)
		}
	}

	p1 = applyDefault(p1, d.defaultPropertyIndex, hasP1)
	p2 = applyDefault(p2, p1, hasP2)
	p3 = applyDefault(p3, p1, hasP3)
	pid = applyDefault(pid, d.defaultPropertyID, hasPID)

	d.mesh.Triangles = append(d.mesh.Triangles, Triangle{
		Indices:  [3]uint32{v1, v2, v3},
		PID:      pid,
		PIndices: [3]uint32{p1, p2, p3},
	})
}

func applyDefault(val, defVal uint32, noDef bool) uint32 {
	if noDef {
		return val
	}
	return defVal
}

type objectDecoder struct {
	baseDecoder
	resource Object
}

func (d *objectDecoder) End() {
	d.Scanner.AddObject(&d.resource)
}

func (d *objectDecoder) Start(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space == "" {
			d.parseCoreAttr(a)
		} else if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, &d.resource, a)
		}
	}
}

func (d *objectDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace {
		if name.Local == attrMesh {
			child = &meshDecoder{resource: &d.resource}
		} else if name.Local == attrComponents {
			child = &componentsDecoder{resource: &d.resource}
		} else if name.Local == attrMetadataGroup {
			child = &metadataGroupDecoder{metadatas: &d.resource.Metadata}
		}
	}
	return
}

func (d *objectDecoder) parseCoreAttr(a xml.Attr) {
	switch a.Name.Local {
	case attrID:
		id, err := strconv.ParseUint(a.Value, 10, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
		}
		d.resource.ID, d.Scanner.ResourceID = uint32(id), uint32(id)
	case attrType:
		var ok bool
		d.resource.ObjectType, ok = newObjectType(a.Value)
		if !ok {
			d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
		}
	case attrThumbnail:
		d.resource.Thumbnail = a.Value
	case attrName:
		d.resource.Name = a.Value
	case attrPartNumber:
		d.resource.PartNumber = a.Value
	case attrPID:
		val, err := strconv.ParseUint(a.Value, 10, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
		}
		d.resource.DefaultPID = uint32(val)
	case attrPIndex:
		val, err := strconv.ParseUint(a.Value, 10, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
		}
		d.resource.DefaultPIndex = uint32(val)
	}
}

type componentsDecoder struct {
	baseDecoder
	resource         *Object
	componentDecoder componentDecoder
}

func (d *componentsDecoder) Start(_ []xml.Attr) {
	d.resource.Components = make([]*Component, 0)
	d.componentDecoder.resource = d.resource
}

func (d *componentsDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionSpace && name.Local == attrComponent {
		child = &d.componentDecoder
	}
	return
}

type componentDecoder struct {
	baseDecoder
	resource *Object
}

func (d *componentDecoder) Start(attrs []xml.Attr) {
	var component Component
	for _, a := range attrs {
		if a.Name.Space == "" {
			if a.Name.Local == attrObjectID {
				val, err := strconv.ParseUint(a.Value, 10, 32)
				if err != nil {
					d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
				}
				component.ObjectID = uint32(val)
			} else if a.Name.Local == attrTransform {
				var ok bool
				component.Transform, ok = ParseMatrix(a.Value)
				if !ok {
					d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
				}
			}
		} else if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, &component, a)
		}
	}
	d.resource.Components = append(d.resource.Components, &component)
}
