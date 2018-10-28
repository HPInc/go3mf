package meshinfo

import (
	"reflect"
	"testing"
)

func TestBaseMaterial_Invalidate(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseMaterial
	}{
		{"base", &BaseMaterial{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Invalidate()
			want := new(BaseMaterial)
			if !reflect.DeepEqual(tt.b, want) {
				t.Errorf("BaseMaterial.Invalidate() = %v, want %v", tt.b, want)
			}
		})
	}
}
