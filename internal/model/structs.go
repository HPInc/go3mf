package model

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/gofrs/uuid"
)

var identityTransform = mgl32.Ident4()

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
	ClippingMeshID          *PackageResourceID
	RepresentationMeshID    *PackageResourceID
}

func registerUUID(old, new uuid.UUID, model *Model) error {
	err := model.registerUUID(new)
	if err == nil {
		model.unregisterUUID(old)
	}
	return err
}
