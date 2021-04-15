// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package go3mf

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"sort"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	Register(fakeSpec.Namespace, new(qmExtension))
	tests := []struct {
		name  string
		model *Model
		want  []string
	}{
		{"empty", new(Model), nil},
		{"rels", &Model{Attachments: []Attachment{{Path: "/a.png"}}, Relationships: []Relationship{
			{}, {Path: "/.png"}, {Path: "/a.png"}, {Path: "a.png"}, {Path: "/b.png"}, {Path: "/a.png"},
			{Path: "/a.png", Type: RelTypePrintTicket}, {Path: "/a.png", Type: RelTypePrintTicket},
		}}, []string{
			fmt.Sprintf("/3D/3dmodel.model@Relationship#0: %v", errors.ErrOPCPartName),
			fmt.Sprintf("/3D/3dmodel.model@Relationship#1: %v", errors.ErrOPCPartName),
			fmt.Sprintf("/3D/3dmodel.model@Relationship#3: %v", errors.ErrOPCPartName),
			fmt.Sprintf("/3D/3dmodel.model@Relationship#4: %v", errors.ErrOPCRelTarget),
			fmt.Sprintf("/3D/3dmodel.model@Relationship#5: %v", errors.ErrOPCDuplicatedRel),
			fmt.Sprintf("/3D/3dmodel.model@Relationship#6: %v", errors.ErrOPCContentType),
			fmt.Sprintf("/3D/3dmodel.model@Relationship#7: %v", errors.ErrOPCDuplicatedRel),
			fmt.Sprintf("/3D/3dmodel.model@Relationship#7: %v", errors.ErrOPCContentType),
			fmt.Sprintf("/3D/3dmodel.model@Relationship#7: %v", errors.ErrOPCDuplicatedTicket),
		}},
		{"namespaces", &Model{Extensions: []Extension{{Namespace: "fake", LocalName: "f", IsRequired: true}}}, []string{
			errors.ErrRequiredExt.Error(),
		}},
		{"metadata", &Model{Extensions: []Extension{{Namespace: "fake", LocalName: "f"}}, Metadata: []Metadata{
			{Name: xml.Name{Space: "fake", Local: "issue"}}, {Name: xml.Name{Space: "f", Local: "issue"}}, {Name: xml.Name{Space: "fake", Local: "issue"}}, {Name: xml.Name{Local: "issue"}}, {},
		}}, []string{
			fmt.Sprintf("Metadata#1: %v", errors.ErrMetadataNamespace),
			fmt.Sprintf("Metadata#2: %v", errors.ErrMetadataDuplicated),
			fmt.Sprintf("Metadata#3: %v", errors.ErrMetadataName),
			fmt.Sprintf("Metadata#4: %v", &errors.MissingFieldError{Name: attrName}),
		}},
		{"build", &Model{Resources: Resources{Assets: []Asset{&BaseMaterials{ID: 1, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}}}}, Objects: []*Object{
			{ID: 2, Type: ObjectTypeOther, Mesh: &Mesh{Vertices: []Point3D{{}, {}, {}, {}}, Triangles: []Triangle{
				NewTriangle(0, 1, 2), NewTriangle(0, 3, 1), NewTriangle(0, 2, 3), NewTriangle(1, 3, 2),
			}}}}}, Build: Build{AnyAttr: AnyAttr{&fakeAttr{}}, Items: []*Item{
			{},
			{ObjectID: 2},
			{ObjectID: 100},
			{ObjectID: 1, Metadata: []Metadata{{Name: xml.Name{Local: "issue"}}}},
		}}}, []string{
			fmt.Sprintf("Build: fake"),
			fmt.Sprintf("Build@Item#0: %v", &errors.MissingFieldError{Name: attrObjectID}),
			fmt.Sprintf("Build@Item#1: %v", errors.ErrOtherItem),
			fmt.Sprintf("Build@Item#2: %v", errors.ErrMissingResource),
			fmt.Sprintf("Build@Item#3: %v", errors.ErrMissingResource),
			fmt.Sprintf("Build@Item#3@Metadata#0: %v", errors.ErrMetadataName),
		}},
		{"childs", &Model{Childs: map[string]*ChildModel{DefaultModelPath: {}, "/a.model": {
			Relationships: make([]Relationship, 1), Resources: Resources{Objects: []*Object{{}}}}}},
			[]string{
				errors.ErrOPCDuplicatedModelName.Error(),
				fmt.Sprintf("/a.model@Relationship#0: %v", errors.ErrOPCPartName),
				fmt.Sprintf("/a.model@Resources@Object#0: %v", errors.ErrMissingID),
				fmt.Sprintf("/a.model@Resources@Object#0: %v", errors.ErrInvalidObject),
			}},
		{"assets", &Model{Resources: Resources{Assets: []Asset{
			&BaseMaterials{Materials: []Base{{Color: color.RGBA{}}}},
			&BaseMaterials{ID: 1, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}}},
			&BaseMaterials{ID: 1},
		}}}, []string{
			fmt.Sprintf("Resources@BaseMaterials#0: %v", errors.ErrMissingID),
			fmt.Sprintf("Resources@BaseMaterials#0@Base#0: %v", &errors.MissingFieldError{Name: attrName}),
			fmt.Sprintf("Resources@BaseMaterials#0@Base#0: %v", &errors.MissingFieldError{Name: attrDisplayColor}),
			fmt.Sprintf("Resources@BaseMaterials#2: %v", errors.ErrDuplicatedID),
			fmt.Sprintf("Resources@BaseMaterials#2: %v", errors.ErrEmptyResourceProps),
		}},
		{"objects", &Model{Resources: Resources{Assets: []Asset{
			&BaseMaterials{ID: 1, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}, {Name: "b", Color: color.RGBA{A: 1}}}},
			&BaseMaterials{ID: 5, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}, {Name: "b", Color: color.RGBA{A: 1}}}},
		}, Objects: []*Object{
			{},
			{ID: 1, PIndex: 1, Mesh: &Mesh{}, Components: []*Component{{ObjectID: 1}}},
			{ID: 2, Mesh: &Mesh{Vertices: []Point3D{{}, {}, {}, {}}, Triangles: []Triangle{
				NewTriangle(0, 1, 2), NewTriangle(0, 3, 1), NewTriangle(0, 2, 3), NewTriangle(1, 3, 2),
			}}},
			{ID: 3, PID: 5, Components: []*Component{
				{ObjectID: 3}, {ObjectID: 2}, {}, {ObjectID: 5}, {ObjectID: 100},
			}},
			{ID: 4, PID: 100, Mesh: &Mesh{Vertices: make([]Point3D, 2), Triangles: make([]Triangle, 3)}},
			{ID: 6, PID: 5, PIndex: 2, Mesh: &Mesh{Vertices: []Point3D{{}, {}, {}, {}},
				Triangles: []Triangle{
					NewTrianglePID(0, 1, 2, 5, 2, 0, 0),
					NewTrianglePID(0, 1, 4, 5, 2, 2, 2),
					NewTrianglePID(0, 2, 3, 5, 1, 1, 0),
					NewTrianglePID(1, 2, 3, 100, 0, 0, 0),
				}}},
		}}}, []string{
			fmt.Sprintf("Resources@Object#0: %v", errors.ErrMissingID),
			fmt.Sprintf("Resources@Object#0: %v", errors.ErrInvalidObject),
			fmt.Sprintf("Resources@Object#1: %v", errors.ErrDuplicatedID),
			fmt.Sprintf("Resources@Object#1: %v", &errors.MissingFieldError{Name: attrPID}),
			fmt.Sprintf("Resources@Object#1: %v", errors.ErrInvalidObject),
			fmt.Sprintf("Resources@Object#1@Mesh: %v", errors.ErrInsufficientVertices),
			fmt.Sprintf("Resources@Object#1@Mesh: %v", errors.ErrInsufficientTriangles),
			fmt.Sprintf("Resources@Object#1@Component#0: %v", errors.ErrRecursion),
			fmt.Sprintf("Resources@Object#3: %v", errors.ErrComponentsPID),
			fmt.Sprintf("Resources@Object#3@Component#0: %v", errors.ErrRecursion),
			fmt.Sprintf("Resources@Object#3@Component#2: %v", &errors.MissingFieldError{Name: attrObjectID}),
			fmt.Sprintf("Resources@Object#3@Component#3: %v", errors.ErrMissingResource),
			fmt.Sprintf("Resources@Object#3@Component#4: %v", errors.ErrMissingResource),
			fmt.Sprintf("Resources@Object#4: %v", errors.ErrMissingResource),
			fmt.Sprintf("Resources@Object#4@Mesh: %v", errors.ErrInsufficientVertices),
			fmt.Sprintf("Resources@Object#4@Mesh: %v", errors.ErrInsufficientTriangles),
			fmt.Sprintf("Resources@Object#4@Mesh@Triangle#0: %v", errors.ErrDuplicatedIndices),
			fmt.Sprintf("Resources@Object#4@Mesh@Triangle#1: %v", errors.ErrDuplicatedIndices),
			fmt.Sprintf("Resources@Object#4@Mesh@Triangle#2: %v", errors.ErrDuplicatedIndices),
			fmt.Sprintf("Resources@Object#5: %v", errors.ErrIndexOutOfBounds),
			fmt.Sprintf("Resources@Object#5@Mesh@Triangle#0: %v", errors.ErrIndexOutOfBounds),
			fmt.Sprintf("Resources@Object#5@Mesh@Triangle#1: %v", errors.ErrIndexOutOfBounds),
			fmt.Sprintf("Resources@Object#5@Mesh@Triangle#3: %v", errors.ErrMissingResource),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.model.Extensions = append(tt.model.Extensions, fakeSpec)
			got := tt.model.Validate()
			if tt.want == nil {
				if got != nil {
					t.Errorf("Model.Validate() err = %v", got)
				}
				return
			}
			if got == nil {
				t.Errorf("Model.Validate() err nil = want %v", tt.want)
				return
			}
			var errs []string
			for _, err := range got.(*errors.List).Errors {
				errs = append(errs, err.Error())
			}
			if diff := deep.Equal(errs, tt.want); diff != nil {
				t.Errorf("Model.Validate() = %v", diff)
			}
		})
	}
}

