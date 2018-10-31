package meshinfo

import "reflect"

var (
	// NodeColorType is the type of a NodeColor.
	NodeColorType = reflect.TypeOf((*NodeColor)(nil))
	// TextureCoordsType is the type of a TextureCoords.
	TextureCoordsType = reflect.TypeOf((*TextureCoords)(nil))
	// BaseMaterialType is the type of a BaseMaterial.
	BaseMaterialType = reflect.TypeOf((*BaseMaterial)(nil))
)

// NewHandler creates a default handler.
func NewHandler() Handler {
	return newlookupHandler()
}

// NewNodeColorInfo creates a default node color mesh info.
func NewNodeColorInfo(currentFaceCount uint32) MeshInfo {
	return newInfo(currentFaceCount, NodeColorType)
}

// NewTextureCoordsInfo creates a default texture coordinates mesh info.
func NewTextureCoordsInfo(currentFaceCount uint32) MeshInfo {
	return newInfo(currentFaceCount, TextureCoordsType)
}

// NewBaseMaterialInfo creates a default base material mesh info.
func NewBaseMaterialInfo(currentFaceCount uint32) MeshInfo {
	return newInfo(currentFaceCount, BaseMaterialType)
}

func newInfo(currentFaceCount uint32, infoType reflect.Type) MeshInfo {
	return newgenericMeshInfo(newmemoryContainer(currentFaceCount, infoType))
}
