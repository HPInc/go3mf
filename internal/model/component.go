package model

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gofrs/uuid"
	"github.com/qmuntal/go3mf/internal/mesh"
)

// Object defines a composable object.
type Object interface {
	RootModel() *Model
	MergeToMesh(*mesh.Mesh, mgl32.Mat4) error
	ID() uint64
	IsValid() bool
	IsValidForSlices(mgl32.Mat4) bool
}

// An ObjectResource is an in memory representation of the 3MF model object.
type ObjectResource struct {
	Resource
	Name            string
	PartNumber      string
	SliceStackID    *ResourceID
	SliceResoultion SliceResolution
	Thumbnail       string
	DefaultProperty interface{}
	Type            ObjectType
	uuid            uuid.UUID
}

func newObjectResource(id uint64, model *Model) (*ObjectResource, error) {
	r, err := newResource(id, model)
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

// RootModel returns the model of the object.
func (o *ObjectResource) RootModel() *Model {
	return o.Model
}

// ID returns the id of the object.
func (o *ObjectResource) ID() uint64 {
	return o.ResourceID.UniqueID()
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
func (c *Component) MergeToMesh(m *mesh.Mesh, transform mgl32.Mat4) error {
	return c.Object.MergeToMesh(m, c.Transform.Mul4(transform))
}

// A ComponentResource resource is an in memory representation of the 3MF component object.
type ComponentResource struct {
	ObjectResource
	Components []*Component
}

// NewComponentResource returns a new component resource.
func NewComponentResource(id uint64, model *Model) (*ComponentResource, error) {
	r, err := newObjectResource(id, model)
	if err != nil {
		return nil, err
	}
	return &ComponentResource{
		ObjectResource: *r,
	}, nil
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

// NewMeshResource returns a new mesh resource.
func NewMeshResource(id uint64, model *Model) (*MeshResource, error) {
	r, err := newObjectResource(id, model)
	if err != nil {
		return nil, err
	}
	return &MeshResource{
		ObjectResource: *r,
	}, nil
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
	switch c.Type {
	case ModelType:
		return c.Mesh.IsManifoldAndOriented()
	case SupportType:
		return c.Mesh.BeamCount() == 0
	case SolidSupportType:
		return c.Mesh.IsManifoldAndOriented()
	case SurfaceType:
		return c.Mesh.BeamCount() == 0
	}

	return false
}

// IsValidForSlices checks if the mesh resource are valid for slices.
func (c *MeshResource) IsValidForSlices(t mgl32.Mat4) bool {
	return c.SliceStackID == nil || t[2] == 0 && t[6] == 0 && t[8] == 0 && t[9] == 0 && t[10] == 1
}
