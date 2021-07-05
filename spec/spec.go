// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package spec

import (
	"encoding/xml"
)

// Spec is the interface that must be implemented by a 3mf spec.
//
// Specs may implement ValidateSpec.
type Spec interface {
	DecodeAttribute(parent interface{}, attr Attr) error
	CreateElementDecoder(parent interface{}, name string) ElementDecoder
}

type PropertyGroup interface {
	Len() int
}

// If a Spec implemented ValidateSpec, then model.Validate will call
// Validate and aggregate the resulting erros.
//
// model is guaranteed to be a *go3mf.Model
// element can be a *go3mf.Model, go3mf.Asset or *go3mf.Object.
type ValidateSpec interface {
	Spec
	Validate(model interface{}, path string, element interface{}) error
}

// An Attr represents an attribute in an XML element (Name=Value).
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

type ErrorWrapper interface {
	Wrap(error) error
}

// ElementDecoder defines the minimum contract to decode a 3MF node.
type ElementDecoder interface {
	Start([]Attr) error
	End()
}

// ChildElementDecoder must be implemented by element decoders
// that need decoding nested elements.
type ChildElementDecoder interface {
	ElementDecoder
	Child(xml.Name) ElementDecoder
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

// An UnknownAttr represents an attribute
// that is not supported by any loaded Spec.
type UnknownAttr []xml.Attr

func (u *UnknownAttr) AppendAttr(att Attr) {
	*u = append(*u, xml.Attr{Name: att.Name, Value: string(att.Value)})
}

func (u UnknownAttr) Marshal3MFAttr(enc Encoder) ([]xml.Attr, error) {
	return u, nil
}

// UnknownTokens represents a section of an xml
// that cannot be decoded by any loaded Spec.
type UnknownTokens []xml.Token

func (u UnknownTokens) Marshal3MF(enc Encoder) error {
	for _, t := range u {
		enc.EncodeToken(t)
	}
	return nil
}

// UnknownTokensDecoder can be used by spec decoders to maintain the
// xml tree elements of unknown extensions.
type UnknownTokensDecoder struct {
	Name xml.Name

	tokens UnknownTokens
}

func (d *UnknownTokensDecoder) Start(attrs []Attr) error {
	var xattrs []xml.Attr
	if len(attrs) > 0 {
		xattrs = make([]xml.Attr, len(attrs))
		for i, att := range attrs {
			xattrs[i] = xml.Attr{Name: att.Name, Value: string(att.Value)}
		}
	}
	d.AppendToken(xml.StartElement{
		Name: d.Name,
		Attr: xattrs,
	})
	return nil
}

func (d *UnknownTokensDecoder) End() {
	d.AppendToken(xml.EndElement{Name: d.Name})
}

func (d *UnknownTokensDecoder) AppendToken(t xml.Token) {
	d.tokens = append(d.tokens, t)
}

func (d UnknownTokensDecoder) Tokens() UnknownTokens {
	return d.tokens
}
