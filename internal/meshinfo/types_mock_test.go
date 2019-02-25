package meshinfo

import (
	"github.com/stretchr/testify/mock"
)

// MockFaceData is a mock of FaceData interface
type MockFaceData struct {
	mock.Mock
}

// Copy mocks base method
func (m *MockFaceData) Copy(arg0 interface{}) {
	m.Called(arg0)
}

// HasData mocks base method
func (m *MockFaceData) HasData() bool {
	args := m.Called()
	return args.Bool(0)
}

// Invalidate mocks base method
func (m *MockFaceData) Invalidate() {
	m.Called()
}

// Merge mocks base method
func (m *MockFaceData) Merge(arg0 interface{}) {
	m.Called(arg0)
}

// Permute mocks base method
func (m *MockFaceData) Permute(arg0, arg1, arg2 uint32) {
	m.Called(arg0, arg1, arg2)
}

// MockContainer is a mock of Container interface
type MockContainer struct {
	mock.Mock
}

// AddFaceData mocks base method
func (m *MockContainer) AddFaceData(arg0 uint32) FaceData {
	args := m.Called(arg0)
	return args.Get(0).(FaceData)
}

// Clear mocks base method
func (m *MockContainer) Clear() {
	m.Called()
}

// FaceCount mocks base method
func (m *MockContainer) FaceCount() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

// FaceData mocks base method
func (m *MockContainer) FaceData(arg0 uint32) FaceData {
	args := m.Called(arg0)
	return args.Get(0).(FaceData)
}

// InfoType mocks base method
func (m *MockContainer) InfoType() dataType {
	args := m.Called()
	return args.Get(0).(dataType)
}

// clone mocks base method
func (m *MockContainer) clone(arg0 uint32) Container {
	args := m.Called(arg0)
	return args.Get(0).(Container)
}

// MockTypedInformer is a mock of TypedInformer interface
type MockTypedInformer struct {
	mock.Mock
}

// infoTypes mocks base method
func (m *MockTypedInformer) infoTypes() []dataType {
	args := m.Called()
	return args.Get(0).([]dataType)
}

// informationByType mocks base method
func (m *MockTypedInformer) informationByType(arg0 dataType) (Handleable, bool) {
	args := m.Called(arg0)
	return args.Get(0).(Handleable), args.Bool(1)
}

// MockFaceQuerier is a mock of FaceQuerier interface
type MockFaceQuerier struct {
	mock.Mock
}

// FaceData mocks base method
func (m *MockFaceQuerier) FaceData(arg0 uint32) FaceData {
	args := m.Called(arg0)
	return args.Get(0).(FaceData)
}

// MockHandleable is a mock of Handleable interface
type MockHandleable struct {
	mock.Mock
}

// AddFaceData mocks base method
func (m *MockHandleable) AddFaceData(arg0 uint32) FaceData {
	args := m.Called(arg0)
	return args.Get(0).(FaceData)
}

// FaceData mocks base method
func (m *MockHandleable) FaceData(arg0 uint32) FaceData {
	args := m.Called(arg0)
	return args.Get(0).(FaceData)
}

// InfoType mocks base method
func (m *MockHandleable) InfoType() dataType {
	args := m.Called()
	return args.Get(0).(dataType)
}

// clone mocks base method
func (m *MockHandleable) clone(arg0 uint32) Handleable {
	args := m.Called(arg0)
	return args.Get(0).(Handleable)
}

// copyFaceInfosFrom mocks base method
func (m *MockHandleable) copyFaceInfosFrom(arg0 uint32, arg1 FaceQuerier, arg2 uint32) {
	m.Called(arg0, arg1, arg2)
}

// permuteNodeInformation mocks base method
func (m *MockHandleable) permuteNodeInformation(arg0, arg1, arg2, arg3 uint32) {
	m.Called(arg0, arg1, arg2, arg3)
}

// resetFaceInformation mocks base method
func (m *MockHandleable) resetFaceInformation(arg0 uint32) {
	m.Called(arg0)
}

// setInternalID mocks base method
func (m *MockHandleable) setInternalID(arg0 uint64) {
	m.Called(arg0)
}
