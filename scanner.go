package go3mf

import (
	"encoding/xml"
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// NodeDecoder defines the minimum contract to decode a 3MF node.
type NodeDecoder interface {
	Open()
	Attributes([]xml.Attr)
	Text([]byte)
	Child(xml.Name) NodeDecoder
	Close()
	SetScanner(*Scanner)
}

type baseDecoder struct {
	Scanner *Scanner
}

func (d *baseDecoder) Open()                      {}
func (d *baseDecoder) Attributes([]xml.Attr)      {}
func (d *baseDecoder) Text([]byte)                {}
func (d *baseDecoder) Child(xml.Name) NodeDecoder { return nil }
func (d *baseDecoder) Close()                     {}
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

// A Scanner is a 3mf model file scanning state machine.
type Scanner struct {
	Resources        []Resource
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

func newScanner() *Scanner {
	return &Scanner{
		extensionDecoder: make(map[string]*extensionDecoderWrapper),
	}
}

// Namespace returns the space of the associated local, if existing.
func (s *Scanner) Namespace(local string) (string, bool) {
	for _, name := range s.Namespaces {
		if name.Local == local {
			return name.Space, true
		}
	}
	return "", false
}

// AddResource adds a new resource to the resource cache.
func (s *Scanner) AddResource(r Resource) {
	s.Resources = append(s.Resources, r)
	s.closeResource()
}

// GenericError adds the error to the warnings.
// Returns false if scanning cannot continue.
func (s *Scanner) GenericError(strict bool, msg string) {
	err := GenericError{ResourceID: s.ResourceID, Element: s.Element, ModelPath: s.ModelPath, Message: msg}
	s.Warnings = append(s.Warnings, err)
	if strict && s.Strict {
		s.Err = err
	}
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

// ParseResourceID parses the ID as a uint32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func (s *Scanner) ParseResourceID(v string) uint32 {
	n, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		s.InvalidAttr(attrID, v, true)
		return 0
	}
	s.ResourceID = uint32(n)
	return s.ResourceID
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

// FormatMatrix converts a matrix to a string.
func FormatMatrix(t Matrix) string {
	return fmt.Sprintf("%.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f",
		t[0], t[1], t[2], t[4], t[5], t[6], t[8], t[9], t[10], t[12], t[13], t[14])
}

// ParseMatrix parses s as a Matrix.
func ParseMatrix(s string) (Matrix, bool) {
	values := strings.Fields(s)
	if len(values) != 12 {
		return Matrix{}, false
	}
	var t [12]float32
	for i := 0; i < 12; i++ {
		val, err := strconv.ParseFloat(values[i], 32)
		if err != nil {
			return Matrix{}, false
		}
		t[i] = float32(val)
	}
	return Matrix{t[0], t[1], t[2], 0.0,
		t[3], t[4], t[5], 0.0,
		t[6], t[7], t[8], 0.0,
		t[9], t[10], t[11], 1.0}, true
}

// ParseRGB parses s as a RGBA color.
func ParseRGB(s string) (c color.RGBA, err error) {
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
