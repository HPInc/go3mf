// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package stl

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/hpinc/go3mf"
)

// asciiDecoder can create a Model from a Read stream that is feeded with a ASCII STL.
type asciiDecoder struct {
	r     io.Reader
	units float32
}

func (d *asciiDecoder) decode(ctx context.Context, m *go3mf.Mesh) (err error) {
	mb := go3mf.NewMeshBuilder(m)
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
			nodes[position] = mb.AddVertex(go3mf.Point3D{float32(f[0]), float32(f[1]), float32(f[2])})
			position++

			if position == 3 {
				position = 0
				m.Triangles = append(m.Triangles, go3mf.Triangle{V1: nodes[0], V2: nodes[1], V3: nodes[2]})
				if len(m.Triangles) > nextFaceCheck {
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
