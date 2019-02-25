package model

import (
	"errors"

	"github.com/gofrs/uuid"
)

// Metadata item is an in memory representation of the 3MF metadata,
// and can be attached to any 3MF model node.
type Metadata struct {
	Name  string
	Value string
}

// A Model is an in memory representation of the 3MF file.
type Model struct {
	Path            string
	resourceHandler ResourceHandler
	usedUUIDs       map[uuid.UUID]struct{}
}

func (m *Model) generatePackageResourceID(id uint64) (*PackageResourceID, error) {
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
