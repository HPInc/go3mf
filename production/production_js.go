package production

import (
	"syscall/js"

	"github.com/qmuntal/go3mf"
)

var (
	jsNS       = "PRODUCTION"
	jsSpec     = go3mf.RegisterClass(jsNS, "X")
	jsUUID     = go3mf.RegisterClassExtends("Uuid", "String", "X", jsNS)
	jsPathUUID = go3mf.RegisterClass("PathUuid", "X", jsNS)
)

// JSValue returns a JavaScript value associated with the object.
func (p *UUID) JSValue() js.Value {
	return jsUUID.New(string(*p))
}

// JSValue returns a JavaScript value associated with the object.
func (p *PathUUID) JSValue() js.Value {
	v := jsPathUUID.New()
	if p.Path != "" {
		v.Set(attrPath, js.Undefined())
	} else {
		v.Set(attrPath, p.Path)
	}
	v.Set(attrProdUUID, string(p.UUID))
	return v
}
