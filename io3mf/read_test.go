package io3mf

import (
	"bytes"
	"compress/flate"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/geo"
	"github.com/stretchr/testify/mock"
)

type mockRelationship struct {
	mock.Mock
}

func newMockRelationship(relType, targetURI string) *mockRelationship {
	m := new(mockRelationship)
	m.On("Type").Return(relType).Maybe()
	m.On("TargetURI").Return(targetURI).Maybe()
	return m
}

func (m *mockRelationship) Type() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockRelationship) TargetURI() string {
	args := m.Called()
	return args.String(0)
}

type modelBuilder struct {
	str      strings.Builder
	hasModel bool
}

func (m *modelBuilder) withElement(s string) *modelBuilder {
	m.str.WriteString(s)
	m.str.WriteString("\n")
	return m
}

func (m *modelBuilder) addAttr(prefix, name, value string) *modelBuilder {
	if prefix != "" {
		m.str.WriteString(fmt.Sprintf(`%s:`, prefix))
	}
	if name != "" {
		m.str.WriteString(fmt.Sprintf(`%s="%s" `, name, value))
	}
	return m
}

func (m *modelBuilder) withDefaultModel() *modelBuilder {
	m.withModel("millimeter", "en-US")
	return m
}

func (m *modelBuilder) withModel(unit string, lang string) *modelBuilder {
	m.str.WriteString(`<model `)
	m.addAttr("", "unit", unit).addAttr("xml", "lang", lang)
	m.addAttr("", "xmlns", nsCoreSpec).addAttr("xmlns", "m", nsMaterialSpec).addAttr("xmlns", "p", nsProductionSpec)
	m.addAttr("xmlns", "b", nsBeamLatticeSpec).addAttr("", "requiredextensions", "m p b")
	m.str.WriteString(">\n")
	m.hasModel = true
	return m
}

