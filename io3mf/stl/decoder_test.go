package stl

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestNewDecoder(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want *Decoder
	}{
		{"base", args{new(bytes.Buffer)}, &Decoder{r: new(bytes.Buffer)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDecoder(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDecoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	triangleASCII := createASCIITriangle()
	triangle := createBinaryTriangle()
	triangle[0] = 0x73
	triangle[1] = 0x6f
	triangle[2] = 0x6c
	triangle[3] = 0x69
	triangle[4] = 0x64
	tests := []struct {
		name    string
		d       *Decoder
		want    *go3mf.MeshResource
		wantErr bool
	}{
		{"empty", NewDecoder(new(bytes.Buffer)), nil, true},
		{"binary", NewDecoder(bytes.NewReader(triangle)), &go3mf.MeshResource{
			ObjectResource: go3mf.ObjectResource{ID: 1},
			Mesh:           createMeshTriangle(),
		}, false},
		{"ascii", NewDecoder(bytes.NewBufferString(triangleASCII)), &go3mf.MeshResource{
			ObjectResource: go3mf.ObjectResource{ID: 1},
			Mesh:           createMeshTriangle(),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(go3mf.Model)
			err := tt.d.Decode(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := deep.Equal(got.Resources[0], tt.want); diff != nil {
					t.Errorf("Decoder.Decode() = %v", diff)
					return
				}
			}
		})
	}
}
