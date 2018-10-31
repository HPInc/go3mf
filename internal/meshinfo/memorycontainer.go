package meshinfo

import (
	"reflect"
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
func newmemoryContainer(currentFaceCount uint32, infoType reflect.Type) Container {
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

func (m *memoryContainer) Clone(currentFaceCount uint32) Container {
	return newmemoryContainer(currentFaceCount, m.infoType)
}

// InfoType returns the type of the stored data.
func (m *memoryContainer) InfoType() reflect.Type {
	return m.infoType
}

func (m *memoryContainer) AddFaceData(newFaceCount uint32) (FaceData, error) {
	if m.infoType == nil {
		return nil, &InvalidInfoTypeError{m.infoType}
	}
	faceData := reflect.New(m.infoType.Elem())
	m.dataBlocks = reflect.Append(m.dataBlocks, faceData)
	m.faceCount++
	if m.faceCount != newFaceCount {
		return nil, &FaceCountMissmatchError{m.faceCount, newFaceCount}
	}
	return faceData.Interface().(FaceData), nil
}

func (m *memoryContainer) GetFaceData(faceIndex uint32) (FaceData, error) {
	if faceIndex >= m.faceCount {
		return nil, &FaceDataIndexError{m.faceCount, faceIndex}
	}

	return m.dataBlocks.Index(int(faceIndex)).Interface().(FaceData), nil
}

func (m *memoryContainer) GetCurrentFaceCount() uint32 {
	return m.faceCount
}

// Clear removes all the information stored in the container.
func (m *memoryContainer) Clear() {
	m.dataBlocks = reflect.MakeSlice(reflect.SliceOf(m.infoType), 0, 0)
	m.faceCount = 0
}
