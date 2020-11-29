package go3mf

import (
	"bytes"
	"encoding/xml"
	"image/color"
	"strconv"
	"strings"
	"unsafe"
)

type modelDecoder struct {
	baseDecoder
	model  *Model
	isRoot bool
}

func (d *modelDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == Namespace {
		switch name.Local {
		case attrResources:
			resources, _ := d.model.FindResources(d.Scanner.modelPath)
			child = &resourceDecoder{
				resources: resources,
			}
		case attrBuild:
			if d.isRoot {
				child = &buildDecoder{build: &d.model.Build}
			}
		case attrMetadata:
			if d.isRoot {
				child = &metadataDecoder{metadatas: &d.model.Metadata}
			}
		}
	} else if ext, ok := d.Scanner.extensionDecoder[name.Space]; ok {
		if ext, ok := ext.(modelElementDecoder); ok {
			child = ext.NewModelElementDecoder(d.model, name.Local)
		}
	}
	return
}

func (d *modelDecoder) Start(attrs []XMLAttr) {
	if !d.isRoot {
		return
	}
	var requiredExts []string
	for _, a := range attrs {
		if a.Name.Space == "" {
			switch a.Name.Local {
			case attrUnit:
				var ok bool
				if d.model.Units, ok = newUnits(string(a.Value)); !ok {
					d.Scanner.InvalidAttr(a.Name.Local, false)
				}
			case attrThumbnail:
				d.model.Thumbnail = string(a.Value)
			case attrReqExt:
				requiredExts = strings.Fields(string(a.Value))
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

func (d *modelDecoder) noCoreAttribute(a XMLAttr) {
	switch a.Name.Space {
	case nsXML:
		if a.Name.Local == attrLang {
			d.model.Language = string(a.Value)
		}
	case attrXmlns:
		if ext, ok := d.model.Specs[string(a.Value)]; ok {
			ext.SetLocal(a.Name.Local)
		} else {
			d.model.WithSpec(&UnknownSpec{SpaceName: string(a.Value), LocalName: a.Name.Local})
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
	if name.Space == Namespace && name.Local == attrMetadata {
		child = &metadataDecoder{metadatas: d.metadatas}
	}
	return
}

type metadataDecoder struct {
	baseDecoder
	metadatas *[]Metadata
	metadata  Metadata
}

func (d *metadataDecoder) Start(attrs []XMLAttr) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrName:
			d.metadata.Name = a.Name
			i := bytes.IndexByte(a.Value, ':')
			if i < 0 {
				d.metadata.Name.Local = string(a.Value)
			} else if _, ok := d.Scanner.namespace(string(a.Value[0:i])); ok {
				d.metadata.Name.Space = string(a.Value[0:i])
				d.metadata.Name.Local = string(a.Value[i+1:])
			} else {
				d.metadata.Name.Local = string(a.Value)
			}
		case attrType:
			d.metadata.Type = string(a.Value)
		case attrPreserve:
			d.metadata.Preserve, _ = strconv.ParseBool(string(a.Value))
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
	if name.Space == Namespace && name.Local == attrItem {
		child = &buildItemDecoder{build: d.build}
	}
	return
}

func (d *buildDecoder) Start(attrs []XMLAttr) {
	for _, a := range attrs {
		if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, d.build, a)
		}
	}
}

type buildItemDecoder struct {
	baseDecoder
	build *Build
	item  Item
}

func (d *buildItemDecoder) End() {
	d.build.Items = append(d.build.Items, &d.item)
	d.Scanner.resourceID = 0
}

func (d *buildItemDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == Namespace && name.Local == attrMetadataGroup {
		child = &metadataGroupDecoder{metadatas: &d.item.Metadata}
	}
	return
}

func (d *buildItemDecoder) Start(attrs []XMLAttr) {
	for _, a := range attrs {
		if a.Name.Space == "" {
			d.parseCoreAttr(a)
		} else if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, &d.item, a)
		}
	}
	return
}

func (d *buildItemDecoder) parseCoreAttr(a XMLAttr) {
	switch a.Name.Local {
	case attrObjectID:
		val, err := strconv.ParseUint(string(a.Value), 10, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, true)
		}
		d.item.ObjectID = uint32(val)
		d.Scanner.resourceID = d.item.ObjectID
	case attrPartNumber:
		d.item.PartNumber = string(a.Value)
	case attrTransform:
		var ok bool
		d.item.Transform, ok = ParseMatrix(string(a.Value))
		if !ok {
			d.Scanner.InvalidAttr(a.Name.Local, false)
		}
	}
}

