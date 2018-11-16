package stl

import (
	"io"
	"bufio"
	"strings"
	"strconv"
	"fmt"
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

func (e *asciiEncoder) encode(m *mesh.Mesh) {
	faceCount := m.FaceCount()

	for i := 0; i < int(faceCount); i++ {				
		// First we start by calculating the normal
		normal := m.FaceNormal(uint32(i))

		// Secondly we catch the vertexes
		node1, node2, node3 := m.FaceCoordinates(uint32(i))

		// Lastly we print all the components
		io.WriteString(e.w, fmt.Sprintf("facet normal %f %f %f\nouter loop\n", normal.X(), normal.Y(), normal.Z()))
		io.WriteString(e.w, fmt.Sprintf("vertex %f %f %f\n", node1.X(), node1.Y(), node1.Z()))
		io.WriteString(e.w, fmt.Sprintf("vertex %f %f %f\n", node2.X(), node2.Y(), node2.Z()))
		io.WriteString(e.w, fmt.Sprintf("vertex %f %f %f\n", node3.X(), node3.Y(), node3.Z()))
		io.WriteString(e.w, "endloop\nendfacet\n")
	}
}