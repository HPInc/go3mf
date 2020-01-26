package stl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/qmuntal/go3mf/geo"
)

// asciiDecoder can create a Model from a Read stream that is feeded with a ASCII STL.
type asciiDecoder struct {
	r     io.Reader
	units float32
}

func (d *asciiDecoder) decode(ctx context.Context, m *geo.Mesh) (err error) {
	mb := geo.NewMeshBuilder(m)
	position := 0
	nextFaceCheck := checkEveryFaces
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
			nodes[position] = mb.AddNode(geo.Point3D{float32(f[0]), float32(f[1]), float32(f[2])})
			position++

			if position == 3 {
				position = 0
				m.Faces = append(m.Faces, geo.Face{
					NodeIndices: [3]uint32{nodes[0], nodes[1], nodes[2]},
				})
				if len(m.Faces) > nextFaceCheck {
					select {
					case <-ctx.Done():
						err = ctx.Err()
						break
					default: // Default is must to avoid blocking
					}
					nextFaceCheck += checkEveryFaces
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

type asciiEncoder struct {
	w io.Writer
}

const pstr = "solid\nfacet normal %f %f %f\nouter loop\nvertex %f %f %f\nvertex %f %f %f\nvertex %f %f %f\nendloop\nendfacet\nendsolid\n"

func (e *asciiEncoder) encode(m *geo.Mesh) error {
	for _, f := range m.Faces {
		n1, n2, n3 := m.Nodes[f.NodeIndices[0]], m.Nodes[f.NodeIndices[1]], m.Nodes[f.NodeIndices[2]]
		n := faceNormal(n1, n2, n3)
		_, err := io.WriteString(e.w, fmt.Sprintf(pstr,
			n.X(), n.Y(), n.Z(),
			n1.X(), n1.Y(), n1.Z(),
			n2.X(), n2.Y(), n2.Z(),
			n3.X(), n3.Y(), n3.Z(),
		))

		if err != nil {
			return err
		}
	}

	return nil
}
