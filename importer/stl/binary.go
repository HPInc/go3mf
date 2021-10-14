// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package stl

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/hpinc/go3mf"
)

type binaryHeader struct {
	_         [80]byte
	FaceCount uint32
}

type binaryFace struct {
	_        [3]float32
	Vertices [3][3]float32
	_        uint16
}

// binaryDecoder can create a Mesh from a Read stream that is feeded with a binary STL.
type binaryDecoder struct {
	r io.Reader
}

// decode loads a binary stl from a io.Reader.
func (d *binaryDecoder) decode(ctx context.Context, m *go3mf.Mesh) error {
	mb := go3mf.NewMeshBuilder(m)
	var header binaryHeader
	err := binary.Read(d.r, binary.LittleEndian, &header)
	if err != nil {
		return err
	}
	mb.Mesh.Triangles.Triangle = make([]go3mf.Triangle, 0, header.FaceCount)
	nextFaceCheck := checkEveryFaces
	var facet binaryFace
	for nFace := 0; nFace < int(header.FaceCount); nFace++ {
		err = binary.Read(d.r, binary.LittleEndian, &facet)
		if err != nil {
			break
		}
		d.decodeFace(&facet, mb)
		if len(m.Triangles.Triangle) > nextFaceCheck {
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
		nodes[nVertex] = mb.AddVertex(go3mf.Point3D{pos[0], pos[1], pos[2]})
	}
	mb.Mesh.Triangles.Triangle = append(mb.Mesh.Triangles.Triangle, go3mf.Triangle{V1: nodes[0], V2: nodes[1], V3: nodes[2]})
}
