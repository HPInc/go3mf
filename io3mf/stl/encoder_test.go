package stl

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/qmuntal/go3mf/mesh"
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
		m *mesh.Mesh
	}
	tests := []struct {
		name    string
		e       *Encoder
		args    args
		wantErr bool
	}{
		{"ascii", NewEncoderType(new(bytes.Buffer), ASCII), args{mesh.NewMesh()}, false},
		{"binary", NewEncoderType(new(bytes.Buffer), Binary), args{mesh.NewMesh()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Encode(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Encoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
