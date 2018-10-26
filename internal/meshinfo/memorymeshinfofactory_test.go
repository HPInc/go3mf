package meshinfo

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNewMemoryMeshInfoFactory(t *testing.T) {
	tests := []struct {
		name string
		want *MemoryMeshInfoFactory
	}{
		{"base", &MemoryMeshInfoFactory{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryMeshInfoFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoryMeshInfoFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryMeshInfoFactory_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := NewMockContainer(mockCtrl)
	mockContainer.EXPECT().Clear().MaxTimes(3)
	type args struct {
		infoType         InformationType
		currentFaceCount uint32
	}
	tests := []struct {
		name    string
		f       *MemoryMeshInfoFactory
		args    args
		want    MeshInfo
		wantErr bool
	}{
		{"error", NewMemoryMeshInfoFactory(), args{InfoAbstract, 0}, nil, true},
		{"basematerials", NewMemoryMeshInfoFactory(), args{InfoBaseMaterials, 0}, newbaseMaterialsMeshInfo(mockContainer), false},
		{"nodecolors", NewMemoryMeshInfoFactory(), args{InfoNodeColors, 0}, newnodeColorsMeshInfo(mockContainer), false},
		{"textureCoords", NewMemoryMeshInfoFactory(), args{InfoTextureCoords, 0}, newtextureCoordsMeshInfo(mockContainer), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Create(tt.args.infoType, tt.args.currentFaceCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryMeshInfoFactory.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(reflect.TypeOf(got) == reflect.TypeOf(tt.want)) {
				t.Errorf("MemoryMeshInfoFactory.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
