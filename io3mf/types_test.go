package io3mf

import (
	"reflect"
	"testing"

	go3mf "github.com/qmuntal/go3mf"
	mesh "github.com/qmuntal/go3mf/mesh"
)

func Test_newTextureFilter(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		want   go3mf.TextureFilter
		wantOk bool
	}{
		{"auto", go3mf.TextureFilterAuto, true},
		{"linear", go3mf.TextureFilterLinear, true},
		{"nearest", go3mf.TextureFilterNearest, true},
		{"empty", go3mf.TextureFilterAuto, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := newTextureFilter(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTextureFilter() got = %v, want %v", got, tt.want)
			}
			if got != tt.want {
				t.Errorf("newTextureFilter() got1 = %v, want %v", got1, tt.want)
			}
		})
	}
}

func Test_newTileStyle(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		want  go3mf.TileStyle
		want1 bool
	}{
		{"wrap", go3mf.TileWrap, true},
		{"mirror", go3mf.TileMirror, true},
		{"clamp", go3mf.TileClamp, true},
		{"none", go3mf.TileNone, true},
		{"empty", go3mf.TileWrap, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := newTileStyle(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTileStyle() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("newTileStyle() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_newTexture2DType(t *testing.T) {
	tests := []struct {
		name  string
		want  go3mf.Texture2DType
		want1 bool
	}{
		{"image/png", go3mf.PNGTexture, true},
		{"image/jpeg", go3mf.JPEGTexture, true},
		{"", go3mf.Texture2DType(0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := newTexture2DType(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTexture2DType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("newTexture2DType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_newObjectType(t *testing.T) {
	tests := []struct {
		name   string
		wantO  go3mf.ObjectType
		wantOk bool
	}{
		{"model", go3mf.ObjectTypeModel, true},
		{"other", go3mf.ObjectTypeOther, true},
		{"support", go3mf.ObjectTypeSupport, true},
		{"solidsupport", go3mf.ObjectTypeSolidSupport, true},
		{"surface", go3mf.ObjectTypeSurface, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, gotOk := newObjectType(tt.name)
			if !reflect.DeepEqual(gotO, tt.wantO) {
				t.Errorf("newObjectType() gotO = %v, want %v", gotO, tt.wantO)
			}
			if gotOk != tt.wantOk {
				t.Errorf("newObjectType() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_newSliceResolution(t *testing.T) {
	tests := []struct {
		name   string
		wantR  go3mf.SliceResolution
		wantOk bool
	}{
		{"fullres", go3mf.ResolutionFull, true},
		{"lowres", go3mf.ResolutionLow, true},
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

func Test_newClipMode(t *testing.T) {
	tests := []struct {
		name   string
		wantC  go3mf.ClipMode
		wantOk bool
	}{
		{"none", go3mf.ClipNone, true},
		{"inside", go3mf.ClipInside, true},
		{"outside", go3mf.ClipOutside, true},
		{"empty", go3mf.ClipNone, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, gotOk := newClipMode(tt.name)
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("newClipMode() gotC = %v, want %v", gotC, tt.wantC)
			}
			if gotOk != tt.wantOk {
				t.Errorf("newClipMode() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_newUnits(t *testing.T) {
	tests := []struct {
		name  string
		want  go3mf.Units
		want1 bool
	}{
		{"micron", go3mf.UnitMicrometer, true},
		{"millimeter", go3mf.UnitMillimeter, true},
		{"centimeter", go3mf.UnitCentimeter, true},
		{"inch", go3mf.UnitInch, true},
		{"foot", go3mf.UnitFoot, true},
		{"meter", go3mf.UnitMeter, true},
		{"", go3mf.UnitMillimeter, false},
		{"other", go3mf.UnitMillimeter, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := newUnits(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUnits() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("newUnits() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_newCapMode(t *testing.T) {
	tests := []struct {
		name   string
		wantT  mesh.CapMode
		wantOk bool
	}{
		{"sphere", mesh.CapModeSphere, true},
		{"hemisphere", mesh.CapModeHemisphere, true},
		{"butt", mesh.CapModeButt, true},
		{"empty", mesh.CapModeSphere, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotOk := newCapMode(tt.name)
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("newCapMode() gotT = %v, want %v", gotT, tt.wantT)
			}
			if gotOk != tt.wantOk {
				t.Errorf("newCapMode() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_newBlendMethod(t *testing.T) {
	tests := []struct {
		name   string
		wantB  go3mf.BlendMethod
		wantOk bool
	}{
		{"mix", go3mf.BlendMix, true},
		{"multiply", go3mf.BlendMultiply, true},
		{"empty", go3mf.BlendMix, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, gotOk := newBlendMethod(tt.name)
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("newBlendMethod() gotB = %v, want %v", gotB, tt.wantB)
			}
			if gotOk != tt.wantOk {
				t.Errorf("newBlendMethod() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
