package model

import (
	"errors"
	"io"

	"github.com/gofrs/uuid"
	"github.com/qmuntal/go3mf/internal/mesh"
)

type register interface {
	register(old, new uuid.UUID) error
}

type uuidRegister struct {
	usedUUIDs map[uuid.UUID]struct{}
}

func (r *uuidRegister) register(oldID, newID uuid.UUID) error {
	if _, ok := r.usedUUIDs[newID]; ok {
		return errors.New("go3mf: duplicated UUID")
	}
	if len(r.usedUUIDs) == 0 {
		r.usedUUIDs = make(map[uuid.UUID]struct{})
	}
	delete(r.usedUUIDs, oldID)
	var tmp struct{}
	r.usedUUIDs[newID] = tmp
	return nil
}

// Identifier defines an object than can be uniquely identified.
type Identifier interface {
	UniqueID() uint64
	ResourceID() uint64
	setUniqueID(uint64)
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

	uuidRegister    uuidRegister
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
		Units:              UnitMillimeter,
		Language:           langUS,
		CustomContentTypes: make(map[string]string),
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
	err := m.uuidRegister.register(m.uuid, id)
	if err == nil {
		m.uuid = id
	}
	return err
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

// FindObject returns the object with the target unique ID.
func (m *Model) FindObject(uniqueID uint64) (o Object, ok bool) {
	r, k := m.FindResource(uniqueID)
	if !k {
		return
	}
	o, ok = r.(Object)
	return
}

// FindPackageResourcePath returns the package resource with the target path and ID.
func (m *Model) FindPackageResourcePath(path string, id uint64) (*ResourceID, bool) {
	return m.resourceHandler.FindResourcePath(path, id)
}

// AddResource adds a new resource to the model.
func (m *Model) AddResource(resource Identifier) error {
	id, err := m.generatePackageResourceID(resource.ResourceID())
	if err != nil {
		return err
	}
	uniqueID := id.UniqueID()
	if _, ok := m.FindResource(uniqueID); ok {
		return errors.New("go3mf: duplicated model resource")
	}
	resource.setUniqueID(uniqueID)
	m.resourceMap[uniqueID] = resource
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
