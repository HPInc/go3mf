package mesh

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	meshinfo "github.com/qmuntal/go3mf/internal/meshinfo"
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
	node := new(Node)
	type args struct {
		node1 *Node
		node2 *Node
		node3 *Node
	}
	tests := []struct {
		name      string
		f         *faceStructure
		args      args
		want      *Face
		wantErr   bool
		wantPanic bool
	}{
		{"max", &faceStructure{maxFaceCount: 1, faces: make([]*Face, 1)}, args{new(Node), new(Node), new(Node)}, nil, false, true},
		{"duplicated0-1", new(faceStructure), args{node, node, new(Node)}, nil, true, false},
		{"duplicated0-2", new(faceStructure), args{node, new(Node), node}, nil, true, false},
		{"duplicated1-2", new(faceStructure), args{new(Node), node, node}, nil, true, false},
		{"base", &faceStructure{informationHandler: meshinfo.NewHandler(), faces: []*Face{new(Face)}}, args{new(Node), new(Node), new(Node)}, &Face{Index: 1}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); tt.wantPanic && r == nil {
					t.Error("faceStructure.AddFace() want panic")
				}
			}()
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
		{"max", &faceStructure{maxFaceCount: 1, faces: make([]*Face, 2)}, args{1}, false},
		{"i0==i1", &faceStructure{faces: []*Face{{NodeIndices: [3]uint32{1, 1, 2}}}}, args{3}, false},
		{"i0==i2", &faceStructure{faces: []*Face{{NodeIndices: [3]uint32{1, 2, 1}}}}, args{3}, false},
		{"i1==i2", &faceStructure{faces: []*Face{{NodeIndices: [3]uint32{2, 1, 1}}}}, args{3}, false},
		{"i0big", &faceStructure{faces: []*Face{{NodeIndices: [3]uint32{3, 1, 2}}}}, args{3}, false},
		{"i1big", &faceStructure{faces: []*Face{{NodeIndices: [3]uint32{0, 3, 2}}}}, args{3}, false},
		{"i2big", &faceStructure{faces: []*Face{{NodeIndices: [3]uint32{0, 1, 3}}}}, args{3}, false},
		{"good", &faceStructure{faces: []*Face{{NodeIndices: [3]uint32{0, 1, 2}}}}, args{3}, true},
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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMesh := NewMockMergeableMesh(mockCtrl)
	nodes := []*Node{{Index: 0}, {Index: 1}, {Index: 2}}
	type args struct {
		other    mergeableFaces
		newNodes []*Node
	}
	tests := []struct {
		name    string
		f       *faceStructure
		args    args
		wantErr bool
		times   uint32
	}{
		{"err", &faceStructure{maxFaceCount: 1, faces: make([]*Face, 1)}, args{mockMesh, nodes}, true, 1},
		{"zero", new(faceStructure), args{mockMesh, make([]*Node, 0)}, false, 0},
		{"merged", new(faceStructure), args{mockMesh, []*Node{{Index: 0}, {Index: 1}, {Index: 2}}}, false, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); tt.wantErr && r == nil {
					t.Error("faceStructure.merge() want panic")
				}
			}()
			mockMesh.EXPECT().FaceCount().Return(tt.times)
			mockMesh.EXPECT().InformationHandler().Return(meshinfo.NewHandler()).MaxTimes(int(tt.times))
			tt.f.informationHandler = meshinfo.NewHandler()
			face := &Face{NodeIndices: [3]uint32{0, 1, 2}}
			mockMesh.EXPECT().Face(gomock.Any()).Return(face).Times(int(tt.times))
			if err := tt.f.merge(tt.args.other, tt.args.newNodes); (err != nil) != tt.wantErr {
				t.Errorf("faceStructure.merge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
