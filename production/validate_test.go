package production

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	validMesh := &go3mf.Object{ID: 1, Mesh: &go3mf.Mesh{Nodes: []go3mf.Point3D{{}, {}, {}, {}}, Faces: []go3mf.Face{
		{NodeIndices: [3]uint32{0, 1, 2}}, {NodeIndices: [3]uint32{0, 3, 1}}, {NodeIndices: [3]uint32{0, 2, 3}}, {NodeIndices: [3]uint32{1, 3, 2}},
	}}}
	rootPath := go3mf.DefaultModelPath
	type args struct {
		model *go3mf.Model
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{"buildEmptyUUID", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Build: go3mf.Build{
			ExtensionAttr: go3mf.ExtensionAttr{mustUUID("")}}}}, []error{
			fmt.Errorf("%s@Build: %v", rootPath, specerr.ErrUUID),
		}},
		{"buildNonValidUUID", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Build: go3mf.Build{
			ExtensionAttr: go3mf.ExtensionAttr{mustUUID("a-b-c-d")}}}}, []error{
			fmt.Errorf("%s@Build: %v", rootPath, specerr.ErrUUID),
		}},
		{"extReq", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, RequiredExtensions: []string{ExtensionName},
			Childs: map[string]*go3mf.ChildModel{"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}}},
			Resources: go3mf.Resources{Objects: []*go3mf.Object{
				{ID: 5, ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d481")}, Components: []*go3mf.Component{
					{ObjectID: 1, ExtensionAttr: go3mf.ExtensionAttr{
						&PathUUID{Path: "/other.model", UUID: UUID("f47ac10b-58cc-0372-8567-0e02b2c3d480")},
					}}}}}}, Build: go3mf.Build{
				ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}, Items: []*go3mf.Item{
					{ObjectID: 1, ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: UUID("f47ac10b-58cc-0372-8567-0e02b2c3d478"), Path: "/other.model"}}},
				}}}}, []error{}},
		{"items", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Build: go3mf.Build{
			ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}, Items: []*go3mf.Item{
				{ObjectID: 1, ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: UUID("f47ac10b-58cc-0372-8567-0e02b2c3d478"), Path: "/other.model"}}},
				{ObjectID: 1},
				{ObjectID: 1, ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: ""}}},
				{ObjectID: 1, ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: "a-b-c-d"}}},
			}},
			Childs:    map[string]*go3mf.ChildModel{"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}}},
			Resources: go3mf.Resources{Objects: []*go3mf.Object{{ID: 1, Mesh: validMesh.Mesh}}}}}, []error{
			fmt.Errorf("%s@Build@Item#0: %v", rootPath, specerr.ErrProdExtRequired),
			fmt.Errorf("%s@Build@Item#2: %v", rootPath, &specerr.MissingFieldError{Name: attrProdUUID}),
			fmt.Errorf("%s@Build@Item#3: %v", rootPath, specerr.ErrUUID),
		}},
		{"components", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Resources: go3mf.Resources{
			Objects: []*go3mf.Object{
				{ID: 2, Mesh: validMesh.Mesh, ExtensionAttr: go3mf.ExtensionAttr{mustUUID("a-b-c-d")}},
				{ID: 3, ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d483")}, Components: []*go3mf.Component{
					{ObjectID: 2, ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: UUID("")}}},
					{ObjectID: 2, ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: UUID("a-b-c-d")}}},
				}},
			},
		}, Build: go3mf.Build{ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}}}}, []error{
			fmt.Errorf("%s@Resources@Object#0: %v", rootPath, specerr.ErrUUID),
			fmt.Errorf("%s@Resources@Object#1@Component#0: %v", rootPath, &specerr.MissingFieldError{Name: attrProdUUID}),
			fmt.Errorf("%s@Resources@Object#1@Component#1: %v", rootPath, specerr.ErrUUID),
		}},
		{"child", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}},
			Build: go3mf.Build{ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}},
			Childs: map[string]*go3mf.ChildModel{
				"/b.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}},
				"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{
					{ID: 2, Components: []*go3mf.Component{
						{ObjectID: 1, ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{Path: "/b.model"}}},
					}},
				}}}}}}, []error{
			fmt.Errorf("/other.model@Resources@Object#0@Component#0: %v", &specerr.MissingFieldError{Name: attrProdUUID}),
			fmt.Errorf("/other.model@Resources@Object#0@Component#0: %v", specerr.ErrProdRefInNonRoot),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.model.Validate()
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}
