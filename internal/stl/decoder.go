package stl

import (
	"bufio"
	"io"
	"strings"

	"github.com/qmuntal/go3mf/internal/mesh"
	"golang.org/x/exp/utf8string"
)

const sizeOfHeader = 300 // minimum size of a closed mesh in binary is 384 bytes, corresponding to a triangle.

// Decoder can decode an stl to a mesh.
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
func (d *Decoder) Decode() (*mesh.Mesh, error) {
	b := bufio.NewReader(d.r)
	isASCII, err := d.isASCII(b)
	if err != nil {
		return nil, err
	}
	if isASCII {
		decoder := asciiDecoder{r: b}
		if m, err := decoder.decode(); err == nil {
			return m, nil
		}
	}
	decoder := binaryDecoder{r: b}
	return decoder.decode()
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
	return strings.HasPrefix(header, "solid") && utf8string.NewString(header).IsASCII(), nil
}
