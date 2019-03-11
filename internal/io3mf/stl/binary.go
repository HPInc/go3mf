package stl

import (
	"encoding/binary"
	"io"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/internal/mesh"
)

type binaryHeader struct {
	_         [80]byte
	FaceCount uint32
}

type binaryFace struct {
	Normal   [3]float32
	Vertices [3][3]float32
	_        uint16
}

// binaryDecoder can create a Mesh from a Read stream that is feeded with a binary STL.
type binaryDecoder struct {
	r io.Reader
}

// decode loads a binary stl from a io.Reader.
func (d *binaryDecoder) decode() (*mesh.Mesh, error) {
	newMesh := mesh.NewMesh()
	newMesh.StartCreation(mesh.CreationOptions{CalculateConnectivity: true})
	defer newMesh.EndCreation()
	var header binaryHeader
	err := binary.Read(d.r, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	var facet binaryFace
	for nFace := 0; nFace < int(header.FaceCount); nFace++ {
		err = binary.Read(d.r, binary.LittleEndian, &facet)
		if err != nil {
			break
		}
		d.decodeFace(&facet, newMesh)
	}

	return newMesh, err
}

func (d *binaryDecoder) decodeFace(facet *binaryFace, newMesh *mesh.Mesh) {
	var nodes [3]uint32
	for nVertex := 0; nVertex < 3; nVertex++ {
		pos := facet.Vertices[nVertex]
		nodes[nVertex] = newMesh.AddNode(mgl32.Vec3{pos[0], pos[1], pos[2]}).Index
	}

	newMesh.AddFace(nodes[0], nodes[1], nodes[2])
}

type binaryEncoder struct {
	w io.Writer
}

func (e *binaryEncoder) encode(m *mesh.Mesh) error {
	faceCount := uint32(len(m.Faces))
	header := binaryHeader{FaceCount: faceCount}
	err := binary.Write(e.w, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	for i := uint32(0); i < faceCount; i++ {
		node1, node2, node3 := m.FaceNodes(i)
		normal := m.FaceNormal(i)
		n1, n2, n3 := node1.Position, node2.Position, node3.Position
		facet := binaryFace{
			Normal:   [3]float32{normal.X(), normal.Y(), normal.Z()},
			Vertices: [3][3]float32{{n1.X(), n1.Y(), n1.Z()}, {n2.X(), n2.Y(), n2.Z()}, {n3.X(), n3.Y(), n3.Z()}},
		}
		err := binary.Write(e.w, binary.LittleEndian, facet)
		if err != nil {
			return err
		}
	}
	return nil
}
