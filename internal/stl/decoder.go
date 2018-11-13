package stl

import (
	"bufio"
	"io"

	"github.com/qmuntal/go3mf/internal/mesh"
)

func DecodeUnits(r io.Reader, units float32) (*mesh.Mesh, error) {
	b := bufio.NewReader(r)
	d := binaryDecoder{r: b}
	return d.decode()
}
