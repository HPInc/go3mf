package model

import "errors"

// ResourceID defines the Package Resource ID.
type ResourceID struct {
	path     string
	id       uint64 // Combination of those path and id must be unique
	uniqueID uint64 // Unique Identifier
}

// SetID sets the ID.
func (p *ResourceID) SetID(path string, id uint64) {
	p.path = path
	p.id = id
}

// ID returns the ID.
func (p *ResourceID) ID() (string, uint64) {
	return p.path, p.id
}

// SetUniqueID sets the unique ID.
func (p *ResourceID) SetUniqueID(id uint64) {
	p.uniqueID = id
}

// UniqueID returns the unique ID.
func (p *ResourceID) UniqueID() uint64 {
	return p.uniqueID
}

type resourceHandler struct {
	resourceIDs map[uint64]*ResourceID
}

func newResourceHandler() *resourceHandler {
	return &resourceHandler{
		resourceIDs: make(map[uint64]*ResourceID, 0),
	}
}

// FindResourceID search for an existing resource ID in the handler.
func (r *resourceHandler) FindResourceID(uniqueID uint64) (val *ResourceID, ok bool) {
	val, ok = r.resourceIDs[uniqueID]
	return
}

// FindResourcePath search for an existing resource ID in the handler looking by ID.
func (r *resourceHandler) FindResourcePath(path string, id uint64) (val *ResourceID, ok bool) {
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
func (r *resourceHandler) NewResourceID(path string, id uint64) (*ResourceID, error) {
	if _, ok := r.FindResourcePath(path, id); ok {
		return nil, errors.New("go3mf: Duplicate resource ID")
	}
	if len(r.resourceIDs) == 0 {
		r.resourceIDs = make(map[uint64]*ResourceID, 0)
	}
	p := &ResourceID{
		path:     path,
		id:       id,
		uniqueID: uint64(len(r.resourceIDs)) + 1,
	}
	r.resourceIDs[p.UniqueID()] = p
	return p, nil
}

// Clear resets the internal map of resource IDs.
func (r *resourceHandler) Clear() {
	r.resourceIDs = make(map[uint64]*ResourceID, 0)
}

// A Resource is an in memory representation of the 3MF resource object
type Resource struct {
	Model      *Model
	ResourceID *ResourceID
}

func newResource(id uint64, model *Model) (*Resource, error) {
	packageID, err := model.generatePackageResourceID(id)
	if err != nil {
		return nil, err
	}
	return &Resource{
		Model:      model,
		ResourceID: packageID,
	}, nil
}

// UniqueID returns the unique ID.
func (r Resource) UniqueID() uint64 {
	return r.ResourceID.UniqueID()
}
