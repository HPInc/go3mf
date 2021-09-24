// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package go3mf

import (
	"io"

	"github.com/qmuntal/opc"
)

type opcPart struct {
	io.Writer
	Part *opc.Part
}

func (o *opcPart) AddRelationship(r Relationship) {
	for _, ro := range o.Part.Relationships {
		if ro.Type == r.Type && ro.TargetURI == r.Path {
			return
		}
	}
	o.Part.Relationships = append(o.Part.Relationships, &opc.Relationship{
		ID:        r.ID,
		Type:      r.Type,
		TargetURI: r.Path,
	})
}

type opcWriter struct {
	w *opc.Writer
}

func newOpcWriter(w io.Writer) *opcWriter {
	return &opcWriter{opc.NewWriter(w)}
}

func (o *opcWriter) Create(name, contentType string) (packagePart, error) {
	p := &opc.Part{Name: opc.NormalizePartName(name), ContentType: contentType}
	w, err := o.w.CreatePart(p, opc.CompressionNormal)
	if err != nil {
		return nil, err
	}
	return &opcPart{Writer: w, Part: p}, nil
}

func (o *opcWriter) AddRelationship(r Relationship) {
	for _, ro := range o.w.Relationships {
		if ro.Type == r.Type && ro.TargetURI == r.Path {
			return
		}
	}
	o.w.Relationships = append(o.w.Relationships, &opc.Relationship{
		ID:        r.ID,
		Type:      r.Type,
		TargetURI: r.Path,
	})
}

func (o *opcWriter) Close() error {
	return o.w.Close()
}

func newRelationships(rels []*opc.Relationship) []Relationship {
	pr := make([]Relationship, len(rels))
	for i, r := range rels {
		pr[i] = Relationship{ID: r.ID, Path: r.TargetURI, Type: r.Type}
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

func (o *opcFile) ContentType() string {
	return o.f.ContentType
}

func (o *opcFile) FindFileFromName(name string) (packageFile, bool) {
	name = opc.ResolveRelationship(o.f.Name, name)
	return findOPCFileFromName(name, o.r)
}

func (o *opcFile) Relationships() []Relationship {
	return newRelationships(o.f.Relationships)
}

type opcReader struct {
	ra   io.ReaderAt
	size int64
	r    *opc.Reader // nil until call Open.
}

func (o *opcReader) Open(f func(r io.Reader) io.ReadCloser) (err error) {
	o.r, err = opc.NewReader(o.ra, o.size)
	if f != nil {
		o.r.SetDecompressor(f)
	}
	return
}

func (o *opcReader) Relationships() []Relationship {
	return newRelationships(o.r.Relationships)
}

func (o *opcReader) FindFileFromName(name string) (packageFile, bool) {
	name = opc.ResolveRelationship("/", name)
	return findOPCFileFromName(name, o.r)
}

func resolveRelationship(source, rel string) string {
	return opc.ResolveRelationship(source, rel)
}

func findOPCFileFromName(name string, r *opc.Reader) (packageFile, bool) {
	for _, f := range r.Files {
		if f.Name == name {
			return &opcFile{r, f}, true
		}
	}
	return nil, false
}
