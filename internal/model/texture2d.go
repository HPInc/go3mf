package model

// Texture2DResource Resource defines the Model Texture 2D.
type Texture2DResource struct {
	ID          uint64
	Path        string
	ContentType Texture2DType
	TileStyleU  TileStyle
	TileStyleV  TileStyle
	Filter      TextureFilter
	modelPath   string
	uniqueID    uint64
}

// NewTexture2DResource returns a new texture 2D resource.
func NewTexture2DResource(id uint64) *Texture2DResource {
	return &Texture2DResource{
		ContentType: PNGTexture,
		TileStyleU:  TileWrap,
		TileStyleV:  TileWrap,
		Filter:      TextureFilterAuto,
	}
}

// ResourceID returns the resource ID, which has the same value as ID.
func (t *Texture2DResource) ResourceID() uint64 {
	return t.ID
}

// UniqueID returns the unique ID.
func (t *Texture2DResource) UniqueID() uint64 {
	return t.uniqueID
}

func (t *Texture2DResource) setUniqueID(id uint64) {
	t.uniqueID = id
}

// Copy copies the properties from another texture.
func (t *Texture2DResource) Copy(other *Texture2DResource) {
	t.Path = other.Path
	t.ContentType = other.ContentType
	t.TileStyleU = other.TileStyleU
	t.TileStyleV = other.TileStyleV
}
