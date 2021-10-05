// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package go3mf

import (
	"bytes"
	"encoding/xml"
	"strconv"
	"strings"
	"unsafe"

	specerr "github.com/hpinc/go3mf/errors"
	"github.com/hpinc/go3mf/spec"
)

type modelDecoder struct {
	baseDecoder
	model  *Model
	isRoot bool
	path   string
}

func (d *modelDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace {
		switch name.Local {
		case attrResources:
			resources, _ := d.model.FindResources(d.path)
			child = &resourceDecoder{resources: resources, model: d.model}
		case attrBuild:
			if d.isRoot {
				child = &buildDecoder{build: &d.model.Build, model: d.model}
			}
		case attrMetadata:
			if d.isRoot {
				child = &metadataDecoder{metadatas: &d.model.Metadata, model: d.model}
			}
		}
	} else if ext, ok := spec.Load(name.Space); ok {
		child = ext.NewElementDecoder(d.model, name.Local)
	} else {
		child = spec.NewAnyUnknownDecoder(name, &d.model.Any)
	}
	return
}

func (d *modelDecoder) Start(attrs []spec.XMLAttr) (err error) {
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
					err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, false))
				}
			case attrThumbnail:
				d.model.Thumbnail = string(a.Value)
			case attrReqExt:
				requiredExts = strings.Fields(string(a.Value))
			}
		} else {
			err = specerr.Append(err, d.noCoreAttribute(a))
		}
	}

	for _, local := range requiredExts {
		for i := range d.model.Extensions {
			ext := &d.model.Extensions[i]
			if ext.LocalName == local {
				ext.IsRequired = true
				break
			}
		}
	}
	return
}

func (d *modelDecoder) noCoreAttribute(a spec.XMLAttr) (err error) {
	switch a.Name.Space {
	case nsXML:
		if a.Name.Local == attrLang {
			d.model.Language = string(a.Value)
		}
	case attrXmlns:
		d.model.Extensions = append(d.model.Extensions, Extension{
			Namespace:  string(a.Value),
			LocalName:  a.Name.Local,
			IsRequired: false,
		})
	default:
		var attr spec.AttrGroup
		if attr = d.model.AnyAttr.Get(a.Name.Space); attr == nil {
			attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrModel})
			d.model.AnyAttr = append(d.model.AnyAttr, attr)
		}
		err = specerr.Append(err, attr.Unmarshal3MFAttr(a))
	}
	return
}

type metadataGroupDecoder struct {
	baseDecoder
	metadatas *[]Metadata
	model     *Model
}

func (d *metadataGroupDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrMetadata {
		child = &metadataDecoder{metadatas: d.metadatas, model: d.model}
	}
	return
}

type metadataDecoder struct {
	baseDecoder
	model     *Model
	metadatas *[]Metadata
	metadata  Metadata
}

func (d *metadataDecoder) namespace(local string) (string, bool) {
	for _, ext := range d.model.Extensions {
		if ext.LocalName == local {
			return ext.Namespace, true
		}
	}
	return "", false
}

func (d *metadataDecoder) Start(attrs []spec.XMLAttr) error {
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
			} else if _, ok := d.namespace(string(a.Value[0:i])); ok {
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
	return nil
}

func (d *metadataDecoder) CharData(txt []byte) {
	d.metadata.Value = string(txt)
}

func (d *metadataDecoder) End() {
	*d.metadatas = append(*d.metadatas, d.metadata)
}

type buildDecoder struct {
	baseDecoder
	model *Model
	build *Build
}

func (d *buildDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrItem {
		child = &buildItemDecoder{build: d.build, model: d.model}
	}
	return
}

func (d *buildDecoder) Wrap(err error) error {
	return specerr.Wrap(err, d.build)
}

func (d *buildDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	for _, a := range attrs {
		var attr spec.AttrGroup
		if attr = d.build.AnyAttr.Get(a.Name.Space); attr == nil {
			attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrBuild})
			d.build.AnyAttr = append(d.build.AnyAttr, attr)
		}
		errs = specerr.Append(errs, attr.Unmarshal3MFAttr(a))
	}
	if errs != nil {
		return specerr.Wrap(errs, d.build)
	}
	return errs
}

type buildItemDecoder struct {
	baseDecoder
	model *Model
	build *Build
	item  Item
}

func (d *buildItemDecoder) End() {
	d.build.Items = append(d.build.Items, &d.item)
}

func (d *buildItemDecoder) Wrap(err error) error {
	return specerr.WrapIndex(err, &d.item, len(d.build.Items))
}

func (d *buildItemDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrMetadataGroup {
		child = &metadataGroupDecoder{metadatas: &d.item.Metadata, model: d.model}
	}
	return
}

