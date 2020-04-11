package slices

import (
	"math"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func validTransform(t go3mf.Matrix) bool {
	return t[2] == 0 && t[6] == 0 && t[8] == 0 && t[9] == 0 && t[10] == 1
}

func (e *Spec) ValidateModel(_ *go3mf.Model) error {
	return nil
}

func (e *Spec) ValidateObject(m *go3mf.Model, path string, obj *go3mf.Object) error {
	var sti *SliceStackInfo
	if !obj.AnyAttr.Get(&sti) {
		return nil
	}
	var errs error
	res, _ := m.FindResources(path)
	if sti.SliceStackID == 0 {
		errs = errors.Append(errs, errors.NewMissingFieldError(attrSliceRefID))
	} else if r, ok := res.FindAsset(sti.SliceStackID); ok {
		if r, ok := r.(*SliceStack); ok {
			if !validateBuildTransforms(m, path, obj.ID) {
				errs = errors.Append(errs, errors.ErrSliceInvalidTranform)
			}
			if obj.Type == go3mf.ObjectTypeModel || obj.Type == go3mf.ObjectTypeSolidSupport {
				if !checkAllClosed(m, r) {
					errs = errors.Append(errs, errors.ErrSlicePolygonNotClosed)
				}
			}
		} else {
			errs = errors.Append(errs, errors.ErrNonSliceStack)
		}
	} else {
		errs = errors.Append(errs, errors.ErrMissingResource)
	}
	if sti.MeshResolution == ResolutionLow {
		if !e.Required() {
			errs = errors.Append(errs, errors.ErrSliceExtRequired)
		}
	}
	return errs
}

func (e *Spec) ValidateAsset(m *go3mf.Model, path string, r go3mf.Asset) error {
	var (
		st *SliceStack
		ok bool
	)
	if st, ok = r.(*SliceStack); !ok {
		return nil
	}
	var errs error
	if (len(st.Slices) != 0 && len(st.Refs) != 0) ||
		(len(st.Slices) == 0 && len(st.Refs) == 0) {
		errs = errors.Append(errs, errors.ErrSlicesAndRefs)
	}
	errs = errors.Append(errs, st.validateRefs(m, path))
	errs = errors.Append(errs, st.validateSlices())
	return errs
}

func (r *SliceStack) validateSlices() error {
	var errs error
	lastTopZ := float32(-math.MaxFloat32)
	for j, slice := range r.Slices {
		if slice.TopZ == 0 {
			errs = errors.Append(errs, errors.WrapIndex(errors.NewMissingFieldError(attrZTop), slice, j))
		} else if slice.TopZ < r.BottomZ {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrSliceSmallTopZ, slice, j))
		}
		if slice.TopZ <= lastTopZ {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrSliceNoMonotonic, slice, j))
		}
		lastTopZ = slice.TopZ
		if len(slice.Polygons) == 0 && len(slice.Vertices) == 0 {
			continue
		}
		if len(slice.Vertices) < 2 {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrSliceInsufficientVertices, slice, j))
		}
		if len(slice.Polygons) == 0 {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrSliceInsufficientPolygons, slice, j))
		}
		var perrs error
		for k, p := range slice.Polygons {
			if len(p.Segments) < 1 {
				perrs = errors.Append(perrs, errors.WrapIndex(errors.ErrSliceInsufficientSegments, p, k))
			}
		}
		if perrs != nil {
			errs = errors.Append(errs, errors.WrapIndex(perrs, slice, j))
		}
	}
	return errs
}

func (r *SliceStack) validateRefs(m *go3mf.Model, path string) error {
	var errs error
	lastTopZ := float32(-math.MaxFloat32)
	for i, ref := range r.Refs {
		valid := true
		if ref.Path == "" {
			valid = false
			errs = errors.Append(errs, errors.WrapIndex(errors.NewMissingFieldError(attrSlicePath), ref, i))
		} else if ref.Path == path {
			valid = false
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrSliceRefSamePart, ref, i))
		}
		if ref.SliceStackID == 0 {
			valid = false
			errs = errors.Append(errs, errors.WrapIndex(errors.NewMissingFieldError(attrSliceRefID), ref, i))
		}
		if !valid {
			continue
		}
		if st, ok := m.FindAsset(ref.Path, ref.SliceStackID); ok {
			if st, ok := st.(*SliceStack); ok {
				if len(st.Refs) != 0 {
					errs = errors.Append(errs, errors.WrapIndex(errors.ErrSliceRefRef, ref, i))
				}
				if len(st.Slices) > 0 && st.Slices[0].TopZ <= lastTopZ {
					errs = errors.Append(errs, errors.WrapIndex(errors.ErrSliceNoMonotonic, ref, i))
				}
				if len(st.Slices) > 0 {
					lastTopZ = st.Slices[len(st.Slices)-1].TopZ
				}
			} else {
				errs = errors.Append(errs, errors.WrapIndex(errors.ErrNonSliceStack, ref, i))
			}
		} else {
			errs = errors.Append(errs, errors.WrapIndex(errors.ErrMissingResource, ref, i))
		}
	}
	return errs
}

func isSliceStackClosed(r *SliceStack) bool {
	for _, slice := range r.Slices {
		for _, p := range slice.Polygons {
			if len(p.Segments) > 0 && p.StartV != p.Segments[len(p.Segments)-1].V2 {
				return false
			}
		}
	}
	return true
}

func checkAllClosed(m *go3mf.Model, r *SliceStack) bool {
	if !isSliceStackClosed(r) {
		return false
	}
	for _, ref := range r.Refs {
		if st, ok := m.FindAsset(ref.Path, ref.SliceStackID); ok {
			if st, ok := st.(*SliceStack); ok {
				if !isSliceStackClosed(st) {
					return false
				}
			}
		}
	}

	return true
}

func validateBuildTransforms(m *go3mf.Model, path string, id uint32) bool {
	for _, item := range m.Build.Items {
		targetPath := item.ObjectPath()
		if targetPath == "" {
			targetPath = path
		}
		if item.ObjectID == id && targetPath == path {
			if item.HasTransform() && !validTransform(item.Transform) {
				return false
			}
		}
		if o, ok := m.FindObject(targetPath, item.ObjectID); ok {
			if !validateObjectTransforms(m, o, path, id) {
				return false
			}
		}
	}
	return true
}

func validateObjectTransforms(m *go3mf.Model, o *go3mf.Object, path string, id uint32) bool {
	for _, c := range o.Components {
		if c.ObjectID == id && c.ObjectPath(path) == path {
			if c.HasTransform() && !validTransform(c.Transform) {
				return false
			}
		}
		if c.ObjectID == o.ID && c.ObjectPath(path) == path { // avoid circular references
			break
		} else {
			if o1, ok := m.FindObject(c.ObjectPath(path), c.ObjectID); ok {
				if !validateObjectTransforms(m, o1, path, id) {
					return false
				}
			}
		}
	}
	return true
}
