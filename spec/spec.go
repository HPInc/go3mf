// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package spec

import (
	"encoding/xml"
	"sync"
)

// Spec is the interface that must be implemented by a 3mf spec.
//
// Specs may implement ValidateSpec.
type Spec interface {
	NewAttr3MF(parent string) AttrGroup
	NewElementDecoder(parent interface{}, name string) ElementDecoder
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

// An XMLAttr represents an attribute in an XML element (Name=Value).
type XMLAttr struct {
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
	Marshal3MF(Encoder, *xml.StartElement) error
}

// UnmarshalerAttr is the interface implemented by objects that can unmarshal
// an XML element description of themselves.
type UnmarshalerAttr interface {
	Unmarshal3MFAttr(XMLAttr) error
}

type ErrorWrapper interface {
	Wrap(error) error
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

var (
	specMu sync.RWMutex
	specs  = make(map[string]Spec)
)

// Register makes a spec available by the provided namesoace.
// If Register is called twice with the same name or if spec is nil,
// it panics.
func Register(namespace string, spec Spec) {
	specMu.Lock()
	defer specMu.Unlock()
	specs[namespace] = spec
}

// AttrGroup defines a container for different attributes of the same namespace.
// It supports encoding and decoding to XML.
type AttrGroup interface {
	UnmarshalerAttr
	Marshaler
	Namespace() string
}

type AnyAttr []AttrGroup

func (a AnyAttr) Get(namespace string) AttrGroup {
	for _, v := range a {
		if v.Namespace() == namespace {
			return v
		}
	}
	return nil
}

func (a AnyAttr) Marshal3MF(x Encoder, start *xml.StartElement) error {
	for _, ext := range a {
		err := ext.Marshal3MF(x, start)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewAttr3MF(namespace, parent string) AttrGroup {
	if ext, ok := LoadExtension(namespace); ok {
		return ext.NewAttr3MF(parent)
	}
	return &UnknownAttrs{
		Space: namespace,
	}
}

// Any is an extension point containing <any> information.
type Any []Marshaler

func (e Any) Marshal3MF(x Encoder, start *xml.StartElement) error {
	for _, ext := range e {
		if err := ext.Marshal3MF(x, start); err == nil {
			return err
		}
	}
	return nil
}

func LoadExtension(space string) (Spec, bool) {
	specMu.RLock()
	ext, ok := specs[space]
	specMu.RUnlock()
	return ext, ok
}

func LoadValidator(ns string) (ValidateSpec, bool) {
	specMu.RLock()
	ext, ok := specs[ns]
	specMu.RUnlock()
	if ok {
		ext, ok := ext.(ValidateSpec)
		return ext, ok
	}
	return nil, false
}
