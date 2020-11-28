package production

import (
	"testing"

	"github.com/qmuntal/go3mf"
)

var _ go3mf.SpecDecoder = new(Spec)
var _ go3mf.SpecValidator = new(Spec)
var _ go3mf.MarshalerAttr = new(BuildAttr)
var _ go3mf.MarshalerAttr = new(ItemAttr)
var _ go3mf.ObjectPather = new(ItemAttr)
var _ go3mf.MarshalerAttr = new(ComponentAttr)
var _ go3mf.ObjectPather = new(ComponentAttr)
var _ go3mf.MarshalerAttr = new(ObjectAttr)

func TestComponentAttr_ObjectPath(t *testing.T) {
	tests := []struct {
		name string
		p    *ComponentAttr
		want string
	}{
		{"empty", new(ComponentAttr), ""},
		{"path", &ComponentAttr{Path: "/a.model"}, "/a.model"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ObjectPath(); got != tt.want {
				t.Errorf("ComponentAttr.ObjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemAttr_ObjectPath(t *testing.T) {
	tests := []struct {
		name string
		p    *ItemAttr
		want string
	}{
		{"empty", new(ItemAttr), ""},
		{"path", &ItemAttr{Path: "/a.model"}, "/a.model"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ObjectPath(); got != tt.want {
				t.Errorf("ItemAttr.ObjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
