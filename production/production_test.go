// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package production

import (
	"testing"

	"github.com/hpinc/go3mf"
	"github.com/hpinc/go3mf/spec"
)

var _ spec.Marshaler = new(BuildAttr)
var _ spec.Marshaler = new(ItemAttr)
var _ spec.Marshaler = new(ComponentAttr)
var _ spec.Marshaler = new(ObjectAttr)

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
func TestSetMissingUUIDs(t *testing.T) {
	components := &go3mf.Object{
		ID:         20,
		Components: &go3mf.Components{Component: []*go3mf.Component{{ObjectID: 8}}},
	}
	m := &go3mf.Model{Path: "/3D/3dmodel.model", Build: go3mf.Build{}}
	m.Resources = go3mf.Resources{Objects: []*go3mf.Object{components}}
	m.Build.Items = append(m.Build.Items, &go3mf.Item{ObjectID: 20}, &go3mf.Item{ObjectID: 8})
	SetMissingUUIDs(m)
	if len(m.Build.AnyAttr) == 0 {
		t.Errorf("SetMissingUUIDs() should have filled build attrs")
	}
	if len(m.Build.Items[0].AnyAttr) == 0 {
		t.Errorf("SetMissingUUIDs() should have filled item attrs")
	}
	if len(m.Build.Items[1].AnyAttr) == 0 {
		t.Errorf("SetMissingUUIDs() should have filled item attrs")
	}
	if len(m.Resources.Objects[0].AnyAttr) == 0 {
		t.Errorf("SetMissingUUIDs() should have filled object attrs")
	}
	if len(m.Resources.Objects[0].Components.Component[0].AnyAttr) == 0 {
		t.Errorf("SetMissingUUIDs() should have filled object attrs")
	}
}
