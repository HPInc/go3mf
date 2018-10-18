package geometry

import (
	"reflect"
	"testing"
)

func TestNewVector2D(t *testing.T) {
	type args struct {
		x float32
		y float32
	}
	tests := []struct {
		name string
		args args
		want Vector2D
	}{
		{"4-5", args{4.0, 5.0}, Vector2D{4.0, 5.0}},
		{"2-3", args{2.0, 3.0}, Vector2D{2.0, 3.0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVector2D(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVector2D() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2D_Add(t *testing.T) {
	type args struct {
		b Vector2D
	}
	tests := []struct {
		name string
		a    Vector2D
		args args
		want Vector2D
	}{
		{"zero", NewVector2D(0.0, 5.0), args{NewVector2D(2.0, 3.0)}, NewVector2D(2.0, 8.0)},
		{"random", NewVector2D(4.0, 5.0), args{NewVector2D(2.0, 3.0)}, NewVector2D(6.0, 8.0)},
		{"big", NewVector2D(4000.0, 5000000.0), args{NewVector2D(2000.0, 3000000.0)}, NewVector2D(6000.0, 8000000.0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Add(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector2D.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2D_Sub(t *testing.T) {
	type args struct {
		b Vector2D
	}
	tests := []struct {
		name string
		a    Vector2D
		args args
		want Vector2D
	}{
		{"zero", NewVector2D(0.0, 5.0), args{NewVector2D(2.0, 3.0)}, NewVector2D(-2.0, 2.0)},
		{"random", NewVector2D(4.0, 5.0), args{NewVector2D(2.0, 3.0)}, NewVector2D(2.0, 2.0)},
		{"big", NewVector2D(4000.0, 5000000.0), args{NewVector2D(2000.0, 3000000.0)}, NewVector2D(2000.0, 2000000.0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Sub(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector2D.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2D_Scale(t *testing.T) {
	type args struct {
		b float32
	}
	tests := []struct {
		name string
		a    Vector2D
		args args
		want Vector2D
	}{
		{"zero", NewVector2D(0.0, 5.0), args{0}, NewVector2D(0.0, 0.0)},
		{"negative", NewVector2D(4.0, 5.0), args{-2.0}, NewVector2D(-8.0, -10.0)},
		{"random", NewVector2D(4.0, 5.0), args{2.0}, NewVector2D(8.0, 10.0)},
		{"big", NewVector2D(4000.0, 5000000.0), args{2000}, NewVector2D(8000000.0, 10000000000.0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Scale(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector2D.Scale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2D_Combine(t *testing.T) {
	type args struct {
		factor1 float32
		b       Vector2D
		factor2 float32
	}
	tests := []struct {
		name string
		a    Vector2D
		args args
		want Vector2D
	}{
		{"zero", NewVector2D(2.0, 5.0), args{0.0, NewVector2D(3.0, 4.0), 0.0}, NewVector2D(0.0, 0.0)},
		{"one", NewVector2D(2.0, 5.0), args{1.0, NewVector2D(3.0, 4.0), 1.0}, NewVector2D(5.0, 9.0)},
		{"three", NewVector2D(2.0, 5.0), args{3.0, NewVector2D(3.0, 4.0), 3.0}, NewVector2D(15.0, 27.0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Combine(tt.args.factor1, tt.args.b, tt.args.factor2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector2D.Combine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2D_Dot(t *testing.T) {
	type args struct {
		b Vector2D
	}
	tests := []struct {
		name string
		a    Vector2D
		args args
		want float32
	}{
		{"zero", NewVector2D(0.0, 40.0), args{NewVector2D(4.0, 0.0)}, 0.0},
		{"negative", NewVector2D(-1.5, 3.0), args{NewVector2D(4.0, 6.0)}, 12.0},
		{"all-negative", NewVector2D(-1.5, 3.0), args{NewVector2D(4.0, -6.0)}, -24.0},
		{"random", NewVector2D(1.5, 3.0), args{NewVector2D(4.0, 6.0)}, 24.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Dot(tt.args.b); got != tt.want {
				t.Errorf("Vector2D.Dot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2D_Cross(t *testing.T) {
	type args struct {
		b Vector2D
	}
	tests := []struct {
		name string
		a    Vector2D
		args args
		want float32
	}{
		{"zero", NewVector2D(0.0, 40.0), args{NewVector2D(0.0, 4.0)}, 0.0},
		{"negative", NewVector2D(-1.5, 3.0), args{NewVector2D(4.0, 6.0)}, -21.0},
		{"all-negative", NewVector2D(-1.5, 3.0), args{NewVector2D(4.0, -6.0)}, -3.0},
		{"random", NewVector2D(1.5, 3.0), args{NewVector2D(4.0, 6.0)}, -3.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Cross(tt.args.b); got != tt.want {
				t.Errorf("Vector2D.Cross() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2D_Length(t *testing.T) {
	tests := []struct {
		name string
		a    Vector2D
		want float32
	}{
		{"a", NewVector2D(0.0, 40.0), 40.0},
		{"b", NewVector2D(-40.0, 0.0), 40.0},
		{"c", NewVector2D(0.0, 0.0), 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Length(); got != tt.want {
				t.Errorf("Vector2D.Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2D_Distance(t *testing.T) {
	type args struct {
		b Vector2D
	}
	tests := []struct {
		name string
		a    Vector2D
		args args
		want float32
	}{
		{"a", NewVector2D(0.0, 40.0), args{NewVector2D(0.0, 4.0)}, 36.0},
		{"b", NewVector2D(0.0, -40.0), args{NewVector2D(0.0, 4.0)}, 44.0},
		{"c", NewVector2D(-40.0, 0.0), args{NewVector2D(10.0, 0.0)}, 50.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Distance(tt.args.b); got != tt.want {
				t.Errorf("Vector2D.Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2D_Normalize(t *testing.T) {
	tests := []struct {
		name    string
		a       Vector2D
		want    Vector2D
		wantErr bool
	}{
		{"withError", NewVector2D(0.0, 0.0), NewVector2D(0.0, 0.0), true},
		{"withErrorByLittle", NewVector2D(0.0, 0.000000000001), NewVector2D(0.0, 0.0), true},
		{"alreadyNormalized", NewVector2D(0.0, 1.0), NewVector2D(0.0, 1.0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Normalize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Vector2D.Normalize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector2D.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2D_Floor(t *testing.T) {
	type args struct {
		units float32
	}
	tests := []struct {
		name    string
		a       Vector2D
		args    args
		want    Vector2DI
		wantErr bool
	}{
		{"zeroUnit", NewVector2D(1.0, 2.0), args{0.0}, NewVector2DI(0, 0), true},
		{"newZeroUnit", NewVector2D(1.0, 2.0), args{0.00000001}, NewVector2DI(0, 0), true},
		{"largeUnit", NewVector2D(1.0, 2.0), args{1001}, NewVector2DI(0, 0), true},
		{"1", NewVector2D(1.0, 2.0), args{1.0}, NewVector2DI(1, 2), false},
		{"zero", NewVector2D(1.0, 2.0), args{10.0}, NewVector2DI(0, 0), false},
		{"100", NewVector2D(100.0, 200.0), args{10.0}, NewVector2DI(10, 20), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Floor(tt.args.units)
			if (err != nil) != tt.wantErr {
				t.Errorf("Vector2D.Floor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector2D.Floor() = %v, want %v", got, tt.want)
			}
		})
	}
}
