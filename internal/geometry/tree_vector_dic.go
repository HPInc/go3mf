package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/internal/common"
	"math"
)

// TreeVectorDic implements a n*log(n) lookup tree class to identify vectors by their position
// The units property defines the units of the vectors, where 1.0 mean meters.
type TreeVectorDic struct {
	units   float32
	entries map[Vec3I]uint32
}

// NewTreeVectorDic creates a default vector
func NewTreeVectorDic() *TreeVectorDic {
	return &TreeVectorDic{
		units:   VectorDefaultUnits,
		entries: map[Vec3I]uint32{},
	}
}

// NewTreeVectorDicWithUnits creates a vector with the desired units
// Error cases: See SetUnits
func NewTreeVectorDicWithUnits(units float32) (*TreeVectorDic, error) {
	t := &TreeVectorDic{
		entries: map[Vec3I]uint32{},
	}
	err := t.SetUnits(units)
	return t, err
}

// Units returns the used units.
func (t *TreeVectorDic) Units() float32 {
	return t.units
}

// SetUnits sets the used units.
// Error cases:
// * ErrorInvalidUnits: ((units < VectorMinUnits) || (units > VectorMaxUnits))
// * ErrorCouldNotSetUnits: non-empty tree
func (t *TreeVectorDic) SetUnits(units float32) error {
	if (units < VectorMinUnits) || (units > VectorMaxUnits) {
		return common.NewError(common.ErrorInvalidUnits)
	}
	if len(t.entries) > 0 {
		return common.NewError(common.ErrorCouldNotSetUnits)
	}
	t.units = units
	return nil
}

// AddVector adds a vector to the dictionary.
func (t *TreeVectorDic) AddVector(vec mgl32.Vec3, value uint32) {
	t.entries[newVec3IFromVec3(vec, t.units)] = value
}

// FindVector returns the identifier of the vector.
func (t *TreeVectorDic) FindVector(vec mgl32.Vec3) (val uint32, ok bool) {
	val, ok = t.entries[newVec3IFromVec3(vec, t.units)]
	return
}

// RemoveVector removes the vector from the dictionary.
func (t *TreeVectorDic) RemoveVector(vec mgl32.Vec3) {
	delete(t.entries, newVec3IFromVec3(vec, t.units))
}

func newVec3IFromVec3(vec mgl32.Vec3, units float32) Vec3I {
	return Vec3I{
		X: int32(math.Floor(float64(vec.X() / units))),
		Y: int32(math.Floor(float64(vec.Y() / units))),
		Z: int32(math.Floor(float64(vec.Z() / units))),
	}
}
