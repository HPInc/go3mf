// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package go3mf

import (
	"bytes"
	"encoding/xml"
	"errors"
	"image/color"
	"reflect"
	"strconv"
	"testing"

	"github.com/go-test/deep"
	"github.com/hpinc/go3mf/spec"
	"github.com/stretchr/testify/mock"
)

func (f *fakeAttr) Marshal3MFAttr(_ spec.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{
		{Name: xml.Name{Space: fakeExtension, Local: "value"}, Value: f.Value},
	}, nil
}

// Marshal3MF encodes the resource.
func (f *fakeAsset) Marshal3MF(x spec.Encoder) error {
	xs := xml.StartElement{Name: xml.Name{Space: fakeExtension, Local: "fakeasset"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: attrID}, Value: strconv.FormatUint(uint64(f.ID), 10)},
	}}
	x.EncodeToken(xs)
	x.EncodeToken(xs.End())
	return nil
}

type mockPackagePart struct {
	mock.Mock
}

func (m *mockPackagePart) Write(arg1 []byte) (int, error) {
	args := m.Called(arg1)
	return args.Int(0), args.Error(1)
}

func (m *mockPackagePart) AddRelationship(args0 Relationship) {
	m.Called(args0)
}

func TestMarshalModel(t *testing.T) {
	Register(fakeSpec.Namespace, new(qmExtension))
	m := &Model{
		Units: UnitMillimeter, Language: "en-US", Path: "/3D/3dmodel.model", Thumbnail: "/thumbnail.png",
		Extensions: []Extension{fakeSpec, fooSpec},
		AnyAttr:    AnyAttr{&fakeAttr{Value: "model_fake"}, &spec.UnknownAttrs{{Name: fooName, Value: "foo1"}}},
		Any: Any{spec.UnknownTokens{
			xml.StartElement{Name: fooName},
			xml.EndElement{Name: fooName},
		}},
		Resources: Resources{
			Assets: []Asset{
				&UnknownAsset{UnknownTokens: spec.UnknownTokens{
					xml.StartElement{Name: fooName, Attr: []xml.Attr{{Name: xml.Name{Local: "n1"}, Value: "v1"}}},
					xml.StartElement{Name: xml.Name{Space: fooName.Space, Local: "child"}},
					xml.EndElement{Name: xml.Name{Space: fooName.Space, Local: "child"}},
					xml.EndElement{Name: fooName},
				}},
				&BaseMaterials{ID: 5, Materials: []Base{
					{Name: "Blue PLA", Color: color.RGBA{0, 0, 255, 255}, AnyAttr: AnyAttr{&spec.UnknownAttrs{{Name: fooName, Value: "foo6"}}}},
					{Name: "Red ABS", Color: color.RGBA{255, 0, 0, 255}},
				}, AnyAttr: AnyAttr{&spec.UnknownAttrs{{Name: fooName, Value: "foo2"}}}}, &fakeAsset{ID: 25}},
			Objects: []*Object{
				{
					ID: 8, Name: "Box 1", PartNumber: "11111111-1111-1111-1111-111111111111", Thumbnail: "/a.png",
					AnyAttr: AnyAttr{&fakeAttr{Value: "object_fake"}, &spec.UnknownAttrs{{Name: fooName, Value: "foo3"}}},
					PID:     1, PIndex: 1, Type: ObjectTypeModel, Mesh: &Mesh{
						Any: Any{spec.UnknownTokens{
							xml.StartElement{Name: fooName},
							xml.EndElement{Name: fooName},
						}},
						Vertices: []Point3D{
							{0, 0, 0}, {100, 0, 0}, {100, 100, 0},
							{0, 100, 0}, {0, 0, 100}, {100, 0, 100},
							{100, 100, 100}, {0, 100, 100}},
						Triangles: []Triangle{
							{V1: 3, V2: 2, V3: 1, PID: 5, P1: 0, P2: 0, P3: 0},
							{V1: 1, V2: 0, V3: 3, PID: 5, P1: 0, P2: 0, P3: 0},
							{V1: 4, V2: 5, V3: 6, PID: 5, P1: 1, P2: 1, P3: 1},
							{V1: 6, V2: 7, V3: 4, PID: 5, P1: 1, P2: 1, P3: 1},
							{V1: 0, V2: 1, V3: 5, PID: 5, P1: 0, P2: 1, P3: 2},
							{V1: 5, V2: 4, V3: 0, PID: 5, P1: 3, P2: 0, P3: 2},
							{V1: 1, V2: 2, V3: 6, PID: 5, P1: 0, P2: 1, P3: 2},
							{V1: 6, V2: 5, V3: 1, PID: 5, P1: 2, P2: 1, P3: 3},
							{V1: 2, V2: 3, V3: 7, PID: 5, P1: 0, P2: 0, P3: 0},
							{V1: 7, V2: 6, V3: 2, PID: 5, P1: 0, P2: 0, P3: 0},
							{V1: 3, V2: 0, V3: 4, PID: 5, P1: 0, P2: 0, P3: 0},
							{V1: 4, V2: 7, V3: 3, PID: 5, P1: 0, P2: 0, P3: 0},
						},
					}},
				{
					ID: 20, Type: ObjectTypeSupport,
					Metadata: []Metadata{{Name: xml.Name{Space: "qm", Local: "CustomMetadata3"}, Type: "xs:boolean", Value: "1"}, {Name: xml.Name{Space: "qm", Local: "CustomMetadata4"}, Type: "xs:boolean", Value: "2"}},
					Components: &Components{Component: []*Component{{ObjectID: 8, Transform: Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1},
						AnyAttr: AnyAttr{&fakeAttr{Value: "component_fake"}, &spec.UnknownAttrs{{Name: fooName, Value: "foo8"}}}}}},
				},
			},
		},
		Build: Build{
			AnyAttr: AnyAttr{&fakeAttr{Value: "build_fake"}, &spec.UnknownAttrs{{Name: fooName, Value: "foo4"}, {Name: fooName, Value: "foo6"}}},
			Items: []*Item{
				{
					ObjectID: 20, PartNumber: "bob", Transform: Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
					Metadata: []Metadata{{Name: xml.Name{Space: "qm", Local: "CustomMetadata3"}, Type: "xs:boolean", Value: "1"}},
				},
				{ObjectID: 21, AnyAttr: AnyAttr{&fakeAttr{Value: "item_fake"}, &spec.UnknownAttrs{{Name: fooName, Value: "foo5"}}}},
			}}, Metadata: []Metadata{
			{Name: xml.Name{Local: "Application"}, Value: "go3mf app"},
			{Name: xml.Name{Space: "qm", Local: "CustomMetadata1"}, Preserve: true, Type: "xs:string", Value: "CE8A91FB-C44E-4F00-B634-BAA411465F6A"},
		}}

	t.Run("base", func(t *testing.T) {
		b, err := MarshalModel(m)
		if err != nil {
			t.Errorf("MarshalModel() error = %v", err)
			return
		}
		newModel := new(Model)
		newModel.Path = m.Path
		if err := UnmarshalModel(b, newModel); err != nil {
			t.Errorf("MarshalModel() error decoding = %v, s = %s", err, string(b))
			return
		}
		if diff := deep.Equal(m, newModel); diff != nil {
			t.Errorf("MarshalModel() = %v, s = %s", diff, string(b))
		}
	})
}

