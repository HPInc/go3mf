package model

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestSlice_BeginPolygon(t *testing.T) {
	s := new(Slice)
	tests := []struct {
		name string
		s    *Slice
		want int
	}{
		{"empty", s, 0},
		{"1", s, 1},
		{"2", s, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.BeginPolygon(); got != tt.want {
				t.Errorf("Slice.BeginPolygon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_AddVertex(t *testing.T) {
	s := new(Slice)
	type args struct {
		x float32
		y float32
	}
	tests := []struct {
		name string
		s    *Slice
		args args
		want int
	}{
		{"empty", s, args{1, 2}, 0},
		{"1", s, args{2, 3}, 1},
		{"2", s, args{4, 5}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.AddVertex(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Slice.AddVertex() = %v, want %v", got, tt.want)
				return
			}
			want := mgl32.Vec2{tt.args.x, tt.args.y}
			if !reflect.DeepEqual(tt.s.Vertices[tt.want], want) {
				t.Errorf("Slice.AddVertex() = %v, want %v", tt.s.Vertices[tt.want], want)
			}
		})
	}
}

func TestSlice_AddPolygonIndex(t *testing.T) {
	type args struct {
		polygonIndex int
		index        int
	}
	tests := []struct {
		name    string
		s       *Slice
		args    args
		wantErr bool
	}{
		{"emptyPolygon", new(Slice), args{0, 0}, true},
		{"emptyVertices", &Slice{Polygons: [][]int{{}}}, args{0, 0}, true},
		{"duplicated", &Slice{Polygons: [][]int{{0}}, Vertices: []mgl32.Vec2{{}}}, args{0, 0}, true},
		{"base", &Slice{Polygons: [][]int{{}}, Vertices: []mgl32.Vec2{{}}}, args{0, 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddPolygonIndex(tt.args.polygonIndex, tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("Slice.AddPolygonIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlice_AllPolygonsAreClosed(t *testing.T) {
	tests := []struct {
		name string
		s    *Slice
		want bool
	}{
		{"closed", &Slice{Polygons: [][]int{{0, 1, 0}}}, true},
		{"open", &Slice{Polygons: [][]int{{0, 1, 2}}}, false},
		{"one", &Slice{Polygons: [][]int{{0}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.AllPolygonsAreClosed(); got != tt.want {
				t.Errorf("Slice.AllPolygonsAreClosed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_IsPolygonValid(t *testing.T) {
	type args struct {
		index int
	}
	tests := []struct {
		name string
		s    *Slice
		args args
		want bool
	}{
		{"empty", new(Slice), args{0}, false},
		{"invalid1", &Slice{Polygons: [][]int{{0}}}, args{0}, false},
		{"invalid2", &Slice{Polygons: [][]int{{0, 1}}}, args{0}, false},
		{"valid", &Slice{Polygons: [][]int{{0, 1, 2}}}, args{0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsPolygonValid(tt.args.index); got != tt.want {
				t.Errorf("Slice.IsPolygonValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceStack_AddSlice(t *testing.T) {
	type args struct {
		slice *Slice
	}
	tests := []struct {
		name    string
		s       *SliceStack
		args    args
		want    int
		wantErr bool
	}{
		{"lower", &SliceStack{BottomZ: 1}, args{&Slice{TopZ: 0.5}}, 0, true},
		{"top", &SliceStack{Slices: []*Slice{{TopZ: 1.0}}}, args{&Slice{TopZ: 0.5}}, 0, true},
		{"ok", &SliceStack{BottomZ: 1, Slices: []*Slice{{TopZ: 1.0}}}, args{&Slice{TopZ: 2.0}}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AddSlice(tt.args.slice)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceStack.AddSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SliceStack.AddSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceStackResource_ReferencePath(t *testing.T) {
	s1, _ := NewSliceStackResource(0, new(Model), new(SliceStack))
	s2, _ := NewSliceStackResource(0, new(Model), &SliceStack{UsesSliceRef: true})
	tests := []struct {
		name string
		s    *SliceStackResource
		want string
	}{
		{"noref", s1, ""},
		{"ref", s2, "/2D/2dmodel_1.model"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ReferencePath(); got != tt.want {
				t.Errorf("SliceStackResource.ReferencePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSliceStackResource(t *testing.T) {
	model := new(Model)
	type args struct {
		id    uint64
		model *Model
		stack *SliceStack
	}
	tests := []struct {
		name    string
		args    args
		want    *SliceStackResource
		wantErr bool
	}{
		{"base", args{0, model, nil}, &SliceStackResource{Resource: Resource{Model: model, ResourceID: &ResourceID{"", 0, 1}}}, false},
		{"dup", args{0, model, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSliceStackResource(tt.args.id, tt.args.model, tt.args.stack)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSliceStackResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSliceStackResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
