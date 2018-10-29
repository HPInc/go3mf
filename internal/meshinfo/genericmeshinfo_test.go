package meshinfo

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_newgenericMeshInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	type args struct {
		container Container
	}
	tests := []struct {
		name string
		args args
		want MeshInfo
	}{
		{"new1", args{mockContainer}, &genericMeshInfo{mockContainer, 0}},
		{"new2", args{mockContainer}, &genericMeshInfo{mockContainer, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newgenericMeshInfo(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newgenericMeshInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genericMeshInfo_resetFaceInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)

	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name    string
		b       MeshInfo
		args    args
		wantErr bool
	}{
		{"error", newgenericMeshInfo(mockContainer), args{2}, true},
		{"success", newgenericMeshInfo(mockContainer), args{4}, false},
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

func Test_genericMeshInfo_Clear(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	tests := []struct {
		name    string
		b       MeshInfo
		faceNum uint32
	}{
		{"empty", newgenericMeshInfo(mockContainer), 0},
		{"one", newgenericMeshInfo(mockContainer), 1},
		{"two", newgenericMeshInfo(mockContainer), 2},
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

func Test_genericMeshInfo_setInternalID(t *testing.T) {
	type args struct {
		internalID uint64
	}
	tests := []struct {
		name string
		b    MeshInfo
		args args
	}{
		{"zero", newgenericMeshInfo(nil), args{0}},
		{"one", newgenericMeshInfo(nil), args{1}},
		{"two", newgenericMeshInfo(nil), args{3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.setInternalID(tt.args.internalID)
			if got := tt.b.getInternalID(); got != tt.args.internalID {
				t.Errorf("genericMeshInfo.setInternalID() = %v, want %v", got, tt.args.internalID)
			}
		})
	}
}

func Test_genericMeshInfo_getInternalID(t *testing.T) {
	tests := []struct {
		name string
		b    MeshInfo
		want uint64
	}{
		{"new", newgenericMeshInfo(nil), 0},
		{"one", &genericMeshInfo{nil, 1}, 1},
		{"two", &genericMeshInfo{nil, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.getInternalID(); got != tt.want {
				t.Errorf("genericMeshInfo.getInternalID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genericMeshInfo_Clone(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer2 := NewMockContainer(mockCtrl)
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		b    MeshInfo
		args args
		want MeshInfo
	}{
		{"base", newgenericMeshInfo(mockContainer), args{2}, newgenericMeshInfo(mockContainer2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().Clone(tt.args.currentFaceCount).Return(mockContainer2)
			if got := tt.b.Clone(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genericMeshInfo.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genericMeshInfo_cloneFaceInfosFrom(t *testing.T) {
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
		b            MeshInfo
		args         args
		data1, data2 *MockFaceData
		err1, err2   error
	}{
		{"err1", newgenericMeshInfo(mockContainer), args{1, NewMockMeshInfo(mockCtrl), 2}, NewMockFaceData(mockCtrl), NewMockFaceData(mockCtrl), errors.New(""), nil},
		{"err2", newgenericMeshInfo(mockContainer), args{1, NewMockMeshInfo(mockCtrl), 2}, NewMockFaceData(mockCtrl), NewMockFaceData(mockCtrl), nil, errors.New("")},
		{"success", newgenericMeshInfo(mockContainer), args{1, NewMockMeshInfo(mockCtrl), 2}, NewMockFaceData(mockCtrl), NewMockFaceData(mockCtrl), nil, nil},
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

func Test_genericMeshInfo_permuteNodeInformation(t *testing.T) {
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
		b    MeshInfo
		args args
		data *MockFaceData
		err  error
	}{
		{"err", newgenericMeshInfo(mockContainer), args{1, 2, 3, 4}, NewMockFaceData(mockCtrl), errors.New("")},
		{"success", newgenericMeshInfo(mockContainer), args{1, 2, 3, 4}, NewMockFaceData(mockCtrl), nil},
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

func Test_genericMeshInfo_mergeInformationFrom(t *testing.T) {
	type args struct {
		info MeshInfo
	}
	tests := []struct {
		name string
		b    MeshInfo
		args args
	}{
		{"nothing happens", &genericMeshInfo{nil, 0}, args{nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.mergeInformationFrom(tt.args.info)
		})
	}
}

func Test_genericMeshInfo_FaceHasData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		b    MeshInfo
		args args
		data *MockFaceData
		err  error
		want bool
	}{
		{"err", newgenericMeshInfo(mockContainer), args{1}, NewMockFaceData(mockCtrl), errors.New(""), false},
		{"false", newgenericMeshInfo(mockContainer), args{1}, NewMockFaceData(mockCtrl), nil, false},
		{"true", newgenericMeshInfo(mockContainer), args{1}, NewMockFaceData(mockCtrl), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(tt.data, tt.err)
			if tt.err == nil {
				tt.data.EXPECT().HasData().Return(tt.want)
			}
			if got := tt.b.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("genericMeshInfo.FaceHasData() = %v, want %v", got, tt.want)
			}
		})
	}
}
