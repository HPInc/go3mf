package go3mf

import (
	"reflect"
	"testing"

	"github.com/qmuntal/go3mf/spec/xml"
)

var _ xml.Marshaler = new(BaseMaterials)

func TestResources_FindAsset(t *testing.T) {
	id1 := &BaseMaterials{ID: 0}
	id2 := &BaseMaterials{ID: 1}
	type args struct {
		id uint32
	}
	tests := []struct {
		name  string
		rs    *Resources
		args  args
		want  Asset
		want1 bool
	}{
		{"exist1", &Resources{Assets: []Asset{id1, id2}}, args{0}, id1, true},
		{"exist2", &Resources{Assets: []Asset{id1, id2}}, args{1}, id2, true},
		{"noexistID", &Resources{Assets: []Asset{id1, id2}}, args{100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.rs.FindAsset(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resources.FindAsset() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Resources.FindAsset() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_FindAsset(t *testing.T) {
	model := &Model{Path: "/3D/model.model"}
	id1 := &BaseMaterials{ID: 0}
	id2 := &BaseMaterials{ID: 1}
	id3 := &BaseMaterials{ID: 1}
	model.Resources = Resources{Assets: []Asset{id1, id2}}
	model.Childs = map[string]*ChildModel{
		"/3D/other.model": {Resources: Resources{Assets: []Asset{id3}}},
	}
	type args struct {
		path string
		id   uint32
	}
	tests := []struct {
		name   string
		m      *Model
		args   args
		wantR  Asset
		wantOk bool
	}{
		{"emptyPath1", model, args{"", 0}, id1, true},
		{"emptyPath2", model, args{"", 1}, id2, true},
		{"exist2", model, args{"/3D/model.model", 1}, id2, true},
		{"exist3", model, args{"/3D/other.model", 1}, id3, true},
		{"noexistID", model, args{"/3D/model.model", 100}, nil, false},
		{"noexistPath", model, args{"/3d.model", 1}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotOk := tt.m.FindAsset(tt.args.path, tt.args.id)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Model.FindAsset() gotR = %v, want %v", gotR, tt.wantR)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("Model.FindAsset() gotOk = %v, want %v", gotOk, tt.wantOk)
				return
			}
		})
	}
}

func TestModel_FindObject(t *testing.T) {
	model := &Model{Path: "/3D/model.model"}
	id1 := &Object{ID: 0}
	id2 := &Object{ID: 1}
	id3 := &Object{ID: 1}
	model.Resources = Resources{Objects: []*Object{id1, id2}}
	model.Childs = map[string]*ChildModel{
		"/3D/other.model": {Resources: Resources{Objects: []*Object{id3}}},
	}
	type args struct {
		path string
		id   uint32
	}
	tests := []struct {
		name   string
		m      *Model
		args   args
		wantR  *Object
		wantOk bool
	}{
		{"emptyPath1", model, args{"", 0}, id1, true},
		{"emptyPath2", model, args{"", 1}, id2, true},
		{"exist2", model, args{"/3D/model.model", 1}, id2, true},
		{"exist3", model, args{"/3D/other.model", 1}, id3, true},
		{"noexistID", model, args{"/3D/model.model", 100}, nil, false},
		{"noexistPath", model, args{"/3d.model", 1}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotOk := tt.m.FindObject(tt.args.path, tt.args.id)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Model.FindObject() gotR = %v, want %v", gotR, tt.wantR)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("Model.FindObject() gotOk = %v, want %v", gotOk, tt.wantOk)
				return
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

func TestBaseMaterials_Identify(t *testing.T) {
	tests := []struct {
		name string
		ms   *BaseMaterials
		want uint32
	}{
		{"base", &BaseMaterials{ID: 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ms.Identify()
			if got != tt.want {
				t.Errorf("BaseMaterials.Identify() got = %v, want %v", got, tt.want)
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
		{"one-asset", &Resources{Assets: []Asset{&BaseMaterials{ID: 2}}}, 1},
		{"one-object", &Resources{Objects: []*Object{{ID: 2}}}, 1},
		{"two", &Resources{Assets: []Asset{&BaseMaterials{ID: 1}}}, 2},
		{"sequence", &Resources{Assets: []Asset{&BaseMaterials{ID: 1}}, Objects: []*Object{{ID: 2}}}, 3},
		{"sparce", &Resources{Assets: []Asset{&BaseMaterials{ID: 1}}, Objects: []*Object{{ID: 3}}}, 2},
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

func TestMeshBuilder_AddVertex(t *testing.T) {
	pos := Point3D{1.0, 2.0, 3.0}
	existingStruct := NewMeshBuilder(new(Mesh))
	existingStruct.AddVertex(pos)
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
		{"base", &MeshBuilder{Mesh: &Mesh{Vertices: []Point3D{{}}}, CalculateConnectivity: false}, args{pos}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddVertex(tt.args.position)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MeshBuilder.AddVertex() = %v, want %v", got, tt.want)
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

func TestComponent_ObjectPath(t *testing.T) {
	type args struct {
		defaultPath string
	}
	tests := []struct {
		name string
		c    *Component
		args args
		want string
	}{
		{"emptyattr", &Component{}, args{"/other.model"}, "/other.model"},
		{"emptypath", &Component{AnyAttr: ExtensionsAttr{&fakeAttr{}}}, args{"/other.model"}, "/other.model"},
		{"emptyattr", &Component{AnyAttr: ExtensionsAttr{&fakeAttr{Value: "/3dmodel.model"}}}, args{"/other.model"}, "/3dmodel.model"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.ObjectPath(tt.args.defaultPath); got != tt.want {
				t.Errorf("Component.ObjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_ObjectPath(t *testing.T) {
	tests := []struct {
		name string
		b    *Item
		want string
	}{
		{"emptyattr", &Item{}, ""},
		{"emptypath", &Item{AnyAttr: ExtensionsAttr{&fakeAttr{}}}, ""},
		{"emptyattr", &Item{AnyAttr: ExtensionsAttr{&fakeAttr{Value: "/3dmodel.model"}}}, "/3dmodel.model"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.ObjectPath(); got != tt.want {
				t.Errorf("Item.ObjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_WalkAssets(t *testing.T) {
	tests := []struct {
		name       string
		m          *Model
		wantPath   []string
		wantAssets []Asset
	}{
		{"base", &Model{Childs: map[string]*ChildModel{
			"/other.model": {
				Resources: Resources{Assets: []Asset{&BaseMaterials{ID: 1}}},
			},
			"/a.model": {
				Resources: Resources{Assets: []Asset{&BaseMaterials{ID: 1}}},
			}}, Resources: Resources{Assets: []Asset{&BaseMaterials{ID: 2}}},
		}, []string{"/a.model", "/other.model", ""}, []Asset{&BaseMaterials{ID: 1}, &BaseMaterials{ID: 1}, &BaseMaterials{ID: 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotP []string
			var gotA []Asset
			tt.m.WalkAssets(func(path string, r Asset) error {
				gotP = append(gotP, path)
				gotA = append(gotA, r)
				return nil
			})
			if !reflect.DeepEqual(gotP, tt.wantPath) {
				t.Errorf("Model.WalkAssets() gotPaths = %v, wantPath %v", gotP, tt.wantPath)
			}
			if !reflect.DeepEqual(gotA, tt.wantAssets) {
				t.Errorf("Model.WalkAssets() gotAssets = %v, wantAsset %v", gotA, tt.wantAssets)
			}
		})
	}
}

func TestModel_WalkObjects(t *testing.T) {
	tests := []struct {
		name       string
		m          *Model
		wantPath   []string
		wantObject []*Object
	}{
		{"base", &Model{Childs: map[string]*ChildModel{
			"/other.model": {
				Resources: Resources{Objects: []*Object{{ID: 1}}},
			},
			"/a.model": {
				Resources: Resources{Objects: []*Object{{ID: 1}}},
			}}, Resources: Resources{Objects: []*Object{{ID: 2}}},
		}, []string{"/a.model", "/other.model", ""}, []*Object{{ID: 1}, {ID: 1}, {ID: 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotP []string
			var gotA []*Object
			tt.m.WalkObjects(func(path string, r *Object) error {
				gotP = append(gotP, path)
				gotA = append(gotA, r)
				return nil
			})
			if !reflect.DeepEqual(gotP, tt.wantPath) {
				t.Errorf("Model.WalkObjects() gotPaths = %v, wantPath %v", gotP, tt.wantPath)
			}
			if !reflect.DeepEqual(gotA, tt.wantObject) {
				t.Errorf("Model.WalkObjects() gotObjects = %v, wantObject %v", gotA, tt.wantObject)
			}
		})
	}
}

func TestMesh_BoundingBox(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want Box
	}{
		{"empty", new(Mesh), Box{}},
		{"base", &Mesh{Vertices: []Point3D{{1, 1, 1}, {2, 2, 2}, {-1, 0, 3}}}, Box{Min: Point3D{-1, 0, 1}, Max: Point3D{2, 2, 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.BoundingBox(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mesh.BoundingBox() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_BoundingBox(t *testing.T) {
	tests := []struct {
		name string
		m    *Model
		want Box
	}{
		{"empty", new(Model), Box{}},
		{"base", &Model{
			Build: Build{Items: []*Item{
				{ObjectID: 1, Transform: Identity()},
				{ObjectID: 2},
				{ObjectID: 3},
			}},
			Resources: Resources{Objects: []*Object{
				{ID: 1, Mesh: &Mesh{Vertices: []Point3D{
					{10, 20, 30},
				}}},
				{ID: 2, Components: []*Component{
					{ObjectID: 1, Transform: Identity().Translate(100, 100, 100)},
					{ObjectID: 10},
				}},
			}},
		}, Box{Min: Point3D{10, 20, 30}, Max: Point3D{110, 120, 130}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.BoundingBox(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.BoundingBox() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtensionsAttr_Get(t *testing.T) {
	tests := []struct {
		name   string
		e      ExtensionsAttr
		want   interface{}
		wantOK bool
	}{
		{"nil", nil, new(fakeAttr), false},
		{"empty", ExtensionsAttr{}, new(fakeAttr), false},
		{"non-exist", ExtensionsAttr{nil}, new(fakeAttr), false},
		{"exist", ExtensionsAttr{&fakeAttr{Value: "1"}}, &fakeAttr{Value: "1"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := new(fakeAttr)
			if got := tt.e.Get(&target); got != tt.wantOK {
				t.Errorf("ExtensionsAttr.Get() = %v, wantOK %v", got, tt.wantOK)
				return
			}
			if !reflect.DeepEqual(target, tt.want) {
				t.Errorf("ExtensionsAttr.Get() = %v, want %v", target, tt.want)
			}
		})
	}
}

func TestExtensions_Get(t *testing.T) {
	tests := []struct {
		name   string
		e      Extensions
		want   interface{}
		wantOK bool
	}{
		{"nil", nil, new(fakeAsset), false},
		{"empty", Extensions{}, new(fakeAsset), false},
		{"non-exist", Extensions{nil}, new(fakeAsset), false},
		{"exist", Extensions{&fakeAsset{ID: 1}}, &fakeAsset{ID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := new(fakeAsset)
			if got := tt.e.Get(&target); got != tt.wantOK {
				t.Errorf("Extensions.Get() = %v, wantOK %v", got, tt.wantOK)
				return
			}
			if !reflect.DeepEqual(target, tt.want) {
				t.Errorf("Extensions.Get() = %v, want %v", target, tt.want)
			}
		})
	}
}

func TestExtensions_Get_Panic(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name string
		e    Extensions
		args args
	}{
		{"nil", Extensions{&fakeAsset{ID: 1}}, args{nil}},
		{"int", Extensions{&fakeAsset{ID: 1}}, args{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Error("Extensions.Get() did not panic")
				}
			}()
			if tt.e.Get(tt.args.target) {
				t.Error("Extensions.Get() want false")
			}
		})
	}
}
