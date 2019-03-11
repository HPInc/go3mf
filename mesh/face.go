package mesh

import (
	"github.com/qmuntal/go3mf/mesh/meshinfo"
)

const maxFaceCount = 2147483646

// Face defines a triangle of a mesh.
type Face struct {
	Index       uint32    // Index of the face inside the mesh.
	NodeIndices [3]uint32 // Coordinates of the three nodes that defines the mesh.
}

type faceStructure struct {
	Faces              []Face
	informationHandler *meshinfo.Handler
	maxFaceCount       int
}

func (f *faceStructure) clear() {
	f.Faces = make([]Face, 0)
}

// AddFace adds a face to the mesh that has the target nodes.
func (f *faceStructure) AddFace(node1, node2, node3 uint32) *Face {
	f.Faces = append(f.Faces, Face{
		Index:       uint32(len(f.Faces)),
		NodeIndices: [3]uint32{node1, node2, node3},
	})
	if f.informationHandler != nil {
		f.informationHandler.AddFace(uint32(len(f.Faces)))
	}
	return &f.Faces[len(f.Faces)-1]
}

func (f *faceStructure) checkSanity(nodeCount uint32) bool {
	faceCount := len(f.Faces)
	if faceCount > f.getMaxFaceCount() {
		return false
	}
	for _, face := range f.Faces {
		i0, i1, i2 := face.NodeIndices[0], face.NodeIndices[1], face.NodeIndices[2]
		if i0 == i1 || i0 == i2 || i1 == i2 {
			return false
		}
		if i0 >= nodeCount || i1 >= nodeCount || i2 >= nodeCount {
			return false
		}
	}
	return true
}

func (f *faceStructure) merge(other *faceStructure, newNodes []uint32) {
	faceCount := len(other.Faces)
	if faceCount == 0 {
		return
	}
	otherHandler := other.informationHandler
	for _, face := range other.Faces {
		newFace := f.AddFace(newNodes[face.NodeIndices[0]], newNodes[face.NodeIndices[1]], newNodes[face.NodeIndices[2]])
		if otherHandler != nil {
			f.informationHandler.CopyFaceInfosFrom(newFace.Index, otherHandler, face.Index)
		}
	}
	return
}

func (f *faceStructure) getMaxFaceCount() int {
	if f.maxFaceCount == 0 {
		return maxFaceCount
	}
	return f.maxFaceCount
}
