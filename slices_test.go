package go3mf

import (
	"testing"

	"github.com/qmuntal/go3mf/mesh"
)

func TestBuildItem_IsValidForSlices(t *testing.T) {
	tests := []struct {
		name string
		b    *BuildItem
		want bool
	}{
		{"valid", &BuildItem{Object: NewMockObject(true, true)}, true},
		{"valid", &BuildItem{Object: NewMockObject(true, false)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsValidForSlices(); got != tt.want {
				t.Errorf("BuildItem.IsValidForSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponentsResource_IsValidForSlices(t *testing.T) {
	type args struct {
		transform mesh.Matrix
	}
	tests := []struct {
		name string
		c    *ComponentsResource
		args args
		want bool
	}{
		{"empty", new(ComponentsResource), args{mesh.Identity()}, true},
		{"oneInvalid", &ComponentsResource{Components: []*Component{{Object: NewMockObject(true, true)}, {Object: NewMockObject(true, false)}}}, args{mesh.Identity()}, false},
		{"valid", &ComponentsResource{Components: []*Component{{Object: NewMockObject(true, true)}, {Object: NewMockObject(true, true)}}}, args{mesh.Identity()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValidForSlices(tt.args.transform); got != tt.want {
				t.Errorf("ComponentsResource.IsValidForSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeshResource_IsValidForSlices(t *testing.T) {
	type args struct {
		t mesh.Matrix
	}
	tests := []struct {
		name string
		c    *MeshResource
		args args
		want bool
	}{
		{"empty", new(MeshResource), args{mesh.Matrix{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}, true},
		{"valid", &MeshResource{ObjectResource: ObjectResource{SliceStackID: 0}}, args{mesh.Matrix{1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1}}, true},
		{"invalid", &MeshResource{ObjectResource: ObjectResource{SliceStackID: 1}}, args{mesh.Matrix{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValidForSlices(tt.args.t); got != tt.want {
				t.Errorf("MeshResource.IsValidForSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceStack_AddSlice(t *testing.T) {
	type args struct {
		slice *mesh.Slice
	}
	tests := []struct {
		name    string
		s       *SliceStack
		args    args
		want    int
		wantErr bool
	}{
		{"lower", &SliceStack{BottomZ: 1}, args{&mesh.Slice{TopZ: 0.5}}, 0, true},
		{"top", &SliceStack{Slices: []*mesh.Slice{{TopZ: 1.0}}}, args{&mesh.Slice{TopZ: 0.5}}, 0, true},
		{"ok", &SliceStack{BottomZ: 1, Slices: []*mesh.Slice{{TopZ: 1.0}}}, args{&mesh.Slice{TopZ: 2.0}}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AddSlice(tt.args.slice)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceStack.AddSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SliceStack.AddSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceStackResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		s     *SliceStackResource
		want  string
		want1 uint32
	}{
		{"base", &SliceStackResource{ID: 1, ModelPath: "3d/3dmodel.model"}, "3d/3dmodel.model", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.Identify()
			if got != tt.want {
				t.Errorf("SliceStackResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SliceStackResource.Identify() got = %v, want %v", got1, tt.want1)
			}
		})
	}
}
