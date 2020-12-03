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

// Decoder must be implemented by specs that want to support
// direct decoding from xml.
type Decoder interface {
	Namespace() string
	Local() string
	Required() bool
	DecodeAttribute(parent interface{}, attr Attr) error
	NewElementDecoder(parent interface{}, name string) ElementDecoder
}

// ElementDecoder defines the minimum contract to decode a 3MF node.
type ElementDecoder interface {
	Start([]Attr) error
	End()
}

// ChildElementDecoder must be implemented by element decoders
// that need decoding nested elements.
type ChildElementDecoder interface {
	Child(xml.Name) ElementDecoder
}

// TextElementDecoder must be implemented by element decoders
// that need to decode raw text.
type TextElementDecoder interface {
	Text([]byte)
}

// Encoder provides de necessary methods to encode specs.
// It should not be implemented by spec authors but
// will be provided be go3mf itself.
type Encoder interface {
	AddRelationship(r Relationship)
	FloatPresicion() int
	EncodeToken(t xml.Token)
	Flush() error
	SetAutoClose(autoClose bool)
}

// PreProcessEncoder must be implemented by specs
// that need to do some processing before encoding.
type PreProcessEncoder interface {
	PreProcessEncode()
}

// PostProcessDecode must be implemented by specs
// that need to do some processing after encoding.
type PostProcessorDecoder interface {
	PostProcessDecode()
}
