package meshinfo

import (
	"reflect"
	"testing"
	"github.com/golang/mock/gomock"
)

func Test_baseMaterialInvalidator_Invalidate(t *testing.T) {
	type args struct {
		data FaceData
	}
	tests := []struct {
		name string
		p    baseMaterialInvalidator
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Invalidate(tt.args.data)
		})
	}
}

func TestNewBaseMaterialsInfo(t *testing.T) {
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
		want *BaseMaterialsInfo
	}{
		{"new", args{mockContainer}, &BaseMaterialsInfo{*newBaseMeshInfo(mockContainer, baseMaterialInvalidator{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBaseMaterialsInfo(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBaseMaterialsInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsInfo_GetType(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear()
	tests := []struct {
		name string
		p    *BaseMaterialsInfo
		want InformationType
	}{
		{"InfoBaseMaterials", NewBaseMaterialsInfo(mockContainer), InfoBaseMaterials},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseMaterialsInfo.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsInfo_FaceHasData(t *testing.T) {
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		p    *BaseMaterialsInfo
		args args
		want bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("BaseMaterialsInfo.FaceHasData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsInfo_Clone(t *testing.T) {
	tests := []struct {
		name string
		p    *BaseMaterialsInfo
		want MeshInfo
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseMaterialsInfo.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsInfo_cloneFaceInfosFrom(t *testing.T) {
	type args struct {
		faceIndex      uint32
		otherInfo      MeshInfo
		otherFaceIndex uint32
	}
	tests := []struct {
		name string
		p    *BaseMaterialsInfo
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.cloneFaceInfosFrom(tt.args.faceIndex, tt.args.otherInfo, tt.args.otherFaceIndex)
		})
	}
}

func TestBaseMaterialsInfo_permuteNodeInformation(t *testing.T) {
	type args struct {
		faceIndex  uint32
		nodeIndex1 uint32
		nodeIndex2 uint32
		nodeIndex3 uint32
	}
	tests := []struct {
		name string
		p    *BaseMaterialsInfo
		args args
	}{
		{"nothing happens", &BaseMaterialsInfo{baseMeshInfo{nil, nil, 0}}, args{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.permuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
		})
	}
}

func TestBaseMaterialsInfo_mergeInformationFrom(t *testing.T) {
	type args struct {
		info MeshInfo
	}
	tests := []struct {
		name string
		p    *BaseMaterialsInfo
		args args
	}{
		{"nothing happens", &BaseMaterialsInfo{baseMeshInfo{nil, nil, 0}}, args{nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.mergeInformationFrom(tt.args.info)
		})
	}
}
