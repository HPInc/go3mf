package booleanoperations

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestDecode(t *testing.T) {
	validMesh1 := &go3mf.Object{
		ID:   1,
		Name: "shuttle",
		Type: go3mf.ObjectTypeModel,
		Mesh: &go3mf.Mesh{Vertices: []go3mf.Point3D{{45, 55, 55}, {45, 45, 55}, {45, 55, 45}, {45, 45, 45}, {55, 55, 45}, {55, 55, 55}, {55, 45, 55}, {55, 45, 45}},
			Triangles: []go3mf.Triangle{go3mf.NewTriangle(0, 1, 2), go3mf.NewTriangle(0, 3, 1), go3mf.NewTriangle(0, 2, 3), go3mf.NewTriangle(1, 3, 2)}}}
	validMesh2 := &go3mf.Object{
		ID:   2,
		Name: "label 1",
		Type: go3mf.ObjectTypeModel,
		Mesh: &go3mf.Mesh{Vertices: []go3mf.Point3D{{45, 55, 55}, {45, 45, 55}, {45, 55, 45}, {45, 45, 45}, {55, 55, 45}, {55, 55, 55}, {55, 45, 55}, {55, 45, 45}},
			Triangles: []go3mf.Triangle{go3mf.NewTriangle(0, 1, 2), go3mf.NewTriangle(0, 3, 1), go3mf.NewTriangle(0, 2, 3), go3mf.NewTriangle(1, 3, 2)}}}

	components := &go3mf.Components{
		AnyAttr: go3mf.AnyAttr{
			&AssociationAttr{association: Association_physical},
			&OperationAttr{operation: BooleanOperation_union},
		}, Component: []*go3mf.Component{
			{ObjectID: 1, Transform: go3mf.Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}},
			{ObjectID: 2, Transform: go3mf.Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}},
		}}
	object1 := &go3mf.Object{
		ID:         3,
		Name:       "model with embossed label",
		Type:       go3mf.ObjectTypeModel,
		Components: components,
	}

	want := &go3mf.Model{
		Path: "/3D/3dmodel.model",
		Resources: go3mf.Resources{
			Objects: []*go3mf.Object{validMesh1, validMesh2, object1},
		},
		Build: go3mf.Build{
			Items: []*go3mf.Item{{ObjectID: 3}},
		}, Units: go3mf.UnitMillimeter,
		Language: "en-US",
	}
	want.Extensions = []go3mf.Extension{DefaultExtension}
	got := &go3mf.Model{
		Path: "/3D/3dmodel.model",
	}
	rootFile := `
	<?xml version="1.0" encoding="utf-8" standalone="no"?>
<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:bo="http://www.hp.com/schemas/3dmanufacturing/booleanoperations/2021/02" unit="millimeter" xml:lang="en-US">
	<resources>
		<object id="1" type="model" name="shuttle">
			<mesh>
				<vertices>
					<vertex x="45.00000" y="55.00000" z="55.00000"/>
					<vertex x="45.00000" y="45.00000" z="55.00000"/>
					<vertex x="45.00000" y="55.00000" z="45.00000"/>
					<vertex x="45.00000" y="45.00000" z="45.00000"/>
					<vertex x="55.00000" y="55.00000" z="45.00000"/>
					<vertex x="55.00000" y="55.00000" z="55.00000"/>
					<vertex x="55.00000" y="45.00000" z="55.00000"/>
					<vertex x="55.00000" y="45.00000" z="45.00000"/>
				</vertices>
				<triangles>
					<triangle v1="0" v2="1" v3="2"/>
					<triangle v1="0" v2="3" v3="1"/>
					<triangle v1="0" v2="2" v3="3"/>
					<triangle v1="1" v2="3" v3="2"/>
				</triangles>
			</mesh>
		</object>
		<object id="2" type="model" name="label 1">
			<mesh>
				<vertices>
					<vertex x="45.00000" y="55.00000" z="55.00000"/>
					<vertex x="45.00000" y="45.00000" z="55.00000"/>
					<vertex x="45.00000" y="55.00000" z="45.00000"/>
					<vertex x="45.00000" y="45.00000" z="45.00000"/>
					<vertex x="55.00000" y="55.00000" z="45.00000"/>
					<vertex x="55.00000" y="55.00000" z="55.00000"/>
					<vertex x="55.00000" y="45.00000" z="55.00000"/>
					<vertex x="55.00000" y="45.00000" z="45.00000"/>
				</vertices>
				<triangles>
					<triangle v1="0" v2="1" v3="2"/>
					<triangle v1="0" v2="3" v3="1"/>
					<triangle v1="0" v2="2" v3="3"/>
					<triangle v1="1" v2="3" v3="2"/>
				</triangles>
			</mesh>
		</object>
		<object id="3" type="model" name="model with embossed label">
			<components bo:association="physical" bo:operation="union">
				<component objectid="1" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
				<component objectid="2" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
			</components>
		</object>
	</resources>
	<build>
		<item objectid="3"/>
	</build>
</model>
		`

	t.Run("base", func(t *testing.T) {
		if err := go3mf.UnmarshalModel([]byte(rootFile), got); err != nil {
			t.Errorf("DecodeRawModel() unexpected error = %v", err)
			return
		}
		if diff := deep.Equal(got, want); diff != nil {
			t.Errorf("DecodeRawModel() = %v", diff)
			return
		}
	})
}
