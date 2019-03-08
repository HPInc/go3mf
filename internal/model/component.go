package model

import (
	"errors"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gofrs/uuid"
	"github.com/qmuntal/go3mf/internal/mesh"
)

// Object defines a composable object.
type Object interface {
	MergeToMesh(*mesh.Mesh, mgl32.Mat4) error
	IsValid() bool
	IsValidForSlices(mgl32.Mat4) bool
	Type() ObjectType
}

// An ObjectResource is an in memory representation of the 3MF model object.
type ObjectResource struct {
	ID              uint64
	Name            string
	PartNumber      string
	SliceStackID    *ResourceID
	SliceResoultion SliceResolution
	Thumbnail       string
	DefaultProperty interface{}
	ObjectType      ObjectType
	uuid            uuid.UUID
	uuidRegister    register
	modelPath       string
	uniqueID        uint64
}

// ResourceID returns the resource ID, which has the same value as ID.
func (o *ObjectResource) ResourceID() uint64 {
	return o.ID
}

// UniqueID returns the unique ID.
func (o *ObjectResource) UniqueID() uint64 {
	return o.uniqueID
}

func (o *ObjectResource) setUniqueID(id uint64) {
	o.uniqueID = id
}

// Type returns the type of the object.
func (o *ObjectResource) Type() ObjectType {
	return o.ObjectType
}

// UUID returns the object UUID.
func (o *ObjectResource) UUID() uuid.UUID {
	return o.uuid
}

// SetUUID sets the object UUID
func (o *ObjectResource) SetUUID(id uuid.UUID) error {
	if o.uuidRegister == nil {
		return errors.New("go3mf: object resource uuid cannot be set as it is not inside any model")
	}
	err := o.uuidRegister.register(o.uuid, id)
	if err == nil {
		o.uuid = id
	}
	return err
}

// MergeToMesh left on purpose empty to be redefined in embedding class.
func (o *ObjectResource) MergeToMesh(m *mesh.Mesh, transform mgl32.Mat4) error {
	return nil
}

// IsValid should be redefined in embedding class.
func (o *ObjectResource) IsValid() bool {
	return false
}

// IsValidForSlices should be redefined in embedding class.
func (o *ObjectResource) IsValidForSlices(transform mgl32.Mat4) bool {
	return false
}

// A Component is an in memory representation of the 3MF component.
type Component struct {
	Object       Object
	Transform    mgl32.Mat4
	uuid         uuid.UUID
	uuidRegister register
}

// UUID returns the object UUID.
func (c *Component) UUID() uuid.UUID {
	return c.uuid
}

// SetUUID sets the object UUID
func (c *Component) SetUUID(id uuid.UUID) error {
	if c.uuidRegister == nil {
		return errors.New("go3mf: component uuid cannot be set as it is not inside any model")
	}
	err := c.uuidRegister.register(c.uuid, id)
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
func (c *Component) MergeToMesh(m *mesh.Mesh, transform mgl32.Mat4) error {
	return c.Object.MergeToMesh(m, c.Transform.Mul4(transform))
}

// A ComponentResource resource is an in memory representation of the 3MF component object.
type ComponentResource struct {
	ObjectResource
	Components []*Component
}

// MergeToMesh merges the mesh with all the components.
func (c *ComponentResource) MergeToMesh(m *mesh.Mesh, transform mgl32.Mat4) error {
	for _, comp := range c.Components {
		if err := comp.MergeToMesh(m, transform); err != nil {
			return err
		}
	}
	return nil
}

// IsValid checks if the component resource and all its child are valid.
func (c *ComponentResource) IsValid() bool {
	if len(c.Components) == 0 {
		return false
	}

	for _, comp := range c.Components {
		if !comp.Object.IsValid() {
			return false
		}
	}
	return true
}

// IsValidForSlices checks if the component resource and all its child are valid to be used with slices.
func (c *ComponentResource) IsValidForSlices(transform mgl32.Mat4) bool {
	if len(c.Components) == 0 {
		return true
	}

	for _, comp := range c.Components {
		if !comp.Object.IsValidForSlices(transform.Mul4(comp.Transform)) {
			return false
		}
	}
	return true
}

// A MeshResource is an in memory representation of the 3MF mesh object.
type MeshResource struct {
	ObjectResource
	Mesh                  *mesh.Mesh
	BeamLatticeAttributes BeamLatticeAttributes
}

// MergeToMesh merges the resource with the mesh.
func (c *MeshResource) MergeToMesh(m *mesh.Mesh, transform mgl32.Mat4) error {
	return c.Mesh.Merge(m, transform)
}

// IsValid checks if the mesh resource are valid.
func (c *MeshResource) IsValid() bool {
	if c.Mesh == nil {
		return false
	}
	switch c.ObjectType {
	case ObjectTypeModel:
		return c.Mesh.IsManifoldAndOriented()
	case ObjectTypeSupport:
		return c.Mesh.BeamCount() == 0
	case ObjectTypeSolidSupport:
		return c.Mesh.IsManifoldAndOriented()
	case ObjectTypeSurface:
		return c.Mesh.BeamCount() == 0
	}

	return false
}

// IsValidForSlices checks if the mesh resource are valid for slices.
func (c *MeshResource) IsValidForSlices(t mgl32.Mat4) bool {
	return c.SliceStackID == nil || t[2] == 0 && t[6] == 0 && t[8] == 0 && t[9] == 0 && t[10] == 1
}
