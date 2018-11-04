package mesh

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	gomock "github.com/golang/mock/gomock"
)

func Test_nodeStructure_clear(t *testing.T) {
	tests := []struct {
		name string
		n    *nodeStructure
	}{
		{"base", &nodeStructure{nodes: make([]*Node, 2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.clear()
			if got := tt.n.NodeCount(); got != 0 {
				t.Errorf("nodeStructure.NodeCount() = %v, want %v", got, 0)
			}
		})
	}
}

func Test_nodeStructure_NodeCount(t *testing.T) {
	tests := []struct {
		name string
		n    *nodeStructure
		want uint32
	}{
		{"zero", new(nodeStructure), 0},
		{"base", &nodeStructure{nodes: make([]*Node, 2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.NodeCount(); got != tt.want {
				t.Errorf("nodeStructure.NodeCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeStructure_Node(t *testing.T) {
	n := new(nodeStructure)
	n.nodes = append(n.nodes, new(Node))
	n.nodes = append(n.nodes, new(Node))
	type args struct {
		index uint32
	}
	tests := []struct {
		name string
		n    *nodeStructure
		args args
		want *Node
	}{
		{"zero", n, args{0}, n.nodes[0]},
		{"one", n, args{1}, n.nodes[1]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Node(tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeStructure.Node() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeStructure_AddNode(t *testing.T) {
	type args struct {
		position mgl32.Vec3
	}
	tests := []struct {
		name    string
		n       *nodeStructure
		args    args
		want    *Node
		wantErr bool
	}{
		{"max", &nodeStructure{maxNodeCount: 1, nodes: []*Node{new(Node)}}, args{mgl32.Vec3{}}, nil, true},
		{"base", &nodeStructure{nodes: []*Node{new(Node)}}, args{mgl32.Vec3{1.0, 2.0, 3.0}}, &Node{
			Index:    1,
			Position: mgl32.Vec3{1.0, 2.0, 3.0},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AddNode(tt.args.position)
			if (err != nil) != tt.wantErr {
				t.Errorf("nodeStructure.AddNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeStructure.AddNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeStructure_checkSanity(t *testing.T) {
	tests := []struct {
		name string
		n    *nodeStructure
		want bool
	}{
		{"max", &nodeStructure{maxNodeCount: 1, nodes: []*Node{new(Node), new(Node)}}, false},
		{"badindex", &nodeStructure{nodes: []*Node{new(Node), &Node{Index: 2}}}, false},
		{"good", &nodeStructure{nodes: []*Node{new(Node), &Node{Index: 1}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.checkSanity(); got != tt.want {
				t.Errorf("nodeStructure.checkSanity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeStructure_merge(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMesh := NewMockMergeableMesh(mockCtrl)
	type args struct {
		other  mergeableNodes
		matrix mgl32.Mat4
	}
	tests := []struct {
		name    string
		n       *nodeStructure
		args    args
		want    []*Node
		wantErr bool
		times   uint32
	}{
		{"zero", new(nodeStructure), args{mockMesh, mgl32.Ident4()}, make([]*Node, 0), false, 0},
		{"err", &nodeStructure{maxNodeCount: 1, nodes: []*Node{new(Node)}}, args{mockMesh, mgl32.Ident4()}, nil, true, 1},
		{"merged", new(nodeStructure), args{mockMesh, mgl32.Translate3D(1.0, 2.0, 3.0)}, []*Node{
			&Node{Index: 0, Position: mgl32.Vec3{1.0, 2.0, 3.0}},
			&Node{Index: 1, Position: mgl32.Vec3{1.0, 2.0, 3.0}}},
			false, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMesh.EXPECT().NodeCount().Return(tt.times)
			mockMesh.EXPECT().Node(gomock.Any()).Return(new(Node)).Times(int(tt.times))
			got, err := tt.n.merge(tt.args.other, tt.args.matrix)
			if (err != nil) != tt.wantErr {
				t.Errorf("nodeStructure.merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeStructure.merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
