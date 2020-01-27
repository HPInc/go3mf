package stl

import (
	"io"

	"github.com/qmuntal/go3mf"
)

// EncodingType is the type of encoding used in the file.
type EncodingType int

const (
	// Binary when the STL is encoded as a binary file.
	Binary EncodingType = iota
	// ASCII when the STL is encoded as an ASCII file.
	ASCII
)

// Encoder can encode a mesh as a binary or an ASCII file.
type Encoder struct {
	w            io.Writer
	encodingType EncodingType
}

// NewEncoder creates a new binary encoder.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:            w,
		encodingType: Binary,
	}
}

// NewEncoderType creates a new encoder of the desired type..
func NewEncoderType(w io.Writer, encodingType EncodingType) *Encoder {
	return &Encoder{
		w:            w,
		encodingType: encodingType,
	}
}

// Encode encodes a mesh to the writer.
func (e *Encoder) Encode(m *go3mf.Mesh) error {
	switch e.encodingType {
	case ASCII:
		encoder := asciiEncoder{w: e.w}
		return encoder.encode(m)
	default:
		encoder := binaryEncoder{w: e.w}
		return encoder.encode(m)
	}
}

func faceNormal(n1, n2, n3 go3mf.Point3D) go3mf.Point3D {
	return n2.Sub(n1).Cross(n3.Sub(n1)).Normalize()
}
