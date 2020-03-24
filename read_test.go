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
	"strconv"
	"strings"
	"testing"

	"github.com/go-test/deep"
	specerr "github.com/qmuntal/go3mf/errors"
	"github.com/stretchr/testify/mock"
)

const fakeExtension = "http://dummy.com/fake_ext"

type fakeSpec struct {
}

func (f *fakeSpec) Name() string { return fakeExtension }

func (f *fakeSpec) Required() bool { return true }

func (f *fakeSpec) Local() string { return "f" }

func (f *fakeSpec) ValidateModel(m *Model) []error {
	var errs []error
	var a *fakeAttr
	if m.Build.ExtensionAttr.Get(&a) {
		errs = append(errs, errors.New("Build: fake"))
	}
	return errs
}

type fakeAsset struct {
	ID uint32
}

func (f *fakeAsset) Identify() uint32 {
	return f.ID
}

type fakeAttr struct {
	Value string
}

func (f *fakeAttr) ObjectPath() string { return f.Value }

type fakeAssetDecoder struct {
	baseDecoder
}

func (f *fakeAssetDecoder) Start(att []xml.Attr) {
	id, _ := strconv.ParseUint(att[0].Value, 10, 32)
	f.Scanner.ResourceID = uint32(id)
	f.Scanner.AddAsset(&fakeAsset{ID: uint32(id)})
}

func nodeDecoder(_ interface{}, nodeName string) NodeDecoder {
	return &fakeAssetDecoder{}
}

func decodeAttribute(s *Scanner, parentNode interface{}, attr xml.Attr) {
	switch t := parentNode.(type) {
	case *Object:
		t.ExtensionAttr = append(t.ExtensionAttr, &fakeAttr{attr.Value})
	case *Build:
		t.ExtensionAttr = append(t.ExtensionAttr, &fakeAttr{attr.Value})
	case *Model:
		t.ExtensionAttr = append(t.ExtensionAttr, &fakeAttr{attr.Value})
	case *Item:
		t.ExtensionAttr = append(t.ExtensionAttr, &fakeAttr{attr.Value})
	case *Component:
		t.ExtensionAttr = append(t.ExtensionAttr, &fakeAttr{attr.Value})
	}
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
	m.withModel("millimeter", "en-US", "/thumbnail.png")
	return m
}

func (m *modelBuilder) withModel(unit string, lang string, thumbnail string) *modelBuilder {
	m.str.WriteString(`<model `)
	m.addAttr("", "unit", unit).addAttr("xml", "lang", lang)
	m.addAttr("", "xmlns", ExtensionName).addAttr("xmlns", "qm", fakeExtension)
	m.addAttr("", "requiredextensions", "qm")
	if thumbnail != "" {
		m.addAttr("", "thumbnail", thumbnail)
	}
	m.str.WriteString(">\n")
	m.hasModel = true
	return m
}

