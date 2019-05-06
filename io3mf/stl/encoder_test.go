package stl

import (
	"bytes"
	"reflect"
	"testing"

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
		n1 geo.Point3D
		n2 geo.Point3D
		n3 geo.Point3D
	}
	tests := []struct {
		name string
		args args
		want geo.Point3D
	}{
		{"X", args{geo.Point3D{0.0, 0.0, 0.0}, geo.Point3D{0.0, 20.0, -20.0}, geo.Point3D{0.0, 0.0019989014, 0.0019989014}}, geo.Point3D{1, 0, 0}},
		{"-Y", args{geo.Point3D{0.0, 0.0, 0.0}, geo.Point3D{20.0, 0.0, -20.0}, geo.Point3D{0.0019989014, 0.0, 0.0019989014}}, geo.Point3D{0, -1, 0}},
		{"Z", args{geo.Point3D{0.0, 0.0, 0.0}, geo.Point3D{20.0, -20.0, 0.0}, geo.Point3D{0.0019989014, 0.0019989014, 0.0}}, geo.Point3D{0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := faceNormal(tt.args.n1, tt.args.n2, tt.args.n3); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("faceNormal() = %v, want %v", got, tt.want)
			}
		})
	}
}
