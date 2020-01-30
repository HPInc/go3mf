package slices

import (
	"errors"

	"github.com/qmuntal/go3mf"
)

// ExtensionName is the canonical name of this extension.
const ExtensionName = "http://schemas.microsoft.com/3dmanufacturing/slice/2015/07"

// Slice defines the resource object for slices.
type Slice struct {
	Vertices []go3mf.Point2D
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
	s.Vertices = append(s.Vertices, go3mf.Point2D{x, y})
	return len(s.Vertices) - 1
}

// AddPolygonIndex adds a new index to the polygon.
func (s *Slice) AddPolygonIndex(polygonIndex, index int) error {
	if polygonIndex >= len(s.Polygons) {
		return errors.New("invalid polygon index")
	}

	if index >= len(s.Vertices) {
		return errors.New("invalid slice segment index")
	}

	p := s.Polygons[polygonIndex]
	if len(p) > 0 && p[len(p)-1] == index {
		return errors.New("duplicated slice segment index")
	}
	s.Polygons[polygonIndex] = append(s.Polygons[polygonIndex], index)
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

// SliceResolution defines the resolutions for a slice.
type SliceResolution uint8

// Supported slice resolution.
const (
	ResolutionFull SliceResolution = iota
	ResolutionLow
)

func newSliceResolution(s string) (r SliceResolution, ok bool) {
	r, ok = map[string]SliceResolution{
		"fullres": ResolutionFull,
		"lowres":  ResolutionLow,
	}[s]
	return
}

func (c SliceResolution) String() string {
	return map[SliceResolution]string{
		ResolutionFull: "fullres",
		ResolutionLow:  "lowres",
	}[c]
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
	Slices  []*Slice
	Refs    []SliceRef
}

// AddSlice adds an slice to the stack and returns its index.
func (s *SliceStack) AddSlice(slice *Slice) (int, error) {
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

// SliceStackInfo defines the attributes added to <object>.
type SliceStackInfo struct {
	SliceStackID    uint32
	SliceResolution SliceResolution
}

// ObjectSliceStackInfo extracts the SliceStackInfo attributes from an ObjectResource.
// If it does not exist a new one is added.
func ObjectSliceStackInfo(o *go3mf.ObjectResource) *SliceStackInfo {
	if attr, ok := o.Extensions[ExtensionName]; ok {
		return attr.(*SliceStackInfo)
	}
	if o.Extensions == nil {
		o.Extensions = make(go3mf.Extensions)
	}
	attr := &SliceStackInfo{}
	o.Extensions[ExtensionName] = attr
	return attr
}

const (
	attrSliceStack = "slicestack"
	attrID         = "id"
	attrZBottom    = "zbottom"
	attrSlice      = "slice"
	attrSliceRef   = "sliceref"
	attrZTop       = "ztop"
	attrVertices   = "vertices"
	attrVertex     = "vertex"
	attrPolygon    = "polygon"
	attrX          = "x"
	attrY          = "y"
	attrZ          = "z"
	attrSegment    = "segment"
	attrV1         = "v1"
	attrV2         = "v2"
	attrV3         = "v3"
	attrStartV     = "startv"
	attrSliceRefID = "slicestackid"
	attrSlicePath  = "slicepath"
	attrMeshRes    = "meshresolution"
)
