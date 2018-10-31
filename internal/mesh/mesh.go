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
	InformationHandler meshinfo.Handler
}

// NewMesh creates a new default Mesh.
func NewMesh() *Mesh {
	return &Mesh{
		BeamLattice:        *NewBeamLattice(),
		InformationHandler: meshinfo.NewHandler(),
	}
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

func (m *Mesh) GetNode(index uint32) *Node {
	return m.Nodes[uint32(index)]
}

func (m *Mesh) GetFace(index uint32) *Face {
	return m.Faces[uint32(index)]
}

func (m *Mesh) AddNode(position mgl32.Vec3) *Node {
	node := &Node{
		Index:    m.NodeCount(),
		Position: position,
	}
	m.Nodes = append(m.Nodes, node)
	return node
}

func (m *Mesh) AddFace(node1, node2, node3 *Node) *Face {
	face := &Face{
		Index:       m.FaceCount(),
		NodeIndices: [3]uint32{node1.Index, node2.Index, node3.Index},
	}
	m.Faces = append(m.Faces, face)
	m.InformationHandler.AddFace(m.FaceCount())
	return face
}

func (m *Mesh) Merge(mesh *Mesh, matrix mgl32.Mat4) {
	m.InformationHandler.AddInfoFromTable(mesh.InformationHandler, m.FaceCount())
	nodeCount := mesh.NodeCount()
	faceCount := mesh.FaceCount()
	beamCount := mesh.BeamCount()

	if nodeCount == 0 {
		return
	}

	newNodes := make([]*Node, nodeCount)
	for i := 0; i < int(nodeCount); i++ {
		node := m.GetNode(uint32(i))
		position := mgl32.TransformCoordinate(node.Position, matrix)
		newNodes[i] = m.AddNode(position)
	}

	if faceCount != 0 {
		for i := 0; i < int(faceCount); i++ {
			face := mesh.GetFace(uint32(i))
			newFace := m.AddFace(newNodes[face.NodeIndices[0]], newNodes[face.NodeIndices[1]], newNodes[face.NodeIndices[2]])
			m.InformationHandler.CloneFaceInfosFrom(newFace.Index, mesh.InformationHandler, face.Index)
		}
	}

	if beamCount != 0 {

	}
}
