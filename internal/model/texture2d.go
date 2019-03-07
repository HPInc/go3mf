package model

// Texture2DResource Resource defines the Model Texture 2D.
type Texture2DResource struct {
	Resource
	Path        string
	ContentType Texture2DType
	TileStyleU  TileStyle
	TileStyleV  TileStyle
	Filter      TextureFilter
}

// NewTexture2DResource returns a new texture 2D resource.
func NewTexture2DResource(id uint64, model *Model) (*Texture2DResource, error) {
	r, err := newResource(id, model)
	if err != nil {
		return nil, err
	}
	return &Texture2DResource{
		Resource:    *r,
		ContentType: PNGTexture,
		TileStyleU:  TileWrap,
		TileStyleV:  TileWrap,
		Filter:      TextureFilterAuto,
	}, nil
}

// Copy copies the properties from another texture.
func (t *Texture2DResource) Copy(other *Texture2DResource) {
	t.Path = other.Path
	t.ContentType = other.ContentType
	t.TileStyleU = other.TileStyleU
	t.TileStyleV = other.TileStyleV
}
