package go3mf

import (
	"encoding/xml"
)

type Spec interface {
	Namespace() string
	Local() string
	Required() bool
	SetRequired(bool)
	SetLocal(string)
}

type UnknownSpec struct {
	SpaceName  string
	LocalName  string
	IsRequired bool
}

func (u *UnknownSpec) Namespace() string  { return u.SpaceName }
func (u *UnknownSpec) Local() string      { return u.LocalName }
func (u *UnknownSpec) Required() bool     { return u.IsRequired }
func (u *UnknownSpec) SetLocal(l string)  { u.LocalName = l }
func (u *UnknownSpec) SetRequired(r bool) { u.IsRequired = r }

type objectPather interface {
	ObjectPath() string
}

type preProcessEncoder interface {
	PreProcessEncode(m *Model)
}

type postProcessorSpecDecoder interface {
	PostProcessDecode(m *Model)
}

type specDecoder interface {
	Namespace() string
	Local() string
	Required() bool
	NewNodeDecoder(interface{}, string) NodeDecoder
	DecodeAttribute(*Scanner, interface{}, XMLAttr)
}

type propertyGroup interface {
	Len() int
}

type marshaler interface {
	Marshal3MF(*XMLEncoder) error
}

type marshalerAttr interface {
	Marshal3MFAttr(*XMLEncoder) ([]xml.Attr, error)
}

type modelValidator interface {
	ValidateModel(*Model) error
}

type assetValidator interface {
	ValidateAsset(*Model, string, Asset) error
}

type objectValidator interface {
	ValidateObject(*Model, string, *Object) error
}
