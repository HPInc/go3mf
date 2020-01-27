package io3mf

import (
	"reflect"
	"testing"

	go3mf "github.com/qmuntal/go3mf"
)

func Test_newObjectType(t *testing.T) {
	tests := []struct {
		name   string
		wantO  go3mf.ObjectType
		wantOk bool
	}{
		{"model", go3mf.ObjectTypeModel, true},
		{"other", go3mf.ObjectTypeOther, true},
		{"support", go3mf.ObjectTypeSupport, true},
		{"solidsupport", go3mf.ObjectTypeSolidSupport, true},
		{"surface", go3mf.ObjectTypeSurface, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, gotOk := newObjectType(tt.name)
			if !reflect.DeepEqual(gotO, tt.wantO) {
				t.Errorf("newObjectType() gotO = %v, want %v", gotO, tt.wantO)
			}
			if gotOk != tt.wantOk {
				t.Errorf("newObjectType() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_newUnits(t *testing.T) {
	tests := []struct {
		name  string
		want  go3mf.Units
		want1 bool
	}{
		{"micron", go3mf.UnitMicrometer, true},
		{"millimeter", go3mf.UnitMillimeter, true},
		{"centimeter", go3mf.UnitCentimeter, true},
		{"inch", go3mf.UnitInch, true},
		{"foot", go3mf.UnitFoot, true},
		{"meter", go3mf.UnitMeter, true},
		{"", go3mf.UnitMillimeter, false},
		{"other", go3mf.UnitMillimeter, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := newUnits(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUnits() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("newUnits() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
