package geo

// Matrix is a 4x4 matrix in row major order.
//
// m[4*r + c] is the element in the r'th row and c'th column.
type Matrix [16]float32

// Identity returns the 4x4 identity matrix.
// The identity matrix is a square matrix with the value 1 on its
// diagonals. The characteristic property of the identity matrix is that
// any matrix multiplied by it is itself. (MI = M; IN = N)
func Identity() Matrix {
	return Matrix{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// CreationOptions defines a set of options for helping in the mesh creation process
type CreationOptions struct {
	// True to automatically check if a node with the same coordinates already exists in the mesh
	// when calling AddNode. If it exists, the return value will be the existing node and no node will be added.
	// Using this option produces an speed penalty.
	CalculateConnectivity bool
}

// Mesh is not really a mesh, since it lacks the component edges and the
// topological information. It only holds the nodes and the faces (triangles).
// Each node,  and face have a ID, which allows to identify them. Each face have an
// orientation (i.e. the face can look up or look down) and have three nodes.
// The orientation is defined by the order of its nodes.
type Mesh struct {
	nodeStructure
	faceStructure
	beamLattice
}

// StartCreation can be called before populating the mesh.
// If so, the connectivity will be automatically calculated but producing and speed penalty.
// When the creationg process is finished EndCreation() must be called in order to clean temporary data.
func (m *Mesh) StartCreation(opts CreationOptions) {
	if opts.CalculateConnectivity {
		m.nodeStructure.vectorTree = newVectorTree()
	}
}

// EndCreation cleans temporary data associated to creating a mesh.
func (m *Mesh) EndCreation() {
	m.nodeStructure.vectorTree = nil
}

// Merge merges the mesh with another mesh. This includes the nodes, faces, beams and all the informations.
func (m *Mesh) Merge(mesh *Mesh, matrix Matrix) {
	m.StartCreation(CreationOptions{CalculateConnectivity: true})
	defer m.EndCreation()

	newNodes := m.nodeStructure.merge(&mesh.nodeStructure, matrix)
	if len(newNodes) == 0 {
		return
	}
	m.faceStructure.merge(&mesh.faceStructure, newNodes)
	m.beamLattice.merge(&mesh.beamLattice, newNodes)
}

// CheckSanity checks if the mesh is well formated.
func (m *Mesh) CheckSanity() bool {
	return m.faceStructure.checkSanity(uint32(len(m.Nodes))) && m.beamLattice.checkSanity(uint32(len(m.Nodes)))
}

// FaceNodes returns the three nodes of a face.
func (m *Mesh) FaceNodes(i uint32) (*Point3D, *Point3D, *Point3D) {
	face := m.Faces[i]
	return &m.Nodes[face.NodeIndices[0]], &m.Nodes[face.NodeIndices[1]], &m.Nodes[face.NodeIndices[2]]
}

// IsManifoldAndOriented returns true if the mesh is manifold and oriented.
func (m *Mesh) IsManifoldAndOriented() bool {
	if len(m.Nodes) < 3 || len(m.Faces) < 3 || !m.CheckSanity() {
		return false
	}

	var edgeCounter uint32
	pairMatching := newPairMatch()
	for _, face := range m.Faces {
		for j := uint32(0); j < 3; j++ {
			n1, n2 := face.NodeIndices[j], face.NodeIndices[(j+1)%3]
			if _, ok := pairMatching.CheckMatch(n1, n2); !ok {
				pairMatching.AddMatch(n1, n2, edgeCounter)
				edgeCounter++
			}
		}
	}

	positive, negative := make([]uint32, edgeCounter), make([]uint32, edgeCounter)
	for _, face := range m.Faces {
		for j := uint32(0); j < 3; j++ {
			n1, n2 := face.NodeIndices[j], face.NodeIndices[(j+1)%3]
			edgeIndex, _ := pairMatching.CheckMatch(n1, n2)
			if n1 <= n2 {
				positive[edgeIndex]++
			} else {
				negative[edgeIndex]++
			}
		}
	}

	for i := uint32(0); i < edgeCounter; i++ {
		if positive[i] != 1 || negative[i] != 1 {
			return false
		}
	}

	return true
}

type pairEntry struct {
	a, b uint32
}

// pairMatch implements a n-log-n tree class which is able to identify
// duplicate pairs of numbers in a given data set.
type pairMatch struct {
	entries map[pairEntry]uint32
}

func newPairMatch() *pairMatch {
	return &pairMatch{map[pairEntry]uint32{}}
}

// AddMatch adds a match to the set.
// If the match exists it is overridden.
func (t *pairMatch) AddMatch(data1, data2, param uint32) {
	t.entries[newPairEntry(data1, data2)] = param
}

// CheckMatch check if a match is in the set.
func (t *pairMatch) CheckMatch(data1, data2 uint32) (val uint32, ok bool) {
	val, ok = t.entries[newPairEntry(data1, data2)]
	return
}

// DeleteMatch deletes a match from the set.
// If match doesn't exist it bevavhe as a no-op
func (t *pairMatch) DeleteMatch(data1, data2 uint32) {
	delete(t.entries, newPairEntry(data1, data2))
}

func newPairEntry(data1, data2 uint32) pairEntry {
	entry := pairEntry{}
	if data1 < data2 {
		entry.a = data1
		entry.b = data2
	} else {
		entry.a = data2
		entry.b = data1
	}
	return entry
}
