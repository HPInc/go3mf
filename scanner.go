package go3mf

import (
	"encoding/xml"
	"errors"
	"fmt"
	"image/color"
	"strings"

	specerr "github.com/qmuntal/go3mf/errors"
)

type XMLAttr struct {
	Name  xml.Name
	Value []byte
}

// NodeDecoder defines the minimum contract to decode a 3MF node.
type NodeDecoder interface {
	Start([]XMLAttr) error
	Text([]byte)
	Child(xml.Name) NodeDecoder
	End()
}

type baseDecoder struct {
}

func (d *baseDecoder) Start([]XMLAttr) error      { return nil }
func (d *baseDecoder) Text([]byte)                {}
func (d *baseDecoder) Child(xml.Name) NodeDecoder { return nil }
func (d *baseDecoder) End()                       {}

// A decoderContext is a 3mf model file scanning state machine.
type decoderContext struct {
	Err              specerr.List
	extensionDecoder map[string]specDecoder
	modelPath        string
	isRoot           bool
	contex           []xml.Name
	resourceID       uint32
}

func (s *decoderContext) namespace(local string) (string, bool) {
	for _, ext := range s.extensionDecoder {
		if ext.Local() == local {
			return ext.Namespace(), true
		}
	}
	return "", false
}

func (s *decoderContext) addErrContext(err error) {
	ct := make([]string, len(s.contex))
	ct[0] = s.modelPath
	for i, s := range s.contex[1:] {
		ct[i+1] = s.Local
	}
	if s.isRoot {
		ct = ct[1:] // don't add path in case happend in root file
	}
	ctx := strings.Join(ct, "@")
	if err, ok := err.(*specerr.List); ok {
		for _, err := range err.Errors {
			specerr.Append(&s.Err, &specerr.ResourceError{
				Context: ctx, ResourceID: s.resourceID, Err: err,
			})
		}
	} else {
		specerr.Append(&s.Err, &specerr.ResourceError{
			Context: ctx, ResourceID: s.resourceID, Err: err,
		})
	}
}

// ParseRGBA parses s as a RGBA color.
func ParseRGBA(s string) (c color.RGBA, err error) {
	var errInvalidFormat = errors.New("go3mf: invalid color format")

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
	if c.A == 255 {
		return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}
