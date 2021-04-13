package booleanoperations

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	validMesh1 := &go3mf.Object{ID: 1, Mesh: &go3mf.Mesh{Vertices: []go3mf.Point3D{{}, {}, {}, {}}, Triangles: []go3mf.Triangle{
		go3mf.NewTriangle(0, 1, 2), go3mf.NewTriangle(0, 3, 1), go3mf.NewTriangle(0, 2, 3), go3mf.NewTriangle(1, 3, 2),
	}}}
	validMesh2 := &go3mf.Object{ID: 2, Mesh: &go3mf.Mesh{Vertices: []go3mf.Point3D{{}, {}, {}, {}}, Triangles: []go3mf.Triangle{
		go3mf.NewTriangle(0, 1, 2), go3mf.NewTriangle(0, 3, 1), go3mf.NewTriangle(0, 2, 3), go3mf.NewTriangle(1, 3, 2),
	}}}

	tests := []struct {
		name  string
		model *go3mf.Model
		want  []string
	}{
		{"Missing Operation", &go3mf.Model{
			Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh1, validMesh2,
				{ID: 3,
					Components: &go3mf.Components{
						AnyAttr: go3mf.AnyAttr{
							&BooleanOperationAttr{Association: Association_physical},
						}, Component: []*go3mf.Component{
							{ObjectID: 1, Transform: go3mf.Matrix{1.0000, 0.0000, 0.0000, 0.0000, 1.0000, 0.0000, 0.0000, 0.0000, 1.0000, 34.1020, 35.1070, 5.1000}},
							{ObjectID: 2, Transform: go3mf.Matrix{1.00000, 0.0000, 0.0000, 0.0000, 1.0000, 0.0000, 0.0000, 0.0000, 1.000, 35.7020, 35.7070, 5.7000}},
						}},
				}}},
		}, []string{
			fmt.Sprintf("Resources@Object#2@Components: %v", &errors.MissingFieldError{Name: attrCompsBoolOperOperation}),
		}},
		{"Missing Association", &go3mf.Model{
			Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh1, validMesh2,
				{ID: 3,
					Components: &go3mf.Components{
						AnyAttr: go3mf.AnyAttr{
							&BooleanOperationAttr{Operation: BooleanOperation_union},
						}, Component: []*go3mf.Component{
							{ObjectID: 1, Transform: go3mf.Matrix{1.0000, 0.0000, 0.0000, 0.0000, 1.0000, 0.0000, 0.0000, 0.0000, 1.0000, 34.1020, 35.1070, 5.1000}},
							{ObjectID: 2, Transform: go3mf.Matrix{1.00000, 0.0000, 0.0000, 0.0000, 1.0000, 0.0000, 0.0000, 0.0000, 1.000, 35.7020, 35.7070, 5.7000}},
						}},
				}}},
		}, []string{
			fmt.Sprintf("Resources@Object#2@Components: %v", &errors.MissingFieldError{Name: attrCompsBoolOperAssociation}),
		}},
		{"Complete", &go3mf.Model{
			Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh1, validMesh2,
				{ID: 3,
					Components: &go3mf.Components{
						AnyAttr: go3mf.AnyAttr{
							&BooleanOperationAttr{Association: Association_physical, Operation: BooleanOperation_union},
						}, Component: []*go3mf.Component{
							{ObjectID: 1, Transform: go3mf.Matrix{1.0000, 0.0000, 0.0000, 0.0000, 1.0000, 0.0000, 0.0000, 0.0000, 1.0000, 34.1020, 35.1070, 5.1000}},
							{ObjectID: 2, Transform: go3mf.Matrix{1.00000, 0.0000, 0.0000, 0.0000, 1.0000, 0.0000, 0.0000, 0.0000, 1.000, 35.7020, 35.7070, 5.7000}},
						}},
				}}},
		}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.model.Extensions = []go3mf.Extension{DefaultExtension}
			err := tt.model.Validate()
			if len(tt.want) == 0 && err == nil {
				t.Log("Sucerss")
			} else {
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
			}
		})
	}
}
