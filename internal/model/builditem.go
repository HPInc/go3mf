package model

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gofrs/uuid"
	"github.com/qmuntal/go3mf/internal/mesh"
)

// A BuildItem is an in memory representation of the 3MF build item.
type BuildItem struct {
	Object     Object
	Transform  mgl32.Mat4
	PartNumber string
	Path       string
	uuid       uuid.UUID
}

// UUID returns the object UUID.
func (b *BuildItem) UUID() uuid.UUID {
	return b.uuid
}

// SetUUID sets the object UUID
func (b *BuildItem) SetUUID(id uuid.UUID) error {
	err := registerUUID(b.uuid, id, b.Object.RootModel())
	if err == nil {
		b.uuid = id
	}
	return err
}

// HasTransform returns true if the transform is different than the identity.
func (b *BuildItem) HasTransform() bool {
	return !b.Transform.ApproxEqual(identityTransform)
}

// IsValidForSlices checks if the build object is valid to be used with slices.
func (b *BuildItem) IsValidForSlices() bool {
	return b.Object.IsValidForSlices(b.Transform)
}

// MergeToMesh merges the build object with the mesh.
func (b *BuildItem) MergeToMesh(m *mesh.Mesh) error {
	return b.Object.MergeToMesh(m, b.Transform)
}