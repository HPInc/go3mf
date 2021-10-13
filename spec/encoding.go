// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package spec

import "encoding/xml"

// Marshaler is the interface implemented by objects
// that can marshal themselves into valid XML elements.
type Marshaler interface {
	Marshal3MF(Encoder, *xml.StartElement) error
}

// UnmarshalerAttr is the interface implemented by objects that can unmarshal
// an XML element description of themselves.
type UnmarshalerAttr interface {
	Unmarshal3MFAttr(XMLAttr) error
}

type GetterElementDecoder interface {
	ElementDecoder
	Element() interface{}
}

// ElementDecoder defines the minimum contract to decode a 3MF node.
type ElementDecoder interface {
	Start([]XMLAttr) error
	End()
}

// ChildElementDecoder must be implemented by element decoders
// that need decoding nested elements.
type ChildElementDecoder interface {
	ElementDecoder
	Child(xml.Name) (int, ElementDecoder)
}

// CharDataElementDecoder must be implemented by element decoders
// that need to decode raw text.
type CharDataElementDecoder interface {
	ElementDecoder
	CharData([]byte)
}

// AppendTokenElementDecoder must be implemented by element decoders
// that need to accumulate tokens to support loseless encoding.
type AppendTokenElementDecoder interface {
	ElementDecoder
	AppendToken(xml.Token)
}

// Encoder provides de necessary methods to encode specs.
// It should not be implemented by spec authors but
// will be provided be go3mf itself.
type Encoder interface {
	AddRelationship(Relationship)
	FloatPresicion() int
	EncodeToken(xml.Token)
	Flush() error
	SetAutoClose(bool)
	// Use SetSkipAttrEscape(true) when there is no need to escape
	// StartElement attribute values, such as as when all attributes
	// are filled using strconv.
	SetSkipAttrEscape(bool)
}
