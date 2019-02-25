package model

import (
	"testing"
)

func newTexture2D() *Texture2DResource {
	r, _ := NewTexture2DResource(0, new(Model))
	return r
}

func TestTexture2DResource_Box(t *testing.T) {
	tests := []struct {
		name       string
		t          *Texture2DResource
		wantU      float32
		wantV      float32
		wantWidth  float32
		wantHeight float32
		wantHasBox bool
	}{
		{"base", newTexture2D(), 0, 0, 1, 1, false},
		{"set", newTexture2D().SetBox(1, 2, 3, 4), 1, 2, 3, 4, true},
		{"clear", newTexture2D().SetBox(1, 2, 3, 4).ClearBox(), 0, 0, 1, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, gotV, gotWidth, gotHeight, gotHasBox := tt.t.Box()
			if gotU != tt.wantU {
				t.Errorf("Texture2DResource.Box() gotU = %v, want %v", gotU, tt.wantU)
			}
			if gotV != tt.wantV {
				t.Errorf("Texture2DResource.Box() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotWidth != tt.wantWidth {
				t.Errorf("Texture2DResource.Box() gotWidth = %v, want %v", gotWidth, tt.wantWidth)
			}
			if gotHeight != tt.wantHeight {
				t.Errorf("Texture2DResource.Box() gotHeight = %v, want %v", gotHeight, tt.wantHeight)
			}
			if gotHasBox != tt.wantHasBox {
				t.Errorf("Texture2DResource.Box() gotHasBox = %v, want %v", gotHasBox, tt.wantHasBox)
			}
		})
	}
}
