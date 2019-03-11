package meshinfo

import (
	"github.com/go-gl/mathgl/mgl32"
)

// TextureCoords informs about the coordinates of a texture.
type TextureCoords struct {
	TextureID uint32 // Identifier of the texture.
	// Each texture coordinate is, at a minimum, a (U,V) pair,
	// which is the horizontal and vertical location in texture space, respectively.
	Coords [3]mgl32.Vec2 // Coordinates of the boundaries of the texture.
}

// Invalidate sets to zero all the properties.
func (t *TextureCoords) Invalidate() {
	t.TextureID = 0
	t.Coords[0] = mgl32.Vec2{0.0, 0.0}
	t.Coords[1] = mgl32.Vec2{0.0, 0.0}
	t.Coords[2] = mgl32.Vec2{0.0, 0.0}
}

// Copy copy the properties of another texture coords.
func (t *TextureCoords) Copy(from FaceData) {
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
func (t *TextureCoords) Merge(other FaceData) {
	// nothing to merge
}

type textureCoordsContainer struct {
	dataBlocks []*TextureCoords
}

func newtextureCoordsContainer(currentFaceCount uint32) *textureCoordsContainer {
	m := &textureCoordsContainer{
		dataBlocks: make([]*TextureCoords, 0, int(currentFaceCount)),
	}
	for i := uint32(1); i <= currentFaceCount; i++ {
		m.AddFaceData(i)
	}
	return m
}

func (m *textureCoordsContainer) clone(currentFaceCount uint32) Container {
	return newtextureCoordsContainer(currentFaceCount)
}

func (m *textureCoordsContainer) InfoType() DataType {
	return TextureCoordsType
}

func (m *textureCoordsContainer) AddFaceData(newFaceCount uint32) FaceData {
	faceData := new(TextureCoords)
	m.dataBlocks = append(m.dataBlocks, faceData)
	if len(m.dataBlocks) != int(newFaceCount) {
		panic(errFaceCountMissmatch)
	}
	return faceData
}

func (m *textureCoordsContainer) FaceData(faceIndex uint32) FaceData {
	return m.dataBlocks[int(faceIndex)]
}

func (m *textureCoordsContainer) FaceCount() uint32 {
	return uint32(len(m.dataBlocks))
}

func (m *textureCoordsContainer) Clear() {
	m.dataBlocks = m.dataBlocks[:0]
}
