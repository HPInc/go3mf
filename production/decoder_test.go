package production

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestDecode(t *testing.T) {
	components := &go3mf.ObjectResource{
		Extensions: map[string]interface{}{ExtensionName: UUID("cb828680-8895-4e08-a1fc-be63e033df15")},
		ID:         20, ModelPath: "/3d/3dmodel.model",
		Components: []*go3mf.Component{{
			Path:       "/3d/other.model",
			Extensions: map[string]interface{}{ExtensionName: UUID("cb828680-8895-4e08-a1fc-be63e033df16")},
			ObjectID:   8, Transform: go3mf.Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}},
		},
	}

	want := &go3mf.Model{Path: "/3d/3dmodel.model"}
	otherMesh := &go3mf.ObjectResource{Mesh: new(go3mf.Mesh), ID: 8, ModelPath: "/3d/other.model"}
	want.Resources = append(want.Resources, otherMesh, components)
	SetBuildUUID(&want.Build, MustNewUUID("e9e25302-6428-402e-8633-cc95528d0ed3"))
	want.Build.Items = append(want.Build.Items, &go3mf.Item{ObjectID: 20,
		Extensions: map[string]interface{}{ExtensionName: UUID("e9e25302-6428-402e-8633-cc95528d0ed2")},
		Transform:  go3mf.Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
	}, &go3mf.Item{ObjectID: 8, Path: "/3d/other.model",
		Extensions: map[string]interface{}{ExtensionName: UUID("e9e25302-6428-402e-8633-cc95528d0ed4")},
	})
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
	got.Resources = append(got.Resources, otherMesh)
	rootFile := `
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:p="http://schemas.microsoft.com/3dmanufacturing/production/2015/06">
		<resources>
			<object id="20" p:UUID="cb828680-8895-4e08-a1fc-be63e033df15">
				<components>
					<component objectid="8" p:UUID="cb828680-8895-4e08-a1fc-be63e033df16" p:path="/3d/other.model" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
				</components>
			</object>
		</resources>
		<build p:UUID="e9e25302-6428-402e-8633-cc95528d0ed3">
			<item objectid="20" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed2" transform="1 0 0 0 2 0 0 0 3 -66.4 -87.1 8.8" />
			<item objectid="8" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed4" p:path="/3d/other.model" />
		</build>
		</model>
		`
	t.Run("base", func(t *testing.T) {
		d := new(go3mf.Decoder)
		RegisterExtension(d)
		d.Strict = true
		if err := d.DecodeRawModel(context.Background(), got, rootFile); err != nil {
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
		go3mf.ParsePropertyError{ResourceID: 20, Element: "object", ModelPath: "/3d/3dmodel.model", Name: "UUID", Value: "cb8286808895-4e08-a1fc-be63e033df15", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Name: "UUID", Value: "cb8286808895-4e08-a1fc-be63e033df16", Type: go3mf.PropertyRequired},
		//go3mf.MissingPropertyError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Name: "UUID"},
		//go3mf.MissingPropertyError{ResourceID: 0, Element: "build", ModelPath: "/3d/3dmodel.model", Name: "UUID"},
		//go3mf.MissingPropertyError{ResourceID: 8, Element: "item", ModelPath: "/3d/3dmodel.model", Name: "UUID"},
		go3mf.ParsePropertyError{ResourceID: 0, Element: "build", Name: "UUID", Value: "e9e25302-6428-402e-8633ed2", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 20, Element: "item", ModelPath: "/3d/3dmodel.model", Name: "UUID", Value: "invalid-uuid", Type: go3mf.PropertyRequired},
	}
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
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
			<item objectid="8" p:path="/3d/other.model"/>
			<item objectid="5" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed4"/>
		</build>
		</model>`

	t.Run("base", func(t *testing.T) {
		d := new(go3mf.Decoder)
		RegisterExtension(d)
		d.Strict = false
		if err := d.DecodeRawModel(context.Background(), got, rootFile); err != nil {
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

func Test_fileFilter(t *testing.T) {
	type args struct {
		relType     string
		isRootModel bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"accepted", args{go3mf.RelTypeModel3D, true}, true},
		{"rejected-noroot", args{go3mf.RelTypeModel3D, false}, false},
		{"rejected-nomodel3d", args{"other", true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileFilter(tt.args.relType, tt.args.isRootModel); got != tt.want {
				t.Errorf("fileFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
