package slices

import (
	"encoding/xml"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestMarshalModel(t *testing.T) {
	otherSlices := &SliceStackResource{
		ID: 10, ModelPath: "/3D/3dmodel.model",
		BottomZ: 2,
		Slices: []*Slice{
			{
				TopZ:     1.2,
				Vertices: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: []Polygon{{StartV: 0, Segments: []Segment{{V2: 1}, {V2: 2}, {V2: 3}, {V2: 0}}}},
			},
		},
	}
	sliceStack := &SliceStackResource{ID: 3, ModelPath: "/3D/3dmodel.model", BottomZ: 1,
		Slices: []*Slice{
			{
				TopZ:     0,
				Vertices: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: []Polygon{{StartV: 0, Segments: []Segment{{V2: 1}, {V2: 2}, {V2: 3}, {V2: 0}}}},
			},
			{
				TopZ:     0.1,
				Vertices: []go3mf.Point2D{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: []Polygon{{StartV: 0, Segments: []Segment{{V2: 2}, {V2: 1}, {V2: 3}, {V2: 0}}}},
			},
		},
	}
	sliceStackRef := &SliceStackResource{ID: 7, ModelPath: "/3D/3dmodel.model", BottomZ: 1.1, Refs: []SliceRef{{SliceStackID: 10, Path: "/3D/3dmodel.model"}}}
	meshRes := &go3mf.ObjectResource{
		Mesh: new(go3mf.Mesh),
		ID:   8, Name: "Box 1", ModelPath: "/3D/3dmodel.model",
		Extensions: go3mf.Extensions{ExtensionName: &SliceStackInfo{SliceStackID: 3, SliceResolution: ResolutionLow}},
	}

	m := &go3mf.Model{Path: "/3D/3dmodel.model", Namespaces: []xml.Name{{Space: ExtensionName, Local: "s"}}}
	m.Resources = append(m.Resources, otherSlices, sliceStack, sliceStackRef, meshRes)

	t.Run("base", func(t *testing.T) {
		b, err := go3mf.MarshalModel(m)
		if err != nil {
			t.Errorf("slices.MarshalModel() error = %v", err)
			return
		}
		d := go3mf.NewDecoder(nil, 0)
		RegisterExtension(d)
		newModel := new(go3mf.Model)
		newModel.Path = m.Path
		if err := d.UnmarshalModel(b, newModel); err != nil {
			t.Errorf("slices.MarshalModel() error decoding = %v, s = %s", err, string(b))
			return
		}
		if diff := deep.Equal(m, newModel); diff != nil {
			t.Errorf("slices.MarshalModel() = %v, s = %s", diff, string(b))
		}
	})
}
