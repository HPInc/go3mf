package mesh

import (
	"math"

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
		BeamLattice: *newBeamLattice(),
	}
}

// NewMeshCloned creates a new mesh that is a clone of another mesh.
func NewMeshCloned(mesh *Mesh) (*Mesh, error) {
	m := NewMesh()
	m.Merge(mesh, mgl32.Ident4())
	return m, nil
}

// Clear resets all the mesh nodes, faces, beams and informations.
func (m *Mesh) Clear() {
	m.informationHandler = nil
	m.Nodes = make([]*Node, 0)
	m.Faces = make([]*Face, 0)
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

// ClearBeamLattice clears all the beam lattice data.
func (m *Mesh) ClearBeamLattice() {
	m.BeamLattice.Clear()
}

// FaceCount returns the number of faces in the mesh.
func (m *Mesh) FaceCount() uint32 {
	return uint32(len(m.Faces))
}

// NodeCount returns the number of nodes in the mesh.
func (m *Mesh) NodeCount() uint32 {
	return uint32(len(m.Nodes))
}

// BeamCount returns the number of beams in the mesh.
func (m *Mesh) BeamCount() uint32 {
	return uint32(len(m.BeamLattice.Beams))
}

// Node retrieve the node with the target index.
func (m *Mesh) Node(index uint32) *Node {
	return m.Nodes[uint32(index)]
}

// Face retrieve the face with the target index.
func (m *Mesh) Face(index uint32) *Face {
	return m.Faces[uint32(index)]
}

// Beam retrieve the beam with the target index.
func (m *Mesh) Beam(index uint32) *Beam {
	return m.BeamLattice.Beams[uint32(index)]
}

// BeamSet retrieve the beam set with the target index.
func (m *Mesh) BeamSet(index uint32) *BeamSet {
	return m.BeamLattice.BeamSets[int(index)]
}

// SetBeamLatticeMinLength sets the minimum allowed length of a beam.
func (m *Mesh) SetBeamLatticeMinLength(minLength float64) {
	m.BeamLattice.SetMinLength(minLength)
}

// BeamLatticeMinLength returns the minimum allowed length of a beam.
func (m *Mesh) BeamLatticeMinLength() float64 {
	return m.BeamLattice.MinLength()
}

// SetBeamLatticeCapMode sets the capping mode of the beams.
func (m *Mesh) SetBeamLatticeCapMode(capMode BeamCapMode) {
	m.BeamLattice.SetCapMode(capMode)
}

// BeamLatticeCapMode returns the capping mode of the beams.
func (m *Mesh) BeamLatticeCapMode() BeamCapMode {
	return m.BeamLattice.CapMode()
}

// SetDefaultBeamRadius sets the default beam radius.
func (m *Mesh) SetDefaultBeamRadius(radius float64) {
	m.BeamLattice.SetRadius(radius)
}

// DefaultBeamLatticeRadius returns the default beam radius.
func (m *Mesh) DefaultBeamLatticeRadius() float64 {
	return m.BeamLattice.Radius()
}

// AddNode adds a node the the mesh at the target position.
func (m *Mesh) AddNode(position mgl32.Vec3) (*Node, error) {
	x, y, z := math.Abs(float64(position.X())), math.Abs(float64(position.Y())), math.Abs(float64(position.Z()))
	if x > MaxCoordinate || y > MaxCoordinate || z > MaxCoordinate {
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

// AddFace adds a face to the mesh that has the target nodes.
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

// AddBeam adds a beam to the mesh with the desried parameters.
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

// AddBeamSet adds a new beam set to the mesh.
func (m *Mesh) AddBeamSet() *BeamSet {
	set := new(BeamSet)
	m.BeamLattice.BeamSets = append(m.BeamLattice.BeamSets, set)
	return set
}

// Merge merges the mesh with another mesh. This includes the nodes, faces, beams and all the informations.
func (m *Mesh) Merge(mesh *Mesh, matrix mgl32.Mat4) error {
	otherHandler := mesh.InformationHandler()
	if otherHandler != nil {
		m.CreateInformationHandler()
		m.informationHandler.AddInfoFrom(otherHandler, m.FaceCount())
	}

	nodeCount := mesh.NodeCount()
	if nodeCount == 0 {
		return nil
	}

	var err error
	newNodes := make([]*Node, nodeCount)
	for i := 0; i < int(nodeCount); i++ {
		node := m.Node(uint32(i))
		position := mgl32.TransformCoordinate(node.Position, matrix)
		newNodes[i], err = m.AddNode(position)
		if err != nil {
			return err
		}
	}

	faceCount := mesh.FaceCount()
	if faceCount != 0 {
		for i := 0; i < int(faceCount); i++ {
			face := mesh.Face(uint32(i))
			newFace, err := m.AddFace(newNodes[face.NodeIndices[0]], newNodes[face.NodeIndices[1]], newNodes[face.NodeIndices[2]])
			if err != nil {
				return err
			}
			if otherHandler != nil {
				m.informationHandler.CloneFaceInfosFrom(newFace.Index, otherHandler, face.Index)
			}
		}
	}

	beamCount := mesh.BeamCount()
	if beamCount != 0 {
		for i := 0; i < int(beamCount); i++ {
			beam := mesh.Beam(uint32(i))
			n1, n2 := newNodes[beam.NodeIndices[0]], newNodes[beam.NodeIndices[1]]
			_, err = m.AddBeam(n1, n2, beam.Radius[0], beam.Radius[1], beam.CapMode[0], beam.CapMode[1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CheckSanity checks if the mesh is well formated.
func (m *Mesh) CheckSanity() bool {
	nodeCount := m.NodeCount()
	faceCount := m.FaceCount()
	beamCount := m.BeamCount()

	if nodeCount > MaxNodeCount {
		return false
	}
	if faceCount > MaxFaceCount {
		return false
	}
	if beamCount > MaxBeamCount {
		return false
	}
	for i := 0; i < int(nodeCount); i++ {
		node := m.Node(uint32(i))
		if node.Index != uint32(i) {
			return false
		}
		position := node.Position
		x, y, z := math.Abs(float64(position.X())), math.Abs(float64(position.Y())), math.Abs(float64(position.Z()))
		if x > MaxCoordinate || y > MaxCoordinate || z > MaxCoordinate {
			return false
		}
	}
	for i := 0; i < int(faceCount); i++ {
		face := m.Face(uint32(i))
		i0, i1, i2 := face.NodeIndices[0], face.NodeIndices[1], face.NodeIndices[2]
		if i0 == i1 || i0 == i2 || i1 == i2 {
			return false
		}
		if i0 < 0 || i1 < 0 || i2 < 0 || i0 >= nodeCount || i1 >= nodeCount || i2 >= nodeCount {
			return false
		}
	}
	for i := 0; i < int(beamCount); i++ {
		beam := m.Beam(uint32(i))
		i0, i1 := beam.NodeIndices[0], beam.NodeIndices[1]
		if i0 == i1 {
			return false
		}
		if i0 < 0 || i1 < 0 || i0 >= nodeCount || i1 >= nodeCount {
			return false
		}
	}
	return true
}
