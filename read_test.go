package go3mf

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
	m.addAttr("", "xmlns", ExtensionName).addAttr("xmlns", "p", nsProductionSpec)
	m.addAttr("", "requiredextensions", "p")
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

func TestDecoder_processOPC(t *testing.T) {
	thumbFile := newMockFile("/a.png", nil, nil, nil, false)
	thumbErr := newMockFile("/a.png", nil, nil, nil, true)
	tests := []struct {
		name    string
		d       *Decoder
		want    *Model
		wantErr bool
	}{
		{"noRoot", &Decoder{p: newMockPackage(nil)}, &Model{}, true},
		{"noRels", &Decoder{p: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}, &Model{Path: "/a.model"}, false},
		{"withThumb", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbFile, thumbFile, false)),
		}, &Model{
			Path:        "/a.model",
			Thumbnail:   &Attachment{RelationshipType: relTypeThumbnail, Path: "/Metadata/thumbnail.png", Stream: new(bytes.Buffer)},
			Attachments: []*Attachment{{RelationshipType: relTypeThumbnail, Path: "/a.png", Stream: new(bytes.Buffer)}},
		}, false},
		{"withThumbErr", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbErr, thumbErr, false)),
		}, &Model{Path: "/a.model"}, false},
		{"withOtherRel", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("other", "/a.png")}, nil, nil, false)),
		}, &Model{Path: "/a.model"}, false},
		{"withModelAttachment", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeModel3D, "/a.model")}, nil, newMockFile("/a.model", nil, nil, nil, false), false)),
		}, &Model{
			Path:                  "/a.model",
			ProductionAttachments: []*ProductionAttachment{{RelationshipType: relTypeModel3D, Path: "/a.model"}},
		}, false},
		{"withAttRel", &Decoder{AttachmentRelations: []string{"b"},
			p: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("b", "/a.xml")}, nil, newMockFile("/a.xml", nil, nil, nil, false), false)),
		}, &Model{
			Path:        "/a.model",
			Attachments: []*Attachment{{RelationshipType: "b", Path: "/a.xml", Stream: new(bytes.Buffer)}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := new(Model)
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
			if err := new(Decoder).processRootModel(context.Background(), tt.f, new(Model)); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.processRootModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDecoder_processRootModel(t *testing.T) {
	baseMaterials := &BaseMaterialsResource{ID: 5, ModelPath: "/3d/3dmodel.model", Materials: []BaseMaterial{
		{Name: "Blue PLA", Color: color.RGBA{0, 0, 255, 255}},
		{Name: "Red ABS", Color: color.RGBA{255, 0, 0, 255}},
	}}
	meshRes := &Mesh{
		ObjectResource: ObjectResource{
			ID: 8, Name: "Box 1", ModelPath: "/3d/3dmodel.model", Thumbnail: "/a.png", DefaultPropertyID: 5, PartNumber: "11111111-1111-1111-1111-111111111111"},
	}
	meshRes.Nodes = append(meshRes.Nodes, []Point3D{
		{0, 0, 0},
		{100, 0, 0},
		{100, 100, 0},
		{0, 100, 0},
		{0, 0, 100},
		{100, 0, 100},
		{100, 100, 100},
		{0, 100, 100},
	}...)
	meshRes.Faces = append(meshRes.Faces, []Face{
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

	components := &Components{
		ObjectResource: ObjectResource{
			ID: 20, UUID: "cb828680-8895-4e08-a1fc-be63e033df15", ModelPath: "/3d/3dmodel.model",
			Metadata: []Metadata{{Name: nsProductionSpec + ":CustomMetadata3", Type: "xs:boolean", Value: "1"}, {Name: nsProductionSpec + ":CustomMetadata4", Type: "xs:boolean", Value: "2"}},
		},
		Components: []*Component{{UUID: "cb828680-8895-4e08-a1fc-be63e033df16", Object: meshRes,
			Transform: Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}}},
	}

	want := &Model{Units: UnitMillimeter, Language: "en-US", Path: "/3d/3dmodel.model", UUID: "e9e25302-6428-402e-8633-cc95528d0ed3"}
	otherMesh := &Mesh{ObjectResource: ObjectResource{ID: 8, ModelPath: "/3d/other.model"}}
	want.Resources = append(want.Resources, otherMesh, baseMaterials, meshRes, components)
	want.BuildItems = append(want.BuildItems, &BuildItem{Object: components, PartNumber: "bob", UUID: "e9e25302-6428-402e-8633-cc95528d0ed2",
		Transform: Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
	})
	want.BuildItems = append(want.BuildItems, &BuildItem{Object: otherMesh, UUID: "e9e25302-6428-402e-8633-cc95528d0ed3", Metadata: []Metadata{{Name: nsProductionSpec + ":CustomMetadata3", Type: "xs:boolean", Value: "1"}}})
	want.Metadata = append(want.Metadata, []Metadata{
		{Name: "Application", Value: "go3mf app"},
		{Name: nsProductionSpec + ":CustomMetadata1", Preserve: true, Type: "xs:string", Value: "CE8A91FB-C44E-4F00-B634-BAA411465F6A"},
	}...)
	got := new(Model)
	got.Path = "/3d/3dmodel.model"
	got.Resources = append(got.Resources, otherMesh)
	rootFile := new(modelBuilder).withDefaultModel().withElement(`
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
		model   *Model
		d       *Decoder
		wantErr bool
		want    *Model
	}{
		{"base", &Model{ProductionAttachments: []*ProductionAttachment{
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
					<basematerials id="6" />
				</resources>
			`).build(),
		}}, false, &Model{
			ProductionAttachments: []*ProductionAttachment{
				{Path: "3d/new.model"},
				{Path: "3d/other.model"},
			}, Resources: []Resource{
				&BaseMaterialsResource{ID: 5, ModelPath: "3d/new.model", Materials: []BaseMaterial{
					{Name: "Blue PLA", Color: color.RGBA{0, 0, 255, 255}},
					{Name: "Red ABS", Color: color.RGBA{255, 0, 0, 255}},
				}},
				&BaseMaterialsResource{ID: 6, ModelPath: "3d/other.model"},
			},
		}},
		{"noAtt", new(Model), new(Decoder), false, new(Model)},
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
			if err := tt.d.Decode(new(Model)); (err != nil) != tt.wantErr {
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
			if err := tt.d.Decode(tt.args.ctx, tt.args.x, new(Model), "", true, false); (err != nil) != tt.wantErr {
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
		ParsePropertyError{ResourceID: 0, Element: "base", Name: "displaycolor", Value: "0000FF", ModelPath: "/3d/3dmodel.model", Type: PropertyRequired},
		MissingPropertyError{ResourceID: 0, Element: "base", ModelPath: "/3d/3dmodel.model", Name: "name"},
		MissingPropertyError{ResourceID: 0, Element: "base", ModelPath: "/3d/3dmodel.model", Name: "displaycolor"},
		MissingPropertyError{ResourceID: 0, Element: "basematerials", ModelPath: "/3d/3dmodel.model", Name: "id"},
		ParsePropertyError{ResourceID: 0, Element: "basematerials", Name: "id", Value: "a", ModelPath: "/3d/3dmodel.model", Type: PropertyRequired},
		MissingPropertyError{ResourceID: 0, Element: "basematerials", ModelPath: "/3d/3dmodel.model", Name: "id"},
		GenericError{ResourceID: 8, Element: "triangle", ModelPath: "/3d/3dmodel.model", Message: "duplicated triangle indices"},
		GenericError{ResourceID: 8, Element: "triangle", ModelPath: "/3d/3dmodel.model", Message: "triangle indices are out of range"},
		ParsePropertyError{ResourceID: 22, Element: "object", ModelPath: "/3d/3dmodel.model", Name: "type", Value: "invalid", Type: PropertyOptional},
		ParsePropertyError{ResourceID: 20, Element: "object", ModelPath: "/3d/3dmodel.model", Name: "UUID", Value: "cb8286808895-4e08-a1fc-be63e033df15", Type: PropertyRequired},
		GenericError{ResourceID: 20, Element: "object", ModelPath: "/3d/3dmodel.model", Message: "default PID is not supported for component objects"},
		ParsePropertyError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Name: "UUID", Value: "cb8286808895-4e08-a1fc-be63e033df16", Type: PropertyRequired},
		ParsePropertyError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Name: "transform", Value: "0 0 0 1 0 0 0 2 -66.4 -87.1 8.8", Type: PropertyOptional},
		MissingPropertyError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Name: "UUID"},
		GenericError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Message: "non-existent referenced object"},
		GenericError{ResourceID: 20, Element: "component", ModelPath: "/3d/3dmodel.model", Message: "non-object referenced resource"},
		MissingPropertyError{ResourceID: 0, Element: "build", ModelPath: "/3d/3dmodel.model", Name: "UUID"},
		ParsePropertyError{ResourceID: 20, Element: "item", Name: "transform", Value: "1 0 0 0 2 0 0 0 3 -66.4 -87.1", ModelPath: "/3d/3dmodel.model", Type: PropertyOptional},
		GenericError{ResourceID: 20, Element: "item", ModelPath: "/3d/3dmodel.model", Message: "referenced object cannot be have OTHER type"},
		MissingPropertyError{ResourceID: 8, Element: "item", ModelPath: "/3d/3dmodel.model", Name: "UUID"},
		GenericError{ResourceID: 8, Element: "item", ModelPath: "/3d/3dmodel.model", Message: "non-existent referenced object"},
		GenericError{ResourceID: 5, Element: "item", ModelPath: "/3d/3dmodel.model", Message: "non-object referenced resource"},
		ParsePropertyError{ResourceID: 0, Element: "build", Name: "UUID", Value: "e9e25302-6428-402e-8633ed2", ModelPath: "/3d/3dmodel.model", Type: PropertyRequired},
	}
	got := new(Model)
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
			<object id="22" p:UUID="cb828680-8895-4e08-a1fc-be63e033df15" type="invalid" />
			<object id="20" pid="3" p:UUID="cb8286808895-4e08-a1fc-be63e033df15" type="other">
				<components>
					<component objectid="8" p:path="/2d/2d.model" p:UUID="cb8286808895-4e08-a1fc-be63e033df16" transform="0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
					<component objectid="5" p:UUID="cb828680-8895-4e08-a1fc-be63e033df16"/>
				</components>
			</object>
		</resources>
		<build>
			<item partnumber="bob" objectid="20" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed2" transform="1 0 0 0 2 0 0 0 3 -66.4 -87.1" />
			<item objectid="8" p:path="/3d/other.model"/>
			<item objectid="5" p:UUID="e9e25302-6428-402e-8633-cc95528d0ed4"/>
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
