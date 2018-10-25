package meshinfo

import (
	"reflect"

	"github.com/qmuntal/go3mf/internal/common"
)

// inMemoryMeshInformationContainer implements MeshInformationContainer
// and provides a generic memory container for holding the texture information state of a complete mesh structure
// using reflection to infer slyce type.
type inMemoryMeshInformationContainer struct {
	elemType   reflect.Type
	faceCount  uint32
	dataBlocks reflect.Value
}

// newInMemoryMeshInformationContainer creates a new container that holds the specified element types.
// Error cases:
// * ErrorInvalidRecordSize: The element type is not defined.
func newInMemoryMeshInformationContainer(currentFaceCount uint32, elemType reflect.Type) *inMemoryMeshInformationContainer {
	m := &inMemoryMeshInformationContainer{
		faceCount:  0,
		elemType:   elemType,
		dataBlocks: reflect.MakeSlice(reflect.SliceOf(elemType), 0, int(currentFaceCount)),
	}
	for i := 1; i <= int(currentFaceCount); i++ {
		m.AddFaceData(uint32(i))
	}
	return m
}

// Clone creates a copy of the container with all the faces invalidated.
func (m *inMemoryMeshInformationContainer) clone() MeshInformationContainer {
	return newInMemoryMeshInformationContainer(m.faceCount, m.elemType)
}

// AddFaceData returns the pointer to the data of the added face.
// The parameter newFaceCount should indicate the faces information stored in the container, including the new one.
// Error cases:
// * ErrorInvalidRecordSize: The element type is not defined.
// * ErrorMeshInformationCountMismatch: The number of faces in the container does not match with the input parameter.
func (m *inMemoryMeshInformationContainer) AddFaceData(newFaceCount uint32) (FaceData, error) {
	if m.elemType == nil {
		return nil, common.NewError(common.ErrorInvalidRecordSize)
	}
	faceData := reflect.New(m.elemType)
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
func (m *inMemoryMeshInformationContainer) GetFaceData(faceIndex uint32) (FaceData, error) {
	if faceIndex >= m.faceCount {
		return nil, common.NewError(common.ErrorInvalidMeshInformationIndex)
	}

	return m.dataBlocks.Index(int(faceIndex)).Addr().Interface().(FaceData), nil
}

// GetCurrentFaceCount returns the number of faces information stored in the container.
func (m *inMemoryMeshInformationContainer) GetCurrentFaceCount() uint32 {
	return m.faceCount
}

// Clear removes all the information stored in the container.
func (m *inMemoryMeshInformationContainer) Clear() {
	m.dataBlocks = reflect.MakeSlice(reflect.SliceOf(m.elemType), 0, 0)
	m.faceCount = 0
}
