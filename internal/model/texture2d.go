package model

// Texture2DResource Resource defines the Model Texture 2D.
type Texture2DResource struct {
	ID          uint64
	ModelPath   string
	Path        string
	ContentType Texture2DType
	TileStyleU  TileStyle
	TileStyleV  TileStyle
	Filter      TextureFilter
}

// NewTexture2DResource returns a new texture 2D resource.
func NewTexture2DResource(id uint64) *Texture2DResource {
	return &Texture2DResource{
		ID:          id,
		ContentType: PNGTexture,
		TileStyleU:  TileWrap,
		TileStyleV:  TileWrap,
		Filter:      TextureFilterAuto,
	}
}

// Identify returns the resource ID and the ModelPath.
func (t *Texture2DResource) Identify() (uint64, string) {
	return t.ID, t.ModelPath
}

// Copy copies the properties from another texture.
func (t *Texture2DResource) Copy(other *Texture2DResource) {
	t.Path = other.Path
	t.ContentType = other.ContentType
	t.TileStyleU = other.TileStyleU
	t.TileStyleV = other.TileStyleV
}
