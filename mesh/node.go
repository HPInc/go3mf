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

func newvec3IFromVec3(vec Point3D) vec3I {
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
func (t *vectorTree) AddVector(vec Point3D, value uint32) {
	t.entries[newvec3IFromVec3(vec)] = value
}

// FindVector returns the identifier of the vector.
func (t *vectorTree) FindVector(vec Point3D) (val uint32, ok bool) {
	val, ok = t.entries[newvec3IFromVec3(vec)]
	return
}

// RemoveVector removes the vector from the dictionary.
func (t *vectorTree) RemoveVector(vec Point3D) {
	delete(t.entries, newvec3IFromVec3(vec))
}

const maxNodeCount = 2147483646

// Point3D defines a node of a mesh as an array of 3 coordinates: x, y and z.
type Point3D [3]float32

// X returns the x coordinate.
func (n Point3D) X() float32 {
	return n[0]
}

// Y returns the y coordinate.
func (n Point3D) Y() float32 {
	return n[1]
}

// Z returns the z coordinate.
func (n Point3D) Z() float32 {
	return n[2]
}

type nodeStructure struct {
	Nodes        []Point3D
	vectorTree   *vectorTree
	maxNodeCount int
}

// AddNode adds a node the the mesh at the target position.
func (n *nodeStructure) AddNode(node Point3D) uint32 {
	if n.vectorTree != nil {
		if index, ok := n.vectorTree.FindVector(node); ok {
			return index
		}
	}
	n.Nodes = append(n.Nodes, node)
	index := uint32(len(n.Nodes)) - 1
	if n.vectorTree != nil {
		n.vectorTree.AddVector(node, index)
	}
	return index
}

func (n *nodeStructure) checkSanity() bool {
	return len(n.Nodes) <= n.getMaxNodeCount()
}

func (n *nodeStructure) merge(other *nodeStructure, matrix Matrix) []uint32 {
	newNodes := make([]uint32, len(other.Nodes))
	if len(other.Nodes) == 0 {
		return newNodes
	}

	m := mgl32.Mat4(matrix)
	for i := range other.Nodes {
		position := mgl32.TransformCoordinate(mgl32.Vec3(other.Nodes[i]), m)
		newNodes[i] = n.AddNode(Point3D(position))
	}
	return newNodes
}

func (n *nodeStructure) getMaxNodeCount() int {
	if n.maxNodeCount == 0 {
		return maxNodeCount
	}
	return n.maxNodeCount
}
