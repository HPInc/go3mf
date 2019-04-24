package go3mf

import (
	"errors"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/mesh"
)

// SliceResolution defines the resolutions for a slice.
type SliceResolution uint8

const (
	// ResolutionFull defines a full resolution slice.
	ResolutionFull SliceResolution = iota
	// ResolutionLow defines a low resolution slice.
	ResolutionLow
)

func (c SliceResolution) String() string {
	return map[SliceResolution]string{
		ResolutionFull: "fullres",
		ResolutionLow:  "lowres",
	}[c]
}

// IsValidForSlices checks if the component resource and all its child are valid to be used with slices.
func (c *ComponentsResource) IsValidForSlices(transform mesh.Matrix) bool {
	if len(c.Components) == 0 {
		return true
	}

	matrix := mgl32.Mat4(transform)
	for _, comp := range c.Components {
		if !comp.Object.IsValidForSlices(mesh.Matrix(matrix.Mul4(mgl32.Mat4(comp.Transform)))) {
			return false
		}
	}
	return true
}

// IsValidForSlices checks if the build object is valid to be used with slices.
func (b *BuildItem) IsValidForSlices() bool {
	return b.Object.IsValidForSlices(b.Transform)
}

// IsValidForSlices checks if the mesh resource are valid for slices.
func (c *MeshResource) IsValidForSlices(t mesh.Matrix) bool {
	return c.SliceStackID == 0 || t[2] == 0 && t[6] == 0 && t[8] == 0 && t[9] == 0 && t[10] == 1
}

// SliceRef reference to a slice stack.
type SliceRef struct {
	SliceStackID uint32
	Path         string
}

// SliceStack defines an stack of slices.
// It can either contain Slices or a Refs.
type SliceStack struct {
	BottomZ float32
	Slices  []*mesh.Slice
	Refs    []SliceRef
}

// AddSlice adds an slice to the stack and returns its index.
func (s *SliceStack) AddSlice(slice *mesh.Slice) (int, error) {
	if slice.TopZ < s.BottomZ || (len(s.Slices) != 0 && slice.TopZ < s.Slices[0].TopZ) {
		return 0, errors.New("the z-coordinates of slices within a slicestack are not increasing")
	}
	s.Slices = append(s.Slices, slice)
	return len(s.Slices) - 1, nil
}

// SliceStackResource defines a slice stack resource.
// It can either contain a SliceStack or a Refs slice.
type SliceStackResource struct {
	Stack     SliceStack
	ID        uint32
	ModelPath string
}

// Identify returns the unique ID of the resource.
func (s *SliceStackResource) Identify() (string, uint32) {
	return s.ModelPath, s.ID
}