func TestObject_ValidateMesh(t *testing.T) {
	tests := []struct {
		name    string
		r       *Mesh
		wantErr bool
	}{
		{"few vertices", &Mesh{Vertices: make([]Point3D, 1), Triangles: make([]Triangle, 3)}, true},
		{"few triangles", &Mesh{Vertices: make([]Point3D, 3), Triangles: make([]Triangle, 3)}, true},
		{"wrong orientation", &Mesh{Vertices: []Point3D{{}, {}, {}, {}},
			Triangles: []Triangle{
				NewTriangle(0, 1, 2),
				NewTriangle(0, 3, 1),
				NewTriangle(0, 2, 3),
				NewTriangle(1, 2, 3),
			}}, true},
		{"correct", &Mesh{Vertices: []Point3D{{}, {}, {}, {}},
			Triangles: []Triangle{
				NewTriangle(0, 1, 2),
				NewTriangle(0, 3, 1),
				NewTriangle(0, 2, 3),
				NewTriangle(1, 3, 2),
			}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.ValidateCoherency(); (err != nil) != tt.wantErr {
				t.Errorf("Object.ValidateCoherency() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModel_ValidateCoherency(t *testing.T) {
	validMesh := &Mesh{Vertices: []Point3D{{}, {}, {}, {}}, Triangles: []Triangle{
		NewTriangle(0, 1, 2), NewTriangle(0, 3, 1),
		NewTriangle(0, 2, 3), NewTriangle(1, 3, 2),
	}}
	invalidMesh := &Mesh{Vertices: []Point3D{{}, {}, {}, {}}, Triangles: []Triangle{
		NewTriangle(0, 1, 2), NewTriangle(0, 3, 1),
		NewTriangle(0, 2, 3), NewTriangle(1, 2, 3),
	}}
	tests := []struct {
		name string
		m    *Model
		want []string
	}{
		{"empty", new(Model), nil},
		{"valid", &Model{Resources: Resources{Objects: []*Object{
			{Mesh: validMesh},
		}}, Childs: map[string]*ChildModel{"/other.model": {Resources: Resources{Objects: []*Object{
			{Mesh: validMesh},
		}}}}}, nil},
		{"invalid", &Model{Resources: Resources{Objects: []*Object{
			{Mesh: invalidMesh},
		}}, Childs: map[string]*ChildModel{"/other.model": {Resources: Resources{Objects: []*Object{
			{Mesh: invalidMesh},
		}}}}}, []string{
			fmt.Sprintf("/other.model@Resources@Object#0@Mesh: %v", errors.ErrMeshConsistency),
			fmt.Sprintf("Resources@Object#0@Mesh: %v", errors.ErrMeshConsistency),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.ValidateCoherency()
			if tt.want == nil {
				if got != nil {
					t.Errorf("Model.ValidateCoherency() err = %v", got)
				}
				return
			}
			if got == nil {
				t.Errorf("Model.ValidateCoherency() err nil = want %v", tt.want)
				return
			}
			var errs []string
			for _, err := range got.(*errors.List).Errors {
				errs = append(errs, err.Error())
			}
			sort.Strings(errs)
			if diff := deep.Equal(errs, tt.want); diff != nil {
				t.Errorf("Model.ValidateCoherency() = %v", diff)
			}
		})
	}
}
