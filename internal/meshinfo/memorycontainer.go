package meshinfo

import (
	"reflect"

	"github.com/qmuntal/go3mf/internal/common"
)

// memoryContainer implements Container
// and provides a generic memory container for holding mesh information state of a complete mesh structure
// using reflection to infer slice type.
type memoryContainer struct {
	infoType   reflect.Type
	faceCount  uint32
	dataBlocks reflect.Value
}

// newmemoryContainer creates a new container that holds the specified element types.
func newmemoryContainer(currentFaceCount uint32, infoType reflect.Type) *memoryContainer {
	m := &memoryContainer{
		faceCount:  0,
		infoType:   infoType,
		dataBlocks: reflect.MakeSlice(reflect.SliceOf(infoType), 0, int(currentFaceCount)),
	}
	for i := 1; i <= int(currentFaceCount); i++ {
		m.AddFaceData(uint32(i))
	}
	return m
}

// Clone creates a copy of the container with all the faces invalidated.
func (m *memoryContainer) Clone(currentFaceCount uint32) Container {
	return newmemoryContainer(currentFaceCount, m.infoType)
}

	// InfoType returns the type of the stored data.
func (m *memoryContainer) InfoType() reflect.Type {
	return m.infoType
}

// AddFaceData returns the pointer to the data of the added face.
// The parameter newFaceCount should indicate the faces information stored in the container, including the new one.
// Error cases:
// * ErrorInvalidRecordSize: The element type is not defined.
// * ErrorMeshInformationCountMismatch: The number of faces in the container does not match with the input parameter.
func (m *memoryContainer) AddFaceData(newFaceCount uint32) (FaceData, error) {
	if m.infoType == nil {
		return nil, common.NewError(common.ErrorInvalidRecordSize)
	}
	faceData := reflect.New(m.infoType)
	m.dataBlocks = reflect.Append(m.dataBlocks, faceData.Elem())
	m.faceCount++
	if m.faceCount != newFaceCount {
		return nil, common.NewError(common.ErrorMeshInformationCountMismatch)
	}
	return faceData.Interface().(FaceData), nil
}

// GetFaceData returns the data of the face with the target index.
// Error cases:
// * ErrorInvalidMeshInformationIndex: Index is higher than the number of faces
func (m *memoryContainer) GetFaceData(faceIndex uint32) (FaceData, error) {
	if faceIndex >= m.faceCount {
		return nil, common.NewError(common.ErrorInvalidMeshInformationIndex)
	}

	return m.dataBlocks.Index(int(faceIndex)).Addr().Interface().(FaceData), nil
}

// GetCurrentFaceCount returns the number of faces information stored in the container.
func (m *memoryContainer) GetCurrentFaceCount() uint32 {
	return m.faceCount
}

// Clear removes all the information stored in the container.
func (m *memoryContainer) Clear() {
	m.dataBlocks = reflect.MakeSlice(reflect.SliceOf(m.infoType), 0, 0)
	m.faceCount = 0
}
