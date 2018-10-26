package meshinfo

import (
	"reflect"

	"github.com/qmuntal/go3mf/internal/common"
)

// MemoryMeshInfoFactory creates mesh info types with containers in memory.
type MemoryMeshInfoFactory struct {
}

func NewMemoryMeshInfoFactory() *MemoryMeshInfoFactory {
	return &MemoryMeshInfoFactory{}
}

// Create creates a new MeshInfo of the desired type.
func (f *MemoryMeshInfoFactory) Create(infoType InformationType, currentFaceCount uint32) (MeshInfo, error) {
	switch infoType {
	case InfoBaseMaterials:
		return newbaseMaterialsMeshInfo(newmemoryContainer(currentFaceCount, reflect.TypeOf((*BaseMaterial)(nil)).Elem())), nil
	case InfoNodeColors:
		return newnodeColorsMeshInfo(newmemoryContainer(currentFaceCount, reflect.TypeOf((*NodeColor)(nil)).Elem())), nil
	case InfoTextureCoords:
		return newtextureCoordsMeshInfo(newmemoryContainer(currentFaceCount, reflect.TypeOf((*TextureCoords)(nil)).Elem())), nil
	}
	return nil, common.NewError(common.ErrorInvalidInformationType)
}
