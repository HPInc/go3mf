package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
)

// VectorTree implements a n*log(n) lookup tree class to identify vectors by their position
type VectorTree struct {
	entries map[vec3I]uint32
}

// NewVectorTree creates a default vector
func NewVectorTree() *VectorTree {
	return &VectorTree{
		entries: map[vec3I]uint32{},
	}
}

// AddVector adds a vector to the dictionary.
func (t *VectorTree) AddVector(vec mgl32.Vec3, value uint32) {
	t.entries[newvec3IFromVec3(vec)] = value
}

// FindVector returns the identifier of the vector.
func (t *VectorTree) FindVector(vec mgl32.Vec3) (val uint32, ok bool) {
	val, ok = t.entries[newvec3IFromVec3(vec)]
	return
}

// RemoveVector removes the vector from the dictionary.
func (t *VectorTree) RemoveVector(vec mgl32.Vec3) {
	delete(t.entries, newvec3IFromVec3(vec))
}
