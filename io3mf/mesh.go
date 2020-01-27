package io3mf

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/geo"
	"github.com/qmuntal/go3mf/iohelper"
)

type meshDecoder struct {
	iohelper.EmptyDecoder
	resource go3mf.MeshResource
}

func (d *meshDecoder) Open() {
	d.resource.Mesh = new(geo.Mesh)
}

func (d *meshDecoder) Close() bool {
	d.Scanner.AddResource(&d.resource)
	return true
}

func (d *meshDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec {
		if name.Local == attrVertices {
			child = &verticesDecoder{resource: &d.resource}
		} else if name.Local == attrTriangles {
			child = &trianglesDecoder{resource: &d.resource}
		}
	} else if ext, ok := extensionDecoder[name.Space]; ok {
		child = ext.NodeDecoder(&d.resource, name.Local)
	}
	return
}

type verticesDecoder struct {
	iohelper.EmptyDecoder
	resource      *go3mf.MeshResource
	vertexDecoder vertexDecoder
}

func (d *verticesDecoder) Open() {
	d.vertexDecoder.resource = d.resource
}

func (d *verticesDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrVertex {
		child = &d.vertexDecoder
	}
	return
}

type vertexDecoder struct {
	iohelper.EmptyDecoder
	resource *go3mf.MeshResource
}

func (d *vertexDecoder) Attributes(attrs []xml.Attr) bool {
	var x, y, z float32
	ok := true
	for _, a := range attrs {
		switch a.Name.Local {
		case attrX:
			x, ok = d.Scanner.ParseFloat32Required(attrX, a.Value)
		case attrY:
			y, ok = d.Scanner.ParseFloat32Required(attrY, a.Value)
		case attrZ:
			z, ok = d.Scanner.ParseFloat32Required(attrZ, a.Value)
		}
		if !ok {
			return false
		}
	}
	d.resource.Mesh.Nodes = append(d.resource.Mesh.Nodes, geo.Point3D{x, y, z})
	return true
}

type trianglesDecoder struct {
	iohelper.EmptyDecoder
	resource        *go3mf.MeshResource
	triangleDecoder triangleDecoder
}

func (d *trianglesDecoder) Open() {
	d.triangleDecoder.resource = d.resource

	if len(d.resource.Mesh.Faces) == 0 && len(d.resource.Mesh.Nodes) > 0 {
		d.resource.Mesh.Faces = make([]geo.Face, 0, len(d.resource.Mesh.Nodes)-1)
	}
}

func (d *trianglesDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrTriangle {
		child = &d.triangleDecoder
	}
	return
}

type triangleDecoder struct {
	iohelper.EmptyDecoder
	resource *go3mf.MeshResource
}

func (d *triangleDecoder) Attributes(attrs []xml.Attr) bool {
	var v1, v2, v3, pid, p1, p2, p3 uint32
	var hasPID, hasP1, hasP2, hasP3 bool
	ok := true
	for _, a := range attrs {
		switch a.Name.Local {
		case attrV1:
			v1, ok = d.Scanner.ParseUint32Required(attrV1, a.Value)
		case attrV2:
			v2, ok = d.Scanner.ParseUint32Required(attrV2, a.Value)
		case attrV3:
			v3, ok = d.Scanner.ParseUint32Required(attrV3, a.Value)
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
		if !ok {
			return false
		}
	}

	p1 = applyDefault(p1, d.resource.DefaultPropertyIndex, hasP1)
	p2 = applyDefault(p2, p1, hasP2)
	p3 = applyDefault(p3, p1, hasP3)
	pid = applyDefault(pid, d.resource.DefaultPropertyID, hasPID)

	return d.addTriangle(v1, v2, v3, pid, p1, p2, p3)
}

func (d *triangleDecoder) addTriangle(v1, v2, v3, pid, p1, p2, p3 uint32) bool {
	if v1 == v2 || v1 == v3 || v2 == v3 {
		return d.Scanner.GenericError(true, "duplicated triangle indices")
	}
	nodeCount := uint32(len(d.resource.Mesh.Nodes))
	if v1 >= nodeCount || v2 >= nodeCount || v3 >= nodeCount {
		return d.Scanner.GenericError(true, "triangle indices are out of range")
	}
	d.resource.Mesh.Faces = append(d.resource.Mesh.Faces, geo.Face{
		NodeIndices:     [3]uint32{v1, v2, v3},
		Resource:        pid,
		ResourceIndices: [3]uint32{p1, p2, p3},
	})
	return true
}

func applyDefault(val, defVal uint32, noDef bool) uint32 {
	if noDef {
		return val
	}
	return defVal
}
