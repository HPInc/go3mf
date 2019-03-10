package model

import (
	"io"

	"github.com/qmuntal/go3mf/internal/mesh"
)

// Identifier defines an object than can be uniquely identified.
type Identifier interface {
	Identify() (uint64, string)
}

// Metadata item is an in memory representation of the 3MF metadata,
// and can be attached to any 3MF model node.
type Metadata struct {
	Name  string
	Value string
}

// Attachment defines the Model Attachment.
type Attachment struct {
	Stream           io.Reader
	RelationshipType string
	Path             string
}

// A Model is an in memory representation of the 3MF file.
type Model struct {
	Path                  string
	RootPath              string
	Language              string
	UUID                  string
	Units                 Units
	Thumbnail             *Attachment
	Metadata              []Metadata
	Resources             []Identifier
	BuildItems            []*BuildItem
	Attachments           []*Attachment
	ProductionAttachments []*Attachment
}

// SetThumbnail sets the package thumbnail.
func (m *Model) SetThumbnail(r io.Reader) *Attachment {
	m.Thumbnail = &Attachment{Stream: r, Path: thumbnailPath, RelationshipType: "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"}
	return m.Thumbnail
}

// MergeToMesh merges the build with the mesh.
func (m *Model) MergeToMesh(msh *mesh.Mesh) error {
	for _, b := range m.BuildItems {
		if err := b.MergeToMesh(msh); err != nil {
			return err
		}
	}
	return nil
}

// FindResource returns the resource with the target unique ID.
func (m *Model) FindResource(id uint64, path string) (i Identifier, ok bool) {
	for _, value := range m.Resources {
		cid, cpath := value.Identify()
		if cid == id && cpath == path {
			i = value
			ok = true
			break
		}
	}
	return
}
