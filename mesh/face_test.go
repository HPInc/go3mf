package mesh

import (
	"reflect"
	"testing"
)

func Test_faceStructure_AddFace(t *testing.T) {
	type args struct {
		node1 uint32
		node2 uint32
		node3 uint32
	}
	tests := []struct {
		name string
		f    *faceStructure
		args args
		want *Face
	}{
		{"base", &faceStructure{Faces: []Face{{}}}, args{0, 1, 2}, &Face{NodeIndices: [3]uint32{0, 1, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.AddFace(tt.args.node1, tt.args.node2, tt.args.node3)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("faceStructure.AddFace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_faceStructure_checkSanity(t *testing.T) {
	type args struct {
		nodeCount uint32
	}
	tests := []struct {
		name string
		f    *faceStructure
		args args
		want bool
	}{
		{"max", &faceStructure{maxFaceCount: 1, Faces: make([]Face, 2)}, args{1}, false},
		{"i0==i1", &faceStructure{Faces: []Face{{NodeIndices: [3]uint32{1, 1, 2}}}}, args{3}, false},
		{"i0==i2", &faceStructure{Faces: []Face{{NodeIndices: [3]uint32{1, 2, 1}}}}, args{3}, false},
		{"i1==i2", &faceStructure{Faces: []Face{{NodeIndices: [3]uint32{2, 1, 1}}}}, args{3}, false},
		{"i0big", &faceStructure{Faces: []Face{{NodeIndices: [3]uint32{3, 1, 2}}}}, args{3}, false},
		{"i1big", &faceStructure{Faces: []Face{{NodeIndices: [3]uint32{0, 3, 2}}}}, args{3}, false},
		{"i2big", &faceStructure{Faces: []Face{{NodeIndices: [3]uint32{0, 1, 3}}}}, args{3}, false},
		{"good", &faceStructure{Faces: []Face{{NodeIndices: [3]uint32{0, 1, 2}}}}, args{3}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.checkSanity(tt.args.nodeCount); got != tt.want {
				t.Errorf("faceStructure.checkSanity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_faceStructure_merge(t *testing.T) {
	type args struct {
		newNodes []uint32
	}
	tests := []struct {
		name  string
		f     *faceStructure
		args  args
		times int
	}{
		{"zero", new(faceStructure), args{make([]uint32, 0)}, 0},
		{"merged", new(faceStructure), args{[]uint32{0, 1, 2}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			face := Face{NodeIndices: [3]uint32{0, 1, 2}}
			mockMesh := new(Mesh)
			for i := 0; i < tt.times; i++ {
				mockMesh.Faces = append(mockMesh.Faces, face)
			}
			tt.f.merge(&mockMesh.faceStructure, tt.args.newNodes)
		})
	}
}
