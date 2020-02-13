package stl

import (
	"bufio"
	"context"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/qmuntal/go3mf"
)

var checkEveryFaces = 1000

const sizeOfHeader = 300 // minimum size of a closed mesh in binary is 384 bytes, corresponding to a triangle.

// Decoder can decode a stl.
// It supports automatic detection of binary or ascii stl encoding.
type Decoder struct {
	r io.Reader
}

// NewDecoder creates a new decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

// Decode creates a mesh from a read stream.
func (d *Decoder) Decode(m *go3mf.Model) error {
	return d.DecodeContext(context.Background(), m)
}

// DecodeContext creates a mesh from a read stream.
func (d *Decoder) DecodeContext(ctx context.Context, m *go3mf.Model) error {
	b := bufio.NewReader(d.r)
	isASCII, err := d.isASCII(b)
	if err != nil {
		return err
	}
	newMesh := go3mf.NewMeshResource()
	if isASCII {
		decoder := asciiDecoder{r: b}
		err = decoder.decode(ctx, newMesh.Mesh)
	} else {
		decoder := binaryDecoder{r: b}
		err = decoder.decode(ctx, newMesh.Mesh)
	}
	if err == nil {
		newMesh.ModelPath = m.Path
		if len(m.Resources) == 0 {
			m.Resources = []*go3mf.Resources{new(go3mf.Resources)}
		}
		rs := m.Resources[len(m.Resources)-1]
		newMesh.ID = rs.UnusedID()
		rs.Objects = append(rs.Objects, newMesh)
	}
	return err
}

func (d *Decoder) isASCII(r *bufio.Reader) (bool, error) {
	var header string
	for {
		buff, err := r.Peek(sizeOfHeader)
		if err == io.EOF {
			return false, err
		}
		if len(buff) >= sizeOfHeader {
			header = strings.ToLower(string(buff))
			break
		}
	}
	return strings.HasPrefix(header, "solid") && isASCII(header), nil
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}
