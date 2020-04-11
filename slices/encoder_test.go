package slices

import (
	"image/color"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestMarshalModel(t *testing.T) {
	sliceStack := &SliceStack{ID: 3, BottomZ: 1,
		Slices: []*Slice{
			{
				TopZ:     0,
				Vertices: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: []Polygon{{StartV: 0, Segments: []Segment{{V2: 1, PID: 10}, {V2: 2, PID: 10, P2: 1}, {V2: 3}, {V2: 0}}}},
			},
			{
				TopZ:     0.1,
				Vertices: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: []Polygon{{StartV: 0, Segments: []Segment{{V2: 2}, {V2: 1}, {V2: 3}, {V2: 0}}}},
			},
		},
	}
	sliceStackRef := &SliceStack{ID: 7, BottomZ: 1.1, Refs: []SliceRef{{SliceStackID: 10, Path: "/2D/2dmodel.model"}}}
	meshRes := &go3mf.Object{
		Mesh: new(go3mf.Mesh),
		ID:   8, Name: "Box 1",
		AnyAttr: go3mf.AttrMarshalers{&SliceStackInfo{SliceStackID: 3, MeshResolution: ResolutionLow}},
	}
	baseMaterial := &go3mf.BaseMaterials{ID: 10, Materials: []go3mf.Base{{Name: "a", Color: color.RGBA{R: 1}}, {Name: "b", Color: color.RGBA{R: 1}}}}

	m := &go3mf.Model{Path: "/3D/3dmodel.model", Resources: go3mf.Resources{
		Assets: []go3mf.Asset{baseMaterial, sliceStack, sliceStackRef}, Objects: []*go3mf.Object{meshRes},
	}}

	t.Run("base", func(t *testing.T) {
		m.WithSpec(&Spec{LocalName: "s"})
		b, err := go3mf.MarshalModel(m)
		if err != nil {
			t.Errorf("slices.MarshalModel() error = %v", err)
			return
		}
		newModel := new(go3mf.Model)
		newModel.Path = m.Path
		newModel.WithSpec(&Spec{LocalName: "s"})
		if err := go3mf.UnmarshalModel(b, newModel); err != nil {
			t.Errorf("slices.MarshalModel() error decoding = %v, s = %s", err, string(b))
			return
		}
		if diff := deep.Equal(m, newModel); diff != nil {
			t.Errorf("slices.MarshalModel() = %v, s = %s", diff, string(b))
		}
	})
}