type resourceDecoder struct {
	baseDecoder
	resources *Resources
}

func (d *resourceDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == Namespace {
		switch name.Local {
		case attrObject:
			child = &objectDecoder{resources: d.resources}
		case attrBaseMaterials:
			child = &baseMaterialsDecoder{resources: d.resources}
		}
	} else if ext, ok := d.Scanner.extensionDecoder[name.Space]; ok {
		if ext, ok := ext.(resourcesElementDecoder); ok {
			child = ext.NewResourcesElementDecoder(d.resources, name.Local)
		}
	}
	if child != nil {
		child = &resourceDecoderWrapper{NodeDecoder: child}
	}
	return
}

type resourceDecoderWrapper struct {
	NodeDecoder
	Scanner *Scanner
}

func (d *resourceDecoderWrapper) SetScanner(s *Scanner) {
	d.Scanner = s
	d.NodeDecoder.SetScanner(s)
}

func (d *resourceDecoderWrapper) Start(attrs []XMLAttr) {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			id, _ := strconv.ParseUint(string(a.Value), 10, 32)
			d.Scanner.resourceID = uint32(id)
			break
		}
	}
	d.NodeDecoder.Start(attrs)
}

func (d *resourceDecoderWrapper) End() {
	d.NodeDecoder.End()
	d.Scanner.resourceID = 0
}

type baseMaterialsDecoder struct {
	baseDecoder
	resources           *Resources
	resource            BaseMaterials
	baseMaterialDecoder baseMaterialDecoder
}

func (d *baseMaterialsDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *baseMaterialsDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == Namespace && name.Local == attrBase {
		child = &d.baseMaterialDecoder
	}
	return
}

func (d *baseMaterialsDecoder) Start(attrs []XMLAttr) {
	d.baseMaterialDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			id, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, true)
			}
			d.resource.ID = uint32(id)
			break
		}
	}
}

type baseMaterialDecoder struct {
	baseDecoder
	resource *BaseMaterials
}

