package production

import (
	"encoding/xml"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestMarshalModel(t *testing.T) {
	components := &go3mf.ObjectResource{
		ExtensionAttr: go3mf.ExtensionAttr{ExtensionName: mustUUID("cb828680-8895-4e08-a1fc-be63e033df15")},
		ID:            20, ModelPath: "/3D/3dmodel.model",
		Components: []*go3mf.Component{{
			ExtensionAttr: go3mf.ExtensionAttr{ExtensionName: &PathUUID{
				Path: "/3D/other.model",
				UUID: UUID("cb828680-8895-4e08-a1fc-be63e033df16"),
			}},
			ObjectID: 8, Transform: go3mf.Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}},
		},
	}
	*ObjectAttr(components) = UUID("cb828680-8895-4e08-a1fc-be63e033df15")
	m := &go3mf.Model{Path: "/3D/3dmodel.model", Namespaces: []xml.Name{{Space: ExtensionName, Local: "p"}}}
	m.Resources = append(m.Resources, components)
	*BuildAttr(&m.Build) = UUID("e9e25302-6428-402e-8633-cc95528d0ed3")
	m.Build.Items = append(m.Build.Items, &go3mf.Item{ObjectID: 20,
		ExtensionAttr: go3mf.ExtensionAttr{ExtensionName: &PathUUID{UUID: UUID("e9e25302-6428-402e-8633-cc95528d0ed2")}},
		Transform:     go3mf.Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
	}, &go3mf.Item{ObjectID: 8,
		ExtensionAttr: go3mf.ExtensionAttr{ExtensionName: &PathUUID{
			Path: "/3D/other.model",
			UUID: UUID("e9e25302-6428-402e-8633-cc95528d0ed4"),
		}},
	})
	t.Run("base", func(t *testing.T) {
		b, err := go3mf.MarshalModel(m)
		if err != nil {
			t.Errorf("production.MarshalModel() error = %v", err)
			return
		}
		d := go3mf.NewDecoder(nil, 0)
		RegisterExtension(d)
		newModel := new(go3mf.Model)
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
