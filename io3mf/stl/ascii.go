package stl

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/mesh"
)

// asciiDecoder can create a Mesh from a Read stream that is feeded with a ASCII STL.
type asciiDecoder struct {
	r     io.Reader
	units float32
}

func (d *asciiDecoder) decode() (*mesh.Mesh, error) {
	newMesh := mesh.NewMesh()
	newMesh.StartCreation(mesh.CreationOptions{CalculateConnectivity: true})
	defer newMesh.EndCreation()
	position := 0
	var nodes [3]uint32
	scanner := bufio.NewScanner(d.r)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 4 && fields[0] == "vertex" {
			var f [3]float64
			f[0], _ = strconv.ParseFloat(fields[1], 32)
			f[1], _ = strconv.ParseFloat(fields[2], 32)
			f[2], _ = strconv.ParseFloat(fields[3], 32)
			nodes[position] = newMesh.AddNode(mgl32.Vec3{float32(f[0]), float32(f[1]), float32(f[2])})
			position++

			if position == 3 {
				position = 0
				newMesh.AddFace(nodes[0], nodes[1], nodes[2])
			}
		}
	}

	return newMesh, scanner.Err()
}

type asciiEncoder struct {
	w io.Writer
}

const pstr = "solid\nfacet normal %f %f %f\nouter loop\nvertex %f %f %f\nvertex %f %f %f\nvertex %f %f %f\nendloop\nendfacet\nendsolid\n"

func (e *asciiEncoder) encode(m *mesh.Mesh) error {
	for i := range m.Faces {
		node1, node2, node3 := m.FaceNodes(uint32(i))
		n := m.FaceNormal(uint32(i))
		n1, n2, n3 := node1.Position, node2.Position, node3.Position
		_, err := io.WriteString(e.w, fmt.Sprintf(pstr, n.X(), n.Y(), n.Z(), n1.X(), n1.Y(), n1.Z(), n2.X(), n2.Y(), n2.Z(), n3.X(), n3.Y(), n3.Z()))

		if err != nil {
			return err
		}
	}

	return nil
}
