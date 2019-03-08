package model

import (
	"image/color"
	"reflect"
	"testing"
)

func TestBaseMaterial_ColotString(t *testing.T) {
	tests := []struct {
		name string
		m    *BaseMaterial
		want string
	}{
		{"base", &BaseMaterial{Color: color.RGBA{200, 250, 60, 80}}, "#c8fa3c50"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ColorString(); got != tt.want {
				t.Errorf("BaseMaterial.ColotString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseMaterialsResource_Merge(t *testing.T) {
	type args struct {
		other []*BaseMaterial
	}
	tests := []struct {
		name string
		ms   *BaseMaterialsResource
		args args
	}{
		{"base", &BaseMaterialsResource{Materials: []*BaseMaterial{{Name: "1", Color: color.RGBA{200, 250, 60, 80}}}}, args{
			[]*BaseMaterial{{Name: "2", Color: color.RGBA{200, 250, 60, 80}}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := append(tt.ms.Materials, tt.args.other...)
			tt.ms.Merge(tt.args.other)
			if !reflect.DeepEqual(tt.ms.Materials, want) {
				t.Errorf("BaseMaterialsResource.Merge() = %v, want %v", tt.ms.Materials, want)
			}
		})
	}
}
