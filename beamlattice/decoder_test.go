package beamlattice

import (
	"encoding/xml"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
)

func TestDecode(t *testing.T) {
	meshLattice := &go3mf.ObjectResource{
		ID: 15, Name: "Box",
		Mesh: &go3mf.Mesh{
			Extension: go3mf.Extension{
				ExtensionName: &BeamLattice{ClipMode: ClipInside, ClippingMeshID: 8, RepresentationMeshID: 8},
			}},
	}
	beamLattice := MeshBeamLattice(meshLattice.Mesh)
	beamLattice.MinLength = 0.0001
	beamLattice.CapMode = CapModeHemisphere
	beamLattice.DefaultRadius = 1
	meshLattice.Mesh.Nodes = append(meshLattice.Mesh.Nodes, []go3mf.Point3D{
		{45, 55, 55},
		{45, 45, 55},
		{45, 55, 45},
		{45, 45, 45},
		{55, 55, 45},
		{55, 55, 55},
		{55, 45, 55},
		{55, 45, 45},
	}...)
	beamLattice.BeamSets = append(beamLattice.BeamSets, BeamSet{Name: "test", Identifier: "set_id", Refs: []uint32{1}})
	beamLattice.Beams = append(beamLattice.Beams, []Beam{
		{NodeIndices: [2]uint32{0, 1}, Radius: [2]float32{1.5, 1.6}, CapMode: [2]CapMode{CapModeSphere, CapModeButt}},
		{NodeIndices: [2]uint32{2, 0}, Radius: [2]float32{3, 1.5}, CapMode: [2]CapMode{CapModeSphere, CapModeHemisphere}},
		{NodeIndices: [2]uint32{1, 3}, Radius: [2]float32{1.6, 3}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{NodeIndices: [2]uint32{3, 2}, Radius: [2]float32{1, 1}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{NodeIndices: [2]uint32{2, 4}, Radius: [2]float32{3, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{NodeIndices: [2]uint32{4, 5}, Radius: [2]float32{2, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{NodeIndices: [2]uint32{5, 6}, Radius: [2]float32{2, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{NodeIndices: [2]uint32{7, 6}, Radius: [2]float32{2, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{NodeIndices: [2]uint32{1, 6}, Radius: [2]float32{1.6, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{NodeIndices: [2]uint32{7, 4}, Radius: [2]float32{2, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{NodeIndices: [2]uint32{7, 3}, Radius: [2]float32{2, 3}, CapMode: [2]CapMode{CapModeHemisphere, CapModeHemisphere}},
		{NodeIndices: [2]uint32{0, 5}, Radius: [2]float32{1.5, 2}, CapMode: [2]CapMode{CapModeHemisphere, CapModeButt}},
	}...)

	want := &go3mf.Model{Path: "/3D/3dmodel.model", Namespaces: []xml.Name{{Space: ExtensionName, Local: "b"}}, Resources: go3mf.Resources{
		Objects: []*go3mf.ObjectResource{meshLattice},
	}}
	got := new(go3mf.Model)
	got.Path = "/3D/3dmodel.model"
	rootFile := `
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:b="http://schemas.microsoft.com/3dmanufacturing/beamlattice/2017/02">
		<resources>
			<object id="15" name="Box" type="model">
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
					<b:other/>
					<b:beamlattice radius="1" minlength="0.0001" cap="hemisphere" clippingmode="inside" clippingmesh="8" representationmesh="8">
						<b:beams>
							<b:beam v1="0" v2="1" r1="1.50000" r2="1.60000" cap1="sphere" cap2="butt"/>
							<b:beam v1="2" v2="0" r1="3.00000" r2="1.50000" cap1="sphere"/>
							<b:beam v1="1" v2="3" r1="1.60000" r2="3.00000"/>
							<b:beam v1="3" v2="2" />
							<b:beam v1="2" v2="4" r1="3.00000" r2="2.00000"/>
							<b:beam v1="4" v2="5" r1="2.00000"/>
							<b:beam v1="5" v2="6" r1="2.00000"/>
							<b:beam v1="7" v2="6" r1="2.00000"/>
							<b:beam v1="1" v2="6" r1="1.60000" r2="2.00000"/>
							<b:beam v1="7" v2="4" r1="2.00000"/>
							<b:beam v1="7" v2="3" r1="2.00000" r2="3.00000"/>
							<b:beam v1="0" v2="5" r1="1.50000" r2="2.00000" cap2="butt"/>
						</b:beams>
						<b:beamsets>
							<b:beamset name="test" identifier="set_id">
								<b:ref index="1"/>
							</b:beamset>
						</b:beamsets>
					</b:beamlattice>
				</mesh>
			</object>
		</resources>
		<build>
		</build>
		</model>
		`

	t.Run("base", func(t *testing.T) {
		d := new(go3mf.Decoder)
		RegisterExtension(d)
		d.Strict = true
		if err := d.UnmarshalModel([]byte(rootFile), got); err != nil {
			t.Errorf("DecodeRawModel() unexpected error = %v", err)
			return
		}
		deep.CompareUnexportedFields = true
		deep.MaxDepth = 20
		if diff := deep.Equal(got, want); diff != nil {
			t.Errorf("DecodeRawModel() = %v", diff)
			return
		}
	})
}

func TestDecode_warns(t *testing.T) {
	want := []error{
		go3mf.MissingPropertyError{ResourceID: 15, Element: "beamlattice", ModelPath: "/3D/3dmodel.model", Name: "radius"},
		go3mf.MissingPropertyError{ResourceID: 15, Element: "beamlattice", ModelPath: "/3D/3dmodel.model", Name: "minlength"},
		go3mf.ParsePropertyError{ResourceID: 15, Element: "beamlattice", ModelPath: "/3D/3dmodel.model", Name: "cap", Value: "invalid", Type: go3mf.PropertyOptional},
		go3mf.ParsePropertyError{ResourceID: 15, Element: "beamlattice", ModelPath: "/3D/3dmodel.model", Name: "clippingmode", Value: "invalid2", Type: go3mf.PropertyOptional},
		go3mf.MissingPropertyError{ResourceID: 15, Element: "beam", ModelPath: "/3D/3dmodel.model", Name: "v1"},
		go3mf.MissingPropertyError{ResourceID: 15, Element: "beam", ModelPath: "/3D/3dmodel.model", Name: "v2"},
		go3mf.MissingPropertyError{ResourceID: 15, Element: "ref", ModelPath: "/3D/3dmodel.model", Name: "index"},
		go3mf.ParsePropertyError{ResourceID: 15, Element: "ref", Name: "index", Value: "a", ModelPath: "/3D/3dmodel.model", Type: go3mf.PropertyRequired},
	}
	got := new(go3mf.Model)
	got.Path = "/3D/3dmodel.model"
	rootFile := `
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:b="http://schemas.microsoft.com/3dmanufacturing/beamlattice/2017/02">
		<resources>
			<object id="15" name="Box" type="model">
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
					<b:beamlattice />
					<b:beamlattice qm:mq="other" radius="1" minlength="0.0001" cap="invalid" clippingmode="invalid2" clippingmesh="8" representationmesh="8">
						<b:beams>
							<b:beam qm:mq="other" v1="0" v2="1" r1="1.50000" r2="1.60000" cap1="sphere" cap2="butt"/>
							<b:beam v1="2" v2="0" r1="3.00000" r2="1.50000" cap1="sphere"/>
							<b:beam v1="1" v2="3" r1="1.60000" r2="3.00000"/>
							<b:beam v1="3" v2="2" />
							<b:beam />
							<b:beam v1="2" v2="4" r1="3.00000" r2="2.00000"/>
							<b:beam v1="4" v2="5" r1="2.00000"/>
							<b:beam v1="5" v2="6" r1="2.00000"/>
							<b:beam v1="7" v2="6" r1="2.00000"/>
							<b:beam v1="1" v2="6" r1="1.60000" r2="2.00000"/>
							<b:beam v1="7" v2="4" r1="2.00000"/>
							<b:beam v1="7" v2="3" r1="2.00000" r2="3.00000"/>
							<b:beam v1="0" v2="5" r1="1.50000" r2="2.00000" cap2="butt"/>
						</b:beams>
						<b:beamsets>
							<b:beamset qm:mq="other" name="test" identifier="set_id">
								<b:ref index="1"/>
								<b:ref />
								<b:ref index="a"/>
							</b:beamset>
						</b:beamsets>
					</b:beamlattice>
				</mesh>
			</object>
		</resources>
		<build>
		</build>
		</model>
		`

	t.Run("base", func(t *testing.T) {
		d := new(go3mf.Decoder)
		RegisterExtension(d)
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
