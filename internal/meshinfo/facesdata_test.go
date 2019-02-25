package meshinfo

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func Test_newFacesData(t *testing.T) {
	mockContainer := new(MockContainer)
	type args struct {
		container Container
	}
	tests := []struct {
		name string
		args args
		want *FacesData
	}{
		{"new1", args{mockContainer}, &FacesData{mockContainer, 0}},
		{"new2", args{mockContainer}, &FacesData{mockContainer, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newFacesData(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFacesData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFacesData_resetFaceInformation(t *testing.T) {
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		b    *FacesData
		args args
	}{
		{"success", new(FacesData), args{4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer := new(MockContainer)
			tt.b.Container = mockContainer
			mockInvalidator := new(MockFaceData)
			mockContainer.On("FaceData", tt.args.faceIndex).Return(mockInvalidator)
			mockInvalidator.On("Invalidate")
			tt.b.resetFaceInformation(tt.args.faceIndex)
			mockInvalidator.AssertExpectations(t)
			mockContainer.AssertExpectations(t)

		})
	}
}

func TestFacesData_Clear(t *testing.T) {
	tests := []struct {
		name    string
		b       *FacesData
		faceNum uint32
	}{
		{"empty", new(FacesData), 0},
		{"one", new(FacesData), 1},
		{"two", new(FacesData), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer := new(MockContainer)
			tt.b.Container = mockContainer
			mockInvalidator := new(MockFaceData)
			mockContainer.On("FaceCount").Return(tt.faceNum)
			mockContainer.On("FaceData", mock.Anything).Maybe().Return(mockInvalidator).Times(int(tt.faceNum))
			mockInvalidator.On("Invalidate").Maybe().Times(int(tt.faceNum))
			tt.b.Clear()
			mockInvalidator.AssertExpectations(t)
			mockContainer.AssertExpectations(t)
		})
	}
}

func TestFacesData_setInternalID(t *testing.T) {
	type args struct {
		internalID uint64
	}
	tests := []struct {
		name string
		b    *FacesData
		args args
	}{
		{"zero", newFacesData(nil), args{0}},
		{"one", newFacesData(nil), args{1}},
		{"two", newFacesData(nil), args{3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.setInternalID(tt.args.internalID)
			if got := tt.b.getInternalID(); got != tt.args.internalID {
				t.Errorf("FacesData.setInternalID() = %v, want %v", got, tt.args.internalID)
			}
		})
	}
}

func TestFacesData_getInternalID(t *testing.T) {
	tests := []struct {
		name string
		b    *FacesData
		want uint64
	}{
		{"new", newFacesData(nil), 0},
		{"one", &FacesData{nil, 1}, 1},
		{"two", &FacesData{nil, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.getInternalID(); got != tt.want {
				t.Errorf("FacesData.getInternalID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFacesData_clone(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		b    *FacesData
		args args
		want *FacesData
	}{
		{"base", new(FacesData), args{2}, new(FacesData)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer := new(MockContainer)
			mockContainer2 := new(MockContainer)
			tt.b.Container = mockContainer
			tt.want.Container = mockContainer2
			mockContainer.On("clone", tt.args.currentFaceCount).Return(mockContainer2)
			if got := tt.b.clone(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FacesData.clone() = %v, want %v", got, tt.want)
			}
			mockContainer.AssertExpectations(t)
			mockContainer2.AssertExpectations(t)
		})
	}
}

func TestFacesData_copyFaceInfosFrom(t *testing.T) {
	type args struct {
		faceIndex      uint32
		otherInfo      *MockHandleable
		otherFaceIndex uint32
	}
	tests := []struct {
		name         string
		b            *FacesData
		args         args
		data1, data2 *MockFaceData
	}{
		{"success", new(FacesData), args{1, new(MockHandleable), 2}, new(MockFaceData), new(MockFaceData)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer := new(MockContainer)
			tt.b.Container = mockContainer
			mockContainer.On("FaceData", tt.args.faceIndex).Return(tt.data1)
			tt.args.otherInfo.On("FaceData", tt.args.otherFaceIndex).Return(tt.data2)
			tt.data1.On("Copy", tt.data2)
			tt.b.copyFaceInfosFrom(tt.args.faceIndex, tt.args.otherInfo, tt.args.otherFaceIndex)
			tt.args.otherInfo.AssertExpectations(t)
			tt.data1.AssertExpectations(t)
			mockContainer.AssertExpectations(t)
		})
	}
}

func TestFacesData_permuteNodeInformation(t *testing.T) {
	type args struct {
		faceIndex  uint32
		nodeIndex1 uint32
		nodeIndex2 uint32
		nodeIndex3 uint32
	}
	tests := []struct {
		name string
		b    *FacesData
		args args
		data *MockFaceData
	}{
		{"success", new(FacesData), args{1, 2, 3, 4}, new(MockFaceData)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer := new(MockContainer)
			tt.b.Container = mockContainer
			mockContainer.On("FaceData", tt.args.faceIndex).Return(tt.data)
			tt.data.On("Permute", tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
			tt.b.permuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
			tt.data.AssertExpectations(t)
			mockContainer.AssertExpectations(t)
		})
	}
}

func TestFacesData_FaceHasData(t *testing.T) {
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		b    *FacesData
		args args
		data *MockFaceData
		want bool
	}{
		{"false", new(FacesData), args{1}, new(MockFaceData), false},
		{"true", new(FacesData), args{1}, new(MockFaceData), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer := new(MockContainer)
			tt.b.Container = mockContainer
			mockContainer.On("FaceData", tt.args.faceIndex).Return(tt.data)
			tt.data.On("HasData").Return(tt.want)
			if got := tt.b.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("FacesData.FaceHasData() = %v, want %v", got, tt.want)
				return
			}
			tt.data.AssertExpectations(t)
			mockContainer.AssertExpectations(t)
		})
	}
}