func (d *buildItemDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	for _, a := range attrs {
		if a.Name.Space == "" {
			errs = specerr.Append(errs, d.parseCoreAttr(a))
		} else {
			var attr spec.AttrGroup
			if attr = d.item.AnyAttr.Get(a.Name.Space); attr == nil {
				attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrItem})
				d.item.AnyAttr = append(d.item.AnyAttr, attr)
			}
			errs = specerr.Append(errs, attr.Unmarshal3MFAttr(a))
		}
	}
	if errs != nil {
		return specerr.WrapIndex(errs, &d.item, len(d.build.Items))
	}
	return errs
}

func (d *buildItemDecoder) parseCoreAttr(a spec.XMLAttr) (errs error) {
	switch a.Name.Local {
	case attrObjectID:
		val, err := strconv.ParseUint(string(a.Value), 10, 32)
		if err != nil {
			errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
		}
		d.item.ObjectID = uint32(val)
	case attrPartNumber:
		d.item.PartNumber = string(a.Value)
	case attrTransform:
		var ok bool
		d.item.Transform, ok = spec.ParseMatrix(string(a.Value))
		if !ok {
			errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
		}
	}
	return
}

type resourceDecoder struct {
	baseDecoder
	model     *Model
	resources *Resources
}

func (d *resourceDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	for _, a := range attrs {
		var attr spec.AttrGroup
		if attr = d.resources.AnyAttr.Get(a.Name.Space); attr == nil {
			attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrResources})
			d.resources.AnyAttr = append(d.resources.AnyAttr, attr)
		}
		errs = specerr.Append(errs, attr.Unmarshal3MFAttr(a))
	}
	if errs != nil {
		return specerr.Wrap(errs, d.resources)
	}
	return errs
}

func (d *resourceDecoder) Wrap(err error) error {
	return specerr.Wrap(err, d.resources)
}

func (d *resourceDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace {
		switch name.Local {
		case attrObject:
			child = &objectDecoder{resources: d.resources, model: d.model}
		case attrBaseMaterials:
			child = &baseMaterialsDecoder{resources: d.resources}
		}
	} else if ext, ok := spec.Load(name.Space); ok {
		child = ext.NewElementDecoder(d.resources, name.Local)
	} else {
		child = &unknownAssetDecoder{UnknownTokensDecoder: *spec.NewUnknownDecoder(name), resources: d.resources}
	}
	return
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

func (d *baseMaterialsDecoder) Wrap(err error) error {
	return specerr.WrapIndex(err, &d.resource, len(d.resources.Assets))
}

func (d *baseMaterialsDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrBase {
		child = &d.baseMaterialDecoder
	}
	return
}

func (d *baseMaterialsDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	d.baseMaterialDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space == "" {
			if a.Name.Local == attrID {
				id, err := strconv.ParseUint(string(a.Value), 10, 32)
				if err != nil {
					errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
				}
				d.resource.ID = uint32(id)
			}
		} else {
			var attr spec.AttrGroup
			if attr = d.resource.AnyAttr.Get(a.Name.Space); attr == nil {
				attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrBaseMaterials})
				d.resource.AnyAttr = append(d.resource.AnyAttr, attr)
			}
			errs = specerr.Append(errs, attr.Unmarshal3MFAttr(a))
		}
	}
	if errs != nil {
		return specerr.WrapIndex(errs, &d.resource, len(d.resources.Assets))
	}
	return nil
}

type baseMaterialDecoder struct {
	baseDecoder
	resource *BaseMaterials
}

func (d *baseMaterialDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		base Base
		errs error
	)
	for _, a := range attrs {
		if a.Name.Space == "" {
			switch a.Name.Local {
			case attrName:
				base.Name = string(a.Value)
			case attrDisplayColor:
				var err error
				base.Color, err = spec.ParseRGBA(string(a.Value))
				if err != nil {
					errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
				}
			}
		} else {
			var attr spec.AttrGroup
			if attr = base.AnyAttr.Get(a.Name.Space); attr == nil {
				attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrBase})
				base.AnyAttr = append(base.AnyAttr, attr)
			}
			errs = specerr.Append(errs, attr.Unmarshal3MFAttr(a))
		}
	}
	d.resource.Materials = append(d.resource.Materials, base)
	if errs != nil {
		return specerr.WrapIndex(errs, base, len(d.resource.Materials)-1)
	}
	return nil
}

type meshDecoder struct {
	baseDecoder
	resource *Object
}

