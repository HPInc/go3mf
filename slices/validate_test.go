//+build integration

package slices

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	// rootPath := go3mf.DefaultModelPath
	tests := []struct {
		name string
		model *go3mf.Model
		want []error
	}{
		{"empty", new(go3mf.Model), []error{}},
		{"child", &go3mf.Model{Childs: map[string]*go3mf.ChildModel{
			"/other.model": &go3mf.ChildModel{Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&SliceStack{ID: 1},
			}}},
			"/that.model": &go3mf.ChildModel{Resources: go3mf.Resources{Assets: []go3mf.Asset{
				&SliceStack{ID: 2},
			}}},
		}}, []error{
			fmt.Errorf("/other.model@Resources@SliceStack#0: %v", specerr.ErrSlicesAndRefs),
			fmt.Errorf("/that.model@Resources@SliceStack#0: %v", specerr.ErrSlicesAndRefs),
		}},
		{"slicestack", &go3mf.Model{Resources: go3mf.Resources{
			Assets: []go3mf.Asset{&SliceStack{
				ID: 1,
			}},
		}}, []error{

		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.Validate()
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}
