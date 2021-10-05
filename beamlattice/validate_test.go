// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package beamlattice

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/hpinc/go3mf"
	"github.com/hpinc/go3mf/errors"
	"github.com/hpinc/go3mf/spec"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		model *go3mf.Model
		want  []string
	}{
		{"error in child", &go3mf.Model{Childs: map[string]*go3mf.ChildModel{
			"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{
				{ID: 1, Mesh: &go3mf.Mesh{Any: spec.Any{&BeamLattice{}}}},
			}}},
		}}, []string{
			fmt.Sprintf("/other.model@Resources@Object#0@Mesh: %v", errors.ErrInsufficientVertices),
			fmt.Sprintf("/other.model@Resources@Object#0@Mesh@BeamLattice: %v", &errors.MissingFieldError{Name: attrMinLength}),
			fmt.Sprintf("/other.model@Resources@Object#0@Mesh@BeamLattice: %v", &errors.MissingFieldError{Name: attrRadius}),
			fmt.Sprintf("/other.model@Resources@Object#0@Mesh@BeamLattice: %v", ErrLatticeClippedNoMesh),
		}},
		{"object without beamlattice", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 1, Mesh: &go3mf.Mesh{}},
		}}}, []string{
			fmt.Sprintf("Resources@Object#0@Mesh: %v", errors.ErrInsufficientVertices),
			fmt.Sprintf("Resources@Object#0@Mesh: %v", errors.ErrInsufficientTriangles),
		}},
		{"object with components", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 1, Components: &go3mf.Components{Component: []*go3mf.Component{{ObjectID: 2}}}},
		}}}, []string{
			fmt.Sprintf("Resources@Object#0@Components@Component#0: %v", errors.ErrMissingResource),
		}},
		{"object incorret type", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 1, Type: go3mf.ObjectTypeOther, Mesh: &go3mf.Mesh{Any: spec.Any{&BeamLattice{
				MinLength: 1, Radius: 1, ClipMode: ClipInside,
			}}}},
			{ID: 2, Type: go3mf.ObjectTypeSurface, Mesh: &go3mf.Mesh{Any: spec.Any{&BeamLattice{
				MinLength: 1, Radius: 1, ClipMode: ClipInside,
			}}}},
			{ID: 3, Type: go3mf.ObjectTypeSupport, Mesh: &go3mf.Mesh{Any: spec.Any{&BeamLattice{
				MinLength: 1, Radius: 1, ClipMode: ClipInside,
			}}}},
		}}}, []string{
			fmt.Sprintf("Resources@Object#0@Mesh@BeamLattice: %v", ErrLatticeObjType),
			fmt.Sprintf("Resources@Object#1@Mesh@BeamLattice: %v", ErrLatticeObjType),
			fmt.Sprintf("Resources@Object#2@Mesh@BeamLattice: %v", ErrLatticeObjType),
		}},
		{"incorrect mesh references", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 1, Mesh: &go3mf.Mesh{Vertices: go3mf.Vertices{Vertex: []go3mf.Point3D{{}, {}, {}}}, Any: spec.Any{nil}}},
			{ID: 2, Mesh: &go3mf.Mesh{Vertices: go3mf.Vertices{Vertex: []go3mf.Point3D{{}, {}, {}}}, Any: spec.Any{&BeamLattice{
				MinLength: 1, Radius: 1, ClippingMeshID: 100, RepresentationMeshID: 2,
			}}}},
			{ID: 3, Mesh: &go3mf.Mesh{Vertices: go3mf.Vertices{Vertex: []go3mf.Point3D{{}, {}, {}}}, Any: spec.Any{&BeamLattice{
				MinLength: 1, Radius: 1, ClippingMeshID: 1, RepresentationMeshID: 2,
			}}}},
		}}}, []string{
			fmt.Sprintf("Resources@Object#1@Mesh@BeamLattice: %v", errors.ErrMissingResource),
			fmt.Sprintf("Resources@Object#1@Mesh@BeamLattice: %v", errors.ErrRecursion),
			fmt.Sprintf("Resources@Object#2@Mesh@BeamLattice: %v", ErrLatticeInvalidMesh),
		}},
		{"incorrect beams", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 2, Mesh: &go3mf.Mesh{Vertices: go3mf.Vertices{Vertex: []go3mf.Point3D{{}, {}, {}}}, Any: spec.Any{&BeamLattice{
				MinLength: 1, Radius: 1, ClipMode: ClipInside, Beams: Beams{Beam: []Beam{
					{}, {Indices: [2]uint32{1, 1}, Radius: [2]float32{0.5, 0}}, {Indices: [2]uint32{1, 3}},
				},
				}}}}},
		}}}, []string{
			fmt.Sprintf("Resources@Object#0@Mesh@BeamLattice@Beam#0: %v", ErrLatticeSameVertex),
			fmt.Sprintf("Resources@Object#0@Mesh@BeamLattice@Beam#1: %v", ErrLatticeSameVertex),
			fmt.Sprintf("Resources@Object#0@Mesh@BeamLattice@Beam#1: %v", ErrLatticeBeamR2),
			fmt.Sprintf("Resources@Object#0@Mesh@BeamLattice@Beam#2: %v", errors.ErrIndexOutOfBounds),
		}},
		{"incorrect beamseat", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 2, Mesh: &go3mf.Mesh{Vertices: go3mf.Vertices{Vertex: []go3mf.Point3D{{}, {}, {}}}, Any: spec.Any{&BeamLattice{
				MinLength: 1, Radius: 1, ClipMode: ClipInside, Beams: Beams{Beam: []Beam{
					{Indices: [2]uint32{1, 2}},
				}}, BeamSets: BeamSets{BeamSet: []BeamSet{{Refs: []uint32{0, 2, 3}}}},
			}}}},
		}}}, []string{
			fmt.Sprintf("Resources@Object#0@Mesh@BeamLattice@BeamSet#0: %v", errors.ErrIndexOutOfBounds),
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
