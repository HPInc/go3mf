package slices

import (
	"reflect"
	"testing"
)

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
