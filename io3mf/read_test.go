package io3mf

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/mesh"
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
	m.addAttr("xmlns", "b", nsBeamLatticeSpec).addAttr("xmlns", "s", nsSliceSpec).addAttr("", "requiredextensions", "m p b s")
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
	m.On("FindFileFromRel", mock.Anything).Return(other, other != nil).Maybe()
	m.On("FindFileFromName", mock.Anything).Return(other, other != nil).Maybe()
	return m
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

func TestReader_processOPC(t *testing.T) {
	abortReader := &Reader{Model: new(go3mf.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}
	abortReader.SetProgressCallback(callbackFalse, nil)
	thumbFile := newMockFile("/a.png", nil, nil, nil, false)
	thumbErr := newMockFile("/a.png", nil, nil, nil, true)
	tests := []struct {
		name    string
		d       *Reader
		want    *go3mf.Model
		wantErr bool
	}{
		{"noRoot", &Reader{Model: new(go3mf.Model), r: newMockPackage(nil)}, &go3mf.Model{}, true},
		{"abort", abortReader, &go3mf.Model{}, true},
		{"noRels", &Reader{Model: new(go3mf.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}, &go3mf.Model{Path: "/a.model"}, false},
		{"withThumb", &Reader{Model: new(go3mf.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbFile, thumbFile, false)),
		}, &go3mf.Model{
			Path:        "/a.model",
			Thumbnail:   &go3mf.Attachment{RelationshipType: relTypeThumbnail, Path: "/Metadata/thumbnail.png", Stream: new(bytes.Buffer)},
			Attachments: []*go3mf.Attachment{{RelationshipType: relTypeThumbnail, Path: "/a.png", Stream: new(bytes.Buffer)}},
		}, false},
		{"withThumbErr", &Reader{Model: new(go3mf.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbErr, thumbErr, false)),
		}, &go3mf.Model{Path: "/a.model"}, false},
		{"withOtherRel", &Reader{Model: new(go3mf.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("other", "/a.png")}, nil, nil, false)),
		}, &go3mf.Model{Path: "/a.model"}, false},
		{"withModelAttachment", &Reader{Model: new(go3mf.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeModel3D, "/a.model")}, nil, newMockFile("/a.model", nil, nil, nil, false), false)),
		}, &go3mf.Model{
			Path:                  "/a.model",
			ProductionAttachments: []*go3mf.ProductionAttachment{{RelationshipType: relTypeModel3D, Path: "/a.model"}},
		}, false},
		{"withAttRel", &Reader{Model: new(go3mf.Model), AttachmentRelations: []string{"b"},
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("b", "/a.xml")}, nil, newMockFile("/a.xml", nil, nil, nil, false), false)),
		}, &go3mf.Model{
			Path:        "/a.model",
			Attachments: []*go3mf.Attachment{{RelationshipType: "b", Path: "/a.xml", Stream: new(bytes.Buffer)}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.d.processOPC()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.processOPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(tt.d.Model, tt.want); diff != nil {
				t.Errorf("Reader.processOPC() = %v", diff)
				return
			}
		})
	}
}

