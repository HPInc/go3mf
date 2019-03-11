package mesh

import (
	"math"
	"errors"
	"github.com/go-gl/mathgl/mgl32"
)

// vec3I represents a 3D vector typed as int32
type vec3I struct {
	X int32 // X coordinate
	Y int32 // Y coordinate
	Z int32 // Z coordinate
}

const micronsAccuracy = 1E-6

func newvec3IFromVec3(vec mgl32.Vec3) vec3I {
	a := vec3I{
		X: int32(math.Floor(float64(vec.X() / micronsAccuracy))),
		Y: int32(math.Floor(float64(vec.Y() / micronsAccuracy))),
		Z: int32(math.Floor(float64(vec.Z() / micronsAccuracy))),
	}
	return a
}

// vectorTree implements a n*log(n) lookup tree class to identify vectors by their position
type vectorTree struct {
	entries map[vec3I]uint32
}

func newVectorTree() *vectorTree {
	return &vectorTree{
		entries: make(map[vec3I]uint32),
	}
}

// AddVector adds a vector to the dictionary.
func (t *vectorTree) AddVector(vec mgl32.Vec3, value uint32) {
	t.entries[newvec3IFromVec3(vec)] = value
}

// FindVector returns the identifier of the vector.
func (t *vectorTree) FindVector(vec mgl32.Vec3) (val uint32, ok bool) {
	val, ok = t.entries[newvec3IFromVec3(vec)]
	return
}

// RemoveVector removes the vector from the dictionary.
func (t *vectorTree) RemoveVector(vec mgl32.Vec3) {
	delete(t.entries, newvec3IFromVec3(vec))
}

// MaxNodeCount is the maximum number of nodes allowed.
const MaxNodeCount = 2147483646

// Node defines a node of a mesh.
type Node struct {
	Index    uint32     // Index of the node inside the mesh.
	Position mgl32.Vec3 // Coordinates of the node.
}

type nodeStructure struct {
	vectorTree   *vectorTree
	nodes        []Node
	maxNodeCount uint32 // If 0 MaxNodeCount will be used.
}

func (n *nodeStructure) clear() {
	n.nodes = make([]Node, 0)
}

// NodeCount returns the number of nodes in the mesh.
func (n *nodeStructure) NodeCount() uint32 {
	return uint32(len(n.nodes))
}

// Node retrieve the node with the target index.
func (n *nodeStructure) Node(index uint32) *Node {
	return &n.nodes[uint32(index)]
}

// AddNode adds a node the the mesh at the target position.
func (n *nodeStructure) AddNode(position mgl32.Vec3) *Node {
	if n.vectorTree != nil {
		if index, ok := n.vectorTree.FindVector(position); ok {
			return n.Node(index)
		}
	}
	nodeCount := n.NodeCount()
	if nodeCount >= n.getMaxNodeCount() {
		panic(errors.New("go3mf: too many nodes has been tried to add to a mesh"))
	}

	n.nodes = append(n.nodes, Node{
		Index:    nodeCount,
		Position: position,
	})
	if n.vectorTree != nil {
		n.vectorTree.AddVector(position, nodeCount)
	}
	return &n.nodes[len(n.nodes)-1]
}

func (n *nodeStructure) checkSanity() bool {
	nodeCount := n.NodeCount()
	if nodeCount > n.getMaxNodeCount() {
		return false
	}
	for i := uint32(0); i < nodeCount; i++ {
		if n.Node(i).Index != i {
			return false
		}
	}
	return true
}

func (n *nodeStructure) merge(other mergeableNodes, matrix mgl32.Mat4) []uint32 {
	nodeCount := other.NodeCount()
	newNodes := make([]uint32, nodeCount)
	if nodeCount == 0 {
		return newNodes
	}

	for i := uint32(0); i < nodeCount; i++ {
		node := other.Node(i)
		position := mgl32.TransformCoordinate(node.Position, matrix)
		newNodes[i] = n.AddNode(position).Index
	}
	return newNodes
}

func (n *nodeStructure) getMaxNodeCount() uint32 {
	if n.maxNodeCount == 0 {
		return MaxNodeCount
	}
	return n.maxNodeCount
}
