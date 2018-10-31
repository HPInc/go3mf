package meshinfo

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewGenericMeshInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	type args struct {
		container Container
	}
	tests := []struct {
		name string
		args args
		want *GenericMeshInfo
	}{
		{"new1", args{mockContainer}, &GenericMeshInfo{mockContainer, 0}},
		{"new2", args{mockContainer}, &GenericMeshInfo{mockContainer, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGenericMeshInfo(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGenericMeshInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenericMeshInfo_resetFaceInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)

	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name    string
		b       *GenericMeshInfo
		args    args
		wantErr bool
	}{
		{"error", NewGenericMeshInfo(mockContainer), args{2}, true},
		{"success", NewGenericMeshInfo(mockContainer), args{4}, false},
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

func TestGenericMeshInfo_Clear(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	tests := []struct {
		name    string
		b       *GenericMeshInfo
		faceNum uint32
	}{
		{"empty", NewGenericMeshInfo(mockContainer), 0},
		{"one", NewGenericMeshInfo(mockContainer), 1},
		{"two", NewGenericMeshInfo(mockContainer), 2},
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

func TestGenericMeshInfo_setInternalID(t *testing.T) {
	type args struct {
		internalID uint64
	}
	tests := []struct {
		name string
		b    *GenericMeshInfo
		args args
	}{
		{"zero", NewGenericMeshInfo(nil), args{0}},
		{"one", NewGenericMeshInfo(nil), args{1}},
		{"two", NewGenericMeshInfo(nil), args{3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.setInternalID(tt.args.internalID)
			if got := tt.b.getInternalID(); got != tt.args.internalID {
				t.Errorf("GenericMeshInfo.setInternalID() = %v, want %v", got, tt.args.internalID)
			}
		})
	}
}

func TestGenericMeshInfo_getInternalID(t *testing.T) {
	tests := []struct {
		name string
		b    *GenericMeshInfo
		want uint64
	}{
		{"new", NewGenericMeshInfo(nil), 0},
		{"one", &GenericMeshInfo{nil, 1}, 1},
		{"two", &GenericMeshInfo{nil, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.getInternalID(); got != tt.want {
				t.Errorf("GenericMeshInfo.getInternalID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenericMeshInfo_clone(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer2 := NewMockContainer(mockCtrl)
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		b    *GenericMeshInfo
		args args
		want *GenericMeshInfo
	}{
		{"base", NewGenericMeshInfo(mockContainer), args{2}, NewGenericMeshInfo(mockContainer2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().clone(tt.args.currentFaceCount).Return(mockContainer2)
			if got := tt.b.clone(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenericMeshInfo.clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenericMeshInfo_cloneFaceInfosFrom(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)

	type args struct {
		faceIndex      uint32
		otherInfo      *MockMeshInfo
		otherFaceIndex uint32
	}
	tests := []struct {
		name         string
		b            *GenericMeshInfo
		args         args
		data1, data2 *MockFaceData
		err1, err2   error
	}{
		{"err1", NewGenericMeshInfo(mockContainer), args{1, NewMockMeshInfo(mockCtrl), 2}, NewMockFaceData(mockCtrl), NewMockFaceData(mockCtrl), errors.New(""), nil},
		{"err2", NewGenericMeshInfo(mockContainer), args{1, NewMockMeshInfo(mockCtrl), 2}, NewMockFaceData(mockCtrl), NewMockFaceData(mockCtrl), nil, errors.New("")},
		{"success", NewGenericMeshInfo(mockContainer), args{1, NewMockMeshInfo(mockCtrl), 2}, NewMockFaceData(mockCtrl), NewMockFaceData(mockCtrl), nil, nil},
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

func TestGenericMeshInfo_permuteNodeInformation(t *testing.T) {
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
		b    *GenericMeshInfo
		args args
		data *MockFaceData
		err  error
	}{
		{"err", NewGenericMeshInfo(mockContainer), args{1, 2, 3, 4}, NewMockFaceData(mockCtrl), errors.New("")},
		{"success", NewGenericMeshInfo(mockContainer), args{1, 2, 3, 4}, NewMockFaceData(mockCtrl), nil},
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

func TestGenericMeshInfo_FaceHasData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		b    *GenericMeshInfo
		args args
		data *MockFaceData
		err  error
		want bool
	}{
		{"err", NewGenericMeshInfo(mockContainer), args{1}, NewMockFaceData(mockCtrl), errors.New(""), false},
		{"false", NewGenericMeshInfo(mockContainer), args{1}, NewMockFaceData(mockCtrl), nil, false},
		{"true", NewGenericMeshInfo(mockContainer), args{1}, NewMockFaceData(mockCtrl), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(tt.data, tt.err)
			if tt.err == nil {
				tt.data.EXPECT().HasData().Return(tt.want)
			}
			if got := tt.b.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("GenericMeshInfo.FaceHasData() = %v, want %v", got, tt.want)
			}
		})
	}
}
