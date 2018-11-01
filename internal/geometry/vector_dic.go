package geometry

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// VectorDic implements a n*log(n) lookup tree class to identify vectors by their position
// The units property defines the units of the vectors, where 1.0 mean meters.
type VectorDic struct {
	units   float32
	entries map[Vec3I]uint32
}

// NewVectorDic creates a default vector
func NewVectorDic() *VectorDic {
	return &VectorDic{
		units:   VectorDefaultUnits,
		entries: map[Vec3I]uint32{},
	}
}

// NewVectorDicWithUnits creates a vector with the desired units
// Error cases: See SetUnits
func NewVectorDicWithUnits(units float32) (*VectorDic, error) {
	t := &VectorDic{
		entries: map[Vec3I]uint32{},
	}
	err := t.SetUnits(units)
	return t, err
}

// Units returns the used units.
func (t *VectorDic) Units() float32 {
	return t.units
}

// SetUnits sets the used units.
// Error cases:
// * ErrorInvalidUnits: ((units < VectorMinUnits) || (units > VectorMaxUnits))
// * ErrorCouldNotSetUnits: non-empty tree
func (t *VectorDic) SetUnits(units float32) error {
	if (units < VectorMinUnits) || (units > VectorMaxUnits) {
		return &InvalidUnitsError{units}
	}
	if len(t.entries) > 0 {
		return new(UnitsNotSettedError)
	}
	t.units = units
	return nil
}

// AddVector adds a vector to the dictionary.
func (t *VectorDic) AddVector(vec mgl32.Vec3, value uint32) {
	t.entries[newVec3IFromVec3(vec, t.units)] = value
}

// FindVector returns the identifier of the vector.
func (t *VectorDic) FindVector(vec mgl32.Vec3) (val uint32, ok bool) {
	val, ok = t.entries[newVec3IFromVec3(vec, t.units)]
	return
}

// RemoveVector removes the vector from the dictionary.
func (t *VectorDic) RemoveVector(vec mgl32.Vec3) {
	delete(t.entries, newVec3IFromVec3(vec, t.units))
}

func newVec3IFromVec3(vec mgl32.Vec3, units float32) Vec3I {
	return Vec3I{
		X: int32(math.Floor(float64(vec.X() / units))),
		Y: int32(math.Floor(float64(vec.Y() / units))),
		Z: int32(math.Floor(float64(vec.Z() / units))),
	}
}