func (d *baseMaterialDecoder) Start(attrs []XMLAttr) {
	var name string
	var baseColor color.RGBA
	for _, a := range attrs {
		switch a.Name.Local {
		case attrName:
			name = string(a.Value)
		case attrDisplayColor:
			var err error
			baseColor, err = ParseRGBA(string(a.Value))
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, true)
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

func (d *meshDecoder) Start(_ []XMLAttr) {
	d.resource.Mesh = new(Mesh)
}

func (d *meshDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == Namespace {
		if name.Local == attrVertices {
			child = &verticesDecoder{mesh: d.resource.Mesh}
		} else if name.Local == attrTriangles {
			child = &trianglesDecoder{resource: d.resource}
		}
	} else if ext, ok := d.Scanner.extensionDecoder[name.Space]; ok {
		if ext, ok := ext.(meshElementDecoder); ok {
			child = ext.NewMeshElementDecoder(d.resource.Mesh, name.Local)
		}
	}
	return
}

type verticesDecoder struct {
	baseDecoder
	mesh          *Mesh
	vertexDecoder vertexDecoder
}

func (d *verticesDecoder) Start(_ []XMLAttr) {
	d.vertexDecoder.mesh = d.mesh
}

func (d *verticesDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == Namespace && name.Local == attrVertex {
		child = &d.vertexDecoder
	}
	return
}

type vertexDecoder struct {
	baseDecoder
	mesh *Mesh
}

func (d *vertexDecoder) Start(attrs []XMLAttr) {
	var x, y, z float32
	for _, a := range attrs {
		val, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&a.Value)), 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, true)
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

func (d *trianglesDecoder) Start(_ []XMLAttr) {
	d.triangleDecoder.mesh = d.resource.Mesh
	d.triangleDecoder.defaultPropertyID = d.resource.PID
	d.triangleDecoder.defaultPropertyIndex = d.resource.PIndex

	if len(d.resource.Mesh.Triangles) == 0 && len(d.resource.Mesh.Vertices) > 0 {
		d.resource.Mesh.Triangles = make([]Triangle, 0, len(d.resource.Mesh.Vertices)*2)
	}
}

func (d *trianglesDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == Namespace && name.Local == attrTriangle {
		child = &d.triangleDecoder
	}
	return
}

type triangleDecoder struct {
	baseDecoder
	mesh                                    *Mesh
	defaultPropertyIndex, defaultPropertyID uint32
}

func (d *triangleDecoder) Start(attrs []XMLAttr) {
	var t Triangle
	var pid, p1, p2, p3 uint32
	var hasPID, hasP1, hasP2, hasP3 bool
	for _, a := range attrs {
		required := true
		val, err := strconv.ParseUint(string(a.Value), 10, 24)
		switch a.Name.Local {
		case attrV1:
			t[0] = ToUint24(uint32(val))
		case attrV2:
			t[1] = ToUint24(uint32(val))
		case attrV3:
			t[2] = ToUint24(uint32(val))
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
			d.Scanner.InvalidAttr(a.Name.Local, required)
		}
	}

	p1 = applyDefault(p1, d.defaultPropertyIndex, hasP1)
	p2 = applyDefault(p2, p1, hasP2)
	p3 = applyDefault(p3, p1, hasP3)
	pid = applyDefault(pid, d.defaultPropertyID, hasPID)
	t.SetPID(pid)
	t.SetPIndices(p1, p2, p3)
	d.mesh.Triangles = append(d.mesh.Triangles, t)
}

func applyDefault(val, defVal uint32, noDef bool) uint32 {
	if noDef {
		return val
	}
	return defVal
}

type objectDecoder struct {
	baseDecoder
	resources *Resources
	resource  Object
}

func (d *objectDecoder) End() {
	d.resources.Objects = append(d.resources.Objects, &d.resource)
}

func (d *objectDecoder) Start(attrs []XMLAttr) {
	for _, a := range attrs {
		if a.Name.Space == "" {
			d.parseCoreAttr(a)
		} else if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, &d.resource, a)
		}
	}
}

func (d *objectDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == Namespace {
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

func (d *objectDecoder) parseCoreAttr(a XMLAttr) {
	switch a.Name.Local {
	case attrID:
		id, err := strconv.ParseUint(string(a.Value), 10, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, true)
		}
		d.resource.ID = uint32(id)
	case attrType:
		var ok bool
		d.resource.Type, ok = newObjectType(string(a.Value))
		if !ok {
			d.Scanner.InvalidAttr(a.Name.Local, false)
		}
	case attrThumbnail:
		d.resource.Thumbnail = string(a.Value)
	case attrName:
		d.resource.Name = string(a.Value)
	case attrPartNumber:
		d.resource.PartNumber = string(a.Value)
	case attrPID:
		val, err := strconv.ParseUint(string(a.Value), 10, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, false)
		}
		d.resource.PID = uint32(val)
	case attrPIndex:
		val, err := strconv.ParseUint(string(a.Value), 10, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, false)
		}
		d.resource.PIndex = uint32(val)
	}
}

type componentsDecoder struct {
	baseDecoder
	resource         *Object
	componentDecoder componentDecoder
}

func (d *componentsDecoder) Start(_ []XMLAttr) {
	d.resource.Components = make([]*Component, 0)
	d.componentDecoder.resource = d.resource
}

func (d *componentsDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == Namespace && name.Local == attrComponent {
		child = &d.componentDecoder
	}
	return
}

type componentDecoder struct {
	baseDecoder
	resource *Object
}

func (d *componentDecoder) Start(attrs []XMLAttr) {
	var component Component
	for _, a := range attrs {
		if a.Name.Space == "" {
			if a.Name.Local == attrObjectID {
				val, err := strconv.ParseUint(string(a.Value), 10, 32)
				if err != nil {
					d.Scanner.InvalidAttr(a.Name.Local, true)
				}
				component.ObjectID = uint32(val)
			} else if a.Name.Local == attrTransform {
				var ok bool
				component.Transform, ok = ParseMatrix(string(a.Value))
				if !ok {
					d.Scanner.InvalidAttr(a.Name.Local, false)
				}
			}
		} else if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, &component, a)
		}
	}
	d.resource.Components = append(d.resource.Components, &component)
}
