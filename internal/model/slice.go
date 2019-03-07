package model

import (
	"errors"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// Slice defines the resource object for slices.
type Slice struct {
	Vertices []mgl32.Vec2
	Polygons [][]int
	TopZ     float32
}

// BeginPolygon adds a new polygon and return its index.
func (s *Slice) BeginPolygon() int {
	s.Polygons = append(s.Polygons, make([]int, 0))
	return len(s.Polygons) - 1
}

// AddVertex adds a new vertex to the slice and returns its index.
func (s *Slice) AddVertex(x, y float32) int {
	s.Vertices = append(s.Vertices, mgl32.Vec2{x, y})
	return len(s.Vertices) - 1
}

// AddPolygonIndex adds a new index to the polygon.
func (s *Slice) AddPolygonIndex(polygonIndex, index int) error {
	if polygonIndex >= len(s.Polygons) {
		return errors.New("go3mf: invalid polygon index")
	}

	if index >= len(s.Vertices) {
		return errors.New("go3mf: invalid slice segment index")
	}

	p := s.Polygons[polygonIndex]
	if len(p) > 0 && p[len(p)-1] == index {
		return errors.New("go3mf: duplicated slice segment index")
	}
	p = append(p, index)
	return nil
}

// AllPolygonsAreClosed returns true if all the polygons are closed.
func (s *Slice) AllPolygonsAreClosed() bool {
	for _, p := range s.Polygons {
		if len(p) > 1 && p[0] != p[len(p)-1] {
			return false
		}
	}
	return true
}

// IsPolygonValid returns true if the polygon is valid.
func (s *Slice) IsPolygonValid(index int) bool {
	if index >= len(s.Polygons) {
		return false
	}
	p := s.Polygons[index]
	return len(p) > 2
}

// SliceStack defines an stack of slices
type SliceStack struct {
	BottomZ      float32
	Slices       []*Slice
	UsesSliceRef bool
}

// AddSlice adds an slice to the stack and returns its index.
func (s *SliceStack) AddSlice(slice *Slice) (int, error) {
	if slice.TopZ < s.BottomZ || (len(s.Slices) != 0 && slice.TopZ < s.Slices[0].TopZ) {
		return 0, errors.New("go3mf: The z-coordinates of slices within a slicestack are not increasing")
	}
	s.Slices = append(s.Slices, slice)
	return len(s.Slices) - 1, nil
}

// SliceStackResource defines a slice stack resource.
type SliceStackResource struct {
	Resource
	*SliceStack
	TimesRefered int
}

// NewSliceStackResource returns a new SliceStackResource.
func NewSliceStackResource(id uint64, model *Model, stack *SliceStack) (*SliceStackResource, error) {
	r, err := newResource(id, model)
	if err != nil {
		return nil, err
	}
	return &SliceStackResource{SliceStack: stack, Resource: *r}, nil
}

// ReferencePath returns the path to the file defining the slice stack and
// empty if UsesSliceRef is false.
func (s *SliceStackResource) ReferencePath() string {
	if s.UsesSliceRef {
		return fmt.Sprintf("/2D/2dmodel_%d.model", s.ResourceID.UniqueID())
	}
	return ""
}
