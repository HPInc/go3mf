package go3mf

import (
	"bytes"
	"encoding/xml"
	"image/color"
	"testing"

	"github.com/go-test/deep"
)

func TestEncoder_writeModel(t *testing.T) {
	want := `
	<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" unit="millimeter" lang="en-US">
		<resources>
			<basematerials id="5">
				<base name="Blue PLA" displaycolor="#0000FF" />
				<base name="Red ABS" displaycolor="#FF0000" />
			</basematerials>
			<object id="8" name="Box 1" pid="5" pindex="0" thumbnail="/a.png" partnumber="11111111-1111-1111-1111-111111111111" type="model">
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
						<triangle v1="4" v2="5" v3="6" p1="1" />
						<triangle v1="6" v2="7" v3="4" pid="5" p1="1" />
						<triangle v1="0" v2="1" v3="5" pid="5" p1="0" p2="1" p3="2"/>
						<triangle v1="5" v2="4" v3="0" pid="5" p1="3" p2="0" p3="2"/>
						<triangle v1="1" v2="2" v3="6" pid="5" p1="0" p2="1" p3="2"/>
						<triangle v1="6" v2="5" v3="1" pid="5" p1="2" p2="1" p3="3"/>
						<triangle v1="2" v2="3" v3="7" />
						<triangle v1="7" v2="6" v3="2" />
						<triangle v1="3" v2="0" v3="4" />
						<triangle v1="4" v2="7" v3="3" />
					</triangles>
				</mesh>
			</object>
			<object id="20">
				<metadatagroup>
					<metadata name="qm:CustomMetadata3" type="xs:boolean">1</metadata>
					<metadata name="qm:CustomMetadata4" type="xs:boolean">2</metadata>
				</metadatagroup>
				<components>
					<component objectid="8" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
				</components>
			</object>
		</resources>
		<build>
			<item partnumber="bob" objectid="20" transform="1 0 0 0 2 0 0 0 3 -66.4 -87.1 8.8">
				<metadatagroup>
					<metadata name="qm:CustomMetadata3" type="xs:boolean">1</metadata>
				</metadatagroup>
			</item>
		</build>
		<metadata name="Application">go3mf app</metadata>
		<metadata name="qm:CustomMetadata1" type="xs:string" preserve="1">CE8A91FB-C44E-4F00-B634-BAA411465F6A</metadata>
		<other />
	</model>`
	type args struct {
		w *bytes.Buffer
		m *Model
	}
	tests := []struct {
		name string
		e    *Encoder
		args args
		want string
	}{
		{"base", new(Encoder), args{new(bytes.Buffer), &Model{
			Units: UnitMillimeter, Language: "en-US", Path: "/3d/3dmodel.model", Thumbnail: "/thumbnail.png",
			Resources: []Resource{
				&BaseMaterialsResource{ID: 5, ModelPath: "/3d/3dmodel.model", Materials: []BaseMaterial{
					{Name: "Blue PLA", Color: color.RGBA{0, 0, 255, 255}},
					{Name: "Red ABS", Color: color.RGBA{255, 0, 0, 255}},
				}}, &ObjectResource{ID: 8, Name: "Box 1", PartNumber: "11111111-1111-1111-1111-111111111111", Thumbnail: "/a.png",
					DefaultPropertyID: 1, DefaultPropertyIndex: 1, ObjectType: ObjectTypeModel, Mesh: &Mesh{
						Nodes: []Point3D{
							{0, 0, 0}, {100, 0, 0}, {100, 100, 0},
							{0, 100, 0}, {0, 0, 100}, {100, 0, 100},
							{100, 100, 100}, {0, 100, 100}},
						Faces: []Face{
							{NodeIndices: [3]uint32{3, 2, 1}, Resource: 5},
							{NodeIndices: [3]uint32{1, 0, 3}, Resource: 5},
							{NodeIndices: [3]uint32{4, 5, 6}, Resource: 5, ResourceIndices: [3]uint32{1, 1, 1}},
							{NodeIndices: [3]uint32{6, 7, 4}, Resource: 5, ResourceIndices: [3]uint32{1, 1, 1}},
							{NodeIndices: [3]uint32{0, 1, 5}, Resource: 5, ResourceIndices: [3]uint32{0, 1, 2}},
							{NodeIndices: [3]uint32{5, 4, 0}, Resource: 5, ResourceIndices: [3]uint32{3, 0, 2}},
							{NodeIndices: [3]uint32{1, 2, 6}, Resource: 5, ResourceIndices: [3]uint32{0, 1, 2}},
							{NodeIndices: [3]uint32{6, 5, 1}, Resource: 5, ResourceIndices: [3]uint32{2, 1, 3}},
							{NodeIndices: [3]uint32{2, 3, 7}, Resource: 5},
							{NodeIndices: [3]uint32{7, 6, 2}, Resource: 5},
							{NodeIndices: [3]uint32{3, 0, 4}, Resource: 5},
							{NodeIndices: [3]uint32{4, 7, 3}, Resource: 5},
						},
					}},
				&ObjectResource{
					ID: 20, ModelPath: "/3d/3dmodel.model",
					Metadata:   []Metadata{{Name: "fake_ext:CustomMetadata3", Type: "xs:boolean", Value: "1"}, {Name: "fake_ext:CustomMetadata4", Type: "xs:boolean", Value: "2"}},
					Components: []*Component{{ObjectID: 8, Transform: Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}}},
				}}, Build: Build{Items: []*Item{{
				ObjectID: 20, PartNumber: "bob", Transform: Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
				Metadata: []Metadata{{Name: "fake_ext:CustomMetadata3", Type: "xs:boolean", Value: "1"}},
			}}}, Metadata: []Metadata{
				{Name: "Application", Value: "go3mf app"},
				{Name: "fake_ext:CustomMetadata1", Preserve: true, Type: "xs:string", Value: "CE8A91FB-C44E-4F00-B634-BAA411465F6A"},
			}}}, want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := xml.NewEncoder(tt.args.w)
			if err := tt.e.writeModel(x, tt.args.m); err != nil {
				t.Errorf("Encoder.writeModel() error = %v", err)
				return
			}
			x.Flush()
			if diff := deep.Equal(tt.args.w.String(), tt.want); diff != nil {
				t.Errorf("Encoder.writeModel() = %v", diff)
			}
		})
	}
}
