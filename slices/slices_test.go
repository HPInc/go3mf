package slices

import (
	"reflect"
	"testing"

	"github.com/qmuntal/go3mf"
)

var _ go3mf.SpecDecoder = new(Spec)
var _ go3mf.SpecValidator = new(Spec)
var _ go3mf.Asset = new(SliceStack)
var _ go3mf.Marshaler = new(SliceStack)
var _ go3mf.AttrMarshaler = new(SliceStackInfo)

func TestSliceStack_Identify(t *testing.T) {
	tests := []struct {
		name string
		s    *SliceStack
		want uint32
	}{
		{"base", &SliceStack{ID: 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Identify()
			if got != tt.want {
				t.Errorf("SliceStack.Identify() got = %v, want %v", got, tt.want)
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
