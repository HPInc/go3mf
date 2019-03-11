package mesh

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func Test_nodeStructure_clear(t *testing.T) {
	tests := []struct {
		name string
		n    *nodeStructure
	}{
		{"base", &nodeStructure{Nodes: make([]Node, 2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.clear()
			if got := len(tt.n.Nodes); got != 0 {
				t.Errorf("nodeStructure.clear() = %v, want %v", got, 0)
			}
		})
	}
}

func Test_nodeStructure_AddNode(t *testing.T) {
	pos := mgl32.Vec3{1.0, 2.0, 3.0}
	existingStruct := &nodeStructure{vectorTree: newVectorTree()}
	existingStruct.AddNode(pos)
	type args struct {
		position mgl32.Vec3
	}
	tests := []struct {
		name      string
		n         *nodeStructure
		args      args
		want      *Node
	}{
		{"existing", existingStruct, args{pos}, &Node{Index: 0, Position: pos}},
		{"base", &nodeStructure{Nodes: []Node{{}}}, args{mgl32.Vec3{1.0, 2.0, 3.0}}, &Node{
			Index:    1,
			Position: mgl32.Vec3{1.0, 2.0, 3.0},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		{"max", &nodeStructure{maxNodeCount: 1, Nodes: []Node{{}, {}}}, false},
		{"badindex", &nodeStructure{Nodes: []Node{{}, {Index: 2}}}, false},
		{"good", &nodeStructure{Nodes: []Node{{}, {Index: 1}}}, true},
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
	type args struct {
		matrix mgl32.Mat4
	}
	tests := []struct {
		name  string
		n     *nodeStructure
		args  args
		want  []uint32
		times int
	}{
		{"zero", new(nodeStructure), args{mgl32.Ident4()}, make([]uint32, 0), 0},
		{"merged", new(nodeStructure), args{mgl32.Translate3D(1.0, 2.0, 3.0)}, []uint32{0, 1}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Node{}
			mockMesh := NewMesh()
			for i := 0; i < tt.times; i++ {
				mockMesh.Nodes = append(mockMesh.Nodes, node)
			}
			got := tt.n.merge(&mockMesh.nodeStructure, tt.args.matrix)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeStructure.merge() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func Test_newvec3IFromVec3(t *testing.T) {
	type args struct {
		vec mgl32.Vec3
	}
	tests := []struct {
		name string
		args args
		want vec3I
	}{
		{"base", args{mgl32.Vec3{1.2, 2.3, 3.4}}, vec3I{1200000, 2300000, 3400000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newvec3IFromVec3(tt.args.vec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newvec3IFromVec3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newVectorTree(t *testing.T) {
	tests := []struct {
		name string
		want *vectorTree
	}{
		{"new", &vectorTree{map[vec3I]uint32{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newVectorTree(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newVectorTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vectorTree_AddFindVector(t *testing.T) {
	p := newVectorTree()
	type args struct {
		vec   mgl32.Vec3
		value uint32
	}
	tests := []struct {
		name string
		t    *vectorTree
		args args
	}{
		{"new", p, args{mgl32.Vec3{10000.3, 20000.2, 1}, 2}},
		{"old", p, args{mgl32.Vec3{10000.3, 20000.2, 1}, 4}},
		{"new2", p, args{mgl32.Vec3{2, 1, 3.4}, 5}},
		{"old2", p, args{mgl32.Vec3{2, 1, 3.4}, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.AddVector(tt.args.vec, tt.args.value)
		})
		got, ok := p.FindVector(tt.args.vec)
		if !ok {
			t.Error("vectorTree.AddMatch() haven't added the match")
			return
		}
		if got != tt.args.value {
			t.Errorf("vectorTree.FindVector() = %v, want %v", got, tt.args.value)
		}
	}
}

func Test_vectorTree_RemoveVector(t *testing.T) {
	p := newVectorTree()
	p.AddVector(mgl32.Vec3{1, 2, 5.3}, 1)
	type args struct {
		vec mgl32.Vec3
	}
	tests := []struct {
		name string
		t    *vectorTree
		args args
	}{
		{"nil", p, args{mgl32.Vec3{2, 3, 4}}},
		{"old", p, args{mgl32.Vec3{1, 2, 5.3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.RemoveVector(tt.args.vec)
		})
	}
}
