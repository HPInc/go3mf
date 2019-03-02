package io3mf

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

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

func (o *opcFile) FindFileFromRel(relType string) (packageFile, bool) {
	name := findOPCFileURIFromRel(relType, o.f.Relationships)
	if !strings.HasPrefix(name, "/") && !strings.HasPrefix(name, "\\") {
		base := strings.Replace(filepath.Dir(o.f.Name), "\\", "/", -1)
		name = fmt.Sprintf("%s/%s", base, name)
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

func (o *opcReader) FindFileFromRel(relType string) (packageFile, bool) {
	name := findOPCFileURIFromRel(relType, o.r.Relationships)
	return o.FindFileFromName(name)
}

func (o *opcReader) FindFileFromName(name string) (packageFile, bool) {
	return findOPCFileFromName(name, o.r)
}

func findOPCFileFromName(name string, r *opc.Reader) (packageFile, bool) {
	for _, f := range r.Files {
		if f.Name == name {
			return &opcFile{r, f}, true
		}
	}
	return nil, false
}

func findOPCFileURIFromRel(relType string, rels []*opc.Relationship) string {
	for _, r := range rels {
		if r.Type == relType {
			return r.TargetURI
		}
	}
	return ""
}
