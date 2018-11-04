package mesh

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// MaxNodeCount is the maximum number of nodes allowed.
const MaxNodeCount = 2147483646

// MaxCoordinate is the maximum value of a coordinate.
const MaxCoordinate = 1000000000.0

// Node defines a node of a mesh.
type Node struct {
	Index    uint32     // Index of the node inside the mesh.
	Position mgl32.Vec3 // Coordinates of the node.
}

type nodeStructure struct {
	nodes []*Node
}

func (n *nodeStructure) clear() {
	n.nodes = make([]*Node, 0)
}

// NodeCount returns the number of nodes in the mesh.
func (n *nodeStructure) NodeCount() uint32 {
	return uint32(len(n.nodes))
}

// Node retrieve the node with the target index.
func (n *nodeStructure) Node(index uint32) *Node {
	return n.nodes[uint32(index)]
}

// AddNode adds a node the the mesh at the target position.
func (n *nodeStructure) AddNode(position mgl32.Vec3) (*Node, error) {
	x, y, z := math.Abs(float64(position.X())), math.Abs(float64(position.Y())), math.Abs(float64(position.Z()))
	if x > MaxCoordinate || y > MaxCoordinate || z > MaxCoordinate {
		return nil, &MaxCoordinateError{position}
	}

	nodeCount := n.NodeCount()
	if nodeCount > MaxNodeCount {
		return nil, new(MaxNodeError)
	}

	node := &Node{
		Index:    nodeCount,
		Position: position,
	}
	n.nodes = append(n.nodes, node)
	return node, nil
}

func (n *nodeStructure) checkSanity() bool {
	nodeCount := n.NodeCount()
	if nodeCount > MaxNodeCount {
		return false
	}
	for i := 0; i < int(nodeCount); i++ {
		node := n.Node(uint32(i))
		if node.Index != uint32(i) {
			return false
		}
		position := node.Position
		x, y, z := math.Abs(float64(position.X())), math.Abs(float64(position.Y())), math.Abs(float64(position.Z()))
		if x > MaxCoordinate || y > MaxCoordinate || z > MaxCoordinate {
			return false
		}
	}
	return true
}

func (n *nodeStructure) merge(other MergeableNodes, matrix mgl32.Mat4) ([]*Node, error) {
	nodeCount := other.NodeCount()
	newNodes := make([]*Node, nodeCount)
	if nodeCount == 0 {
		return newNodes, nil
	}

	var err error
	for i := 0; i < int(nodeCount); i++ {
		node := other.Node(uint32(i))
		position := mgl32.TransformCoordinate(node.Position, matrix)
		newNodes[i], err = n.AddNode(position)
		if err != nil {
			return nil, err
		}
	}
	return newNodes, nil
}
