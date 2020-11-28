package production

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	validMesh := &go3mf.Object{ID: 1, Mesh: &go3mf.Mesh{Vertices: []go3mf.Point3D{{}, {}, {}, {}}, Triangles: []go3mf.Triangle{
		go3mf.NewTriangle(0, 1, 2), go3mf.NewTriangle(0, 3, 1), go3mf.NewTriangle(0, 2, 3), go3mf.NewTriangle(1, 3, 2),
	}}}
	tests := []struct {
		name  string
		model *go3mf.Model
		want  []error
	}{
		{"buildNoUUID", &go3mf.Model{Build: go3mf.Build{}}, []error{
			fmt.Errorf("Build: %v", &errors.MissingFieldError{Name: attrProdUUID}),
		}},
		{"buildEmptyUUID", &go3mf.Model{Build: go3mf.Build{
			AnyAttr: go3mf.ExtensionsAttr{mustUUID("")}}}, []error{
			fmt.Errorf("Build: %v", errors.ErrUUID),
		}},
		{"buildNonValidUUID", &go3mf.Model{Build: go3mf.Build{
			AnyAttr: go3mf.ExtensionsAttr{mustUUID("a-b-c-d")}}}, []error{
			fmt.Errorf("Build: %v", errors.ErrUUID),
		}},
		{"extReq", &go3mf.Model{Specs: map[string]go3mf.Spec{Namespace: &Spec{IsRequired: true}},
			Childs: map[string]*go3mf.ChildModel{"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}}},
			Resources: go3mf.Resources{Objects: []*go3mf.Object{
				{ID: 5, AnyAttr: go3mf.ExtensionsAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d481")}, Components: []*go3mf.Component{
					{ObjectID: 1, AnyAttr: go3mf.ExtensionsAttr{
						&PathUUID{Path: "/other.model", UUID: UUID("f47ac10b-58cc-0372-8567-0e02b2c3d480")},
					}}}}}}, Build: go3mf.Build{
				AnyAttr: go3mf.ExtensionsAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}, Items: []*go3mf.Item{
					{ObjectID: 1, AnyAttr: go3mf.ExtensionsAttr{&PathUUID{UUID: UUID("f47ac10b-58cc-0372-8567-0e02b2c3d478"), Path: "/other.model"}}},
				}}}, []error{
			fmt.Errorf("/other.model@Resources@Object#0: %v", &errors.MissingFieldError{Name: attrProdUUID}),
		}},
		{"items", &go3mf.Model{Build: go3mf.Build{
			AnyAttr: go3mf.ExtensionsAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}, Items: []*go3mf.Item{
				{ObjectID: 1, AnyAttr: go3mf.ExtensionsAttr{&PathUUID{UUID: UUID("f47ac10b-58cc-0372-8567-0e02b2c3d478"), Path: "/other.model"}}},
				{ObjectID: 1},
				{ObjectID: 1, AnyAttr: go3mf.ExtensionsAttr{&PathUUID{UUID: ""}}},
				{ObjectID: 1, AnyAttr: go3mf.ExtensionsAttr{&PathUUID{UUID: "a-b-c-d"}}},
			}},
			Childs:    map[string]*go3mf.ChildModel{"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}}},
			Resources: go3mf.Resources{Objects: []*go3mf.Object{{ID: 1, Mesh: validMesh.Mesh}}}}, []error{
			fmt.Errorf("Build@Item#0: %v", errors.ErrProdExtRequired),
			fmt.Errorf("Build@Item#1: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Errorf("Build@Item#2: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Errorf("Build@Item#3: %v", errors.ErrUUID),
			fmt.Errorf("/other.model@Resources@Object#0: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Errorf("Resources@Object#0: %v", &errors.MissingFieldError{Name: attrProdUUID}),
		}},
		{"components", &go3mf.Model{Resources: go3mf.Resources{
			Objects: []*go3mf.Object{
				{ID: 2, Mesh: validMesh.Mesh, AnyAttr: go3mf.ExtensionsAttr{mustUUID("a-b-c-d")}},
				{ID: 3, AnyAttr: go3mf.ExtensionsAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d483")}, Components: []*go3mf.Component{
					{ObjectID: 2, AnyAttr: go3mf.ExtensionsAttr{&PathUUID{UUID: UUID("")}}},
					{ObjectID: 2, AnyAttr: go3mf.ExtensionsAttr{&PathUUID{UUID: UUID("a-b-c-d")}}},
					{ObjectID: 2},
				}},
			},
		}, Build: go3mf.Build{AnyAttr: go3mf.ExtensionsAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}}}, []error{
			fmt.Errorf("Resources@Object#0: %v", errors.ErrUUID),
			fmt.Errorf("Resources@Object#1@Component#0: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Errorf("Resources@Object#1@Component#1: %v", errors.ErrUUID),
			fmt.Errorf("Resources@Object#1@Component#2: %v", &errors.MissingFieldError{Name: attrProdUUID}),
		}},
		{"child", &go3mf.Model{Build: go3mf.Build{AnyAttr: go3mf.ExtensionsAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}},
			Childs: map[string]*go3mf.ChildModel{
				"/b.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}},
				"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{
					{ID: 2, Components: []*go3mf.Component{
						{ObjectID: 1, AnyAttr: go3mf.ExtensionsAttr{&PathUUID{Path: "/b.model"}}},
					}},
				}}}}}, []error{
			fmt.Errorf("/b.model@Resources@Object#0: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Errorf("/other.model@Resources@Object#0: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Errorf("/other.model@Resources@Object#0@Component#0: %v", &errors.MissingFieldError{Name: attrProdUUID}),
			fmt.Errorf("/other.model@Resources@Object#0@Component#0: %v", errors.ErrProdRefInNonRoot),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.model.Specs) == 0 {
				tt.model.WithSpec(&Spec{})
			}
			got := tt.model.Validate()
			if diff := deep.Equal(got.(*errors.List).Errors, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}

func TestSpec_ValidateAsset(t *testing.T) {
	type args struct {
		in0 *go3mf.Model
		in1 string
		in2 go3mf.Asset
	}
	tests := []struct {
		name    string
		e       *Spec
		args    args
		wantErr bool
	}{
		{"base", &Spec{}, args{nil, "", nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.ValidateAsset(tt.args.in0, tt.args.in1, tt.args.in2); (err != nil) != tt.wantErr {
				t.Errorf("Spec.ValidateAsset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
