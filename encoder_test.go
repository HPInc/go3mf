package go3mf

import (
	"bytes"
	"encoding/xml"
	"errors"
	"image/color"
	"io"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/mock"
)

func TestMarshalModel(t *testing.T) {
	m := &Model{
		Units: UnitMillimeter, Language: "en-US", Path: "/3D/3dmodel.model", Thumbnail: "/thumbnail.png",
		Namespaces:         []xml.Name{{Space: fakeExtension, Local: "qm"}},
		RequiredExtensions: []string{fakeExtension},
		Resources: Resources{
			Assets: []Asset{&BaseMaterialsResource{ID: 5, Materials: []BaseMaterial{
				{Name: "Blue PLA", Color: color.RGBA{0, 0, 255, 255}},
				{Name: "Red ABS", Color: color.RGBA{255, 0, 0, 255}},
			}}},
			Objects: []*Object{
				{
					ID: 8, Name: "Box 1", PartNumber: "11111111-1111-1111-1111-111111111111", Thumbnail: "/a.png",
					DefaultPID: 1, DefaultPIndex: 1, ObjectType: ObjectTypeModel, Mesh: &Mesh{
						Nodes: []Point3D{
							{0, 0, 0}, {100, 0, 0}, {100, 100, 0},
							{0, 100, 0}, {0, 0, 100}, {100, 0, 100},
							{100, 100, 100}, {0, 100, 100}},
						Faces: []Face{
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
						},
					}},
				{
					ID: 20, ObjectType: ObjectTypeSupport,
					Metadata:   []Metadata{{Name: "qm:CustomMetadata3", Type: "xs:boolean", Value: "1"}, {Name: "qm:CustomMetadata4", Type: "xs:boolean", Value: "2"}},
					Components: []*Component{{ObjectID: 8, Transform: Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}}},
				},
			},
		},
		Build: Build{Items: []*Item{
			{
				ObjectID: 20, PartNumber: "bob", Transform: Matrix{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, -66.4, -87.1, 8.8, 1},
				Metadata: []Metadata{{Name: "qm:CustomMetadata3", Type: "xs:boolean", Value: "1"}},
			},
			{ObjectID: 21},
		}}, Metadata: []Metadata{
			{Name: "Application", Value: "go3mf app"},
			{Name: "qm:CustomMetadata1", Preserve: true, Type: "xs:string", Value: "CE8A91FB-C44E-4F00-B634-BAA411465F6A"},
		}}

	t.Run("base", func(t *testing.T) {
		b, err := MarshalModel(m)
		if err != nil {
			t.Errorf("MarshalModel() error = %v", err)
			return
		}
		d := NewDecoder(nil, 0)
		d.RegisterNodeDecoderExtension(fakeExtension, nil)
		d.RegisterDecodeAttributeExtension(fakeExtension, nil)
		newModel := new(Model)
		newModel.Path = m.Path
		if err := d.UnmarshalModel(b, newModel); err != nil {
			t.Errorf("MarshalModel() error decoding = %v, s = %s", err, string(b))
			return
		}
		if diff := deep.Equal(m, newModel); diff != nil {
			t.Errorf("MarshalModel() = %v", diff)
		}
	})
}

func TestEncoder_writeAttachements(t *testing.T) {
	type args struct {
		m *Model
		f io.Writer
	}
	tests := []struct {
		name    string
		e       *Encoder
		args    args
		wantErr bool
	}{
		{"empty", &Encoder{}, args{new(Model), nil}, false},
		{"err-create", &Encoder{}, args{&Model{Attachments: []*Attachment{{}}}, new(bytes.Buffer)}, true},
		{"base", &Encoder{}, args{&Model{Attachments: []*Attachment{{Stream: new(bytes.Buffer)}}}, new(bytes.Buffer)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var argErr error
			if tt.wantErr {
				argErr = errors.New("")
			}
			m := new(mockPackage)
			m.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(tt.args.f, argErr)
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
