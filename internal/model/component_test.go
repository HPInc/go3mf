package model

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/gofrs/uuid"
	"github.com/qmuntal/go3mf/internal/mesh"
)

func newObject() *ObjectResource {
	o, _ := NewObjectResource(0, new(Model))
	return o
}

func TestObjectResource_UUID(t *testing.T) {
	tests := []struct {
		name string
		o    *ObjectResource
		want uuid.UUID
	}{
		{"base", &ObjectResource{uuid: uuid.UUID{}}, uuid.UUID{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.UUID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectResource.UUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectResource_SetUUID(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name string
		o    *ObjectResource
		args args
	}{
		{"base", newObject(), args{uuid.Must(uuid.NewV4())}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.SetUUID(tt.args.id)
		})
	}
}

func TestComponent_UUID(t *testing.T) {
	tests := []struct {
		name string
		c    *Component
		want uuid.UUID
	}{
		{"base", &Component{uuid: uuid.UUID{}}, uuid.UUID{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.UUID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Component.UUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponent_SetUUID(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		c       *Component
		args    args
		wantErr bool
	}{
		{"base", &Component{Object: newObject()}, args{uuid.UUID{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetUUID(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Component.SetUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComponent_HasTransform(t *testing.T) {
	tests := []struct {
		name string
		c    *Component
		want bool
	}{
		{"identity", &Component{Transform: mgl32.Ident4()}, false},
		{"base", &Component{Transform: mgl32.Mat4{2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.HasTransform(); got != tt.want {
				t.Errorf("Component.HasTransform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponent_MergeToMesh(t *testing.T) {
	type args struct {
		m         *mesh.Mesh
		transform mgl32.Mat4
	}
	tests := []struct {
		name string
		c    *Component
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.MergeToMesh(tt.args.m, tt.args.transform)
		})
	}
}
