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
	mdl "github.com/qmuntal/go3mf/internal/model"
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

func (m *modelBuilder) withModel(unit mdl.Units, lang string) *modelBuilder {
	m.str.WriteString(`<model `)
	m.addAttr("", "unit", string(unit)).addAttr("xml", "lang", lang)
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
	abortReader := &Reader{Model: new(mdl.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}
	abortReader.SetProgressCallback(callbackFalse, nil)
	thumbFile := newMockFile("/a.png", nil, nil, nil, false)
	thumbErr := newMockFile("/a.png", nil, nil, nil, true)
	tests := []struct {
		name    string
		d       *Reader
		want    *mdl.Model
		wantErr bool
	}{
		{"noRoot", &Reader{Model: new(mdl.Model), r: newMockPackage(nil)}, &mdl.Model{}, true},
		{"abort", abortReader, &mdl.Model{}, true},
		{"noRels", &Reader{Model: new(mdl.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}, &mdl.Model{RootPath: "/a.model"}, false},
		{"withThumb", &Reader{Model: new(mdl.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbFile, thumbFile, false)),
		}, &mdl.Model{
			RootPath:    "/a.model",
			Thumbnail:   &mdl.Attachment{RelationshipType: relTypeThumbnail, Path: "/Metadata/thumbnail.png", Stream: new(bytes.Buffer)},
			Attachments: []*mdl.Attachment{{RelationshipType: relTypeThumbnail, Path: "/a.png", Stream: new(bytes.Buffer)}},
		}, false},
		{"withThumbErr", &Reader{Model: new(mdl.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbErr, thumbErr, false)),
		}, &mdl.Model{RootPath: "/a.model"}, false},
		{"withOtherRel", &Reader{Model: new(mdl.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("other", "/a.png")}, nil, nil, false)),
		}, &mdl.Model{RootPath: "/a.model"}, false},
		{"withModelAttachment", &Reader{Model: new(mdl.Model),
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship(relTypeModel3D, "/a.model")}, nil, newMockFile("/a.model", nil, nil, nil, false), false)),
		}, &mdl.Model{
			RootPath:              "/a.model",
			ProductionAttachments: []*mdl.Attachment{{RelationshipType: relTypeModel3D, Path: "/a.model", Stream: new(bytes.Buffer)}},
		}, false},
		{"withAttRel", &Reader{Model: new(mdl.Model), AttachmentRelations: []string{"b"},
			r: newMockPackage(newMockFile("/a.model", []relationship{newMockRelationship("b", "/a.xml")}, nil, newMockFile("/a.xml", nil, nil, nil, false), false)),
		}, &mdl.Model{
			RootPath:    "/a.model",
			Attachments: []*mdl.Attachment{{RelationshipType: "b", Path: "/a.xml", Stream: new(bytes.Buffer)}},
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

func TestReader_processRootModel(t *testing.T) {
	abortReader := &Reader{Model: new(mdl.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, false))}
	abortReader.SetProgressCallback(callbackFalse, nil)
	baseModel := mdl.NewModel()
	baseMaterials, _ := mdl.NewBaseMaterialsResource(5, baseModel)
	baseTetxure, _ := mdl.NewTexture2DResource(6, baseModel)
	sliceStack, _ := mdl.NewSliceStackResource(3, baseModel, &mdl.SliceStack{
		BottomZ: 1,
		Slices: []*mdl.Slice{
			{
				TopZ:     0,
				Vertices: []mgl32.Vec2{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{1, 2, 3, 0}},
			},
			{
				TopZ:     0.1,
				Vertices: []mgl32.Vec2{{1.01, 1.02}, {9.03, 1.04}, {9.05, 9.06}, {1.07, 9.08}},
				Polygons: [][]int{{2, 1, 3, 0}},
			},
		},
	})
	baseTetxure.Path = "/3D/Texture/msLogo.png"
	baseTetxure.ContentType = mdl.PNGTexture
	baseTetxure.TileStyleV = mdl.TileMirror
	baseMaterials.Materials = []*mdl.BaseMaterial{
		{Name: "Blue PLA", Color: color.RGBA{0, 0, 85, 255}},
		{Name: "Red ABS", Color: color.RGBA{85, 0, 0, 255}},
	}
	baseModel.Resources = []mdl.Identifier{baseMaterials, baseTetxure, sliceStack}
	tests := []struct {
		name    string
		r       *Reader
		want    *mdl.Model
		wantErr bool
	}{
		{"noRoot", &Reader{Model: new(mdl.Model), r: newMockPackage(nil)}, new(mdl.Model), true},
		{"abort", abortReader, new(mdl.Model), true},
		{"errOpen", &Reader{Model: new(mdl.Model), r: newMockPackage(newMockFile("/a.model", nil, nil, nil, true))}, new(mdl.Model), true},
		{"errEncode", &Reader{Model: new(mdl.Model), r: newMockPackage(new(modelBuilder).withEncoding("utf16").build())}, new(mdl.Model), true},
		{"invalidUnits", &Reader{Model: new(mdl.Model), r: newMockPackage(new(modelBuilder).withModel("other", "en-US").build())}, &mdl.Model{}, true},
		{"base", &Reader{Model: mdl.NewModel(), r: newMockPackage(new(modelBuilder).withDefaultModel().withElement(`
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
		</resources>`).build())}, baseModel, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.processRootModel(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.processRootModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(tt.r.Model, tt.want); diff != nil {
				t.Errorf("Reader.processRootModel() = %v", diff)
				return
			}
		})
	}
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
		{"rgb", args{"#112233"}, color.RGBA{17, 34, 51, 255}, false},
		{"rgb", args{"#000233"}, color.RGBA{0, 2, 51, 255}, false},
		{"rgba", args{"#00023311"}, color.RGBA{0, 2, 51, 17}, false},
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
