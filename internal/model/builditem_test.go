package model

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/gofrs/uuid"
	"github.com/qmuntal/go3mf/internal/mesh"
)

func TestBuildItem_UUID(t *testing.T) {
	tests := []struct {
		name string
		b    *BuildItem
		want uuid.UUID
	}{
		{"base", &BuildItem{uuid: uuid.UUID{}}, uuid.UUID{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.UUID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildItem.UUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildItem_SetUUID(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		b       *BuildItem
		args    args
		wantErr bool
	}{
		{"base", &BuildItem{Object: newObject()}, args{uuid.Must(uuid.NewV4())}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.SetUUID(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("BuildItem.SetUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildItem_HasTransform(t *testing.T) {
	tests := []struct {
		name string
		b    *BuildItem
		want bool
	}{
		{"identity", &BuildItem{Transform: mgl32.Ident4()}, false},
		{"base", &BuildItem{Transform: mgl32.Mat4{2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.HasTransform(); got != tt.want {
				t.Errorf("BuildItem.HasTransform() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestBuildItem_MergeToMesh(t *testing.T) {
	type args struct {
		m *mesh.Mesh
	}
	tests := []struct {
		name string
		b    *BuildItem
		args args
	}{
		{"base", &BuildItem{Object: newObject()}, args{new(mesh.Mesh)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.MergeToMesh(tt.args.m)
		})
	}
}