func (d *meshDecoder) Start(attrs []spec.XMLAttr) error {
	d.resource.Mesh = new(Mesh)
	var errs error
	for _, a := range attrs {
		var attr spec.AttrGroup
		if attr = d.resource.Mesh.AnyAttr.Get(a.Name.Space); attr == nil {
			attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrMesh})
			d.resource.Mesh.AnyAttr = append(d.resource.Mesh.AnyAttr, attr)
		}
		errs = specerr.Append(errs, attr.Unmarshal3MFAttr(a))
	}
	if errs != nil {
		return specerr.Wrap(errs, d.resource.Mesh)
	}
	return errs
}

func (d *meshDecoder) Wrap(err error) error {
	return specerr.Wrap(err, d.resource.Mesh)
}

func (d *meshDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace {
		if name.Local == attrVertices {
			child = &verticesDecoder{mesh: d.resource.Mesh}
		} else if name.Local == attrTriangles {
			child = &trianglesDecoder{resource: d.resource}
		}
	} else if ext, ok := spec.Load(name.Space); ok {
		child = ext.NewElementDecoder(d.resource.Mesh, name.Local)
	} else {
		child = spec.NewAnyUnknownDecoder(name, &d.resource.Mesh.Any)
	}
	return
}

type verticesDecoder struct {
	baseDecoder
	mesh          *Mesh
	vertexDecoder vertexDecoder
}

func (d *verticesDecoder) Start(_ []spec.XMLAttr) error {
	d.vertexDecoder.mesh = d.mesh
	return nil
}

func (d *verticesDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrVertex {
		child = &d.vertexDecoder
	}
	return
}

type vertexDecoder struct {
	baseDecoder
	mesh *Mesh
}

func (d *vertexDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		x, y, z float32
		errs    error
	)
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		val, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&a.Value)), 32)
		if err != nil {
			errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
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
	if errs != nil {
		return specerr.WrapIndex(errs, Point3D{x, y, z}, len(d.mesh.Vertices)-1)
	}
	return nil
}

type trianglesDecoder struct {
	baseDecoder
	resource        *Object
	triangleDecoder triangleDecoder
}

func (d *trianglesDecoder) Start(_ []spec.XMLAttr) error {
	d.triangleDecoder.mesh = d.resource.Mesh
	d.triangleDecoder.defaultPropertyID = d.resource.PID
	d.triangleDecoder.defaultPropertyIndex = d.resource.PIndex

	if len(d.resource.Mesh.Triangles) == 0 && len(d.resource.Mesh.Vertices) > 0 {
		d.resource.Mesh.Triangles = make([]Triangle, 0, len(d.resource.Mesh.Vertices)*2)
	}
	return nil
}

func (d *trianglesDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
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

func (d *triangleDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		t                           Triangle
		pid, p1, p2, p3             uint32
		hasPID, hasP1, hasP2, hasP3 bool
		errs                        error
	)

	for _, a := range attrs {
		if a.Name.Space == "" {
			required := true
			val, err := strconv.ParseUint(string(a.Value), 10, 32)
			switch a.Name.Local {
			case attrV1:
				t.V1 = uint32(val)
			case attrV2:
				t.V2 = uint32(val)
			case attrV3:
				t.V3 = uint32(val)
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
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, required))
			}
		} else {
			var attr spec.AttrGroup
			if attr = t.AnyAttr.Get(a.Name.Space); attr == nil {
				attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrTriangle})
				t.AnyAttr = append(t.AnyAttr, attr)
			}
			errs = specerr.Append(errs, attr.Unmarshal3MFAttr(a))
		}
	}

	p1 = applyDefault(p1, d.defaultPropertyIndex, hasP1)
	p2 = applyDefault(p2, p1, hasP2)
	p3 = applyDefault(p3, p1, hasP3)
	pid = applyDefault(pid, d.defaultPropertyID, hasPID)
	t.PID = pid
	t.P1, t.P2, t.P3 = p1, p2, p3
	d.mesh.Triangles = append(d.mesh.Triangles, t)
	if errs != nil {
		return specerr.WrapIndex(errs, t, len(d.mesh.Triangles)-1)
	}
	return nil
}

func applyDefault(val, defVal uint32, noDef bool) uint32 {
	if noDef {
		return val
	}
	return defVal
}

type objectDecoder struct {
	baseDecoder
	model     *Model
	resources *Resources
	resource  Object
}

func (d *objectDecoder) End() {
	d.resources.Objects = append(d.resources.Objects, &d.resource)
}

func (d *objectDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	for _, a := range attrs {
		if a.Name.Space == "" {
			errs = specerr.Append(errs, d.parseCoreAttr(a))
		} else {
			var attr spec.AttrGroup
			if attr = d.resource.AnyAttr.Get(a.Name.Space); attr == nil {
				attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrObject})
				d.resource.AnyAttr = append(d.resource.AnyAttr, attr)
			}
			errs = specerr.Append(errs, attr.Unmarshal3MFAttr(a))
		}
	}
	if errs != nil {
		return specerr.WrapIndex(errs, &d.resource, len(d.resources.Objects))
	}
	return errs
}

