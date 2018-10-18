package geometry

import (
	"reflect"
	"testing"
)

func TestNewVector2DII(t *testing.T) {
	type args struct {
		x int32
		y int32
	}
	tests := []struct {
		name string
		args args
		want Vector2DI
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVector2DI(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVector2DII() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2DI_Uncast(t *testing.T) {
	type args struct {
		units float32
	}
	tests := []struct {
		name    string
		a       Vector2DI
		args    args
		want    Vector2D
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Uncast(tt.args.units)
			if (err != nil) != tt.wantErr {
				t.Errorf("Vector2DI.Uncast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector2DI.Uncast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2DI_Add(t *testing.T) {
	type args struct {
		b Vector2DI
	}
	tests := []struct {
		name string
		a    Vector2DI
		args args
		want Vector2DI
	}{
		{"zero", NewVector2DI(0, 5), args{NewVector2DI(2, 3)}, NewVector2DI(2, 8)},
		{"random", NewVector2DI(4, 5), args{NewVector2DI(2, 3)}, NewVector2DI(6, 8)},
		{"negative", NewVector2DI(4, 5), args{NewVector2DI(-2, 3)}, NewVector2DI(2, 8)},
		{"big", NewVector2DI(4000, 5000000), args{NewVector2DI(2000, 3000000)}, NewVector2DI(6000, 8000000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Add(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector2DI.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2DI_Sub(t *testing.T) {
	type args struct {
		b Vector2DI
	}
	tests := []struct {
		name string
		a    Vector2DI
		args args
		want Vector2DI
	}{
		{"zero", NewVector2DI(0, 5), args{NewVector2DI(2, 3)}, NewVector2DI(-2, 2)},
		{"random", NewVector2DI(4, 5), args{NewVector2DI(2, 3)}, NewVector2DI(2, 2)},
		{"negative", NewVector2DI(4, -5), args{NewVector2DI(2, 3)}, NewVector2DI(2, -8)},
		{"big", NewVector2DI(4000, 5000000), args{NewVector2DI(2000, 3000000)}, NewVector2DI(2000, 2000000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Sub(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector2DI.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2DI_Scale(t *testing.T) {
	type args struct {
		b int32
	}
	tests := []struct {
		name string
		a    Vector2DI
		args args
		want Vector2DI
	}{
		{"zero", NewVector2DI(0, 5), args{0}, NewVector2DI(0, 0)},
		{"negative", NewVector2DI(4, 5), args{-2}, NewVector2DI(-8, -10)},
		{"random", NewVector2DI(4, 5), args{2}, NewVector2DI(8, 10)},
		{"big", NewVector2DI(4000, 500000), args{2000}, NewVector2DI(8000000, 1000000000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Scale(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vector2DI.Scale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2DI_Dot(t *testing.T) {
	type args struct {
		b Vector2DI
	}
	tests := []struct {
		name string
		a    Vector2DI
		args args
		want int64
	}{
		{"zero", NewVector2DI(0, 40), args{NewVector2DI(4, 0)}, 0},
		{"negative", NewVector2DI(-1, 3), args{NewVector2DI(4, 6)}, 14},
		{"all-negative", NewVector2DI(-1, 3), args{NewVector2DI(4, -6)}, -22},
		{"random", NewVector2DI(1, 3), args{NewVector2DI(4, 6)}, 22},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Dot(tt.args.b); got != tt.want {
				t.Errorf("Vector2DI.Dot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2DI_Length(t *testing.T) {
	tests := []struct {
		name string
		a    Vector2DI
		want float32
	}{
		{"a", NewVector2DI(0, 40), 40},
		{"b", NewVector2DI(-40, 0), 40},
		{"c", NewVector2DI(0, 0), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Length(); got != tt.want {
				t.Errorf("Vector2DI.Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector2DI_Distance(t *testing.T) {
	type args struct {
		b Vector2DI
	}
	tests := []struct {
		name string
		a    Vector2DI
		args args
		want float32
	}{
		{"a", NewVector2DI(0, 40), args{NewVector2DI(0, 4)}, 36},
		{"b", NewVector2DI(0, -40), args{NewVector2DI(0, 4)}, 44},
		{"c", NewVector2DI(-40, 0), args{NewVector2DI(10, 0)}, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Distance(tt.args.b); got != tt.want {
				t.Errorf("Vector2DI.Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}
