package io3mf

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/iohelper"
)

type meshDecoder struct {
	iohelper.EmptyDecoder
	mesh go3mf.Mesh
}

func (d *meshDecoder) Close() bool {
	d.Scanner.AddResource(&d.mesh)
	return true
}

func (d *meshDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec {
		if name.Local == attrVertices {
			child = &verticesDecoder{mesh: &d.mesh}
		} else if name.Local == attrTriangles {
			child = &trianglesDecoder{mesh: &d.mesh}
		}
	} else if ext, ok := extensionDecoder[name.Space]; ok {
		child = ext.NodeDecoder(&d.mesh, name.Local)
	}
	return
}

type verticesDecoder struct {
	iohelper.EmptyDecoder
	mesh          *go3mf.Mesh
	vertexDecoder vertexDecoder
}

func (d *verticesDecoder) Open() {
	d.vertexDecoder.mesh = d.mesh
}

func (d *verticesDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrVertex {
		child = &d.vertexDecoder
	}
	return
}

type vertexDecoder struct {
	iohelper.EmptyDecoder
	mesh *go3mf.Mesh
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
	d.mesh.Nodes = append(d.mesh.Nodes, go3mf.Point3D{x, y, z})
	return true
}

type trianglesDecoder struct {
	iohelper.EmptyDecoder
	mesh            *go3mf.Mesh
	triangleDecoder triangleDecoder
}

func (d *trianglesDecoder) Open() {
	d.triangleDecoder.mesh = d.mesh

	if len(d.mesh.Faces) == 0 && len(d.mesh.Nodes) > 0 {
		d.mesh.Faces = make([]go3mf.Face, 0, len(d.mesh.Nodes)-1)
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
	mesh *go3mf.Mesh
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

	p1 = applyDefault(p1, d.mesh.DefaultPropertyIndex, hasP1)
	p2 = applyDefault(p2, p1, hasP2)
	p3 = applyDefault(p3, p1, hasP3)
	pid = applyDefault(pid, d.mesh.DefaultPropertyID, hasPID)

	return d.addTriangle(v1, v2, v3, pid, p1, p2, p3)
}

func (d *triangleDecoder) addTriangle(v1, v2, v3, pid, p1, p2, p3 uint32) bool {
	if v1 == v2 || v1 == v3 || v2 == v3 {
		return d.Scanner.GenericError(true, "duplicated triangle indices")
	}
	nodeCount := uint32(len(d.mesh.Nodes))
	if v1 >= nodeCount || v2 >= nodeCount || v3 >= nodeCount {
		return d.Scanner.GenericError(true, "triangle indices are out of range")
	}
	d.mesh.Faces = append(d.mesh.Faces, go3mf.Face{
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
