package mesh

import (
	"errors"
	"github.com/qmuntal/go3mf/internal/meshinfo"
)

const maxFaceCount = 2147483646

// Face defines a triangle of a mesh.
type Face struct {
	Index       uint32    // Index of the face inside the mesh.
	NodeIndices [3]uint32 // Coordinates of the three nodes that defines the mesh.
}

type faceStructure struct {
	faces              []Face
	informationHandler *meshinfo.Handler
	maxFaceCount       uint32
}

func (f *faceStructure) clear() {
	f.faces = make([]Face, 0)
}

// FaceCount returns the number of faces in the mesh.
func (f *faceStructure) FaceCount() uint32 {
	return uint32(len(f.faces))
}

// Face retrieve the face with the target index.
func (f *faceStructure) Face(index uint32) *Face {
	return &f.faces[uint32(index)]
}

// AddFace adds a face to the mesh that has the target nodes.
func (f *faceStructure) AddFace(node1, node2, node3 uint32) (*Face, error) {
	if (node1 == node2) || (node1 == node3) || (node2 == node3) {
		return nil, errors.New("go3mf: a beam with two identical nodes has been tried to add to a mesh")
	}

	faceCount := f.FaceCount()
	if faceCount >= f.getMaxFaceCount() {
		panic(errors.New("go3mf: a face with too many face has been tried to add to a meshs"))
	}

	f.faces = append(f.faces, Face{
		Index:       faceCount,
		NodeIndices: [3]uint32{node1, node2, node3},
	})
	if f.informationHandler != nil {
		f.informationHandler.AddFace(f.FaceCount())
	}
	return &f.faces[len(f.faces)-1], nil
}

func (f *faceStructure) checkSanity(nodeCount uint32) bool {
	faceCount := f.FaceCount()
	if faceCount > f.getMaxFaceCount() {
		return false
	}
	for i := uint32(0); i < faceCount; i++ {
		face := f.Face(i)
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

func (f *faceStructure) merge(other MergeableFaces, newNodes []uint32) error {
	faceCount := other.FaceCount()
	if faceCount == 0 {
		return nil
	}
	otherHandler := other.InformationHandler()
	for i := uint32(0); i < faceCount; i++ {
		face := other.Face(i)
		newFace, err := f.AddFace(newNodes[face.NodeIndices[0]], newNodes[face.NodeIndices[1]], newNodes[face.NodeIndices[2]])
		if err != nil {
			return err
		}
		if otherHandler != nil {
			f.informationHandler.CopyFaceInfosFrom(newFace.Index, otherHandler, face.Index)
		}
	}
	return nil
}

func (f *faceStructure) getMaxFaceCount() uint32 {
	if f.maxFaceCount == 0 {
		return maxFaceCount
	}
	return f.maxFaceCount
}
