package encoding

import (
	"encoding/xml"
)

type Name = xml.Name
type EndElement = xml.EndElement
type CharData = xml.CharData

type Attr struct {
	Name  xml.Name
	Value []byte
}

type Relationship struct {
	Path string
	Type string
	ID   string
}

// Marshaler is the interface implemented by objects
// that can marshal themselves into valid XML elements.
type Marshaler interface {
	Marshal3MF(Encoder) error
}

// MarshalerAttr is the interface implemented by objects that can marshal
// themselves into valid XML attributes.
type MarshalerAttr interface {
	Marshal3MFAttr(Encoder) ([]xml.Attr, error)
}

// NodeDecoder defines the minimum contract to decode a 3MF node.
type NodeDecoder interface {
	Start([]Attr) error
	End()
}

type Decoder interface {
	Namespace() string
	Local() string
	Required() bool
	DecodeAttribute(interface{}, Attr) error
}

type ElementDecoder interface {
	NewElementDecoder(interface{}, string) NodeDecoder
}

type ChildNodeDecoder interface {
	Child(xml.Name) NodeDecoder
}

type TextNodeDecoder interface {
	Text([]byte)
}

type Encoder interface {
	AddRelationship(r Relationship)
	FloatPresicion() int
	EncodeToken(t xml.Token)
	Flush() error
	SetAutoClose(autoClose bool)
}

type PreProcessEncoder interface {
	PreProcessEncode()
}

type PostProcessorDecoder interface {
	PostProcessDecode()
}
