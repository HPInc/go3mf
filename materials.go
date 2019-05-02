package go3mf

import "image/color"

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

// BlendMethod defines the equation to use when blending a layer with the previous layer.
type BlendMethod uint8

const (
	// BlendMix to mix the layer properties.
	BlendMix BlendMethod = iota
	// BlendMultiply multiply the layer properties.
	BlendMultiply
)

func (b BlendMethod) String() string {
	return map[BlendMethod]string{
		BlendMix:      "mix",
		BlendMultiply: "multiply",
	}[b]
}

// Texture2DResource defines the Model Texture 2D.
type Texture2DResource struct {
	ID          uint32
	ModelPath   string
	Path        string
	ContentType Texture2DType
	TileStyleU  TileStyle
	TileStyleV  TileStyle
	Filter      TextureFilter
}

// Identify returns the unique ID of the resource.
func (t *Texture2DResource) Identify() (string, uint32) {
	return t.ModelPath, t.ID
}

// Copy copies the properties from another texture.
func (t *Texture2DResource) Copy(other *Texture2DResource) {
	t.Path = other.Path
	t.ContentType = other.ContentType
	t.TileStyleU = other.TileStyleU
	t.TileStyleV = other.TileStyleV
}

// TextureCoord map a vertex of a triangle to a position in image space (U, V coordinates)
type TextureCoord [2]float32

// U returns the first coordinate.
func (t TextureCoord) U() float32 {
	return t[0]
}

// V returns the second coordinate.
func (t TextureCoord) V() float32 {
	return t[1]
}

// Texture2DGroupResource acts as a container for texture coordinate properties.
type Texture2DGroupResource struct {
	ID        uint32
	ModelPath string
	TextureID uint32
	Coords    []TextureCoord
}

// Identify returns the unique ID of the resource.
func (t *Texture2DGroupResource) Identify() (string, uint32) {
	return t.ModelPath, t.ID
}

// ColorGroupResource acts as a container for color properties.
type ColorGroupResource struct {
	ID        uint32
	ModelPath string
	Colors    []color.RGBA
}

// Identify returns the unique ID of the resource.
func (c *ColorGroupResource) Identify() (string, uint32) {
	return c.ModelPath, c.ID
}

// A Composite specifies the proportion of the overall mixture for each material.
type Composite struct {
	Values []float64
}

// CompositeMaterialsResource defines materials derived by mixing 2 or more base materials in defined ratios.
type CompositeMaterialsResource struct {
	ID         uint32
	ModelPath  string
	MaterialID uint32
	Indices    []uint32
	Composites []Composite
}

// Identify returns the unique ID of the resource.
func (c *CompositeMaterialsResource) Identify() (string, uint32) {
	return c.ModelPath, c.ID
}

// The Multi element combines the constituent materials and properties.
type Multi struct {
	ResourceIndices []uint32
}

// A MultiPropertiesResource element acts as a container for Multi
// elements which are indexable groups of property indices.
type MultiPropertiesResource struct {
	ID           uint32
	ModelPath    string
	Resources    []uint32
	BlendMethods []BlendMethod
	Multis       []Multi
}

// Identify returns the unique ID of the resource.
func (c *MultiPropertiesResource) Identify() (string, uint32) {
	return c.ModelPath, c.ID
}
