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
	ID        uint64
	Materials []*BaseMaterial
	modelPath string
	uniqueID  uint64
}

// ResourceID returns the resource ID, which has the same value as ID.
func (ms *BaseMaterialsResource) ResourceID() uint64 {
	return ms.ID
}

// UniqueID returns the unique ID.
func (ms *BaseMaterialsResource) UniqueID() uint64 {
	return ms.uniqueID
}

func (ms *BaseMaterialsResource) setUniqueID(id uint64) {
	ms.uniqueID = id
}

// Merge appends all the other base materials.
func (ms *BaseMaterialsResource) Merge(other []*BaseMaterial) {
	for _, m := range other {
		ms.Materials = append(ms.Materials, &BaseMaterial{m.Name, m.Color})
	}
}
