package geometry

import (
	"math"
	"github.com/qmuntal/go3mf/internal/common"
)

// Vector2DI is defined as a vector in the 2D space using int32
type Vector2DI struct {
	X int32 // X coordinate
	Y int32 // Y coordinate
}

// NewVector2DI created a new Vector2DI
func NewVector2DI(x, y int32) Vector2DI {
	return Vector2DI{x, y}
}

//Cast returns the Vector2D version of the vector.
func (a Vector2DI) Uncast(units float32) (Vector2D, error) {
	if units < VectorMinUnits || units > VectorMaxUnits {
		return Vector2D{}, common.NewError(common.ErrorInvalidUnits)
	}
	return NewVector2D(float32(a.X)*units, float32(a.Y)*units), nil
}

// Add returns a new vector that is the sum of both vectors.
func (a Vector2DI) Add(b Vector2DI) Vector2DI {
	return Vector2DI{a.X + b.X, a.Y + b.Y}
}

// Sub returns a new vector that is the substraction of both vectors.
func (a Vector2DI) Sub(b Vector2DI) Vector2DI {
	return Vector2DI{a.X - b.X, a.Y - b.Y}
}

// Scale returns a the vector scaled by the given factor.
func (a Vector2DI) Scale(b int32) Vector2DI {
	return Vector2DI{a.X * b, a.Y * b}
}

// Dot returns the dot product of both vectors.
func (a Vector2DI) Dot(b Vector2DI) int64 {
	return int64(a.X*b.X + a.Y*b.Y)
}

// Length returns the length of the vector.
func (a Vector2DI) Length() float32 {
	return float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y)))
}

// Distance returns the distance between both vectors.
func (a Vector2DI) Distance(b Vector2DI) float32 {
	return a.Sub(b).Length()
}
