package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/mesh"
)

type meshDecoder struct {
	emptyDecoder
	resource go3mf.MeshResource
}

func (d *meshDecoder) Open() error {
	d.resource.Mesh = new(mesh.Mesh)
	d.resource.Mesh.StartCreation(mesh.CreationOptions{CalculateConnectivity: false})
	return nil
}

func (d *meshDecoder) Close() error {
	d.resource.Mesh.EndCreation()
	d.ModelFile().AddResource(&d.resource)
	return nil
}

func (d *meshDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec {
		if name.Local == attrVertices {
			child = &verticesDecoder{resource: &d.resource}
		} else if name.Local == attrTriangles {
			child = &trianglesDecoder{resource: &d.resource}
		}
	} else if name.Space == nsBeamLatticeSpec && name.Local == attrBeamLattice {
		child = &beamLatticeDecoder{resource: &d.resource}
	}
	return
}

type verticesDecoder struct {
	emptyDecoder
	resource      *go3mf.MeshResource
	vertexDecoder vertexDecoder
}

func (d *verticesDecoder) Open() error {
	d.vertexDecoder.resource = d.resource
	return nil
}

func (d *verticesDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrVertex {
		child = &d.vertexDecoder
	}
	return
}

type vertexDecoder struct {
	emptyDecoder
	resource *go3mf.MeshResource
}

func (d *vertexDecoder) Attributes(attrs []xml.Attr) (err error) {
	var x, y, z float64
	for _, a := range attrs {
		switch a.Name.Local {
		case attrX:
			x, err = strconv.ParseFloat(a.Value, 32)
		case attrY:
			y, err = strconv.ParseFloat(a.Value, 32)
		case attrZ:
			z, err = strconv.ParseFloat(a.Value, 32)
		}
		if err != nil {
			break
		}
	}
	if err == nil {
		d.resource.Mesh.AddNode(mesh.Node{float32(x), float32(y), float32(z)})
	}
	return
}

type trianglesDecoder struct {
	emptyDecoder
	resource        *go3mf.MeshResource
	triangleDecoder triangleDecoder
}

func (d *trianglesDecoder) Open() error {
	d.triangleDecoder.resource = d.resource

	if len(d.resource.Mesh.Faces) == 0 && len(d.resource.Mesh.Nodes) > 0 {
		d.resource.Mesh.Faces = make([]mesh.Face, 0, len(d.resource.Mesh.Nodes)-1)
	}
	return nil
}

func (d *trianglesDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrTriangle {
		child = &d.triangleDecoder
	}
	return
}

type triangleDecoder struct {
	emptyDecoder
	resource *go3mf.MeshResource
}

func (d *triangleDecoder) Attributes(attrs []xml.Attr) (err error) {
	var v1, v2, v3, pid, p1, p2, p3 uint64
	var hasPID, hasP1, hasP2, hasP3 bool
	for _, a := range attrs {
		switch a.Name.Local {
		case attrV1:
			v1, err = strconv.ParseUint(a.Value, 10, 32)
		case attrV2:
			v2, err = strconv.ParseUint(a.Value, 10, 32)
		case attrV3:
			v3, err = strconv.ParseUint(a.Value, 10, 32)
		case attrPID:
			pid, err = strconv.ParseUint(a.Value, 10, 32)
			hasPID = true
		case attrP1:
			p1, err = strconv.ParseUint(a.Value, 10, 32)
			hasP1 = true
		case attrP2:
			p2, err = strconv.ParseUint(a.Value, 10, 32)
			hasP2 = true
		case attrP3:
			p3, err = strconv.ParseUint(a.Value, 10, 32)
			hasP3 = true
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		return
	}

	p1 = applyDefault(p1, uint64(d.resource.DefaultPropertyIndex), hasP1)
	p2 = applyDefault(p2, p1, hasP2)
	p3 = applyDefault(p3, p1, hasP3)
	pid = applyDefault(pid, uint64(d.resource.DefaultPropertyID), hasPID)

	return d.addTriangle(uint32(v1), uint32(v2), uint32(v3), uint32(pid), uint32(p1), uint32(p2), uint32(p3))
}

func (d *triangleDecoder) addTriangle(v1, v2, v3, pid, p1, p2, p3 uint32) error {
	if v1 == v2 || v1 == v3 || v2 == v3 {
		return errors.New("go3mf: invalid triangle indices")
	}
	nodeCount := uint32(len(d.resource.Mesh.Nodes))
	if v1 >= nodeCount || v2 >= nodeCount || v3 >= nodeCount {
		return errors.New("go3mf: triangle index is out of range")
	}
	d.resource.Mesh.Faces = append(d.resource.Mesh.Faces, mesh.Face{
		NodeIndices: [3]uint32{v1, v2, v3},
		Resource: pid,
		ResourceIndices: [3]uint32{p1, p2, p3},
	})
	return nil
}

func applyDefault(val, defVal uint64, noDef bool) uint64 {
	if noDef {
		return val
	}
	return defVal
}
