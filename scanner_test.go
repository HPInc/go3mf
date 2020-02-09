package go3mf

import (
	"encoding/xml"
	"errors"
	"image/color"
	"reflect"
	"testing"
)

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

func TestReadRGB(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantC   color.RGBA
		wantErr bool
	}{
		{"empty", args{""}, color.RGBA{}, true},
		{"nohashrgb", args{"101010"}, color.RGBA{}, true},
		{"nohashrgba", args{"10101010"}, color.RGBA{}, true},
		{"invalidChar", args{"#€0101010"}, color.RGBA{}, true},
		{"invalidChar", args{"#T0101010"}, color.RGBA{0, 16, 16, 16}, true},
		{"rgb", args{"#112233"}, color.RGBA{17, 34, 51, 255}, false},
		{"rgb", args{"#ff0033"}, color.RGBA{255, 0, 51, 255}, false},
		{"rgba", args{"#000233ff"}, color.RGBA{0, 2, 51, 255}, false},
		{"rgbaLetter", args{"#ff0233AB"}, color.RGBA{255, 2, 51, 171}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := ParseRGB(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRGB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("ParseRGB() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestFormatRGBA(t *testing.T) {
	tests := []struct {
		name string
		c    color.RGBA
		want string
	}{
		{"base", color.RGBA{200, 250, 60, 80}, "#c8fa3c50"},
		{"red", color.RGBA{255, 0, 0, 255}, "#ff0000ff"},
		{"green", color.RGBA{0, 255, 0, 255}, "#00ff00ff"},
		{"blue", color.RGBA{0, 0, 255, 255}, "#0000ffff"},
		{"transparent", color.RGBA{0, 0, 0, 0}, "#00000000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatRGBA(tt.c); got != tt.want {
				t.Errorf("FormatRGBA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_strictError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		p       *Scanner
		args    args
		wantErr bool
	}{
		{"empty", new(Scanner), args{errors.New("fake_err")}, false},
		{"lax", &Scanner{Strict: false}, args{errors.New("fake_err")}, false},
		{"strict", &Scanner{Strict: true}, args{errors.New("fake_err")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.strictError(tt.args.err)
			if (tt.p.Err != nil) != tt.wantErr {
				t.Errorf("Scanner.strictError() error = %v, wantErr %v", tt.p.Err, tt.wantErr)
				return
			}
		})
	}
}

func Test_baseDecoder_Child(t *testing.T) {
	type args struct {
		in0 xml.Name
	}
	tests := []struct {
		name string
		d    *baseDecoder
		args args
		want NodeDecoder
	}{
		{"base", new(baseDecoder), args{xml.Name{}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Child(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseDecoder.Child() = %v, want %v", got, tt.want)
			}
		})
	}
}
