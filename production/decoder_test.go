package production

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

func mustUUID(u string) *UUID {
	v := UUID(u)
	return &v
}

func TestDecode(t *testing.T) {
	components := &go3mf.Object{
		AnyAttr: go3mf.AttrMarshalers{mustUUID("cb828680-8895-4e08-a1fc-be63e033df15")},
		ID:      20,
		Components: []*go3mf.Component{{
			AnyAttr: go3mf.AttrMarshalers{&PathUUID{
				Path: "/3D/other.model",
				UUID: UUID("cb828680-8895-4e08-a1fc-be63e033df16"),
			}},
			ObjectID: 8, Transform: go3mf.Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}},
		},
	}

	want := &go3mf.Model{Path: "/3D/3dmodel.model", Resources: go3mf.Resources{
		Objects: []*go3mf.Object{components},
	}, Build: go3mf.Build{AnyAttr: go3mf.AttrMarshalers{mustUUID("e9e25302-6428-402e-8633-cc95528d0ed3")}},
	}
	want.Build.Items = append(want.Build.Items, &go3mf.Item{ObjectID: 20,
		AnyAttr:   go3mf.AttrMarshalers{&PathUUID{UUID: UUID("e9e25302-6428-402e-8633-cc95528d0ed2")}},
		Transform: go3mf.Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
	}, &go3mf.Item{ObjectID: 8,
		AnyAttr: go3mf.AttrMarshalers{&PathUUID{
			Path: "/3D/other.model",
			UUID: UUID("e9e25302-6428-402e-8633-cc95528d0ed4"),
		}},
	})
	got := new(go3mf.Model)
	got.Path = "/3D/3dmodel.model"
	rootFile := `
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:p="http://schemas.microsoft.com/3dmanufacturing/production/2015/06">
		<resources>
			<object id="20" p:UUID="cb828680-8895-4e08-a1fc-be63e033df15">
				<components>
					<component objectid="8" p:UUID="cb828680-8895-4e08-a1fc-be63e033df16" p:path="/3D/other.model" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
				</components>
			</object>
		</resources>
		<build p:UUID="e9e25302-6428-402e-8633-cc95528d0ed3">
			<item objectid="20" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed2" transform="1 0 0 0 2 0 0 0 3 -66.4 -87.1 8.8" />
			<item objectid="8" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed4" p:path="/3D/other.model" />
		</build>
		</model>
		`
	t.Run("base", func(t *testing.T) {
		want.WithSpec(&Spec{LocalName: "p"})
		got.WithSpec(&Spec{LocalName: "p"})
		d := new(go3mf.Decoder)
		d.Strict = true
		if err := d.UnmarshalModel([]byte(rootFile), got); err != nil {
			t.Errorf("DecodeRawModel() unexpected error = %v", err)
			return
		}
		deep.CompareUnexportedFields = true
		deep.MaxDepth = 20
		if diff := deep.Equal(got, want); diff != nil {
			t.Errorf("DecodeRawModell() = %v", diff)
			return
		}
	})
}

func TestDecode_warns(t *testing.T) {
	want := []error{
		&specerr.ParseFieldError{ResourceID: 20, Element: "object", ModelPath: "/3D/3dmodel.model", Name: "UUID", Value: "cb8286808895-4e08-a1fc-be63e033df15", Required: true},
		&specerr.ParseFieldError{ResourceID: 20, Element: "component", ModelPath: "/3D/3dmodel.model", Name: "UUID", Value: "cb8286808895-4e08-a1fc-be63e033df16", Required: true},
		&specerr.ParseFieldError{ResourceID: 0, Element: "build", Name: "UUID", Value: "e9e25302-6428-402e-8633ed2", ModelPath: "/3D/3dmodel.model", Required: true},
		&specerr.ParseFieldError{ResourceID: 20, Element: "item", ModelPath: "/3D/3dmodel.model", Name: "UUID", Value: "invalid-uuid", Required: true},
	}
	got := new(go3mf.Model)
	got.Path = "/3D/3dmodel.model"
	rootFile := `
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:p="http://schemas.microsoft.com/3dmanufacturing/production/2015/06">
		<resources>
			<object id="22" p:UUID="cb828680-8895-4e08-a1fc-be63e033df15" />
			<object id="20" p:UUID="cb8286808895-4e08-a1fc-be63e033df15">
				<components>
					<component objectid="8" p:path="/2d/2d.model" p:UUID="cb8286808895-4e08-a1fc-be63e033df16"/>
					<component objectid="5" p:UUID="cb828680-8895-4e08-a1fc-be63e033df16"/>
				</components>
			</object>
		</resources>
		<build p:UUID="e9e25302-6428-402e-8633ed2">
			<item partnumber="bob" objectid="20" p:UUID="invalid-uuid" />
			<item objectid="8" p:path="/3D/other.model"/>
			<item objectid="5" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed4"/>
		</build>
		</model>`

	t.Run("base", func(t *testing.T) {
		got.WithSpec(&Spec{LocalName: "p"})
		d := new(go3mf.Decoder)
		d.Strict = false
		if err := d.UnmarshalModel([]byte(rootFile), got); err != nil {
			t.Errorf("DecodeRawModel_warn() unexpected error = %v", err)
			return
		}
		deep.MaxDiff = 1
		if diff := deep.Equal(d.Warnings, want); diff != nil {
			t.Errorf("DecodeRawModel_warn() = %v", diff)
			return
		}
	})
}
