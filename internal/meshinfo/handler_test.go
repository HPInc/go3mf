package meshinfo

import (
	"reflect"
	"testing"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name string
		want *Handler
	}{
		{"base", &Handler{genericHandler: *newgenericHandler()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_AddBaseMaterialInfo(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		h    *Handler
		args args
	}{
		{"base", NewHandler(), args{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.h.AddBaseMaterialInfo(tt.args.currentFaceCount)
			want, _ := tt.h.BaseMaterialInfo()
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Handler.AddBaseMaterialInfo() = %v, want %v", got, want)
			}
		})
	}
}

func TestHandler_AddTextureCoordsInfo(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		h    *Handler
		args args
	}{
		{"base", NewHandler(), args{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.h.AddTextureCoordsInfo(tt.args.currentFaceCount)
			want, _ := tt.h.TextureCoordsInfo()
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Handler.AddTextureCoordsInfo() = %v, want %v", got, want)
			}
		})
	}
}

func TestHandler_AddNodeColorInfo(t *testing.T) {
	type args struct {
		currentFaceCount uint32
	}
	tests := []struct {
		name string
		h    *Handler
		args args
	}{
		{"base", NewHandler(), args{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.h.AddNodeColorInfo(tt.args.currentFaceCount)
			want, _ := tt.h.NodeColorInfo()
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Handler.AddNodeColorInfo() = %v, want %v", got, want)
			}
		})
	}
}
