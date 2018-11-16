package mesh

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	gomock "github.com/golang/mock/gomock"
	"github.com/qmuntal/go3mf/internal/geometry"
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
				t.Errorf("nodeStructure.clear() = %v, want %v", got, 0)
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
	pos := mgl32.Vec3{1.0, 2.0, 3.0}
	existingStruct := &nodeStructure{vectorTree: geometry.NewVectorTree()}
	existingStruct.AddNode(pos)
	type args struct {
		position mgl32.Vec3
	}
	tests := []struct {
		name      string
		n         *nodeStructure
		args      args
		want      *Node
		wantPanic bool
	}{
		{"existing", existingStruct, args{pos}, &Node{Index: 0, Position: pos}, false},
		{"max", &nodeStructure{maxNodeCount: 1, nodes: []*Node{new(Node)}}, args{mgl32.Vec3{}}, nil, true},
		{"base", &nodeStructure{nodes: []*Node{new(Node)}}, args{mgl32.Vec3{1.0, 2.0, 3.0}}, &Node{
			Index:    1,
			Position: mgl32.Vec3{1.0, 2.0, 3.0},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); tt.wantPanic && r == nil {
					t.Error("nodeStructure.AddNode() want panic")
				}
			}()
			got := tt.n.AddNode(tt.args.position)
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
		{"badindex", &nodeStructure{nodes: []*Node{new(Node), {Index: 2}}}, false},
		{"good", &nodeStructure{nodes: []*Node{new(Node), {Index: 1}}}, true},
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
		name  string
		n     *nodeStructure
		args  args
		want  []*Node
		times uint32
	}{
		{"zero", new(nodeStructure), args{mockMesh, mgl32.Ident4()}, make([]*Node, 0), 0},
		{"merged", new(nodeStructure), args{mockMesh, mgl32.Translate3D(1.0, 2.0, 3.0)}, []*Node{
			{Index: 0, Position: mgl32.Vec3{1.0, 2.0, 3.0}},
			{Index: 1, Position: mgl32.Vec3{1.0, 2.0, 3.0}}},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMesh.EXPECT().NodeCount().Return(tt.times)
			mockMesh.EXPECT().Node(gomock.Any()).Return(new(Node)).Times(int(tt.times))
			got := tt.n.merge(tt.args.other, tt.args.matrix)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeStructure.merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
