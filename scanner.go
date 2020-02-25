package go3mf

import (
	"encoding/xml"
	"errors"
	"fmt"
	"image/color"
)

// NodeDecoder defines the minimum contract to decode a 3MF node.
type NodeDecoder interface {
	Start([]xml.Attr)
	Text([]byte)
	Child(xml.Name) NodeDecoder
	End()
	SetScanner(*Scanner)
}

type baseDecoder struct {
	Scanner *Scanner
}

func (d *baseDecoder) Start([]xml.Attr)           {}
func (d *baseDecoder) Text([]byte)                {}
func (d *baseDecoder) Child(xml.Name) NodeDecoder { return nil }
func (d *baseDecoder) End()                       {}
func (d *baseDecoder) SetScanner(s *Scanner)      { d.Scanner = s }

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

// A Scanner is a 3mf model file scanning state machine.
type Scanner struct {
	Resources        Resources
	BuildItems       []*Item
	Strict           bool
	ModelPath        string
	IsRoot           bool
	Element          string
	ResourceID       uint32
	Err              error
	Warnings         []error
	Namespaces       []xml.Name
	extensionDecoder map[string]*extensionDecoderWrapper
}

// Namespace returns the space of the associated local, if existing.
func (s *Scanner) namespace(local string) (string, bool) {
	for _, name := range s.Namespaces {
		if name.Local == local {
			return name.Space, true
		}
	}
	return "", false
}

// AddAsset adds a new resource to the resource cache.
func (s *Scanner) AddAsset(r Asset) {
	s.Resources.Assets = append(s.Resources.Assets, r)
	s.closeResource()
}

// AddObject adds a new resource to the resource cache.
func (s *Scanner) AddObject(r *Object) {
	s.Resources.Objects = append(s.Resources.Objects, r)
	s.closeResource()
}

// InvalidAttr adds the error to the warnings.
// Returns false if scanning cannot continue.
func (s *Scanner) InvalidAttr(attr string, val string, required bool) {
	tp := PropertyRequired
	if !required {
		tp = PropertyOptional
	}
	s.strictError(ParsePropertyError{ResourceID: s.ResourceID, Element: s.Element, Name: attr, Value: val, ModelPath: s.ModelPath, Type: tp})
}

// MissingAttr adds the error to the warnings.
func (s *Scanner) MissingAttr(attr string) {
	s.strictError(MissingPropertyError{ResourceID: s.ResourceID, Element: s.Element, Name: attr, ModelPath: s.ModelPath})
}

func (s *Scanner) strictError(err error) {
	s.Warnings = append(s.Warnings, err)
	if s.Strict {
		s.Err = err
	}
}

// closeResource closes the current resource.
// If there is no resource to close MissingPropertyError is added to the warnings.
func (s *Scanner) closeResource() {
	if s.ResourceID == 0 {
		s.MissingAttr(attrID)
		return
	}
	s.ResourceID = 0
}

// ParseRGBA parses s as a RGBA color.
func ParseRGBA(s string) (c color.RGBA, err error) {
	var errInvalidFormat = errors.New("gltf: invalid color format")

	if len(s) == 0 || s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 9:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
		c.A = hexToByte(s[7])<<4 + hexToByte(s[8])
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
		c.A = 0xff
	default:
		err = errInvalidFormat
	}
	return
}

// FormatRGBA returns the color as a hex string with the format #rrggbbaa.
func FormatRGBA(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}
