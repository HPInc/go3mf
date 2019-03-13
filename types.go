package go3mf

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
		"micron":     UnitMicrometer,
		"centimeter": UnitCentimeter,
		"inch":       UnitInch,
		"foot":       UnitFoot,
		"meter":      UnitMeter,
	}[s]
	return
}

func (u Units) String() string {
	return map[Units]string{
		UnitMillimeter: "millimeter",
		UnitMicrometer: "micron",
		UnitCentimeter: "centimeter",
		UnitInch:       "inch",
		UnitFoot:       "foot",
		UnitMeter:      "meter",
	}[u]
}

// ClipMode defines the clipping modes for the beam lattices.
type ClipMode uint8

const (
	// ClipNone defines a beam lattice without clipping.
	ClipNone ClipMode = iota
	// ClipInside defines a beam lattice with clipping inside.
	ClipInside
	// ClipOutside defines a beam lattice with clipping outside.
	ClipOutside
)

// NewClipMode returns a new ClipMode from a string.
func NewClipMode(s string) (c ClipMode, ok bool) {
	c, ok = map[string]ClipMode{
		"none":    ClipNone,
		"inside":  ClipInside,
		"outside": ClipOutside,
	}[s]
	return
}

func (c ClipMode) String() string {
	return map[ClipMode]string{
		ClipNone:    "none",
		ClipInside:  "inside",
		ClipOutside: "outside",
	}[c]
}

// SliceResolution defines the resolutions for a slice.
type SliceResolution uint8

const (
	// ResolutionFull defines a full resolution slice.
	ResolutionFull SliceResolution = iota
	// ResolutionLow defines a low resolution slice.
	ResolutionLow
)

// NewSliceResolution returns a new SliceResolution from a string.
func NewSliceResolution(s string) (r SliceResolution, ok bool) {
	r, ok = map[string]SliceResolution{
		"fullres": ResolutionFull,
		"lowres":  ResolutionLow,
	}[s]
	return
}

func (c SliceResolution) String() string {
	return map[SliceResolution]string{
		ResolutionFull: "fullres",
		ResolutionLow:  "lowres",
	}[c]
}

// ObjectType defines the allowed object types.
type ObjectType int8

const (
	// ObjectTypeModel defines a model object type.
	ObjectTypeModel ObjectType = iota
	// ObjectTypeOther defines a generic object type.
	ObjectTypeOther
	// ObjectTypeSupport defines a support object type.
	ObjectTypeSupport
	// ObjectTypeSolidSupport defines a solid support object type.
	ObjectTypeSolidSupport
	// ObjectTypeSurface defines a surface object type.
	ObjectTypeSurface
)

// NewObjectType returns a new ObjectType from a string.
func NewObjectType(s string) (o ObjectType, ok bool) {
	o, ok = map[string]ObjectType{
		"model":        ObjectTypeModel,
		"other":        ObjectTypeOther,
		"support":      ObjectTypeSupport,
		"solidsupport": ObjectTypeSolidSupport,
		"surface":      ObjectTypeSurface,
	}[s]
	return
}

func (o ObjectType) String() string {
	return map[ObjectType]string{
		ObjectTypeModel:        "model",
		ObjectTypeOther:        "other",
		ObjectTypeSupport:      "support",
		ObjectTypeSolidSupport: "solidsupport",
		ObjectTypeSurface:      "surface",
	}[o]
}

// Texture2DType defines the allowed texture 2D types.
type Texture2DType uint8

const (
	// PNGTexture defines a png texture type.
	PNGTexture Texture2DType = iota + 1
	// JPEGTexture defines a jpeg texture type.
	JPEGTexture
)

// NewTexture2DType returns a new Texture2DType from a string.
func NewTexture2DType(s string) (t Texture2DType, ok bool) {
	t, ok = map[string]Texture2DType{
		"image/png":  PNGTexture,
		"image/jpeg": JPEGTexture,
	}[s]
	return
}

func (t Texture2DType) String() string {
	return map[Texture2DType]string{
		PNGTexture:  "image/png",
		JPEGTexture: "image/jpeg",
	}[t]
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
		"wrap":   TileWrap,
		"mirror": TileMirror,
		"clamp":  TileClamp,
		"none":   TileNone,
	}[s]
	return
}

func (t TileStyle) String() string {
	return map[TileStyle]string{
		TileWrap:   "wrap",
		TileMirror: "mirror",
		TileClamp:  "clamp",
		TileNone:   "none",
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
		"auto":    TextureFilterAuto,
		"linear":  TextureFilterLinear,
		"nearest": TextureFilterNearest,
	}[s]
	return
}

func (t TextureFilter) String() string {
	return map[TextureFilter]string{
		TextureFilterAuto:    "auto",
		TextureFilterLinear:  "linear",
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
