package beamlattice

import (
	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

func (e *Spec) ValidateModel(_ *go3mf.Model) []error {
	return nil
}

func (e *Spec) ValidateAsset(_ *go3mf.Model, _ string, _ go3mf.Asset) []error {
	return nil
}

func (e *Spec) ValidateObject(m *go3mf.Model, path string, obj *go3mf.Object) []error {
	if obj.Mesh == nil {
		return nil
	}

	var errs []error
	var bl *BeamLattice
	if !obj.Mesh.Any.Get(&bl) {
		return errs
	}

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
		if err := e.validateRefMesh(m, path, bl.ClippingMeshID, obj.ID); err != nil {
			errs = append(errs, err)
		}
	}
	if bl.RepresentationMeshID != 0 {
		if err := e.validateRefMesh(m, path, bl.RepresentationMeshID, obj.ID); err != nil {
			errs = append(errs, err)
		}
	}

	for i, b := range bl.Beams {
		if b.NodeIndices[0] == b.NodeIndices[1] {
			errs = append(errs, specerr.NewIndexed(b, i, specerr.ErrLatticeSameVertex))
		} else {
			l := len(obj.Mesh.Vertices)
			if int(b.NodeIndices[0]) >= l || int(b.NodeIndices[1]) >= l {
				errs = append(errs, specerr.NewIndexed(b, i, specerr.ErrIndexOutOfBounds))
			}
		}
		if b.Radius[0] != 0 && b.Radius[0] != bl.DefaultRadius && b.Radius[0] != b.Radius[1] {
			errs = append(errs, specerr.NewIndexed(b, i, specerr.ErrLatticeBeamR2))
		}
	}
	for i, set := range bl.BeamSets {
		for _, ref := range set.Refs {
			if int(ref) >= len(set.Refs) {
				errs = append(errs, specerr.NewIndexed(set, i, specerr.ErrIndexOutOfBounds))
				break
			}
		}
	}

	for i, err := range errs {
		errs[i] = specerr.New(obj.Mesh, specerr.New(bl, err))
	}

	return errs
}

func (e *Spec) validateRefMesh(m *go3mf.Model, path string, meshID, selfID uint32) error {
	if meshID == selfID {
		return specerr.ErrRecursion
	}
	if res, ok := m.FindResources(path); ok {
		for _, r := range res.Objects {
			if r.ID == selfID {
				return specerr.ErrMissingResource
			}
			var lattice *BeamLattice
			if r.ID == meshID {
				if r.Mesh == nil || r.ObjectType != go3mf.ObjectTypeModel || r.Mesh.Any.Get(&lattice) {
					return specerr.ErrLatticeInvalidMesh
				}
				break
			}
		}
	}
	return nil
}
