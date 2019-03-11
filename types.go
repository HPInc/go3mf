package model

import (
	"image/color"
	"io"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	thumbnailPath = "/Metadata/thumbnail.png"
	langUS        = "en-US"
)

// Units define the allowed model units.
type Units uint8

const (
	// UnitMillimeter for millimeter
	UnitMillimeter Units = iota
	// UnitMicrometer for microns
	UnitMicrometer
	// UnitCentimeter for centimeter
	UnitCentimeter
	// UnitInch for inch
	UnitInch
	// UnitFoot for foot
	UnitFoot
	// UnitMeter for meter
	UnitMeter
)

// NewUnits returns a new unit from a string.
func NewUnits(s string) (u Units, ok bool) {
	u, ok = map[string]Units{
		"millimeter": UnitMillimeter,
		"micron": UnitMicrometer,
		"centimeter": UnitCentimeter,
		"inch": UnitInch,
		"foot": UnitFoot,
		"meter": UnitMeter,
	}[s]
	return
}

func (u Units) String() string {
	return map[Units]string{
		UnitMillimeter: "millimeter",
		UnitMicrometer: "micron",
		UnitCentimeter: "centimeter",
		UnitInch: "inch",
		UnitFoot: "foot",
		UnitMeter: "meter",
	}[u]
}

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
	// ObjectTypeOther defines a generic object type.
	ObjectTypeOther ObjectType = "other"
	// ObjectTypeModel defines a model object type.
	ObjectTypeModel = "model"
	// ObjectTypeSupport defines a support object type.
	ObjectTypeSupport = "support"
	// ObjectTypeSolidSupport defines a solid support object type.
	ObjectTypeSolidSupport = "solidsupport"
	// ObjectTypeSurface defines a surface object type.
	ObjectTypeSurface = "surface"
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

// NewTexture2DType returns a new Texture2DType from a string.
func NewTexture2DType(s string) (Texture2DType, bool) {
	u := Texture2DType(s)
	if u == PNGTexture || u == JPEGTexture {
		return u, true
	}
	return "", false
}

// TileStyle defines the allowed tile styles.
type TileStyle uint8

const (
	// TileWrap wraps the tile.
	TileWrap TileStyle = iota
	// TileMirror mirrors the tile.
	TileMirror
	// TileClamp clamps the tile.
	TileClamp
	// TileNone apply no style.
	TileNone
)

// NewTileStyle returns a new TileStyle from a string.
func NewTileStyle(s string) (t TileStyle, ok bool) {
	t, ok = map[string]TileStyle{
		"wrap": TileWrap,
		"mirror": TileMirror,
		"clamp": TileClamp,
		"none": TileNone,
	}[s]
	return
}

func (t TileStyle) String() string {
	return map[TileStyle]string{
		TileWrap: "wrap",
		TileMirror: "mirror",
		TileClamp: "clamp",
		TileNone: "none",
	}[t]
}

// TextureFilter defines the allowed texture filters.
type TextureFilter uint8

const (
	// TextureFilterAuto applies an automatic filter.
	TextureFilterAuto TextureFilter = iota
	// TextureFilterLinear applies a linear filter.
	TextureFilterLinear
	// TextureFilterNearest applies an nearest filter.
	TextureFilterNearest
)

// NewTextureFilter returns a new TextureFilter from a string.
func NewTextureFilter(s string) (t TextureFilter, ok bool) {
	t, ok = map[string]TextureFilter{
		"auto": TextureFilterAuto,
		"linear": TextureFilterLinear,
		"nearest": TextureFilterNearest,
	}[s]
	return
}

func (t TextureFilter) String() string {
	return map[TextureFilter]string{
		TextureFilterAuto: "auto",
		TextureFilterLinear: "linear",
		TextureFilterNearest: "nearest",
	}[t]
}

var identityTransform = mgl32.Ident4()

// Metadata item is an in memory representation of the 3MF metadata,
// and can be attached to any 3MF model node.
type Metadata struct {
	Name  string
	Value string
}

// Attachment defines the Model Attachment.
type Attachment struct {
	Stream           io.Reader
	RelationshipType string
	Path             string
}

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
	ClippingMeshID          uint64
	RepresentationMeshID    uint64
}
