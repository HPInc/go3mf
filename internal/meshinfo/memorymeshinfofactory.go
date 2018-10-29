package meshinfo

import (
	"reflect"
)

// MemoryMeshInfoFactory creates mesh info types with containers in memory.
type MemoryMeshInfoFactory struct {
}

// Create creates a new MeshInfo of the desired type.
func (f *MemoryMeshInfoFactory) Create(infoType reflect.Type, currentFaceCount uint32) (MeshInfo, error) {
	return newgenericMeshInfo(newmemoryContainer(currentFaceCount, infoType)), nil
}
