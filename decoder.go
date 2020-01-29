package go3mf

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"strings"
)

type modelDecoder struct {
	BaseDecoder
	model                *Model
	withinIgnoredElement bool
}

func (d *modelDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName {
		switch name.Local {
		case attrResources:
			d.withinIgnoredElement = false
			child = &resourceDecoder{}
		case attrBuild:
			if !d.Scanner.IsRoot {
				d.withinIgnoredElement = true
			} else {
				d.withinIgnoredElement = false
				child = &buildDecoder{build: &d.model.Build}
			}
		case attrMetadata:
			if !d.Scanner.IsRoot {
				d.withinIgnoredElement = true
			} else {
				d.withinIgnoredElement = true
				child = &metadataDecoder{metadatas: &d.model.Metadata}
			}
		}
	}
	return
}

func (d *modelDecoder) Attributes(attrs []xml.Attr) {
	var requiredExts string
	for _, a := range attrs {
		if a.Name.Space == "" {
			switch a.Name.Local {
			case attrUnit:
				if d.Scanner.IsRoot {
					var ok bool
					if d.model.Units, ok = newUnits(a.Value); !ok {
						d.Scanner.InvalidOptionalAttr(attrUnit, a.Value)
					}
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
		ext = d.Scanner.Namespaces[ext]
		if ext != ExtensionName {
			if _, ok := d.Scanner.extensionDecoder[ext]; !ok {
				d.Scanner.GenericError(true, fmt.Sprintf("'%s' extension is not supported", ext))
			}
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
		d.Scanner.Namespaces[a.Name.Local] = a.Value
	}
}

type metadataGroupDecoder struct {
	BaseDecoder
	metadatas *[]Metadata
}

func (d *metadataGroupDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrMetadata {
		child = &metadataDecoder{metadatas: d.metadatas}
	}
	return
}

type metadataDecoder struct {
	BaseDecoder
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
			} else if ns, ok := d.Scanner.Namespaces[a.Value[0:i]]; ok {
				d.metadata.Name = ns + ":" + a.Value[i+1:]
			} else {
				d.Scanner.GenericError(true, "unregistered namespace")
			}
		case attrType:
			d.metadata.Type = a.Value
		case attrPreserve:
			if a.Value != "0" {
				d.metadata.Preserve = true
			}
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
	BaseDecoder
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
	BaseDecoder
	item Item
}

func (d *buildItemDecoder) Close() {
	d.Scanner.BuildItems = append(d.Scanner.BuildItems, &d.item)
	d.Scanner.CloseResource()
}

func (d *buildItemDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrMetadataGroup {
		child = &metadataGroupDecoder{metadatas: &d.item.Metadata}
	}
	return
}

// TODO: validate coeherence after parsing
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
	return
}

func (d *buildItemDecoder) parseCoreAttr(a xml.Attr) {
	switch a.Name.Local {
	case attrObjectID:
		d.item.ObjectID = d.Scanner.ParseResourceID(a.Value)
	case attrPartNumber:
		d.item.PartNumber = a.Value
	case attrTransform:
		var ok bool
		d.item.Transform, ok = ParseToMatrix(a.Value)
		if !ok {
			d.Scanner.InvalidOptionalAttr(a.Name.Local, a.Value)
		}
	}
}

type resourceDecoder struct {
	BaseDecoder
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
	BaseDecoder
	resource            BaseMaterialsResource
	baseMaterialDecoder baseMaterialDecoder
}

func (d *baseMaterialsDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
	d.baseMaterialDecoder.resource = &d.resource
}

func (d *baseMaterialsDecoder) Close() {
	d.Scanner.CloseResource()
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
	BaseDecoder
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
		case attrBaseMaterialColor:
			var err error
			baseColor, err = ParseRGB(a.Value)
			withColor = true
			if err != nil {
				d.Scanner.InvalidRequiredAttr(attrBaseMaterialColor, a.Value)
			}
		}
	}
	if name == "" {
		d.Scanner.MissingAttr(attrName)
	}
	if !withColor {
		d.Scanner.MissingAttr(attrBaseMaterialColor)
	}
	d.resource.Materials = append(d.resource.Materials, BaseMaterial{Name: name, Color: baseColor})
	return
}

type meshDecoder struct {
	BaseDecoder
	mesh Mesh
}

func (d *meshDecoder) Close() {
	d.Scanner.AddResource(&d.mesh)
}

func (d *meshDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName {
		if name.Local == attrVertices {
			child = &verticesDecoder{mesh: &d.mesh}
		} else if name.Local == attrTriangles {
			child = &trianglesDecoder{mesh: &d.mesh}
		}
	} else if ext, ok := d.Scanner.extensionDecoder[name.Space]; ok {
		child = ext.NewNodeDecoder(&d.mesh, name.Local)
	}
	return
}

