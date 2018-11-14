package stl

import (
	"encoding/binary"
	"io"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/internal/mesh"
)

type binaryHeader struct {
	_   [80]byte
	FaceCount uint32
}

type binaryFace struct {
	_        [3]float32
	Vertices [3][3]float32
	_        uint16
}

// binaryDecoder can create a Mesh from a Read stream that is feeded with a binary STL.
type binaryDecoder struct {
	r     io.Reader
	units float32 // Units of the stream where 1.0 mean meters.
}

// decode loads a binary stl from a io.Reader.
func (d *binaryDecoder) decode() (*mesh.Mesh, error) {
	newMesh := mesh.NewMesh()
	err := newMesh.StartCreation(d.units)
	defer newMesh.EndCreation()
	if err != nil {
		return nil, err
	}

	var header binaryHeader
	err = binary.Read(d.r, binary.LittleEndian, &header)
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
	var nodes [3]*mesh.Node
	for nVertex := 0; nVertex < 3; nVertex++ {
		pos := facet.Vertices[nVertex]
		nodes[nVertex] = newMesh.AddNode(mgl32.Vec3{pos[0], pos[1], pos[2]})
	}

	newMesh.AddFace(nodes[0], nodes[1], nodes[2])
}
