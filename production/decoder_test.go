package production

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestDecode(t *testing.T) {
	meshRes := &go3mf.Mesh{
		ObjectResource: go3mf.ObjectResource{
			ID: 8, Name: "Box 1", ModelPath: "/3d/3dmodel.model",
			Extensions: map[string]interface{}{ExtensionName: &ObjectAttr{UUID: "11111111-1111-1111-1111-111111111111"}},
		},
	}
	meshRes.Nodes = append(meshRes.Nodes, []go3mf.Point3D{
		{0, 0, 0},
		{100, 0, 0},
		{100, 100, 0},
		{0, 100, 0},
		{0, 0, 100},
		{100, 0, 100},
		{100, 100, 100},
		{0, 100, 100},
	}...)
	meshRes.Faces = append(meshRes.Faces, []go3mf.Face{
		{NodeIndices: [3]uint32{3, 2, 1}},
		{NodeIndices: [3]uint32{1, 0, 3}},
		{NodeIndices: [3]uint32{4, 5, 6}},
		{NodeIndices: [3]uint32{6, 7, 4}},
		{NodeIndices: [3]uint32{0, 1, 5}},
		{NodeIndices: [3]uint32{5, 4, 0}},
		{NodeIndices: [3]uint32{1, 2, 6}},
		{NodeIndices: [3]uint32{6, 5, 1}},
		{NodeIndices: [3]uint32{2, 3, 7}},
		{NodeIndices: [3]uint32{7, 6, 2}},
		{NodeIndices: [3]uint32{3, 0, 4}},
		{NodeIndices: [3]uint32{4, 7, 3}},
	}...)

	components := &go3mf.Components{
		ObjectResource: go3mf.ObjectResource{
			Extensions: map[string]interface{}{ExtensionName: &ObjectAttr{UUID: "cb828680-8895-4e08-a1fc-be63e033df15"}},
			ID:         20, ModelPath: "/3d/3dmodel.model",
		},
		Components: []*go3mf.Component{{
			Extensions: map[string]interface{}{ExtensionName: &ComponentAttr{UUID: "cb828680-8895-4e08-a1fc-be63e033df16"}},
			ObjectID:   8, Transform: go3mf.Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}},
		},
	}

	want := &go3mf.Model{Path: "/3d/3dmodel.model"}
	otherMesh := &go3mf.Mesh{ObjectResource: go3mf.ObjectResource{ID: 8, ModelPath: "/3d/other.model"}}
	want.Resources = append(want.Resources, otherMesh, meshRes, components)
	ExtensionBuild(&want.Build).UUID = "e9e25302-6428-402e-8633-cc95528d0ed3"
	want.Build.Items = append(want.Build.Items, &go3mf.Item{ObjectID: 20,
		Extensions: map[string]interface{}{ExtensionName: &ItemAttr{UUID: "e9e25302-6428-402e-8633-cc95528d0ed2"}},
		Transform:  go3mf.Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
	}, &go3mf.Item{ObjectID: 8,
		Extensions: map[string]interface{}{ExtensionName: &ItemAttr{UUID: "e9e25302-6428-402e-8633-cc95528d0ed4", Path: "/3d/other.model"}},
	})
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
	got.Resources = append(got.Resources, otherMesh)
	rootFile := `
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:p="http://schemas.microsoft.com/3dmanufacturing/production/2015/06">
		<resources>
			<object id="8" name="Box 1" p:UUID="11111111-1111-1111-1111-111111111111" type="model">
				<mesh>
					<vertices>
						<vertex x="0" y="0" z="0" />
						<vertex x="100.00000" y="0" z="0" />
						<vertex x="100.00000" y="100.00000" z="0" />
						<vertex x="0" y="100.00000" z="0" />
						<vertex x="0" y="0" z="100.00000" />
						<vertex x="100.00000" y="0" z="100.00000" />
						<vertex x="100.00000" y="100.00000" z="100.00000" />
						<vertex x="0" y="100.00000" z="100.00000" />
					</vertices>
					<triangles>
						<triangle v1="3" v2="2" v3="1" />
						<triangle v1="1" v2="0" v3="3" />
						<triangle v1="4" v2="5" v3="6" />
						<triangle v1="6" v2="7" v3="4" />
						<triangle v1="0" v2="1" v3="5" />
						<triangle v1="5" v2="4" v3="0" />
						<triangle v1="1" v2="2" v3="6" />
						<triangle v1="6" v2="5" v3="1" />
						<triangle v1="2" v2="3" v3="7" />
						<triangle v1="7" v2="6" v3="2" />
						<triangle v1="3" v2="0" v3="4" />
						<triangle v1="4" v2="7" v3="3" />
					</triangles>
				</mesh>
			</object>
			<object id="20" p:UUID="cb828680-8895-4e08-a1fc-be63e033df15">
				<components>
					<component objectid="8" p:UUID="cb828680-8895-4e08-a1fc-be63e033df16" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
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
	}
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
	rootFile := `
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:p="http://schemas.microsoft.com/3dmanufacturing/production/2015/06">
		<resources>		
			<object id="8" name="Box 1" pid="5" pindex="0" partnumber="11111111-1111-1111-1111-111111111111" type="model">
				<mesh>
					<vertices>
						<vertex x="0" y="0" z="0" />
						<vertex x="100.00000" y="0" z="0" />
						<vertex x="100.00000" y="100.00000" z="0" />
						<vertex x="0" y="100.00000" z="0" />
						<vertex x="0" y="0" z="100.00000" />
						<vertex x="100.00000" y="0" z="100.00000" />
						<vertex x="100.00000" y="100.00000" z="100.00000" />
						<vertex x="0" y="100.00000" z="100.00000" />
					</vertices>
					<triangles>
						<triangle v1="2" v2="3" v3="1" />
						<triangle v1="1" v2="2" v3="3" />
						<triangle v1="3" v2="2" v3="1" />
						<triangle v1="1" v2="0" v3="3" />
						<triangle v1="4" v2="5" v3="6" />
						<triangle v1="6" v2="7" v3="4" />
						<triangle v1="0" v2="1" v3="5" />
						<triangle v1="5" v2="4" v3="0" />
						<triangle v1="1" v2="2" v3="6" />
						<triangle v1="6" v2="5" v3="1" />
						<triangle v1="2" v2="3" v3="7" />
						<triangle v1="7" v2="6" v3="2" />
						<triangle v1="3" v2="0" v3="4" />
						<triangle v1="4" v2="7" v3="3" />
					</triangles>
				</mesh>
			</object>
			<object id="22" p:UUID="cb828680-8895-4e08-a1fc-be63e033df15" />
			<object id="20" p:UUID="cb8286808895-4e08-a1fc-be63e033df15">
				<components>
					<component objectid="8" p:path="/2d/2d.model" p:UUID="cb8286808895-4e08-a1fc-be63e033df16"/>
					<component objectid="5" p:UUID="cb828680-8895-4e08-a1fc-be63e033df16"/>
				</components>
			</object>
		</resources>
		<build p:UUID="e9e25302-6428-402e-8633ed2">
			<item partnumber="bob" objectid="20" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed2" />
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
