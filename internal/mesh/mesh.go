package mesh

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/internal/meshinfo"
)

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
		beamLattice: *newbeamLattice(),
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
	m.ClearBeamLattice()
}

// InformationHandler returns the information handler of the mesh.
// If CreateInformationHandler() has not been called, it will always be nil.
func (m *Mesh) InformationHandler() *meshinfo.Handler {
	return &m.informationHandler
}

// ClearInformationHandler sets the information handler to nil.
func (m *Mesh) ClearInformationHandler() {
	m.informationHandler.RemoveAllInformations()
}

func (m *Mesh) Equal(mesh *Mesh) bool {
	if len(m.nodes) != len(mesh.nodes) {
		return false
	}
	if len(m.faces) != len(mesh.faces) {
		return false
	}
	for i := 0; i < len(m.nodes); i++ {
		if !m.nodes[i].Position.ApproxEqualThreshold(mesh.nodes[i].Position, 0.0001) {
			return false
		}
	}
	for i := 0; i < len(m.faces); i++ {
		indices := m.faces[i].NodeIndices
		other := mesh.faces[i].NodeIndices
		if indices[0] != other[0] || indices[1] != other[1] || indices[2] != other[2] {
			return false
		}
	}
	return true
}

// Merge merges the mesh with another mesh. This includes the nodes, faces, beams and all the informations.
func (m *Mesh) Merge(mesh MergeableMesh, matrix mgl32.Mat4) error {
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
