package go3mf

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

type modelDecoder struct {
	baseDecoder
	model *Model
}

func (d *modelDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName {
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
	}
	return
}

func (d *modelDecoder) Attributes(attrs []xml.Attr) {
	if !d.Scanner.IsRoot {
		return
	}
	var requiredExts string
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
				requiredExts = a.Value
			}
		} else {
			d.noCoreAttribute(a)
		}
	}

	d.checkRequiredExt(requiredExts)
}

func (d *modelDecoder) checkRequiredExt(requiredExts string) {
	for _, ext := range strings.Fields(requiredExts) {
		if _, ok := d.Scanner.Namespace(ext); !ok {
			d.Scanner.GenericError(true, fmt.Sprintf("'%s' extension is not supported", ext))
		}
	}
}

func (d *modelDecoder) noCoreAttribute(a xml.Attr) {
	switch a.Name.Space {
	case nsXML:
		if d.Scanner.IsRoot {
			if a.Name.Local == attrLang {
				d.model.Language = a.Value
			}
		}
	case attrXmlns:
		d.Scanner.Namespaces = append(d.Scanner.Namespaces, xml.Name{Space: a.Value, Local: a.Name.Local})
	}
}

type metadataGroupDecoder struct {
	baseDecoder
	metadatas *[]Metadata
}

func (d *metadataGroupDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrMetadata {
		child = &metadataDecoder{metadatas: d.metadatas}
	}
	return
}

type metadataDecoder struct {
	baseDecoder
	metadatas *[]Metadata
	metadata  Metadata
}

func (d *metadataDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrName:
			i := strings.IndexByte(a.Value, ':')
			if i < 0 {
				d.metadata.Name = a.Value
			} else if _, ok := d.Scanner.Namespace(a.Value[0:i]); ok {
				d.metadata.Name = a.Value[0:i] + ":" + a.Value[i+1:]
			} else {
				d.Scanner.GenericError(true, "unregistered namespace")
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

func (d *metadataDecoder) Close() {
	*d.metadatas = append(*d.metadatas, d.metadata)
}

type buildDecoder struct {
	baseDecoder
	build *Build
}

func (d *buildDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrItem {
		child = &buildItemDecoder{}
	}
	return
}

func (d *buildDecoder) Attributes(attrs []xml.Attr) {
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

func (d *buildItemDecoder) Close() {
	d.Scanner.BuildItems = append(d.Scanner.BuildItems, &d.item)
	d.Scanner.ResourceID = 0
}

func (d *buildItemDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrMetadataGroup {
		child = &metadataGroupDecoder{metadatas: &d.item.Metadata}
	}
	return
}

// TODO: validate coeherence after decoding
// func (d *buildItemDecoder) processItem() {
// 	resource, ok := d.Scanner.FindResource(d.objectPath, uint32(d.objectID))
// 	if !ok {
// 		d.Scanner.GenericError(true, "non-existent referenced object")
// 	} else if d.item.Object, ok = resource.(Object); !ok {
// 		d.Scanner.GenericError(true, "non-object referenced resource")
// 	}
// 	if ok {
// 		if d.item.Object != nil && d.item.Object.Type() == ObjectTypeOther {
// 			d.Scanner.GenericError(true, "referenced object cannot be have OTHER type")
// 		}
// 	}
// 	if ok {
// 	}
// }

func (d *buildItemDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space == "" {
			d.parseCoreAttr(a)
		} else if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, &d.item, a)
		}
	}
	if d.item.ObjectID == 0 {
		d.Scanner.MissingAttr(attrObjectID)
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
	if name.Space == ExtensionName {
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
	resource            BaseMaterialsResource
	baseMaterialDecoder baseMaterialDecoder
}

func (d *baseMaterialsDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.baseMaterialDecoder.resource = &d.resource
}

func (d *baseMaterialsDecoder) Close() {
	d.Scanner.AddResource(&d.resource)
}

func (d *baseMaterialsDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrBase {
		child = &d.baseMaterialDecoder
	}
	return
}

func (d *baseMaterialsDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			d.resource.ID = d.Scanner.ParseResourceID(a.Value)
			break
		}
	}
}

type baseMaterialDecoder struct {
	baseDecoder
	resource *BaseMaterialsResource
}

func (d *baseMaterialDecoder) Attributes(attrs []xml.Attr) {
	var name string
	var withColor bool
	baseColor := color.RGBA{}
	for _, a := range attrs {
		switch a.Name.Local {
		case attrName:
			name = a.Value
		case attrDisplayColor:
			var err error
			baseColor, err = ParseRGB(a.Value)
			withColor = true
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
		}
	}
	if name == "" {
		d.Scanner.MissingAttr(attrName)
	}
	if !withColor {
		d.Scanner.MissingAttr(attrDisplayColor)
	}
	d.resource.Materials = append(d.resource.Materials, BaseMaterial{Name: name, Color: baseColor})
	return
}

