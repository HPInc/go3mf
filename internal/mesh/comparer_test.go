package mesh

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func Test_comparer_CompareGeometry(t *testing.T) {
	type args struct {
		m1 *Mesh
		m2 *Mesh
	}
	tests := []struct {
		name string
		c    comparer
		args args
		want bool
	}{
		{"base", comparer{}, args{NewMesh(), NewMesh()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.CompareGeometry(tt.args.m1, tt.args.m2); got != tt.want {
				t.Errorf("comparer.CompareGeometry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_comparer_fastCheck(t *testing.T) {
	msh := NewMesh()
	type args struct {
		m1 *Mesh
		m2 *Mesh
	}
	tests := []struct {
		name string
		c    comparer
		args args
		want bool
	}{
		{"nils", comparer{}, args{nil, nil}, false},
		{"nil1", comparer{}, args{nil, NewMesh()}, false},
		{"nil2", comparer{}, args{NewMesh(), nil}, false},
		{"nodes", comparer{}, args{&Mesh{nodeStructure: nodeStructure{nodes: make([]Node, 2)}}, NewMesh()}, false},
		{"faces", comparer{}, args{&Mesh{faceStructure: faceStructure{faces: make([]Face, 2)}}, NewMesh()}, false},
		{"beams", comparer{}, args{&Mesh{beamLattice: beamLattice{beams: make([]Beam, 2)}}, NewMesh()}, false},
		{"samepinter", comparer{}, args{msh, msh}, true},
		{"same", comparer{}, args{NewMesh(), NewMesh()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.fastCheck(tt.args.m1, tt.args.m2); got != tt.want {
				t.Errorf("comparer.fastCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_comparer_compareNodes(t *testing.T) {
	type args struct {
		m1 *Mesh
		m2 *Mesh
	}
	tests := []struct {
		name string
		c    comparer
		args args
		want bool
	}{
		{"diff", comparer{}, args{
			&Mesh{nodeStructure: nodeStructure{nodes: []Node{Node{Position: mgl32.Vec3{1.0, 2.5, 3.33}}}}},
			&Mesh{nodeStructure: nodeStructure{nodes: []Node{Node{Position: mgl32.Vec3{1.0, 3.5, 3.33}}}}},
		}, false},
		{"same", comparer{}, args{
			&Mesh{nodeStructure: nodeStructure{nodes: []Node{Node{Position: mgl32.Vec3{1.0, 2.5, 3.33}}}}},
			&Mesh{nodeStructure: nodeStructure{nodes: []Node{Node{Position: mgl32.Vec3{1.0, 2.5, 3.33}}}}},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.compareNodes(tt.args.m1, tt.args.m2); got != tt.want {
				t.Errorf("comparer.compareNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_comparer_compareFaces(t *testing.T) {
	type args struct {
		m1 *Mesh
		m2 *Mesh
	}
	tests := []struct {
		name string
		c    comparer
		args args
		want bool
	}{
		{"diff", comparer{}, args{
			&Mesh{faceStructure: faceStructure{faces: []Face{Face{NodeIndices: [3]uint32{0, 1, 2}}}}},
			&Mesh{faceStructure: faceStructure{faces: []Face{Face{NodeIndices: [3]uint32{0, 1, 3}}}}},
		}, false},
		{"same", comparer{}, args{
			&Mesh{faceStructure: faceStructure{faces: []Face{Face{NodeIndices: [3]uint32{0, 1, 2}}}}},
			&Mesh{faceStructure: faceStructure{faces: []Face{Face{NodeIndices: [3]uint32{0, 1, 2}}}}},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.compareFaces(tt.args.m1, tt.args.m2); got != tt.want {
				t.Errorf("comparer.compareFaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_comparer_compareBeams(t *testing.T) {
	type args struct {
		m1 *Mesh
		m2 *Mesh
	}
	tests := []struct {
		name string
		c    comparer
		args args
		want bool
	}{
		{"diff", comparer{}, args{
			&Mesh{beamLattice: beamLattice{beams: []Beam{{NodeIndices: [2]uint32{0, 1}}}}},
			&Mesh{beamLattice: beamLattice{beams: []Beam{{NodeIndices: [2]uint32{0, 2}}}}},
		}, false},
		{"same", comparer{}, args{
			&Mesh{beamLattice: beamLattice{beams: []Beam{{NodeIndices: [2]uint32{0, 1}}}}},
			&Mesh{beamLattice: beamLattice{beams: []Beam{{NodeIndices: [2]uint32{0, 1}}}}},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.compareBeams(tt.args.m1, tt.args.m2); got != tt.want {
				t.Errorf("comparer.compareBeams() = %v, want %v", got, tt.want)
			}
		})
	}
}
