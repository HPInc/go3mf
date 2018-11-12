package meshimporter

import (
	"bytes"
	"encoding/binary"
	"io"
	"unsafe"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/qmuntal/go3mf/internal/geometry"
	"github.com/qmuntal/go3mf/internal/mesh"
)

type stlBinaryFace struct {
	normal [3]float32
	vertices [3][3]float32
	attribute uint16
}

// STLBinary can create a Mesh from a Read stream that is feeded with a binary STL.
// The struct is idempontent so can be reused for different streams and goroutines.
type STLBinary struct {
	Units              float32 // Units of the stream where 1.0 mean meters.
	IgnoreInvalidFaces bool    // True to ignore invalid faces, false to do a fast fail.
}

// LoadMesh loads a binary stl from a io.Reader.
func (s *STLBinary) LoadMesh(stream io.Reader) (*mesh.Mesh, error) {
	newMesh := mesh.NewMesh()
	vectorTree := geometry.NewVectorTree()
	err := vectorTree.SetUnits(s.Units)
	if err != nil {
		return nil, err
	}

	// Read header
	buff := make([]byte, 80)
	_, err = stream.Read(buff)
	if err != nil {
		return nil, err
	}

	var faceCount uint32
	err = s.readBytes(stream, faceCount)
	if err != nil {
		return nil, err
	}

	for nFace := 0; nFace < int(faceCount); nFace++ {
		var facet stlBinaryFace 
		err = s.readBytes(stream, facet)
		if err != nil {
			return nil, err
		}

		var nodes [3]*mesh.Node
		for nVertex := 0; nVertex< 3; nVertex++ {
			pos := facet.vertices[nVertex]
			vec := mgl32.Vec3{pos[0], pos[1], pos[2]}
			if index, ok := vectorTree.FindVector(vec); ok {
				nodes[nVertex] = newMesh.Node(index)
			} else {
				newNode := newMesh.AddNode(vec)
				vectorTree.AddVector(newNode.Position, newNode.Index)
				nodes[nVertex] = newNode
			}
		}

		_, err := newMesh.AddFace(nodes[0], nodes[1], nodes[2])
		if err != nil && !s.IgnoreInvalidFaces {
			return nil, err
		}
	}

	return newMesh, nil
}

func (s *STLBinary) readBytes(stream io.Reader, data interface{}) error {
	buff := make([]byte, unsafe.Sizeof(data))
	_, err := stream.Read(buff)
	if err != nil {
		return err
	}
	buf := bytes.NewReader(buff)
	return binary.Read(buf, binary.LittleEndian, &data)
}
