package go3mf

import (
	"reflect"
	"testing"
)

func TestExtensionsAttr_Get(t *testing.T) {
	tests := []struct {
		name   string
		e      ExtensionsAttr
		want   interface{}
		wantOK bool
	}{
		{"nil", nil, new(fakeAttr), false},
		{"empty", ExtensionsAttr{}, new(fakeAttr), false},
		{"non-exist", ExtensionsAttr{nil}, new(fakeAttr), false},
		{"exist", ExtensionsAttr{&fakeAttr{Value: "1"}}, &fakeAttr{Value: "1"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := new(fakeAttr)
			if got := tt.e.Get(&target); got != tt.wantOK {
				t.Errorf("ExtensionsAttr.Get() = %v, wantOK %v", got, tt.wantOK)
				return
			}
			if !reflect.DeepEqual(target, tt.want) {
				t.Errorf("ExtensionsAttr.Get() = %v, want %v", target, tt.want)
			}
		})
	}
}

func TestExtensions_Get(t *testing.T) {
	tests := []struct {
		name   string
		e      Extensions
		want   interface{}
		wantOK bool
	}{
		{"nil", nil, new(fakeAsset), false},
		{"empty", Extensions{}, new(fakeAsset), false},
		{"non-exist", Extensions{nil}, new(fakeAsset), false},
		{"exist", Extensions{&fakeAsset{ID: 1}}, &fakeAsset{ID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := new(fakeAsset)
			if got := tt.e.Get(&target); got != tt.wantOK {
				t.Errorf("Extensions.Get() = %v, wantOK %v", got, tt.wantOK)
				return
			}
			if !reflect.DeepEqual(target, tt.want) {
				t.Errorf("Extensions.Get() = %v, want %v", target, tt.want)
			}
		})
	}
}

func TestExtensions_Get_Panic(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name string
		e    Extensions
		args args
	}{
		{"nil", Extensions{&fakeAsset{ID: 1}}, args{nil}},
		{"int", Extensions{&fakeAsset{ID: 1}}, args{1}},
		{"nonPtrToPtr", Extensions{&fakeAsset{ID: 1}}, args{new(fakeAsset)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Error("Extensions.Get() did not panic")
				}
			}()
			if tt.e.Get(tt.args.target) {
				t.Error("Extensions.Get() want false")
			}
		})
	}
}

func TestExtensionsAttr_Get_Panic(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name string
		e    ExtensionsAttr
		args args
	}{
		{"nil", ExtensionsAttr{&fakeAttr{Value: "1"}}, args{nil}},
		{"int", ExtensionsAttr{&fakeAttr{Value: "1"}}, args{1}},
		{"nonPtrToPtr", ExtensionsAttr{&fakeAttr{Value: "1"}}, args{new(fakeAttr)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Error("ExtensionsAttr.Get() did not panic")
				}
			}()
			if tt.e.Get(tt.args.target) {
				t.Error("ExtensionsAttr.Get() want false")
			}
		})
	}
}
