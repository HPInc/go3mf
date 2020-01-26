package geo

// CreationOptions defines a set of options for helping in the mesh creation process
type CreationOptions struct {
	// True to automatically check if a node with the same coordinates already exists in the mesh
	// when calling AddNode. If it exists, the return value will be the existing node and no node will be added.
	// Using this option produces an speed penalty.
	CalculateConnectivity bool
}

// Face defines a triangle of a mesh.
type Face struct {
	NodeIndices     [3]uint32 // Coordinates of the three nodes that defines the face.
	Resource        uint32
	ResourceIndices [3]uint32 // Resource subindex of the three nodes that defines the face.
}

// BeamSet defines a set of beams.
type BeamSet struct {
	Refs       []uint32
	Name       string
	Identifier string
}

// A CapMode is an enumerable for the different capping modes.
type CapMode uint8

const (
	// CapModeSphere when the capping is an sphere.
	CapModeSphere CapMode = iota
	// CapModeHemisphere when the capping is an hemisphere.
	CapModeHemisphere
	// CapModeButt when the capping is an butt.
	CapModeButt
)

func (b CapMode) String() string {
	return map[CapMode]string{
		CapModeSphere:     "sphere",
		CapModeHemisphere: "hemisphere",
		CapModeButt:       "butt",
	}[b]
}

// Beam defines a single beam.
type Beam struct {
	NodeIndices [2]uint32  // Indices of the two nodes that defines the beam.
	Radius      [2]float64 // Radius of both ends of the beam.
	CapMode     [2]CapMode // Capping mode.
}

// Mesh is not really a mesh, since it lacks the component edges and the
// topological information. It only holds the nodes and the faces (triangles).
// Each node,  and face have a ID, which allows to identify them. Each face have an
// orientation (i.e. the face can look up or look down) and have three nodes.
// The orientation is defined by the order of its nodes.
type Mesh struct {
	Nodes                    []Point3D
	Faces                    []Face
	Beams                    []Beam
	BeamSets                 []BeamSet
	MinLength, DefaultRadius float64
	CapMode                  CapMode
	vectorTree               vectorTree
}

// StartCreation can be called before populating the mesh.
// If so, the connectivity will be automatically calculated but producing and speed penalty.
// When the creationg process is finished EndCreation() must be called in order to clean temporary data.
func (m *Mesh) StartCreation(opts CreationOptions) {
	if opts.CalculateConnectivity {
		m.vectorTree = vectorTree{}
	}
}

// EndCreation cleans temporary data associated to creating a mesh.
func (m *Mesh) EndCreation() {
	m.vectorTree = nil
}

// AddNode adds a node the the mesh at the target position.
func (m *Mesh) AddNode(node Point3D) uint32 {
	if m.vectorTree != nil {
		if index, ok := m.vectorTree.FindVector(node); ok {
			return index
		}
	}
	m.Nodes = append(m.Nodes, node)
	index := uint32(len(m.Nodes)) - 1
	if m.vectorTree != nil {
		m.vectorTree.AddVector(node, index)
	}
	return index
}

// CheckSanity checks if the mesh is well formated.
func (m *Mesh) CheckSanity() bool {
	return m.checkFacesSanity() && m.checkBeamsSanity()
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

func (m *Mesh) checkBeamsSanity() bool {
	nodeCount := uint32(len(m.Nodes))
	for _, beam := range m.Beams {
		i0, i1 := beam.NodeIndices[0], beam.NodeIndices[1]
		if i0 == i1 {
			return false
		}
		if i0 >= nodeCount || i1 >= nodeCount {
			return false
		}
	}
	return true
}

func (m *Mesh) checkFacesSanity() bool {
	nodeCount := uint32(len(m.Nodes))
	for _, face := range m.Faces {
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
