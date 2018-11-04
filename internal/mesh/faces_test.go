package mesh

import (
	"reflect"
	"testing"
)

func Test_faceStructure_clear(t *testing.T) {
	tests := []struct {
		name string
		f    *faceStructure
	}{
		{"base", &faceStructure{faces: make([]*Face, 2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.clear()
			if got := tt.f.FaceCount(); got != 0 {
				t.Errorf("faceStructure.clear() = %v, want %v", got, 0)
			}
		})
	}
}

func Test_faceStructure_FaceCount(t *testing.T) {
	tests := []struct {
		name string
		f    *faceStructure
		want uint32
	}{
		{"zero", new(faceStructure), 0},
		{"base", &faceStructure{faces: make([]*Face, 2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.FaceCount(); got != tt.want {
				t.Errorf("faceStructure.FaceCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_faceStructure_Face(t *testing.T) {
	f := new(faceStructure)
	f.faces = append(f.faces, new(Face))
	f.faces = append(f.faces, new(Face))
	type args struct {
		index uint32
	}
	tests := []struct {
		name string
		f    *faceStructure
		args args
		want *Face
	}{
		{"zero", f, args{0}, f.faces[0]},
		{"one", f, args{1}, f.faces[1]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Face(tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("faceStructure.Face() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_faceStructure_AddFace(t *testing.T) {
	type args struct {
		node1 *Node
		node2 *Node
		node3 *Node
	}
	tests := []struct {
		name    string
		f       *faceStructure
		args    args
		want    *Face
		wantErr bool
	}{
		{"max", &faceStructure{maxFaceCount: 1, faces: []*Face{new(Face)}}, args{new(Node), new(Node), new(Node)}, nil, true},
		{"base", &faceStructure{faces: []*Face{new(Face)}}, args{new(Node), new(Node), new(Node)}, &Face{Index: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.AddFace(tt.args.node1, tt.args.node2, tt.args.node3)
			if (err != nil) != tt.wantErr {
				t.Errorf("faceStructure.AddFace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		// TODO: Add test cases.
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
		other    mergeableFaces
		newNodes []*Node
	}
	tests := []struct {
		name    string
		f       *faceStructure
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.merge(tt.args.other, tt.args.newNodes); (err != nil) != tt.wantErr {
				t.Errorf("faceStructure.merge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
