package mesh

import (
	"github.com/qmuntal/go3mf/internal/meshinfo"
	"github.com/stretchr/testify/mock"
)

// MockMergeableMesh is a mock of MergeableMesh interface
type MockMergeableMesh struct {
	mock.Mock
}

// Beam mocks base method
func (m *MockMergeableMesh) Beam(arg0 uint32) *Beam {
	args := m.Called(arg0)
	return args.Get(0).(*Beam)
}

// BeamCount mocks base method
func (m *MockMergeableMesh) BeamCount() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

// Face mocks base method
func (m *MockMergeableMesh) Face(arg0 uint32) *Face {
	args := m.Called(arg0)
	return args.Get(0).(*Face)
}

// FaceCount mocks base method
func (m *MockMergeableMesh) FaceCount() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

// InformationHandler mocks base method
func (m *MockMergeableMesh) InformationHandler() *meshinfo.Handler {
	args := m.Called()
	return args.Get(0).(*meshinfo.Handler)
}

// Node mocks base method
func (m *MockMergeableMesh) Node(arg0 uint32) *Node {
	args := m.Called(arg0)
	return args.Get(0).(*Node)
}

// NodeCount mocks base method
func (m *MockMergeableMesh) NodeCount() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}
