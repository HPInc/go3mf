package meshinfo

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_newFacesData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)

	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		b    *FacesData
		args args
	}{
		{"success", newFacesData(mockContainer), args{4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInvalidator := NewMockFaceData(mockCtrl)
			mockContainer.EXPECT().FaceData(tt.args.faceIndex).Return(mockInvalidator)
			mockInvalidator.EXPECT().Invalidate()
			tt.b.resetFaceInformation(tt.args.faceIndex)
		})
	}
}

func TestFacesData_Clear(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	tests := []struct {
		name    string
		b       *FacesData
		faceNum uint32
	}{
		{"empty", newFacesData(mockContainer), 0},
		{"one", newFacesData(mockContainer), 1},
		{"two", newFacesData(mockContainer), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInvalidator := NewMockFaceData(mockCtrl)
			mockContainer.EXPECT().FaceCount().Return(tt.faceNum)
			mockContainer.EXPECT().FaceData(gomock.Any()).Return(mockInvalidator).Times(int(tt.faceNum))
			mockInvalidator.EXPECT().Invalidate().Times(int(tt.faceNum))
			tt.b.Clear()
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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer2 := NewMockContainer(mockCtrl)
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		b    *FacesData
		args args
		want *FacesData
	}{
		{"base", newFacesData(mockContainer), args{2}, newFacesData(mockContainer2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().clone(tt.args.currentFaceCount).Return(mockContainer2)
			if got := tt.b.clone(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FacesData.clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFacesData_copyFaceInfosFrom(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)

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
		{"success", newFacesData(mockContainer), args{1, NewMockHandleable(mockCtrl), 2}, NewMockFaceData(mockCtrl), NewMockFaceData(mockCtrl)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().FaceData(tt.args.faceIndex).Return(tt.data1)
			tt.args.otherInfo.EXPECT().FaceData(tt.args.otherFaceIndex).Return(tt.data2)
			tt.data1.EXPECT().Copy(tt.data2)
			tt.b.copyFaceInfosFrom(tt.args.faceIndex, tt.args.otherInfo, tt.args.otherFaceIndex)
		})
	}
}

func TestFacesData_permuteNodeInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
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
		{"success", newFacesData(mockContainer), args{1, 2, 3, 4}, NewMockFaceData(mockCtrl)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().FaceData(tt.args.faceIndex).Return(tt.data)
			tt.data.EXPECT().Permute(tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
			tt.b.permuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
		})
	}
}

func TestFacesData_FaceHasData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
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
		{"false", newFacesData(mockContainer), args{1}, NewMockFaceData(mockCtrl), false},
		{"true", newFacesData(mockContainer), args{1}, NewMockFaceData(mockCtrl), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().FaceData(tt.args.faceIndex).Return(tt.data)
			tt.data.EXPECT().HasData().Return(tt.want)
			if got := tt.b.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("FacesData.FaceHasData() = %v, want %v", got, tt.want)
			}
		})
	}
}
