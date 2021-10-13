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
	"github.com/hpinc/go3mf/errors"
	"github.com/hpinc/go3mf/spec"
)

func TestValidate(t *testing.T) {
	spec.Register(fakeSpec.Namespace, new(qmExtension))
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
			fmt.Sprintf("go3mf: Path: /3D/3dmodel.model XPath: /model/relationship[0]: %v", errors.ErrOPCPartName),
			fmt.Sprintf("go3mf: Path: /3D/3dmodel.model XPath: /model/relationship[1]: %v", errors.ErrOPCPartName),
			fmt.Sprintf("go3mf: Path: /3D/3dmodel.model XPath: /model/relationship[3]: %v", errors.ErrOPCPartName),
			fmt.Sprintf("go3mf: Path: /3D/3dmodel.model XPath: /model/relationship[4]: %v", errors.ErrOPCRelTarget),
			fmt.Sprintf("go3mf: Path: /3D/3dmodel.model XPath: /model/relationship[5]: %v", errors.ErrOPCDuplicatedRel),
			fmt.Sprintf("go3mf: Path: /3D/3dmodel.model XPath: /model/relationship[6]: %v", errors.ErrOPCContentType),
			fmt.Sprintf("go3mf: Path: /3D/3dmodel.model XPath: /model/relationship[7]: %v", errors.ErrOPCDuplicatedRel),
			fmt.Sprintf("go3mf: Path: /3D/3dmodel.model XPath: /model/relationship[7]: %v", errors.ErrOPCContentType),
			fmt.Sprintf("go3mf: Path: /3D/3dmodel.model XPath: /model/relationship[7]: %v", errors.ErrOPCDuplicatedTicket),
		}},
		{"namespaces", &Model{Extensions: []Extension{{Namespace: "fake", LocalName: "f", IsRequired: true}}}, []string{
			fmt.Sprintf("go3mf: XPath: /model: %v", errors.ErrRequiredExt),
		}},
		{"metadata", &Model{Extensions: []Extension{{Namespace: "fake", LocalName: "f"}}, Metadata: []Metadata{
			{Name: xml.Name{Space: "fake", Local: "issue"}}, {Name: xml.Name{Space: "f", Local: "issue"}}, {Name: xml.Name{Space: "fake", Local: "issue"}}, {Name: xml.Name{Local: "issue"}}, {},
		}}, []string{
			fmt.Sprintf("go3mf: XPath: /model/metadata[1]: %v", errors.ErrMetadataNamespace),
			fmt.Sprintf("go3mf: XPath: /model/metadata[2]: %v", errors.ErrMetadataDuplicated),
			fmt.Sprintf("go3mf: XPath: /model/metadata[3]: %v", errors.ErrMetadataName),
			fmt.Sprintf("go3mf: XPath: /model/metadata[4]: %v", &errors.MissingFieldError{Name: attrName}),
		}},
		{"build", &Model{Resources: Resources{Assets: []Asset{&BaseMaterials{ID: 1, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}}}}, Objects: []*Object{
			{ID: 2, Type: ObjectTypeOther, Mesh: &Mesh{Vertices: Vertices{Vertex: []Point3D{{}, {}, {}, {}}}, Triangles: Triangles{Triangle: []Triangle{
				{V1: 0, V2: 1, V3: 2}, {V1: 0, V2: 3, V3: 1}, {V1: 0, V2: 2, V3: 3}, {V1: 1, V2: 3, V3: 2},
			}}}}}}, Build: Build{AnyAttr: spec.AnyAttr{&fakeAttr{}}, Items: []*Item{
			{},
			{ObjectID: 2},
			{ObjectID: 100},
			{ObjectID: 1, Metadata: MetadataGroup{Metadata: []Metadata{{Name: xml.Name{Local: "issue"}}}}},
		}}}, []string{
			"go3mf: XPath: /model: Build: fake",
			fmt.Sprintf("go3mf: XPath: /model/build/item[0]: %v", &errors.MissingFieldError{Name: attrObjectID}),
			fmt.Sprintf("go3mf: XPath: /model/build/item[1]: %v", errors.ErrOtherItem),
			fmt.Sprintf("go3mf: XPath: /model/build/item[2]: %v", errors.ErrMissingResource),
			fmt.Sprintf("go3mf: XPath: /model/build/item[3]: %v", errors.ErrMissingResource),
			fmt.Sprintf("go3mf: XPath: /model/build/item[3]/metadata[0]: %v", errors.ErrMetadataName),
		}},
		{"childs", &Model{Childs: map[string]*ChildModel{DefaultModelPath: {}, "/a.model": {
			Relationships: make([]Relationship, 1), Resources: Resources{Objects: []*Object{{}}}}}},
			[]string{
				fmt.Sprintf("go3mf: XPath: /model: %v", errors.ErrOPCDuplicatedModelName),
				fmt.Sprintf("go3mf: Path: /a.model XPath: /model/relationship[0]: %v", errors.ErrOPCPartName),
				fmt.Sprintf("go3mf: Path: /a.model XPath: /model/resources/object[0]: %v", errors.ErrMissingID),
				fmt.Sprintf("go3mf: Path: /a.model XPath: /model/resources/object[0]: %v", errors.ErrInvalidObject),
			}},
		{"assets", &Model{Resources: Resources{Assets: []Asset{
			&BaseMaterials{Materials: []Base{{Color: color.RGBA{}}}},
			&BaseMaterials{ID: 1, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}}},
			&BaseMaterials{ID: 1},
		}}}, []string{
			fmt.Sprintf("go3mf: XPath: /model/resources/basematerials[0]: %v", errors.ErrMissingID),
			fmt.Sprintf("go3mf: XPath: /model/resources/basematerials[0]/base[0]: %v", &errors.MissingFieldError{Name: attrName}),
			fmt.Sprintf("go3mf: XPath: /model/resources/basematerials[0]/base[0]: %v", &errors.MissingFieldError{Name: attrDisplayColor}),
			fmt.Sprintf("go3mf: XPath: /model/resources/basematerials[2]: %v", errors.ErrDuplicatedID),
			fmt.Sprintf("go3mf: XPath: /model/resources/basematerials[2]: %v", errors.ErrEmptyResourceProps),
		}},
		{"objects", &Model{Resources: Resources{Assets: []Asset{
			&BaseMaterials{ID: 1, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}, {Name: "b", Color: color.RGBA{A: 1}}}},
			&BaseMaterials{ID: 5, Materials: []Base{{Name: "a", Color: color.RGBA{A: 1}}, {Name: "b", Color: color.RGBA{A: 1}}}},
		}, Objects: []*Object{
			{},
			{ID: 1, PIndex: 1, Mesh: &Mesh{}, Components: &Components{Component: []*Component{{ObjectID: 1}}}},
			{ID: 2, Mesh: &Mesh{Vertices: Vertices{Vertex: []Point3D{{}, {}, {}, {}}}, Triangles: Triangles{Triangle: []Triangle{
				{V1: 0, V2: 1, V3: 2}, {V1: 0, V2: 3, V3: 1}, {V1: 0, V2: 2, V3: 3}, {V1: 1, V2: 3, V3: 2},
			}}}},
			{ID: 3, PID: 5, Components: &Components{Component: []*Component{
				{ObjectID: 3}, {ObjectID: 2}, {}, {ObjectID: 5}, {ObjectID: 100},
			}}},
			{ID: 4, PID: 100, Mesh: &Mesh{Vertices: Vertices{Vertex: make([]Point3D, 2)}, Triangles: Triangles{Triangle: make([]Triangle, 3)}}},
			{ID: 6, PID: 5, PIndex: 2, Mesh: &Mesh{Vertices: Vertices{Vertex: []Point3D{{}, {}, {}, {}}},
				Triangles: Triangles{Triangle: []Triangle{
					{V1: 0, V2: 1, V3: 2, PID: 5, P1: 2, P2: 0, P3: 0},
					{V1: 0, V2: 1, V3: 4, PID: 5, P1: 2, P2: 2, P3: 2},
					{V1: 0, V2: 2, V3: 3, PID: 5, P1: 1, P2: 1, P3: 0},
					{V1: 1, V2: 2, V3: 3, PID: 100, P1: 0, P2: 0, P3: 0},
				}}}},
		}}}, []string{
			fmt.Sprintf("go3mf: XPath: /model/resources/object[0]: %v", errors.ErrMissingID),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[0]: %v", errors.ErrInvalidObject),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[1]: %v", errors.ErrDuplicatedID),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[1]: %v", &errors.MissingFieldError{Name: attrPID}),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[1]: %v", errors.ErrInvalidObject),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[1]/mesh: %v", errors.ErrInsufficientVertices),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[1]/mesh: %v", errors.ErrInsufficientTriangles),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[1]/components/component[0]: %v", errors.ErrRecursion),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[3]: %v", errors.ErrComponentsPID),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[3]/components/component[0]: %v", errors.ErrRecursion),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[3]/components/component[2]: %v", &errors.MissingFieldError{Name: attrObjectID}),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[3]/components/component[3]: %v", errors.ErrMissingResource),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[3]/components/component[4]: %v", errors.ErrMissingResource),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[4]: %v", errors.ErrMissingResource),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[4]/mesh: %v", errors.ErrInsufficientVertices),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[4]/mesh: %v", errors.ErrInsufficientTriangles),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[4]/mesh/triangle[0]: %v", errors.ErrDuplicatedIndices),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[4]/mesh/triangle[1]: %v", errors.ErrDuplicatedIndices),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[4]/mesh/triangle[2]: %v", errors.ErrDuplicatedIndices),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[5]: %v", errors.ErrIndexOutOfBounds),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[5]/mesh/triangle[0]: %v", errors.ErrIndexOutOfBounds),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[5]/mesh/triangle[1]: %v", errors.ErrIndexOutOfBounds),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[5]/mesh/triangle[3]: %v", errors.ErrMissingResource),
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
		{"few vertices", &Mesh{Vertices: Vertices{Vertex: make([]Point3D, 1)}, Triangles: Triangles{Triangle: make([]Triangle, 3)}}, true},
		{"few triangles", &Mesh{Vertices: Vertices{Vertex: make([]Point3D, 3)}, Triangles: Triangles{Triangle: make([]Triangle, 3)}}, true},
		{"wrong orientation", &Mesh{Vertices: Vertices{Vertex: []Point3D{{}, {}, {}, {}}},
			Triangles: Triangles{Triangle: []Triangle{
				{V1: 0, V2: 1, V3: 2},
				{V1: 0, V2: 3, V3: 1},
				{V1: 0, V2: 2, V3: 3},
				{V1: 1, V2: 2, V3: 3},
			}}}, true},
		{"correct", &Mesh{Vertices: Vertices{Vertex: []Point3D{{}, {}, {}, {}}},
			Triangles: Triangles{Triangle: []Triangle{
				{V1: 0, V2: 1, V3: 2},
				{V1: 0, V2: 3, V3: 1},
				{V1: 0, V2: 2, V3: 3},
				{V1: 1, V2: 3, V3: 2},
			}}}, false},
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
	validMesh := &Mesh{Vertices: Vertices{Vertex: []Point3D{{}, {}, {}, {}}}, Triangles: Triangles{Triangle: []Triangle{
		{V1: 0, V2: 1, V3: 2}, {V1: 0, V2: 3, V3: 1},
		{V1: 0, V2: 2, V3: 3}, {V1: 1, V2: 3, V3: 2},
	}}}
	invalidMesh := &Mesh{Vertices: Vertices{Vertex: []Point3D{{}, {}, {}, {}}}, Triangles: Triangles{Triangle: []Triangle{
		{V1: 0, V2: 1, V3: 2}, {V1: 0, V2: 3, V3: 1},
		{V1: 0, V2: 2, V3: 3}, {V1: 1, V2: 2, V3: 3},
	}}}
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
			fmt.Sprintf("go3mf: Path: /other.model XPath: /model/resources/object[0]/mesh: %v", errors.ErrMeshConsistency),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[0]/mesh: %v", errors.ErrMeshConsistency),
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
