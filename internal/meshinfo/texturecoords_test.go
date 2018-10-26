package meshinfo

import (
	"reflect"
	"testing"
	"errors"
	"github.com/go-gl/mathgl/mgl32"
	gomock "github.com/golang/mock/gomock"
)

func TestNewTextureCoords(t *testing.T) {
	type args struct {
		textureID uint32
	}
	tests := []struct {
		name string
		args args
		want *TextureCoords
	}{
		{"new", args{1}, &TextureCoords{1, [3]mgl32.Vec2{mgl32.Vec2{0.0, 0.0}, mgl32.Vec2{0.0, 0.0}, mgl32.Vec2{0.0, 0.0}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextureCoords(tt.args.textureID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextureCoords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_textureCoordsInvalidator_Invalidate(t *testing.T) {
	type args struct {
		data FaceData
	}
	tests := []struct {
		name string
		p    textureCoordsInvalidator
		args args
	}{
		{"generic", textureCoordsInvalidator{}, args{&fakeFaceData{}}},
		{"specific", textureCoordsInvalidator{}, args{NewTextureCoords(4)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Invalidate(tt.args.data)
		})
	}
}

func TestNewTextureCoordsMeshInfo(t *testing.T) {
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
		want *TextureCoordsMeshInfo
	}{
		{"new", args{mockContainer}, &TextureCoordsMeshInfo{*newBaseMeshInfo(mockContainer, textureCoordsInvalidator{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextureCoordsMeshInfo(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextureCoordsMeshInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextureCoordsMeshInfo_GetType(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear()
	tests := []struct {
		name string
		p    *TextureCoordsMeshInfo
		want InformationType
	}{
		{"InfoTextureCoords", NewTextureCoordsMeshInfo(mockContainer), InfoTextureCoords},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TextureCoordsMeshInfo.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextureCoordsMeshInfo_FaceHasData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear().MaxTimes(3)
	type args struct {
		faceIndex uint32
	}
	tests := []struct {
		name string
		p    *TextureCoordsMeshInfo
		args args
		wantErr bool
		coords   *TextureCoords
		want bool
	}{
		{"error", NewTextureCoordsMeshInfo(mockContainer), args{0}, true, NewTextureCoords(0), false},
		{"nodata", NewTextureCoordsMeshInfo(mockContainer), args{0}, false, NewTextureCoords(0), false},
		{"data", NewTextureCoordsMeshInfo(mockContainer), args{0}, false, NewTextureCoords(1), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.wantErr {
				err = errors.New("")
			}
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(tt.coords, err)
			if got := tt.p.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("TextureCoordsMeshInfo.FaceHasData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextureCoordsMeshInfo_Clone(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer2 := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear()
	mockContainer2.EXPECT().Clear()
	mockContainer.EXPECT().Clone().Return(mockContainer2)
	tests := []struct {
		name string
		p    *TextureCoordsMeshInfo
		want MeshInfo
	}{
		{"base", NewTextureCoordsMeshInfo(mockContainer), &TextureCoordsMeshInfo{*newBaseMeshInfo(mockContainer2, textureCoordsInvalidator{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TextureCoordsMeshInfo.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextureCoordsMeshInfo_cloneFaceInfosFrom(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer1 := NewMockContainer(mockCtrl)
	mockContainer2 := NewMockContainer(mockCtrl)
	mockContainer1.EXPECT().Clear().MaxTimes(3)
	mockContainer2.EXPECT().Clear().MaxTimes(3)
	source := NewTextureCoords(4)
	source.Coords[0] = mgl32.Vec2{1.0,3.0}
	source.Coords[1] = mgl32.Vec2{0.0,2.0}
	source.Coords[2] = mgl32.Vec2{0.0,0.0}
	type args struct {
		faceIndex      uint32
		otherInfo      MeshInfo
		otherFaceIndex uint32
	}
	tests := []struct {
		name string
		p    *TextureCoordsMeshInfo
		args args
		want1, want2 *TextureCoords
		err1, err2   error
	}{
		{"err1", NewTextureCoordsMeshInfo(mockContainer1), args{1, NewTextureCoordsMeshInfo(mockContainer2), 2}, NewTextureCoords(1), source, errors.New(""), nil},
		{"err2", NewTextureCoordsMeshInfo(mockContainer1), args{1, NewTextureCoordsMeshInfo(mockContainer2), 2}, NewTextureCoords(1), source, nil, errors.New("")},
		{"permuted", NewTextureCoordsMeshInfo(mockContainer1), args{1, NewTextureCoordsMeshInfo(mockContainer2), 2}, NewTextureCoords(1), source, nil, nil},
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
					t.Error("TextureCoordsMeshInfo.cloneFaceInfosFrom() modified face data when it shouldn't (1)")
				}
			} else if tt.err2 != nil {
				if reflect.DeepEqual(tt.want1, tt.want2) {
					t.Error("TextureCoordsMeshInfo.cloneFaceInfosFrom() modified face data when it shouldn't (2)")
				}
			} else if !reflect.DeepEqual(tt.want1, tt.want2) {
				t.Errorf("TextureCoordsMeshInfo.cloneFaceInfosFrom() = %v, want %v", tt.want1, tt.want2)
			}
		})
	}
}

func TestTextureCoordsMeshInfo_permuteNodeInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear().MaxTimes(6)
	source := NewTextureCoords(4)
	source.Coords[0] = mgl32.Vec2{1.0,3.0}
	source.Coords[1] = mgl32.Vec2{0.0,2.0}
	source.Coords[2] = mgl32.Vec2{0.0,0.0}
	target := NewTextureCoords(4)
	target.Coords[0] = mgl32.Vec2{0.0,0.0}
	target.Coords[1] = mgl32.Vec2{1.0,3.0}
	target.Coords[2] = mgl32.Vec2{0.0,2.0}
	type args struct {
		faceIndex  uint32
		nodeIndex1 uint32
		nodeIndex2 uint32
		nodeIndex3 uint32
	}
	tests := []struct {
		name string
		p    *TextureCoordsMeshInfo
		args args
		wantErr bool
		data    *TextureCoords
		want    *TextureCoords
	}{
		{"err", NewTextureCoordsMeshInfo(mockContainer), args{1, 2, 1, 0}, true, NewTextureCoords(1), NewTextureCoords(1)},
		{"index1", NewTextureCoordsMeshInfo(mockContainer), args{1, 3, 1, 0}, false, NewTextureCoords(1), NewTextureCoords(1)},
		{"index2", NewTextureCoordsMeshInfo(mockContainer), args{1, 2, 3, 0}, false, NewTextureCoords(1), NewTextureCoords(1)},
		{"index3", NewTextureCoordsMeshInfo(mockContainer), args{1, 2, 2, 3}, false, NewTextureCoords(1), NewTextureCoords(1)},
		{"equal", NewTextureCoordsMeshInfo(mockContainer), args{1, 0, 1, 2}, false, NewTextureCoords(1), NewTextureCoords(1)},
		{"diff", NewTextureCoordsMeshInfo(mockContainer), args{1, 2, 0, 1}, false, source, target},
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
				t.Errorf("NodeColorsMeshInfo.permuteNodeInformation() = %v, want %v", tt.data, tt.want)
			}
		})
	}
}

func TestTextureCoordsMeshInfo_mergeInformationFrom(t *testing.T) {
	type args struct {
		info MeshInfo
	}
	tests := []struct {
		name string
		p    *TextureCoordsMeshInfo
		args args
	}{
		{"nothing happens", &TextureCoordsMeshInfo{baseMeshInfo{nil, nil, 0}}, args{nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.mergeInformationFrom(tt.args.info)
		})
	}
}
