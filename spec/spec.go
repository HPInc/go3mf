package spec

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
)

type ObjectPather interface {
	ObjectPath() string
}

type PreProcessEncoder interface {
	PreProcessEncode(m *go3mf.Model)
}

// Marshaler is the interface implemented by objects
// that can marshal themselves into valid XML elements.
type Marshaler interface {
	Marshal3MF(*go3mf.XMLEncoder) error
}

// MarshalerAttr is the interface implemented by objects that can marshal
// themselves into valid XML attributes.
type MarshalerAttr interface {
	Marshal3MFAttr(*go3mf.XMLEncoder) ([]xml.Attr, error)
}

type Decoder interface {
	Namespace() string
	Local() string
	Required() bool
	NewNodeDecoder(interface{}, string) go3mf.NodeDecoder
	DecodeAttribute(*go3mf.Scanner, interface{}, go3mf.XMLAttr)
}

type PostProcessorDecoder interface {
	PostProcessDecode(m *go3mf.Model)
}

type PropertyGroup interface {
	Len() int
}

type ModelValidator interface {
	ValidateModel(*go3mf.Model) error
}

type AssetValidator interface {
	ValidateAsset(*go3mf.Model, string, go3mf.Asset) error
}

type ObjectValidator interface {
	ValidateObject(*go3mf.Model, string, *go3mf.Object) error
}
