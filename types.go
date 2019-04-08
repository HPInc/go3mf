package go3mf

const thumbnailPath = "/Metadata/thumbnail.png"

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

func (t TextureFilter) String() string {
	return map[TextureFilter]string{
		TextureFilterAuto:    "auto",
		TextureFilterLinear:  "linear",
		TextureFilterNearest: "nearest",
	}[t]
}
