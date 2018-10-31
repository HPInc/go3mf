package meshinfo

import (
	"reflect"
	"testing"
)

func TestNewNodeColorFacesData(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{"base", args{1}, reflect.TypeOf((*NodeColor)(nil))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNodeColorFacesData(tt.args.currentFaceCount).InfoType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNodeColorFacesData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTextureCoordsFacesData(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{"base", args{1}, reflect.TypeOf((*TextureCoords)(nil))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextureCoordsFacesData(tt.args.currentFaceCount).InfoType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextureCoordsFacesData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBaseMaterialFacesData(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{"base", args{1}, reflect.TypeOf((*BaseMaterial)(nil))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBaseMaterialFacesData(tt.args.currentFaceCount).InfoType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBaseMaterialFacesData() = %v, want %v", got, tt.want)
			}
		})
	}
}
