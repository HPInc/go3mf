package booleanoperations

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestMarshalModel(t *testing.T) {
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
			&BooleanOperationAttr{association: Association_physical, operation: BooleanOperation_union},
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

	m := &go3mf.Model{
		Path:       "/3D/3dmodel.model",
		Extensions: []go3mf.Extension{DefaultExtension},
		Resources: go3mf.Resources{
			Objects: []*go3mf.Object{validMesh1, validMesh2, object1},
		},
		Build: go3mf.Build{
			Items: []*go3mf.Item{{ObjectID: 3}},
		}, Units: go3mf.UnitMillimeter,
		Language: "en-US",
	}
	m.Extensions = []go3mf.Extension{DefaultExtension}
	b, err := go3mf.MarshalModel(m)
	if err != nil {
		t.Errorf("booleanoperations.MarshalModel() error = %v", err)
		return
	}
	newModel := new(go3mf.Model)
	newModel.Path = m.Path
	if err := go3mf.UnmarshalModel(b, newModel); err != nil {
		t.Errorf("booleanoperations.MarshalModel() error decoding = %v, s = %s", err, string(b))
		return
	}

	if diff := deep.Equal(m, newModel); diff != nil {
		t.Errorf("booleanoperations.MarshalModel() = %v, s = %s", diff, string(b))
	}
}
