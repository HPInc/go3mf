package mesh

import (
	"math"
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

const maxNodeCount = 2147483646

// Node defines a node of a mesh.
type Node struct {
	Index    uint32     // Index of the node inside the mesh.
	Position mgl32.Vec3 // Coordinates of the node.
}

type nodeStructure struct {
	Nodes        []Node
	vectorTree   *vectorTree
	maxNodeCount int
}

func (n *nodeStructure) clear() {
	n.Nodes = make([]Node, 0)
}

// AddNode adds a node the the mesh at the target position.
func (n *nodeStructure) AddNode(position mgl32.Vec3) *Node {
	if n.vectorTree != nil {
		if index, ok := n.vectorTree.FindVector(position); ok {
			return &n.Nodes[index]
		}
	}
	nodeCount := uint32(len(n.Nodes))
	n.Nodes = append(n.Nodes, Node{
		Index:    nodeCount,
		Position: position,
	})
	if n.vectorTree != nil {
		n.vectorTree.AddVector(position, nodeCount)
	}
	return &n.Nodes[len(n.Nodes)-1]
}

func (n *nodeStructure) checkSanity() bool {
	if len(n.Nodes) > n.getMaxNodeCount() {
		return false
	}
	for i := range n.Nodes {
		if n.Nodes[i].Index != uint32(i) {
			return false
		}
	}
	return true
}

func (n *nodeStructure) merge(other *nodeStructure, matrix mgl32.Mat4) []uint32 {
	newNodes := make([]uint32, len(other.Nodes))
	if len(other.Nodes) == 0 {
		return newNodes
	}

	for i := range other.Nodes {
		position := mgl32.TransformCoordinate(other.Nodes[i].Position, matrix)
		newNodes[i] = n.AddNode(position).Index
	}
	return newNodes
}

func (n *nodeStructure) getMaxNodeCount() int {
	if n.maxNodeCount == 0 {
		return maxNodeCount
	}
	return n.maxNodeCount
}
