package go3mf

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockObject is a mock of Object interface
type MockObject struct {
	mock.Mock
}

func NewMockObject() *MockObject {
	o := new(MockObject)
	return o
}

func (o *MockObject) Identify() (string, uint32) {
	args := o.Called()
	return args.String(0), uint32(args.Int(1))
}

func (o *MockObject) Type() ObjectType {
	return ObjectTypeOther
}

func TestModel_FindObject(t *testing.T) {
	model := &Model{Path: "/3D/model.model"}
	id1 := &ObjectResource{ID: 0, ModelPath: ""}
	id2 := &ObjectResource{ID: 1, ModelPath: "/3D/model.model"}
	model.Resources = append(model.Resources,
		&Resources{Path: "", Objects: []*ObjectResource{id1}},
		&Resources{Path: "/3D/model.model", Objects: []*ObjectResource{id2}},
	)
	type args struct {
		path string
		id   uint32
	}
	tests := []struct {
		name   string
		m      *Model
		args   args
		wantR  *ObjectResource
		wantOk bool
	}{
		{"emptyPathExist", model, args{"", 1}, id2, true},
		{"emptyPathNoExist", model, args{"", 0}, nil, false},
		{"exist2", model, args{"/3D/model.model", 1}, id2, true},
		{"noexist", model, args{"/3D/model.model", 100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotOk := tt.m.FindObject(tt.args.path, tt.args.id)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Model.FindObject() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Model.FindObject() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestBuildItem_HasTransform(t *testing.T) {
	tests := []struct {
		name string
		b    *Item
		want bool
	}{
		{"zero", &Item{}, false},
		{"identity", &Item{Transform: Identity()}, false},
		{"base", &Item{Transform: Matrix{2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.HasTransform(); got != tt.want {
				t.Errorf("Item.HasTransform() = %v, want %v", got, tt.want)
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
		{"zero", &Component{}, false},
		{"identity", &Component{Transform: Identity()}, false},
		{"base", &Component{Transform: Matrix{2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.HasTransform(); got != tt.want {
				t.Errorf("Component.HasTransform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectResource_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    *ObjectResource
		want bool
	}{
		{"empty", new(ObjectResource), false},
		{"both", &ObjectResource{Mesh: new(Mesh), Components: make([]*Component, 0)}, false},
		{"other", &ObjectResource{Mesh: new(Mesh), ObjectType: ObjectTypeOther}, false},
		//{"surface", &Mesh{ObjectResource: ObjectResource{ObjectType: ObjectTypeSurface}}, true},
		//{"support", &Mesh{ObjectResource: ObjectResource{ObjectType: ObjectTypeSupport}}, true},
		{"solidsupport", &ObjectResource{Mesh: new(Mesh), ObjectType: ObjectTypeSolidSupport}, false},
		{"model", &ObjectResource{Mesh: new(Mesh), ObjectType: ObjectTypeModel}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("ObjectResource.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		ms    *BaseMaterialsResource
		want  string
		want1 uint32
	}{
		{"base", &BaseMaterialsResource{ID: 1, ModelPath: "/3D/3dmodel.model"}, "/3D/3dmodel.model", 1},
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

func TestResources_UnusedID(t *testing.T) {
	tests := []struct {
		name string
		m    *Resources
		want uint32
	}{
		{"empty", new(Resources), 1},
		{"one-asset", &Resources{Assets: []Resource{&BaseMaterialsResource{ID: 2}}}, 1},
		{"one-object", &Resources{Objects: []*ObjectResource{{ID: 2}}}, 1},
		{"two", &Resources{Assets: []Resource{&BaseMaterialsResource{ID: 1}}}, 2},
		{"sequence", &Resources{Assets: []Resource{&BaseMaterialsResource{ID: 1}}, Objects: []*ObjectResource{{ID: 2}}}, 3},
		{"sparce", &Resources{Assets: []Resource{&BaseMaterialsResource{ID: 1}}, Objects: []*ObjectResource{{ID: 3}}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.UnusedID(); got != tt.want {
				t.Errorf("Resources.UnusedID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectType_String(t *testing.T) {
	tests := []struct {
		name string
		o    ObjectType
	}{
		{"model", ObjectTypeModel},
		{"other", ObjectTypeOther},
		{"support", ObjectTypeSupport},
		{"solidsupport", ObjectTypeSolidSupport},
		{"surface", ObjectTypeSurface},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.String(); got != tt.name {
				t.Errorf("ObjectType.String() = %v, want %v", got, tt.name)
			}
		})
	}
}
func TestUnits_String(t *testing.T) {
	tests := []struct {
		name string
		u    Units
	}{
		{"micron", UnitMicrometer},
		{"millimeter", UnitMillimeter},
		{"centimeter", UnitCentimeter},
		{"inch", UnitInch},
		{"foot", UnitFoot},
		{"meter", UnitMeter},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.name {
				t.Errorf("Units.String() = %v, want %v", got, tt.name)
			}
		})
	}
}

func TestMesh_CheckSanity(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want bool
	}{
		{"new", new(Mesh), true},
		{"facefail", &Mesh{Faces: make([]Face, 2)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.CheckSanity(); got != tt.want {
				t.Errorf("Mesh.CheckSanity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_IsManifoldAndOriented(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want bool
	}{
		{"valid", &Mesh{
			Nodes: []Point3D{{}, {}, {}, {}},
			Faces: []Face{
				{NodeIndices: [3]uint32{0, 1, 2}},
				{NodeIndices: [3]uint32{0, 3, 1}},
				{NodeIndices: [3]uint32{0, 2, 3}},
				{NodeIndices: [3]uint32{1, 3, 2}},
			},
		}, true},
		{"nonmanifold", &Mesh{
			Nodes: []Point3D{{}, {}, {}, {}},
			Faces: []Face{
				{NodeIndices: [3]uint32{0, 1, 2}},
				{NodeIndices: [3]uint32{0, 1, 3}},
				{NodeIndices: [3]uint32{0, 2, 3}},
				{NodeIndices: [3]uint32{1, 2, 3}},
			},
		}, false},
		{"empty", new(Mesh), false},
		{"2nodes", &Mesh{
			Nodes: make([]Point3D, 2),
			Faces: make([]Face, 3),
		}, false},
		{"2faces", &Mesh{
			Nodes: make([]Point3D, 3),
			Faces: make([]Face, 2),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsManifoldAndOriented(); got != tt.want {
				t.Errorf("Mesh.IsManifoldAndOriented() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeshBuilder_AddNode(t *testing.T) {
	pos := Point3D{1.0, 2.0, 3.0}
	existingStruct := NewMeshBuilder(new(Mesh))
	existingStruct.AddNode(pos)
	type args struct {
		position Point3D
	}
	tests := []struct {
		name string
		m    *MeshBuilder
		args args
		want uint32
	}{
		{"existing", existingStruct, args{pos}, 0},
		{"base", &MeshBuilder{Mesh: &Mesh{Nodes: []Point3D{{}}}, CalculateConnectivity: false}, args{pos}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddNode(tt.args.position)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MeshBuilder.AddNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_checkFacesSanity(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want bool
	}{
		{"i0==i1", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{1, 1, 2}}}}, false},
		{"i0==i2", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{1, 2, 1}}}}, false},
		{"i1==i2", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{2, 1, 1}}}}, false},
		{"i0big", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{3, 1, 2}}}}, false},
		{"i1big", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{0, 3, 2}}}}, false},
		{"i2big", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{0, 1, 3}}}}, false},
		{"good", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{0, 1, 2}}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.checkFacesSanity(); got != tt.want {
				t.Errorf("Mesh.checkFacesSanity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newObjectType(t *testing.T) {
	tests := []struct {
		name   string
		wantO  ObjectType
		wantOk bool
	}{
		{"model", ObjectTypeModel, true},
		{"other", ObjectTypeOther, true},
		{"support", ObjectTypeSupport, true},
		{"solidsupport", ObjectTypeSolidSupport, true},
		{"surface", ObjectTypeSurface, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, gotOk := newObjectType(tt.name)
			if !reflect.DeepEqual(gotO, tt.wantO) {
				t.Errorf("newObjectType() gotO = %v, want %v", gotO, tt.wantO)
			}
			if gotOk != tt.wantOk {
				t.Errorf("newObjectType() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_newUnits(t *testing.T) {
	tests := []struct {
		name  string
		want  Units
		want1 bool
	}{
		{"micron", UnitMicrometer, true},
		{"millimeter", UnitMillimeter, true},
		{"centimeter", UnitCentimeter, true},
		{"inch", UnitInch, true},
		{"foot", UnitFoot, true},
		{"meter", UnitMeter, true},
		{"", UnitMillimeter, false},
		{"other", UnitMillimeter, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := newUnits(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUnits() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("newUnits() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewMeshResource(t *testing.T) {
	tests := []struct {
		name string
		want *ObjectResource
	}{
		{"base", &ObjectResource{Mesh: new(Mesh)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMeshResource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMeshResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewComponentsResource(t *testing.T) {
	tests := []struct {
		name string
		want *ObjectResource
	}{
		{"base", &ObjectResource{Components: make([]*Component, 0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewComponentsResource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewComponentsResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
