package model

import "errors"

// PackageResourceID defines the Package Resource ID.
type PackageResourceID struct {
	path     string
	id       uint64 // Combination of those path and id must be unique
	uniqueID uint64 // Unique Identifier
}

// SetID sets the ID.
func (p *PackageResourceID) SetID(path string, id uint64) {
	p.path = path
	p.id = id
}

// ID returns the ID.
func (p *PackageResourceID) ID() (string, uint64) {
	return p.path, p.id
}

// SetUniqueID sets the unique ID.
func (p *PackageResourceID) SetUniqueID(id uint64) {
	p.uniqueID = id
}

// UniqueID returns the unique ID.
func (p *PackageResourceID) UniqueID() uint64 {
	return p.uniqueID
}

// ResourceHandler helps creating new resource identifiers.
type ResourceHandler struct {
	resourceIDs map[uint64]*PackageResourceID
}

func newResourceHandler() *ResourceHandler {
	return &ResourceHandler{
		resourceIDs: make(map[uint64]*PackageResourceID, 0),
	}
}

// FindResourceID search for an existing resource ID in the handler.
func (r *ResourceHandler) FindResourceID(uniqueID uint64) (val *PackageResourceID, ok bool) {
	val, ok = r.resourceIDs[uniqueID]
	return
}

// FindResourceIDByID search for an existing resource ID in the handler looking by ID.
func (r *ResourceHandler) FindResourceIDByID(path string, id uint64) (val *PackageResourceID, ok bool) {
	for _, value := range r.resourceIDs {
		cpath, cid := value.ID()
		if cpath == path && cid == id {
			val = value
			ok = true
			break
		}
	}
	return val, ok
}

// NewResourceID creates a new unique resource ID.
func (r *ResourceHandler) NewResourceID(path string, id uint64) (*PackageResourceID, error) {
	if _, ok := r.FindResourceIDByID(path, id); ok {
		return nil, errors.New("go3mf: Duplicate Resource ID")
	}
	if len(r.resourceIDs) == 0 {
		r.resourceIDs = make(map[uint64]*PackageResourceID, 0)
	}
	p := &PackageResourceID{
		path:     path,
		id:       id,
		uniqueID: uint64(len(r.resourceIDs)) + 1,
	}
	r.resourceIDs[p.UniqueID()] = p
	return p, nil
}

// Clear resets the internal map of resource IDs.
func (r *ResourceHandler) Clear() {
	r.resourceIDs = make(map[uint64]*PackageResourceID, 0)
}

// A Resource is an in memory representation of the 3MF resource object
type Resource struct {
	Model      *Model
	ResourceID *PackageResourceID
}

func newResource(model *Model, id uint64) (*Resource, error) {
	packageID, err := model.generatePackageResourceID(id)
	if err != nil {
		return nil, err
	}
	return &Resource{
		Model:      model,
		ResourceID: packageID,
	}, nil
}
