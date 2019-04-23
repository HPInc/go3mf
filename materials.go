package go3mf

import "image/color"

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
type Composite []float64

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
