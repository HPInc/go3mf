package io3mf

import (
	"fmt"
	"strconv"
)

// A MissingAttrError represents a missing required attribute error.
// If MissingAttrError is 0 means that the error took place while parsing the resource attributes before the ID appeared.
type MissingAttrError struct {
	ResourceID uint32
	ModelPath  string
	Element    string
	Attr       string
}

func (e MissingAttrError) Error() string {
	return fmt.Sprintf("go3mf: missing required required attribute '%s' of element '%s' in resource '%d'", e.Attr, e.Element, e.ResourceID)
}

// A RequiredAttrError represents an error while decoding a required attribute.
// If ResourceID is 0 means that the error took place while parsing the resource attributes before the ID appeared.
type RequiredAttrError struct {
	ResourceID uint32
	ModelPath  string
	Element    string
	Attr       string
}

func (e RequiredAttrError) Error() string {
	return fmt.Sprintf("go3mf: error decoding a required attribute '%s' of element '%s' in resource '%d'", e.Attr, e.Element, e.ResourceID)
}

// A OptionalAttrError represents an error while decoding aan optional attribute.
// If ResourceID is 0 means that the error took place while parsing the resource attributes before the ID appeared.
type OptionalAttrError struct {
	ResourceID uint32
	ModelPath  string
	Element    string
	Attr       string
}

func (e OptionalAttrError) Error() string {
	return fmt.Sprintf("go3mf: error decoding an optional attribute '%s' of element '%s' in resource '%d'", e.Attr, e.Element, e.ResourceID)
}

// A GenericError represents a generic error.
// If ResourceID is 0 means that the error took place while parsing the resource attributes before the ID appeared.
type GenericError struct {
	ResourceID uint32
	ModelPath  string
	Element    string
	Message    string
}

func (e GenericError) Error() string {
	return fmt.Sprintf("go3mf: error at element '%s' in resource '%d' with messages '%s'", e.Element, e.ResourceID, e.Message)
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

func (p *parser) InvalidRequiredAttr(attr string) bool {
	return p.strictError(RequiredAttrError{ResourceID: p.ResourceID, Element: p.Element, Attr: attr, ModelPath: p.ModelPath})
}

func (p *parser) InvalidOptionalAttr(attr string) {
	p.Warnings = append(p.Warnings, OptionalAttrError{ResourceID: p.ResourceID, Element: p.Element, Attr: attr, ModelPath: p.ModelPath})
}

func (p *parser) MissingAttr(attr string) bool {
	return p.strictError(MissingAttrError{ResourceID: p.ResourceID, Element: p.Element, Attr: attr, ModelPath: p.ModelPath})
}

func (p *parser) ParseResourceID(s string) (uint32, bool) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, p.strictError(RequiredAttrError{ResourceID: p.ResourceID, Element: p.Element, Attr: attrID, ModelPath: p.ModelPath})
	}
	p.ResourceID = uint32(n)
	return p.ResourceID, true
}

func (p *parser) CloseResource() bool {
	if p.ResourceID == 0 {
		return p.strictError(MissingAttrError{ResourceID: 0, Element: p.Element, Attr: attrID, ModelPath: p.ModelPath})
	}
	p.ResourceID = 0
	return true
}

func (p *parser) ParseUint32Required(attr string, s string) (uint32, bool) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, p.strictError(RequiredAttrError{ResourceID: p.ResourceID, Element: p.Element, Attr: attr, ModelPath: p.ModelPath})
	}
	return uint32(n), true
}

func (p *parser) ParseUint32Optional(attr string, s string) uint32 {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		p.InvalidOptionalAttr(attr)
	}
	return uint32(n)
}

func (p *parser) ParseFloat32Required(attr string, s string) (float32, bool) {
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, p.strictError(RequiredAttrError{ResourceID: p.ResourceID, Element: p.Element, Attr: attr, ModelPath: p.ModelPath})
	}
	return float32(n), true
}

func (p *parser) ParseFloat32Optional(attr string, s string) float32 {
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		p.InvalidOptionalAttr(attr)
	}
	return float32(n)
}

func (p *parser) ParseFloat64Required(attr string, s string) (float64, bool) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, p.strictError(RequiredAttrError{ResourceID: p.ResourceID, Element: p.Element, Attr: attr, ModelPath: p.ModelPath})
	}
	return n, true
}

func (p *parser) ParseFloat64Optional(attr string, s string) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		p.InvalidOptionalAttr(attr)
	}
	return n
}
