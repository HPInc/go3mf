package go3mf

import (
	"image/color"
	"io"
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/mesh"
	"github.com/stretchr/testify/mock"
)

// MockObject is a mock of Object interface
type MockObject struct {
	mock.Mock
}

func NewMockObject(isValid, isValidForSlices bool) *MockObject {
	o := new(MockObject)
	o.On("IsValid").Return(isValid)
	o.On("IsValidForSlices", mock.Anything).Return(isValidForSlices)
	return o
}

func (o *MockObject) Identify() (string, uint64) {
	args := o.Called()
	return args.String(0), uint64(args.Int(1))
}

func (o *MockObject) Type() ObjectType {
	return ObjectTypeOther
}

func (o *MockObject) MergeToMesh(args0 *mesh.Mesh, args1 mgl32.Mat4) {
	o.Called(args0, args1)
	return
}

func (o *MockObject) IsValid() bool {
	args := o.Called()
	return args.Bool(0)
}

func (o *MockObject) IsValidForSlices(args0 mgl32.Mat4) bool {
	args := o.Called(args0)
	return args.Bool(0)
}

func TestModel_SetThumbnail(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		m    *Model
		args args
		want *Attachment
	}{
		{"base", new(Model), args{nil}, &Attachment{Path: thumbnailPath, RelationshipType: "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.SetThumbnail(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.SetThumbnail() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestModel_MergeToMesh(t *testing.T) {
	type args struct {
		msh *mesh.Mesh
	}
	tests := []struct {
		name string
		m    *Model
		args args
	}{
		{"base", &Model{BuildItems: []*BuildItem{{Object: new(ComponentsResource)}}}, args{new(mesh.Mesh)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.MergeToMesh(tt.args.msh)
		})
	}
}

func TestModel_FindResource(t *testing.T) {
	model := &Model{Path: "/3D/model.model"}
	id1 := &ObjectResource{ID: 0, ModelPath: ""}
	id2 := &ObjectResource{ID: 1, ModelPath: "/3D/model.model"}
	model.Resources = append(model.Resources, id1, id2)
	type args struct {
		path string
		id   uint64
	}
	tests := []struct {
		name   string
		m      *Model
		args   args
		wantR  Resource
		wantOk bool
	}{
		{"emptyPathExist", model, args{"", 1}, id2, true},
		{"emptyPathNoExist", model, args{"", 0}, nil, false},
		{"exist2", model, args{"/3D/model.model", 1}, id2, true},
		{"noexist", model, args{"/3D/model.model", 100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotOk := tt.m.FindResource(tt.args.path, tt.args.id)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Model.FindResource() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Model.FindResource() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestBaseMaterial_ColotString(t *testing.T) {
	tests := []struct {
		name string
		m    *BaseMaterial
		want string
	}{
		{"base", &BaseMaterial{Color: color.RGBA{200, 250, 60, 80}}, "#c8fa3c50"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ColorString(); got != tt.want {
				t.Errorf("BaseMaterial.ColotString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsResource_Merge(t *testing.T) {
	type args struct {
		other []BaseMaterial
	}
	tests := []struct {
		name string
		ms   *BaseMaterialsResource
		args args
	}{
		{"base", &BaseMaterialsResource{Materials: []BaseMaterial{{Name: "1", Color: color.RGBA{200, 250, 60, 80}}}}, args{
			[]BaseMaterial{{Name: "2", Color: color.RGBA{200, 250, 60, 80}}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := append(tt.ms.Materials, tt.args.other...)
			tt.ms.Merge(tt.args.other)
			if !reflect.DeepEqual(tt.ms.Materials, want) {
				t.Errorf("BaseMaterialsResource.Merge() = %v, want %v", tt.ms.Materials, want)
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
		{"base", &BuildItem{Object: new(ComponentsResource)}, args{new(mesh.Mesh)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.MergeToMesh(tt.args.m)
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

func TestComponentsResource_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    *ComponentsResource
		want bool
	}{
		{"empty", new(ComponentsResource), false},
		{"oneInvalid", &ComponentsResource{Components: []*Component{{Object: NewMockObject(true, true)}, {Object: NewMockObject(false, true)}}}, false},
		{"valid", &ComponentsResource{Components: []*Component{{Object: NewMockObject(true, true)}, {Object: NewMockObject(true, true)}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("ComponentsResource.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponentsResource_IsValidForSlices(t *testing.T) {
	type args struct {
		transform mgl32.Mat4
	}
	tests := []struct {
		name string
		c    *ComponentsResource
		args args
		want bool
	}{
		{"empty", new(ComponentsResource), args{mgl32.Ident4()}, true},
		{"oneInvalid", &ComponentsResource{Components: []*Component{{Object: NewMockObject(true, true)}, {Object: NewMockObject(true, false)}}}, args{mgl32.Ident4()}, false},
		{"valid", &ComponentsResource{Components: []*Component{{Object: NewMockObject(true, true)}, {Object: NewMockObject(true, true)}}}, args{mgl32.Ident4()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValidForSlices(tt.args.transform); got != tt.want {
				t.Errorf("ComponentsResource.IsValidForSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponentsResource_MergeToMesh(t *testing.T) {
	type args struct {
		m         *mesh.Mesh
		transform mgl32.Mat4
	}
	tests := []struct {
		name string
		c    *ComponentsResource
		args args
	}{
		{"empty", new(ComponentsResource), args{nil, mgl32.Ident4()}},
		{"base", &ComponentsResource{Components: []*Component{{Object: new(ComponentsResource)}}}, args{nil, mgl32.Ident4()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.MergeToMesh(tt.args.m, tt.args.transform)
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
		{"valid", &MeshResource{ObjectResource: ObjectResource{SliceStackID: 0}}, args{mgl32.Mat4{1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1}}, true},
		{"invalid", &MeshResource{ObjectResource: ObjectResource{SliceStackID: 1}}, args{mgl32.Mat4{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValidForSlices(tt.args.t); got != tt.want {
				t.Errorf("MeshResource.IsValidForSlices() = %v, want %v", got, tt.want)
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
		{"other", &MeshResource{Mesh: new(mesh.Mesh), ObjectResource: ObjectResource{ObjectType: ObjectTypeOther}}, false},
		{"surface", &MeshResource{Mesh: new(mesh.Mesh), ObjectResource: ObjectResource{ObjectType: ObjectTypeSurface}}, true},
		{"support", &MeshResource{Mesh: new(mesh.Mesh), ObjectResource: ObjectResource{ObjectType: ObjectTypeSupport}}, true},
		{"solidsupport", &MeshResource{Mesh: new(mesh.Mesh), ObjectResource: ObjectResource{ObjectType: ObjectTypeSolidSupport}}, false},
		{"model", &MeshResource{Mesh: new(mesh.Mesh), ObjectResource: ObjectResource{ObjectType: ObjectTypeModel}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("MeshResource.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeshResource_MergeToMesh(t *testing.T) {
	type args struct {
		m         *mesh.Mesh
		transform mgl32.Mat4
	}
	tests := []struct {
		name string
		c    *MeshResource
		args args
	}{
		{"base", &MeshResource{Mesh: new(mesh.Mesh)}, args{new(mesh.Mesh), mgl32.Ident4()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.MergeToMesh(tt.args.m, tt.args.transform)
		})
	}
}

func TestObjectResource_Type(t *testing.T) {
	tests := []struct {
		name string
		o    *ObjectResource
		want ObjectType
	}{
		{"base", &ObjectResource{ObjectType: ObjectTypeModel}, ObjectTypeModel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectResource.Type() = %v, want %v", got, tt.want)
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

func TestTexture2DResource_Copy(t *testing.T) {
	type args struct {
		other *Texture2DResource
	}
	tests := []struct {
		name string
		t    *Texture2DResource
		args args
	}{
		{"equal", &Texture2DResource{Path: "/a.png", ContentType: PNGTexture}, args{&Texture2DResource{Path: "/a.png", ContentType: PNGTexture}}},
		{"diff", &Texture2DResource{Path: "/b.png", ContentType: PNGTexture}, args{&Texture2DResource{Path: "/a.png", ContentType: JPEGTexture}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Copy(tt.args.other)
			if tt.t.Path != tt.args.other.Path {
				t.Errorf("Texture2DResource.Copy() gotPath = %v, want %v", tt.t.Path, tt.args.other.Path)
			}
			if tt.t.ContentType != tt.args.other.ContentType {
				t.Errorf("Texture2DResource.Copy() gotContentType = %v, want %v", tt.t.ContentType, tt.args.other.ContentType)
			}
		})
	}
}

func TestBaseMaterialsResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		ms    *BaseMaterialsResource
		want  string
		want1 uint64
	}{
		{"base", &BaseMaterialsResource{ID: 1, ModelPath: "3d/3dmodel.model"}, "3d/3dmodel.model", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.ms.Identify()
			if got != tt.want {
				t.Errorf("BaseMaterialsResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BaseMaterialsResource.Identify() got = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestObjectResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		o     *ObjectResource
		want  string
		want1 uint64
	}{
		{"base", &ObjectResource{ID: 1, ModelPath: "3d/3dmodel.model"}, "3d/3dmodel.model", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.o.Identify()
			if got != tt.want {
				t.Errorf("ObjectResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ObjectResource.Identify() got = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSliceStackResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		s     *SliceStackResource
		want  string
		want1 uint64
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

func TestTexture2DResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		t     *Texture2DResource
		want  string
		want1 uint64
	}{
		{"base", &Texture2DResource{ID: 1, ModelPath: "3d/3dmodel.model"}, "3d/3dmodel.model", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.t.Identify()
			if got != tt.want {
				t.Errorf("Texture2DResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Texture2DResource.Identify() got = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_UnusedID(t *testing.T) {
	tests := []struct {
		name string
		m    *Model
		want uint64
	}{
		{"empty", new(Model), 1},
		{"one", &Model{Resources: []Resource{&ColorGroupResource{ID: 2}}}, 1},
		{"two", &Model{Resources: []Resource{&ColorGroupResource{ID: 1}}}, 2},
		{"sequence", &Model{Resources: []Resource{&ColorGroupResource{ID: 1}, &ColorGroupResource{ID: 2}}}, 3},
		{"sparce", &Model{Resources: []Resource{&ColorGroupResource{ID: 1}, &ColorGroupResource{ID: 3}}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.UnusedID(); got != tt.want {
				t.Errorf("Model.UnusedID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextureCoord_U(t *testing.T) {
	tests := []struct {
		name string
		t    TextureCoord
		want float32
	}{
		{"base", TextureCoord{1, 2}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.U(); got != tt.want {
				t.Errorf("TextureCoord.U() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextureCoord_V(t *testing.T) {
	tests := []struct {
		name string
		t    TextureCoord
		want float32
	}{
		{"base", TextureCoord{1, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.V(); got != tt.want {
				t.Errorf("TextureCoord.V() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTexture2DGroupResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		t     *Texture2DGroupResource
		want  string
		want1 uint64
	}{
		{"base", &Texture2DGroupResource{ID: 1, ModelPath: "3d/3dmodel"}, "3d/3dmodel", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.t.Identify()
			if got != tt.want {
				t.Errorf("Texture2DGroupResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Texture2DGroupResource.Identify() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestColorGroupResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		c     *ColorGroupResource
		want  string
		want1 uint64
	}{
		{"base", &ColorGroupResource{ID: 1, ModelPath: "3d/3dmodel"}, "3d/3dmodel", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Identify()
			if got != tt.want {
				t.Errorf("ColorGroupResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ColorGroupResource.Identify() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
