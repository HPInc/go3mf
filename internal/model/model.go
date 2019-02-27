package model

import (
	"errors"
	"io"

	"github.com/gofrs/uuid"
	"github.com/qmuntal/go3mf/internal/mesh"
)

// Identifier defines an object than can be uniquely identified.
type Identifier interface {
	UniqueID() uint64
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
	Units                 Units
	Thumbnail             *Attachment
	Metadata              []Metadata
	Resources             []Identifier
	BuildItems            []*BuildItem
	Attachments           []*Attachment
	ProductionAttachments []*Attachment
	CustomContentTypes    map[string]string

	usedUUIDs       map[uuid.UUID]struct{}
	resourceHandler resourceHandler
	resourceMap     map[uint64]Identifier
	uuid            uuid.UUID
	objects         []interface{}
	baseMaterials   []*BaseMaterialsResource
	textures        []*Texture2DResource
	sliceStacks     []*SliceStackResource
}

// NewModel returns a new initialized model.
func NewModel() *Model {
	m := &Model{
		Units:              Millimeter,
		Language:           langUS,
		CustomContentTypes: make(map[string]string),
		usedUUIDs:          make(map[uuid.UUID]struct{}),
		resourceMap:        make(map[uint64]Identifier),
	}
	m.SetUUID(uuid.Must(uuid.NewV4()))
	return m
}

// UUID returns the build UUID.
func (m *Model) UUID() uuid.UUID {
	return m.uuid
}

// SetUUID sets the build UUID
func (m *Model) SetUUID(id uuid.UUID) error {
	err := registerUUID(m.uuid, id, m)
	if err == nil {
		m.uuid = id
	}
	return err
}

// SetThumbnail sets the package thumbnail.
func (m *Model) SetThumbnail(r io.Reader) *Attachment {
	m.Thumbnail = &Attachment{Stream: r, Path: thumbnailPath, RelationshipType: relTypeThumbnail}
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

// FindResourcePath returns the resource with the target path and ID.
func (m *Model) FindResourcePath(path string, id uint64) (r Identifier, ok bool) {
	rID, ok := m.FindPackageResourcePath(path, id)
	if ok {
		return m.FindResource(rID.UniqueID())
	}
	return nil, false
}

// FindResource returns the resource with the target unique ID.
func (m *Model) FindResource(uniqueID uint64) (i Identifier, ok bool) {
	i, ok = m.resourceMap[uniqueID]
	return
}

// FindPackageResourceID returns the package resource with the target unique ID.
func (m *Model) FindPackageResourceID(uniqueID uint64) (*ResourceID, bool) {
	return m.resourceHandler.FindResourceID(uniqueID)
}

// FindPackageResourcePath returns the package resource with the target path and ID.
func (m *Model) FindPackageResourcePath(path string, id uint64) (*ResourceID, bool) {
	return m.resourceHandler.FindResourcePath(path, id)
}

// AddResource adds a new resource to the model.
func (m *Model) AddResource(resource Identifier) error {
	id := resource.UniqueID()
	if _, ok := m.FindResource(id); ok {
		return errors.New("go3mf: duplicated model resource")
	}

	m.resourceMap[id] = resource
	m.Resources = append(m.Resources, resource)
	m.addResourceToLookupTable(resource)
	return nil
}

func (m *Model) addResourceToLookupTable(resource Identifier) {
	switch resource.(type) {
	case *ComponentResource:
		m.objects = append(m.objects, resource)
	case *MeshResource:
		m.objects = append(m.objects, resource)
	case *BaseMaterialsResource:
		m.baseMaterials = append(m.baseMaterials, resource.(*BaseMaterialsResource))
	case *Texture2DResource:
		m.textures = append(m.textures, resource.(*Texture2DResource))
	case *SliceStackResource:
		m.sliceStacks = append(m.sliceStacks, resource.(*SliceStackResource))
	}
}

func (m *Model) generatePackageResourceID(id uint64) (*ResourceID, error) {
	return m.resourceHandler.NewResourceID(m.Path, id)
}

func (m *Model) unregisterUUID(id uuid.UUID) {
	delete(m.usedUUIDs, id)
}

func (m *Model) registerUUID(id uuid.UUID) error {
	if _, ok := m.usedUUIDs[id]; ok {
		return errors.New("go3mf: duplicated UUID")
	}
	if len(m.usedUUIDs) == 0 {
		m.usedUUIDs = make(map[uuid.UUID]struct{})
	}
	var tmp struct{}
	m.usedUUIDs[id] = tmp
	return nil
}
