package geometry

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// VectorTree implements a n*log(n) lookup tree class to identify vectors by their position
// The units property defines the units of the vectors, where 1.0 mean meters.
type VectorTree struct {
	units   float32
	entries map[Vec3I]uint32
}

// NewVectorTree creates a default vector
func NewVectorTree() *VectorTree {
	return &VectorTree{
		units:   VectorDefaultUnits,
		entries: map[Vec3I]uint32{},
	}
}

// Units returns the used units.
func (t *VectorTree) Units() float32 {
	if t.units == 0 {
		return VectorDefaultUnits
	}
	return t.units
}

// SetUnits sets the used units.
// Error cases:
// * ErrorInvalidUnits: ((units < VectorMinUnits) || (units > VectorMaxUnits))
// * ErrorCouldNotSetUnits: non-empty tree
func (t *VectorTree) SetUnits(units float32) error {
	if units == 0 {
		units = VectorDefaultUnits
	} else if (units < VectorMinUnits) || (units > VectorMaxUnits) {
		return &InvalidUnitsError{units}
	}
	if len(t.entries) > 0 {
		return new(UnitsNotSettedError)
	}
	t.units = units
	return nil
}

// AddVector adds a vector to the dictionary.
func (t *VectorTree) AddVector(vec mgl32.Vec3, value uint32) {
	t.entries[newVec3IFromVec3(vec, t.Units())] = value
}

// FindVector returns the identifier of the vector.
func (t *VectorTree) FindVector(vec mgl32.Vec3) (val uint32, ok bool) {
	val, ok = t.entries[newVec3IFromVec3(vec, t.Units())]
	return
}

// RemoveVector removes the vector from the dictionary.
func (t *VectorTree) RemoveVector(vec mgl32.Vec3) {
	delete(t.entries, newVec3IFromVec3(vec, t.Units()))
}

func newVec3IFromVec3(vec mgl32.Vec3, units float32) Vec3I {
	return Vec3I{
		X: int32(math.Floor(float64(vec.X() / units))),
		Y: int32(math.Floor(float64(vec.Y() / units))),
		Z: int32(math.Floor(float64(vec.Z() / units))),
	}
}
