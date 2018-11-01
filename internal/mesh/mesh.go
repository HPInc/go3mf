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
	Nodes              []*Node
	Faces              []*Face
	BeamLattice        BeamLattice
	informationHandler *meshinfo.Handler
}

// NewMesh creates a new default Mesh.
func NewMesh() *Mesh {
	return &Mesh{
		BeamLattice: *NewBeamLattice(),
	}
}

func (m *Mesh) Clear() {
	m.informationHandler = nil
	m.Nodes = make([]*Node, 0)
	m.Faces = make([]*Face, 0)
	m.ClearBeamLattice()
}

func (m *Mesh) InformationHandler() *meshinfo.Handler {
	return m.informationHandler
}

func (m *Mesh) CreateInformationHandler() *meshinfo.Handler {
	if m.informationHandler == nil {
		m.informationHandler = meshinfo.NewHandler()
	}

	return m.informationHandler
}

func (m *Mesh) ClearInformationHandler() {
	m.informationHandler = nil
}

func (m *Mesh) ClearBeamLattice() {
	m.BeamLattice.Clear()
}

func (m *Mesh) FaceCount() uint32 {
	return uint32(len(m.Faces))
}

func (m *Mesh) NodeCount() uint32 {
	return uint32(len(m.Nodes))
}

func (m *Mesh) BeamCount() uint32 {
	return uint32(len(m.BeamLattice.Beams))
}

func (m *Mesh) Node(index uint32) *Node {
	return m.Nodes[uint32(index)]
}

func (m *Mesh) Face(index uint32) *Face {
	return m.Faces[uint32(index)]
}

func (m *Mesh) BeamSet(index uint32) *BeamSet {
	return m.BeamLattice.BeamSets[int(index)]
}

func (m *Mesh) SetBeamLatticeMinLength(minLength float64) {
	m.BeamLattice.SetMinLength(minLength)
}

func (m *Mesh) BeamLatticeMinLength() float64 {
	return m.BeamLattice.MinLength()
}

func (m *Mesh) SetBeamLatticeCapMode(capMode BeamCapMode) {
	m.BeamLattice.SetCapMode(capMode)
}

func (m *Mesh) BeamLatticeCapMode() BeamCapMode {
	return m.BeamLattice.CapMode()
}

func (m *Mesh) SetDefaultBeamRadius(radius float64) {
	m.BeamLattice.SetRadius(radius)
}

func (m *Mesh) BeamLatticeRadius() float64 {
	return m.BeamLattice.Radius()
}

func (m *Mesh) AddNode(position mgl32.Vec3) (*Node, error) {
	if position.X() > MaxCoordinate || position.Y() > MaxCoordinate || position.Z() > MaxCoordinate {
		return nil, &MaxCoordinateError{position}
	}

	nodeCount := m.NodeCount()
	if nodeCount > MaxNodeCount {
		return nil, new(MaxNodeError)
	}

	node := &Node{
		Index:    nodeCount,
		Position: position,
	}
	m.Nodes = append(m.Nodes, node)
	return node, nil
}

func (m *Mesh) AddFace(node1, node2, node3 *Node) (*Face, error) {
	if (node1 == node2) || (node1 == node3) || (node2 == node3) {
		return nil, new(DuplicatedNodeError)
	}

	faceCount := m.FaceCount()
	if faceCount > MaxFaceCount {
		return nil, new(MaxFaceError)
	}

	face := &Face{
		Index:       faceCount,
		NodeIndices: [3]uint32{node1.Index, node2.Index, node3.Index},
	}
	m.Faces = append(m.Faces, face)
	if m.informationHandler != nil {
		m.informationHandler.AddFace(m.FaceCount())
	}
	return face, nil
}

func (m *Mesh) AddBeam(node1, node2 *Node, radius1, radius2 float64, capMode1, capMode2 BeamCapMode) (*Beam, error) {
	if node1 == node2 {
		return nil, new(DuplicatedNodeError)
	}

	beamCount := m.BeamCount()
	if beamCount > MaxBeamCount {
		return nil, new(MaxBeamError)
	}

	beam := &Beam{
		NodeIndices: [2]uint32{node1.Index, node2.Index},
		Index:       beamCount,
		Radius:      [2]float64{radius1, radius2},
		CapMode:     [2]BeamCapMode{capMode1, capMode2},
	}

	return beam, nil
}

func (m *Mesh) AddBeamSet() *BeamSet {
	set := new(BeamSet)
	m.BeamLattice.BeamSets = append(m.BeamLattice.BeamSets, set)
	return set
}

func (m *Mesh) Merge(mesh *Mesh, matrix mgl32.Mat4) {
	m.informationHandler.AddInfoFromTable(mesh.informationHandler, m.FaceCount())
	nodeCount := mesh.NodeCount()
	faceCount := mesh.FaceCount()
	beamCount := mesh.BeamCount()

	if nodeCount == 0 {
		return
	}

	newNodes := make([]*Node, nodeCount)
	for i := 0; i < int(nodeCount); i++ {
		node := m.Node(uint32(i))
		position := mgl32.TransformCoordinate(node.Position, matrix)
		newNodes[i], _ = m.AddNode(position)
	}

	if faceCount != 0 {
		for i := 0; i < int(faceCount); i++ {
			face := mesh.Face(uint32(i))
			newFace, _ := m.AddFace(newNodes[face.NodeIndices[0]], newNodes[face.NodeIndices[1]], newNodes[face.NodeIndices[2]])
			m.informationHandler.CloneFaceInfosFrom(newFace.Index, mesh.informationHandler, face.Index)
		}
	}

	if beamCount != 0 {

	}
}
