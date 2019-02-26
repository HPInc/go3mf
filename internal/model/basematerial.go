package model

import (
	"fmt"
	"image/color"
)

// BaseMaterial defines the Model Base Material Resource.
// A model material resource is an in memory representation of the 3MF
// material resource object.
type BaseMaterial struct {
	Name  string
	Color color.RGBA
}

// ColorString returns the color as a hex string with the format #rrggbbaa.
func (m *BaseMaterial) ColorString() string {
	return fmt.Sprintf("#%x%x%x%x", m.Color.R, m.Color.G, m.Color.B, m.Color.A)
}

// BaseMaterialsResource defines a slice of BaseMaterial.
type BaseMaterialsResource struct {
	Resource
	Materials []*BaseMaterial
}

// NewBaseMaterialsResource returns a new texture 2D resource.
func NewBaseMaterialsResource(id uint64, model *Model) (*BaseMaterialsResource, error) {
	r, err := newResource(id, model)
	if err != nil {
		return nil, err
	}
	return &BaseMaterialsResource{
		Resource: *r,
	}, nil
}

// Merge appends all the other base materials.
func (ms *BaseMaterialsResource) Merge(other []*BaseMaterial) {
	for _, m := range other {
		ms.Materials = append(ms.Materials, &BaseMaterial{m.Name, m.Color})
	}
}
