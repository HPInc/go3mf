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

// NewNodeColorFacesData creates a default node color mesh info.
func NewNodeColorFacesData(currentFaceCount uint32) *FacesData {
	return newInfo(currentFaceCount, NodeColorType)
}

// NewTextureCoordsFacesData creates a default texture coordinates mesh info.
func NewTextureCoordsFacesData(currentFaceCount uint32) *FacesData {
	return newInfo(currentFaceCount, TextureCoordsType)
}

// NewBaseMaterialFacesData creates a default base material mesh info.
func NewBaseMaterialFacesData(currentFaceCount uint32) *FacesData {
	return newInfo(currentFaceCount, BaseMaterialType)
}

func newInfo(currentFaceCount uint32, infoType reflect.Type) *FacesData {
	return newFacesData(newmemoryContainer(currentFaceCount, infoType))
}
