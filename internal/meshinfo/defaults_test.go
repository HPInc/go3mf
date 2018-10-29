package meshinfo

import (
	"reflect"
	"testing"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name string
		want Handler
	}{
		{"base", newlookupHandler()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNodeColorInfo(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{"base", args{1}, reflect.TypeOf((*NodeColor)(nil)).Elem()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNodeColorInfo(tt.args.currentFaceCount).InfoType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNodeColorInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTextureCoordsInfo(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{"base", args{1}, reflect.TypeOf((*TextureCoords)(nil)).Elem()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextureCoordsInfo(tt.args.currentFaceCount).InfoType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextureCoordsInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBaseMaterialInfo(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{"base", args{1}, reflect.TypeOf((*BaseMaterial)(nil)).Elem()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBaseMaterialInfo(tt.args.currentFaceCount).InfoType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBaseMaterialInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newInfo(t *testing.T) {
	type args struct {
		currentFaceCount uint32
		infoType         reflect.Type
	}
	tests := []struct {
		name string
		args args
		want MeshInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newInfo(tt.args.currentFaceCount, tt.args.infoType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
