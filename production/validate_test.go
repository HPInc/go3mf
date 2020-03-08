package production

import (
	"encoding/xml"
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
		{"empty", args{new(go3mf.Model)}, []error{}},
		{"buildEmptyUUID", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Build: go3mf.Build{
			ExtensionAttr: go3mf.ExtensionAttr{mustUUID("")}}}}, []error{&specerr.BuildError{Err: specerr.ErrUUID}},
		},
		{"buildNonValidUUID", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Build: go3mf.Build{
			ExtensionAttr: go3mf.ExtensionAttr{mustUUID("a-b-c-d")}}}}, []error{
			&specerr.BuildError{Err: specerr.ErrUUID},
		}},
		{"extReq", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, RequiredExtensions: []string{ExtensionName},
			Childs: map[string]*go3mf.ChildModel{"/other.model": &go3mf.ChildModel{Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}}},
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
			Childs:    map[string]*go3mf.ChildModel{"/other.model": &go3mf.ChildModel{Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}}},
			Resources: go3mf.Resources{Objects: []*go3mf.Object{&go3mf.Object{ID: 1, Mesh: validMesh.Mesh}}}}}, []error{
			&specerr.ItemError{Index: 0, Err: specerr.ErrProdExtRequired},
			&specerr.ItemError{Index: 2, Err: &specerr.MissingFieldError{Name: attrProdUUID}},
			&specerr.ItemError{Index: 3, Err: specerr.ErrUUID},
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
			&specerr.ObjectError{Path: rootPath, Index: 0, Err: specerr.ErrUUID},
			&specerr.ObjectError{Path: rootPath, Index: 1, Err: &specerr.IndexedError{Name: "component", Index: 0, Err: &specerr.MissingFieldError{Name: attrProdUUID}}},
			&specerr.ObjectError{Path: rootPath, Index: 1, Err: &specerr.IndexedError{Name: "component", Index: 1, Err: specerr.ErrUUID}},
		}},
		{"child", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}},
			Build: go3mf.Build{ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}},
			Childs: map[string]*go3mf.ChildModel{
				"/b.model": &go3mf.ChildModel{Resources: go3mf.Resources{Objects: []*go3mf.Object{validMesh}}},
				"/other.model": &go3mf.ChildModel{Resources: go3mf.Resources{Objects: []*go3mf.Object{
					{ID: 2, Components: []*go3mf.Component{
						{ObjectID: 1, ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{Path: "/b.model"}}},
					}},
				}}}}}}, []error{
			&specerr.ObjectError{Path: "/other.model", Index: 0, Err: &specerr.IndexedError{Name: "component", Index: 0, Err: &specerr.MissingFieldError{Name: attrProdUUID}}},
			&specerr.ObjectError{Path: "/other.model", Index: 0, Err: &specerr.IndexedError{Name: "component", Index: 0, Err: specerr.ErrProdRefInNonRoot}},
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
