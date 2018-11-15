package stl

import (
	"io"
	"bufio"
	"strings"
	"strconv"
	"github.com/qmuntal/go3mf/internal/mesh"
	"github.com/go-gl/mathgl/mgl32"
)

// asciiDecoder can create a Mesh from a Read stream that is feeded with a ASCII STL.
type asciiDecoder struct {
	r 	io.Reader
	units	float32
}

func (d* asciiDecoder) decode() (*mesh.Mesh, error) {
	newMesh := mesh.NewMesh()
	err := newMesh.StartCreation(d.units)
	defer newMesh.EndCreation()
	if err != nil {
		return nil, err
	}

	position := 0
	var nodes [3]*mesh.Node
	scanner := bufio.NewScanner(d.r)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 4 && fields[0] == "vertex" {
			var f[3] float64
			f[0], _ = strconv.ParseFloat(fields[1], 32)
			f[1], _ = strconv.ParseFloat(fields[2], 32)
			f[2], _ = strconv.ParseFloat(fields[3], 32)

			nodes[position] = newMesh.AddNode(mgl32.Vec3{float32(f[0]), float32(f[1]), float32(f[2])})
			position++

			if (position == 3) {
				position = 0 
				newMesh.AddFace(nodes[0], nodes[1], nodes[2])
			}
		}
	}

	return newMesh, err
}

type asciiEncoder struct {
	w io.Writer
}

/*func (e *asciiEncoder) encode(m *mesh.Mesh) error {
	


	return nil
}*/