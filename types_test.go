package go3mf

import (
	"testing"
)

func TestUnits_String(t *testing.T) {
	tests := []struct {
		name string
		u    Units
	}{
		{"micron", UnitMicrometer},
		{"millimeter", UnitMillimeter},
		{"centimeter", UnitCentimeter},
		{"inch", UnitInch},
		{"foot", UnitFoot},
		{"meter", UnitMeter},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.name {
				t.Errorf("Units.String() = %v, want %v", got, tt.name)
			}
		})
	}
}

func TestTileStyle_String(t *testing.T) {
	tests := []struct {
		name string
		t    TileStyle
	}{
		{"wrap", TileWrap},
		{"mirror", TileMirror},
		{"clamp", TileClamp},
		{"none", TileNone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.name {
				t.Errorf("TileStyle.String() = %v, want %v", got, tt.name)
			}
		})
	}
}

func TestTextureFilter_String(t *testing.T) {
	tests := []struct {
		name string
		t    TextureFilter
	}{
		{"auto", TextureFilterAuto},
		{"linear", TextureFilterLinear},
		{"nearest", TextureFilterNearest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.name {
				t.Errorf("TextureFilter.String() = %v, want %v", got, tt.name)
			}
		})
	}
}

func TestClipMode_String(t *testing.T) {
	tests := []struct {
		name string
		c    ClipMode
	}{
		{"none", ClipNone},
		{"inside", ClipInside},
		{"outside", ClipOutside},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.name {
				t.Errorf("ClipMode.String() = %v, want %v", got, tt.name)
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

func TestObjectType_String(t *testing.T) {
	tests := []struct {
		name string
		o    ObjectType
	}{
		{"model", ObjectTypeModel},
		{"other", ObjectTypeOther},
		{"support", ObjectTypeSupport},
		{"solidsupport", ObjectTypeSolidSupport},
		{"surface", ObjectTypeSurface},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.String(); got != tt.name {
				t.Errorf("ObjectType.String() = %v, want %v", got, tt.name)
			}
		})
	}
}

func TestTexture2DType_String(t *testing.T) {
	tests := []struct {
		name string
		t    Texture2DType
	}{
		{"image/png", PNGTexture},
		{"image/jpeg", JPEGTexture},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.name {
				t.Errorf("Texture2DType.String() = %v, want %v", got, tt.name)
			}
		})
	}
}
