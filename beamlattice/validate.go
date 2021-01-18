package beamlattice

import (
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func (Spec) Validate(m interface{}, path string, obj interface{}) error {
	if obj, ok := obj.(*go3mf.Object); ok {
		return validateObject(m.(*go3mf.Model), path, obj)
	}
	return nil
}

func validateObject(m *go3mf.Model, path string, obj *go3mf.Object) error {
	if obj.Mesh == nil {
		return nil
	}

	bl := GetBeamLattice(obj.Mesh)
	if bl == nil {
		return nil
	}

	var errs error

	if obj.Type != go3mf.ObjectTypeModel && obj.Type != go3mf.ObjectTypeSolidSupport {
		errs = errors.Append(errs, ErrLatticeObjType)
	}
	if bl.MinLength == 0 {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrMinLength))
	}
	if bl.Radius == 0 {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrRadius))
	}
	if bl.ClipMode == ClipNone && bl.ClippingMeshID == 0 {
		errs = errors.Append(errs, ErrLatticeClippedNoMesh)
	}
	if bl.ClippingMeshID != 0 {
		errs = errors.Append(errs, validateRefMesh(m, path, bl.ClippingMeshID, obj.ID))
	}
	if bl.RepresentationMeshID != 0 {
		errs = errors.Append(errs, validateRefMesh(m, path, bl.RepresentationMeshID, obj.ID))
	}

	for i, b := range bl.Beams {
		if b.Indices[0] == b.Indices[1] {
			errs = errors.Append(errs, errors.WrapIndex(ErrLatticeSameVertex, b, i))
		} else {
			l := len(obj.Mesh.Vertices)
			if int(b.Indices[0]) >= l || int(b.Indices[1]) >= l {
				errs = errors.Append(errs, errors.WrapIndex(errors.ErrIndexOutOfBounds, b, i))
			}
		}
		if b.Radius[0] != 0 && b.Radius[0] != bl.Radius && b.Radius[0] != b.Radius[1] {
			errs = errors.Append(errs, errors.WrapIndex(ErrLatticeBeamR2, b, i))
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

func validateRefMesh(m *go3mf.Model, path string, meshID, selfID uint32) error {
	if meshID == selfID {
		return errors.ErrRecursion
	}
	if res, ok := m.FindResources(path); ok {
		for _, r := range res.Objects {
			if r.ID == selfID {
				return errors.ErrMissingResource
			}
			if r.ID == meshID {
				if r.Mesh == nil || r.Type != go3mf.ObjectTypeModel || GetBeamLattice(r.Mesh) != nil {
					return ErrLatticeInvalidMesh
				}
				break
			}
		}
	}
	return nil
}
