package go3mf

import (
	"reflect"
	"testing"
)

func TestPackageResourceID_SetID(t *testing.T) {
	type args struct {
		path string
		id   uint64
	}
	tests := []struct {
		name string
		p    *PackageResourceID
		args args
		want *PackageResourceID
	}{
		{"base", new(PackageResourceID), args{"a", 6}, &PackageResourceID{path: "a", id: 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.SetID(tt.args.path, tt.args.id)
			if !reflect.DeepEqual(tt.p, tt.want) {
				t.Errorf("PackageResourceID.SetID() = %v, want %v", tt.p, tt.want)
			}
		})
	}
}

func TestPackageResourceID_ID(t *testing.T) {
	tests := []struct {
		name  string
		p     *PackageResourceID
		want  string
		want1 uint64
	}{
		{"new", new(PackageResourceID), "", 0},
		{"base", &PackageResourceID{path: "a", id: 6}, "a", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.p.ID()
			if got != tt.want {
				t.Errorf("PackageResourceID.ID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PackageResourceID.ID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPackageResourceID_SetUniqueID(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name string
		p    *PackageResourceID
		args args
	}{
		{"base", new(PackageResourceID), args{8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.SetUniqueID(tt.args.id)
			if tt.p.uniqueID != tt.args.id {
				t.Errorf("PackageResourceID.SetUniqueID() = %v, want %v", tt.p.uniqueID, tt.args.id)
			}
		})
	}
}

func TestPackageResourceID_UniqueID(t *testing.T) {
	tests := []struct {
		name string
		p    *PackageResourceID
		want uint64
	}{
		{"base", new(PackageResourceID), 0},
		{"base", &PackageResourceID{uniqueID: 4}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.UniqueID(); got != tt.want {
				t.Errorf("PackageResourceID.UniqueID() = %v, want %v", got, tt.want)
			}
		})
	}
}
