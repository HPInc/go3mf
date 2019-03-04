package io3mf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/go-test/deep"
	mdl "github.com/qmuntal/go3mf/internal/model"
	"github.com/qmuntal/go3mf/internal/progress"
	"github.com/stretchr/testify/mock"
)

var callbackFalse = func(progress int, id progress.Stage, data interface{}) bool {
	return false
}

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
	hasBuild bool
	hasModel bool
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
	return m.withModel("millimeter", "en-US", "m p b s", "m", "p", "b", "s")
}

func (m *modelBuilder) withModel(unit mdl.Units, lang, requiredExt, nsMaterial, nsProd, nsLattice, nsSlice string) *modelBuilder {
	m.str.WriteString(`<model `)
	m.addAttr("", "unit", string(unit)).addAttr("xml", "lang", lang)
	m.addAttr("", "xmlns", nsCoreSpec).addAttr("xmlns", nsMaterial, nsMaterialSpec).addAttr("xmlns", nsProd, nsProductionSpec)
	m.addAttr("xmlns", nsLattice, nsBeamLatticeSpec).addAttr("xmlns", nsSlice, nsSliceSpec).addAttr("", "requiredextensions", requiredExt)
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
		{"base", &Reader{Model: new(mdl.Model), r: newMockPackage(new(modelBuilder).withDefaultModel().build())}, &mdl.Model{
			Units: mdl.Millimeter,
		}, false},
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

func TestReader_namespaceAttr(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		r    *Reader
		args args
		want string
	}{
		{"empty", &Reader{defaultNamespace: "http://b.com", namespaces: map[string]string{"xml": "http://xml.com"}}, args{""}, ""},
		{"xml", &Reader{defaultNamespace: "http://b.com", namespaces: map[string]string{"xml": "http:/xml.com"}}, args{"xml"}, "http:/xml.com"},
		{"noexist", &Reader{defaultNamespace: "http://b.com", namespaces: map[string]string{"xml": "http:/xml.com"}}, args{"b"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.namespaceAttr(tt.args.prefix); got != tt.want {
				t.Errorf("Reader.namespaceAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_namespaceContent(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		r    *Reader
		args args
		want string
	}{
		{"empty", &Reader{defaultNamespace: "http://b.com", namespaces: map[string]string{"xml": "http://xml.com"}}, args{""}, "http://b.com"},
		{"xml", &Reader{defaultNamespace: "http://b.com", namespaces: map[string]string{"xml": "http:/xml.com"}}, args{"xml"}, "http:/xml.com"},
		{"noexist", &Reader{defaultNamespace: "http://b.com", namespaces: map[string]string{"xml": "http:/xml.com"}}, args{"b"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.namespaceContent(tt.args.prefix); got != tt.want {
				t.Errorf("Reader.namespaceContent() = %v, want %v", got, tt.want)
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
		{"empty", &Reader{namespaces: map[string]string{"xml": "http://xml.com"}}, args{""}, false},
		{"exist", &Reader{namespaces: map[string]string{"xml": "http://xml.com"}}, args{"http://xml.com"}, true},
		{"noexist", &Reader{namespaces: map[string]string{"xml": "http://xml.com"}}, args{"xmls"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.namespaceRegistered(tt.args.ns); got != tt.want {
				t.Errorf("Reader.namespaceRegistered() = %v, want %v", got, tt.want)
			}
		})
	}
}