type verticesDecoder struct {
	BaseDecoder
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
	BaseDecoder
	mesh *Mesh
}

func (d *vertexDecoder) Attributes(attrs []xml.Attr) {
	var x, y, z float32
	for _, a := range attrs {
		switch a.Name.Local {
		case attrX:
			x = d.Scanner.ParseFloat32Required(attrX, a.Value)
		case attrY:
			y = d.Scanner.ParseFloat32Required(attrY, a.Value)
		case attrZ:
			z = d.Scanner.ParseFloat32Required(attrZ, a.Value)
		}
	}
	d.mesh.Nodes = append(d.mesh.Nodes, Point3D{x, y, z})
}

type trianglesDecoder struct {
	BaseDecoder
	mesh            *Mesh
	triangleDecoder triangleDecoder
}

func (d *trianglesDecoder) Open() {
	d.triangleDecoder.mesh = d.mesh

	if len(d.mesh.Faces) == 0 && len(d.mesh.Nodes) > 0 {
		d.mesh.Faces = make([]Face, 0, len(d.mesh.Nodes)-1)
	}
}

func (d *trianglesDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrTriangle {
		child = &d.triangleDecoder
	}
	return
}

type triangleDecoder struct {
	BaseDecoder
	mesh *Mesh
}

func (d *triangleDecoder) Attributes(attrs []xml.Attr) {
	var v1, v2, v3, pid, p1, p2, p3 uint32
	var hasPID, hasP1, hasP2, hasP3 bool
	for _, a := range attrs {
		switch a.Name.Local {
		case attrV1:
			v1 = d.Scanner.ParseUint32Required(attrV1, a.Value)
		case attrV2:
			v2 = d.Scanner.ParseUint32Required(attrV2, a.Value)
		case attrV3:
			v3 = d.Scanner.ParseUint32Required(attrV3, a.Value)
		case attrPID:
			pid = d.Scanner.ParseUint32Optional(attrPID, a.Value)
			hasPID = true
		case attrP1:
			p1 = d.Scanner.ParseUint32Optional(attrP1, a.Value)
			hasP1 = true
		case attrP2:
			p2 = d.Scanner.ParseUint32Optional(attrP2, a.Value)
			hasP2 = true
		case attrP3:
			p3 = d.Scanner.ParseUint32Optional(attrP3, a.Value)
			hasP3 = true
		}
	}

	p1 = applyDefault(p1, d.mesh.DefaultPropertyIndex, hasP1)
	p2 = applyDefault(p2, p1, hasP2)
	p3 = applyDefault(p3, p1, hasP3)
	pid = applyDefault(pid, d.mesh.DefaultPropertyID, hasPID)

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
		Resource:        pid,
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
	BaseDecoder
	resource ObjectResource
}

func (d *objectDecoder) Open() {
	d.resource.ModelPath = d.Scanner.ModelPath
}

func (d *objectDecoder) Close() {
	d.Scanner.CloseResource()
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
			child = &meshDecoder{mesh: Mesh{ObjectResource: d.resource}}
		} else if name.Local == attrComponents {
			if d.resource.DefaultPropertyID != 0 {
				d.Scanner.GenericError(true, "default PID is not supported for component objects")
			}
			child = &componentsDecoder{resource: Components{ObjectResource: d.resource}}
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
}

type componentsDecoder struct {
	BaseDecoder
	resource         Components
	componentDecoder componentDecoder
}

func (d *componentsDecoder) Open() {
	d.componentDecoder.resource = &d.resource
}
func (d *componentsDecoder) Close() {
	d.Scanner.AddResource(&d.resource)
}

func (d *componentsDecoder) Child(name xml.Name) (child NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrComponent {
		child = &d.componentDecoder
	}
	return
}

type componentDecoder struct {
	BaseDecoder
	resource *Components
}

func (d *componentDecoder) Attributes(attrs []xml.Attr) {
	var component Component
	for _, a := range attrs {
		if a.Name.Space == "" {
			if a.Name.Local == attrObjectID {
				component.ObjectID = d.Scanner.ParseUint32Required(attrObjectID, a.Value)
			} else if a.Name.Local == attrTransform {
				var ok bool
				component.Transform, ok = ParseToMatrix(a.Value)
				if !ok {
					d.Scanner.InvalidOptionalAttr(a.Name.Local, a.Value)
				}
			}
		} else if ext, ok := d.Scanner.extensionDecoder[a.Name.Space]; ok {
			ext.DecodeAttribute(d.Scanner, &component, a)
		}
	}
	d.resource.Components = append(d.resource.Components, &component)
}

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
