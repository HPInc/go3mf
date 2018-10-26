package meshinfo

import (
	"errors"
	"reflect"
	"testing"
	"github.com/golang/mock/gomock"
)

func Test_baseMaterialInvalidator_Invalidate(t *testing.T) {
	expected := &BaseMaterial{0,0}
	type args struct {
		data FaceData
	}
	tests := []struct {
		name string
		p    baseMaterialInvalidator
		args args
	}{
		{"generic", baseMaterialInvalidator{}, args{&fakeFaceData{}}},
		{"specific", baseMaterialInvalidator{}, args{&BaseMaterial{2,1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Invalidate(tt.args.data)
			if got, ok := tt.args.data.(*BaseMaterial); ok {
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("baseMaterialInvalidator.Invalidate expected  = %v, want %v", got, expected)
				}
			}
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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear().MaxTimes(3)
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		p    *BaseMaterialsInfo
		args args
		wantErr bool
		want bool
	}{
		{"error", NewBaseMaterialsInfo(mockContainer), args{0}, true, false},
		{"nodata", NewBaseMaterialsInfo(mockContainer), args{0}, false, false},
		{"nodata", NewBaseMaterialsInfo(mockContainer), args{0}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &BaseMaterial{0,0}
			if tt.want {
				data.MaterialGroupID = 1
			}
			var err error
			if tt.wantErr {
				err = errors.New("")
			}
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(data, err)
			if got := tt.p.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("BaseMaterialsInfo.FaceHasData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsInfo_Clone(t *testing.T) {	
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer2 := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear()
	mockContainer2.EXPECT().Clear()
	mockContainer.EXPECT().Clone().Return(mockContainer2)
	tests := []struct {
		name string
		p    *BaseMaterialsInfo
		want MeshInfo
	}{
		{"base", NewBaseMaterialsInfo(mockContainer), &BaseMaterialsInfo{*newBaseMeshInfo(mockContainer2, baseMaterialInvalidator{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseMaterialsInfo.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsInfo_cloneFaceInfosFrom(t *testing.T) {
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
		name string
		p    *BaseMaterialsInfo
		args args
		want1, want2 *BaseMaterial
		err1, err2 error
	}{
		{"err1", NewBaseMaterialsInfo(mockContainer1), args{1, NewBaseMaterialsInfo(mockContainer2), 2}, &BaseMaterial{2,3}, &BaseMaterial{4,5}, errors.New(""), nil},
		{"err2", NewBaseMaterialsInfo(mockContainer1), args{1, NewBaseMaterialsInfo(mockContainer2), 2}, &BaseMaterial{2,3}, &BaseMaterial{4,5}, nil, errors.New("")},
		{"err2", NewBaseMaterialsInfo(mockContainer1), args{1, NewBaseMaterialsInfo(mockContainer2), 2}, &BaseMaterial{2,3}, &BaseMaterial{4,5}, nil, nil},
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
					t.Error("BaseMaterialsInfo.cloneFaceInfosFrom() modified face data when it shouldn't (1)")
				}
			} else if tt.err2 != nil {
				if reflect.DeepEqual(tt.want1, tt.want2) {
					t.Error("BaseMaterialsInfo.cloneFaceInfosFrom() modified face data when it shouldn't (2)")
				}
			} else if !reflect.DeepEqual(tt.want1, tt.want2) {
				t.Errorf("BaseMaterialsInfo.cloneFaceInfosFrom() = %v, want %v", tt.want1, tt.want2)
			}
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
