package mesh

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/internal/meshinfo"
)

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
	informationHandler meshinfo.Handler
}

// NewMesh creates a new default Mesh.
func NewMesh() *Mesh {
	m := &Mesh{
		beamLattice:        *newbeamLattice(),
		informationHandler: *meshinfo.NewHandler(),
	}
	m.faceStructure.informationHandler = &m.informationHandler
	return m
}

// Clone creates a deep clone of the mesh.
func (m *Mesh) Clone() (*Mesh, error) {
	new := NewMesh()
	err := new.Merge(m, mgl32.Ident4())
	return new, err
}

// Clear resets all the mesh nodes, faces, beams and informations.
func (m *Mesh) Clear() {
	m.ClearInformationHandler()
	m.nodeStructure.clear()
	m.faceStructure.clear()
	m.clearBeamLattice()
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

// InformationHandler returns the information handler of the mesh.
func (m *Mesh) InformationHandler() *meshinfo.Handler {
	return &m.informationHandler
}

// ClearInformationHandler sets the information handler to nil.
func (m *Mesh) ClearInformationHandler() {
	m.informationHandler.RemoveAllInformations()
}

// Merge merges the mesh with another mesh. This includes the nodes, faces, beams and all the informations.
func (m *Mesh) Merge(mesh MergeableMesh, matrix mgl32.Mat4) error {
	m.StartCreation(CreationOptions{CalculateConnectivity: true})
	defer m.EndCreation()
	m.informationHandler.AddInfoFrom(mesh.InformationHandler(), m.FaceCount())

	newNodes := m.nodeStructure.merge(mesh, matrix)
	if len(newNodes) == 0 {
		return nil
	}

	err := m.faceStructure.merge(mesh, newNodes)
	if err != nil {
		return err
	}

	return m.beamLattice.merge(mesh, newNodes)
}

// CheckSanity checks if the mesh is well formated.
func (m *Mesh) CheckSanity() bool {
	if !m.nodeStructure.checkSanity() {
		return false
	}
	if !m.faceStructure.checkSanity(m.NodeCount()) {
		return false
	}
	return m.beamLattice.checkSanity(m.NodeCount())
}

// FaceNodes returns the three nodes of a face.
func (m *Mesh) FaceNodes(i uint32) (*Node, *Node, *Node) {
	face := m.Face(uint32(i))
	return m.Node(face.NodeIndices[0]), m.Node(face.NodeIndices[1]), m.Node(face.NodeIndices[2])
}

// FaceNormal returns the normal of a face.
func (m *Mesh) FaceNormal(i uint32) mgl32.Vec3 {
	node1, node2, node3 := m.FaceNodes(i)
	return faceNormal(node1.Position, node2.Position, node3.Position)
}

func faceNormal(n1, n2, n3 mgl32.Vec3) mgl32.Vec3 {
	return n2.Sub(n1).Cross(n3.Sub(n1)).Normalize()
}

// IsManifoldAndOriented returns true if the mesh is manifold and oriented.
func (m *Mesh) IsManifoldAndOriented() bool {
	if m.NodeCount() < 3 || m.FaceCount() < 3 || !m.CheckSanity() {
		return false
	}

	var edgeCounter uint32
	pairMatching := newPairMatch()
	for i := uint32(0); i < m.FaceCount(); i++ {
		face := m.Face(i)
		for j := uint32(0); j < 3; j++ {
			n1, n2 := face.NodeIndices[j], face.NodeIndices[(j+1)%3]
			if _, ok := pairMatching.CheckMatch(n1, n2); !ok {
				pairMatching.AddMatch(n1, n2, edgeCounter)
				edgeCounter++
			}
		}
	}

	positive, negative := make([]uint32, edgeCounter), make([]uint32, edgeCounter)
	for i := uint32(0); i < m.FaceCount(); i++ {
		face := m.Face(i)
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
