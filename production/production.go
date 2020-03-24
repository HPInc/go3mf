package production

// ExtensionName is the canonical name of this extension.
const ExtensionName = "http://schemas.microsoft.com/3dmanufacturing/production/2015/06"

type Extension struct {
	LocalName  string
	IsRequired bool
}

func (e Extension) Name() string { return ExtensionName }

func (e Extension) Local() string {
	if e.LocalName != "" {
		return e.LocalName
	}
	return "p"
}

func (e Extension) Required() bool {
	return e.IsRequired
}

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

// ObjectPath returns the Path extension attribute.
func (p *PathUUID) ObjectPath() string {
	return p.Path
}

const (
	attrProdUUID = "UUID"
	attrPath     = "path"
)
