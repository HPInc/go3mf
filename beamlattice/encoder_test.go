// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package beamlattice

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/hpinc/go3mf"
	"github.com/hpinc/go3mf/spec"
)

func TestMarshalModel(t *testing.T) {
	beamLattice := &BeamLattice{ClipMode: ClipInside, ClippingMeshID: 8, RepresentationMeshID: 8}
	meshLattice := &go3mf.Object{
		ID: 15, Name: "Box",
		Mesh: &go3mf.Mesh{
			Triangles: []go3mf.Triangle{},
			Any:       spec.Any{beamLattice}},
	}
	beamLattice.MinLength = 0.0001
	beamLattice.CapMode = CapModeHemisphere
	beamLattice.Radius = 1
	meshLattice.Mesh.Vertices = append(meshLattice.Mesh.Vertices, []go3mf.Point3D{
		{45, 55, 55},
		{45, 45, 55},
		{45, 55, 45},
		{45, 45, 45},
		{55, 55, 45},
		{55, 55, 55},
		{55, 45, 55},
		{55, 45, 45},
	}...)
	beamLattice.BeamSets.BeamSet = append(beamLattice.BeamSets.BeamSet, BeamSet{Name: "test", Identifier: "set_id", Refs: []uint32{1}})
	beamLattice.Beams.Beam = append(beamLattice.Beams.Beam, []Beam{
		{Indices: [2]uint32{0, 1}, Radius: [2]float32{1.5, 1.6}, CapMode: [2]CapMode{CapModeSphere, CapModeButt}},
		{Indices: [2]uint32{2, 0}, Radius: [2]float32{3, 1.5}, CapMode: [2]CapMode{CapModeSphere, CapModeHemisphere}},
		{Indices: [2]uint32{1, 3}, Radius: [2]float32{1.6, 3}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{Indices: [2]uint32{3, 2}, Radius: [2]float32{1, 1}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{Indices: [2]uint32{2, 4}, Radius: [2]float32{3, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{Indices: [2]uint32{4, 5}, Radius: [2]float32{2, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{Indices: [2]uint32{5, 6}, Radius: [2]float32{2, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{Indices: [2]uint32{7, 6}, Radius: [2]float32{2, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{Indices: [2]uint32{1, 6}, Radius: [2]float32{1.6, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{Indices: [2]uint32{7, 4}, Radius: [2]float32{2, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{Indices: [2]uint32{7, 3}, Radius: [2]float32{2, 3}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{Indices: [2]uint32{0, 5}, Radius: [2]float32{1.5, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeButt}},
	}...)

	m := &go3mf.Model{
		Path:       "/3D/3dmodel.model",
		Extensions: []go3mf.Extension{DefaultExtension},
		Resources: go3mf.Resources{
			Objects: []*go3mf.Object{meshLattice},
		},
	}

	t.Run("base", func(t *testing.T) {
		b, err := go3mf.MarshalModel(m)
		if err != nil {
			t.Errorf("beamlattice.MarshalModel() error = %v", err)
			return
		}
		newModel := new(go3mf.Model)
		newModel.Path = m.Path
		if err := go3mf.UnmarshalModel(b, newModel); err != nil {
			t.Errorf("beamlattice.MarshalModel() error decoding = %v, s = %s", err, string(b))
			return
		}
		if diff := deep.Equal(m, newModel); diff != nil {
			t.Errorf("beamlattice.MarshalModel() = %v, s = %s", diff, string(b))
		}
	})
}
