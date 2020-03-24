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
	tests := []struct {
		name  string
		model *Model
		want  []error
	}{
		{"empty", new(Model), []error{}},
		{"rels", &Model{Attachments: []Attachment{{Path: "/a.png"}}, Relationships: []Relationship{
			{}, {Path: "/.png"}, {Path: "/a.png"}, {Path: "a.png"}, {Path: "/b.png"}, {Path: "/a.png"},
			{Path: "/a.png", Type: RelTypePrintTicket}, {Path: "/a.png", Type: RelTypePrintTicket},
		}}, []error{
			fmt.Errorf("Relationship#0: %v", specerr.ErrOPCPartName),
			fmt.Errorf("Relationship#1: %v", specerr.ErrOPCPartName),
			fmt.Errorf("Relationship#3: %v", specerr.ErrOPCPartName),
			fmt.Errorf("Relationship#4: %v", specerr.ErrOPCRelTarget),
			fmt.Errorf("Relationship#5: %v", specerr.ErrOPCDuplicatedRel),
			fmt.Errorf("Relationship#6: %v", specerr.ErrOPCContentType),
			fmt.Errorf("Relationship#7: %v", specerr.ErrOPCDuplicatedRel),
			fmt.Errorf("Relationship#7: %v", specerr.ErrOPCContentType),
			fmt.Errorf("Relationship#7: %v", specerr.ErrOPCDuplicatedTicket),
		}},
		{"namespaces", &Model{RequiredExtensions: []string{"fake", "other"}, Namespaces: []xml.Name{{Space: "fake", Local: "f"}}}, []error{
			specerr.ErrRequiredExt,
		}},
		{"metadata", &Model{Namespaces: []xml.Name{{Space: "fake", Local: "f"}}, Metadata: []Metadata{
			{Name: xml.Name{Space: "fake", Local: "issue"}}, {Name: xml.Name{Space: "f", Local: "issue"}}, {Name: xml.Name{Space: "fake", Local: "issue"}}, {Name: xml.Name{Local: "issue"}}, {},
		}}, []error{
			fmt.Errorf("Metadata#1: %v", specerr.ErrMetadataNamespace),
			fmt.Errorf("Metadata#2: %v", specerr.ErrMetadataDuplicated),
			fmt.Errorf("Metadata#3: %v", specerr.ErrMetadataName),
			fmt.Errorf("Metadata#4: %v", &specerr.MissingFieldError{Name: attrName}),
		}},
		{"build", &Model{Resources: Resources{Assets: []Asset{&BaseMaterials{ID: 1, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}}}}, Objects: []*Object{
			{ID: 2, ObjectType: ObjectTypeOther, Mesh: &Mesh{Nodes: []Point3D{{}, {}, {}, {}}, Faces: []Face{
				{NodeIndices: [3]uint32{0, 1, 2}}, {NodeIndices: [3]uint32{0, 3, 1}}, {NodeIndices: [3]uint32{0, 2, 3}}, {NodeIndices: [3]uint32{1, 3, 2}},
			}}}}}, Build: Build{ExtensionAttr: ExtensionAttr{&fakeAttr{}}, Items: []*Item{
			{},
			{ObjectID: 2},
			{ObjectID: 100},
			{ObjectID: 1, Metadata: []Metadata{{Name: xml.Name{Local: "issue"}}}},
		}}}, []error{
			fmt.Errorf("Build: fake"),
			fmt.Errorf("Build@Item#0: %v", &specerr.MissingFieldError{Name: attrObjectID}),
			fmt.Errorf("Build@Item#1: %v", specerr.ErrOtherItem),
			fmt.Errorf("Build@Item#2: %v", specerr.ErrMissingResource),
			fmt.Errorf("Build@Item#3: %v", specerr.ErrMissingResource),
			fmt.Errorf("Build@Item#3@Metadata#0: %v", specerr.ErrMetadataName),
		}},
		{"childs", &Model{Childs: map[string]*ChildModel{DefaultModelPath: {}, "/a.model": {
			Relationships: make([]Relationship, 1), Resources: Resources{Objects: []*Object{{}}}}}},
			[]error{
				specerr.ErrOPCDuplicatedModelName,
				fmt.Errorf("/a.model@Relationship#0: %v", specerr.ErrOPCPartName),
				fmt.Errorf("/a.model@Resources@Object#0: %v", specerr.ErrMissingID),
				fmt.Errorf("/a.model@Resources@Object#0: %v", specerr.ErrInvalidObject),
			}},
		{"assets", &Model{Resources: Resources{Assets: []Asset{
			&BaseMaterials{Materials: []Base{{Color: color.RGBA{}}}},
			&BaseMaterials{ID: 1, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}}},
			&BaseMaterials{ID: 1},
		}}}, []error{
			fmt.Errorf("Resources@BaseMaterials#0: %v", specerr.ErrMissingID),
			fmt.Errorf("Resources@BaseMaterials#0@Base#0: %v", &specerr.MissingFieldError{Name: attrName}),
			fmt.Errorf("Resources@BaseMaterials#0@Base#0: %v", &specerr.MissingFieldError{Name: attrDisplayColor}),
			fmt.Errorf("Resources@BaseMaterials#2: %v", specerr.ErrDuplicatedID),
			fmt.Errorf("Resources@BaseMaterials#2: %v", specerr.ErrEmptyResourceProps),
		}},
		{"objects", &Model{Resources: Resources{Assets: []Asset{
			&BaseMaterials{ID: 1, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}, {Name: "b", Color: color.RGBA{A: 1}}}},
			&BaseMaterials{ID: 5, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}, {Name: "b", Color: color.RGBA{A: 1}}}},
		}, Objects: []*Object{
			{},
			{ID: 1, DefaultPIndex: 1, Mesh: &Mesh{}, Components: []*Component{{ObjectID: 1}}},
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
		}}}, []error{
			fmt.Errorf("Resources@Object#0: %v", specerr.ErrMissingID),
			fmt.Errorf("Resources@Object#0: %v", specerr.ErrInvalidObject),
			fmt.Errorf("Resources@Object#1: %v", specerr.ErrDuplicatedID),
			fmt.Errorf("Resources@Object#1: %v", &specerr.MissingFieldError{Name: attrPID}),
			fmt.Errorf("Resources@Object#1: %v", specerr.ErrInvalidObject),
			fmt.Errorf("Resources@Object#1@Mesh: %v", specerr.ErrInsufficientVertices),
			fmt.Errorf("Resources@Object#1@Mesh: %v", specerr.ErrInsufficientTriangles),
			fmt.Errorf("Resources@Object#1@Component#0: %v", specerr.ErrRecursiveComponent),
			fmt.Errorf("Resources@Object#3: %v", specerr.ErrComponentsPID),
			fmt.Errorf("Resources@Object#3@Component#0: %v", specerr.ErrRecursiveComponent),
			fmt.Errorf("Resources@Object#3@Component#2: %v", &specerr.MissingFieldError{Name: attrObjectID}),
			fmt.Errorf("Resources@Object#3@Component#3: %v", specerr.ErrMissingResource),
			fmt.Errorf("Resources@Object#3@Component#4: %v", specerr.ErrMissingResource),
			fmt.Errorf("Resources@Object#4: %v", specerr.ErrMissingResource),
			fmt.Errorf("Resources@Object#4@Mesh: %v", specerr.ErrInsufficientVertices),
			fmt.Errorf("Resources@Object#4@Mesh: %v", specerr.ErrInsufficientTriangles),
			fmt.Errorf("Resources@Object#4@Mesh@Face#0: %v", specerr.ErrDuplicatedIndices),
			fmt.Errorf("Resources@Object#4@Mesh@Face#1: %v", specerr.ErrDuplicatedIndices),
			fmt.Errorf("Resources@Object#4@Mesh@Face#2: %v", specerr.ErrDuplicatedIndices),
			fmt.Errorf("Resources@Object#5: %v", specerr.ErrIndexOutOfBounds),
			fmt.Errorf("Resources@Object#5@Mesh@Face#0: %v", specerr.ErrIndexOutOfBounds),
			fmt.Errorf("Resources@Object#5@Mesh@Face#1: %v", specerr.ErrIndexOutOfBounds),
			fmt.Errorf("Resources@Object#5@Mesh@Face#3: %v", specerr.ErrMissingResource),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.model.ExtensionSpecs = append(tt.model.ExtensionSpecs, &fakeSpec{})
			got := tt.model.Validate()
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Model.Validate() = %v", diff)
			}
		})
	}
}
