package geometry

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestFaceNormal(t *testing.T) {
	type args struct {
		n1 mgl32.Vec3
		n2 mgl32.Vec3
		n3 mgl32.Vec3
	}
	tests := []struct {
		name string
		args args
		want mgl32.Vec3
	}{
		{"X", args{mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 20.0, -20.0}, mgl32.Vec3{0.0, 0.0019989014, 0.0019989014}}, mgl32.Vec3{1, 0, 0}},
		{"-Y", args{mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{20.0, 0.0, -20.0}, mgl32.Vec3{0.0019989014, 0.0, 0.0019989014}}, mgl32.Vec3{0, -1, 0}},
		{"Z", args{mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{20.0, -20.0, 0.0}, mgl32.Vec3{0.0019989014, 0.0019989014, 0.0}}, mgl32.Vec3{0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FaceNormal(tt.args.n1, tt.args.n2, tt.args.n3); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FaceNormal() = %v, want %v", got, tt.want)
			}
		})
	}
}
