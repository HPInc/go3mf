package model

import (
	"io"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/qmuntal/go3mf/internal/mesh"
)

func Test_uuidRegister_register(t *testing.T) {
	var a struct{}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		r       *uuidRegister
		args    args
		wantErr bool
	}{
		{"duplicated", &uuidRegister{usedUUIDs: map[uuid.UUID]struct{}{{}: a}}, args{uuid.UUID{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.register(uuid.UUID{}, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("uuidRegister.register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewModel(t *testing.T) {
	tests := []struct {
		name string
		want *Model
	}{
		{"base", &Model{
			Units:              UnitMillimeter,
			Language:           langUS,
			CustomContentTypes: make(map[string]string),
			resourceMap:        make(map[uint64]Identifier),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewModel()
			tt.want.SetUUID(got.UUID())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewModel() = %v, want %v", got, tt.want)
			}
		})
	}
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
		{"base", NewModel(), args{nil}, &Attachment{Path: thumbnailPath, RelationshipType: "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"}},
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

func TestModel_AddResource(t *testing.T) {
	m := NewModel()
	type args struct {
		resource Identifier
	}
	tests := []struct {
		name          string
		m             *Model
		args          args
		wantResources int
		wantErr       bool
	}{
		{"baseMaterial", m, args{&BaseMaterialsResource{uniqueID: 0}}, 1, false},
		{"sliceStack", m, args{&SliceStackResource{uniqueID: 1}}, 2, false},
		{"texture2D", m, args{&Texture2DResource{uniqueID: 2}}, 3, false},
		{"component", m, args{&ComponentResource{ObjectResource: ObjectResource{uniqueID: 3}}}, 4, false},
		{"mesh", m, args{&MeshResource{ObjectResource: ObjectResource{uniqueID: 4}}}, 5, false},
		{"duplicated", m, args{&MeshResource{ObjectResource: ObjectResource{uniqueID: 4}}}, 5, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.AddResource(tt.args.resource); (err != nil) != tt.wantErr {
				t.Errorf("Model.AddResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResources != len(tt.m.Resources) {
				t.Errorf("Model.AddResource() resource count error = %v, wantErr %v", len(tt.m.Resources), tt.wantResources)
				return
			}
		})
	}
}

func TestModel_MergeToMesh(t *testing.T) {
	type args struct {
		msh *mesh.Mesh
	}
	tests := []struct {
		name    string
		m       *Model
		args    args
		wantErr bool
	}{
		{"base", &Model{BuildItems: []*BuildItem{{Object: new(ObjectResource)}}}, args{new(mesh.Mesh)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.MergeToMesh(tt.args.msh); (err != nil) != tt.wantErr {
				t.Errorf("Model.MergeToMesh() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModel_generatePackageResourceID(t *testing.T) {
	model := &Model{Path: "/3D/model.model"}
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		m       *Model
		args    args
		want    *ResourceID
		wantErr bool
	}{
		{"new", new(Model), args{0}, &ResourceID{"", 0, 1}, false},
		{"path1", model, args{1}, &ResourceID{"/3D/model.model", 1, 1}, false},
		{"path2", model, args{3}, &ResourceID{"/3D/model.model", 3, 2}, false},
		{"error", model, args{1}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.generatePackageResourceID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Model.generatePackageResourceID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.generatePackageResourceID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_FindResource(t *testing.T) {
	model := new(Model)
	id1 := &ObjectResource{uniqueID: 0}
	id2 := &ObjectResource{uniqueID: 1}
	model.resourceMap = map[uint64]Identifier{id1.UniqueID(): id1, id2.UniqueID(): id2}
	type args struct {
		id uint64
	}
	tests := []struct {
		name   string
		m      *Model
		args   args
		wantR  Identifier
		wantOk bool
	}{
		{"exist1", model, args{id1.UniqueID()}, id1, true},
		{"exist2", model, args{id2.UniqueID()}, id2, true},
		{"noexist", model, args{100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotOk := tt.m.FindResource(tt.args.id)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Model.FindResource() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Model.FindResource() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestModel_FindResourcePath(t *testing.T) {
	model := &Model{Path: "/3D/model.model"}
	id1 := &ObjectResource{uniqueID: 0}
	id2 := &ObjectResource{uniqueID: 1}
	model.resourceMap = map[uint64]Identifier{id1.UniqueID(): id1, id2.UniqueID(): id2}
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
		{"exist1", model, args{"/3D/model.model", 0}, id1, true},
		{"exist2", model, args{"/3D/model.model", 1}, id2, true},
		{"noexistpath", model, args{"", 1}, nil, false},
		{"noexist", model, args{"/3D/model.model", 100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotOk := tt.m.FindResourcePath(tt.args.path, tt.args.id)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Model.FindResourcePath() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Model.FindResourcePath() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestModel_FindObject(t *testing.T) {
	model := new(Model)
	id1 := &ObjectResource{uniqueID: 10}
	id2 := &ObjectResource{uniqueID: 11}
	id3 := &Texture2DResource{uniqueID: 12}
	model.resourceMap = map[uint64]Identifier{id1.UniqueID(): id1, id2.UniqueID(): id2, id3.UniqueID(): id3}
	type args struct {
		uniqueID uint64
	}
	tests := []struct {
		name   string
		m      *Model
		args   args
		wantO  Object
		wantOk bool
	}{
		{"exist1", model, args{id1.UniqueID()}, id1, true},
		{"exist2", model, args{id2.UniqueID()}, id2, true},
		{"noobj", model, args{id3.UniqueID()}, nil, false},
		{"noexist", model, args{100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, gotOk := tt.m.FindObject(tt.args.uniqueID)
			if !reflect.DeepEqual(gotO, tt.wantO) {
				t.Errorf("Model.FindObject() gotO = %v, want %v", gotO, tt.wantO)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("Model.FindObject() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
