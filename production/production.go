// Package production handles new non-object resources, 
// as well as attributes to the build section for uniquely identifying parts within a particular 3MF package
// Despite item and component paths are production attributes, they are also handled by
// the core package, to avoid duplications they won't be stored in the Extension map
// but the core properties will be updated.
package production

import "github.com/qmuntal/go3mf"

// ExtensionName is the canonical name of this extension.
const ExtensionName = "http://schemas.microsoft.com/3dmanufacturing/production/2015/06"

// UUID must be any of the four UUID variants described in IETF RFC 4122,
// which includes Microsoft GUIDs as well as time-based UUIDs.
type UUID string

func extractUUID(ext map[string]interface{}) UUID {
	if attr, ok := ext[ExtensionName]; ok {
		return attr.(UUID)
	}
	return UUID("")
}

func setUUID(ext map[string]interface{}, u UUID) {
	if ext == nil {
		ext = make(map[string]interface{})
	}
	ext[ExtensionName] = &u
}

// BuildUUID extracts the UUID attributes from a Build.
// Returns an empty UUID if it does not exist.
func BuildUUID(b *go3mf.Build) UUID {
	return extractUUID(b.Extensions)
}

// SetBuildUUID sets the UUID.
func SetBuildUUID(b *go3mf.Build, u UUID) {
	if b.Extensions == nil {
		b.Extensions = make(map[string]interface{})
	}
	b.Extensions[ExtensionName] = u
}

// ItemUUID extracts the UUID attributes from an Item.
// Returns an empty UUID if it does not exist.
func ItemUUID(o *go3mf.Item) UUID {
	return extractUUID(o.Extensions)
}

// SetItemdUUID sets the UUID.
func SetItemdUUID(i *go3mf.Item, u UUID) {
	if i.Extensions == nil {
		i.Extensions = make(map[string]interface{})
	}
	i.Extensions[ExtensionName] = u
}

// ComponentUUID extracts the UUID attributes from a Component.
// Returns an empty UUID if it does not exist.
func ComponentUUID(c *go3mf.Component) UUID {
	return extractUUID(c.Extensions)
}

// SetComponentUUID sets the UUID.
func SetComponentUUID(c *go3mf.Component, u UUID) {
	if c.Extensions == nil {
		c.Extensions = make(map[string]interface{})
	}
	c.Extensions[ExtensionName] = u
}

// ObjectUUID extracts the UUID attributes from a ObjectResource.
// Returns an empty UUID if it does not exist.
func ObjectUUID(o *go3mf.ObjectResource) UUID {
	return extractUUID(o.Extensions)
}

// SetObjectUUID sets the UUID.
func SetObjectUUID(o *go3mf.ObjectResource, u UUID) {
	if o.Extensions == nil {
		o.Extensions = make(map[string]interface{})
	}
	o.Extensions[ExtensionName] = u
}

const (
	attrProdUUID = "UUID"
	attrPath     = "path"
)
