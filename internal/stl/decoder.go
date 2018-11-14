package stl

import (
	"bufio"
	"io"
	"strings"

	"github.com/qmuntal/go3mf/internal/mesh"
	"golang.org/x/exp/utf8string"
)

const sizeOfHeader = 300 // minimum size of a closed mesh in binary is 384 bytes, corresponding to a triangle

func DecodeUnits(r io.Reader, units float32) (*mesh.Mesh, error) {
	b := bufio.NewReader(r)
	ascii, err := isASCII(b)
	if err != nil {
		return nil, err
	}
	if ascii {
		return nil, nil
	}
	d := binaryDecoder{r: b}
	return d.decode()
}

func isASCII(r *bufio.Reader) (bool, error) {
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
