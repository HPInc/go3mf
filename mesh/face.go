package mesh

const maxFaceCount = 2147483646

// Face defines a triangle of a mesh.
type Face struct {
	NodeIndices     [3]uint32 // Coordinates of the three nodes that defines the face.
	Resource        uint32
	ResourceIndices [3]uint32 // Resource subindex of the three nodes that defines the face.
}

type faceStructure struct {
	Faces        []Face
	maxFaceCount int
}

// AddFace adds a face to the mesh that has the target nodes.
func (f *faceStructure) AddFace(node1, node2, node3 uint32) *Face {
	f.Faces = append(f.Faces, Face{
		NodeIndices: [3]uint32{node1, node2, node3},
	})
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
	for _, face := range other.Faces {
		f.Faces = append(f.Faces, Face{
			NodeIndices:     [3]uint32{face.NodeIndices[0], newNodes[face.NodeIndices[1]], newNodes[face.NodeIndices[2]]},
			Resource:        face.Resource,
			ResourceIndices: [3]uint32{face.ResourceIndices[0], newNodes[face.ResourceIndices[1]], newNodes[face.ResourceIndices[2]]},
		})
	}
	return
}

func (f *faceStructure) getMaxFaceCount() int {
	if f.maxFaceCount == 0 {
		return maxFaceCount
	}
	return f.maxFaceCount
}
