package meshinfo

import (
	"github.com/go-gl/mathgl/mgl32"
)

const maxInternalID = 9223372036854775808

// Color represents a RGB color.
type Color uint32

// FaceData defines a generic information of a face. Implementations could by NodeColor or TextureCoords.
type FaceData interface{}

// InformationType is an enumerator that identifies different information types.
type InformationType int

const (
	// InfoAbstract defines abstract information.
	InfoAbstract InformationType = iota
	// InfoBaseMaterials defines base materials information.
	InfoBaseMaterials
	// InfoNodeColors defines node colors information.
	InfoNodeColors
	// InfoTextureCoords defines texture coordinates information.
	InfoTextureCoords
	// InfoCompositeMaterials defines composite materials information.
	InfoCompositeMaterials
	// InfoMultiProperties defines multiple properties information.
	InfoMultiProperties
	infoLastType
)

// NodeColor informs about the color of a node.
type NodeColor struct {
	Colors [3]Color // Colors of every vertex in a node.
}

// TextureCoords informs about the coordinates of a texture.
type TextureCoords struct {
	TextureID uint32        // Identifier of the texture.
	Coords    [3]mgl32.Vec2 // Coordinates of the boundaries of the texture.
}

// BaseMaterial informs about a base material.
type BaseMaterial struct {
	MaterialGroupID uint32 // Identifier of the group.
	MaterialIndex   uint32 // Index of the base material used in the group.
}

// MultiProperties informs about different properties.
type MultiProperties struct {
	MultiPropertyID uint32 // Encoded properties
}

// Composites informs about the properties of a composite.
type Composites struct {
	CompositeID uint32 // Identifier of the composite.
}

// meshInformationContainer provides a container for holding the texture information state of a complete mesh structure.
type meshInformationContainer interface {
	// addFaceData returns the pointer to the data of the added face.
	// The parameter newFaceCount should indicate the faces information stored in the container, including the new one.
	// If the count is not equal to the one returned by GetCurrentFaceCount an error will be returned.
	addFaceData(newFaceCount uint32) (val FaceData, err error)
	// getFaceData returns the data of the face with the target index.
	getFaceData(index uint32) (val FaceData, err error)
	// getCurrentFaceCount returns the number of faces information stored in the container.
	getCurrentFaceCount() uint32
	// clear removes all the information stored in the container.
	clear()
}

// MeshInformation defines the Mesh Information Class.
// This is a base class for handling all the mesh-related linear information (like face colors, textures, etc...).
type MeshInformation interface {
	GetFaceData(faceIndex uint32) (val FaceData, err error)
	AddFaceData(newFaceCount uint32) (val FaceData, err error)
	ResetFaceInformation(faceIndex uint32)
	ResetAllFaceInformation()
	GetType() InformationType
	FaceHasData(faceIndex uint32) bool
	cloneFaceInfosFrom(faceIndex uint32, otherInfo FaceData, otherFaceIndex uint32)
	invalidateFace(data FaceData)
	cloneInstance(currentFaceCount uint32) *MeshInformation
	permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32)
	mergeInformationFrom(info *MeshInformation)
	setInternalID(internalID uint64)
	getInternalID() uint64
}