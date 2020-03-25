package go3mf

import (
	"reflect"
	"testing"
)

func TestAttrMarshalers_Get(t *testing.T) {
	tests := []struct {
		name   string
		e      AttrMarshalers
		want   interface{}
		wantOK bool
	}{
		{"nil", nil, new(fakeAttr), false},
		{"empty", AttrMarshalers{}, new(fakeAttr), false},
		{"non-exist", AttrMarshalers{nil}, new(fakeAttr), false},
		{"exist", AttrMarshalers{&fakeAttr{Value: "1"}}, &fakeAttr{Value: "1"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := new(fakeAttr)
			if got := tt.e.Get(&target); got != tt.wantOK {
				t.Errorf("AttrMarshalers.Get() = %v, wantOK %v", got, tt.wantOK)
				return
			}
			if !reflect.DeepEqual(target, tt.want) {
				t.Errorf("AttrMarshalers.Get() = %v, want %v", target, tt.want)
			}
		})
	}
}

func TestMarshalers_Get(t *testing.T) {
	tests := []struct {
		name   string
		e      Marshalers
		want   interface{}
		wantOK bool
	}{
		{"nil", nil, new(fakeAsset), false},
		{"empty", Marshalers{}, new(fakeAsset), false},
		{"non-exist", Marshalers{nil}, new(fakeAsset), false},
		{"exist", Marshalers{&fakeAsset{ID: 1}}, &fakeAsset{ID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := new(fakeAsset)
			if got := tt.e.Get(&target); got != tt.wantOK {
				t.Errorf("Marshalers.Get() = %v, wantOK %v", got, tt.wantOK)
				return
			}
			if !reflect.DeepEqual(target, tt.want) {
				t.Errorf("Marshalers.Get() = %v, want %v", target, tt.want)
			}
		})
	}
}

func TestMarshalers_Get_Panic(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name string
		e    Marshalers
		args args
	}{
		{"nil", Marshalers{&fakeAsset{ID: 1}}, args{nil}},
		{"int", Marshalers{&fakeAsset{ID: 1}}, args{1}},
		{"nonPtrToPtr", Marshalers{&fakeAsset{ID: 1}}, args{new(fakeAsset)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Error("Marshalers.Get() did not panic")
				}
			}()
			if tt.e.Get(tt.args.target) {
				t.Error("Marshalers.Get() want false")
			}
		})
	}
}

func TestAttrMarshalers_Get_Panic(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name string
		e    AttrMarshalers
		args args
	}{
		{"nil", AttrMarshalers{&fakeAttr{Value: "1"}}, args{nil}},
		{"int", AttrMarshalers{&fakeAttr{Value: "1"}}, args{1}},
		{"nonPtrToPtr", AttrMarshalers{&fakeAttr{Value: "1"}}, args{new(fakeAttr)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Error("AttrMarshalers.Get() did not panic")
				}
			}()
			if tt.e.Get(tt.args.target) {
				t.Error("AttrMarshalers.Get() want false")
			}
		})
	}
}
