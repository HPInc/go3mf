package go3mf

import (
	"math"
	"reflect"
	"testing"
)

func TestPoint2D_X(t *testing.T) {
	tests := []struct {
		name string
		n    Point2D
		want float32
	}{
		{"base", Point2D{1, 2}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.X(); got != tt.want {
				t.Errorf("Point2D.X() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint2D_Y(t *testing.T) {
	tests := []struct {
		name string
		n    Point2D
		want float32
	}{
		{"base", Point2D{1, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Y(); got != tt.want {
				t.Errorf("Point2D.Y() = %v, want %v", got, tt.want)
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

func Test_vectorTree_AddFindVector(t *testing.T) {
	p := vectorTree{}
	type args struct {
		vec   Point3D
		value uint32
	}
	tests := []struct {
		name string
		t    vectorTree
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
	p := vectorTree{}
	p.AddVector(Point3D{1, 2, 5.3}, 1)
	type args struct {
		vec Point3D
	}
	tests := []struct {
		name string
		t    vectorTree
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

func TestMatrix_Mul(t *testing.T) {
	type args struct {
		m Matrix
	}
	tests := []struct {
		name string
		m1   Matrix
		args args
		want Matrix
	}{
		{"base", Identity(), args{Identity()}, Identity()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m1.Mul(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Matrix.Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_String(t *testing.T) {
	tests := []struct {
		name string
		m    Matrix
		want string
	}{
		{"base", Matrix{
			1.12345, 2.12345, 3.12345, 4.12345, 5.12345, 6.12345, 7.12345, 8.12345,
			9.12345, 10.12345, 11.12345, 12.12345, 13.12345, 14.12345, 15.12345, 16.12345,
		}, "1.123 2.123 3.123 5.123 6.123 7.123 9.123 10.123 11.123 13.123 14.123 15.123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("Matrix.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseToMatrixOptional(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want Matrix
	}{
		{"empty", args{""}, Matrix{}},
		{"11values", args{"1 1 1 1 1 1 1 1 1 1 1"}, Matrix{}},
		{"13values", args{"1 1 1 1 1 1 1 1 1 1 1 1 1"}, Matrix{}},
		{"char", args{"1 1 a 1 1 1 1 1 1 1 1 1"}, Matrix{}},
		{"base", args{"1 1 1 1 1 1 1 1 1 1 1 1"}, Matrix{1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1}},
		{"other", args{"0 1 2 10 11 12 20 21 22 30 31 32"}, Matrix{0, 1, 2, 0, 10, 11, 12, 0, 20, 21, 22, 0, 30, 31, 32, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ParseMatrix(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scanner.ParseMatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_Translate(t *testing.T) {
	type args struct {
		x float32
		y float32
		z float32
	}
	tests := []struct {
		name string
		m    Matrix
		args args
		want Matrix
	}{
		{"zero", Matrix{}, args{0, 0, 0}, Matrix{}},
		{"identity", Identity(), args{1, 2, 3}, Matrix{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 2, 3, 1}},
		{"other", Matrix{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 0}, args{1, 2, 3}, Matrix{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 4, 6, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Translate(tt.args.x, tt.args.y, tt.args.z); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Matrix.Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_Mul3D(t *testing.T) {
	type args struct {
		v Point3D
	}
	tests := []struct {
		name string
		m1   Matrix
		args args
		want Point3D
	}{
		{"zero", Matrix{}, args{Point3D{}}, Point3D{}},
		{"identity", Identity(), args{Point3D{}}, Point3D{}},
		{"identity", Identity(), args{Point3D{1, 2, 3}}, Point3D{1, 2, 3}},
		{"translate", Identity().Translate(1, 2, 3), args{Point3D{}}, Point3D{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m1.Mul3D(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Matrix.Mul3D() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_Mul2D(t *testing.T) {
	type args struct {
		v Point2D
	}
	tests := []struct {
		name string
		m1   Matrix
		args args
		want Point2D
	}{
		{"zero", Matrix{}, args{Point2D{}}, Point2D{}},
		{"identity", Identity(), args{Point2D{}}, Point2D{}},
		{"identity", Identity(), args{Point2D{1, 2}}, Point2D{1, 2}},
		{"translate", Identity().Translate(1, 2, 0), args{Point2D{}}, Point2D{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m1.Mul2D(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Matrix.Mul2D() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pairMatch_AddMatch(t *testing.T) {
	p := make(pairMatch)
	type args struct {
		data1 uint32
		data2 uint32
		param uint32
	}
	tests := []struct {
		name string
		t    pairMatch
		args args
	}{
		{"new", p, args{1, 1, 2}},
		{"old", p, args{1, 1, 4}},
		{"new2", p, args{2, 1, 5}},
		{"old2", p, args{2, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.AddMatch(tt.args.data1, tt.args.data2, tt.args.param)
			got, ok := p.CheckMatch(tt.args.data1, tt.args.data2)
			if !ok {
				t.Error("pairMatch.AddMatch() haven't added the match")
				return
			}
			if got != tt.args.param {
				t.Errorf("pairMatch.CheckMatch() = %v, want %v", got, tt.args.param)
			}
		})
	}
}

func Test_newPairEntry(t *testing.T) {
	type args struct {
		data1 uint32
		data2 uint32
	}
	tests := []struct {
		name string
		args args
		want pairEntry
	}{
		{"d1=d2", args{1, 1}, pairEntry{1, 1}},
		{"d1>d2", args{2, 1}, pairEntry{1, 2}},
		{"d1<d2", args{1, 2}, pairEntry{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPairEntry(tt.args.data1, tt.args.data2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPairEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}
