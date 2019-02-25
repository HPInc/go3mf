package model

// Texture2DResource Resource defines the Model Texture 2D.
type Texture2DResource struct {
	Resource
	Path                            string
	ContentType                     Texture2DType
	HasBox                          bool
	TileStyleU                      TileStyle
	TileStyleV                      TileStyle
	Filter                          TextureFilter
	boxU, boxV, boxWidth, boxHeight float32
}

// NewTexture2DResource returns a new texture 2D resource.
func NewTexture2DResource(id uint64, model *Model) (*Texture2DResource, error) {
	r, err := newResource(model, id)
	if err != nil {
		return nil, err
	}
	return &Texture2DResource{
		Resource:    *r,
		ContentType: UnknownTexture,
		boxWidth:    1,
		boxHeight:   1,
		TileStyleU:  WrapTile,
		TileStyleV:  WrapTile,
		Filter:      AutoFilter,
	}, nil
}

// Box returns the box of the texture or a default box if it doesn't have one.
func (t *Texture2DResource) Box() (u, v, width, height float32, hasBox bool) {
	if t.HasBox {
		u, v, width, height = t.boxU, t.boxV, t.boxWidth, t.boxHeight
		hasBox = true
	} else {
		width, height = 1, 1
		hasBox = false
	}
	return
}

// SetBox sets the box for the texture.
func (t *Texture2DResource) SetBox(u, v, width, height float32) *Texture2DResource {
	t.boxU, t.boxV, t.boxWidth, t.boxHeight = u, v, width, height
	t.HasBox = true
	return t
}

// ClearBox remove the box from the texture.
func (t *Texture2DResource) ClearBox() *Texture2DResource {
	t.boxU, t.boxV, t.boxWidth, t.boxHeight = 0, 0, 1, 1
	t.HasBox = false
	return t
}
