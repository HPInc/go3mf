// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package production

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/hpinc/go3mf"
	"github.com/hpinc/go3mf/errors"
	"github.com/hpinc/go3mf/spec"
)

func TestValidate(t *testing.T) {
	validMesh := &go3mf.Object{ID: 1, Mesh: &go3mf.Mesh{Vertices: go3mf.Vertices{Vertex: []go3mf.Point3D{{}, {}, {}, {}}}, Triangles: go3mf.Triangles{Triangle: []go3mf.Triangle{
		{V1: 0, V2: 1, V3: 2}, {V1: 0, V2: 3, V3: 1}, {V1: 0, V2: 2, V3: 3}, {V1: 1, V2: 3, V3: 2},
	}}}}
	tests := []struct {
		name  string
		model *go3mf.Model
		want  []string
	}{
		{"buildNoUUID", &go3mf.Model{Build: go3mf.Build{}}, []string{
			fmt.Sprintf("go3mf: XPath: /model/build: %v", &errors.MissingFieldError{Name: attrProdUUID}),
		}},
		{"buildEmptyUUID", &go3mf.Model{Build: go3mf.Build{
			AnyAttr: spec.AnyAttr{&BuildAttr{}}}}, []string{
			fmt.Sprintf("go3mf: XPath: /model/build: %v", ErrUUID),
		}},
		{"buildNonValidUUID", &go3mf.Model{Build: go3mf.Build{
			AnyAttr: spec.AnyAttr{&BuildAttr{"a-b-c-d"}}}}, []string{
			fmt.Sprintf("go3mf: XPath: /model/build: %v", ErrUUID),
		}},
		{"extReq", &go3mf.Model{
			Childs: map[string]*go3mf.ChildModel{"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}}},
			Resources: go3mf.Resources{Objects: []*go3mf.Object{
				{ID: 5, AnyAttr: spec.AnyAttr{&ObjectAttr{UUID: "f47ac10b-58cc-0372-8567-0e02b2c3d481"}}, Components: &go3mf.Components{Component: []*go3mf.Component{
					{ObjectID: 1, AnyAttr: spec.AnyAttr{
						&ComponentAttr{Path: "/other.model", UUID: "f47ac10b-58cc-0372-8567-0e02b2c3d480"},
					}}}}}}}, Build: go3mf.Build{
				AnyAttr: spec.AnyAttr{&BuildAttr{UUID: "f47ac10b-58cc-0372-8567-0e02b2c3d479"}}, Items: []*go3mf.Item{
					{ObjectID: 1, AnyAttr: spec.AnyAttr{&ItemAttr{UUID: "f47ac10b-58cc-0372-8567-0e02b2c3d478", Path: "/other.model"}}},
				}}}, []string{
			fmt.Sprintf("go3mf: Path: /other.model XPath: /model/resources/object[0]: %v", &errors.MissingFieldError{Name: attrProdUUID}),
		}},
		{"items", &go3mf.Model{Build: go3mf.Build{
			AnyAttr: spec.AnyAttr{&BuildAttr{UUID: "f47ac10b-58cc-0372-8567-0e02b2c3d479"}}, Items: []*go3mf.Item{
				{ObjectID: 1, AnyAttr: spec.AnyAttr{&ItemAttr{UUID: "f47ac10b-58cc-0372-8567-0e02b2c3d478", Path: "/other.model"}}},
				{ObjectID: 1},
				{ObjectID: 1, AnyAttr: spec.AnyAttr{&ItemAttr{}}},
				{ObjectID: 1, AnyAttr: spec.AnyAttr{&ItemAttr{UUID: "a-b-c-d"}}},
			}},
			Childs:    map[string]*go3mf.ChildModel{"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}}},
			Resources: go3mf.Resources{Objects: []*go3mf.Object{{ID: 1, Mesh: validMesh.Mesh}}}}, []string{
			fmt.Sprintf("go3mf: XPath: /model/build/item[1]: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Sprintf("go3mf: XPath: /model/build/item[2]: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Sprintf("go3mf: XPath: /model/build/item[3]: %v", ErrUUID),
			fmt.Sprintf("go3mf: Path: /other.model XPath: /model/resources/object[0]: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[0]: %v", &errors.MissingFieldError{Name: attrProdUUID}),
		}},
		{"components", &go3mf.Model{Resources: go3mf.Resources{
			Objects: []*go3mf.Object{
				{ID: 2, Mesh: validMesh.Mesh, AnyAttr: spec.AnyAttr{&ObjectAttr{UUID: "a-b-c-d"}}},
				{ID: 3, AnyAttr: spec.AnyAttr{&ObjectAttr{UUID: "f47ac10b-58cc-0372-8567-0e02b2c3d483"}}, Components: &go3mf.Components{Component: []*go3mf.Component{
					{ObjectID: 2, AnyAttr: spec.AnyAttr{&ComponentAttr{}}},
					{ObjectID: 2, AnyAttr: spec.AnyAttr{&ComponentAttr{UUID: "a-b-c-d"}}},
					{ObjectID: 2},
				}}},
			},
		}, Build: go3mf.Build{AnyAttr: spec.AnyAttr{&BuildAttr{UUID: "f47ac10b-58cc-0372-8567-0e02b2c3d479"}}}}, []string{
			fmt.Sprintf("go3mf: XPath: /model/resources/object[0]: %v", ErrUUID),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[1]/components/component[0]: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[1]/components/component[1]: %v", ErrUUID),
			fmt.Sprintf("go3mf: XPath: /model/resources/object[1]/components/component[2]: %v", &errors.MissingFieldError{Name: attrProdUUID}),
		}},
		{"child", &go3mf.Model{Build: go3mf.Build{AnyAttr: spec.AnyAttr{&BuildAttr{UUID: "f47ac10b-58cc-0372-8567-0e02b2c3d479"}}},
			Childs: map[string]*go3mf.ChildModel{
				"/b.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}},
				"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{
					{ID: 2, Components: &go3mf.Components{Component: []*go3mf.Component{
						{ObjectID: 1, AnyAttr: spec.AnyAttr{&ComponentAttr{Path: "/b.model"}}},
					}}},
				}}}}}, []string{
			fmt.Sprintf("go3mf: Path: /b.model XPath: /model/resources/object[0]: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Sprintf("go3mf: Path: /other.model XPath: /model/resources/object[0]: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Sprintf("go3mf: Path: /other.model XPath: /model/resources/object[0]/components/component[0]: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Sprintf("go3mf: Path: /other.model XPath: /model/resources/object[0]/components/component[0]: %v", ErrProdRefInNonRoot),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.model.Extensions = []go3mf.Extension{DefaultExtension}
			err := tt.model.Validate()
			if err == nil {
				t.Fatal("error expected")
			}
			var errs []string
			for _, err := range err.(*errors.List).Errors {
				errs = append(errs, err.Error())
			}
			if diff := deep.Equal(errs, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}
