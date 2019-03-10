package model

import (
	"io"
	"reflect"
	"testing"

	"github.com/qmuntal/go3mf/internal/mesh"
)

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
