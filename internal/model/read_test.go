package model

import (
	"bytes"
	"errors"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"testing"
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

func newMockPackage(relationships []relationship, other *mockFile) *mockPackage {
	m := new(mockPackage)
	m.On("Relationships").Return(relationships).Maybe()
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

func (m *mockPackage) Relationships() []relationship {
	args := m.Called()
	return args.Get(0).([]relationship)
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
	type args struct {
		model *Model
	}
	tests := []struct {
		name    string
		d       *Decoder
		args    args
		want    *Model
		wantErr bool
	}{
		{"noRoot", &Decoder{r: newMockPackage(nil, nil)}, args{new(Model)}, new(Model), true},
		{"noRelsFail", &Decoder{r: newMockPackage(nil, newMockFile("/a.model", nil, nil, nil, true))}, args{new(Model)}, &Model{RootPath: "/a.model"}, true},
		{"noRels", &Decoder{r: newMockPackage(nil, newMockFile("/a.model", nil, nil, nil, false))}, args{new(Model)}, &Model{RootPath: "/a.model"}, false},
		{"withThumb", &Decoder{r: newMockPackage(nil,
			newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbFile, thumbFile, false)),
		}, args{new(Model)}, &Model{
			RootPath:    "/a.model",
			Thumbnail:   &Attachment{RelationshipType: relTypeThumbnail, Path: "/Metadata/thumbnail.png", Stream: new(bytes.Buffer)},
			Attachments: []*Attachment{{RelationshipType: relTypeThumbnail, Path: "/a.png", Stream: new(bytes.Buffer)}},
		}, false},
		{"withThumbErr", &Decoder{r: newMockPackage(nil,
			newMockFile("/a.model", []relationship{newMockRelationship(relTypeThumbnail, "/a.png")}, thumbErr, thumbErr, false)),
		}, args{new(Model)}, &Model{RootPath: "/a.model"}, false},
		{"withOtherRel", &Decoder{r: newMockPackage(nil,
			newMockFile("/a.model", []relationship{newMockRelationship("other", "/a.png")}, nil, nil, false)),
		}, args{new(Model)}, &Model{RootPath: "/a.model"}, false},
		{"withModelAttachment", &Decoder{r: newMockPackage(nil,
			newMockFile("/a.model", []relationship{newMockRelationship(relTypeModel3D, "/a.model")}, nil, newMockFile("/a.model", nil, nil, nil, false), false)),
		}, args{new(Model)}, &Model{
			RootPath:              "/a.model",
			ProductionAttachments: []*Attachment{{RelationshipType: relTypeModel3D, Path: "/a.model", Stream: new(bytes.Buffer)}},
		}, false},
		{"withAttRel", &Decoder{AttachmentRelations: []string{"b"}, r: newMockPackage(nil,
			newMockFile("/a.model", []relationship{newMockRelationship("b", "/a.xml")}, nil, newMockFile("/a.xml", nil, nil, nil, false), false)),
		}, args{new(Model)}, &Model{
			RootPath:    "/a.model",
			Attachments: []*Attachment{{RelationshipType: "b", Path: "/a.xml", Stream: new(bytes.Buffer)}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.d.processOPC(tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.processOPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(tt.args.model, tt.want); diff != nil {
				t.Errorf("Decoder.processOPC() = %v", diff)
				return
			}
		})
	}
}
