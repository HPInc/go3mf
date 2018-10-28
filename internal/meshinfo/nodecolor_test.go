package meshinfo

import (
	"errors"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNodeColor_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		n    *NodeColor
	}{
		{"base", new(NodeColor)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.Colors[0] = 1
			tt.n.Colors[1] = 2
			tt.n.Colors[2] = 3
			tt.n.Invalidate()
			want := new(NodeColor)
			if !reflect.DeepEqual(tt.n, want) {
				t.Errorf("NodeColor.Invalidate() = %v, want %v", tt.n, want)
			}
		})
	}
}

func TestNewnodeColorsMeshInfo(t *testing.T) {
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
		want *nodeColorsMeshInfo
	}{
		{"new", args{mockContainer}, &nodeColorsMeshInfo{*newbaseMeshInfo(mockContainer)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newnodeColorsMeshInfo(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newnodeColorsMeshInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeColorsMeshInfo_GetType(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear()
	tests := []struct {
		name string
		p    *nodeColorsMeshInfo
		want InformationType
	}{
		{"InfoNodeColors", newnodeColorsMeshInfo(mockContainer), InfoNodeColors},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeColorsMeshInfo.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeColorsMeshInfo_FaceHasData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear().MaxTimes(6)
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name    string
		p       *nodeColorsMeshInfo
		args    args
		wantErr bool
		color   *NodeColor
		want    bool
	}{
		{"error", newnodeColorsMeshInfo(mockContainer), args{0}, true, &NodeColor{[3]Color{1, 2, 3}}, false},
		{"nocolor1", newnodeColorsMeshInfo(mockContainer), args{0}, false, &NodeColor{[3]Color{0, 0, 0}}, false},
		{"nocolor1", newnodeColorsMeshInfo(mockContainer), args{0}, false, &NodeColor{[3]Color{0, 2, 3}}, true},
		{"nocolor2", newnodeColorsMeshInfo(mockContainer), args{0}, false, &NodeColor{[3]Color{1, 0, 3}}, true},
		{"nocolor3", newnodeColorsMeshInfo(mockContainer), args{0}, false, &NodeColor{[3]Color{1, 2, 0}}, true},
		{"data", newnodeColorsMeshInfo(mockContainer), args{0}, false, &NodeColor{[3]Color{1, 2, 3}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.wantErr {
				err = errors.New("")
			}
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(tt.color, err)
			if got := tt.p.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("nodeColorsMeshInfo.FaceHasData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeColorsMeshInfo_Clone(t *testing.T) {
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
		p    *nodeColorsMeshInfo
		args args
		want MeshInfo
	}{
		{"base", newnodeColorsMeshInfo(mockContainer), args{2}, &nodeColorsMeshInfo{*newbaseMeshInfo(mockContainer2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().Clone(tt.args.currentFaceCount).Return(mockContainer2)
			if got := tt.p.Clone(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeColorsMeshInfo.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeColorsMeshInfo_cloneFaceInfosFrom(t *testing.T) {
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
		p            *nodeColorsMeshInfo
		args         args
		want1, want2 *NodeColor
		err1, err2   error
	}{
		{"err1", newnodeColorsMeshInfo(mockContainer1), args{1, newnodeColorsMeshInfo(mockContainer2), 2}, &NodeColor{[3]Color{1, 2, 3}}, &NodeColor{[3]Color{4, 5, 6}}, errors.New(""), nil},
		{"err2", newnodeColorsMeshInfo(mockContainer1), args{1, newnodeColorsMeshInfo(mockContainer2), 2}, &NodeColor{[3]Color{1, 2, 3}}, &NodeColor{[3]Color{4, 5, 6}}, nil, errors.New("")},
		{"permuted", newnodeColorsMeshInfo(mockContainer1), args{1, newnodeColorsMeshInfo(mockContainer2), 2}, &NodeColor{[3]Color{1, 2, 3}}, &NodeColor{[3]Color{4, 5, 6}}, nil, nil},
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
					t.Error("nodeColorsMeshInfo.cloneFaceInfosFrom() modified face data when it shouldn't (1)")
				}
			} else if tt.err2 != nil {
				if reflect.DeepEqual(tt.want1, tt.want2) {
					t.Error("nodeColorsMeshInfo.cloneFaceInfosFrom() modified face data when it shouldn't (2)")
				}
			} else if !reflect.DeepEqual(tt.want1, tt.want2) {
				t.Errorf("nodeColorsMeshInfo.cloneFaceInfosFrom() = %v, want %v", tt.want1, tt.want2)
			}
		})
	}
}

func TestNodeColorsMeshInfo_permuteNodeInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear().MaxTimes(6)
	type args struct {
		faceIndex  uint32
		nodeIndex1 uint32
		nodeIndex2 uint32
		nodeIndex3 uint32
	}
	tests := []struct {
		name    string
		p       *nodeColorsMeshInfo
		args    args
		wantErr bool
		data    *NodeColor
		want    *NodeColor
	}{
		{"err", newnodeColorsMeshInfo(mockContainer), args{1, 2, 1, 0}, true, &NodeColor{[3]Color{1, 2, 0}}, &NodeColor{[3]Color{1, 2, 0}}},
		{"index1", newnodeColorsMeshInfo(mockContainer), args{1, 3, 1, 0}, false, &NodeColor{[3]Color{1, 2, 0}}, &NodeColor{[3]Color{1, 2, 0}}},
		{"index2", newnodeColorsMeshInfo(mockContainer), args{1, 2, 3, 0}, false, &NodeColor{[3]Color{1, 2, 0}}, &NodeColor{[3]Color{1, 2, 0}}},
		{"index3", newnodeColorsMeshInfo(mockContainer), args{1, 2, 2, 3}, false, &NodeColor{[3]Color{1, 2, 0}}, &NodeColor{[3]Color{1, 2, 0}}},
		{"equal", newnodeColorsMeshInfo(mockContainer), args{1, 0, 1, 2}, false, &NodeColor{[3]Color{1, 2, 0}}, &NodeColor{[3]Color{1, 2, 0}}},
		{"diff", newnodeColorsMeshInfo(mockContainer), args{1, 2, 0, 1}, false, &NodeColor{[3]Color{4, 3, 1}}, &NodeColor{[3]Color{1, 4, 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.wantErr {
				err = errors.New("")
			}
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(tt.data, err)
			tt.p.permuteNodeInformation(tt.args.faceIndex, tt.args.nodeIndex1, tt.args.nodeIndex2, tt.args.nodeIndex3)
			if !reflect.DeepEqual(tt.data, tt.want) {
				t.Errorf("nodeColorsMeshInfo.permuteNodeInformation() = %v, want %v", tt.data, tt.want)
			}
		})
	}
}

func TestNodeColorsMeshInfo_mergeInformationFrom(t *testing.T) {
	type args struct {
		info MeshInfo
	}
	tests := []struct {
		name string
		p    *nodeColorsMeshInfo
		args args
	}{
		{"nothing happens", &nodeColorsMeshInfo{baseMeshInfo{nil, 0}}, args{nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.mergeInformationFrom(tt.args.info)
		})
	}
}
