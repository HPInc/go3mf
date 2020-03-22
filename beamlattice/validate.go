package beamlattice

import (
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

func (bl *BeamLattice) Validate(m *go3mf.Model, path string, e interface{}) []error {
	var (
		obj *go3mf.Object
		ok  bool
	)
	if obj, ok = e.(*go3mf.Object); !ok || obj.Mesh == nil {
		return nil
	}

	var errs []error
	if obj.ObjectType != go3mf.ObjectTypeModel && obj.ObjectType != go3mf.ObjectTypeSolidSupport {
		errs = append(errs, specerr.ErrLatticeObjType)
	}
	if bl.MinLength == 0 {
		errs = append(errs, &specerr.MissingFieldError{Name: attrMinLength})
	}
	if bl.DefaultRadius == 0 {
		errs = append(errs, &specerr.MissingFieldError{Name: attrRadius})
	}
	if bl.ClipMode == ClipNone && bl.ClippingMeshID == 0 {
		errs = append(errs, specerr.ErrLatticeClippedNoMesh)
	}
	if bl.ClippingMeshID != 0 {
		if err := validateRefMesh(m, path, bl.ClippingMeshID, obj.ID); err != nil {
			errs = append(errs, err)
		}
	}
	if bl.RepresentationMeshID != 0 {
		if err := validateRefMesh(m, path, bl.RepresentationMeshID, obj.ID); err != nil {
			errs = append(errs, err)
		}
	}

	for i, b := range bl.Beams {
		if b.NodeIndices[0] == b.NodeIndices[1] {
			errs = append(errs, specerr.NewIndexed(path, b, i, specerr.ErrLatticeSameVertex))
		} else {
			l := len(obj.Mesh.Nodes)
			if int(b.NodeIndices[0]) >= l || int(b.NodeIndices[1]) >= l {
				errs = append(errs, specerr.NewIndexed(path, b, i, specerr.ErrIndexOutOfBounds))
			}
		}
		if b.Radius[0] != 0 && b.Radius[0] != bl.DefaultRadius && b.Radius[0] != b.Radius[1] {
			errs = append(errs, specerr.NewIndexed(path, b, i, specerr.ErrLatticeBeamR2))
		}
	}
	for i, set := range bl.BeamSets {
		var setErrs []error
		for j, ref := range set.Refs {
			if int(ref) >= len(set.Refs) {
				setErrs = append(setErrs, specerr.NewIndexed(path, ref, j, specerr.ErrIndexOutOfBounds))
			}
		}
		for _, err := range setErrs {
			errs = append(errs, specerr.NewIndexed(path, set, i, err))
		}
	}

	return errs
}

func validateRefMesh(m *go3mf.Model, path string, meshID, selfID uint32) error {
	if meshID == selfID {
		return specerr.ErrLatticeSelfReference
	}
	if res, ok := m.FindResources(path); ok {
		for _, r := range res.Objects {
			if r.ID == selfID {
				return specerr.ErrMissingResource
			}
			var lattice *BeamLattice
			if r.ID == meshID {
				if r.Mesh == nil || r.ObjectType != go3mf.ObjectTypeModel || r.Mesh.Extension.Get(&lattice) {
					return specerr.ErrLatticeInvalidMesh
				}
			}
		}
	}
	return nil
}
