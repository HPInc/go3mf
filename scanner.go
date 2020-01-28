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

// BaseDecoder defines a base class for all decoders.
// It is not mandatory for a Decoder to embed this struct,
// but if embedded any struct automatically fulfills the NodeDecoder interface.
// It is typically used when creating extension decoders.
type BaseDecoder struct {
	Scanner *Scanner
}

// Open do nothing.
func (d *BaseDecoder) Open() { return }

// Attributes do nothing.
func (d *BaseDecoder) Attributes([]xml.Attr) { return }

// Text do nothing.
func (d *BaseDecoder) Text([]byte) { return }

// Child returns nil.
func (d *BaseDecoder) Child(xml.Name) NodeDecoder { return nil }

// Close do nothing.
func (d *BaseDecoder) Close() { return }

// SetScanner sets the scanner.
func (d *BaseDecoder) SetScanner(s *Scanner) { d.Scanner = s }

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
	Resources    []Resource
	BuildItems   []*BuildItem
	UUID         string
	Strict       bool
	ModelPath    string
	IsRoot       bool
	Element      string
	ResourceID   uint32
	Err          error
	Warnings     []error
	Namespaces   map[string]string
	model        *Model
	resourcesMap map[uint32]Resource
}

// NewScanner returns an initialized scanner.
func NewScanner(model *Model) *Scanner {
	return &Scanner{
		model:        model,
		Namespaces:   make(map[string]string),
		resourcesMap: make(map[uint32]Resource),
	}
}

// AddResource adds a new resource to the resource cache.
func (p *Scanner) AddResource(r Resource) {
	_, id := r.Identify()
	p.resourcesMap[id] = r
	p.Resources = append(p.Resources, r)
}

// FindResource returns the resource with the target unique ID.
func (p *Scanner) FindResource(path string, id uint32) (r Resource, ok bool) {
	if path == "" {
		path = p.model.Path
	}
	if path == p.ModelPath {
		r, ok = p.resourcesMap[id]
	} else {
		r, ok = p.model.FindResource(path, id)
	}
	return
}

// NamespaceRegistered checks if the namespace is registered.
func (p *Scanner) NamespaceRegistered(ns string) bool {
	for _, space := range p.Namespaces {
		if ns == space {
			return true
		}
	}
	return false
}

func (p *Scanner) strictError(err error) {
	p.Warnings = append(p.Warnings, err)
	if p.Strict {
		p.Err = err
	}
}

// GenericError adds the error to the warnings.
// Returns false if scanning cannot continue.
func (p *Scanner) GenericError(strict bool, msg string) {
	err := GenericError{ResourceID: p.ResourceID, Element: p.Element, ModelPath: p.ModelPath, Message: msg}
	p.Warnings = append(p.Warnings, err)
	if strict && p.Strict {
		p.Err = err
	}
}

// InvalidRequiredAttr adds the error to the warnings.
// Returns false if scanning cannot continue.
func (p *Scanner) InvalidRequiredAttr(attr string, val string) {
	p.strictError(ParsePropertyError{ResourceID: p.ResourceID, Element: p.Element, Name: attr, Value: val, ModelPath: p.ModelPath, Type: PropertyRequired})
}

// InvalidOptionalAttr adds the error to the warnings.
func (p *Scanner) InvalidOptionalAttr(attr string, val string) {
	p.Warnings = append(p.Warnings, ParsePropertyError{ResourceID: p.ResourceID, Element: p.Element, Name: attr, Value: val, ModelPath: p.ModelPath, Type: PropertyOptional})
}

// MissingAttr adds the error to the warnings.
func (p *Scanner) MissingAttr(attr string) {
	p.strictError(MissingPropertyError{ResourceID: p.ResourceID, Element: p.Element, Name: attr, ModelPath: p.ModelPath})
}

// ParseResourceID parses the ID as a uint32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func (p *Scanner) ParseResourceID(s string) uint32 {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		p.InvalidRequiredAttr("id", s)
		return 0
	}
	p.ResourceID = uint32(n)
	return p.ResourceID
}

// CloseResource closes the current resource.
// If there is no resource to close MissingPropertyError is added to the warnings.
func (p *Scanner) CloseResource() {
	if p.ResourceID == 0 {
		p.MissingAttr("id")
		return
	}
	p.ResourceID = 0
}

// ParseUint32Required parses s as a uint32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func (p *Scanner) ParseUint32Required(attr string, s string) uint32 {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		p.InvalidRequiredAttr(attr, s)
		return 0
	}
	return uint32(n)
}

// ParseUint32Optional parses s as a uint32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func (p *Scanner) ParseUint32Optional(attr string, s string) uint32 {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		p.InvalidOptionalAttr(attr, s)
	}
	return uint32(n)
}

// ParseFloat32Required parses s as a float32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func (p *Scanner) ParseFloat32Required(attr string, s string) float32 {
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		p.InvalidRequiredAttr(attr, s)
		return 0
	}
	return float32(n)
}

// ParseFloat32Optional parses s as a float32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func (p *Scanner) ParseFloat32Optional(attr string, s string) float32 {
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		p.InvalidOptionalAttr(attr, s)
	}
	return float32(n)
}

// ParseToMatrix parses s as a Matrix.
func ParseToMatrix(s string) (Matrix, bool) {
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
