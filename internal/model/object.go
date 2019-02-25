package model

import (
	"github.com/satori/go.uuid"
)

// An Object is an in memory representation of the 3MF model object.
type Object struct {
	Resource
	Name            string
	PartNumber      string
	SliceStackID    *PackageResourceID
	SliceResoultion SliceResolution
	Thumbnail string
	DefaultProperty interface{}
	Type ObjectType
	uuid uuid.UUID
}

// UUID returns the object UUID.
func (o *Object) UUID() uuid.UUID {
	return o.uuid
}

// SetUUID sets the object UUID
func (o *Object) SetUUID(id uuid.UUID) {
	o.uuid = id
}