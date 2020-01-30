// Package production handles new non-object resources,
// as well as attributes to the build section for uniquely identifying parts within a particular 3MF package
// Despite item and component paths are production attributes, they are also handled by
// the core package, to avoid duplications they won't be stored in the Extension map
// but the core properties will be updated.
package production

import (
	"reflect"
	"testing"

	"github.com/qmuntal/go3mf"
)

func TestBuildUUID(t *testing.T) {
	type args struct {
		b *go3mf.Build
	}
	tests := []struct {
		name string
		args args
		want UUID
	}{
		{"exists", args{&go3mf.Build{Extensions: go3mf.Extensions{ExtensionName: UUID("fake")}}}, "fake"},
		{"no-exists", args{&go3mf.Build{}}, UUID("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildUUID(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemUUID(t *testing.T) {
	type args struct {
		o *go3mf.Item
	}
	tests := []struct {
		name string
		args args
		want UUID
	}{
		{"exists", args{&go3mf.Item{Extensions: go3mf.Extensions{ExtensionName: UUID("fake")}}}, "fake"},
		{"no-exists", args{&go3mf.Item{}}, UUID("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ItemUUID(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponentUUID(t *testing.T) {
	type args struct {
		c *go3mf.Component
	}
	tests := []struct {
		name string
		args args
		want UUID
	}{
		{"exists", args{&go3mf.Component{Extensions: go3mf.Extensions{ExtensionName: UUID("fake")}}}, "fake"},
		{"no-exists", args{&go3mf.Component{}}, UUID("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComponentUUID(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComponentUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectUUID(t *testing.T) {
	type args struct {
		o *go3mf.ObjectResource
	}
	tests := []struct {
		name string
		args args
		want UUID
	}{
		{"exists", args{&go3mf.ObjectResource{Extensions: go3mf.Extensions{ExtensionName: UUID("fake")}}}, "fake"},
		{"no-exists", args{&go3mf.ObjectResource{}}, UUID("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ObjectUUID(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
