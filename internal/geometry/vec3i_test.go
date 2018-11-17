package geometry

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func Test_newvec3IFromVec3(t *testing.T) {
	type args struct {
		vec mgl32.Vec3
	}
	tests := []struct {
		name string
		args args
		want vec3I
	}{
		{"base", args{mgl32.Vec3{1.2, 2.3, 3.4}}, vec3I{1200000, 2300000, 3400000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newvec3IFromVec3(tt.args.vec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newvec3IFromVec3() = %v, want %v", got, tt.want)
			}
		})
	}
}
