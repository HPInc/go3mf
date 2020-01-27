package iohelper

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/qmuntal/go3mf/geo"
)

func TestParser_NamespaceRegistered(t *testing.T) {
	type args struct {
		ns string
	}
	tests := []struct {
		name string
		p    *Scanner
		args args
		want bool
	}{
		{"empty", &Scanner{Namespaces: map[string]string{"p": "http://xml.com"}}, args{""}, false},
		{"exist", &Scanner{Namespaces: map[string]string{"p": "http://xml.com"}}, args{"http://xml.com"}, true},
		{"noexist", &Scanner{Namespaces: map[string]string{"p": "http://xml.com"}}, args{"xmls"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.NamespaceRegistered(tt.args.ns); got != tt.want {
				t.Errorf("Parser.NamespaceRegistered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_ParseToMatrixOptional(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want geo.Matrix
	}{
		{"empty", args{""}, geo.Matrix{}},
		{"11values", args{"1 1 1 1 1 1 1 1 1 1 1"}, geo.Matrix{}},
		{"13values", args{"1 1 1 1 1 1 1 1 1 1 1 1 1"}, geo.Matrix{}},
		{"char", args{"1 1 a 1 1 1 1 1 1 1 1 1"}, geo.Matrix{}},
		{"base", args{"1 1 1 1 1 1 1 1 1 1 1 1"}, geo.Matrix{1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1}},
		{"other", args{"0 1 2 10 11 12 20 21 22 30 31 32"}, geo.Matrix{0, 1, 2, 0, 10, 11, 12, 0, 20, 21, 22, 0, 30, 31, 32, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(Scanner).ParseToMatrixOptional("", tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scanner.ParseToMatrixOptional() = %v, want %v", got, tt.want)
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
			gotC, err := ReadRGB(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadRGB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("ReadRGB() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
