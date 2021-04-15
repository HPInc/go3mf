// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package go3mf

import (
	"fmt"
	"math"
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

// Matrix is a 4x4 matrix in row major order.
//
// m[4*r + c] is the element in the r'th row and c'th column.
type Matrix [16]float32

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

func (b Box) extend(v Box) Box {
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

func (b Box) extendPoint(v Point3D) Box {
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
