package obj

import (
	"bytes"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf/mesh"
)

func TestDecoderDecode(t *testing.T) {
	tests := []struct {
		name    string
		d       *Decoder
		want    *mesh.Mesh
		wantErr bool
	}{
		{"base", NewDecoder(bytes.NewBufferString(objMesh())), createMeshTriangle(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Decode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := deep.Equal(got, tt.want); diff != nil {
					t.Errorf("Decoder.Decode() = %v", diff)
					return
				}
			}
		})
	}
}

func createMeshTriangle() *mesh.Mesh {
	m := mesh.NewMesh()
	n1 := m.AddNode(mesh.Node{-20.0, -20.0, 0.0})
	n2 := m.AddNode(mesh.Node{20.0, -20.0, 0.0})
	n3 := m.AddNode(mesh.Node{0.0019989014, 0.0019989014, 39.998})
	n4 := m.AddNode(mesh.Node{-20.0, 20.0, 0.0})
	n5 := m.AddNode(mesh.Node{20.0, 20.0, 0.0})
	m.AddFace(n1, n2, n3)
	m.AddFace(n4, n2, n1)
	m.AddFace(n1, n3, n4)
	m.AddFace(n2, n5, n3)
	m.AddFace(n5, n4, n3)
	m.AddFace(n5, n2, n4)
	return m
}

func objMesh() string {
	return `
# Exported from 3D Builder
mtllib Piràmide.mtl

o Object.1
v -40.000000 -40.000000 0.000000 2.0
v 20.000000 -20.000000 0.000000 255.0 155.0 100.0
v 0.0019989014 0.0019989014 39.998 83.0 98.0 100.0
v -20.000000 20.000000 0.000000
v 20.000000 20.000000 0.000000 255.0 155.0 100.0

vt 1.0000 0.0000 0.0000

usemtl Yellow_0
f 1 2 3
f 4 2 1
f 1 3 4
f 2/1 5/1 3/1
f -1 4 3
f -1 2 4
`
}
