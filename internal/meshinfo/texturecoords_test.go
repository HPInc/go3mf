package meshinfo

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	gomock "github.com/golang/mock/gomock"
)

func TestTextureCoords_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		t    *TextureCoords
	}{
		{"base", new(TextureCoords)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.TextureID = 2
			tt.t.Coords[0] = mgl32.Vec2{1.0, 2.0}
			tt.t.Coords[1] = mgl32.Vec2{5.0, 3.0}
			tt.t.Coords[2] = mgl32.Vec2{6.0, 4.0}
			tt.t.Invalidate()
			want := new(TextureCoords)
			if !reflect.DeepEqual(tt.t, want) {
				t.Errorf("TextureCoords.Invalidate() = %v, want %v", tt.t, want)
			}
		})
	}
}

func TestNewtextureCoordsMeshInfo(t *testing.T) {
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
		want *textureCoordsMeshInfo
	}{
		{"new", args{mockContainer}, &textureCoordsMeshInfo{*newbaseMeshInfo(mockContainer)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newtextureCoordsMeshInfo(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newtextureCoordsMeshInfo() = %v, want %v", got, tt.want)
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
		p    *textureCoordsMeshInfo
		want InformationType
	}{
		{"InfoTextureCoords", newtextureCoordsMeshInfo(mockContainer), InfoTextureCoords},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("textureCoordsMeshInfo.GetType() = %v, want %v", got, tt.want)
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
		name    string
		p       *textureCoordsMeshInfo
		args    args
		wantErr bool
		coords  *TextureCoords
		want    bool
	}{
		{"error", newtextureCoordsMeshInfo(mockContainer), args{0}, true, &TextureCoords{}, false},
		{"nodata", newtextureCoordsMeshInfo(mockContainer), args{0}, false, &TextureCoords{}, false},
		{"data", newtextureCoordsMeshInfo(mockContainer), args{0}, false, &TextureCoords{TextureID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.wantErr {
				err = errors.New("")
			}
			mockContainer.EXPECT().GetFaceData(tt.args.faceIndex).Return(tt.coords, err)
			if got := tt.p.FaceHasData(tt.args.faceIndex); got != tt.want {
				t.Errorf("textureCoordsMeshInfo.FaceHasData() = %v, want %v", got, tt.want)
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
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		p    *textureCoordsMeshInfo
		args args
		want MeshInfo
	}{
		{"base", newtextureCoordsMeshInfo(mockContainer), args{2}, &textureCoordsMeshInfo{*newbaseMeshInfo(mockContainer2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContainer.EXPECT().Clone(tt.args.currentFaceCount).Return(mockContainer2)
			if got := tt.p.Clone(tt.args.currentFaceCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("textureCoordsMeshInfo.Clone() = %v, want %v", got, tt.want)
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
	source := &TextureCoords{TextureID: 4}
	source.Coords[0] = mgl32.Vec2{1.0, 3.0}
	source.Coords[1] = mgl32.Vec2{0.0, 2.0}
	source.Coords[2] = mgl32.Vec2{0.0, 0.0}
	type args struct {
		faceIndex      uint32
		otherInfo      MeshInfo
		otherFaceIndex uint32
	}
	tests := []struct {
		name         string
		p            *textureCoordsMeshInfo
		args         args
		want1, want2 *TextureCoords
		err1, err2   error
	}{
		{"err1", newtextureCoordsMeshInfo(mockContainer1), args{1, newtextureCoordsMeshInfo(mockContainer2), 2}, &TextureCoords{TextureID: 1}, source, errors.New(""), nil},
		{"err2", newtextureCoordsMeshInfo(mockContainer1), args{1, newtextureCoordsMeshInfo(mockContainer2), 2}, &TextureCoords{TextureID: 1}, source, nil, errors.New("")},
		{"permuted", newtextureCoordsMeshInfo(mockContainer1), args{1, newtextureCoordsMeshInfo(mockContainer2), 2}, &TextureCoords{TextureID: 1}, source, nil, nil},
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
					t.Error("textureCoordsMeshInfo.cloneFaceInfosFrom() modified face data when it shouldn't (1)")
				}
			} else if tt.err2 != nil {
				if reflect.DeepEqual(tt.want1, tt.want2) {
					t.Error("textureCoordsMeshInfo.cloneFaceInfosFrom() modified face data when it shouldn't (2)")
				}
			} else if !reflect.DeepEqual(tt.want1, tt.want2) {
				t.Errorf("textureCoordsMeshInfo.cloneFaceInfosFrom() = %v, want %v", tt.want1, tt.want2)
			}
		})
	}
}

func TestTextureCoordsMeshInfo_permuteNodeInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear().MaxTimes(6)
	source := &TextureCoords{TextureID: 4}
	source.Coords[0] = mgl32.Vec2{1.0, 3.0}
	source.Coords[1] = mgl32.Vec2{0.0, 2.0}
	source.Coords[2] = mgl32.Vec2{0.0, 0.0}
	target := &TextureCoords{TextureID: 4}
	target.Coords[0] = mgl32.Vec2{0.0, 0.0}
	target.Coords[1] = mgl32.Vec2{1.0, 3.0}
	target.Coords[2] = mgl32.Vec2{0.0, 2.0}
	type args struct {
		faceIndex  uint32
		nodeIndex1 uint32
		nodeIndex2 uint32
		nodeIndex3 uint32
	}
	tests := []struct {
		name    string
		p       *textureCoordsMeshInfo
		args    args
		wantErr bool
		data    *TextureCoords
		want    *TextureCoords
	}{
		{"err", newtextureCoordsMeshInfo(mockContainer), args{1, 2, 1, 0}, true, &TextureCoords{TextureID: 1}, &TextureCoords{TextureID: 1}},
		{"index1", newtextureCoordsMeshInfo(mockContainer), args{1, 3, 1, 0}, false, &TextureCoords{TextureID: 1}, &TextureCoords{TextureID: 1}},
		{"index2", newtextureCoordsMeshInfo(mockContainer), args{1, 2, 3, 0}, false, &TextureCoords{TextureID: 1}, &TextureCoords{TextureID: 1}},
		{"index3", newtextureCoordsMeshInfo(mockContainer), args{1, 2, 2, 3}, false, &TextureCoords{TextureID: 1}, &TextureCoords{TextureID: 1}},
		{"equal", newtextureCoordsMeshInfo(mockContainer), args{1, 0, 1, 2}, false, &TextureCoords{TextureID: 1}, &TextureCoords{TextureID: 1}},
		{"diff", newtextureCoordsMeshInfo(mockContainer), args{1, 2, 0, 1}, false, source, target},
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

func TestTextureCoordsMeshInfo_mergeInformationFrom(t *testing.T) {
	type args struct {
		info MeshInfo
	}
	tests := []struct {
		name string
		p    *textureCoordsMeshInfo
		args args
	}{
		{"nothing happens", &textureCoordsMeshInfo{baseMeshInfo{nil, 0}}, args{nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.mergeInformationFrom(tt.args.info)
		})
	}
}
