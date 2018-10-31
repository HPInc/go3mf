package meshinfo

import (
	"errors"
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
		name    string
		b       *FacesData
		args    args
		wantErr bool
	}{
		{"error", newFacesData(mockContainer), args{2}, true},
		{"success", newFacesData(mockContainer), args{4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInvalidator := NewMockFaceData(mockCtrl)
			var (
				err   error
				times int
			)
			if tt.wantErr {
				err = errors.New("")
			} else {
				times = 1
			}

			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(mockInvalidator, err)
			mockInvalidator.EXPECT().Invalidate().Times(times)
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
			mockContainer.EXPECT().GetCurrentFaceCount().Return(tt.faceNum)
			mockContainer.EXPECT().GetFaceData(gomock.Any()).Return(mockInvalidator, nil).Times(int(tt.faceNum))
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

func TestFacesData_cloneFaceInfosFrom(t *testing.T) {
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
		err1, err2   error
	}{
		{"err1", newFacesData(mockContainer), args{1, NewMockHandleable(mockCtrl), 2}, NewMockFaceData(mockCtrl), NewMockFaceData(mockCtrl), errors.New(""), nil},
		{"err2", newFacesData(mockContainer), args{1, NewMockHandleable(mockCtrl), 2}, NewMockFaceData(mockCtrl), NewMockFaceData(mockCtrl), nil, errors.New("")},
		{"success", newFacesData(mockContainer), args{1, NewMockHandleable(mockCtrl), 2}, NewMockFaceData(mockCtrl), NewMockFaceData(mockCtrl), nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(tt.data1, tt.err1)
			if tt.err1 == nil {
				tt.args.otherInfo.EXPECT().GetFaceData(tt.args.otherFaceIndex).Return(tt.data2, tt.err2)
			}
			if tt.err1 == nil && tt.err2 == nil {
				tt.data1.EXPECT().Copy(tt.data2)
			}
			tt.b.cloneFaceInfosFrom(tt.args.faceIndex, tt.args.otherInfo, tt.args.otherFaceIndex)
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
		err  error
	}{
		{"err", newFacesData(mockContainer), args{1, 2, 3, 4}, NewMockFaceData(mockCtrl), errors.New("")},
		{"success", newFacesData(mockContainer), args{1, 2, 3, 4}, NewMockFaceData(mockCtrl), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(tt.data, tt.err)
			if tt.err == nil {
				tt.data.EXPECT().Permute(tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
			}
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
		err  error
		want bool
	}{
		{"err", newFacesData(mockContainer), args{1}, NewMockFaceData(mockCtrl), errors.New(""), false},
		{"false", newFacesData(mockContainer), args{1}, NewMockFaceData(mockCtrl), nil, false},
		{"true", newFacesData(mockContainer), args{1}, NewMockFaceData(mockCtrl), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(tt.data, tt.err)
			if tt.err == nil {
				tt.data.EXPECT().HasData().Return(tt.want)
			}
			if got := tt.b.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("FacesData.FaceHasData() = %v, want %v", got, tt.want)
			}
		})
	}
}
