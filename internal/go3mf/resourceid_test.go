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
		{"new", new(PackageResourceID), 0},
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

func TestNewResourceHandler(t *testing.T) {
	tests := []struct {
		name string
		want *ResourceHandler
	}{
		{"base", &ResourceHandler{
			resourceIDs: make(map[uint64]*PackageResourceID, 0),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceHandler_FindResourceID(t *testing.T) {
	rh := NewResourceHandler()
	r1, _ := rh.NewResourceID("a", 11)
	r2, _ := rh.NewResourceID("b", 12)
	type args struct {
		uniqueID uint64
	}
	tests := []struct {
		name    string
		r       *ResourceHandler
		args    args
		wantVal *PackageResourceID
		wantOk  bool
	}{
		{"nook", rh, args{123}, nil, false},
		{"r1", rh, args{1}, r1, true},
		{"r2", rh, args{2}, r2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := tt.r.FindResourceID(tt.args.uniqueID)
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("ResourceHandler.FindResourceID() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOk != tt.wantOk {
				t.Errorf("ResourceHandler.FindResourceID() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestResourceHandler_FindResourceIDByID(t *testing.T) {
	rh := NewResourceHandler()
	r1, _ := rh.NewResourceID("a", 11)
	r2, _ := rh.NewResourceID("b", 12)
	type args struct {
		path string
		id   uint64
	}
	tests := []struct {
		name    string
		r       *ResourceHandler
		args    args
		wantVal *PackageResourceID
		wantOk  bool
	}{
		{"nook", rh, args{"abc", 11}, nil, false},
		{"r1", rh, args{"a", 11}, r1, true},
		{"r2", rh, args{"b", 12}, r2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := tt.r.FindResourceIDByID(tt.args.path, tt.args.id)
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("ResourceHandler.FindResourceIDByID() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOk != tt.wantOk {
				t.Errorf("ResourceHandler.FindResourceIDByID() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestResourceHandler_NewResourceID(t *testing.T) {
	rh := NewResourceHandler()
	type args struct {
		path string
		id   uint64
	}
	tests := []struct {
		name    string
		r       *ResourceHandler
		args    args
		want    *PackageResourceID
		wantErr bool
	}{
		{"add1", rh, args{"a", 12}, &PackageResourceID{"a", 12, 1}, false},
		{"add2", rh, args{"b", 13}, &PackageResourceID{"b", 13, 2}, false},
		{"err", rh, args{"b", 13}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.NewResourceID(tt.args.path, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceHandler.NewResourceID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceHandler.NewResourceID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceHandler_Clear(t *testing.T) {
	rh := NewResourceHandler()
	rh.NewResourceID("a", 11)
	rh.NewResourceID("b", 12)
	tests := []struct {
		name string
		r    *ResourceHandler
	}{
		{"new", NewResourceHandler()},
		{"base", rh},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Clear()
			if len(tt.r.resourceIDs) != 0 {
				t.Error("ResourceHandler.Clear() should clear uniqueIDs")
			}
		})
	}
}
