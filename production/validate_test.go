package production

import (
	"encoding/xml"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	rootPath := go3mf.DefaultModelPath
	type args struct {
		model *go3mf.Model
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{"empty", args{new(go3mf.Model)}, nil},
		{"noSpace", args{&go3mf.Model{Build: go3mf.Build{ExtensionAttr: go3mf.ExtensionAttr{
			mustUUID(""),
		}}}}, nil},
		{"buildNoUUID", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}}},
			[]error{&specerr.BuildError{Err: &specerr.MissingFieldError{Name: attrProdUUID}}},
		},
		{"buildEmptyUUID", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Build: go3mf.Build{
			ExtensionAttr: go3mf.ExtensionAttr{mustUUID("")}}}}, []error{&specerr.BuildError{Err: ErrUUID}},
		},
		{"buildNonValidUUID", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Build: go3mf.Build{
			ExtensionAttr: go3mf.ExtensionAttr{mustUUID("a-b-c-d")}}}}, []error{
			&specerr.BuildError{Err: ErrUUID},
		}},
		{"extReq", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, RequiredExtensions: []string{ExtensionName},
			Resources: go3mf.Resources{Objects: []*go3mf.Object{{ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d481")},
				Components: []*go3mf.Component{{ExtensionAttr: go3mf.ExtensionAttr{
					&PathUUID{Path: "/other.model", UUID: UUID("f47ac10b-58cc-0372-8567-0e02b2c3d480")},
				}}}}}}, Build: go3mf.Build{
				ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}, Items: []*go3mf.Item{
					{ObjectID: 1, ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: UUID("f47ac10b-58cc-0372-8567-0e02b2c3d478"), Path: "/other.model"}}},
				}}}}, []error{}},
		{"items", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Build: go3mf.Build{
			ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}, Items: []*go3mf.Item{
				{ObjectID: 1, ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: UUID("f47ac10b-58cc-0372-8567-0e02b2c3d478"), Path: "/other.model"}}},
				{},
				{ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: ""}}},
				{ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: "a-b-c-d"}}},
			}}}}, []error{
			&specerr.ItemError{Index: 1, Err: &specerr.MissingFieldError{Name: attrProdUUID}},
			&specerr.ItemError{Index: 2, Err: &specerr.MissingFieldError{Name: attrProdUUID}},
			&specerr.ItemError{Index: 3, Err: ErrUUID},
			ErrExtRequired,
		}},
		{"components", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}}, Resources: go3mf.Resources{
			Objects: []*go3mf.Object{
				{},
				{ExtensionAttr: go3mf.ExtensionAttr{mustUUID("a-b-c-d")}},
				{ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d483")}, Components: []*go3mf.Component{
					{},
					{ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: UUID("")}}},
					{ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{UUID: UUID("a-b-c-d")}}},
				}},
			},
		}, Build: go3mf.Build{ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}}}}, []error{
			&specerr.ObjectError{Path: rootPath, Index: 0, Err: &specerr.MissingFieldError{Name: attrProdUUID}},
			&specerr.ObjectError{Path: rootPath, Index: 1, Err: ErrUUID},
			&specerr.ObjectError{Path: rootPath, Index: 2, Err: &specerr.ComponentError{Index: 0, Err: &specerr.MissingFieldError{Name: attrProdUUID}}},
			&specerr.ObjectError{Path: rootPath, Index: 2, Err: &specerr.ComponentError{Index: 1, Err: &specerr.MissingFieldError{Name: attrProdUUID}}},
			&specerr.ObjectError{Path: rootPath, Index: 2, Err: &specerr.ComponentError{Index: 2, Err: ErrUUID}},
		}},
		{"child", args{&go3mf.Model{Namespaces: []xml.Name{{Space: ExtensionName}},
			Build: go3mf.Build{ExtensionAttr: go3mf.ExtensionAttr{mustUUID("f47ac10b-58cc-0372-8567-0e02b2c3d479")}},
			Childs: map[string]*go3mf.ChildModel{
				"/other.model": &go3mf.ChildModel{Resources: go3mf.Resources{Objects: []*go3mf.Object{
					{Components: []*go3mf.Component{
						{ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{Path: "/b.model"}}},
					}},
				}}}}}}, []error{
			&specerr.ObjectError{Path: "/other.model", Index: 0, Err: &specerr.MissingFieldError{Name: attrProdUUID}},
			&specerr.ObjectError{Path: "/other.model", Index: 0, Err: &specerr.ComponentError{Index: 0, Err: &specerr.MissingFieldError{Name: attrProdUUID}}},
			&specerr.ObjectError{Path: "/other.model", Index: 0, Err: &specerr.ComponentError{Index: 0, Err: ErrRefInNonRoot}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Validate(tt.args.model)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}
