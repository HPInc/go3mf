package geometry

import (
	"math"

	"github.com/qmuntal/go3mf/internal/common"
)

// Vector2D is defined as a vector in the 2D space using float32
type Vector2D struct {
	X float32 // X coordinate
	Y float32 // Y coordinate
}

// NewVector2D created a new Vector2D
func NewVector2D(x, y float32) Vector2D {
	return Vector2D{x, y}
}

// Add returns a new vector that is the sum of both vectors
func (a Vector2D) Add(b Vector2D) Vector2D {
	return Vector2D{a.X + b.X, a.Y + b.Y}
}

func (a Vector2D) Sub(b Vector2D) Vector2D {
	return Vector2D{a.X - b.X, a.Y - b.Y}
}

func (a Vector2D) Scale(b float32) Vector2D {
	return Vector2D{a.X * b, a.Y * b}
}

func (a Vector2D) Combine(factor1 float32, b Vector2D, factor2 float32) Vector2D {
	return Vector2D{a.X*factor1 + b.X*factor2, a.Y*factor1 + b.Y*factor2}
}

func (a Vector2D) Dot(b Vector2D) float32 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vector2D) Cross(b Vector2D) float32 {
	return a.X*b.Y - a.Y*b.X
}

func (a Vector2D) Length() float32 {
	return float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y)))
}

func (a Vector2D) Distance(b Vector2D) float32 {
	return a.Sub(b).Length()
}

func (a Vector2D) Normalize() (Vector2D, error) {
	l := a.Length()
	if l < VectorMinNormalizeLength {
		return Vector2D{}, common.NewError(common.ErrorNormalizedZeroVector)
	}
	return a.Scale(1.0/l), nil
}

func (a Vector2D) Floor(units float32) (Vector2DI, error) {
	if units < VectorMinUnits || units > VectorMaxUnits {
		return Vector2DI{}, common.NewError(common.ErrorInvalidUnits)
	}
	return NewVector2DI(int32(math.Floor(float64(a.X/units))), int32(math.Floor(float64(a.Y/units)))), nil
}
