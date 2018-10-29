package meshinfo

import (
	"github.com/go-gl/mathgl/mgl32"
)

// TextureCoords informs about the coordinates of a texture.
type TextureCoords struct {
	TextureID uint32        // Identifier of the texture.
	Coords    [3]mgl32.Vec2 // Coordinates of the boundaries of the texture.
}

// NewTextureCoords creates a new NewTextureCoords.
func NewTextureCoords(textureID uint32, coord1, coord2, coord3 mgl32.Vec2) *TextureCoords {
	return &TextureCoords{textureID, [3]mgl32.Vec2{coord1, coord2, coord3}}
}

// Invalidate sets to zero all the properties.
func (t *TextureCoords) Invalidate() {
	t.TextureID = 0
	t.Coords[0] = mgl32.Vec2{0.0, 0.0}
	t.Coords[1] = mgl32.Vec2{0.0, 0.0}
	t.Coords[2] = mgl32.Vec2{0.0, 0.0}
}

// Copy copy the properties of another texture coords.
func (t *TextureCoords) Copy(from interface{}) {
	other, ok := from.(*TextureCoords)
	if !ok {
		return
	}
	t.TextureID = other.TextureID
	t.Coords[0], t.Coords[1], t.Coords[2] = other.Coords[0], other.Coords[1], other.Coords[2]
}

// HasData returns true if the texture id is different from zero.
func (t *TextureCoords) HasData() bool {
	return t.TextureID != 0
}

// Permute swap the coordinates using the given indexes. Do nothing if any of the indexes is bigger than 2.
func (t *TextureCoords) Permute(index1, index2, index3 uint32) {
	if (index1 > 2) || (index2 > 2) || (index3 > 2) {
		return
	}
	t.Coords[0], t.Coords[1], t.Coords[2] = t.Coords[index1], t.Coords[index2], t.Coords[index3]
}

// Merge is not necessary for a tetxure coords.
func (t *TextureCoords) Merge(other interface{}) {
	// nothing to merge
}