type meshDecoder struct {
	baseDecoder
	resource *ObjectResource
}

func (d *meshDecoder) Open() {
	d.resource.Mesh = new(Mesh)
}

func (d *meshDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName {
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

func (d *verticesDecoder) Open() {
	d.vertexDecoder.mesh = d.mesh
}

func (d *verticesDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrVertex {
		child = &d.vertexDecoder
	}
	return
}

type vertexDecoder struct {
	baseDecoder
	mesh *Mesh
}

func (d *vertexDecoder) Attributes(attrs []xml.Attr) {
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
	d.mesh.Nodes = append(d.mesh.Nodes, Point3D{x, y, z})
}

type trianglesDecoder struct {
	baseDecoder
	resource        *ObjectResource
	triangleDecoder triangleDecoder
}

func (d *trianglesDecoder) Open() {
	d.triangleDecoder.mesh = d.resource.Mesh
	d.triangleDecoder.defaultPropertyID = d.resource.DefaultPropertyID
	d.triangleDecoder.defaultPropertyIndex = d.resource.DefaultPropertyIndex

	if len(d.resource.Mesh.Faces) == 0 && len(d.resource.Mesh.Nodes) > 0 {
		d.resource.Mesh.Faces = make([]Face, 0, len(d.resource.Mesh.Nodes)-1)
	}
}

func (d *trianglesDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrTriangle {
		child = &d.triangleDecoder
	}
	return
}

type triangleDecoder struct {
	baseDecoder
	mesh                                    *Mesh
	defaultPropertyIndex, defaultPropertyID uint32
}

func (d *triangleDecoder) Attributes(attrs []xml.Attr) {
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

	d.addTriangle(v1, v2, v3, pid, p1, p2, p3)
}

func (d *triangleDecoder) addTriangle(v1, v2, v3, pid, p1, p2, p3 uint32) {
	if v1 == v2 || v1 == v3 || v2 == v3 {
		d.Scanner.GenericError(true, "duplicated triangle indices")
	}
	nodeCount := uint32(len(d.mesh.Nodes))
	if v1 >= nodeCount || v2 >= nodeCount || v3 >= nodeCount {
		d.Scanner.GenericError(true, "triangle indices are out of range")
	}
	d.mesh.Faces = append(d.mesh.Faces, Face{
		NodeIndices:     [3]uint32{v1, v2, v3},
		PID:             pid,
		ResourceIndices: [3]uint32{p1, p2, p3},
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
	resource ObjectResource
}

func (d *objectDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
}

func (d *objectDecoder) Close() {
	d.Scanner.AddResource(&d.resource)
}

func (d *objectDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space == "" {
			d.parseCoreAttr(a)
		} else if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, &d.resource, a)
		}
	}
}

func (d *objectDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName {
		if name.Local == attrMesh {
			child = &meshDecoder{resource: &d.resource}
		} else if name.Local == attrComponents {
			if d.resource.DefaultPropertyID != 0 {
				d.Scanner.GenericError(true, "default PID is not supported for component objects")
			}
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
		d.resource.ID = d.Scanner.ParseResourceID(a.Value)
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
		d.resource.DefaultPropertyID = uint32(val)
	case attrPIndex:
		val, err := strconv.ParseUint(a.Value, 10, 32)
		if err != nil {
			d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
		}
		d.resource.DefaultPropertyIndex = uint32(val)
	}
}

type componentsDecoder struct {
	baseDecoder
	resource         *ObjectResource
	componentDecoder componentDecoder
}

func (d *componentsDecoder) Open() {
	d.resource.Components = make([]*Component, 0)
	d.componentDecoder.resource = d.resource
}

func (d *componentsDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrComponent {
		child = &d.componentDecoder
	}
	return
}

type componentDecoder struct {
	baseDecoder
	resource *ObjectResource
}

func (d *componentDecoder) Attributes(attrs []xml.Attr) {
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
	if component.ObjectID == 0 {
		d.Scanner.MissingAttr(attrObjectID)
	}
	d.resource.Components = append(d.resource.Components, &component)
}

// TODO: validate coeherence after decoding
// func (d *componentDecoder) addComponent(component *Component, path string, objectID uint32) {
// if path != "" && !d.Scanner.IsRoot {
// d.Scanner.GenericError(true, "path attribute in a non-root file is not supported")
// }

// resource, ok := d.Scanner.FindResource(path, uint32(objectID))
// if !ok {
// d.Scanner.GenericError(true, "non-existent referenced object")
// } else if component.Object, ok = resource.(Object); !ok {
// d.Scanner.GenericError(true, "non-object referenced resource")
// }
// d.resource.Components = append(d.resource.Components, component)
// }
