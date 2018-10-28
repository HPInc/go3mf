package meshinfo

import (
	"github.com/go-gl/mathgl/mgl32"
)

// TextureCoords informs about the coordinates of a texture.
type TextureCoords struct {
	TextureID uint32        // Identifier of the texture.
	Coords    [3]mgl32.Vec2 // Coordinates of the boundaries of the texture.
}

func (t *TextureCoords) Invalidate() {
	t.TextureID = 0
	t.Coords[0] = mgl32.Vec2{0.0, 0.0}
	t.Coords[1] = mgl32.Vec2{0.0, 0.0}
	t.Coords[2] = mgl32.Vec2{0.0, 0.0}
}

func (t *TextureCoords) Copy(from interface{}) {
	other, ok := from.(*TextureCoords)
	if !ok {
		return
	}
	t.Coords[0], t.Coords[1], t.Coords[2] = other.Coords[0], other.Coords[1], other.Coords[2]
}

func (t *TextureCoords) HasData() bool {
	return t.TextureID != 0
}

func (t *TextureCoords) Permute(index1, index2, index3 uint32) {
	if (index1 > 2) || (index2 > 2) || (index3 > 2) {
		return
	}
	t.Coords[0], t.Coords[1], t.Coords[2] = t.Coords[index1], t.Coords[index2], t.Coords[index3]
}

func (t *TextureCoords) Merge(other interface{}) {
	// nothing to merge
}
