package iohelper

import (
	"encoding/xml"
	"errors"
	"image/color"
	"strconv"
	"strings"

	"github.com/qmuntal/go3mf"
)

// NodeDecoder defines the minimum contract to decode a 3MF node.
type NodeDecoder interface {
	Open()
	Attributes([]xml.Attr) bool
	Text([]byte) bool
	Child(xml.Name) NodeDecoder
	Close() bool
	SetScanner(*Scanner)
}

// EmptyDecoder defines a base class for all decoders.
type EmptyDecoder struct {
	Scanner *Scanner
}

// Open do nothing.
func (d *EmptyDecoder) Open() { return }

// Attributes returns true.
func (d *EmptyDecoder) Attributes([]xml.Attr) bool { return true }

// Text returns true.
func (d *EmptyDecoder) Text([]byte) bool { return true }

// Child returns nil.
func (d *EmptyDecoder) Child(xml.Name) NodeDecoder { return nil }

// Close returns true.
func (d *EmptyDecoder) Close() bool { return true }

// SetScanner sets the scanner.
func (d *EmptyDecoder) SetScanner(s *Scanner) { d.Scanner = s }

// A Scanner is a 3mf model file scanning state machine.
type Scanner struct {
	Resources    []go3mf.Resource
	BuildItems   []*go3mf.BuildItem
	UUID         string
	Strict       bool
	ModelPath    string
	IsRoot       bool
	Element      string
	ResourceID   uint32
	Err          error
	Warnings     []error
	Namespaces   map[string]string
	model        *go3mf.Model
	resourcesMap map[uint32]go3mf.Resource
}

// NewScanner returns an initialized scanner.
func NewScanner(model *go3mf.Model) *Scanner {
	return &Scanner{
		model:        model,
		Namespaces:   make(map[string]string),
		resourcesMap: make(map[uint32]go3mf.Resource),
	}
}

// AddResource adds a new resource to the resource cache.
func (p *Scanner) AddResource(r go3mf.Resource) {
	_, id := r.Identify()
	p.resourcesMap[id] = r
	p.Resources = append(p.Resources, r)
}

// FindResource returns the resource with the target unique ID.
func (p *Scanner) FindResource(path string, id uint32) (r go3mf.Resource, ok bool) {
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

func (p *Scanner) strictError(err error) bool {
	p.Warnings = append(p.Warnings, err)
	if p.Strict {
		p.Err = err
	}
	return !p.Strict
}

// GenericError adds the error to the warnings.
// Returns false if scanning cannot continue.
func (p *Scanner) GenericError(strict bool, msg string) bool {
	err := go3mf.GenericError{ResourceID: p.ResourceID, Element: p.Element, ModelPath: p.ModelPath, Message: msg}
	p.Warnings = append(p.Warnings, err)
	if strict && p.Strict {
		p.Err = err
		return false
	}
	return true
}

// InvalidRequiredAttr adds the error to the warnings.
// Returns false if scanning cannot continue.
func (p *Scanner) InvalidRequiredAttr(attr string, val string) bool {
	return p.strictError(go3mf.ParsePropertyError{ResourceID: p.ResourceID, Element: p.Element, Name: attr, Value: val, ModelPath: p.ModelPath, Type: go3mf.PropertyRequired})
}

// InvalidOptionalAttr adds the error to the warnings.
func (p *Scanner) InvalidOptionalAttr(attr string, val string) {
	p.Warnings = append(p.Warnings, go3mf.ParsePropertyError{ResourceID: p.ResourceID, Element: p.Element, Name: attr, Value: val, ModelPath: p.ModelPath, Type: go3mf.PropertyOptional})
}

// MissingAttr adds the error to the warnings.
// Returns false if scanning cannot continue.
func (p *Scanner) MissingAttr(attr string) bool {
	return p.strictError(go3mf.MissingPropertyError{ResourceID: p.ResourceID, Element: p.Element, Name: attr, ModelPath: p.ModelPath})
}

// ParseResourceID parses the ID as a uint32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
// Returns false if scanning cannot continue.
func (p *Scanner) ParseResourceID(s string) (uint32, bool) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, p.InvalidRequiredAttr("id", s)
	}
	p.ResourceID = uint32(n)
	return p.ResourceID, true
}

// CloseResource closes the current resource.
// If there is no resource to close MissingPropertyError is added to the warnings.
// Returns false if scanning cannot continue.
func (p *Scanner) CloseResource() bool {
	if p.ResourceID == 0 {
		return p.MissingAttr("id")
	}
	p.ResourceID = 0
	return true
}

// ParseUint32Required parses s as a uint32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
// Returns false if scanning cannot continue.
func (p *Scanner) ParseUint32Required(attr string, s string) (uint32, bool) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, p.InvalidRequiredAttr(attr, s)
	}
	return uint32(n), true
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
// Returns false if scanning cannot continue.
func (p *Scanner) ParseFloat32Required(attr string, s string) (float32, bool) {
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, p.InvalidRequiredAttr(attr, s)
	}
	return float32(n), true
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

// ParseFloat64Required parses s as a float64.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
// Returns false if scanning cannot continue.
func (p *Scanner) ParseFloat64Required(attr string, s string) (float64, bool) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, p.InvalidRequiredAttr(attr, s)
	}
	return n, true
}

// ParseFloat64Optional parses s as a float64.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func (p *Scanner) ParseFloat64Optional(attr string, s string) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		p.InvalidOptionalAttr(attr, s)
	}
	return n
}

// ParseToMatrixRequired parses s as a go3mf.Matrix.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
// Returns false if scanning cannot continue.
func (p *Scanner) ParseToMatrixRequired(attr string, s string) (go3mf.Matrix, bool) {
	var matrix go3mf.Matrix
	values := strings.Fields(s)
	if len(values) != 12 {
		return matrix, p.InvalidRequiredAttr(attr, s)
	}
	var t [12]float32
	for i := 0; i < 12; i++ {
		val, err := strconv.ParseFloat(values[i], 32)
		if err != nil {
			return matrix, p.InvalidRequiredAttr(attr, s)
		}
		t[i] = float32(val)
	}
	return go3mf.Matrix{t[0], t[1], t[2], 0.0,
		t[3], t[4], t[5], 0.0,
		t[6], t[7], t[8], 0.0,
		t[9], t[10], t[11], 1.0}, true
}

// ParseToMatrixOptional parses s as a go3mf.Matrix.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func (p *Scanner) ParseToMatrixOptional(attr string, s string) go3mf.Matrix {
	values := strings.Fields(s)
	if len(values) != 12 {
		p.InvalidOptionalAttr(attr, s)
		return go3mf.Matrix{}
	}
	var t [12]float32
	for i := 0; i < 12; i++ {
		val, err := strconv.ParseFloat(values[i], 32)
		if err != nil {
			p.InvalidOptionalAttr(attr, s)
			return go3mf.Matrix{}
		}
		t[i] = float32(val)
	}
	return go3mf.Matrix{t[0], t[1], t[2], 0.0,
		t[3], t[4], t[5], 0.0,
		t[6], t[7], t[8], 0.0,
		t[9], t[10], t[11], 1.0}
}

// ReadRGB parses s as a RGBA color.
func ReadRGB(s string) (c color.RGBA, err error) {
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
