package geo

import (
	"math"
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

// Point3D defines a node of a mesh as an array of 3 coordinates: x, y and z.
type Point3D [3]float32

// X returns the x coordinate.
func (v1 Point3D) X() float32 {
	return v1[0]
}

// Y returns the y coordinate.
func (v1 Point3D) Y() float32 {
	return v1[1]
}

// Z returns the z coordinate.
func (v1 Point3D) Z() float32 {
	return v1[2]
}

// Add performs element-wise addition between two vectors.
func (v1 Point3D) Add(v2 Point3D) Point3D {
	return Point3D{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}

// Sub performs element-wise subtraction between two vectors.
func (v1 Point3D) Sub(v2 Point3D) Point3D {
	return Point3D{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

// Len returns the vector's length. Note that this is NOT the dimension of
// the vector (len(v)), but the mathematical length.
func (v1 Point3D) Len() float32 {
	return float32(math.Sqrt(float64(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2])))

}

// Normalize normalizes the vector. Normalization is (1/|v|)*v,
// making this equivalent to v.Scale(1/v.Len()). If the len is 0.0,
// this function will return an infinite value for all elements due
// to how floating point division works in Go (n/0.0 = math.Inf(Sign(n))).
func (v1 Point3D) Normalize() Point3D {
	l := 1.0 / v1.Len()
	return Point3D{v1[0] * l, v1[1] * l, v1[2] * l}
}

// Cross product is most often used for finding surface normals. The cross product of vectors
// will generate a vector that is perpendicular to the plane they form.
func (v1 Point3D) Cross(v2 Point3D) Point3D {
	return Point3D{v1[1]*v2[2] - v1[2]*v2[1], v1[2]*v2[0] - v1[0]*v2[2], v1[0]*v2[1] - v1[1]*v2[0]}
}

type nodeStructure struct {
	Nodes      []Point3D
	vectorTree *vectorTree
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
