package model

import (
	"io"
	"strings"
	"path/filepath"
	"github.com/qmuntal/opc"
)

type opcRelationship struct {
	rel *opc.Relationship
}

func (o *opcRelationship) Type() string {
	return o.rel.Type
}

func (o *opcRelationship) TargetURI() string {
	return o.rel.TargetURI
}

func newRelationships(rels []*opc.Relationship) []relationship {
	pr := make([]relationship, len(rels))
	for i, r := range rels {
		pr[i] = &opcRelationship{r}
	}
	return pr
}

type opcFile struct {
	r *opc.Reader
	f *opc.File
}

func (o *opcFile) Open() (io.ReadCloser, error) {
	return o.f.Open()
}

func (o *opcFile) Name() string {
	return o.f.Name
}

func (o *opcFile) FindFileFromRel(relType string) packageFile {
	name := findOPCFileURIFromRel(relType, o.f.Relationships)
	return o.FindFileFromName(name)
}


func (o *opcFile) FindFileFromName(name string) packageFile {
	if strings.HasPrefix(name, "/") || strings.HasPrefix(name, "\\") {
		name = filepath.Dir(o.f.Name) + name
	}
	return findOPCFileFromName(name, o.r)
}

func (o *opcFile) Relationships() []relationship {
	return newRelationships(o.f.Relationships)
}

type opcReader struct {
	r *opc.Reader
}

func newOPCReader(r io.ReaderAt, size int64) (*opcReader, error) {
	opcr, err := opc.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	return &opcReader{opcr}, nil
}

func (o *opcReader) FindFileFromRel(relType string) packageFile {
	name := findOPCFileURIFromRel(relType, o.r.Relationships)
	return o.FindFileFromName(name)
}

func (o *opcReader) FindFileFromName(name string) packageFile {
	return findOPCFileFromName(name, o.r)
}

func (o *opcReader) Relationships() []relationship {
	return newRelationships(o.r.Relationships)
}

func findOPCFileFromName(name string, r *opc.Reader) packageFile {
	for _, f := range r.Files {
		if f.Name == name {
			return &opcFile{r, f}
		}
	}
	return nil
}

func findOPCFileURIFromRel(relType string, rels []*opc.Relationship) string {
	for _, r := range rels {
		if r.Type == relType {
			return r.TargetURI
		}
	}
	return ""
}