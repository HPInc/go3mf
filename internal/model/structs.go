package model

import (
	"github.com/gofrs/uuid"
	"image/color"
)

// DefaultBaseMaterial defines the default base material property.
type DefaultBaseMaterial struct {
	ResourceID    uint64
	ResourceIndex uint64
}

// DefaultColor defines the default color property.
type DefaultColor struct {
	Color color.RGBA
}

// DefaultTexCoord2D defines the default textture coordinates property.
type DefaultTexCoord2D struct {
	ResourceID uint64
	U, V       float32
}

// BeamLatticeAttributes defines the Model Mesh BeamLattice Attributes class and is part of the BeamLattice extension to 3MF.
type BeamLatticeAttributes struct {
	ClipMode                ClipMode
	HasClippingMeshID       bool
	HasRepresentationMeshID bool
	ClippingMeshID          PackageResourceID
	RepresentationMeshID    PackageResourceID
}

// An Object is an in memory representation of the 3MF model object.
type Object struct {
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

// UUID returns the object UUID.
func (o *Object) UUID() uuid.UUID {
	return o.uuid
}

// SetUUID sets the object UUID
func (o *Object) SetUUID(id uuid.UUID) {
	o.uuid = id
}
