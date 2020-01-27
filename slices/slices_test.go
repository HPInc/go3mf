package slices

import (
	"reflect"
	"testing"
	"github.com/qmuntal/go3mf"
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
			want := go3mf.Point2D{tt.args.x, tt.args.y}
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
		{"duplicated", &Slice{Polygons: [][]int{{0}}, Vertices: []go3mf.Point2D{{}}}, args{0, 0}, true},
		{"base", &Slice{Polygons: [][]int{{}}, Vertices: []go3mf.Point2D{{}}}, args{0, 0}, false},
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

func TestSliceStackResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		s     *SliceStackResource
		want  string
		want1 uint32
	}{
		{"base", &SliceStackResource{ID: 1, ModelPath: "3d/3dmodel.model"}, "3d/3dmodel.model", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.Identify()
			if got != tt.want {
				t.Errorf("SliceStackResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SliceStackResource.Identify() got = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSliceResolution_String(t *testing.T) {
	tests := []struct {
		name string
		c    SliceResolution
	}{
		{"fullres", ResolutionFull},
		{"lowres", ResolutionLow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.name {
				t.Errorf("SliceResolution.String() = %v, want %v", got, tt.name)
			}
		})
	}
}

func Test_newSliceResolution(t *testing.T) {
	tests := []struct {
		name   string
		wantR  SliceResolution
		wantOk bool
	}{
		{"fullres", ResolutionFull, true},
		{"lowres", ResolutionLow, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotOk := newSliceResolution(tt.name)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("newSliceResolution() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotOk != tt.wantOk {
				t.Errorf("newSliceResolution() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
