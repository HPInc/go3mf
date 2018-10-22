package meshinfo

import (
	"github.com/qmuntal/go3mf/internal/common"
	"reflect"
)

// InMemoryMeshInformationContainer implements MeshInformationContainer
// and provides a memory container for holding the texture information state of a complete mesh structure.
type InMemoryMeshInformationContainer struct {
	elemType   reflect.Type
	faceCount  uint32
	dataBlocks reflect.Value
}

// NewInMemoryMeshInformationContainer creates a new container.
func NewInMemoryMeshInformationContainer(currentFaceCount uint32, elemExample FaceData) *InMemoryMeshInformationContainer {
	elemType := reflect.TypeOf(elemExample)
	m := &InMemoryMeshInformationContainer{
		faceCount:  0,
		elemType:   elemType,
		dataBlocks: reflect.MakeSlice(reflect.SliceOf(elemType), 0, int(currentFaceCount)),
	}
	for i := 1; i <= int(currentFaceCount); i++ {
		m.AddFaceData(uint32(i))
	}
	return m
}

// AddFaceData returns the pointer to the data of the added face.
// The parameter newFaceCount should indicate the faces information stored in the container, including the new one.
// Error cases:
// * ErrorMeshInformationCountMismatch: The number of faces in the container does not match with the input parameter.
func (m *InMemoryMeshInformationContainer) AddFaceData(newFaceCount uint32) (val *FaceData, err error) {
	faceData := reflect.New(m.elemType)
	m.dataBlocks = reflect.Append(m.dataBlocks, faceData)
	m.faceCount++
	if m.faceCount != newFaceCount {
		return nil, common.NewError(common.ErrorMeshInformationCountMismatch)
	}
	result := faceData.Elem().Interface().(FaceData)
	return &result, nil
}

// GetFaceData returns the data of the face with the target index.
// Error cases:
// * ErrorInvalidMeshInformationIndex: Index is higher than the number of faces
func (m *InMemoryMeshInformationContainer) GetFaceData(index uint32) (val *FaceData, err error) {
	if index >= m.faceCount {
		return nil, common.NewError(common.ErrorInvalidMeshInformationIndex)
	}

	result := m.dataBlocks.Field(int(index)).Elem().Interface().(FaceData)
	return &result, nil
}

// GetCurrentFaceCount returns the number of faces information stored in the container.
func (m *InMemoryMeshInformationContainer) GetCurrentFaceCount() uint32 {
	return m.faceCount
}

// Clear removes all the information stored in the container.
func (m *InMemoryMeshInformationContainer) Clear() {
	m.dataBlocks = reflect.MakeSlice(reflect.SliceOf(m.elemType), 0, 0)
	m.faceCount = 0
}