func (d *objectDecoder) Wrap(err error) error {
	return specerr.WrapIndex(err, &d.resource, len(d.resources.Objects))
}

func (d *objectDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace {
		if name.Local == attrMesh {
			child = &meshDecoder{resource: &d.resource}
		} else if name.Local == attrComponents {
			child = &componentsDecoder{resource: &d.resource}
		} else if name.Local == attrMetadataGroup {
			child = &metadataGroupDecoder{metadatas: &d.resource.Metadata, model: d.model}
		}
	}
	return
}

func (d *objectDecoder) parseCoreAttr(a spec.XMLAttr) (errs error) {
	switch a.Name.Local {
	case attrID:
		id, err := strconv.ParseUint(string(a.Value), 10, 32)
		if err != nil {
			errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
		}
		d.resource.ID = uint32(id)
	case attrType:
		var ok bool
		d.resource.Type, ok = newObjectType(string(a.Value))
		if !ok {
			errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
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
			errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
		}
		d.resource.PID = uint32(val)
	case attrPIndex:
		val, err := strconv.ParseUint(string(a.Value), 10, 32)
		if err != nil {
			errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
		}
		d.resource.PIndex = uint32(val)
	}
	return
}

type componentsDecoder struct {
	baseDecoder
	resource         *Object
	componentDecoder componentDecoder
}

func (d *componentsDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	components := new(Components)
	d.componentDecoder.resource = d.resource

	for _, a := range attrs {
		var attr spec.AttrGroup
		if attr = components.AnyAttr.Get(a.Name.Space); attr == nil {
			attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrComponents})
			components.AnyAttr = append(components.AnyAttr, attr)
		}
		errs = specerr.Append(errs, attr.Unmarshal3MFAttr(a))
	}
	d.resource.Components = components
	if errs != nil {
		return specerr.Wrap(errs, d.resource.Components)
	}
	return nil
}

func (d *componentsDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrComponent {
		child = &d.componentDecoder
	}
	return
}

func (d *componentsDecoder) Wrap(err error) error {
	return specerr.Wrap(err, d.resource.Components)
}

type componentDecoder struct {
	baseDecoder
	resource *Object
}

func (d *componentDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		component Component
		errs      error
	)
	for _, a := range attrs {
		if a.Name.Space == "" {
			if a.Name.Local == attrObjectID {
				val, err := strconv.ParseUint(string(a.Value), 10, 32)
				if err != nil {
					errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
				}
				component.ObjectID = uint32(val)
			} else if a.Name.Local == attrTransform {
				var ok bool
				component.Transform, ok = spec.ParseMatrix(string(a.Value))
				if !ok {
					errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
				}
			}
		} else {
			var attr spec.AttrGroup
			if attr = component.AnyAttr.Get(a.Name.Space); attr == nil {
				attr = spec.NewAttrGroup(a.Name.Space, xml.Name{Space: Namespace, Local: attrComponent})
				component.AnyAttr = append(component.AnyAttr, attr)
			}
			errs = specerr.Append(errs, attr.Unmarshal3MFAttr(a))
		}
	}
	d.resource.Components.Component = append(d.resource.Components.Component, &component)
	if errs != nil {
		return specerr.WrapIndex(errs, &component, len(d.resource.Components.Component)-1)
	}
	return nil
}

type baseDecoder struct {
}

func (d *baseDecoder) Start([]spec.XMLAttr) error { return nil }
func (d *baseDecoder) End()                       {}

type topLevelDecoder struct {
	baseDecoder
	model  *Model
	isRoot bool
	path   string
}

func (d *topLevelDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	modelName := xml.Name{Space: Namespace, Local: attrModel}
	if name == modelName {
		child = &modelDecoder{model: d.model, isRoot: d.isRoot, path: d.path}
	}
	return
}

type unknownAssetDecoder struct {
	spec.UnknownTokensDecoder
	resources *Resources
	resource  UnknownAsset
}

func (d *unknownAssetDecoder) Start(attrs []spec.XMLAttr) (errs error) {
	d.UnknownTokensDecoder.Start(attrs)
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			id, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.id = uint32(id)
			break
		}
	}
	if errs != nil {
		return specerr.WrapIndex(errs, &d.resource, len(d.resources.Assets))
	}
	return nil
}

func (d *unknownAssetDecoder) End() {
	d.UnknownTokensDecoder.End()
	d.resource.UnknownTokens = d.UnknownTokensDecoder.Tokens()
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}