func TestEncoder_writeAttachements(t *testing.T) {
	type args struct {
		m *Model
	}
	tests := []struct {
		name    string
		e       *Encoder
		args    args
		wantErr bool
	}{
		{"empty", &Encoder{}, args{new(Model)}, false},
		{"err-create", &Encoder{}, args{&Model{Attachments: []Attachment{{}}}}, true},
		{"base", &Encoder{}, args{&Model{Attachments: []Attachment{{Stream: new(bytes.Buffer)}}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var argErr error
			if tt.wantErr {
				argErr = errors.New("")
			}
			m := new(mockPackage)
			mp := new(mockPackagePart)
			mp.On("Write", mock.Anything).Return(mock.Anything, mock.Anything)
			m.On("Create", mock.Anything, mock.Anything).Return(mp, argErr)
			m.On("AddRelationship", mock.Anything).Return()
			tt.e.w = m
			if err := tt.e.writeAttachements(tt.args.m.Attachments); (err != nil) != tt.wantErr {
				t.Errorf("Encoder.writeAttachements() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEncoder(t *testing.T) {
	tests := []struct {
		name string
		want *Encoder
	}{
		{"base", &Encoder{FloatPrecision: defaultFloatPrecision, w: newOpcWriter(nil)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEncoder(nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEncoder() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestEncoder_Encode_Normalize(t *testing.T) {
	type args struct {
		m *Model
	}
	tests := []struct {
		name string
		args args
		want *Model
	}{
		{"empty", args{new(Model)}, &Model{Path: DefaultModelPath}},
		{"withAttrs", args{&Model{Path: "a/other.ml", Thumbnail: "/Metadata/thumbnail.png", Attachments: []Attachment{
			{ContentType: "image/png", Path: "Metadata/thumbnail.png", Stream: bytes.NewBufferString("fake")},
		}}}, &Model{Path: "/a/other.ml", Units: UnitMillimeter, Thumbnail: "/Metadata/thumbnail.png", Attachments: []Attachment{
			{ContentType: "image/png", Path: "/Metadata/thumbnail.png", Stream: bytes.NewBufferString("fake")},
		}, RootRelationships: []Relationship{
			{Path: "/Metadata/thumbnail.png", Type: RelTypeThumbnail, ID: "rId1"},
		}}},
		{"withRootRel", args{&Model{
			RootRelationships: []Relationship{
				{Path: "Metadata/thumbnail.png", Type: RelTypeThumbnail, ID: "2"},
				{Path: "Metadata/thumbnail.png", Type: RelTypeThumbnail, ID: "2"},
			},
			Attachments: []Attachment{
				{ContentType: "application/vnd.ms-printing.printticket+xml", Path: "/3D/Metadata/pt.xml", Stream: bytes.NewBufferString("other")},
				{ContentType: "image/png", Path: "/Metadata/thumbnail.png", Stream: bytes.NewBufferString("fake")},
			}}},
			&Model{Path: DefaultModelPath,
				RootRelationships: []Relationship{
					{Path: "Metadata/thumbnail.png", Type: RelTypeThumbnail, ID: "2"},
				},
				Attachments: []Attachment{
					{ContentType: "image/png", Path: "/Metadata/thumbnail.png", Stream: bytes.NewBufferString("fake")},
				}}},
		{"withChildModel", args{&Model{
			Childs: map[string]*ChildModel{
				"empty.model": {},
				"/other.model": {Relationships: []Relationship{
					{Path: "/3D/Metadata/pt.xml", Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/printticket", ID: "1"}},
				},
			}}}, &Model{Path: DefaultModelPath,
			Childs: map[string]*ChildModel{
				"/3D/empty.model": {},
				"/other.model":    {},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buff := new(bytes.Buffer)
			if err := NewEncoder(buff).Encode(tt.args.m); err != nil {
				t.Errorf("Encoder.Encode() error = %v", err)
				return
			}
			newModel := new(Model)
			err := NewDecoder(bytes.NewReader(buff.Bytes()), int64(buff.Len())).Decode(newModel)
			if err != nil {
				t.Errorf("Encoder.Encode() malformed = %v", err)
				return
			}
			if tt.args.m.Path == "" {
				tt.args.m.Path = DefaultModelPath
			}
			if diff := deep.Equal(newModel, tt.want); diff != nil {
				t.Errorf("MarshalModel() = %v", diff)
			}
		})
	}
}

func TestEncoder_Encode_Roundtrip(t *testing.T) {
	type args struct {
		m *Model
	}
	tests := []struct {
		name string
		args args
	}{
		{"empty", args{new(Model)}},
		{"withAttrs", args{&Model{Path: "/a/other.ml", Language: "un", Units: UnitFoot, Thumbnail: "/thumb.png"}}},
		{"withMetdata", args{&Model{Metadata: []Metadata{
			{Name: xml.Name{Local: "a"}, Value: "b", Type: "tp", Preserve: true},
			{Name: xml.Name{Local: "ab"}, Value: "bb", Type: "tpb", Preserve: false},
		}}}},
		{"withRootRel", args{&Model{
			RootRelationships: []Relationship{
				{Path: "/3D/Metadata/pt.xml", Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/printticket", ID: "1"},
				{Path: "/Metadata/thumbnail.png", Type: "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail", ID: "2"},
			},
			Attachments: []Attachment{
				{ContentType: "application/vnd.ms-printing.printticket+xml", Path: "/3D/Metadata/pt.xml", Stream: bytes.NewBufferString("other")},
				{ContentType: "image/png", Path: "/Metadata/thumbnail.png", Stream: bytes.NewBufferString("fake")},
			}}},
		},
		{"withChildModel", args{&Model{
			Attachments: []Attachment{
				{ContentType: "application/vnd.ms-printing.printticket+xml", Path: "/3D/Metadata/pt.xml", Stream: bytes.NewBufferString("other")},
			},
			Childs: map[string]*ChildModel{
				"/empty.model": {},
				"/other.model": {Relationships: []Relationship{
					{Path: "/3D/Metadata/pt.xml", Type: "http://schemas.microsoft.com/3dmanufacturing/2013/01/printticket", ID: "1"}},
				},
			}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buff := new(bytes.Buffer)
			if err := NewEncoder(buff).Encode(tt.args.m); err != nil {
				t.Errorf("Encoder.Encode() error = %v", err)
				return
			}
			newModel := new(Model)
			err := NewDecoder(bytes.NewReader(buff.Bytes()), int64(buff.Len())).Decode(newModel)
			if err != nil {
				t.Errorf("Encoder.Encode() malformed = %v", err)
				return
			}
			if tt.args.m.Path == "" {
				tt.args.m.Path = DefaultModelPath
			}
			if diff := deep.Equal(newModel, tt.args.m); diff != nil {
				t.Errorf("MarshalModel() = %v", diff)
			}
		})
	}
}
