package slices

import (
	"reflect"
	"testing"
)

func TestSliceStackResource_Identify(t *testing.T) {
	tests := []struct {
		name  string
		s     *SliceStackResource
		want  string
		want1 uint32
	}{
		{"base", &SliceStackResource{ID: 1, ModelPath: "/3D/3dmodel.model"}, "/3D/3dmodel.model", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.Identify()
			if got != tt.want {
				t.Errorf("SliceStackResource.Identify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SliceStackResource.Identify() got = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSliceResolution_String(t *testing.T) {
	tests := []struct {
		name string
		c    SliceResolution
	}{
		{"fullres", ResolutionFull},
		{"lowres", ResolutionLow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.name {
				t.Errorf("SliceResolution.String() = %v, want %v", got, tt.name)
			}
		})
	}
}

func Test_newSliceResolution(t *testing.T) {
	tests := []struct {
		name   string
		wantR  SliceResolution
		wantOk bool
	}{
		{"fullres", ResolutionFull, true},
		{"lowres", ResolutionLow, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotOk := newSliceResolution(tt.name)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("newSliceResolution() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotOk != tt.wantOk {
				t.Errorf("newSliceResolution() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
