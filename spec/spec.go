// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package spec

import (
	"encoding/xml"
	"sync"
)

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

func Load(space string) (Spec, bool) {
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

// Spec is the interface that must be implemented by a 3mf spec.
//
// Specs may implement ValidateSpec.
type Spec interface {
	NewAttrGroup(parent xml.Name) AttrGroup
	NewElementDecoder(name xml.Name) GetterElementDecoder
}

// If a Spec implemented ValidateSpec, then model.Validate will call
// Validate and aggregate the resulting erros.
//
// model is guaranteed to be a *go3mf.Model
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

// AttrGroup defines a container for different attributes of the same namespace.
// It supports encoding and decoding to XML.
type AttrGroup interface {
	UnmarshalerAttr
	Marshaler
	Namespace() string
}

type PropertyGroup interface {
	Len() int
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

func NewAttrGroup(namespace string, parent xml.Name) AttrGroup {
	if ext, ok := Load(namespace); ok {
		return ext.NewAttrGroup(parent)
	}
	return &UnknownAttrs{
		Space: namespace,
	}
}

func NewElementDecoder(name xml.Name) GetterElementDecoder {
	if ext, ok := Load(name.Space); ok {
		return ext.NewElementDecoder(name)
	}
	return &UnknownTokensDecoder{XMLName: name}
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
