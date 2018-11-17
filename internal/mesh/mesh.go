package mesh

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/internal/meshinfo"
	"github.com/qmuntal/go3mf/internal/geometry"
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

// StartCreation can be called before populating the mesh. 
// If so, the connectivity will be autmatically calculated but producing and speed penalty.
// When the creationg process is finished EndCreation() must be called in order to clean temporary data.
func (m *Mesh) StartCreation(units float32) error {
	m.nodeStructure.vectorTree = geometry.NewVectorTree()
	return m.nodeStructure.vectorTree.SetUnits(units)
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

// ApproxEqual compares the geometry of two meshes to check if they are equal.The mesh properties are not compared.
func (m *Mesh) ApproxEqual(mesh *Mesh) bool {
	return comparer{}.CompareGeometry(m, mesh)
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

// FaceNodes returns the three nodes of a face
func (m *Mesh) FaceNodes(i uint32) (*Node, *Node, *Node) {
	face := m.Face(uint32(i))
	return m.Node(face.NodeIndices[0]), m.Node(face.NodeIndices[1]), m.Node(face.NodeIndices[2])
}