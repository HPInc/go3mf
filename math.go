package go3mf

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

type pairEntry struct {
	a, b uint32
}

// pairMatch implements a n-log-n tree class which is able to identify
// duplicate pairs of numbers in a given data set.
type pairMatch map[pairEntry]uint32

// AddMatch adds a match to the set.
// If the match exists it is overridden.
func (t pairMatch) AddMatch(data1, data2, param uint32) {
	t[newPairEntry(data1, data2)] = param
}

// CheckMatch check if a match is in the set.
func (t pairMatch) CheckMatch(data1, data2 uint32) (val uint32, ok bool) {
	val, ok = t[newPairEntry(data1, data2)]
	return
}

func newPairEntry(data1, data2 uint32) pairEntry {
	if data1 < data2 {
		return pairEntry{data1, data2}
	}
	return pairEntry{data2, data1}
}

// vec3I represents a 3D vector typed as int32
type vec3I struct {
	X int32 // X coordinate
	Y int32 // Y coordinate
	Z int32 // Z coordinate
}

const micronsAccuracy = 1e-6

func newvec3IFromVec3(vec Point3D) vec3I {
	a := vec3I{
		X: int32(math.Floor(float64(vec.X() / micronsAccuracy))),
		Y: int32(math.Floor(float64(vec.Y() / micronsAccuracy))),
		Z: int32(math.Floor(float64(vec.Z() / micronsAccuracy))),
	}
	return a
}

// vectorTree implements a n*log(n) lookup tree class to identify vectors by their position
type vectorTree map[vec3I]uint32

// AddVector adds a vector to the dictionary.
func (t vectorTree) AddVector(vec Point3D, value uint32) {
	t[newvec3IFromVec3(vec)] = value
}

// FindVector returns the identifier of the vector.
func (t vectorTree) FindVector(vec Point3D) (val uint32, ok bool) {
	val, ok = t[newvec3IFromVec3(vec)]
	return
}

// RemoveVector removes the vector from the dictionary.
func (t vectorTree) RemoveVector(vec Point3D) {
	delete(t, newvec3IFromVec3(vec))
}

// Point2D defines a node of a slice as an array of 2 coordinates: x and y.
type Point2D [2]float32

// X returns the x coordinate.
func (n Point2D) X() float32 {
	return n[0]
}

