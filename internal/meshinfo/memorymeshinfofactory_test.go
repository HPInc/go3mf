package meshinfo

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

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
		{"error", new(MemoryMeshInfoFactory), args{InfoAbstract, 0}, nil, true},
		{"basematerials", new(MemoryMeshInfoFactory), args{InfoBaseMaterials, 0}, newgenericMeshInfo(mockContainer, InfoBaseMaterials), false},
		{"nodecolors", new(MemoryMeshInfoFactory), args{InfoNodeColors, 0}, newgenericMeshInfo(mockContainer, InfoNodeColors), false},
		{"textureCoords", new(MemoryMeshInfoFactory), args{InfoTextureCoords, 0}, newgenericMeshInfo(mockContainer, InfoTextureCoords), false},
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