func TestReader_processRootModel_Fail(t *testing.T) {
	abortReader := &Reader{Model: new(go3mf.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}
	abortReader.SetProgressCallback(callbackFalse, nil)
	tests := []struct {
		name    string
		r       *Reader
		wantErr bool
	}{
		{"noRoot", &Reader{Model: new(go3mf.Model), r: newMockPackage(nil)}, true},
		{"abort", abortReader, true},
		{"errOpen", &Reader{Model: new(go3mf.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, true))}, true},
		{"errEncode", &Reader{Model: new(go3mf.Model), r: newMockPackage(new(modelBuilder).withEncoding("utf16").build())}, true},
		{"invalidUnits", &Reader{Model: new(go3mf.Model), r: newMockPackage(new(modelBuilder).withModel("other", "en-US").build())}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.processRootModel(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.processRootModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReader_processRootModel(t *testing.T) {
	baseMaterials := &go3mf.BaseMaterialsResource{ID: 5, ModelPath: "/3d/3dmodel.model", Materials: []go3mf.BaseMaterial{
		{Name: "Blue PLA", Color: color.RGBA{0, 0, 85, 255}},
		{Name: "Red ABS", Color: color.RGBA{85, 0, 0, 255}},
	}}
	baseTexture := &go3mf.Texture2DResource{ID: 6, ModelPath: "/3d/3dmodel.model", Path: "/3D/Texture/msLogo.png", ContentType: go3mf.PNGTexture, TileStyleU: go3mf.TileWrap, TileStyleV: go3mf.TileMirror, Filter: go3mf.TextureFilterAuto}
	otherSlices := &go3mf.SliceStack{
		BottomZ: 2,
		Slices: []*mesh.Slice{
			{
				TopZ:     1.2,
				Vertices: []mgl32.Vec2{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{0, 1, 2, 3, 0}},
			},
		},
	}
	sliceStack := &go3mf.SliceStackResource{ID: 3, ModelPath: "/3d/3dmodel.model", SliceStack: &go3mf.SliceStack{
		BottomZ: 1,
		Slices: []*mesh.Slice{
			{
				TopZ:     0,
				Vertices: []mgl32.Vec2{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{0, 1, 2, 3, 0}},
			},
			{
				TopZ:     0.1,
				Vertices: []mgl32.Vec2{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{0, 2, 1, 3, 0}},
			},
		},
	}}
	sliceStackRef := &go3mf.SliceStackResource{ID: 7, ModelPath: "/3d/3dmodel.model", SliceStack: otherSlices}
	sliceStackRef.BottomZ = 1.1
	sliceStackRef.UsesSliceRef = true
	sliceStackRef.Slices = append(sliceStackRef.Slices, otherSlices.Slices...)
	meshRes := &go3mf.MeshResource{
		ObjectResource: go3mf.ObjectResource{ID: 8, Name: "Box 1", ModelPath: "/3d/3dmodel.model", SliceStackID: 3, DefaultPropertyID: 5, SliceResoultion: go3mf.ResolutionLow, PartNumber: "11111111-1111-1111-1111-111111111111"},
		Mesh:           new(mesh.Mesh),
	}
	meshRes.Mesh.Nodes = append(meshRes.Mesh.Nodes, []mesh.Node{
		{0, 0, 0},
		{100, 0, 0},
		{100, 100, 0},
		{0, 100, 0},
		{0, 0, 100},
		{100, 0, 100},
		{100, 100, 100},
		{0, 100, 100},
	}...)
	meshRes.Mesh.Faces = append(meshRes.Mesh.Faces, []mesh.Face{
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
		ObjectResource:        go3mf.ObjectResource{ID: 15, Name: "Box", ModelPath: "/3d/3dmodel.model", PartNumber: "e1ef01d4-cbd4-4a62-86b6-9634e2ca198b"},
		BeamLatticeAttributes: go3mf.BeamLatticeAttributes{ClipMode: go3mf.ClipInside, ClippingMeshID: 8, RepresentationMeshID: 8},
		Mesh:                  new(mesh.Mesh),
	}
	meshLattice.Mesh.MinLength = 0.0001
	meshLattice.Mesh.CapMode = mesh.CapModeHemisphere
	meshLattice.Mesh.DefaultRadius = 1
	meshLattice.Mesh.Nodes = append(meshLattice.Mesh.Nodes, []mesh.Node{
		{45, 55, 55},
		{45, 45, 55},
		{45, 55, 45},
		{45, 45, 45},
		{55, 55, 45},
		{55, 55, 55},
		{55, 45, 55},
		{55, 45, 45},
	}...)
	meshLattice.Mesh.BeamSets = append(meshLattice.Mesh.BeamSets, mesh.BeamSet{Name: "test", Identifier: "set_id", Refs: []uint32{1}})
	meshLattice.Mesh.Beams = append(meshLattice.Mesh.Beams, []mesh.Beam{
		{NodeIndices: [2]uint32{0, 1}, Radius: [2]float64{1.5, 1.6}, CapMode: [2]mesh.CapMode{mesh.CapModeSphere, mesh.CapModeButt}},
		{NodeIndices: [2]uint32{2, 0}, Radius: [2]float64{3, 1.5}, CapMode: [2]mesh.CapMode{mesh.CapModeSphere, mesh.CapModeHemisphere}},
		{NodeIndices: [2]uint32{1, 3}, Radius: [2]float64{1.6, 3}, CapMode: [2]mesh.CapMode{mesh.CapModeHemisphere, mesh.CapModeHemisphere}},
		{NodeIndices: [2]uint32{3, 2}, Radius: [2]float64{1, 1}, CapMode: [2]mesh.CapMode{mesh.CapModeHemisphere, mesh.CapModeHemisphere}},
		{NodeIndices: [2]uint32{2, 4}, Radius: [2]float64{3, 2}, CapMode: [2]mesh.CapMode{mesh.CapModeHemisphere, mesh.CapModeHemisphere}},
		{NodeIndices: [2]uint32{4, 5}, Radius: [2]float64{2, 2}, CapMode: [2]mesh.CapMode{mesh.CapModeHemisphere, mesh.CapModeHemisphere}},
		{NodeIndices: [2]uint32{5, 6}, Radius: [2]float64{2, 2}, CapMode: [2]mesh.CapMode{mesh.CapModeHemisphere, mesh.CapModeHemisphere}},
		{NodeIndices: [2]uint32{7, 6}, Radius: [2]float64{2, 2}, CapMode: [2]mesh.CapMode{mesh.CapModeHemisphere, mesh.CapModeHemisphere}},
		{NodeIndices: [2]uint32{1, 6}, Radius: [2]float64{1.6, 2}, CapMode: [2]mesh.CapMode{mesh.CapModeHemisphere, mesh.CapModeHemisphere}},
		{NodeIndices: [2]uint32{7, 4}, Radius: [2]float64{2, 2}, CapMode: [2]mesh.CapMode{mesh.CapModeHemisphere, mesh.CapModeHemisphere}},
		{NodeIndices: [2]uint32{7, 3}, Radius: [2]float64{2, 3}, CapMode: [2]mesh.CapMode{mesh.CapModeHemisphere, mesh.CapModeHemisphere}},
		{NodeIndices: [2]uint32{0, 5}, Radius: [2]float64{1.5, 2}, CapMode: [2]mesh.CapMode{mesh.CapModeHemisphere, mesh.CapModeButt}},
	}...)

	components := &go3mf.ComponentsResource{
		ObjectResource: go3mf.ObjectResource{ID: 20, UUID: "cb828680-8895-4e08-a1fc-be63e033df15", ModelPath: "/3d/3dmodel.model"},
		Components: []*go3mf.Component{{UUID: "cb828680-8895-4e08-a1fc-be63e033df16", Object: meshRes,
			Transform: mgl32.Mat4{3, 0, 0, -66.4, 0, 1, 0, -87.1, 0, 0, 2, 8.8, 0, 0, 0, 1}}},
	}

	want := &go3mf.Model{Units: go3mf.UnitMillimeter, Language: "en-US", Path: "/3d/3dmodel.model", UUID: "e9e25302-6428-402e-8633-cc95528d0ed3"}
	otherMesh := &go3mf.MeshResource{ObjectResource: go3mf.ObjectResource{ID: 8, ModelPath: "/3d/other.model"}, Mesh: new(mesh.Mesh)}
	colorGroup := &go3mf.ColorGroupResource{ID: 1, ModelPath: "/3d/3dmodel.model", Colors: []color.RGBA{{R: 85, G: 85, B: 85, A: 255}, {R: 0, G: 0, B: 0, A: 255}, {R: 16, G: 21, B: 103, A: 255}, {R: 53, G: 4, B: 80, A: 255}}}
	texGroup := &go3mf.Texture2DGroupResource{ID: 2, ModelPath: "/3d/3dmodel.model", TextureID: 6, Coords: []go3mf.TextureCoord{{0.3, 0.5}, {0.3, 0.8}, {0.5, 0.8}, {0.5, 0.5}}}
	want.Resources = append(want.Resources, &go3mf.SliceStackResource{ID: 10, ModelPath: "/2D/2Dmodel.model", SliceStack: otherSlices, TimesRefered: 1})
	want.Resources = append(want.Resources, []go3mf.Identifier{otherMesh, baseMaterials, baseTexture, colorGroup, texGroup, sliceStack, sliceStackRef, meshRes, meshLattice, components}...)
	want.BuildItems = append(want.BuildItems, &go3mf.BuildItem{Object: components, PartNumber: "bob", UUID: "e9e25302-6428-402e-8633-cc95528d0ed2",
		Transform: mgl32.Mat4{1, 0, 0, -66.4, 0, 2, 0, -87.1, 0, 0, 3, 8.8, 0, 0, 0, 1},
	})
	want.BuildItems = append(want.BuildItems, &go3mf.BuildItem{Object: otherMesh})
	got := new(go3mf.Model)
	got.Path = "/3d/3dmodel.model"
	got.Resources = append(got.Resources, &go3mf.SliceStackResource{ID: 10, ModelPath: "/2D/2Dmodel.model", SliceStack: otherSlices}, otherMesh)
	r := &Reader{
		Model: got,
		r: newMockPackage(new(modelBuilder).withDefaultModel().withElement(`
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
				<s:slicestack id="3" zbottom="1">
					<s:slice ztop="0">
						<s:vertices>
							<s:vertex x="1.01" y="1.02" /> <s:vertex x="9.03" y="1.04" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
						</s:vertices>
						<s:polygon startv="0">
							<s:segment v2="1"></s:segment> <s:segment v2="2"></s:segment> <s:segment v2="3"></s:segment> <s:segment v2="0"></s:segment>
						</s:polygon>
					</s:slice>
					<s:slice ztop="0.1">
						<s:vertices>
							<s:vertex x="1.01" y="1.02" /> <s:vertex x="9.03" y="1.04" /> <s:vertex x="9.05" y="9.06" /> <s:vertex x="1.07" y="9.08" />
						</s:vertices>
						<s:polygon startv="0"> 
							<s:segment v2="2"></s:segment> <s:segment v2="1"></s:segment> <s:segment v2="3"></s:segment> <s:segment v2="0"></s:segment>
						</s:polygon>
					</s:slice>
				</s:slicestack>
				<s:slicestack id="7" zbottom="1.1">
					<s:sliceref slicestackid="10" slicepath="/2D/2Dmodel.model" />
				</s:slicestack>
				<object id="8" name="Box 1" pid="5" pindex="0" s:meshresolution="lowres" s:slicestackid="3" partnumber="11111111-1111-1111-1111-111111111111" type="model">
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
					<components>
                		<component objectid="8" p:UUID="cb828680-8895-4e08-a1fc-be63e033df16" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
            		</components>
				</object>
			</resources>
			<build p:UUID="e9e25302-6428-402e-8633-cc95528d0ed3">
				<item partnumber="bob" objectid="20" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed2" transform="1 0 0 0 2 0 0 0 3 -66.4 -87.1 8.8" />
				<item objectid="8" p:path="/3d/other.model" />
			</build>
			<other />
		`).build()),
	}

	t.Run("base", func(t *testing.T) {
		if err := r.processRootModel(); err != nil {
			t.Errorf("Reader.processRootModel() unexpected error = %v", err)
			return
		}
		deep.CompareUnexportedFields = true
		deep.MaxDepth = 20
		if diff := deep.Equal(r.Model, want); diff != nil {
			t.Errorf("Reader.processRootModel() = %v", diff)
			return
		}
	})
}

func TestReader_namespaceRegistered(t *testing.T) {
	type args struct {
		ns string
	}
	tests := []struct {
		name string
		r    *Reader
		args args
		want bool
	}{
		{"empty", &Reader{namespaces: []string{"http://xml.com"}}, args{""}, false},
		{"exist", &Reader{namespaces: []string{"http://xml.com"}}, args{"http://xml.com"}, true},
		{"noexist", &Reader{namespaces: []string{"http://xml.com"}}, args{"xmls"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.namespaceRegistered(tt.args.ns); got != tt.want {
				t.Errorf("Reader.namespaceRegistered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strToMatrix(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    mgl32.Mat4
		wantErr bool
	}{
		{"empty", args{""}, mgl32.Mat4{}, true},
		{"11values", args{"1 1 1 1 1 1 1 1 1 1 1"}, mgl32.Mat4{}, true},
		{"13values", args{"1 1 1 1 1 1 1 1 1 1 1 1 1"}, mgl32.Mat4{}, true},
		{"char", args{"1 1 a 1 1 1 1 1 1 1 1 1"}, mgl32.Mat4{}, true},
		{"base", args{"1 1 1 1 1 1 1 1 1 1 1 1"}, mgl32.Mat4{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1}, false},
		{"other", args{"0 1 2 10 11 12 20 21 22 30 31 32"}, mgl32.Mat4{0, 10, 20, 30, 1, 11, 21, 31, 2, 12, 22, 32, 0, 0, 0, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToMatrix(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("strToMatrix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("strToMatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strToSRGB(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantC   color.RGBA
		wantErr bool
	}{
		{"empty", args{""}, color.RGBA{}, true},
		{"nohashrgb", args{"101010"}, color.RGBA{}, true},
		{"nohashrgba", args{"10101010"}, color.RGBA{}, true},
		{"invalidChar", args{"#€0101010"}, color.RGBA{}, true},
		{"invalidChar", args{"#T0101010"}, color.RGBA{0, 16, 16, 16}, true},
		{"rgb", args{"#112233"}, color.RGBA{17, 34, 51, 255}, false},
		{"rgb", args{"#000233"}, color.RGBA{0, 2, 51, 255}, false},
		{"rgba", args{"#00023311"}, color.RGBA{0, 2, 51, 17}, false},
		{"rgbaLetter", args{"#ff0233AB"}, color.RGBA{85, 2, 51, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := strToSRGB(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("strToSRGB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("strToSRGB() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestReader_processNonRootModels(t *testing.T) {
	abortReader := &Reader{Model: &go3mf.Model{ProductionAttachments: []*go3mf.ProductionAttachment{{}}}}
	abortReader.SetProgressCallback(callbackFalse, nil)
	tests := []struct {
		name    string
		r       *Reader
		wantErr bool
		want    *go3mf.Model
	}{
		{"base", &Reader{Model: &go3mf.Model{ProductionAttachments: []*go3mf.ProductionAttachment{
			{Path: "3d/new.model"},
			{Path: "3d/other.model"},
		}}, productionModels: map[string]packageFile{
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
			Language: "en-US",
			ProductionAttachments: []*go3mf.ProductionAttachment{
				{Path: "3d/new.model"},
				{Path: "3d/other.model"},
			}, Resources: []go3mf.Identifier{
				&go3mf.Texture2DResource{ID: 6, ModelPath: "3d/other.model", Path: "/3D/Texture/msLogo.png", ContentType: go3mf.PNGTexture, TileStyleU: go3mf.TileWrap, TileStyleV: go3mf.TileMirror, Filter: go3mf.TextureFilterAuto},
				&go3mf.BaseMaterialsResource{ID: 5, ModelPath: "3d/new.model", Materials: []go3mf.BaseMaterial{
					{Name: "Blue PLA", Color: color.RGBA{0, 0, 85, 255}},
					{Name: "Red ABS", Color: color.RGBA{85, 0, 0, 255}},
				}},
			},
		}},
		{"noAtt", &Reader{Model: new(go3mf.Model)}, false, new(go3mf.Model)},
		{"abort", abortReader, true, abortReader.Model},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.processNonRootModels(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.processNonRootModels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			deep.CompareUnexportedFields = true
			deep.MaxDepth = 20
			if diff := deep.Equal(tt.r.Model, tt.want); diff != nil {
				t.Errorf("Reader.processNonRootModels() = %v", diff)
				return
			}
		})
	}
}

func TestReader_Decode(t *testing.T) {
	tests := []struct {
		name    string
		r       *Reader
		wantErr bool
	}{
		{"base", &Reader{Model: new(go3mf.Model), AttachmentRelations: []string{"b"},
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("b", "/a.xml")}, nil, newMockFile("/a.xml", nil, nil, nil, false), false)),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Decode(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
