package stl

import (
	"bytes"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf/internal/mesh"
)

func Test_asciiDecoder_decode(t *testing.T) {
	triangle := createASCIITriangle()
	tests := []struct {
		name    string
		d       *asciiDecoder
		want    *mesh.Mesh
		wantErr bool
	}{
		{"eof", &asciiDecoder{r: bytes.NewReader(make([]byte, 0))}, mesh.NewMesh(), false},
		{"base", &asciiDecoder{r: bytes.NewBufferString(triangle)}, createMeshTriangle(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.decode()
			if (err != nil) != tt.wantErr {
				t.Errorf("asciiDecoder.decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := deep.Equal(got, tt.want); diff != nil {
					t.Errorf("asciiDecoder.decode() = %v", diff)
					return
				}
			}
		})
	}
}

func Test_asciiEncoder_encode(t *testing.T) {
	triangle := createMeshTriangle()
	type args struct {
		m *mesh.Mesh
	}
	tests := []struct {
		name    string
		e       *asciiEncoder
		args    args
		wantErr bool
	}{
		{"base", &asciiEncoder{w: new(bytes.Buffer)}, args{triangle}, false},
		{"error", &asciiEncoder{w: new(errorWriter)}, args{triangle}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.encode(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("asciiEncoder.encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// We do decoder and then encoder again, and the result must be the same
				decoder := &asciiDecoder{r: tt.e.w.(*bytes.Buffer)}
				got, _ := decoder.decode()
				if diff := deep.Equal(got, tt.args.m); diff != nil {
					t.Errorf("asciiDecoder.encode() = %v", diff)
					return
				}
			}
		})
	}
}

func createASCIITriangle() string {
	return `solid 
  		facet normal 0 0 0
    		outer loop
      			vertex -20.0 -20.0 0.0
      			vertex 20.0 -20.0 0.0
      			vertex 0.0019989014 0.0019989014 39.998
    		endloop
  		endfacet
  		facet normal 0 0 0
			outer loop
			vertex -20.0 20.0 0.0
			vertex 20.0 -20.0 0.0
			vertex -20.0 -20.0 0.0
			endloop
		endfacet
		facet normal 0 0 0
			outer loop
			vertex -20.0 -20.0 0.0
			vertex 0.0 0.0019989014 39.998
			vertex -20.0 20.0 0.0
			endloop
		endfacet
		facet normal 0 0 0
			outer loop
			vertex 20.0 -20.0 0.0
			vertex 20.0 20.0 0.0
			vertex 0.0019989014 0.0019989014 39.998
			endloop
		endfacet
		facet normal 0 0 0
			outer loop
			vertex 20.0 20.0 0.0
			vertex -20.0 20.0 0.0
			vertex 0.0019989014 0.0019989014 39.998
			endloop
		endfacet
		facet normal 0 0 0
			outer loop
			vertex 20.0 20.0 0.0
			vertex 20.0 -20.0 0.0
			vertex -20.0 20.0 0.0
			endloop
		endfacet
	endsolid`
}
