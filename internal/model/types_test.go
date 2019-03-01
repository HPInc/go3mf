package model

import (
	"reflect"
	"testing"
)

func TestNewUnits(t *testing.T) {
	tests := []struct {
		name  string
		want  Units
		want1 bool
	}{
		{"micron", Micrometer, true},
		{"millimeter", Millimeter, true},
		{"centimeter", Centimeter, true},
		{"inch", Inch, true},
		{"foot", Foot, true},
		{"meter", Meter, true},
		{"empty", Units(""), false},
		{"none", Units(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewUnits(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnits() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NewUnits() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
