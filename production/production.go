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

const (
	attrProdUUID = "UUID"
	attrPath     = "path"
)

// PathObject returns the Path extension attributes if exists and is not empty.
// Else it returns defaultValue.
func PathObject(o *go3mf.Object, defaultValue string) string {
	var ext *PathUUID
	if o.ExtensionAttr.Get(&ext) {
		if ext.Path != "" {
			return ext.Path
		}
	}
	return defaultValue
}

// PathObject returns the Path extension attributes if exists and is not empty.
// Else it returns defaultValue.
func PathItem(o *go3mf.Item, defaultValue string) string {
	var ext *PathUUID
	if o.ExtensionAttr.Get(&ext) {
		if ext.Path != "" {
			return ext.Path
		}
	}
	return defaultValue
}
