package meshinfo

import (
	"reflect"
	"testing"
)

func TestMemoryMeshInfoFactory_Create(t *testing.T) {
	type args struct {
		infoType         reflect.Type
		currentFaceCount uint32
	}
	tests := []struct {
		name    string
		f       *MemoryMeshInfoFactory
		args    args
		want    MeshInfo
		wantErr bool
	}{
		{"basematerial", new(MemoryMeshInfoFactory), args{reflect.TypeOf((*BaseMaterial)(nil)).Elem(), 3}, newgenericMeshInfo(newmemoryContainer(3, reflect.TypeOf((*BaseMaterial)(nil)).Elem())), false},
		{"basematerial", new(MemoryMeshInfoFactory), args{reflect.TypeOf((*NodeColor)(nil)).Elem(), 3}, newgenericMeshInfo(newmemoryContainer(3, reflect.TypeOf((*NodeColor)(nil)).Elem())), false},
		{"basematerial", new(MemoryMeshInfoFactory), args{reflect.TypeOf((*TextureCoords)(nil)).Elem(), 3}, newgenericMeshInfo(newmemoryContainer(3, reflect.TypeOf((*TextureCoords)(nil)).Elem())), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Create(tt.args.infoType, tt.args.currentFaceCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryMeshInfoFactory.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.InfoType(), tt.want.InfoType()) {
				t.Errorf("MemoryMeshInfoFactory.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
