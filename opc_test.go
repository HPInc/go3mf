// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package go3mf

import (
	"reflect"
	"testing"

	"github.com/qmuntal/opc"
)

func Test_newRelationships(t *testing.T) {
	type args struct {
		rels []*opc.Relationship
	}
	tests := []struct {
		name string
		args args
		want []Relationship
	}{
		{"base", args{[]*opc.Relationship{{}, {TargetURI: "a.xml"}}}, []Relationship{{}, {Path: "a.xml"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRelationships(tt.args.rels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRelationships() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opcFile_Name(t *testing.T) {
	tests := []struct {
		name string
		o    *opcFile
		want string
	}{
		{"empty", &opcFile{nil, &opc.File{Part: new(opc.Part)}}, ""},
		{"base", &opcFile{nil, &opc.File{Part: &opc.Part{Name: "a.xml"}}}, "a.xml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Name(); got != tt.want {
				t.Errorf("opcFile.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opcFile_Relationships(t *testing.T) {
	tests := []struct {
		name string
		o    *opcFile
		want []Relationship
	}{
		{"empty", &opcFile{nil, &opc.File{Part: new(opc.Part)}}, []Relationship{}},
		{"base", &opcFile{nil, &opc.File{Part: &opc.Part{Relationships: []*opc.Relationship{
			{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture", TargetURI: "/a.xml"},
			{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel", TargetURI: "/b.xml"},
		}}}}, []Relationship{
			{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture", Path: "/a.xml"},
			{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel", Path: "/b.xml"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Relationships(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("opcFile.Relationships() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opcReader_FindFileFromName(t *testing.T) {
	reader := &opc.Reader{Files: []*opc.File{{Part: &opc.Part{Name: "/a.xml"}}, {Part: &opc.Part{Name: "/b.xml"}}}}
	type args struct {
		name string
	}
	tests := []struct {
		name string
		o    *opcReader
		args args
		want packageFile
	}{
		{"foundA", &opcReader{nil, 0, reader}, args{"/a.xml"}, &opcFile{reader, &opc.File{Part: &opc.Part{Name: "/a.xml"}}}},
		{"foundB", &opcReader{nil, 0, reader}, args{"/b.xml"}, &opcFile{reader, &opc.File{Part: &opc.Part{Name: "/b.xml"}}}},
		{"notfound", &opcReader{nil, 0, reader}, args{"/c.xml"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.o.FindFileFromName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("opcReader.FindFileFromName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opcFile_ContentType(t *testing.T) {
	tests := []struct {
		name string
		o    *opcFile
		want string
	}{
		{"base", &opcFile{f: &opc.File{Part: &opc.Part{ContentType: "fake_type"}}}, "fake_type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.ContentType(); got != tt.want {
				t.Errorf("opcFile.ContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opcWriter_AddRelationship(t *testing.T) {
	type args struct {
		r Relationship
	}
	tests := []struct {
		name string
		o    *opcWriter
		args args
		want []*opc.Relationship
	}{
		{"base", newOpcWriter(nil), args{Relationship{ID: "id_1", Path: "fake_uri", Type: "fake_type"}}, []*opc.Relationship{
			{ID: "id_1", TargetURI: "fake_uri", Type: "fake_type"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.AddRelationship(tt.args.r)
			if !reflect.DeepEqual(tt.o.w.Relationships, tt.want) {
				t.Errorf("opcWriter.AddRelationship() = %v, want %v", tt.o.w.Relationships, tt.want)
			}
		})
	}
}
