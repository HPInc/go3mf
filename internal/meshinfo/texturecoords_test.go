package meshinfo

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestTextureCoords_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		t    *TextureCoords
	}{
		{"base", new(TextureCoords)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.TextureID = 2
			tt.t.Coords[0] = mgl32.Vec2{1.0, 2.0}
			tt.t.Coords[1] = mgl32.Vec2{5.0, 3.0}
			tt.t.Coords[2] = mgl32.Vec2{6.0, 4.0}
			tt.t.Invalidate()
			want := new(TextureCoords)
			if !reflect.DeepEqual(tt.t, want) {
				t.Errorf("TextureCoords.Invalidate() = %v, want %v", tt.t, want)
			}
		})
	}
}
