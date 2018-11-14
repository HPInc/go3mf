package stl

import (
	"bytes"
	"io"
	"testing"

	"github.com/qmuntal/go3mf/internal/mesh"
)

func TestDecodeUnits(t *testing.T) {
	triangle := createBinaryTriangle()
	triangle[0] = 0x73
	triangle[1] = 0x6f
	triangle[2] = 0x6c
	triangle[3] = 0x69
	triangle[4] = 0x64
	type args struct {
		r     io.Reader
		units float32
	}
	tests := []struct {
		name    string
		args    args
		want    *mesh.Mesh
		wantErr bool
	}{
		{"base", args{bytes.NewReader(triangle), 0.0}, createMeshTriangle(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeUnits(tt.args.r, tt.args.units)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeUnits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.ApproxEqual(tt.want) {
				t.Errorf("DecodeUnits() = %v, want %v", got, tt.want)
			}
		})
	}
}
