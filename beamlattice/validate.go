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

	errs := new(errors.ErrorList)

	if obj.ObjectType != go3mf.ObjectTypeModel && obj.ObjectType != go3mf.ObjectTypeSolidSupport {
		errs.Append(errors.ErrLatticeObjType)
	}
	if bl.MinLength == 0 {
		errs.Append(&errors.MissingFieldError{Name: attrMinLength})
	}
	if bl.DefaultRadius == 0 {
		errs.Append(&errors.MissingFieldError{Name: attrRadius})
	}
	if bl.ClipMode == ClipNone && bl.ClippingMeshID == 0 {
		errs.Append(errors.ErrLatticeClippedNoMesh)
	}
	if bl.ClippingMeshID != 0 {
		errs.Append(e.validateRefMesh(m, path, bl.ClippingMeshID, obj.ID))
	}
	if bl.RepresentationMeshID != 0 {
		errs.Append(e.validateRefMesh(m, path, bl.RepresentationMeshID, obj.ID))
	}

	for i, b := range bl.Beams {
		if b.Indices[0] == b.Indices[1] {
			errs.Append(errors.NewIndexed(b, i, errors.ErrLatticeSameVertex))
		} else {
			l := len(obj.Mesh.Vertices)
			if int(b.Indices[0]) >= l || int(b.Indices[1]) >= l {
				errs.Append(errors.NewIndexed(b, i, errors.ErrIndexOutOfBounds))
			}
		}
		if b.Radius[0] != 0 && b.Radius[0] != bl.DefaultRadius && b.Radius[0] != b.Radius[1] {
			errs.Append(errors.NewIndexed(b, i, errors.ErrLatticeBeamR2))
		}
	}
	for i, set := range bl.BeamSets {
		for _, ref := range set.Refs {
			if int(ref) >= len(set.Refs) {
				errs.Append(errors.NewIndexed(set, i, errors.ErrIndexOutOfBounds))
				break
			}
		}
	}

	return errors.New(obj.Mesh, errors.New(bl, errs))
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
				if r.Mesh == nil || r.ObjectType != go3mf.ObjectTypeModel || r.Mesh.Any.Get(&lattice) {
					return errors.ErrLatticeInvalidMesh
				}
				break
			}
		}
	}
	return nil
}
