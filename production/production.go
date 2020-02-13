package production

import "github.com/qmuntal/go3mf"

// ExtensionName is the canonical name of this extension.
const ExtensionName = "http://schemas.microsoft.com/3dmanufacturing/production/2015/06"

// UUID must be any of the four UUID variants described in IETF RFC 4122,
// which includes Microsoft GUIDs as well as time-based UUIDs.
type UUID string

// NewUUID creates a UUID from s.
func NewUUID(s string) (UUID, error) {
	if err := validateUUID(s); err != nil {
		return UUID(""), err
	}
	return UUID(s), nil
}

type PathUUID struct {
	UUID UUID
	Path string
}

func extractUUID(ext go3mf.ExtensionAttr) *UUID {
	if attr, ok := ext[ExtensionName]; ok {
		return attr.(*UUID)
	}
	attr := UUID("")
	pa := &attr
	ext[ExtensionName] = pa
	return pa
}

func extractPathUUID(ext go3mf.ExtensionAttr) *PathUUID {
	if attr, ok := ext[ExtensionName]; ok {
		return attr.(*PathUUID)
	}
	attr := &PathUUID{}
	ext[ExtensionName] = attr
	return attr
}

// BuildAttr extracts the UUID attributes from a Build.
// Returns an empty UUID if it does not exist, never nil.
func BuildAttr(b *go3mf.Build) *UUID {
	if b.ExtensionAttr == nil {
		b.ExtensionAttr = make(go3mf.ExtensionAttr)
	}
	return extractUUID(b.ExtensionAttr)
}

// ItemAttr extracts the Path and UUID attributes from an Item.
// Returns an empty PathUUID if it does not exist, never nil.
func ItemAttr(item *go3mf.Item) *PathUUID {
	if item.ExtensionAttr == nil {
		item.ExtensionAttr = make(go3mf.ExtensionAttr)
	}
	return extractPathUUID(item.ExtensionAttr)
}

// ComponentAttr extracts the Pathn and UUID attributes from a Component.
// Returns an empty PathUUID if it does not exist, never nil.
func ComponentAttr(c *go3mf.Component) *PathUUID {
	if c.ExtensionAttr == nil {
		c.ExtensionAttr = make(go3mf.ExtensionAttr)
	}
	return extractPathUUID(c.ExtensionAttr)
}

// ObjectAttr extracts the UUID attributes from a ObjectResource.
// Returns an empty UUID if it does not exist, never nil.
func ObjectAttr(o *go3mf.ObjectResource) *UUID {
	if o.ExtensionAttr == nil {
		o.ExtensionAttr = make(go3mf.ExtensionAttr)
	}
	return extractUUID(o.ExtensionAttr)
}

const (
	attrProdUUID = "UUID"
	attrPath     = "path"
)
