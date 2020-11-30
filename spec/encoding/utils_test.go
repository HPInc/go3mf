package encoding

import (
	"image/color"
	"reflect"
	"testing"
)

func TestReadRGBA(t *testing.T) {
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
			gotC, err := ParseRGBA(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRGBA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("ParseRGBA() = %v, want %v", gotC, tt.wantC)
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
		{"red", color.RGBA{255, 0, 0, 255}, "#ff0000"},
		{"green", color.RGBA{0, 255, 0, 255}, "#00ff00"},
		{"blue", color.RGBA{0, 0, 255, 255}, "#0000ff"},
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
