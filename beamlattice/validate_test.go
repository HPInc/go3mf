package beamlattice

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		model *go3mf.Model
		want  []error
	}{
		{"error in child", &go3mf.Model{Childs: map[string]*go3mf.ChildModel{
			"/other.model": {Resources: go3mf.Resources{Objects: []*go3mf.Object{
				{ID: 1, Mesh: &go3mf.Mesh{Any: go3mf.Marshalers{&BeamLattice{}}}},
			}}},
		}}, []error{
			fmt.Errorf("/other.model@Resources@Object#0@Mesh: %v", errors.ErrInsufficientVertices),
			fmt.Errorf("/other.model@Resources@Object#0@Mesh@BeamLattice: %v", &errors.MissingFieldError{Name: attrMinLength}),
			fmt.Errorf("/other.model@Resources@Object#0@Mesh@BeamLattice: %v", &errors.MissingFieldError{Name: attrRadius}),
			fmt.Errorf("/other.model@Resources@Object#0@Mesh@BeamLattice: %v", errors.ErrLatticeClippedNoMesh),
		}},
		{"object without beamlattice", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 1, Mesh: &go3mf.Mesh{}},
		}}}, []error{
			fmt.Errorf("Resources@Object#0@Mesh: %v", errors.ErrInsufficientVertices),
			fmt.Errorf("Resources@Object#0@Mesh: %v", errors.ErrInsufficientTriangles),
		}},
		{"object with components", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 1, Components: []*go3mf.Component{{ObjectID: 2}}},
		}}}, []error{
			fmt.Errorf("Resources@Object#0@Component#0: %v", errors.ErrMissingResource),
		}},
		{"object incorret type", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 1, ObjectType: go3mf.ObjectTypeOther, Mesh: &go3mf.Mesh{Any: go3mf.Marshalers{&BeamLattice{
				MinLength: 1, DefaultRadius: 1, ClipMode: ClipInside,
			}}}},
			{ID: 2, ObjectType: go3mf.ObjectTypeSurface, Mesh: &go3mf.Mesh{Any: go3mf.Marshalers{&BeamLattice{
				MinLength: 1, DefaultRadius: 1, ClipMode: ClipInside,
			}}}},
			{ID: 3, ObjectType: go3mf.ObjectTypeSupport, Mesh: &go3mf.Mesh{Any: go3mf.Marshalers{&BeamLattice{
				MinLength: 1, DefaultRadius: 1, ClipMode: ClipInside,
			}}}},
		}}}, []error{
			fmt.Errorf("Resources@Object#0@Mesh@BeamLattice: %v", errors.ErrLatticeObjType),
			fmt.Errorf("Resources@Object#1@Mesh@BeamLattice: %v", errors.ErrLatticeObjType),
			fmt.Errorf("Resources@Object#2@Mesh@BeamLattice: %v", errors.ErrLatticeObjType),
		}},
		{"incorrect mesh references", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 1, Mesh: &go3mf.Mesh{Vertices: []go3mf.Point3D{{}, {}, {}}, Any: go3mf.Marshalers{nil}}},
			{ID: 2, Mesh: &go3mf.Mesh{Vertices: []go3mf.Point3D{{}, {}, {}}, Any: go3mf.Marshalers{&BeamLattice{
				MinLength: 1, DefaultRadius: 1, ClippingMeshID: 100, RepresentationMeshID: 2,
			}}}},
			{ID: 3, Mesh: &go3mf.Mesh{Vertices: []go3mf.Point3D{{}, {}, {}}, Any: go3mf.Marshalers{&BeamLattice{
				MinLength: 1, DefaultRadius: 1, ClippingMeshID: 1, RepresentationMeshID: 2,
			}}}},
		}}}, []error{
			fmt.Errorf("Resources@Object#1@Mesh@BeamLattice: %v", errors.ErrMissingResource),
			fmt.Errorf("Resources@Object#1@Mesh@BeamLattice: %v", errors.ErrRecursion),
			fmt.Errorf("Resources@Object#2@Mesh@BeamLattice: %v", errors.ErrLatticeInvalidMesh),
		}},
		{"incorrect beams", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 2, Mesh: &go3mf.Mesh{Vertices: []go3mf.Point3D{{}, {}, {}}, Any: go3mf.Marshalers{&BeamLattice{
				MinLength: 1, DefaultRadius: 1, ClipMode: ClipInside, Beams: []Beam{
					{}, {Indices: [2]uint32{1, 1}, Radius: [2]float32{0.5, 0}}, {Indices: [2]uint32{1, 3}},
				},
			}}}},
		}}}, []error{
			fmt.Errorf("Resources@Object#0@Mesh@BeamLattice@Beam#0: %v", errors.ErrLatticeSameVertex),
			fmt.Errorf("Resources@Object#0@Mesh@BeamLattice@Beam#1: %v", errors.ErrLatticeSameVertex),
			fmt.Errorf("Resources@Object#0@Mesh@BeamLattice@Beam#1: %v", errors.ErrLatticeBeamR2),
			fmt.Errorf("Resources@Object#0@Mesh@BeamLattice@Beam#2: %v", errors.ErrIndexOutOfBounds),
		}},
		{"incorrect beamseat", &go3mf.Model{Resources: go3mf.Resources{Objects: []*go3mf.Object{
			{ID: 2, Mesh: &go3mf.Mesh{Vertices: []go3mf.Point3D{{}, {}, {}}, Any: go3mf.Marshalers{&BeamLattice{
				MinLength: 1, DefaultRadius: 1, ClipMode: ClipInside, Beams: []Beam{
					{Indices: [2]uint32{1, 2}},
				}, BeamSets: []BeamSet{{Refs: []uint32{0, 2, 3}}},
			}}}},
		}}}, []error{
			fmt.Errorf("Resources@Object#0@Mesh@BeamLattice@BeamSet#0: %v", errors.ErrIndexOutOfBounds),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.model.WithSpec(&Spec{})
			got := tt.model.Validate()
			if diff := deep.Equal(got.(*errors.ErrorList).Errors, tt.want); diff != nil {
				t.Errorf("Validate() = %v", diff)
			}
		})
	}
}
