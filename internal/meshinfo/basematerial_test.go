package meshinfo

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestBaseMaterial_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseMaterial
	}{
		{"base", &BaseMaterial{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Invalidate()
			want := new(BaseMaterial)
			if !reflect.DeepEqual(tt.b, want) {
				t.Errorf("BaseMaterial.Invalidate() = %v, want %v", tt.b, want)
			}
		})
	}
}

func TestNewbaseMaterialsMeshInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear()
	type args struct {
		container Container
	}
	tests := []struct {
		name string
		args args
		want *baseMaterialsMeshInfo
	}{
		{"new", args{mockContainer}, &baseMaterialsMeshInfo{*newbaseMeshInfo(mockContainer)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newbaseMaterialsMeshInfo(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newbaseMaterialsMeshInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsMeshInfo_GetType(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear()
	tests := []struct {
		name string
		p    *baseMaterialsMeshInfo
		want InformationType
	}{
		{"InfoBaseMaterials", newbaseMaterialsMeshInfo(mockContainer), InfoBaseMaterials},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseMaterialsMeshInfo.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsMeshInfo_FaceHasData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear().MaxTimes(3)
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name    string
		p       *baseMaterialsMeshInfo
		args    args
		wantErr bool
		want    bool
	}{
		{"error", newbaseMaterialsMeshInfo(mockContainer), args{0}, true, false},
		{"nodata", newbaseMaterialsMeshInfo(mockContainer), args{0}, false, false},
		{"data", newbaseMaterialsMeshInfo(mockContainer), args{0}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &BaseMaterial{0, 0}
			if tt.want {
				data.MaterialGroupID = 1
			}
			var err error
			if tt.wantErr {
				err = errors.New("")
			}
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(data, err)
			if got := tt.p.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("baseMaterialsMeshInfo.FaceHasData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsMeshInfo_Clone(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer2 := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear()
	mockContainer2.EXPECT().Clear()
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		p    *baseMaterialsMeshInfo
		args args
		want MeshInfo
	}{
		{"base", newbaseMaterialsMeshInfo(mockContainer), args{2}, &baseMaterialsMeshInfo{*newbaseMeshInfo(mockContainer2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().Clone(tt.args.currentFaceCount).Return(mockContainer2)
			if got := tt.p.Clone(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseMaterialsMeshInfo.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsMeshInfo_cloneFaceInfosFrom(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer1 := NewMockContainer(mockCtrl)
	mockContainer2 := NewMockContainer(mockCtrl)
	mockContainer1.EXPECT().Clear().MaxTimes(3)
	mockContainer2.EXPECT().Clear().MaxTimes(3)
	type args struct {
		faceIndex      uint32
		otherInfo      MeshInfo
		otherFaceIndex uint32
	}
	tests := []struct {
		name         string
		p            *baseMaterialsMeshInfo
		args         args
		want1, want2 *BaseMaterial
		err1, err2   error
	}{
		{"err1", newbaseMaterialsMeshInfo(mockContainer1), args{1, newbaseMaterialsMeshInfo(mockContainer2), 2}, &BaseMaterial{2, 3}, &BaseMaterial{4, 5}, errors.New(""), nil},
		{"err2", newbaseMaterialsMeshInfo(mockContainer1), args{1, newbaseMaterialsMeshInfo(mockContainer2), 2}, &BaseMaterial{2, 3}, &BaseMaterial{4, 5}, nil, errors.New("")},
		{"err2", newbaseMaterialsMeshInfo(mockContainer1), args{1, newbaseMaterialsMeshInfo(mockContainer2), 2}, &BaseMaterial{2, 3}, &BaseMaterial{4, 5}, nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer1.EXPECT().GetFaceData(tt.args.faceIndex).Return(tt.want1, tt.err1)
			if tt.err1 == nil {
				mockContainer2.EXPECT().GetFaceData(tt.args.otherFaceIndex).Return(tt.want2, tt.err2)
			}

			tt.p.cloneFaceInfosFrom(tt.args.faceIndex, tt.args.otherInfo, tt.args.otherFaceIndex)

			if tt.err1 != nil {
				if reflect.DeepEqual(tt.want1, tt.want2) {
					t.Error("baseMaterialsMeshInfo.cloneFaceInfosFrom() modified face data when it shouldn't (1)")
				}
			} else if tt.err2 != nil {
				if reflect.DeepEqual(tt.want1, tt.want2) {
					t.Error("baseMaterialsMeshInfo.cloneFaceInfosFrom() modified face data when it shouldn't (2)")
				}
			} else if !reflect.DeepEqual(tt.want1, tt.want2) {
				t.Errorf("baseMaterialsMeshInfo.cloneFaceInfosFrom() = %v, want %v", tt.want1, tt.want2)
			}
		})
	}
}

func TestBaseMaterialsMeshInfo_permuteNodeInformation(t *testing.T) {
	type args struct {
		faceIndex  uint32
		nodeIndex1 uint32
		nodeIndex2 uint32
		nodeIndex3 uint32
	}
	tests := []struct {
		name string
		p    *baseMaterialsMeshInfo
		args args
	}{
		{"nothing happens", &baseMaterialsMeshInfo{baseMeshInfo{nil, 0}}, args{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.permuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
		})
	}
}

func TestBaseMaterialsMeshInfo_mergeInformationFrom(t *testing.T) {
	type args struct {
		info MeshInfo
	}
	tests := []struct {
		name string
		p    *baseMaterialsMeshInfo
		args args
	}{
		{"nothing happens", &baseMaterialsMeshInfo{baseMeshInfo{nil, 0}}, args{nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.mergeInformationFrom(tt.args.info)
		})
	}
}
