package model

import (
	"image/color"
	"io"
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/internal/mesh"
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

func mustIdentifier(a Identifier, err error) Identifier {
	return a
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
		{"base", &Model{BuildItems: []*BuildItem{{Object: new(ObjectResource)}}}, args{new(mesh.Mesh)}},
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
		wantR  Identifier
		wantOk bool
	}{
		{"exist1", model, args{"", 0}, id1, true},
		{"exist2", model, args{"/3D/model.model", 1}, id2, true},
		{"noexistpath", model, args{"", 1}, nil, false},
		{"noexist", model, args{"/3D/model.model", 100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotOk := tt.m.FindResource(tt.args.id, tt.args.path)
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
		{"base", &BuildItem{Object: new(ObjectResource)}, args{new(mesh.Mesh)}},
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
		{"base", &Component{Object: new(ObjectResource)}, args{new(mesh.Mesh), mgl32.Ident4()}},
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
		{"empty", new(ComponentResource), args{nil, mgl32.Ident4()}},
		{"base", &ComponentResource{Components: []*Component{{Object: new(ObjectResource)}}}, args{nil, mgl32.Ident4()}},
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

func TestSlice_BeginPolygon(t *testing.T) {
	s := new(Slice)
	tests := []struct {
		name string
		s    *Slice
		want int
	}{
		{"empty", s, 0},
		{"1", s, 1},
		{"2", s, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.BeginPolygon(); got != tt.want {
				t.Errorf("Slice.BeginPolygon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_AddVertex(t *testing.T) {
	s := new(Slice)
	type args struct {
		x float32
		y float32
	}
	tests := []struct {
		name string
		s    *Slice
		args args
		want int
	}{
		{"empty", s, args{1, 2}, 0},
		{"1", s, args{2, 3}, 1},
		{"2", s, args{4, 5}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.AddVertex(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Slice.AddVertex() = %v, want %v", got, tt.want)
				return
			}
			want := mgl32.Vec2{tt.args.x, tt.args.y}
			if !reflect.DeepEqual(tt.s.Vertices[tt.want], want) {
				t.Errorf("Slice.AddVertex() = %v, want %v", tt.s.Vertices[tt.want], want)
			}
		})
	}
}

func TestSlice_AddPolygonIndex(t *testing.T) {
	type args struct {
		polygonIndex int
		index        int
	}
	tests := []struct {
		name    string
		s       *Slice
		args    args
		wantErr bool
	}{
		{"emptyPolygon", new(Slice), args{0, 0}, true},
		{"emptyVertices", &Slice{Polygons: [][]int{{}}}, args{0, 0}, true},
		{"duplicated", &Slice{Polygons: [][]int{{0}}, Vertices: []mgl32.Vec2{{}}}, args{0, 0}, true},
		{"base", &Slice{Polygons: [][]int{{}}, Vertices: []mgl32.Vec2{{}}}, args{0, 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddPolygonIndex(tt.args.polygonIndex, tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("Slice.AddPolygonIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlice_AllPolygonsAreClosed(t *testing.T) {
	tests := []struct {
		name string
		s    *Slice
		want bool
	}{
		{"closed", &Slice{Polygons: [][]int{{0, 1, 0}}}, true},
		{"open", &Slice{Polygons: [][]int{{0, 1, 2}}}, false},
		{"one", &Slice{Polygons: [][]int{{0}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.AllPolygonsAreClosed(); got != tt.want {
				t.Errorf("Slice.AllPolygonsAreClosed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_IsPolygonValid(t *testing.T) {
	type args struct {
		index int
	}
	tests := []struct {
		name string
		s    *Slice
		args args
		want bool
	}{
		{"empty", new(Slice), args{0}, false},
		{"invalid1", &Slice{Polygons: [][]int{{0}}}, args{0}, false},
		{"invalid2", &Slice{Polygons: [][]int{{0, 1}}}, args{0}, false},
		{"valid", &Slice{Polygons: [][]int{{0, 1, 2}}}, args{0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsPolygonValid(tt.args.index); got != tt.want {
				t.Errorf("Slice.IsPolygonValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceStack_AddSlice(t *testing.T) {
	type args struct {
		slice *Slice
	}
	tests := []struct {
		name    string
		s       *SliceStack
		args    args
		want    int
		wantErr bool
	}{
		{"lower", &SliceStack{BottomZ: 1}, args{&Slice{TopZ: 0.5}}, 0, true},
		{"top", &SliceStack{Slices: []*Slice{{TopZ: 1.0}}}, args{&Slice{TopZ: 0.5}}, 0, true},
		{"ok", &SliceStack{BottomZ: 1, Slices: []*Slice{{TopZ: 1.0}}}, args{&Slice{TopZ: 2.0}}, 1, false},
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

func TestNewTexture2DResource(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name string
		args args
		want *Texture2DResource
	}{
		{"base", args{0}, &Texture2DResource{
			ContentType: PNGTexture,
			TileStyleU:  TileWrap,
			TileStyleV:  TileWrap,
			Filter:      TextureFilterAuto,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTexture2DResource(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTexture2DResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
