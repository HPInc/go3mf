package meshinfo

import (
	"reflect"

	"github.com/qmuntal/go3mf/internal/common"
)

// MemoryMeshInfoFactory creates mesh info types with containers in memory.
type MemoryMeshInfoFactory struct {
}

// Create creates a new MeshInfo of the desired type.
func (f *MemoryMeshInfoFactory) Create(infoType InformationType, currentFaceCount uint32) (MeshInfo, error) {
	switch infoType {
	case InfoBaseMaterials:
		return newgenericMeshInfo(newmemoryContainer(currentFaceCount, reflect.TypeOf((*BaseMaterial)(nil)).Elem()), infoType), nil
	case InfoNodeColors:
		return newgenericMeshInfo(newmemoryContainer(currentFaceCount, reflect.TypeOf((*NodeColor)(nil)).Elem()), infoType), nil
	case InfoTextureCoords:
		return newgenericMeshInfo(newmemoryContainer(currentFaceCount, reflect.TypeOf((*TextureCoords)(nil)).Elem()), infoType), nil
	}
	return nil, common.NewError(common.ErrorInvalidInformationType)
}
