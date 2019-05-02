package stl

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/geo"
)

func TestNewEncoder(t *testing.T) {
	tests := []struct {
		name  string
		want  *Encoder
		wantW string
	}{
		{"base", &Encoder{w: new(bytes.Buffer), encodingType: Binary}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if got := NewEncoder(w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEncoder() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NewEncoder() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestNewEncoderType(t *testing.T) {
	type args struct {
		encodingType EncodingType
	}
	tests := []struct {
		name  string
		args  args
		want  *Encoder
		wantW string
	}{
		{"binary", args{Binary}, &Encoder{w: new(bytes.Buffer), encodingType: Binary}, ""},
		{"ascii", args{ASCII}, &Encoder{w: new(bytes.Buffer), encodingType: ASCII}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if got := NewEncoderType(w, tt.args.encodingType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEncoderType() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NewEncoderType() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestEncoder_Encode(t *testing.T) {
	type args struct {
		m *geo.Mesh
	}
	tests := []struct {
		name    string
		e       *Encoder
		args    args
		wantErr bool
	}{
		{"ascii", NewEncoderType(new(bytes.Buffer), ASCII), args{new(geo.Mesh)}, false},
		{"binary", NewEncoderType(new(bytes.Buffer), Binary), args{new(geo.Mesh)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Encode(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Encoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_faceNormal(t *testing.T) {
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
			if got := faceNormal(tt.args.n1, tt.args.n2, tt.args.n3); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("faceNormal() = %v, want %v", got, tt.want)
			}
		})
	}
}