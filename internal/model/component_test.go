package model

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/gofrs/uuid"
	"github.com/qmuntal/go3mf/internal/mesh"
	"github.com/stretchr/testify/mock"
)

// MockMergeableMesh is a mock of MergeableMesh interface
type MockObject struct {
	mock.Mock
}

func NewMockObject(isValid, isValidForSlices bool) *MockObject {
	o := new(MockObject)
	o.On("IsValid").Return(isValid)
	o.On("IsValidForSlices", mock.Anything).Return(isValidForSlices)
	return o
}

func (o *MockObject) RootModel() *Model {
	o.Called()
	return new(Model)
}
func (o *MockObject) MergeToMesh(args0 *mesh.Mesh, args1 mgl32.Mat4) {
	o.Called(args0, args1)
	return
}
func (o *MockObject) ID() uint64 {
	o.Called()
	return 0
}
func (o *MockObject) IsValid() bool {
	args := o.Called()
	return args.Bool(0)
}

func (o *MockObject) IsValidForSlices(args0 mgl32.Mat4) bool {
	args := o.Called(args0)
	return args.Bool(0)
}

func newObject() *ObjectResource {
	o, _ := newObjectResource(0, new(Model))
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
		{"base", &Component{Object: newObject()}, args{new(mesh.Mesh), mgl32.Ident4()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.MergeToMesh(tt.args.m, tt.args.transform)
		})
	}
}

func TestObjectResource_IsValid(t *testing.T) {
	tests := []struct {
		name string
		o    *ObjectResource
		want bool
	}{
		{"base", new(ObjectResource), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.IsValid(); got != tt.want {
				t.Errorf("ObjectResource.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponentResource_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    *ComponentResource
		want bool
	}{
		{"empty", new(ComponentResource), false},
		{"oneInvalid", &ComponentResource{Components: []*Component{{Object: NewMockObject(true, true)}, {Object: NewMockObject(false, true)}}}, false},
		{"valid", &ComponentResource{Components: []*Component{{Object: NewMockObject(true, true)}, {Object: NewMockObject(true, true)}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("ComponentResource.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectResource_IsValidForSlices(t *testing.T) {
	type args struct {
		transform mgl32.Mat4
	}
	tests := []struct {
		name string
		o    *ObjectResource
		args args
		want bool
	}{
		{"base", new(ObjectResource), args{mgl32.Ident4()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.IsValidForSlices(tt.args.transform); got != tt.want {
				t.Errorf("ObjectResource.IsValidForSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponentResource_IsValidForSlices(t *testing.T) {
	type args struct {
		transform mgl32.Mat4
	}
	tests := []struct {
		name string
		c    *ComponentResource
		args args
		want bool
	}{
		{"empty", new(ComponentResource), args{mgl32.Ident4()}, true},
		{"oneInvalid", &ComponentResource{Components: []*Component{{Object: NewMockObject(true, true)}, {Object: NewMockObject(true, false)}}}, args{mgl32.Ident4()}, false},
		{"valid", &ComponentResource{Components: []*Component{{Object: NewMockObject(true, true)}, {Object: NewMockObject(true, true)}}}, args{mgl32.Ident4()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValidForSlices(tt.args.transform); got != tt.want {
				t.Errorf("ComponentResource.IsValidForSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectResource_ID(t *testing.T) {
	tests := []struct {
		name string
		o    *ObjectResource
		want uint64
	}{
		{"base", &ObjectResource{Resource: Resource{ResourceID: &PackageResourceID{uniqueID: 1}}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.ID(); got != tt.want {
				t.Errorf("ObjectResource.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectResource_MergeToMesh(t *testing.T) {
	type args struct {
		m         *mesh.Mesh
		transform mgl32.Mat4
	}
	tests := []struct {
		name string
		o    *ObjectResource
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.MergeToMesh(tt.args.m, tt.args.transform)
		})
	}
}

func TestComponentResource_MergeToMesh(t *testing.T) {
	type args struct {
		m         *mesh.Mesh
		transform mgl32.Mat4
	}
	tests := []struct {
		name string
		c    *ComponentResource
		args args
	}{
		{"base", new(ComponentResource), args{nil, mgl32.Ident4()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.MergeToMesh(tt.args.m, tt.args.transform)
		})
	}
}

func TestNewComponentResource(t *testing.T) {
	model := new(Model)
	type args struct {
		id    uint64
		model *Model
	}
	tests := []struct {
		name    string
		args    args
		want    *ComponentResource
		wantErr bool
	}{
		{"base", args{0, model}, &ComponentResource{ObjectResource: ObjectResource{Resource: Resource{Model: model, ResourceID: &PackageResourceID{"", 0, 1}}}}, false},
		{"dup", args{0, model}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewComponentResource(tt.args.id, tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewComponentResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewComponentResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeshResource_IsValidForSlices(t *testing.T) {
	type args struct {
		t mgl32.Mat4
	}
	tests := []struct {
		name string
		c    *MeshResource
		args args
		want bool
	}{
		{"empty", new(MeshResource), args{mgl32.Mat4{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}, true},
		{"valid", &MeshResource{ObjectResource: ObjectResource{SliceStackID: &PackageResourceID{}}}, args{mgl32.Mat4{1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1}}, true},
		{"invalid", &MeshResource{ObjectResource: ObjectResource{SliceStackID: &PackageResourceID{}}}, args{mgl32.Mat4{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValidForSlices(tt.args.t); got != tt.want {
				t.Errorf("MeshResource.IsValidForSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMeshResource(t *testing.T) {
	model := new(Model)
	type args struct {
		id    uint64
		model *Model
	}
	tests := []struct {
		name    string
		args    args
		want    *MeshResource
		wantErr bool
	}{
		{"base", args{0, model}, &MeshResource{ObjectResource: ObjectResource{Resource: Resource{Model: model, ResourceID: &PackageResourceID{"", 0, 1}}}}, false},
		{"dup", args{0, model}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMeshResource(tt.args.id, tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMeshResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMeshResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeshResource_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    *MeshResource
		want bool
	}{
		{"empty", new(MeshResource), false},
		{"other", &MeshResource{Mesh: new(mesh.Mesh), ObjectResource: ObjectResource{Type: OtherType}}, false},
		{"surface", &MeshResource{Mesh: new(mesh.Mesh), ObjectResource: ObjectResource{Type: SurfaceType}}, true},
		{"support", &MeshResource{Mesh: new(mesh.Mesh), ObjectResource: ObjectResource{Type: SupportType}}, true},
		{"solidsupport", &MeshResource{Mesh: new(mesh.Mesh), ObjectResource: ObjectResource{Type: SolidSupportType}}, false},
		{"model", &MeshResource{Mesh: new(mesh.Mesh), ObjectResource: ObjectResource{Type: ModelType}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("MeshResource.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
