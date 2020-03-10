package go3mf

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"testing"

	"github.com/go-test/deep"
	specerr "github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	path := DefaultModelPath
	type args struct {
		model *Model
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{"empty", args{new(Model)}, []error{}},
		{"rels", args{&Model{Attachments: []Attachment{{Path: "/a.png"}}, Relationships: []Relationship{
			{}, {Path: "/.png"}, {Path: "/a.png"}, {Path: "a.png"}, {Path: "/b.png"}, {Path: "/a.png"},
			{Path: "/a.png", Type: RelTypePrintTicket}, {Path: "/a.png", Type: RelTypePrintTicket},
		}}}, []error{
			fmt.Errorf("go3mf: %s@Relationship#0: %v", path, specerr.ErrOPCPartName),
			fmt.Errorf("go3mf: %s@Relationship#1: %v", path, specerr.ErrOPCPartName),
			fmt.Errorf("go3mf: %s@Relationship#3: %v", path, specerr.ErrOPCPartName),
			fmt.Errorf("go3mf: %s@Relationship#4: %v", path, specerr.ErrOPCRelTarget),
			fmt.Errorf("go3mf: %s@Relationship#5: %v", path, specerr.ErrOPCDuplicatedRel),
			fmt.Errorf("go3mf: %s@Relationship#6: %v", path, specerr.ErrOPCContentType),
			fmt.Errorf("go3mf: %s@Relationship#7: %v", path, specerr.ErrOPCDuplicatedRel),
			fmt.Errorf("go3mf: %s@Relationship#7: %v", path, specerr.ErrOPCContentType),
			fmt.Errorf("go3mf: %s@Relationship#7: %v", path, specerr.ErrOPCDuplicatedTicket),
		}},
		{"namespaces", args{&Model{RequiredExtensions: []string{"fake", "other"}, Namespaces: []xml.Name{{Space: "fake", Local: "f"}}}}, []error{
			fmt.Errorf("go3mf: %s: %v", path, specerr.ErrRequiredExt),
		}},
		{"metadata", args{&Model{Namespaces: []xml.Name{{Space: "fake", Local: "f"}}, Metadata: []Metadata{
			{Name: xml.Name{Space: "fake", Local: "issue"}}, {Name: xml.Name{Space: "f", Local: "issue"}}, {Name: xml.Name{Space: "fake", Local: "issue"}}, {Name: xml.Name{Local: "issue"}}, {},
		}}}, []error{
			fmt.Errorf("go3mf: %s@Metadata#1: %v", path, specerr.ErrMetadataNamespace),
			fmt.Errorf("go3mf: %s@Metadata#2: %v", path, specerr.ErrMetadataDuplicated),
			fmt.Errorf("go3mf: %s@Metadata#3: %v", path, specerr.ErrMetadataName),
			fmt.Errorf("go3mf: %s@Metadata#4: %v", path, &specerr.MissingFieldError{Name: attrName}),
		}},
		{"build", args{&Model{Resources: Resources{Assets: []Asset{&BaseMaterialsResource{ID: 1, Materials: []BaseMaterial{{Name: "a", Color: color.RGBA{A: 1}}}}}, Objects: []*Object{
			{ID: 2, ObjectType: ObjectTypeOther, Mesh: &Mesh{Nodes: []Point3D{{}, {}, {}, {}}, Faces: []Face{
				{NodeIndices: [3]uint32{0, 1, 2}}, {NodeIndices: [3]uint32{0, 3, 1}}, {NodeIndices: [3]uint32{0, 2, 3}}, {NodeIndices: [3]uint32{1, 3, 2}},
			}}}}}, Build: Build{ExtensionAttr: ExtensionAttr{&fakeAttr{}}, Items: []*Item{
			{},
			{ObjectID: 2},
			{ObjectID: 100},
			{ObjectID: 1, Metadata: []Metadata{{Name: xml.Name{Local: "issue"}}}},
		}}}}, []error{
			fmt.Errorf("go3mf: %s@Build: ", path),
			fmt.Errorf("go3mf: %s@Build@Item#0: %v", path, &specerr.MissingFieldError{Name: attrObjectID}),
			fmt.Errorf("go3mf: %s@Build@Item#1: %v", path, specerr.ErrOtherItem),
			fmt.Errorf("go3mf: %s@Build@Item#2: %v", path, specerr.ErrMissingResource),
			fmt.Errorf("go3mf: %s@Build@Item#3: %v", path, specerr.ErrMissingResource),
			fmt.Errorf("go3mf: %s@Build@Item#3@Metadata#0: %v", path, specerr.ErrMetadataName),
		}},
		{"childs", args{&Model{Childs: map[string]*ChildModel{path: &ChildModel{}, "/a.model": &ChildModel{
			Relationships: make([]Relationship, 1), Resources: Resources{Objects: []*Object{{}}}}}}},
			[]error{
				fmt.Errorf("go3mf: %s: %v", path, specerr.ErrOPCDuplicatedModelName),
				fmt.Errorf("go3mf: /a.model@Relationship#0: %v", specerr.ErrOPCPartName),
				fmt.Errorf("go3mf: /a.model@Resources@Object#0: %v", specerr.ErrMissingID),
				fmt.Errorf("go3mf: /a.model@Resources@Object#0: %v", specerr.ErrInvalidObject),
			}},
		{"assets", args{&Model{Resources: Resources{Assets: []Asset{
			&BaseMaterialsResource{Materials: []BaseMaterial{{Color: color.RGBA{}}}},
			&BaseMaterialsResource{ID: 1, Materials: []BaseMaterial{{Name: "a", Color: color.RGBA{A: 1}}}},
			&BaseMaterialsResource{ID: 1},
		}}}}, []error{
			fmt.Errorf("go3mf: %s@Resources@BaseMaterialsResource#0: %v", path, specerr.ErrMissingID),
			fmt.Errorf("go3mf: %s@Resources@BaseMaterialsResource#0@BaseMaterial#0: %v", path, &specerr.MissingFieldError{Name: attrName}),
			fmt.Errorf("go3mf: %s@Resources@BaseMaterialsResource#0@BaseMaterial#0: %v", path, &specerr.MissingFieldError{Name: attrDisplayColor}),
			fmt.Errorf("go3mf: %s@Resources@BaseMaterialsResource#2: %v", path, specerr.ErrDuplicatedID),
			fmt.Errorf("go3mf: %s@Resources@BaseMaterialsResource#2: %v", path, specerr.ErrEmptyResourceProps),
		}},
		{"objects", args{&Model{Resources: Resources{Assets: []Asset{
			&BaseMaterialsResource{ID: 1, Materials: []BaseMaterial{{Name: "a", Color: color.RGBA{A: 1}}, {Name: "b", Color: color.RGBA{A: 1}}}},
			&BaseMaterialsResource{ID: 5, Materials: []BaseMaterial{{Name: "a", Color: color.RGBA{A: 1}}, {Name: "b", Color: color.RGBA{A: 1}}}},
		}, Objects: []*Object{
			{},
			{ID: 1, DefaultPIndex: 1, Mesh: &Mesh{}, Components: []*Component{{ObjectID: 1, ExtensionAttr: ExtensionAttr{&fakeAttr{path}}}}},
			{ID: 2, Mesh: &Mesh{Nodes: []Point3D{{}, {}, {}, {}}, Faces: []Face{
				{NodeIndices: [3]uint32{0, 1, 2}}, {NodeIndices: [3]uint32{0, 3, 1}}, {NodeIndices: [3]uint32{0, 2, 3}}, {NodeIndices: [3]uint32{1, 3, 2}},
			}}},
			{ID: 3, DefaultPID: 5, Components: []*Component{
				{ObjectID: 3}, {ObjectID: 2}, {}, {ObjectID: 5}, {ObjectID: 100},
			}},
			{ID: 4, DefaultPID: 100, Mesh: &Mesh{Nodes: make([]Point3D, 2), Faces: make([]Face, 3)}},
			{ID: 6, DefaultPID: 5, DefaultPIndex: 2, Mesh: &Mesh{Nodes: []Point3D{{}, {}, {}, {}},
				Faces: []Face{
					{NodeIndices: [3]uint32{0, 1, 2}, PID: 5, PIndex: [3]uint32{2, 0, 0}},
					{NodeIndices: [3]uint32{0, 1, 4}, PID: 5, PIndex: [3]uint32{2, 2, 2}},
					{NodeIndices: [3]uint32{0, 2, 3}, PID: 5, PIndex: [3]uint32{1, 1, 0}},
					{NodeIndices: [3]uint32{1, 2, 3}, PID: 100},
				}}},
		}}}}, []error{
			fmt.Errorf("go3mf: %s@Resources@Object#0: %v", path, specerr.ErrMissingID),
			fmt.Errorf("go3mf: %s@Resources@Object#0: %v", path, specerr.ErrInvalidObject),
			fmt.Errorf("go3mf: %s@Resources@Object#1: %v", path, specerr.ErrDuplicatedID),
			fmt.Errorf("go3mf: %s@Resources@Object#1: %v", path, &specerr.MissingFieldError{Name: attrPID}),
			fmt.Errorf("go3mf: %s@Resources@Object#1: %v", path, specerr.ErrInvalidObject),
			fmt.Errorf("go3mf: %s@Resources@Object#1@Mesh: %v", path, specerr.ErrInsufficientVertices),
			fmt.Errorf("go3mf: %s@Resources@Object#1@Mesh: %v", path, specerr.ErrInsufficientTriangles),
			fmt.Errorf("go3mf: %s@Resources@Object#1@Component#0: %v", path, specerr.ErrRecursiveComponent),
			fmt.Errorf("go3mf: %s@Resources@Object#1@Component#0: ", path),
			fmt.Errorf("go3mf: %s@Resources@Object#3: %v", path, specerr.ErrComponentsPID),
			fmt.Errorf("go3mf: %s@Resources@Object#3@Component#0: %v", path, specerr.ErrRecursiveComponent),
			fmt.Errorf("go3mf: %s@Resources@Object#3@Component#2: %v", path, &specerr.MissingFieldError{Name: attrObjectID}),
			fmt.Errorf("go3mf: %s@Resources@Object#3@Component#3: %v", path, specerr.ErrMissingResource),
			fmt.Errorf("go3mf: %s@Resources@Object#3@Component#4: %v", path, specerr.ErrMissingResource),
			fmt.Errorf("go3mf: %s@Resources@Object#4: %v", path, specerr.ErrMissingResource),
			fmt.Errorf("go3mf: %s@Resources@Object#4@Mesh: %v", path, specerr.ErrInsufficientVertices),
			fmt.Errorf("go3mf: %s@Resources@Object#4@Mesh: %v", path, specerr.ErrInsufficientTriangles),
			fmt.Errorf("go3mf: %s@Resources@Object#4@Mesh@Face#0: %v", path, specerr.ErrDuplicatedIndices),
			fmt.Errorf("go3mf: %s@Resources@Object#4@Mesh@Face#1: %v", path, specerr.ErrDuplicatedIndices),
			fmt.Errorf("go3mf: %s@Resources@Object#4@Mesh@Face#2: %v", path, specerr.ErrDuplicatedIndices),
			fmt.Errorf("go3mf: %s@Resources@Object#5: %v", path, specerr.ErrIndexOutOfBounds),
			fmt.Errorf("go3mf: %s@Resources@Object#5@Mesh@Face#0: %v", path, specerr.ErrIndexOutOfBounds),
			fmt.Errorf("go3mf: %s@Resources@Object#5@Mesh@Face#1: %v", path, specerr.ErrIndexOutOfBounds),
			fmt.Errorf("go3mf: %s@Resources@Object#5@Mesh@Face#3: %v", path, specerr.ErrMissingResource),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.model.Validate()
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Model.Validate() = %v", diff)
			}
		})
	}
}
