package production

import (
	"testing"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/spec"
)

var _ spec.Decoder = new(Spec)
var _ spec.PostProcessorDecoder = new(Spec)
var _ spec.PreProcessEncoder = new(Spec)
var _ go3mf.Spec = new(Spec)
var _ spec.ModelValidator = new(Spec)
var _ spec.ObjectValidator = new(Spec)
var _ spec.MarshalerAttr = new(BuildAttr)
var _ spec.MarshalerAttr = new(ItemAttr)
var _ spec.ObjectPather = new(ItemAttr)
var _ spec.MarshalerAttr = new(ComponentAttr)
var _ spec.ObjectPather = new(ComponentAttr)
var _ spec.MarshalerAttr = new(ObjectAttr)

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
