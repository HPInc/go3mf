package model

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
		p    *ResourceID
		args args
		want *ResourceID
	}{
		{"base", new(ResourceID), args{"a", 6}, &ResourceID{path: "a", id: 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.SetID(tt.args.path, tt.args.id)
			if !reflect.DeepEqual(tt.p, tt.want) {
				t.Errorf("ResourceID.SetID() = %v, want %v", tt.p, tt.want)
			}
		})
	}
}

func TestPackageResourceID_ID(t *testing.T) {
	tests := []struct {
		name  string
		p     *ResourceID
		want  string
		want1 uint64
	}{
		{"new", new(ResourceID), "", 0},
		{"base", &ResourceID{path: "a", id: 6}, "a", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.p.ID()
			if got != tt.want {
				t.Errorf("ResourceID.ID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ResourceID.ID() got1 = %v, want %v", got1, tt.want1)
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
		p    *ResourceID
		args args
	}{
		{"base", new(ResourceID), args{8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.SetUniqueID(tt.args.id)
			if tt.p.uniqueID != tt.args.id {
				t.Errorf("ResourceID.SetUniqueID() = %v, want %v", tt.p.uniqueID, tt.args.id)
			}
		})
	}
}

func TestPackageResourceID_UniqueID(t *testing.T) {
	tests := []struct {
		name string
		p    *ResourceID
		want uint64
	}{
		{"new", new(ResourceID), 0},
		{"base", &ResourceID{uniqueID: 4}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.UniqueID(); got != tt.want {
				t.Errorf("ResourceID.UniqueID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newResourceHandler(t *testing.T) {
	tests := []struct {
		name string
		want *resourceHandler
	}{
		{"base", &resourceHandler{
			resourceIDs: make(map[uint64]*ResourceID, 0),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newResourceHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newResourceHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceHandler_FindResourceID(t *testing.T) {
	rh := newResourceHandler()
	r1, _ := rh.NewResourceID("a", 11)
	r2, _ := rh.NewResourceID("b", 12)
	type args struct {
		uniqueID uint64
	}
	tests := []struct {
		name    string
		r       *resourceHandler
		args    args
		wantVal *ResourceID
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
				t.Errorf("resourceHandler.FindResourceID() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOk != tt.wantOk {
				t.Errorf("resourceHandler.FindResourceID() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_resourceHandler_FindResourcePath(t *testing.T) {
	rh := newResourceHandler()
	r1, _ := rh.NewResourceID("a", 11)
	r2, _ := rh.NewResourceID("b", 12)
	type args struct {
		path string
		id   uint64
	}
	tests := []struct {
		name    string
		r       *resourceHandler
		args    args
		wantVal *ResourceID
		wantOk  bool
	}{
		{"nook", rh, args{"abc", 11}, nil, false},
		{"r1", rh, args{"a", 11}, r1, true},
		{"r2", rh, args{"b", 12}, r2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := tt.r.FindResourcePath(tt.args.path, tt.args.id)
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("resourceHandler.FindResourcePath() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOk != tt.wantOk {
				t.Errorf("resourceHandler.FindResourcePath() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_resourceHandler_NewResourceID(t *testing.T) {
	rh := newResourceHandler()
	type args struct {
		path string
		id   uint64
	}
	tests := []struct {
		name    string
		r       *resourceHandler
		args    args
		want    *ResourceID
		wantErr bool
	}{
		{"add1", rh, args{"a", 12}, &ResourceID{"a", 12, 1}, false},
		{"add2", rh, args{"b", 13}, &ResourceID{"b", 13, 2}, false},
		{"err", rh, args{"b", 13}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.NewResourceID(tt.args.path, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("resourceHandler.NewResourceID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceHandler.NewResourceID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceHandler_Clear(t *testing.T) {
	rh := newResourceHandler()
	rh.NewResourceID("a", 11)
	rh.NewResourceID("b", 12)
	tests := []struct {
		name string
		r    *resourceHandler
	}{
		{"new", newResourceHandler()},
		{"base", rh},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Clear()
			if len(tt.r.resourceIDs) != 0 {
				t.Error("resourceHandler.Clear() should clear uniqueIDs")
			}
		})
	}
}
