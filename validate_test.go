package go3mf

import (
	"encoding/xml"
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
		{"empty", args{new(Model)}, nil},
		{"rels", args{&Model{Attachments: []Attachment{{Path: "/a.png"}}, Relationships: []Relationship{
			{}, {Path: "/.png"}, {Path: "/a.png"}, {Path: "a.png"}, {Path: "/b.png"}, {Path: "/a.png"},
			{Path: "/a.png", Type: RelTypePrintTicket}, {Path: "/a.png", Type: RelTypePrintTicket},
		}}}, []error{
			&specerr.RelationshipError{Path: path, Index: 0, Err: specerr.ErrOPCPartName},
			&specerr.RelationshipError{Path: path, Index: 1, Err: specerr.ErrOPCPartName},
			&specerr.RelationshipError{Path: path, Index: 3, Err: specerr.ErrOPCPartName},
			&specerr.RelationshipError{Path: path, Index: 4, Err: specerr.ErrOPCRelTarget},
			&specerr.RelationshipError{Path: path, Index: 5, Err: specerr.ErrOPCDuplicatedRel},
			&specerr.RelationshipError{Path: path, Index: 6, Err: specerr.ErrOPCContentType},
			&specerr.RelationshipError{Path: path, Index: 7, Err: specerr.ErrOPCDuplicatedRel},
			&specerr.RelationshipError{Path: path, Index: 7, Err: specerr.ErrOPCContentType},
			&specerr.RelationshipError{Path: path, Index: 7, Err: specerr.ErrOPCDuplicatedTicket},
		}},
		{"namespaces", args{&Model{RequiredExtensions: []string{"fake", "other"}, Namespaces: []xml.Name{{Space: "fake", Local: "f"}}}}, []error{
			specerr.ErrRequiredExt,
		}},
		{"metadata", args{&Model{Namespaces: []xml.Name{{Space: "fake", Local: "f"}}, Metadata: []Metadata{
			{Name: xml.Name{Space: "fake", Local: "issue"}}, {Name: xml.Name{Space: "f", Local: "issue"}}, {Name: xml.Name{Space: "fake", Local: "issue"}}, {Name: xml.Name{Local: "issue"}}, {},
		}}}, []error{
			&specerr.IndexedError{Name: attrMetadata, Index: 1, Err: specerr.ErrMetadataNamespace},
			&specerr.IndexedError{Name: attrMetadata, Index: 2, Err: specerr.ErrMetadataDuplicated},
			&specerr.IndexedError{Name: attrMetadata, Index: 3, Err: specerr.ErrMetadataName},
			&specerr.IndexedError{Name: attrMetadata, Index: 4, Err: &specerr.MissingFieldError{Name: attrName}},
		}},
		{"build", args{&Model{Resources: Resources{Assets: []Asset{&BaseMaterialsResource{ID: 1, Materials: []BaseMaterial{{Name: "a", Color: color.RGBA{A: 1}}}}}, Objects: []*Object{
			{ID: 2, ObjectType: ObjectTypeOther, Mesh: &Mesh{Nodes: []Point3D{{}, {}, {}, {}}, Faces: []Face{
				{NodeIndices: [3]uint32{0, 1, 2}}, {NodeIndices: [3]uint32{0, 3, 1}}, {NodeIndices: [3]uint32{0, 2, 3}}, {NodeIndices: [3]uint32{1, 3, 2}},
			}}}}}, Build: Build{Items: []*Item{
			{},
			{ObjectID: 2},
			{ObjectID: 100},
			{ObjectID: 1, Metadata: []Metadata{{Name: xml.Name{Local: "issue"}}}},
		}}}}, []error{
			&specerr.ItemError{Index: 0, Err: &specerr.MissingFieldError{Name: attrObjectID}},
			&specerr.ItemError{Index: 1, Err: specerr.ErrOtherItem},
			&specerr.ItemError{Index: 2, Err: specerr.ErrMissingResource},
			&specerr.ItemError{Index: 3, Err: specerr.ErrMissingResource},
			&specerr.ItemError{Index: 3, Err: &specerr.IndexedError{Name: attrMetadata, Index: 0, Err: specerr.ErrMetadataName}},
		}},
		{"childs", args{&Model{Childs: map[string]*ChildModel{path: &ChildModel{}, "/a.model": &ChildModel{
			Relationships: make([]Relationship, 1), Resources: Resources{Objects: []*Object{{}}}}}}},
			[]error{
				specerr.ErrOPCDuplicatedModelName,
				&specerr.RelationshipError{Path: "/a.model", Index: 0, Err: specerr.ErrOPCPartName},
				&specerr.ObjectError{Path: "/a.model", Index: 0, Err: specerr.ErrMissingID},
				&specerr.ObjectError{Path: "/a.model", Index: 0, Err: specerr.ErrInvalidObject},
			}},
		{"assets", args{&Model{Resources: Resources{Assets: []Asset{
			&BaseMaterialsResource{Materials: []BaseMaterial{{Color: color.RGBA{}}}},
			&BaseMaterialsResource{ID: 1, Materials: []BaseMaterial{{Name: "a", Color: color.RGBA{A: 1}}}},
			&BaseMaterialsResource{ID: 1},
		}}}}, []error{
			&specerr.AssetError{Path: path, Index: 0, Name: "BaseMaterialsResource", Err: specerr.ErrMissingID},
			&specerr.AssetError{Path: path, Index: 0, Name: "BaseMaterialsResource", Err: &specerr.IndexedError{Name: attrBase, Index: 0, Err: &specerr.MissingFieldError{Name: attrName}}},
			&specerr.AssetError{Path: path, Index: 0, Name: "BaseMaterialsResource", Err: &specerr.IndexedError{Name: attrBase, Index: 0, Err: &specerr.MissingFieldError{Name: attrDisplayColor}}},
			&specerr.AssetError{Path: path, Index: 2, Name: "BaseMaterialsResource", Err: specerr.ErrDuplicatedID},
			&specerr.AssetError{Path: path, Index: 2, Name: "BaseMaterialsResource", Err: specerr.ErrEmptyResourceProps},
		}},
		{"objects", args{&Model{Resources: Resources{Assets: []Asset{
			&BaseMaterialsResource{ID: 1, Materials: []BaseMaterial{{Name: "a", Color: color.RGBA{A: 1}}, {Name: "b", Color: color.RGBA{A: 1}}}},
			&BaseMaterialsResource{ID: 5, Materials: []BaseMaterial{{Name: "a", Color: color.RGBA{A: 1}}, {Name: "b", Color: color.RGBA{A: 1}}}},
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
			{ID: 6, DefaultPID: 5, DefaultPIndex: 3, Mesh: &Mesh{Nodes: []Point3D{{}, {}, {}, {}},
				Faces: []Face{
					{NodeIndices: [3]uint32{0, 1, 2}, PID: 5, PIndex: [3]uint32{4, 0, 0}},
					{NodeIndices: [3]uint32{0, 1, 4}},
					{NodeIndices: [3]uint32{0, 2, 3}, PID: 5, PIndex: [3]uint32{1, 2, 0}},
					{NodeIndices: [3]uint32{1, 2, 3}, PID: 100},
				}}},
		}}}}, []error{
			&specerr.ObjectError{Path: path, Index: 0, Err: specerr.ErrMissingID},
			&specerr.ObjectError{Path: path, Index: 0, Err: specerr.ErrInvalidObject},
			&specerr.ObjectError{Path: path, Index: 1, Err: specerr.ErrDuplicatedID},
			&specerr.ObjectError{Path: path, Index: 1, Err: &specerr.MissingFieldError{Name: attrPID}},
			&specerr.ObjectError{Path: path, Index: 1, Err: specerr.ErrInvalidObject},
			&specerr.ObjectError{Path: path, Index: 1, Err: specerr.ErrInsufficientVertices},
			&specerr.ObjectError{Path: path, Index: 1, Err: specerr.ErrInsufficientTriangles},
			&specerr.ObjectError{Path: path, Index: 1, Err: &specerr.IndexedError{Name: attrComponent, Index: 0, Err: specerr.ErrRecursiveComponent}},
			&specerr.ObjectError{Path: path, Index: 3, Err: specerr.ErrComponentsPID},
			&specerr.ObjectError{Path: path, Index: 3, Err: &specerr.IndexedError{Name: attrComponent, Index: 0, Err: specerr.ErrRecursiveComponent}},
			&specerr.ObjectError{Path: path, Index: 3, Err: &specerr.IndexedError{Name: attrComponent, Index: 2, Err: &specerr.MissingFieldError{Name: attrObjectID}}},
			&specerr.ObjectError{Path: path, Index: 3, Err: &specerr.IndexedError{Name: attrComponent, Index: 3, Err: specerr.ErrMissingResource}},
			&specerr.ObjectError{Path: path, Index: 3, Err: &specerr.IndexedError{Name: attrComponent, Index: 4, Err: specerr.ErrMissingResource}},
			&specerr.ObjectError{Path: path, Index: 4, Err: specerr.ErrMissingResource},
			&specerr.ObjectError{Path: path, Index: 4, Err: specerr.ErrInsufficientVertices},
			&specerr.ObjectError{Path: path, Index: 4, Err: specerr.ErrInsufficientTriangles},
			&specerr.ObjectError{Path: path, Index: 4, Err: &specerr.IndexedError{Name: attrTriangle, Index: 0, Err: specerr.ErrDuplicatedIndices}},
			&specerr.ObjectError{Path: path, Index: 4, Err: &specerr.IndexedError{Name: attrTriangle, Index: 1, Err: specerr.ErrDuplicatedIndices}},
			&specerr.ObjectError{Path: path, Index: 4, Err: &specerr.IndexedError{Name: attrTriangle, Index: 2, Err: specerr.ErrDuplicatedIndices}},
			&specerr.ObjectError{Path: path, Index: 5, Err: specerr.ErrIndexOutOfBounds},
			&specerr.ObjectError{Path: path, Index: 5, Err: &specerr.IndexedError{Name: attrTriangle, Index: 0, Err: specerr.ErrIndexOutOfBounds}},
			&specerr.ObjectError{Path: path, Index: 5, Err: &specerr.IndexedError{Name: attrTriangle, Index: 1, Err: specerr.ErrIndexOutOfBounds}},
			&specerr.ObjectError{Path: path, Index: 5, Err: &specerr.IndexedError{Name: attrTriangle, Index: 2, Err: specerr.ErrBaseMaterialGradient}},
			&specerr.ObjectError{Path: path, Index: 5, Err: &specerr.IndexedError{Name: attrTriangle, Index: 3, Err: specerr.ErrMissingResource}},
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
