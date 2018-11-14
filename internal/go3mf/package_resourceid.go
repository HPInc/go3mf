package go3mf

// PackageResourceID defines the Package Resource ID.
type PackageResourceID struct {
	path     string
	id       uint64 // Combination of those path and id must be unique
	uniqueID uint64 // Unique Identifier
}

// SetID sets the id.
func (p *PackageResourceID) SetID(path string, id uint64) {
	p.path = path
	p.id = id
}

// ID returns the id.
func (p *PackageResourceID) ID() (string, uint64) {
	return p.path, p.id
}

// SetUniqueID sets the unique id.
func (p *PackageResourceID) SetUniqueID(id uint64) {
	p.uniqueID = id
}

// UniqueID returns the unique id.
func (p *PackageResourceID) UniqueID() uint64 {
	return p.uniqueID
}
