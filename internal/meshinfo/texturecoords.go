package meshinfo

import (
	"github.com/go-gl/mathgl/mgl32"
)

// TextureCoords informs about the coordinates of a texture.
type TextureCoords struct {
	TextureID uint32        // Identifier of the texture.
	Coords    [3]mgl32.Vec2 // Coordinates of the boundaries of the texture.
}

// NewTextureCoords creates a new node color form an RGB color.
func NewTextureCoords(textureID uint32) *TextureCoords {
	return &TextureCoords{textureID, [3]mgl32.Vec2{mgl32.Vec2{0.0,0.0}, mgl32.Vec2{0.0,0.0}, mgl32.Vec2{0.0,0.0}}}
}

type textureCoordsInvalidator struct {
}

func (p textureCoordsInvalidator) Invalidate(data FaceData) {
	if node, ok := data.(*TextureCoords); ok {
		node.TextureID = 0
		node.Coords[0] = mgl32.Vec2{0.0,0.0}
		node.Coords[1] = mgl32.Vec2{0.0,0.0}
		node.Coords[2] = mgl32.Vec2{0.0,0.0}
	}
}

// TextureCoordsMeshInfo specializes the baseMeshInfo struct to "textures".
// It implements functions to interpolate and reconstruct texture coordinates while the mesh topology is changing.
type TextureCoordsMeshInfo struct {
	baseMeshInfo
}

// NewTextureCoordsMeshInfo creates a new Node colors mesh information struct.
func NewTextureCoordsMeshInfo(container Container) *TextureCoordsMeshInfo {
	container.Clear()
	return &TextureCoordsMeshInfo{*newBaseMeshInfo(container, textureCoordsInvalidator{})}
}

// GetType returns the type of information stored in this instance.
func (p *TextureCoordsMeshInfo) GetType() InformationType {
	return InfoTextureCoords
}

// FaceHasData checks if the specific face has any associated data.
func (p *TextureCoordsMeshInfo) FaceHasData(faceIndex uint32) bool {
	data, err := p.GetFaceData(faceIndex)
	if err == nil {
		return data.(*TextureCoords).TextureID != 0
	}
	return false
}

// Clone creates a deep copy of this instance.
func (p *TextureCoordsMeshInfo) Clone() MeshInfo {
	return NewTextureCoordsMeshInfo(p.baseMeshInfo.Container.Clone())
}

// cloneFaceInfosFrom clones the data from another face.
func (p *TextureCoordsMeshInfo) cloneFaceInfosFrom(faceIndex uint32, otherInfo MeshInfo, otherFaceIndex uint32) {
	targetData, err := p.GetFaceData(faceIndex)
	if err != nil {
		return
	}
	sourceData, err := otherInfo.GetFaceData(otherFaceIndex)
	if err != nil {
		return
	}
	node1, node2 := targetData.(*TextureCoords), sourceData.(*TextureCoords)
	node1.TextureID = node2.TextureID
	node1.Coords[0], node1.Coords[1], node1.Coords[2] = node2.Coords[0], node2.Coords[1],node2.Coords[2]
}

//permuteNodeInformation swaps the coordinates.
func (p *TextureCoordsMeshInfo) permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	data, err := p.GetFaceData(faceIndex)
	if err == nil && (nodeIndex1 < 3) && (nodeIndex2 < 3) && (nodeIndex3 < 3) {
		node := data.(*TextureCoords)
		node.Coords[0], node.Coords[1], node.Coords[2] = node.Coords[nodeIndex1], node.Coords[nodeIndex2], node.Coords[nodeIndex3]
	}
}

// mergeInformationFrom does nothing.
func (p *TextureCoordsMeshInfo) mergeInformationFrom(info MeshInfo) {
	// nothing to merge
}
