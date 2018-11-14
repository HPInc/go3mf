package go3mf

import (
	"image/color"
	"testing"
)

func TestBaseMaterial_ColotString(t *testing.T) {
	tests := []struct {
		name string
		m    *BaseMaterial
		want string
	}{
		{"base", &BaseMaterial{Color: color.RGBA{200, 250, 60, 80}}, "#c8fa3c50"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ColotString(); got != tt.want {
				t.Errorf("BaseMaterial.ColotString() = %v, want %v", got, tt.want)
			}
		})
	}
}