func (m *modelBuilder) withEncoding(encode string) *modelBuilder {
	m.str.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="%s"?>`, encode))
	m.str.WriteString("\n")
	return m
}

func (m *modelBuilder) build(name string) *mockFile {
	if name == "" {
		name = "/3D/3dmodel.model"
	}
	if m.hasModel {
		m.str.WriteString("</model>\n")
	}
	f := new(mockFile)
	f.On("Name").Return(name).Maybe()
	f.On("Open").Return(ioutil.NopCloser(bytes.NewBufferString(m.str.String())), nil).Maybe()
	return f
}

type mockFile struct {
	mock.Mock
}

func newMockFile(name string, relationships []Relationship, other *mockFile, openErr bool) *mockFile {
	m := new(mockFile)
	m.On("Name").Return(name).Maybe()
	m.On("ContentType").Return("").Maybe()
	m.On("Relationships").Return(relationships).Maybe()
	m.On("FindFileFromName", mock.Anything).Return(other, other != nil).Maybe()
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

func (m *mockFile) ContentType() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockFile) FindFileFromName(args0 string) (packageFile, bool) {
	args := m.Called(args0)
	return args.Get(0).(packageFile), args.Bool(1)
}

func (m *mockFile) Relationships() []Relationship {
	args := m.Called()
	return args.Get(0).([]Relationship)
}

type mockPackage struct {
	mock.Mock
}

func newMockPackage(other *mockFile) *mockPackage {
	m := new(mockPackage)
	m.On("Open", mock.Anything).Return(nil).Maybe()
	m.On("Create", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	m.On("Relationships").Return([]Relationship{{Path: DefaultModelPath, Type: RelType3DModel}}).Maybe()
	m.On("FindFileFromName", mock.Anything).Return(other, other != nil).Maybe()
	return m
}

func (m *mockPackage) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockPackage) Relationships() []Relationship {
	args := m.Called()
	return args.Get(0).([]Relationship)
}

func (m *mockPackage) AddRelationship(args0 Relationship) {
	m.Called(args0)
}

func (m *mockPackage) Create(args0, args1 string) (packagePart, error) {
	args := m.Called(args0, args1)
	return args.Get(0).(packagePart), args.Error(1)
}

func (m *mockPackage) Open(f func(r io.Reader) io.ReadCloser) error {
	args := m.Called(f)
	return args.Error(0)
}

func (m *mockPackage) FindFileFromName(args0 string) (packageFile, bool) {
	args := m.Called(args0)
	return args.Get(0).(packageFile), args.Bool(1)
}

func TestDecoder_processOPC(t *testing.T) {
	extType := "fake_type"
	otherModel := newMockFile("/other.model", nil, nil, false)
	tests := []struct {
		name    string
		d       *Decoder
		want    *Model
		wantErr bool
	}{
		{"noRoot", &Decoder{p: newMockPackage(nil)}, &Model{}, true},
		{"noRels", &Decoder{p: newMockPackage(newMockFile("/a.model", nil, nil, false))}, &Model{Path: "/a.model"}, false},
		{"withThumb", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []Relationship{{Type: RelTypeThumbnail, Path: "/a.png"}}, newMockFile("/a.png", nil, nil, false), false)),
		}, &Model{
			Path:          "/a.model",
			Relationships: []Relationship{{Path: "/a.png", Type: RelTypeThumbnail}},
			Attachments:   []Attachment{{Path: "/a.png", Stream: new(bytes.Buffer)}},
		}, false},
		{"withPrintTicket", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []Relationship{{Type: RelTypePrintTicket, Path: "/pc.png"}}, newMockFile("/pc.png", nil, nil, false), false)),
		}, &Model{
			Path:          "/a.model",
			Relationships: []Relationship{{Path: "/pc.png", Type: RelTypePrintTicket}},
			Attachments:   []Attachment{{Path: "/pc.png", Stream: new(bytes.Buffer)}},
		}, false},
		{"withExtRel", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []Relationship{{Type: extType, Path: "/other.png"}}, newMockFile("/other.png", nil, nil, false), false)),
		}, &Model{
			Path:          "/a.model",
			Relationships: []Relationship{{Path: "/other.png", Type: extType}},
			Attachments:   []Attachment{{Path: "/other.png", Stream: new(bytes.Buffer)}},
		}, false},
		{"withOtherRel", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []Relationship{{Type: "other", Path: "/a.png"}}, nil, false)),
		}, &Model{Path: "/a.model"}, false},
		{"withModelAttachment", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []Relationship{{Type: RelType3DModel, Path: "/other.model"}}, otherModel, false)),
		}, &Model{Path: "/a.model", Childs: map[string]*ChildModel{"/other.model": new(ChildModel)}}, false},
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
		{"errOpen", newMockFile("/a.model", nil, nil, true), true},
		{"errEncode", new(modelBuilder).withEncoding("utf16").build(""), true},
		{"invalidUnits", new(modelBuilder).withModel("other", "en-US", "").build(""), false},
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
	baseMaterials := &BaseMaterials{ID: 5, Materials: []Base{
		{Name: "Blue PLA", Color: color.RGBA{0, 0, 255, 255}},
		{Name: "Red ABS", Color: color.RGBA{255, 0, 0, 255}},
	}}
	meshRes := &Object{
		Mesh: new(Mesh),
		ID:   8, Name: "Box 1", Thumbnail: "/a.png", DefaultPID: 5, PartNumber: "11111111-1111-1111-1111-111111111111",
	}
	meshRes.Mesh.Nodes = append(meshRes.Mesh.Nodes, []Point3D{
		{0, 0, 0},
		{100, 0, 0},
		{100, 100, 0},
		{0, 100, 0},
		{0, 0, 100},
		{100, 0, 100},
		{100, 100, 100},
		{0, 100, 100},
	}...)
	meshRes.Mesh.Faces = append(meshRes.Mesh.Faces, []Face{
		{NodeIndices: [3]uint32{3, 2, 1}, PID: 5},
		{NodeIndices: [3]uint32{1, 0, 3}, PID: 5},
		{NodeIndices: [3]uint32{4, 5, 6}, PID: 5, PIndex: [3]uint32{1, 1, 1}},
		{NodeIndices: [3]uint32{6, 7, 4}, PID: 5, PIndex: [3]uint32{1, 1, 1}},
		{NodeIndices: [3]uint32{0, 1, 5}, PID: 5, PIndex: [3]uint32{0, 1, 2}},
		{NodeIndices: [3]uint32{5, 4, 0}, PID: 5, PIndex: [3]uint32{3, 0, 2}},
		{NodeIndices: [3]uint32{1, 2, 6}, PID: 5, PIndex: [3]uint32{0, 1, 2}},
		{NodeIndices: [3]uint32{6, 5, 1}, PID: 5, PIndex: [3]uint32{2, 1, 3}},
		{NodeIndices: [3]uint32{2, 3, 7}, PID: 5},
		{NodeIndices: [3]uint32{7, 6, 2}, PID: 5},
		{NodeIndices: [3]uint32{3, 0, 4}, PID: 5},
		{NodeIndices: [3]uint32{4, 7, 3}, PID: 5},
	}...)

	components := &Object{
		ID: 20, ObjectType: ObjectTypeSupport,
		Metadata:   []Metadata{{Name: xml.Name{Space: "qm", Local: "CustomMetadata3"}, Type: "xs:boolean", Value: "1"}, {Name: xml.Name{Space: "qm", Local: "CustomMetadata4"}, Type: "xs:boolean", Value: "2"}},
		Components: []*Component{{ObjectID: 8, Transform: Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}}},
	}

	want := &Model{
		Units: UnitMillimeter, Language: "en-US", Path: "/3D/3dmodel.model", Thumbnail: "/thumbnail.png",
		Namespaces:         []xml.Name{{Space: fakeExtension, Local: "qm"}},
		RequiredExtensions: []string{fakeExtension},
		Resources: Resources{
			Assets: []Asset{baseMaterials}, Objects: []*Object{meshRes, components},
		},
	}
	want.Build.Items = append(want.Build.Items, &Item{
		ObjectID: 20, PartNumber: "bob", Transform: Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
		Metadata: []Metadata{{Name: xml.Name{Space: "qm", Local: "CustomMetadata3"}, Type: "xs:boolean", Value: "1"}},
	})
	want.Metadata = append(want.Metadata, []Metadata{
		{Name: xml.Name{Local: "Application"}, Value: "go3mf app"},
		{Name: xml.Name{Space: "qm", Local: "CustomMetadata1"}, Preserve: true, Type: "xs:string", Value: "CE8A91FB-C44E-4F00-B634-BAA411465F6A"},
	}...)
	got := new(Model)
	got.Path = "/3D/3dmodel.model"
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
			<object id="20" type="support">
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
		`).build("")

	t.Run("base", func(t *testing.T) {
		d := new(Decoder)
		d.RegisterNodeDecoderExtension(fakeExtension, nil)
		d.RegisterDecodeAttributeExtension(fakeExtension, nil)
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
		{"base", &Model{Childs: map[string]*ChildModel{"/3D/other.model": new(ChildModel), "/3D/new.model": new(ChildModel)}},
			&Decoder{nonRootModels: []packageFile{
				new(modelBuilder).withDefaultModel().withElement(`
				<resources>
					<basematerials id="5">
						<base name="Blue PLA" displaycolor="#0000FF" />
						<base name="Red ABS" displaycolor="#FF0000" />
					</basematerials>
				</resources>
			`).build("/3D/new.model"),
				new(modelBuilder).withDefaultModel().withElement(`
				<resources>
					<basematerials id="6" />
				</resources>
			`).build("/3D/other.model"),
			}}, false, &Model{
				Childs: map[string]*ChildModel{
					"/3D/other.model": {Resources: Resources{Assets: []Asset{&BaseMaterials{ID: 6}}}},
					"/3D/new.model": {Resources: Resources{Assets: []Asset{
						&BaseMaterials{ID: 5, Materials: []Base{
							{Name: "Blue PLA", Color: color.RGBA{0, 0, 255, 255}},
							{Name: "Red ABS", Color: color.RGBA{255, 0, 0, 255}},
						}}}}},
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
		{"base", &Decoder{
			p: newMockPackage(newMockFile("/a.model", []Relationship{{Type: "b", Path: "/a.xml"}}, nil, false)),
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
		args    args
		wantErr bool
	}{
		{"nochild", args{context.Background(), xml.NewDecoder(bytes.NewBufferString(`
			<a></a>
			<b></b>
		`))}, false},
		{"eof", args{context.Background(), xml.NewDecoder(bytes.NewBufferString(`
			<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02">
				<build></build>
		`))}, true},
		{"canceled", args{ctx, xml.NewDecoder(bytes.NewBufferString(`
			<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02">
				<build></build>
			</model>
		`))}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if sc := decodeModelFile(tt.args.ctx, tt.args.x, new(Model), "", true, false, nil); (sc.Err != nil) != tt.wantErr {
				t.Errorf("modelFile.Decode() error = %v, wantErr %v", sc.Err, tt.wantErr)
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
		{"base", args{nil, 5}, &Decoder{
			Strict:           true,
			p:                &opcReader{ra: nil, size: 5},
			extensionDecoder: make(map[string]*extensionDecoderWrapper),
		}},
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
		&specerr.ParseFieldError{ResourceID: 0, Element: "base", Name: "displaycolor", Value: "0000FF", ModelPath: "/3D/3dmodel.model", Required: true},
		&specerr.ParseFieldError{ResourceID: 0, Element: "basematerials", Name: "id", Value: "a", ModelPath: "/3D/3dmodel.model", Required: true},
		&specerr.ParseFieldError{ResourceID: 8, Element: "vertex", Name: "x", ModelPath: "/3D/3dmodel.model", Value: "a", Required: true},
		&specerr.ParseFieldError{ResourceID: 8, Element: "triangle", ModelPath: "/3D/3dmodel.model", Name: "v1", Value: "a", Required: true},
		&specerr.ParseFieldError{ResourceID: 22, Element: "object", ModelPath: "/3D/3dmodel.model", Name: "pid", Value: "a", Required: false},
		&specerr.ParseFieldError{ResourceID: 22, Element: "object", ModelPath: "/3D/3dmodel.model", Name: "pindex", Value: "a", Required: false},
		&specerr.ParseFieldError{ResourceID: 22, Element: "object", ModelPath: "/3D/3dmodel.model", Name: "type", Value: "invalid", Required: false},
		&specerr.ParseFieldError{ResourceID: 20, Element: "component", ModelPath: "/3D/3dmodel.model", Name: "transform", Value: "0 0 0 1 0 0 0 2 -66.4 -87.1 8.8", Required: false},
		&specerr.ParseFieldError{ResourceID: 20, Element: "component", ModelPath: "/3D/3dmodel.model", Name: "objectid", Value: "a", Required: true},
		&specerr.ParseFieldError{ResourceID: 20, Element: "item", Name: "transform", Value: "1 0 0 0 2 0 0 0 3 -66.4 -87.1", ModelPath: "/3D/3dmodel.model", Required: false},
		&specerr.ParseFieldError{Element: "item", Name: "objectid", Value: "a", ModelPath: "/3D/3dmodel.model", Required: true},
	}
	got := new(Model)
	got.Path = "/3D/3dmodel.model"
	rootFile := new(modelBuilder).withElement(`
		<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" 
		xmlns:qm="http://dummy.com/fake_ext" requiredextensions="qm other">
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
						<vertex x="a" y="100.00000" z="100.00000" />
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
						<triangle v1="a" v2="7" v3="3" />
					</triangles>
				</mesh>
			</object>
			<object id="22" pid="a" pindex="a" type="invalid" />
			<object id="20" pid="3" type="other">
				<components>
					<component objectid="8" transform="0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
					<component objectid="a"/>
				</components>
			</object>
		</resources>
		<build>
			<item partnumber="bob" objectid="20" transform="1 0 0 0 2 0 0 0 3 -66.4 -87.1" />
			<item objectid="8"/>
			<item objectid="5"/>
			<item objectid="a"/>
		</build>
		<metadata name="Application">go3mf app</metadata>
		<metadata name="qm:CustomMetadata1" type="xs:string" preserve="1">CE8A91FB-C44E-4F00-B634-BAA411465F6A</metadata>
		<metadata name="unknown:CustomMetadata1" type="xs:string" preserve="1">CE8A91FB-C44E-4F00-B634-BAA411465F6A</metadata>
		<other />
		</model>
		`).build("")

	t.Run("base", func(t *testing.T) {
		d := new(Decoder)
		d.RegisterNodeDecoderExtension(fakeExtension, nil)
		d.RegisterDecodeAttributeExtension(fakeExtension, nil)
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