// Y returns the y coordinate.
func (n Point2D) Y() float32 {
	return n[1]
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

// Normalize normalizes the vector. If the len is 0.0,
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

// Matrix is a 4x4 matrix in row major order.
//
// m[4*r + c] is the element in the r'th row and c'th column.
type Matrix [16]float32

// ParseMatrix parses s as a Matrix.
func ParseMatrix(s string) (Matrix, bool) {
	values := strings.Fields(s)
	if len(values) != 12 {
		return Matrix{}, false
	}
	var t [12]float32
	for i := 0; i < 12; i++ {
		val, err := strconv.ParseFloat(values[i], 32)
		if err != nil {
			return Matrix{}, false
		}
		t[i] = float32(val)
	}
	return Matrix{t[0], t[1], t[2], 0.0,
		t[3], t[4], t[5], 0.0,
		t[6], t[7], t[8], 0.0,
		t[9], t[10], t[11], 1.0}, true
}

// String returns the string representation of a Matrix.
func (m1 Matrix) String() string {
	return fmt.Sprintf("%.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f",
		m1[0], m1[1], m1[2], m1[4], m1[5], m1[6], m1[8], m1[9], m1[10], m1[12], m1[13], m1[14])
}

// Identity returns the 4x4 identity matrix.
// The identity matrix is a square matrix with the value 1 on its
// diagonals. The characteristic property of the identity matrix is that
// any matrix multiplied by it is itself. (MI = M; IN = N)
func Identity() Matrix {
	return Matrix{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// Translate returns a matrix with a relative translation applied.
func (m1 Matrix) Translate(x, y, z float32) Matrix {
	m1[12] += x
	m1[13] += y
	m1[14] += z
	return m1
}

// Mul performs a "matrix product" between this matrix
// and another matrix.
func (m1 Matrix) Mul(m2 Matrix) Matrix {
	return Matrix{
		m1[0]*m2[0] + m1[4]*m2[1] + m1[8]*m2[2] + m1[12]*m2[3],
		m1[1]*m2[0] + m1[5]*m2[1] + m1[9]*m2[2] + m1[13]*m2[3],
		m1[2]*m2[0] + m1[6]*m2[1] + m1[10]*m2[2] + m1[14]*m2[3],
		m1[3]*m2[0] + m1[7]*m2[1] + m1[11]*m2[2] + m1[15]*m2[3],
		m1[0]*m2[4] + m1[4]*m2[5] + m1[8]*m2[6] + m1[12]*m2[7],
		m1[1]*m2[4] + m1[5]*m2[5] + m1[9]*m2[6] + m1[13]*m2[7],
		m1[2]*m2[4] + m1[6]*m2[5] + m1[10]*m2[6] + m1[14]*m2[7],
		m1[3]*m2[4] + m1[7]*m2[5] + m1[11]*m2[6] + m1[15]*m2[7],
		m1[0]*m2[8] + m1[4]*m2[9] + m1[8]*m2[10] + m1[12]*m2[11],
		m1[1]*m2[8] + m1[5]*m2[9] + m1[9]*m2[10] + m1[13]*m2[11],
		m1[2]*m2[8] + m1[6]*m2[9] + m1[10]*m2[10] + m1[14]*m2[11],
		m1[3]*m2[8] + m1[7]*m2[9] + m1[11]*m2[10] + m1[15]*m2[11],
		m1[0]*m2[12] + m1[4]*m2[13] + m1[8]*m2[14] + m1[12]*m2[15],
		m1[1]*m2[12] + m1[5]*m2[13] + m1[9]*m2[14] + m1[13]*m2[15],
		m1[2]*m2[12] + m1[6]*m2[13] + m1[10]*m2[14] + m1[14]*m2[15],
		m1[3]*m2[12] + m1[7]*m2[13] + m1[11]*m2[14] + m1[15]*m2[15],
	}
}

// Mul3D performs a "matrix product" between this matrix
// and another 3D point.
func (m1 Matrix) Mul3D(v Point3D) Point3D {
	return Point3D{
		m1[0]*v[0] + m1[4]*v[1] + m1[8]*v[2] + m1[12],
		m1[1]*v[0] + m1[5]*v[1] + m1[9]*v[2] + m1[13],
		m1[2]*v[0] + m1[6]*v[1] + m1[10]*v[2] + m1[14],
	}
}

// Mul2D performs a "matrix product" between this matrix
// and another 2D point.
func (m1 Matrix) Mul2D(v Point2D) Point2D {
	return Point2D{
		m1[0]*v[0] + m1[4]*v[1] + m1[12],
		m1[1]*v[0] + m1[5]*v[1] + m1[13],
	}
}

// MulBox performs a "matrix product" between this matrix
// and a box
func (m1 Matrix) MulBox(b Box) Box {
	if m1[15] == 0 {
		return b
	}
	box := Box{
		Min: m1.Mul3D(b.Min),
		Max: m1.Mul3D(b.Max),
	}
	if box.Min.X() > box.Max.X() {
		box.Min[0], box.Max[0] = box.Max[0], box.Min[0]
	}
	if box.Min.Y() > box.Max.Y() {
		box.Min[1], box.Max[1] = box.Max[1], box.Min[1]
	}
	if box.Min.Z() > box.Max.Z() {
		box.Min[2], box.Max[2] = box.Max[2], box.Min[2]
	}
	return box
}

// Box defines a box in the 3D space.
type Box struct {
	Min Point3D
	Max Point3D
}

var emptyBox = Box{}

func newLimitBox() Box {
	return Box{
		Min: Point3D{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32},
		Max: Point3D{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32},
	}
}

// Extends adds v to the box.
func (b Box) Extend(v Box) Box {
	return Box{
		Min: Point3D{
			min(b.Min.X(), v.Min.X()),
			min(b.Min.Y(), v.Min.Y()),
			min(b.Min.Z(), v.Min.Z()),
		},
		Max: Point3D{
			max(b.Max.X(), v.Max.X()),
			max(b.Max.Y(), v.Max.Y()),
			max(b.Max.Z(), v.Max.Z()),
		},
	}
}

// ExtendPoint adds v to the box.
func (b Box) ExtendPoint(v Point3D) Box {
	return Box{
		Min: Point3D{
			min(b.Min.X(), v.X()),
			min(b.Min.Y(), v.Y()),
			min(b.Min.Z(), v.Z()),
		},
		Max: Point3D{
			max(b.Max.X(), v.X()),
			max(b.Max.Y(), v.Y()),
			max(b.Max.Z(), v.Z()),
		},
	}
}

func min(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

func max(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

// ParseRGBA parses s as a RGBA color.
func ParseRGBA(s string) (c color.RGBA, err error) {
	var errInvalidFormat = errors.New("go3mf: invalid color format")

	if len(s) == 0 || s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 9:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
		c.A = hexToByte(s[7])<<4 + hexToByte(s[8])
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
		c.A = 0xff
	default:
		err = errInvalidFormat
	}
	return
}

// FormatRGBA returns the color as a hex string with the format #rrggbbaa.
func FormatRGBA(c color.RGBA) string {
	if c.A == 255 {
		return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}
