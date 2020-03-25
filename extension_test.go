package go3mf

import (
	"reflect"
	"testing"
)

func TestAnyAttr_Get(t *testing.T) {
	tests := []struct {
		name   string
		e      AnyAttr
		want   interface{}
		wantOK bool
	}{
		{"nil", nil, new(fakeAttr), false},
		{"empty", AnyAttr{}, new(fakeAttr), false},
		{"non-exist", AnyAttr{nil}, new(fakeAttr), false},
		{"exist", AnyAttr{&fakeAttr{Value: "1"}}, &fakeAttr{Value: "1"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := new(fakeAttr)
			if got := tt.e.Get(&target); got != tt.wantOK {
				t.Errorf("AnyAttr.Get() = %v, wantOK %v", got, tt.wantOK)
				return
			}
			if !reflect.DeepEqual(target, tt.want) {
				t.Errorf("AnyAttr.Get() = %v, want %v", target, tt.want)
			}
		})
	}
}

func TestExtension_Get(t *testing.T) {
	tests := []struct {
		name   string
		e      Extension
		want   interface{}
		wantOK bool
	}{
		{"nil", nil, new(fakeAsset), false},
		{"empty", Extension{}, new(fakeAsset), false},
		{"non-exist", Extension{nil}, new(fakeAsset), false},
		{"exist", Extension{&fakeAsset{ID: 1}}, &fakeAsset{ID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := new(fakeAsset)
			if got := tt.e.Get(&target); got != tt.wantOK {
				t.Errorf("Extension.Get() = %v, wantOK %v", got, tt.wantOK)
				return
			}
			if !reflect.DeepEqual(target, tt.want) {
				t.Errorf("Extension.Get() = %v, want %v", target, tt.want)
			}
		})
	}
}

func TestExtension_Get_Panic(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name string
		e    Extension
		args args
	}{
		{"nil", Extension{&fakeAsset{ID: 1}}, args{nil}},
		{"int", Extension{&fakeAsset{ID: 1}}, args{1}},
		{"nonPtrToPtr", Extension{&fakeAsset{ID: 1}}, args{new(fakeAsset)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Error("Extension.Get() did not panic")
				}
			}()
			if tt.e.Get(tt.args.target) {
				t.Error("Extension.Get() want false")
			}
		})
	}
}

func TestAnyAttr_Get_Panic(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name string
		e    AnyAttr
		args args
	}{
		{"nil", AnyAttr{&fakeAttr{Value: "1"}}, args{nil}},
		{"int", AnyAttr{&fakeAttr{Value: "1"}}, args{1}},
		{"nonPtrToPtr", AnyAttr{&fakeAttr{Value: "1"}}, args{new(fakeAttr)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Error("AnyAttr.Get() did not panic")
				}
			}()
			if tt.e.Get(tt.args.target) {
				t.Error("AnyAttr.Get() want false")
			}
		})
	}
}
