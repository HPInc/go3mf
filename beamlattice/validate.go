package beamlattice

import (
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func (e *Spec) ValidateModel(_ *go3mf.Model) error {
	return nil
}

func (e *Spec) ValidateAsset(_ *go3mf.Model, _ string, _ go3mf.Asset) error {
	return nil
}

func (e *Spec) ValidateObject(m *go3mf.Model, path string, obj *go3mf.Object) error {
	if obj.Mesh == nil {
		return nil
	}

	var bl *BeamLattice
	if !obj.Mesh.Any.Get(&bl) {
		return nil
	}

	var errs error

	if obj.Type != go3mf.ObjectTypeModel && obj.Type != go3mf.ObjectTypeSolidSupport {
		errs = errors.Append(errs, errors.ErrLatticeObjType)
	}
	if bl.MinLength == 0 {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrMinLength))
	}
	if bl.DefaultRadius == 0 {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrRadius))
	}
	if bl.ClipMode == ClipNone && bl.ClippingMeshID == 0 {
		errs = errors.Append(errs, errors.ErrLatticeClippedNoMesh)
	}
	if bl.ClippingMeshID != 0 {
		errs = errors.Append(errs, e.validateRefMesh(m, path, bl.ClippingMeshID, obj.ID))
	}
	if bl.RepresentationMeshID != 0 {
		errs = errors.Append(errs, e.validateRefMesh(m, path, bl.RepresentationMeshID, obj.ID))
	}

	for i, b := range bl.Beams {
		if b.Indices[0] == b.Indices[1] {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrLatticeSameVertex, b, i))
		} else {
			l := len(obj.Mesh.Vertices)
			if int(b.Indices[0]) >= l || int(b.Indices[1]) >= l {
				errs = errors.Append(errs, errors.WrapIndex(errors.ErrIndexOutOfBounds, b, i))
			}
		}
		if b.Radius[0] != 0 && b.Radius[0] != bl.DefaultRadius && b.Radius[0] != b.Radius[1] {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrLatticeBeamR2, b, i))
		}
	}
	for i, set := range bl.BeamSets {
		for _, ref := range set.Refs {
			if int(ref) >= len(set.Refs) {
				errs = errors.Append(errs, errors.WrapIndex(errors.ErrIndexOutOfBounds, set, i))
				break
			}
		}
	}
	if errs != nil {
		errs = errors.Wrap(errors.Wrap(errs, bl), obj.Mesh)
	}
	return errs
}

func (e *Spec) validateRefMesh(m *go3mf.Model, path string, meshID, selfID uint32) error {
	if meshID == selfID {
		return errors.ErrRecursion
	}
	if res, ok := m.FindResources(path); ok {
		for _, r := range res.Objects {
			if r.ID == selfID {
				return errors.ErrMissingResource
			}
			var lattice *BeamLattice
			if r.ID == meshID {
				if r.Mesh == nil || r.Type != go3mf.ObjectTypeModel || r.Mesh.Any.Get(&lattice) {
					return errors.ErrLatticeInvalidMesh
				}
				break
			}
		}
	}
	return nil
}
