package meshinfo

import (
	"reflect"
	"testing"
)

func TestNodeColor_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		n    *NodeColor
	}{
		{"base", new(NodeColor)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.Colors[0] = 1
			tt.n.Colors[1] = 2
			tt.n.Colors[2] = 3
			tt.n.Invalidate()
			want := new(NodeColor)
			if !reflect.DeepEqual(tt.n, want) {
				t.Errorf("NodeColor.Invalidate() = %v, want %v", tt.n, want)
			}
		})
	}
}
