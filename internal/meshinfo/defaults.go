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

// NewNodeColorInfo creates a default node color mesh info.
func NewNodeColorInfo(currentFaceCount uint32) *GenericMeshInfo {
	return newInfo(currentFaceCount, NodeColorType)
}

// NewTextureCoordsInfo creates a default texture coordinates mesh info.
func NewTextureCoordsInfo(currentFaceCount uint32) *GenericMeshInfo {
	return newInfo(currentFaceCount, TextureCoordsType)
}

// NewBaseMaterialInfo creates a default base material mesh info.
func NewBaseMaterialInfo(currentFaceCount uint32) *GenericMeshInfo {
	return newInfo(currentFaceCount, BaseMaterialType)
}

func newInfo(currentFaceCount uint32, infoType reflect.Type) *GenericMeshInfo {
	return NewGenericMeshInfo(newmemoryContainer(currentFaceCount, infoType))
}
