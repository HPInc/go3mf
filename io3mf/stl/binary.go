package stl

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/qmuntal/go3mf"
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
func (d *binaryDecoder) decode(ctx context.Context, m *go3mf.MeshResource) error {
	mb := go3mf.NewMeshBuilder(m)
	var header binaryHeader
	err := binary.Read(d.r, binary.LittleEndian, &header)
	if err != nil {
		return err
	}

	nextFaceCheck := checkEveryFaces
	var facet binaryFace
	for nFace := 0; nFace < int(header.FaceCount); nFace++ {
		err = binary.Read(d.r, binary.LittleEndian, &facet)
		if err != nil {
			break
		}
		d.decodeFace(&facet, mb)
		if len(m.Faces) > nextFaceCheck {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				break
			default: // Default is must to avoid blocking
			}
			nextFaceCheck += checkEveryFaces
		}
	}

	return err
}

func (d *binaryDecoder) decodeFace(facet *binaryFace, mb *go3mf.MeshBuilder) {
	var nodes [3]uint32
	for nVertex := 0; nVertex < 3; nVertex++ {
		pos := facet.Vertices[nVertex]
		nodes[nVertex] = mb.AddNode(go3mf.Point3D{pos[0], pos[1], pos[2]})
	}
	mb.Mesh.Faces = append(mb.Mesh.Faces, go3mf.Face{
		NodeIndices: [3]uint32{nodes[0], nodes[1], nodes[2]},
	})
}

type binaryEncoder struct {
	w io.Writer
}

func (e *binaryEncoder) encode(m *go3mf.MeshResource) error {
	header := binaryHeader{FaceCount: uint32(len(m.Faces))}
	err := binary.Write(e.w, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	for _, f := range m.Faces {
		n1, n2, n3 := m.Nodes[f.NodeIndices[0]], m.Nodes[f.NodeIndices[1]], m.Nodes[f.NodeIndices[2]]
		normal := faceNormal(n1, n2, n3)
		facet := binaryFace{
			Normal:   [3]float32{normal[0], normal[1], normal[2]},
			Vertices: [3][3]float32{{n1.X(), n1.Y(), n1.Z()}, {n2.X(), n2.Y(), n2.Z()}, {n3.X(), n3.Y(), n3.Z()}},
		}
		err := binary.Write(e.w, binary.LittleEndian, facet)
		if err != nil {
			return err
		}
	}
	return nil
}
