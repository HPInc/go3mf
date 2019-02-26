package model

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gofrs/uuid"
	"github.com/qmuntal/go3mf/internal/mesh"
)

// Object defines a composable object.
type Object interface {
	RootModel() *Model
	MergeToMesh(*mesh.Mesh, mgl32.Mat4)
	ID() uint64
}

// An ObjectResource is an in memory representation of the 3MF model object.
type ObjectResource struct {
	Resource
	Name            string
	PartNumber      string
	SliceStackID    *PackageResourceID
	SliceResoultion SliceResolution
	Thumbnail       string
	DefaultProperty interface{}
	Type            ObjectType
	uuid            uuid.UUID
}

// NewObjectResource returns a new object resource.
func NewObjectResource(id uint64, model *Model) (*ObjectResource, error) {
	r, err := newResource(model, id)
	if err != nil {
		return nil, err
	}
	return &ObjectResource{
		Resource: *r,
	}, nil
}

// UUID returns the object UUID.
func (o *ObjectResource) UUID() uuid.UUID {
	return o.uuid
}

// SetUUID sets the object UUID
func (o *ObjectResource) SetUUID(id uuid.UUID) error {
	err := registerUUID(o.uuid, id, o.Model)
	if err == nil {
		o.uuid = id
	}
	return err
}

// MergeToMesh left on purpose empty to be reimplemented
func (o *ObjectResource) MergeToMesh(m *mesh.Mesh, transform mgl32.Mat4) {
}

// RootModel returns the model of the object.
func (o *ObjectResource) RootModel() *Model {
	return o.Model
}

// ID returns the id of the object.
func (o *ObjectResource) ID() uint64 {
	return o.ResourceID.UniqueID()
}

// A Component is an in memory representation of the 3MF component.
type Component struct {
	Object    Object
	Transform mgl32.Mat4
	uuid      uuid.UUID
}

// UUID returns the object UUID.
func (c *Component) UUID() uuid.UUID {
	return c.uuid
}

// SetUUID sets the object UUID
func (c *Component) SetUUID(id uuid.UUID) error {
	err := registerUUID(c.uuid, id, c.Object.RootModel())
	if err == nil {
		c.uuid = id
	}
	return err
}

// HasTransform returns true if the transform is different than the identity.
func (c *Component) HasTransform() bool {
	return !c.Transform.ApproxEqual(identityTransform)
}

// MergeToMesh merges a mesh with the component.
func (c *Component) MergeToMesh(m *mesh.Mesh, transform mgl32.Mat4) {
	c.Object.MergeToMesh(m, c.Transform.Mul4(transform))
}
