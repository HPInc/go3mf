package geometry

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewVectorTree(t *testing.T) {
	tests := []struct {
		name string
		want *VectorTree
	}{
		{"new", &VectorTree{0.001, map[Vec3I]uint32{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVectorTree(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVectorTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorTree_Units(t *testing.T) {
	p := NewVectorTree()
	if got := p.Units(); got != VectorDefaultUnits {
		t.Errorf("VectorTree.Units() = %v, want %v", got, VectorDefaultUnits)
	}
	p.units = 0
	if got := p.Units(); got != VectorDefaultUnits {
		t.Errorf("VectorTree.Units() = %v, want %v", got, VectorDefaultUnits)
	}
	p.SetUnits(1.0)
	if got := p.Units(); got != 1.0 {
		t.Errorf("VectorTree.Units() = %v, want %v", got, VectorDefaultUnits)
	}
}

func TestVectorTree_SetUnits(t *testing.T) {
	p := NewVectorTree()
	p.AddVector(mgl32.Vec3{}, 1)
	type args struct {
		units float32
	}
	tests := []struct {
		name    string
		t       *VectorTree
		args    args
		wantErr bool
	}{
		{"zero", NewVectorTree(), args{0.0}, false},
		{"minunitsfail", NewVectorTree(), args{0.000009}, true},
		{"minunits", NewVectorTree(), args{VectorMinUnits}, false},
		{"one", NewVectorTree(), args{1.0}, false},
		{"maxunits", NewVectorTree(), args{VectorMaxUnits}, false},
		{"maxunitsfail", NewVectorTree(), args{1001.0}, true},
		{"notempty", p, args{1.0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.SetUnits(tt.args.units); (err != nil) != tt.wantErr {
				t.Errorf("VectorTree.SetUnits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want := tt.args.units
			if want == 0 || tt.wantErr {
				want = VectorDefaultUnits
			}
			if got := tt.t.Units(); got != want {
				t.Errorf("VectorTree.Units() = %v, want %v", got, want)
			}
		})
	}
}

func TestVectorTree_AddFindVector(t *testing.T) {
	p := NewVectorTree()
	type args struct {
		vec   mgl32.Vec3
		value uint32
	}
	tests := []struct {
		name string
		t    *VectorTree
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
			t.Error("VectorTree.AddMatch() haven't added the match")
			return
		}
		if got != tt.args.value {
			t.Errorf("VectorTree.FindVector() = %v, want %v", got, tt.args.value)
		}
	}
}

func TestVectorTree_RemoveVector(t *testing.T) {
	p := NewVectorTree()
	p.AddVector(mgl32.Vec3{1, 2, 5.3}, 1)
	type args struct {
		vec mgl32.Vec3
	}
	tests := []struct {
		name string
		t    *VectorTree
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

func Test_newVec3IFromVec3(t *testing.T) {
	type args struct {
		vec   mgl32.Vec3
		units float32
	}
	tests := []struct {
		name string
		args args
		want Vec3I
	}{
		{"1", args{mgl32.Vec3{1.2, 2.3, 3.4}, 1.0}, Vec3I{1, 2, 3}},
		{"2", args{mgl32.Vec3{1.2, 2.3, 3.4}, 0.001}, Vec3I{1200, 2299, 3400}},
		{"3", args{mgl32.Vec3{1.2, 2.3, 3.4}, 1000.0}, Vec3I{0, 0, 0}},
		{"4", args{mgl32.Vec3{1000.2, 2000.3, 3000.4}, 1000.0}, Vec3I{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newVec3IFromVec3(tt.args.vec, tt.args.units); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newVec3IFromVec3() = %v, want %v", got, tt.want)
			}
		})
	}
}
