package geo

import (
	"math"
	"reflect"
	"testing"
)

func Test_nodeStructure_AddNode(t *testing.T) {
	pos := Point3D{1.0, 2.0, 3.0}
	existingStruct := &nodeStructure{vectorTree: newVectorTree()}
	existingStruct.AddNode(pos)
	type args struct {
		position Point3D
	}
	tests := []struct {
		name string
		n    *nodeStructure
		args args
		want uint32
	}{
		{"existing", existingStruct, args{pos}, 0},
		{"base", &nodeStructure{Nodes: []Point3D{{}}}, args{Point3D{1.0, 2.0, 3.0}}, 1},
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

func Test_newvec3IFromVec3(t *testing.T) {
	type args struct {
		vec Point3D
	}
	tests := []struct {
		name string
		args args
		want vec3I
	}{
		{"base", args{Point3D{1.2, 2.3, 3.4}}, vec3I{1200000, 2300000, 3400000}},
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
		vec   Point3D
		value uint32
	}
	tests := []struct {
		name string
		t    *vectorTree
		args args
	}{
		{"new", p, args{Point3D{10000.3, 20000.2, 1}, 2}},
		{"old", p, args{Point3D{10000.3, 20000.2, 1}, 4}},
		{"new2", p, args{Point3D{2, 1, 3.4}, 5}},
		{"old2", p, args{Point3D{2, 1, 3.4}, 1}},
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
	p.AddVector(Point3D{1, 2, 5.3}, 1)
	type args struct {
		vec Point3D
	}
	tests := []struct {
		name string
		t    *vectorTree
		args args
	}{
		{"nil", p, args{Point3D{2, 3, 4}}},
		{"old", p, args{Point3D{1, 2, 5.3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.RemoveVector(tt.args.vec)
		})
	}
}

func TestPoint3D_Add(t *testing.T) {
	type args struct {
		v2 Point3D
	}
	tests := []struct {
		name string
		v1   Point3D
		args args
		want Point3D
	}{
		{"base", Point3D{1.0, 2.5, 1.1}, args{Point3D{0.0, 1.0, 9.9}}, Point3D{1.0, 3.5, 11.0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Add(tt.args.v2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Point3D.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint3D_Sub(t *testing.T) {
	type args struct {
		v2 Point3D
	}
	tests := []struct {
		name string
		v1   Point3D
		args args
		want Point3D
	}{
		{"base", Point3D{1.0, 2.5, 1.0}, args{Point3D{0.0, 1.0, 9.9}}, Point3D{1.0, 1.5, -8.9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Sub(tt.args.v2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Point3D.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint3D_Len(t *testing.T) {
	tests := []struct {
		name string
		v1   Point3D
		want float32
	}{
		{"base", Point3D{2, -5, 4}, 6.708203932499},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Len(); math.Abs(float64(got-tt.want)) > 1e-6 {
				t.Errorf("Point3D.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint3D_Normalize(t *testing.T) {
	tests := []struct {
		name string
		v1   Point3D
		want Point3D
	}{
		{"x", Point3D{2, 0, 0}, Point3D{1, 0, 0}},
		{"y", Point3D{0, 3, 0}, Point3D{0, 1, 0}},
		{"z", Point3D{0, 0, 4}, Point3D{0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Normalize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Point3D.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint3D_Cross(t *testing.T) {
	type args struct {
		v2 Point3D
	}
	tests := []struct {
		name string
		v1   Point3D
		args args
		want Point3D
	}{
		{"base", Point3D{1, 2, 3}, args{Point3D{10, 11, 12}}, Point3D{-9, 18, -9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Cross(tt.args.v2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Point3D.Cross() = %v, want %v", got, tt.want)
			}
		})
	}
}
