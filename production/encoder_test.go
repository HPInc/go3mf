package production

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestMarshalModel(t *testing.T) {
	components := &go3mf.Object{
		ExtensionAttr: go3mf.ExtensionAttr{mustUUID("cb828680-8895-4e08-a1fc-be63e033df15")},
		ID:            20,
		Components: []*go3mf.Component{{
			ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{
				Path: "/3D/other.model",
				UUID: UUID("cb828680-8895-4e08-a1fc-be63e033df16"),
			}},
			ObjectID: 8, Transform: go3mf.Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}},
		},
	}
	m := &go3mf.Model{Path: "/3D/3dmodel.model", Build: go3mf.Build{
		ExtensionAttr: go3mf.ExtensionAttr{mustUUID("e9e25302-6428-402e-8633-cc95528d0ed3")},
	}}
	m.Resources = go3mf.Resources{Objects: []*go3mf.Object{components}}
	m.Build.Items = append(m.Build.Items, &go3mf.Item{ObjectID: 20,
		ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: UUID("e9e25302-6428-402e-8633-cc95528d0ed2")}},
		Transform:     go3mf.Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
	}, &go3mf.Item{ObjectID: 8,
		ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{
			Path: "/3D/other.model",
			UUID: UUID("e9e25302-6428-402e-8633-cc95528d0ed4"),
		}},
	})
	t.Run("base", func(t *testing.T) {
		m.WithExtension(&Extension{LocalName: "p"})
		b, err := go3mf.MarshalModel(m)
		if err != nil {
			t.Errorf("production.MarshalModel() error = %v", err)
			return
		}
		d := go3mf.NewDecoder(nil, 0)
		newModel := new(go3mf.Model)
		newModel.WithExtension(&Extension{LocalName: "p"})
		newModel.Path = m.Path
		if err := d.UnmarshalModel(b, newModel); err != nil {
			t.Errorf("production.MarshalModel() error decoding = %v, s = %s", err, string(b))
			return
		}
		if diff := deep.Equal(m, newModel); diff != nil {
			t.Errorf("production.MarshalModel() = %v, s = %s", diff, string(b))
		}
	})
}
