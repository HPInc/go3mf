package meshinfo

import "reflect"

// NewHandler creates a default handler.
func NewHandler() Handler {
	return newlookupHandler()
}

// NewNodeColorInfo creates a default node color mesh info.
func NewNodeColorInfo(currentFaceCount uint32) MeshInfo {
	return newInfo(currentFaceCount, reflect.TypeOf((*NodeColor)(nil)).Elem())
}

// NewTextureCoordsInfo creates a default texture coordinates mesh info.
func NewTextureCoordsInfo(currentFaceCount uint32) MeshInfo {
	return newInfo(currentFaceCount, reflect.TypeOf((*TextureCoords)(nil)).Elem())
}

// NewBaseMaterialInfo creates a default base material mesh info.
func NewBaseMaterialInfo(currentFaceCount uint32) MeshInfo {
	return newInfo(currentFaceCount, reflect.TypeOf((*BaseMaterial)(nil)).Elem())
}

func newInfo(currentFaceCount uint32, infoType reflect.Type) MeshInfo {
	return newgenericMeshInfo(newmemoryContainer(currentFaceCount, infoType))
}
