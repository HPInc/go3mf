package production

// Namespace is the canonical name of this extension.
const Namespace = "http://schemas.microsoft.com/3dmanufacturing/production/2015/06"

type Spec struct {
	LocalName  string
	IsRequired bool
}

func (e Spec) Namespace() string   { return Namespace }
func (e Spec) Required() bool      { return e.IsRequired }
func (e *Spec) SetRequired(r bool) { e.IsRequired = r }
func (e *Spec) SetLocal(l string)  { e.LocalName = l }

func (e Spec) Local() string {
	if e.LocalName != "" {
		return e.LocalName
	}
	return "p"
}

// UUID must be any of the four UUID variants described in IETF RFC 4122,
// which includes Microsoft GUIDs as well as time-based UUIDs.
type UUID string

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
