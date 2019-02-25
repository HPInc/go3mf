package model

const (
	relTypeTexture3D = "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture"
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
