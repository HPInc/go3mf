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
	informationHandler *meshinfo.Handler
}

// NewMesh creates a new default Mesh.
func NewMesh() *Mesh {
	m := &Mesh{
		beamLattice: *newbeamLattice(),
	}
	m.faceStructure.informationHandler = m.informationHandler
	return m
}

// NewMeshCloned creates a new mesh that is a clone of another mesh.
func NewMeshCloned(mesh MergeableMesh) (*Mesh, error) {
	m := NewMesh()
	m.Merge(mesh, mgl32.Ident4())
	return m, nil
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
	return m.informationHandler
}

// CreateInformationHandler creates a new information handler if it has not been created yet.
func (m *Mesh) CreateInformationHandler() *meshinfo.Handler {
	if m.informationHandler == nil {
		m.informationHandler = meshinfo.NewHandler()
	}

	return m.informationHandler
}

// ClearInformationHandler sets the information handler to nil.
func (m *Mesh) ClearInformationHandler() {
	m.informationHandler = nil
}

// Merge merges the mesh with another mesh. This includes the nodes, faces, beams and all the informations.
func (m *Mesh) Merge(mesh MergeableMesh, matrix mgl32.Mat4) error {
	otherHandler := mesh.InformationHandler()
	if otherHandler != nil {
		m.CreateInformationHandler()
		m.informationHandler.AddInfoFrom(otherHandler, m.FaceCount())
	}

	newNodes, err := m.nodeStructure.merge(mesh, matrix)
	if len(newNodes) == 0 || err != nil {
		return err
	}

	err = m.faceStructure.merge(mesh, newNodes)
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