func (m *modelBuilder) withEncoding(encode string) *modelBuilder {
	m.str.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="%s"?>`, encode))
	m.str.WriteString("\n")
	return m
}

func (m *modelBuilder) build() *mockFile {
	if m.hasModel {
		m.str.WriteString("</model>\n")
	}
	f := new(mockFile)
	f.On("Name").Return("/3d/3dmodel.model").Maybe()
	f.On("Open").Return(ioutil.NopCloser(bytes.NewBufferString(m.str.String())), nil).Maybe()
	return f
}

type mockFile struct {
	mock.Mock
}

func newMockFile(name string, relationships []relationship, thumb *mockFile, other *mockFile, openErr bool) *mockFile {
	m := new(mockFile)
	m.On("Name").Return(name).Maybe()
	m.On("Relationships").Return(relationships).Maybe()
	m.On("FindFileFromRel", relTypeThumbnail).Return(thumb, thumb != nil).Maybe()
	m.On("FindFileFromRel", mock.Anything).Return(other, other != nil).Maybe()
	var err error
	if openErr {
		err = errors.New("")
	}
	m.On("Open").Return(ioutil.NopCloser(new(bytes.Buffer)), err).Maybe()
	return m
}

func (m *mockFile) Open() (io.ReadCloser, error) {
	args := m.Called()
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockFile) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockFile) FindFileFromRel(args0 string) (packageFile, bool) {
	args := m.Called(args0)
	return args.Get(0).(packageFile), args.Bool(1)
}

func (m *mockFile) Relationships() []relationship {
	args := m.Called()
	return args.Get(0).([]relationship)
}

type mockPackage struct {
	mock.Mock
}

func newMockPackage(other *mockFile) *mockPackage {
	m := new(mockPackage)
	m.On("Open", mock.Anything).Return(nil).Maybe()
	m.On("FindFileFromRel", mock.Anything).Return(other, other != nil).Maybe()
	m.On("FindFileFromName", mock.Anything).Return(other, other != nil).Maybe()
	return m
}

func (m *mockPackage) Open(f func(r io.Reader) io.ReadCloser) error {
	args := m.Called(f)
	return args.Error(0)
}

func (m *mockPackage) FindFileFromName(args0 string) (packageFile, bool) {
	args := m.Called(args0)
	return args.Get(0).(packageFile), args.Bool(1)
}

func (m *mockPackage) FindFileFromRel(args0 string) (packageFile, bool) {
	args := m.Called(args0)
	return args.Get(0).(packageFile), args.Bool(1)
}

func TestReadError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *ReadError
		want string
	}{
		{"new", new(ReadError), ""},
		{"generic", &ReadError{Message: "generic error"}, "generic error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("ReadError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_processOPC(t *testing.T) {
	thumbFile := newMockFile("/a.png", nil, nil, nil, false)
	thumbErr := newMockFile("/a.png", nil, nil, nil, true)
	tests := []struct {
		name    string
		d       *Decoder
		want    *go3mf.Model
		wantErr bool
	}{
		{"noRoot", &Decoder{p: newMockPackage(nil)}, &go3mf.Model{}, true},
		{"noRels", &Decoder{p: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}, &go3mf.Model{Path: "/a.model"}, false},
		{"withThumb", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbFile, thumbFile, false)),
		}, &go3mf.Model{
			Path:        "/a.model",
			Thumbnail:   &go3mf.Attachment{RelationshipType: relTypeThumbnail, Path: "/Metadata/thumbnail.png", Stream: new(bytes.Buffer)},
			Attachments: []*go3mf.Attachment{{RelationshipType: relTypeThumbnail, Path: "/a.png", Stream: new(bytes.Buffer)}},
		}, false},
		{"withThumbErr", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbErr, thumbErr, false)),
		}, &go3mf.Model{Path: "/a.model"}, false},
		{"withOtherRel", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("other", "/a.png")}, nil, nil, false)),
		}, &go3mf.Model{Path: "/a.model"}, false},
		{"withModelAttachment", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeModel3D, "/a.model")}, nil, newMockFile("/a.model", nil, nil, nil, false), false)),
		}, &go3mf.Model{
			Path:                  "/a.model",
			ProductionAttachments: []*go3mf.ProductionAttachment{{RelationshipType: relTypeModel3D, Path: "/a.model"}},
		}, false},
		{"withAttRel", &Decoder{AttachmentRelations: []string{"b"},
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("b", "/a.xml")}, nil, newMockFile("/a.xml", nil, nil, nil, false), false)),
		}, &go3mf.Model{
			Path:        "/a.model",
			Attachments: []*go3mf.Attachment{{RelationshipType: "b", Path: "/a.xml", Stream: new(bytes.Buffer)}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := new(go3mf.Model)
			_, err := tt.d.processOPC(model)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.processOPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(model, tt.want); diff != nil {
				t.Errorf("Decoder.processOPC() = %v", diff)
				return
			}
		})
	}
}

func TestDecoder_processRootModel_Fail(t *testing.T) {
	tests := []struct {
		name    string
		f       *mockFile
		wantErr bool
	}{
		{"errOpen", newMockFile("/a.model", nil, nil, nil, true), true},
		{"errEncode", new(modelBuilder).withEncoding("utf16").build(), true},
		{"invalidUnits", new(modelBuilder).withModel("other", "en-US").build(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := new(Decoder).processRootModel(context.Background(), tt.f, new(go3mf.Model)); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.processRootModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDecoder_processRootModel(t *testing.T) {
	baseMaterials := &go3mf.BaseMaterialsResource{ID: 5, ModelPath: "/3d/3dmodel.model", Materials: []go3mf.BaseMaterial{
		{Name: "Blue PLA", Color: color.RGBA{0, 0, 255, 255}},
		{Name: "Red ABS", Color: color.RGBA{255, 0, 0, 255}},
	}}
	baseTexture := &go3mf.Texture2DResource{ID: 6, ModelPath: "/3d/3dmodel.model", Path: "/3D/Texture/msLogo.png", ContentType: go3mf.TextureTypePNG, TileStyleU: go3mf.TileWrap, TileStyleV: go3mf.TileMirror, Filter: go3mf.TextureFilterAuto}
	meshRes := &go3mf.MeshResource{
		ObjectResource: go3mf.ObjectResource{
			ID: 8, Name: "Box 1", ModelPath: "/3d/3dmodel.model", Thumbnail: "/a.png", DefaultPropertyID: 5, PartNumber: "11111111-1111-1111-1111-111111111111",
			Attr: map[string]interface{}{}},
		Mesh: new(geo.Mesh),
	}
	meshRes.Mesh.Nodes = append(meshRes.Mesh.Nodes, []geo.Point3D{
		{0, 0, 0},
		{100, 0, 0},
		{100, 100, 0},
		{0, 100, 0},
		{0, 0, 100},
		{100, 0, 100},
		{100, 100, 100},
		{0, 100, 100},
	}...)
	meshRes.Mesh.Faces = append(meshRes.Mesh.Faces, []geo.Face{
		{NodeIndices: [3]uint32{3, 2, 1}, Resource: 5},
		{NodeIndices: [3]uint32{1, 0, 3}, Resource: 5},
		{NodeIndices: [3]uint32{4, 5, 6}, Resource: 5, ResourceIndices: [3]uint32{1, 1, 1}},
		{NodeIndices: [3]uint32{6, 7, 4}, Resource: 5, ResourceIndices: [3]uint32{1, 1, 1}},
		{NodeIndices: [3]uint32{0, 1, 5}, Resource: 2, ResourceIndices: [3]uint32{0, 1, 2}},
		{NodeIndices: [3]uint32{5, 4, 0}, Resource: 2, ResourceIndices: [3]uint32{3, 0, 2}},
		{NodeIndices: [3]uint32{1, 2, 6}, Resource: 1, ResourceIndices: [3]uint32{0, 1, 2}},
		{NodeIndices: [3]uint32{6, 5, 1}, Resource: 1, ResourceIndices: [3]uint32{2, 1, 3}},
		{NodeIndices: [3]uint32{2, 3, 7}, Resource: 5},
		{NodeIndices: [3]uint32{7, 6, 2}, Resource: 5},
		{NodeIndices: [3]uint32{3, 0, 4}, Resource: 5},
		{NodeIndices: [3]uint32{4, 7, 3}, Resource: 5},
	}...)

	meshLattice := &go3mf.MeshResource{
		ObjectResource:        go3mf.ObjectResource{ID: 15, Name: "Box", ModelPath: "/3d/3dmodel.model", PartNumber: "e1ef01d4-cbd4-4a62-86b6-9634e2ca198b", Attr: make(map[string]interface{})},
		BeamLatticeAttributes: go3mf.BeamLatticeAttributes{ClipMode: go3mf.ClipInside, ClippingMeshID: 8, RepresentationMeshID: 8},
		Mesh:                  new(geo.Mesh),
	}
	meshLattice.Mesh.MinLength = 0.0001
	meshLattice.Mesh.CapMode = geo.CapModeHemisphere
	meshLattice.Mesh.DefaultRadius = 1
	meshLattice.Mesh.Nodes = append(meshLattice.Mesh.Nodes, []geo.Point3D{
		{45, 55, 55},
		{45, 45, 55},
		{45, 55, 45},
		{45, 45, 45},
		{55, 55, 45},
		{55, 55, 55},
		{55, 45, 55},
		{55, 45, 45},
	}...)
	meshLattice.Mesh.BeamSets = append(meshLattice.Mesh.BeamSets, geo.BeamSet{Name: "test", Identifier: "set_id", Refs: []uint32{1}})
	meshLattice.Mesh.Beams = append(meshLattice.Mesh.Beams, []geo.Beam{
		{NodeIndices: [2]uint32{0, 1}, Radius: [2]float64{1.5, 1.6}, CapMode: [2]geo.CapMode{geo.CapModeSphere, geo.CapModeButt}},
		{NodeIndices: [2]uint32{2, 0}, Radius: [2]float64{3, 1.5}, CapMode: [2]geo.CapMode{geo.CapModeSphere, geo.CapModeHemisphere}},
		{NodeIndices: [2]uint32{1, 3}, Radius: [2]float64{1.6, 3}, CapMode: [2]geo.CapMode{geo.CapModeHemisphere, geo.CapModeHemisphere}},
		{NodeIndices: [2]uint32{3, 2}, Radius: [2]float64{1, 1}, CapMode: [2]geo.CapMode{geo.CapModeHemisphere, geo.CapModeHemisphere}},
		{NodeIndices: [2]uint32{2, 4}, Radius: [2]float64{3, 2}, CapMode: [2]geo.CapMode{geo.CapModeHemisphere, geo.CapModeHemisphere}},
		{NodeIndices: [2]uint32{4, 5}, Radius: [2]float64{2, 2}, CapMode: [2]geo.CapMode{geo.CapModeHemisphere, geo.CapModeHemisphere}},
		{NodeIndices: [2]uint32{5, 6}, Radius: [2]float64{2, 2}, CapMode: [2]geo.CapMode{geo.CapModeHemisphere, geo.CapModeHemisphere}},
		{NodeIndices: [2]uint32{7, 6}, Radius: [2]float64{2, 2}, CapMode: [2]geo.CapMode{geo.CapModeHemisphere, geo.CapModeHemisphere}},
		{NodeIndices: [2]uint32{1, 6}, Radius: [2]float64{1.6, 2}, CapMode: [2]geo.CapMode{geo.CapModeHemisphere, geo.CapModeHemisphere}},
		{NodeIndices: [2]uint32{7, 4}, Radius: [2]float64{2, 2}, CapMode: [2]geo.CapMode{geo.CapModeHemisphere, geo.CapModeHemisphere}},
		{NodeIndices: [2]uint32{7, 3}, Radius: [2]float64{2, 3}, CapMode: [2]geo.CapMode{geo.CapModeHemisphere, geo.CapModeHemisphere}},
		{NodeIndices: [2]uint32{0, 5}, Radius: [2]float64{1.5, 2}, CapMode: [2]geo.CapMode{geo.CapModeHemisphere, geo.CapModeButt}},
	}...)

	components := &go3mf.ComponentsResource{
		ObjectResource: go3mf.ObjectResource{
			ID: 20, UUID: "cb828680-8895-4e08-a1fc-be63e033df15", ModelPath: "/3d/3dmodel.model", Attr: make(map[string]interface{}),
			Metadata: []go3mf.Metadata{{Name: nsProductionSpec + ":CustomMetadata3", Type: "xs:boolean", Value: "1"}, {Name: nsProductionSpec + ":CustomMetadata4", Type: "xs:boolean", Value: "2"}},
		},
		Components: []*go3mf.Component{{UUID: "cb828680-8895-4e08-a1fc-be63e033df16", Object: meshRes,
			Transform: geo.Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}}},
	}

	want := &go3mf.Model{Units: go3mf.UnitMillimeter, Language: "en-US", Path: "/3d/3dmodel.model", UUID: "e9e25302-6428-402e-8633-cc95528d0ed3"}
	otherMesh := &go3mf.MeshResource{ObjectResource: go3mf.ObjectResource{ID: 8, ModelPath: "/3d/other.model", Attr: make(map[string]interface{})}, Mesh: new(geo.Mesh)}
	colorGroup := &go3mf.ColorGroupResource{ID: 1, ModelPath: "/3d/3dmodel.model", Colors: []color.RGBA{{R: 255, G: 255, B: 255, A: 255}, {R: 0, G: 0, B: 0, A: 255}, {R: 26, G: 181, B: 103, A: 255}, {R: 223, G: 4, B: 90, A: 255}}}
	texGroup := &go3mf.Texture2DGroupResource{ID: 2, ModelPath: "/3d/3dmodel.model", TextureID: 6, Coords: []go3mf.TextureCoord{{0.3, 0.5}, {0.3, 0.8}, {0.5, 0.8}, {0.5, 0.5}}}
	compositeGroup := &go3mf.CompositeMaterialsResource{ID: 4, ModelPath: "/3d/3dmodel.model", MaterialID: 5, Indices: []uint32{1, 2}, Composites: []go3mf.Composite{{Values: []float64{0.5, 0.5}}, {Values: []float64{0.2, 0.8}}}}
	multiGroup := &go3mf.MultiPropertiesResource{ID: 9, ModelPath: "/3d/3dmodel.model", BlendMethods: []go3mf.BlendMethod{go3mf.BlendMultiply}, Resources: []uint32{5, 2}, Multis: []go3mf.Multi{{ResourceIndices: []uint32{0, 0}}, {ResourceIndices: []uint32{1, 0}}, {ResourceIndices: []uint32{2, 3}}}}
	want.Resources = append(want.Resources, otherMesh, baseMaterials, baseTexture, colorGroup, texGroup, compositeGroup, multiGroup, meshRes, meshLattice, components)
	want.BuildItems = append(want.BuildItems, &go3mf.BuildItem{Object: components, PartNumber: "bob", UUID: "e9e25302-6428-402e-8633-cc95528d0ed2",
		Transform: geo.Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
	})
	want.BuildItems = append(want.BuildItems, &go3mf.BuildItem{Object: otherMesh, UUID: "e9e25302-6428-402e-8633-cc95528d0ed3", Metadata: []go3mf.Metadata{{Name: nsProductionSpec + ":CustomMetadata3", Type: "xs:boolean", Value: "1"}}})
	want.Metadata = append(want.Metadata, []go3mf.Metadata{
		{Name: "Application", Value: "go3mf app"},
		{Name: nsProductionSpec + ":CustomMetadata1", Preserve: true, Type: "xs:string", Value: "CE8A91FB-C44E-4F00-B634-BAA411465F6A"},
	}...)
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
	got.Resources = append(got.Resources, otherMesh)
	rootFile := new(modelBuilder).withDefaultModel().withElement(`
		<resources>
			<basematerials id="5">
				<base name="Blue PLA" displaycolor="#0000FF" />
				<base name="Red ABS" displaycolor="#FF0000" />
			</basematerials>
			<m:texture2d id="6" path="/3D/Texture/msLogo.png" contenttype="image/png" tilestyleu="wrap" tilestylev="mirror" filter="auto" />
			<m:colorgroup id="1">
				<m:color color="#FFFFFF" /> <m:color color="#000000" /> <m:color color="#1AB567" /> <m:color color="#DF045A" />
			</m:colorgroup>
			<m:texture2dgroup id="2" texid="6">
				<m:tex2coord u="0.3" v="0.5" /> <m:tex2coord u="0.3" v="0.8" />	<m:tex2coord u="0.5" v="0.8" />	<m:tex2coord u="0.5" v="0.5" />
			</m:texture2dgroup>
			<m:compositematerials id="4" matid="5" matindices="1 2">
				<m:composite values="0.5 0.5"/>
				<m:composite values="0.2 0.8"/>
			</m:compositematerials>
			<m:multiproperties id="9" pids="5 2" blendmethods="multiply">
				<m:multi pindices="0 0" />
				<m:multi pindices="1 0" />
				<m:multi pindices="2 3" />
			</m:multiproperties>
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
						<triangle v1="0" v2="1" v3="5" pid="2" p1="0" p2="1" p3="2"/>
						<triangle v1="5" v2="4" v3="0" pid="2" p1="3" p2="0" p3="2"/>
						<triangle v1="1" v2="2" v3="6" pid="1" p1="0" p2="1" p3="2"/>
						<triangle v1="6" v2="5" v3="1" pid="1" p1="2" p2="1" p3="3"/>
						<triangle v1="2" v2="3" v3="7" />
						<triangle v1="7" v2="6" v3="2" />
						<triangle v1="3" v2="0" v3="4" />
						<triangle v1="4" v2="7" v3="3" />
					</triangles>
				</mesh>
			</object>
			<object id="15" name="Box" partnumber="e1ef01d4-cbd4-4a62-86b6-9634e2ca198b" type="model">
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
			<object id="20" p:UUID="cb828680-8895-4e08-a1fc-be63e033df15">
				<metadatagroup>
					<metadata name="p:CustomMetadata3" type="xs:boolean">1</metadata>
					<metadata name="p:CustomMetadata4" type="xs:boolean">2</metadata>
				</metadatagroup>
				<components>
					<component objectid="8" p:UUID="cb828680-8895-4e08-a1fc-be63e033df16" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
				</components>
			</object>
		</resources>
		<build p:UUID="e9e25302-6428-402e-8633-cc95528d0ed3">
			<item partnumber="bob" objectid="20" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed2" transform="1 0 0 0 2 0 0 0 3 -66.4 -87.1 8.8" />
			<item objectid="8" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed3" p:path="/3d/other.model">
				<metadatagroup>
					<metadata name="p:CustomMetadata3" type="xs:boolean">1</metadata>
				</metadatagroup>
			</item>
		</build>
		<metadata name="Application">go3mf app</metadata>
		<metadata name="p:CustomMetadata1" type="xs:string" preserve="1">CE8A91FB-C44E-4F00-B634-BAA411465F6A</metadata>
		<other />
		`).build()

	t.Run("base", func(t *testing.T) {
		d := new(Decoder)
		d.Strict = true
		d.SetDecompressor(func(r io.Reader) io.ReadCloser { return flate.NewReader(r) })
		d.SetXMLDecoder(func(r io.Reader) XMLDecoder { return xml.NewDecoder(r) })
		if err := d.processRootModel(context.Background(), rootFile, got); err != nil {
			t.Errorf("Decoder.processRootModel() unexpected error = %v", err)
			return
		}
		deep.CompareUnexportedFields = true
		deep.MaxDepth = 20
		if diff := deep.Equal(got, want); diff != nil {
			t.Errorf("Decoder.processRootModel() = %v", diff)
			return
		}
	})
}

func TestDecoder_processNonRootModels(t *testing.T) {
	tests := []struct {
		name    string
		model   *go3mf.Model
		d       *Decoder
		wantErr bool
		want    *go3mf.Model
	}{
		{"base", &go3mf.Model{ProductionAttachments: []*go3mf.ProductionAttachment{
			{Path: "3d/new.model"},
			{Path: "3d/other.model"},
		}}, &Decoder{productionModels: map[string]packageFile{
			"3d/new.model": new(modelBuilder).withDefaultModel().withElement(`
				<resources>
					<basematerials id="5">
						<base name="Blue PLA" displaycolor="#0000FF" />
						<base name="Red ABS" displaycolor="#FF0000" />
					</basematerials>
				</resources>
			`).build(),
			"3d/other.model": new(modelBuilder).withDefaultModel().withElement(`
				<resources>
					<m:texture2d id="6" path="/3D/Texture/msLogo.png" contenttype="image/png" tilestyleu="wrap" tilestylev="mirror" filter="auto" />
				</resources>
			`).build(),
		}}, false, &go3mf.Model{
			ProductionAttachments: []*go3mf.ProductionAttachment{
				{Path: "3d/new.model"},
				{Path: "3d/other.model"},
			}, Resources: []go3mf.Resource{
				&go3mf.BaseMaterialsResource{ID: 5, ModelPath: "3d/new.model", Materials: []go3mf.BaseMaterial{
					{Name: "Blue PLA", Color: color.RGBA{0, 0, 255, 255}},
					{Name: "Red ABS", Color: color.RGBA{255, 0, 0, 255}},
				}},
				&go3mf.Texture2DResource{ID: 6, ModelPath: "3d/other.model", Path: "/3D/Texture/msLogo.png", ContentType: go3mf.TextureTypePNG, TileStyleU: go3mf.TileWrap, TileStyleV: go3mf.TileMirror, Filter: go3mf.TextureFilterAuto},
			},
		}},
		{"noAtt", new(go3mf.Model), new(Decoder), false, new(go3mf.Model)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.processNonRootModels(context.Background(), tt.model); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.processNonRootModels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			deep.CompareUnexportedFields = true
			deep.MaxDepth = 20
			if diff := deep.Equal(tt.model, tt.want); diff != nil {
				t.Errorf("Decoder.processNonRootModels() = %v", diff)
				return
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	tests := []struct {
		name    string
		d       *Decoder
		wantErr bool
	}{
		{"base", &Decoder{AttachmentRelations: []string{"b"},
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("b", "/a.xml")}, nil, newMockFile("/a.xml", nil, nil, nil, false), false)),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Decode(new(go3mf.Model)); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_modelFile_Decode(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	checkEveryBytes = 108
	type args struct {
		ctx context.Context
		x   *xml.Decoder
	}
	tests := []struct {
		name    string
		d       *modelFileDecoder
		args    args
		wantErr bool
	}{
		{"nochild", new(modelFileDecoder), args{context.Background(), xml.NewDecoder(bytes.NewBufferString(`
			<a></a>
			<b></b>
		`))}, false},
		{"eof", new(modelFileDecoder), args{context.Background(), xml.NewDecoder(bytes.NewBufferString(`
			<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02">
				<build></build>
		`))}, true},
		{"canceled", new(modelFileDecoder), args{ctx, xml.NewDecoder(bytes.NewBufferString(`
			<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02">
				<build></build>
			</model>
		`))}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Decode(tt.args.ctx, tt.args.x, new(go3mf.Model), "", true, false); (err != nil) != tt.wantErr {
				t.Errorf("modelFile.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDecoder(t *testing.T) {
	type args struct {
		r    io.ReaderAt
		size int64
	}
	tests := []struct {
		name string
		args args
		want *Decoder
	}{
		{"base", args{nil, 5}, &Decoder{Strict: true, p: &opcReader{ra: nil, size: 5}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDecoder(tt.args.r, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDecoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_processRootModel_warns(t *testing.T) {
	want := []error{
		go3mf.ParsePropertyError{ResourceID: 0, Element: "base", Name: "displaycolor", Value: "0000FF", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.MissingPropertyError{ResourceID: 0, Element: "base", ModelPath: "/3d/3dmodel.model", Name: "name"},
		go3mf.MissingPropertyError{ResourceID: 0, Element: "base", ModelPath: "/3d/3dmodel.model", Name: "displaycolor"},
		go3mf.MissingPropertyError{ResourceID: 0, Element: "basematerials", ModelPath: "/3d/3dmodel.model", Name: "id"},
		go3mf.ParsePropertyError{ResourceID: 0, Element: "basematerials", Name: "id", Value: "a", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.MissingPropertyError{ResourceID: 0, Element: "basematerials", ModelPath: "/3d/3dmodel.model", Name: "id"},
		go3mf.ParsePropertyError{ResourceID: 0, Element: "texture2d", Name: "id", Value: "b", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.MissingPropertyError{ResourceID: 0, Element: "texture2d", ModelPath: "/3d/3dmodel.model", Name: "path"},
		go3mf.MissingPropertyError{ResourceID: 0, Element: "texture2d", ModelPath: "/3d/3dmodel.model", Name: "id"},
		go3mf.ParsePropertyError{ResourceID: 1, Element: "color", Name: "color", Value: "#FFFFF", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 2, Element: "tex2coord", Name: "u", Value: "b", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 2, Element: "tex2coord", Name: "v", Value: "c", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.MissingPropertyError{ResourceID: 4, Element: "compositematerials", ModelPath: "/3d/3dmodel.model", Name: "matid"},
		go3mf.MissingPropertyError{ResourceID: 4, Element: "compositematerials", ModelPath: "/3d/3dmodel.model", Name: "matindices"},
		go3mf.MissingPropertyError{ResourceID: 4, Element: "composite", ModelPath: "/3d/3dmodel.model", Name: "values"},
		go3mf.ParsePropertyError{ResourceID: 4, Element: "composite", Name: "values", Value: "a", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 9, Element: "multiproperties", ModelPath: "/3d/3dmodel.model", Name: "pids", Value: "a", Type: go3mf.PropertyRequired},
		go3mf.MissingPropertyError{ResourceID: 9, Element: "multi", ModelPath: "/3d/3dmodel.model", Name: "pindices"},
		go3mf.MissingPropertyError{ResourceID: 19, Element: "multiproperties", ModelPath: "/3d/3dmodel.model", Name: "pids"},
		go3mf.GenericError{ResourceID: 8, Element: "triangle", ModelPath: "/3d/3dmodel.model", Message: "duplicated triangle indices"},
		go3mf.GenericError{ResourceID: 8, Element: "triangle", ModelPath: "/3d/3dmodel.model", Message: "triangle indices are out of range"},
		go3mf.MissingPropertyError{ResourceID: 15, Element: "beamlattice", ModelPath: "/3d/3dmodel.model", Name: "radius"},
		go3mf.MissingPropertyError{ResourceID: 15, Element: "beamlattice", ModelPath: "/3d/3dmodel.model", Name: "minlength"},
		go3mf.MissingPropertyError{ResourceID: 15, Element: "beam", ModelPath: "/3d/3dmodel.model", Name: "v1"},
		go3mf.MissingPropertyError{ResourceID: 15, Element: "beam", ModelPath: "/3d/3dmodel.model", Name: "v2"},
		go3mf.MissingPropertyError{ResourceID: 15, Element: "ref", ModelPath: "/3d/3dmodel.model", Name: "index"},
		go3mf.ParsePropertyError{ResourceID: 15, Element: "ref", Name: "index", Value: "a", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 22, Element: "object", ModelPath: "/3d/3dmodel.model", Name: "type", Value: "invalid", Type: go3mf.PropertyOptional},
		go3mf.ParsePropertyError{ResourceID: 20, Element: "object", ModelPath: "/3d/3dmodel.model", Name: "UUID", Value: "cb8286808895-4e08-a1fc-be63e033df15", Type: go3mf.PropertyRequired},
		go3mf.GenericError{ResourceID: 20, Element: "object", ModelPath: "/3d/3dmodel.model", Message: "default PID is not supported for component objects"},
		go3mf.ParsePropertyError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Name: "UUID", Value: "cb8286808895-4e08-a1fc-be63e033df16", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Name: "transform", Value: "0 0 0 1 0 0 0 2 -66.4 -87.1 8.8", Type: go3mf.PropertyOptional},
		go3mf.MissingPropertyError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Name: "UUID"},
		go3mf.GenericError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Message: "non-existent referenced object"},
		go3mf.GenericError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Message: "non-object referenced resource"},
		go3mf.MissingPropertyError{ResourceID: 0, Element: "build", ModelPath: "/3d/3dmodel.model", Name: "UUID"},
		go3mf.ParsePropertyError{ResourceID: 20, Element: "item", Name: "transform", Value: "1 0 0 0 2 0 0 0 3 -66.4 -87.1", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyOptional},
		go3mf.GenericError{ResourceID: 20, Element: "item", ModelPath: "/3d/3dmodel.model", Message: "referenced object cannot be have OTHER type"},
		go3mf.MissingPropertyError{ResourceID: 8, Element: "item", ModelPath: "/3d/3dmodel.model", Name: "UUID"},
		go3mf.GenericError{ResourceID: 8, Element: "item", ModelPath: "/3d/3dmodel.model", Message: "non-existent referenced object"},
		go3mf.GenericError{ResourceID: 5, Element: "item", ModelPath: "/3d/3dmodel.model", Message: "non-object referenced resource"},
		go3mf.ParsePropertyError{ResourceID: 15, Element: "item", Name: "UUID", Value: "e9e", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
		go3mf.ParsePropertyError{ResourceID: 0, Element: "build", Name: "UUID", Value: "e9e25302-6428-402e-8633ed2", ModelPath: "/3d/3dmodel.model", Type: go3mf.PropertyRequired},
	}
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
	rootFile := new(modelBuilder).withDefaultModel().withElement(`
		<resources>
			<basematerials>
				<base name="Blue PLA" displaycolor="0000FF" />
				<base />
			</basematerials>
			<basematerials id="a"/>
			<basematerials id="5">
				<base name="Blue PLA" displaycolor="#0000FF" />
				<base name="Red ABS" displaycolor="#FF0000" />
			</basematerials>			
			<m:texture2d id="6" qm:mq="other" path="/3D/Texture/msLogo.png" contenttype="image/png" tilestyleu="wrap" tilestylev="mirror" filter="auto" />
			<m:texture2d id="b" contenttype="image/png" tilestyleu="wrap" tilestylev="mirror" filter="auto" />
			<m:colorgroup id="1">
				<m:color color="#FFFFF" /> <m:color color="#000000" /> <m:color color="#1AB567" /> <m:color color="#DF045A" />
			</m:colorgroup>
			<m:texture2dgroup qm:mq="other" id="2" texid="6">
				<m:tex2coord qm:mq="other" u="b" v="0.5" /> <m:tex2coord u="0.3" v="c" />	<m:tex2coord u="0.5" v="0.8" />	<m:tex2coord u="0.5" v="0.5" />
			</m:texture2dgroup>
			<m:compositematerials id="4" qm:mq="other">
				<m:composite/>
				<m:composite values="a 0.8"/>
			</m:compositematerials>
			<m:multiproperties id="9" qm:mq="other" pids="a 2">
				<m:multi />
			</m:multiproperties>
			<m:multiproperties id="19" />
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
						<triangle v1="2" v2="2" v3="1" />
						<triangle v1="30" v2="2" v3="1" />
						<triangle v1="3" v2="2" v3="1" />
						<triangle v1="1" v2="0" v3="3" />
						<triangle v1="4" v2="5" v3="6" p1="1" />
						<triangle v1="6" v2="7" v3="4" pid="5" p1="1" />
						<triangle v1="0" v2="1" v3="5" pid="2" p1="0" p2="1" p3="2"/>
						<triangle v1="5" v2="4" v3="0" pid="2" p1="3" p2="0" p3="2"/>
						<triangle v1="1" v2="2" v3="6" pid="1" p1="0" p2="1" p3="2"/>
						<triangle v1="6" v2="5" v3="1" pid="1" p1="2" p2="1" p3="3"/>
						<triangle v1="2" v2="3" v3="7" />
						<triangle v1="7" v2="6" v3="2" />
						<triangle v1="3" v2="0" v3="4" />
						<triangle v1="4" v2="7" v3="3" />
					</triangles>
				</mesh>
			</object>
			<object id="15" name="Box" partnumber="e1ef01d4-cbd4-4a62-86b6-9634e2ca198b" type="model">
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
					<b:beamlattice qm:mq="other" radius="1" minlength="0.0001" cap="hemisphere" clippingmode="inside" clippingmesh="8" representationmesh="8">
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
			<object id="22" p:UUID="cb828680-8895-4e08-a1fc-be63e033df15" type="invalid" />
			<object id="20" pid="3" p:UUID="cb8286808895-4e08-a1fc-be63e033df15" type="other">
				<components>
					<component objectid="8" p:path="/2d/2d.model" p:UUID="cb8286808895-4e08-a1fc-be63e033df16" transform="0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
					<component objectid="9" p:UUID="cb828680-8895-4e08-a1fc-be63e033df16"/>
				</components>
			</object>
		</resources>
		<build>
			<item partnumber="bob" objectid="20" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed2" transform="1 0 0 0 2 0 0 0 3 -66.4 -87.1" />
			<item objectid="8" p:path="/3d/other.model"/>
			<item objectid="5" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed4"/>
			<item objectid="15" p:UUID="e9e"/>
		</build>
		<build p:UUID="e9e25302-6428-402e-8633ed2"/>
		<metadata name="Application">go3mf app</metadata>
		<metadata name="p:CustomMetadata1" type="xs:string" preserve="1">CE8A91FB-C44E-4F00-B634-BAA411465F6A</metadata>
		<other />
		`).build()

	t.Run("base", func(t *testing.T) {
		d := new(Decoder)
		d.Strict = false
		d.SetDecompressor(func(r io.Reader) io.ReadCloser { return flate.NewReader(r) })
		d.SetXMLDecoder(func(r io.Reader) XMLDecoder { return xml.NewDecoder(r) })
		if err := d.processRootModel(context.Background(), rootFile, got); err != nil {
			t.Errorf("Decoder.processRootModel() unexpected error = %v", err)
			return
		}
		deep.MaxDiff = 1
		if diff := deep.Equal(d.Warnings, want); diff != nil {
			t.Errorf("Decoder.processRootModel() = %v", diff)
			return
		}
	})
}
