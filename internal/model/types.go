package model

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/gofrs/uuid"
)

const (
	relTypeTexture3D = "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture"
	relTypeThumbnail = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"
	relTypeModel3D   = "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel"
	langUS           = "en-US"
	thumbnailPath    = "/Metadata/thumbnail.png"
)

// WarningLevel defines the level of a reader warning.
type WarningLevel int

const (
	// FatalWarning fatal.
	FatalWarning WarningLevel = iota
	// MandatoryWarning mandatory.
	MandatoryWarning
	// OptionalWarning optional.
	OptionalWarning
)

// Units define the allowed model units.
type Units string

const (
	// Micrometer for microns
	Micrometer Units = "micron"
	// Millimeter for millimeter
	Millimeter = "millimeter"
	// Centimeter for centimeter
	Centimeter = "centimeter"
	// Inch for inch
	Inch = "inch"
	// Foot for foot
	Foot = "foot"
	// Meter for meter
	Meter = "meter"
)

// ClipMode defines the clipping modes for the beam lattices.
type ClipMode string

const (
	// ClipNone defines a beam lattice without clipping.
	ClipNone ClipMode = "none"
	// ClipInside defines a beam lattice with clipping inside.
	ClipInside = "inside"
	// ClipOutside defines a beam lattice with clipping outside.
	ClipOutside = "outside"
)

// SliceResolution defines the resolutions for a slice.
type SliceResolution string

const (
	// ResolutionFull defines a full resolution slice.
	ResolutionFull SliceResolution = "fullres"
	// ResolutionLow defines a low resolution slice.
	ResolutionLow = "lowres"
)

// ObjectType defines the allowed object types.
type ObjectType string

const (
	// OtherType defines a generic object type.
	OtherType ObjectType = "other"
	// ModelType defines a model object type.
	ModelType = "model"
	// SupportType defines a support object type.
	SupportType = "support"
	// SolidSupportType defines a solid support object type.
	SolidSupportType = "solidsupport"
	// SurfaceType defines a surface object type.
	SurfaceType = "surface"
)

// Texture2DType defines the allowed texture 2D types.
type Texture2DType string

const (
	// PNGTexture defines a png texture type.
	PNGTexture Texture2DType = "image/png"
	// JPEGTexture defines a jpeg texture type.
	JPEGTexture = "image/jpeg"
	// UnknownTexture defines an unknown texture type.
	UnknownTexture = ""
)

// TileStyle defines the allowed tile styles.
type TileStyle string

const (
	// WrapTile wraps the tile.
	WrapTile TileStyle = "wrap"
	// MirrorTile mirrors the tile.
	MirrorTile = "mirror"
	// ClampTile clamps the tile.
	ClampTile = "clamp"
	// NoneTile apply no style.
	NoneTile = "none"
)

// TextureFilter defines the allowed texture filters.
type TextureFilter string

const (
	// AutoFilter applies an automatic filter.
	AutoFilter TextureFilter = "auto"
	// LinearFilter applies a linear filter.
	LinearFilter = "linear"
	// NearestFilter applies an nearest filter.
	NearestFilter = "nearest"
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
	ClippingMeshID          *ResourceID
	RepresentationMeshID    *ResourceID
}

func registerUUID(old, new uuid.UUID, model *Model) error {
	err := model.registerUUID(new)
	if err == nil {
		model.unregisterUUID(old)
	}
	return err
}
