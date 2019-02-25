package model

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
}

func (m *Model) generatePackageResourceID(id uint64) (*PackageResourceID, error) {
	return m.resourceHandler.NewResourceID(m.Path, id)
}
