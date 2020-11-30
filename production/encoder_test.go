package production

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestSpec_PreProcessEncode_AutoUUIDDisabled(t *testing.T) {
	m := &go3mf.Model{Path: "/3D/3dmodel.model", Build: go3mf.Build{}}
	s := &Spec{LocalName: "p", DisableAutoUUID: true, m: m}
	s.PreProcessEncode()
	if len(m.Build.AnyAttr) != 0 {
		t.Errorf("Spec.PreProcessEncode() shouldn't have filled build attrs")
	}
}

func TestSpec_PreProcessEncode(t *testing.T) {
	components := &go3mf.Object{
		ID:         20,
		Components: []*go3mf.Component{{ObjectID: 8}},
	}
	m := &go3mf.Model{Path: "/3D/3dmodel.model", Build: go3mf.Build{}}
	m.Resources = go3mf.Resources{Objects: []*go3mf.Object{components}}
	m.Build.Items = append(m.Build.Items, &go3mf.Item{ObjectID: 20}, &go3mf.Item{ObjectID: 8})
	s := &Spec{LocalName: "p", m: m}
	s.PreProcessEncode()
	if len(m.Build.AnyAttr) == 0 {
		t.Errorf("Spec.PreProcessEncode() should have filled build attrs")
	}
	if len(m.Build.Items[0].AnyAttr) == 0 {
		t.Errorf("Spec.PreProcessEncode() should have filled item attrs")
	}
	if len(m.Build.Items[1].AnyAttr) == 0 {
		t.Errorf("Spec.PreProcessEncode() should have filled item attrs")
	}
	if len(m.Resources.Objects[0].AnyAttr) == 0 {
		t.Errorf("Spec.PreProcessEncode() should have filled object attrs")
	}
	if len(m.Resources.Objects[0].Components[0].AnyAttr) == 0 {
		t.Errorf("Spec.PreProcessEncode() should have filled object attrs")
	}
}

func TestMarshalModel(t *testing.T) {
	components := &go3mf.Object{
		AnyAttr: go3mf.ExtensionsAttr{&ObjectAttr{UUID: "cb828680-8895-4e08-a1fc-be63e033df15"}},
		ID:      20,
		Components: []*go3mf.Component{{
			ObjectID: 8,
			AnyAttr: go3mf.ExtensionsAttr{&ComponentAttr{
				Path: "/3D/other.model",
				UUID: "cb828680-8895-4e08-a1fc-be63e033df16",
			}}},
		},
	}
	m := &go3mf.Model{Path: "/3D/3dmodel.model", Build: go3mf.Build{
		AnyAttr: go3mf.ExtensionsAttr{&BuildAttr{UUID: "e9e25302-6428-402e-8633-cc95528d0ed3"}},
	}}
	m.Resources = go3mf.Resources{Objects: []*go3mf.Object{components}}
	m.Build.Items = append(m.Build.Items, &go3mf.Item{ObjectID: 20,
		AnyAttr: go3mf.ExtensionsAttr{&ItemAttr{UUID: "e9e25302-6428-402e-8633-cc95528d0ed2"}},
	}, &go3mf.Item{ObjectID: 8,
		AnyAttr: go3mf.ExtensionsAttr{&ItemAttr{
			Path: "/3D/other.model",
			UUID: "e9e25302-6428-402e-8633-cc95528d0ed4",
		}},
	})
	m.WithSpec(&Spec{LocalName: "p"})
	b, err := go3mf.MarshalModel(m)
	if err != nil {
		t.Errorf("production.MarshalModel() error = %v", err)
		return
	}
	newModel := new(go3mf.Model)
	newModel.WithSpec(&Spec{LocalName: "p"})
	newModel.Path = m.Path
	if err := go3mf.UnmarshalModel(b, newModel); err != nil {
		t.Errorf("production.MarshalModel() error decoding = %v, s = %s", err, string(b))
		return
	}
	if diff := deep.Equal(m, newModel); diff != nil {
		t.Errorf("production.MarshalModel() = %v, s = %s", diff, string(b))
	}
}
