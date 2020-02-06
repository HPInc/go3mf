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
		want []*relationship
	}{
		{"base", args{[]*opc.Relationship{{}, {TargetURI: "a.xml"}}}, []*relationship{{}, {TargetURI: "a.xml"}}},
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

func Test_opcFile_FindFileFromRel(t *testing.T) {
	rels := []*opc.Relationship{
		{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture", TargetURI: "/c.xml"},
		{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel", TargetURI: "/b.xml"},
		{Type: "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail", TargetURI: "Metadata/thumbnail.png"},
	}
	reader := &opc.Reader{
		Files: []*opc.File{{Part: &opc.Part{Name: "/c.xml"}}, {Part: &opc.Part{Name: "/props/Metadata/thumbnail.png"}}},
	}
	type args struct {
		relType string
	}
	tests := []struct {
		name string
		o    *opcFile
		args args
		want packageFile
	}{
		{"foundC", &opcFile{reader, &opc.File{Part: &opc.Part{Relationships: rels}}}, args{"http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture"}, &opcFile{reader, &opc.File{Part: &opc.Part{Name: "/c.xml"}}}},
		{"thumbnail", &opcFile{reader, &opc.File{Part: &opc.Part{Name: "/props/a.xml", Relationships: rels}}}, args{"http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"}, &opcFile{reader, &opc.File{Part: &opc.Part{Name: "/props/Metadata/thumbnail.png"}}}},
		{"noFile", &opcFile{reader, &opc.File{Part: &opc.Part{Relationships: rels}}}, args{"http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel"}, nil},
		{"noRel", &opcFile{reader, &opc.File{Part: &opc.Part{Relationships: rels}}}, args{"other"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.o.FindFileFromRel(tt.args.relType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("opcFile.FindFileFromRel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opcFile_Relationships(t *testing.T) {
	tests := []struct {
		name string
		o    *opcFile
		want []*relationship
	}{
		{"empty", &opcFile{nil, &opc.File{Part: new(opc.Part)}}, []*relationship{}},
		{"base", &opcFile{nil, &opc.File{Part: &opc.Part{Relationships: []*opc.Relationship{
			{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture", TargetURI: "/a.xml"},
			{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel", TargetURI: "/b.xml"},
		}}}}, []*relationship{
			{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture", TargetURI: "/a.xml"},
			{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel", TargetURI: "/b.xml"},
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

func Test_opcReader_FindFileFromRel(t *testing.T) {
	reader := &opc.Reader{
		Relationships: []*opc.Relationship{
			{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture", TargetURI: "/a.xml"},
			{Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel", TargetURI: "/b.xml"},
		},
		Files: []*opc.File{{Part: &opc.Part{Name: "/a.xml"}}},
	}
	type args struct {
		relType string
	}
	tests := []struct {
		name string
		o    *opcReader
		args args
		want packageFile
	}{
		{"foundA", &opcReader{nil, 0, reader}, args{"http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture"}, &opcFile{reader, &opc.File{Part: &opc.Part{Name: "/a.xml"}}}},
		{"noFile", &opcReader{nil, 0, reader}, args{"http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel"}, nil},
		{"noRel", &opcReader{nil, 0, reader}, args{"other"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.o.FindFileFromRel(tt.args.relType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("opcReader.FindFileFromRel() = %v, want %v", got, tt.want)
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
