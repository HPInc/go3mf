package model

import (
	"reflect"
	"testing"
)

func TestNewUnits(t *testing.T) {
	tests := []struct {
		name  string
		want  Units
		want1 bool
	}{
		{"micron", Micrometer, true},
		{"millimeter", Millimeter, true},
		{"centimeter", Centimeter, true},
		{"inch", Inch, true},
		{"foot", Foot, true},
		{"meter", Meter, true},
		{"empty", Units(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewUnits(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnits() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NewUnits() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewTexture2DType(t *testing.T) {
	tests := []struct {
		name  string
		want  Texture2DType
		want1 bool
	}{
		{"image/png", PNGTexture, true},
		{"image/jpeg", JPEGTexture, true},
		{"empty", Texture2DType(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewTexture2DType(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTexture2DType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NewTexture2DType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewTileStyle(t *testing.T) {
	tests := []struct {
		name  string
		want  TileStyle
		want1 bool
	}{
		{"wrap", WrapTile, true},
		{"mirror", MirrorTile, true},
		{"clamp", ClampTile, true},
		{"none", NoneTile, true},
		{"empty", TileStyle(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewTileStyle(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTileStyle() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NewTileStyle() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewTextureFilter(t *testing.T) {
	tests := []struct {
		name  string
		want  TextureFilter
		want1 bool
	}{
		{"auto", AutoFilter, true},
		{"linear", LinearFilter, true},
		{"nearest", NearestFilter, true},
		{"empty", TextureFilter(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewTextureFilter(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextureFilter() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NewTextureFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
