package production

import "github.com/qmuntal/go3mf"

// Namespace is the canonical name of this extension.
const Namespace = "http://schemas.microsoft.com/3dmanufacturing/production/2015/06"

type Spec struct {
	LocalName       string
	DisableAutoUUID bool
	m               *go3mf.Model
}

func (e *Spec) SetModel(m *go3mf.Model) { e.m = m }
func (e Spec) Namespace() string        { return Namespace }
func (e Spec) Required() bool           { return true }
func (e *Spec) SetRequired(r bool)      {}
func (e *Spec) SetLocal(l string)       { e.LocalName = l }

func (e Spec) Local() string {
	if e.LocalName != "" {
		return e.LocalName
	}
	return "p"
}

// BuildAttr provides a UUID in the root model file build element to ensure
// that a 3MF package can be tracked across uses by various consumers.
type BuildAttr struct {
	UUID string
}

// ObjectAttr provides a UUID in the item element
// for traceability across 3MF packages.
type ObjectAttr struct {
	UUID string
}

// ItemAttr provides a UUID in the item element to ensure
// that each object can be reliably tracked.
type ItemAttr struct {
	UUID string
	Path string
}

// ObjectPath returns the Path extension attribute.
func (p *ItemAttr) ObjectPath() string {
	return p.Path
}

func (p *ItemAttr) getUUID() string {
	return p.UUID
}

// ObjectAttr provides a UUID in the component element
// for traceability across 3MF packages.
type ComponentAttr struct {
	UUID string
	Path string
}

// ObjectPath returns the Path extension attribute.
func (p *ComponentAttr) ObjectPath() string {
	return p.Path
}

func (p *ComponentAttr) getUUID() string {
	return p.UUID
}

const (
	attrProdUUID = "UUID"
	attrPath     = "path"
)
