package io3mf

import (
	"encoding/xml"
	"errors"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/mesh"
	"github.com/qmuntal/go3mf/mesh/meshinfo"
	"strconv"
)

type meshDecoder struct {
	r                               *Reader
	resource                        go3mf.MeshResource
	colorMapping                    *colorMapping
	texCoordMapping                 *texCoordMapping
	defaultPropID, defaultPropIndex uint64
	triangleCounter, vertexCounter  int
}

func (d *meshDecoder) Decode(x xml.TokenReader) error {
	d.resource.Mesh = mesh.NewMesh()
	d.resource.Mesh.StartCreation(mesh.CreationOptions{CalculateConnectivity: false})
	defer d.resource.Mesh.EndCreation()
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			var err error
			if tp.Name.Space == nsCoreSpec {
				if tp.Name.Local == attrVertices {
					err = d.parseVertices(x)
				} else if tp.Name.Local == attrTriangles {
					err = d.parseTriangles(x)
				}
			}
			if err != nil {
				return err
			}
		case xml.EndElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrMesh {
				d.r.addResource(&d.resource)
				return nil
			}
		}
	}
}

func (d *meshDecoder) parseVertices(x xml.TokenReader) error {
	if len(d.resource.Mesh.Nodes)%readVerticesUpdate == readVerticesUpdate-1 {
		d.vertexCounter++
		if !d.r.progress.progress(0.5-1.0/float64(d.vertexCounter+2), StageReadMesh) {
			return ErrUserAborted
		}
	}
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrVertex {
				if err := d.parseVertex(tp.Attr); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrVertices {
				return nil
			}
		}
	}
}

func (d *meshDecoder) parseVertex(attrs []xml.Attr) (err error) {
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
		d.resource.Mesh.AddNode(mgl32.Vec3{float32(x), float32(y), float32(z)})
	}
	return
}

func (d *meshDecoder) parseTriangles(x xml.TokenReader) error {
	if len(d.resource.Mesh.Faces)%readTrianglesUpdate == readTrianglesUpdate-1 {
		d.triangleCounter++
		if !d.r.progress.progress(1.0-1.0/float64(d.triangleCounter+2), StageReadMesh) {
			return ErrUserAborted
		}
	}
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrTriangle {
				if err := d.parseTriangle(tp.Attr); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrTriangles {
				return nil
			}
		}
	}
}

func (d *meshDecoder) parseTriangle(attrs []xml.Attr) (err error) {
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

	p1 = applyDefault(p1, d.resource.DefaultPropertyIndex, hasP1)
	p2 = applyDefault(p2, p1, hasP2)
	p3 = applyDefault(p3, p1, hasP3)
	pid = applyDefault(pid, d.resource.DefaultPropertyID, hasPID)

	return d.addTriangle(uint32(v1), uint32(v2), uint32(v3), pid, p1, p2, p3)
}

func applyDefault(val, defVal uint64, noDef bool) uint64 {
	if noDef {
		return val
	}
	return defVal
}

func (d *meshDecoder) addTriangle(v1, v2, v3 uint32, pid, p1, p2, p3 uint64) error {
	if v1 == v2 || v1 == v3 || v2 == v3 {
		return errors.New("go3mf: invalid triangle indices")
	}
	nodeCount := uint32(len(d.resource.Mesh.Nodes))
	if v1 >= nodeCount || v2 >= nodeCount || v3 >= nodeCount {
		return errors.New("go3mf: triangle index is out of range")
	}
	face := d.resource.Mesh.AddFace(v1, v2, v3)
	if pid == 0 {
		return nil
	}
	_ = d.checkColor(face, pid, p1, p2, p3) || d.checkBaseMaterial(face, pid, p1) || d.checkTexture(face, pid, p1, p2, p3)
	return nil
}

func (d *meshDecoder) checkBaseMaterial(face *mesh.Face, pid, p1 uint64) bool {
	ref, ok := d.r.Model.FindResource(pid, d.resource.ModelPath)
	if ok {
		if _, ok := ref.(*go3mf.BaseMaterialsResource); ok {
			handler := d.resource.Mesh.InformationHandler()
			var info *meshinfo.FacesData
			if info, ok = handler.BaseMaterialInfo(); !ok {
				info = handler.AddBaseMaterialInfo(uint32(len(d.resource.Mesh.Faces)))
			}
			faceData := info.FaceData(face.Index).(*meshinfo.BaseMaterial)
			faceData.GroupID = uint32(pid)
			faceData.Index = uint32(p1)
		}
		return true
	}
	return false
}

func (d *meshDecoder) checkColor(face *mesh.Face, pid, p1, p2, p3 uint64) (ok bool) {
	if d.colorMapping.hasResource(pid) {
		handler := d.resource.Mesh.InformationHandler()
		var info *meshinfo.FacesData
		if info, ok = handler.NodeColorInfo(); !ok {
			info = handler.AddNodeColorInfo(uint32(len(d.resource.Mesh.Faces)))
		}
		faceData := info.FaceData(face.Index).(*meshinfo.NodeColor)
		faceData.Colors[0], _ = d.colorMapping.find(pid, p1)
		faceData.Colors[1], _ = d.colorMapping.find(pid, p2)
		faceData.Colors[2], _ = d.colorMapping.find(pid, p3)
		ok = true
	}
	return
}

func (d *meshDecoder) checkTexture(face *mesh.Face, pid, p1, p2, p3 uint64) (ok bool) {
	if d.texCoordMapping.hasResource(pid) {
		handler := d.resource.Mesh.InformationHandler()
		var info *meshinfo.FacesData
		if info, ok = handler.TextureCoordsInfo(); !ok {
			info = handler.AddTextureCoordsInfo(uint32(len(d.resource.Mesh.Faces)))
		}
		faceData := info.FaceData(face.Index).(*meshinfo.TextureCoords)
		t0, _ := d.texCoordMapping.find(pid, p1)
		t1, _ := d.texCoordMapping.find(pid, p2)
		t2, _ := d.texCoordMapping.find(pid, p3)
		faceData.TextureID = uint32(t0.id)
		faceData.Coords[0] = mgl32.Vec2{t0.u, t0.v}
		faceData.Coords[1] = mgl32.Vec2{t1.u, t1.v}
		faceData.Coords[2] = mgl32.Vec2{t2.u, t2.v}
		ok = true
	}
	return
}
