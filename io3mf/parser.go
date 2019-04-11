package io3mf

import (
	"fmt"
	"strconv"
)

// A MissingPropertyError represents a missing required property error.
// If MissingPropertyError is 0 means that the error took place while parsing the resource property before the ID appeared.
// When Element is 'item' the ResourceID is the objectID property of a build item.
type MissingPropertyError struct {
	ResourceID uint32
	ModelPath  string
	Element    string
	Name       string
}

func (e MissingPropertyError) Error() string {
	return fmt.Sprintf("go3mf: missing required property '%s' of element '%s' in resource '%s:%d'", e.Name, e.Element, e.ModelPath, e.ResourceID)
}

// PropertyType defines the possible property types.
type PropertyType string

const (
	// PropertyRequired is mandatory.
	PropertyRequired PropertyType = "required"
	// PropertyOptional is optional.
	PropertyOptional = "optional"
)

// A ParsePropertyError represents an error while decoding a required or an optional property.
// If ResourceID is 0 means that the error took place while parsing the resource property before the ID appeared.
// When Element is 'item' the ResourceID is the objectID property of a build item.
type ParsePropertyError struct {
	ResourceID uint32
	ModelPath  string
	Element    string
	Name       string
	Value      string
	Type       PropertyType
}

func (e ParsePropertyError) Error() string {
	return fmt.Sprintf("go3mf: [%s] error parsing property '%s = %s' of element '%s' in resource '%s:%d'", e.Type, e.Name, e.Value, e.Element, e.ModelPath, e.ResourceID)
}

// A GenericError represents a generic error.
// If ResourceID is 0 means that the error took place while parsing the resource property before the ID appeared.
// When Element is 'item' the ResourceID is the objectID property of a build item.
type GenericError struct {
	ResourceID uint32
	ModelPath  string
	Element    string
	Message    string
}

func (e GenericError) Error() string {
	return fmt.Sprintf("go3mf: error at element '%s' in resource '%s:%d' with messages '%s'", e.Element, e.ModelPath, e.ResourceID, e.Message)
}

type parser struct {
	Strict     bool
	ModelPath  string
	Element    string
	ResourceID uint32
	Err        error
	Warnings   []error
}

func (p *parser) strictError(err error) bool {
	p.Warnings = append(p.Warnings, err)
	if p.Strict {
		p.Err = err
	}
	return !p.Strict
}

func (p *parser) GenericError(strict bool, msg string) bool {
	err := GenericError{ResourceID: p.ResourceID, Element: p.Element, ModelPath: p.ModelPath, Message: msg}
	p.Warnings = append(p.Warnings, err)
	if strict && p.Strict {
		p.Err = err
		return false
	}
	return true
}

func (p *parser) InvalidRequiredAttr(attr string, val string) bool {
	return p.strictError(ParsePropertyError{ResourceID: p.ResourceID, Element: p.Element, Name: attr, Value: val, ModelPath: p.ModelPath, Type: PropertyRequired})
}

func (p *parser) InvalidOptionalAttr(attr string, val string) {
	p.Warnings = append(p.Warnings, ParsePropertyError{ResourceID: p.ResourceID, Element: p.Element, Name: attr, Value: val, ModelPath: p.ModelPath, Type: PropertyOptional})
}

func (p *parser) MissingAttr(attr string) bool {
	return p.strictError(MissingPropertyError{ResourceID: p.ResourceID, Element: p.Element, Name: attr, ModelPath: p.ModelPath})
}

func (p *parser) ParseResourceID(s string) (uint32, bool) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, p.InvalidRequiredAttr(attrID, s)
	}
	p.ResourceID = uint32(n)
	return p.ResourceID, true
}

func (p *parser) CloseResource() bool {
	if p.ResourceID == 0 {
		return p.InvalidRequiredAttr(attrID, "0")
	}
	p.ResourceID = 0
	return true
}

func (p *parser) ParseUint32Required(attr string, s string) (uint32, bool) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, p.InvalidRequiredAttr(attr, s)
	}
	return uint32(n), true
}

func (p *parser) ParseUint32Optional(attr string, s string) uint32 {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		p.InvalidOptionalAttr(attr, s)
	}
	return uint32(n)
}

func (p *parser) ParseFloat32Required(attr string, s string) (float32, bool) {
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, p.InvalidRequiredAttr(attr, s)
	}
	return float32(n), true
}

func (p *parser) ParseFloat32Optional(attr string, s string) float32 {
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		p.InvalidOptionalAttr(attr, s)
	}
	return float32(n)
}

func (p *parser) ParseFloat64Required(attr string, s string) (float64, bool) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, p.InvalidRequiredAttr(attr, s)
	}
	return n, true
}

func (p *parser) ParseFloat64Optional(attr string, s string) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		p.InvalidOptionalAttr(attr, s)
	}
	return n
}
