package slices

import (
	"sort"

	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
)

type resource struct {
	path string
	id   uint32
}

// Validate checks that the model is conformant with the 3MF spec.
// Core spec related checks are not reported.
func Validate(model *go3mf.Model) []error {
	var hasExt bool
	for _, ext := range model.Namespaces {
		if ext.Space == ExtensionName {
			hasExt = true
			break
		}
	}
	if !hasExt {
		return nil
	}
	err := make([]error, 0)
	stacks := make(map[resource]struct{})
	var mustRequire1, mustRequire2 bool
	mustRequire1, err = validateChilds(model, stacks, err)
	mustRequire2, err = validateRoot(model, stacks, err)
	if mustRequire1 || mustRequire2 {
		var extRequired bool
		for _, r := range model.RequiredExtensions {
			if r == ExtensionName {
				extRequired = true
				break
			}
		}
		if !extRequired {
			err = append(err, specerr.ErrProdExtRequired)
		}
	}
	return err
}

func validTransform(t go3mf.Matrix) bool {
	return t[2] == 0 && t[6] == 0 && t[8] == 0 && t[9] == 0 && t[10] == 1
}

func validateRoot(model *go3mf.Model, stacks map[resource]struct{}, err []error) (bool, []error) {
	var mustRequire bool
	path := model.PathOrDefault()
	mustRequire, err = validateObjects(path, &model.Resources, stacks, err)
	err = validateAssets(path, model, &model.Resources, stacks, err)
	var ext *SliceStackInfo
	for i, item := range model.Build.Items {
		if item.HasTransform() && !validTransform(item.Transform) {
			if obj, ok := model.FindObject(item.ObjectPath(path), item.ObjectID); ok {
				if obj.ExtensionAttr.Get(&ext) {
					err = append(err, specerr.NewItem(i, specerr.ErrSliceInvalidTranform))
				}
			}
		}
	}
	return mustRequire, err
}

func validateAssets(path string, model *go3mf.Model, res *go3mf.Resources, stacks map[resource]struct{}, err []error) []error {
	var lastTopZ float32
	for i, r := range res.Assets {
		if r, ok := r.(*SliceStackResource); ok {
			if len(r.Slices) != 0 && len(r.Refs) != 0 {
				err = append(err, specerr.NewAsset(path, i, r, specerr.ErrSlicesAndRefs))
			}
			err = validateRefs(path, model, i, r, stacks, err)
			for j, slice := range r.Slices {
				if slice.TopZ == 0 {
					err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: &specerr.MissingFieldError{Name: attrZTop}}))
				} else if slice.TopZ < r.BottomZ {
					err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: specerr.ErrSliceSmallTopZ}))
				}
				if len(slice.Polygons) == 0 && len(slice.Vertices) == 0 {
					continue
				}
				if slice.TopZ < lastTopZ {
					err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: specerr.ErrSliceNoMonotonic}))
				}
				if len(slice.Vertices) < 2 {
					err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: specerr.ErrSliceInsufficientVertices}))
				}
				if len(slice.Polygons) < 2 {
					err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: specerr.ErrSliceInsufficientPolygons}))
				}
				for k, p := range slice.Polygons {
					if len(p.Segments) < 1 {
						err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: &specerr.SliceSegmentError{Index: k, Err: specerr.ErrSliceInsufficientSegments}}))
					} else if _, ok := stacks[resource{path, r.ID}]; ok {
						if p.StartV != p.Segments[len(p.Segments)-1].V2 {
							err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: &specerr.SliceSegmentError{Index: k, Err: specerr.ErrSlicePolygonNotClosed}}))
						}
					}
				}
				lastTopZ = slice.TopZ
			}
		}
	}
	return err
}

func validateRefs(path string, model *go3mf.Model, i int, r *SliceStackResource, stacks map[resource]struct{}, err []error) []error {
	var lastTopZ float32
	for j, ref := range r.Refs {
		valid := true
		if ref.Path == "" {
			valid = false
			err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: &specerr.MissingFieldError{Name: attrSlicePath}}))
		} else if ref.Path == path {
			valid = false
			err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: specerr.ErrSliceRefSamePart}))
		}
		if ref.SliceStackID == 0 {
			valid = false
			err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: &specerr.MissingFieldError{Name: attrSliceRefID}}))
		}
		if !valid {
			continue
		}
		if st, ok := model.FindAsset(ref.Path, ref.SliceStackID); ok {
			if st, ok := st.(*SliceStackResource); ok {
				if _, ok := stacks[resource{path, r.ID}]; ok {
					stacks[resource{ref.Path, ref.SliceStackID}] = struct{}{}
				}
				if len(st.Refs) != 0 {
					err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: specerr.ErrSliceRefRef}))
				}
				if len(st.Slices) > 0 && st.Slices[0].TopZ < lastTopZ {
					err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: specerr.ErrSliceNoMonotonic}))
				}
				if len(st.Slices) > 0 {
					lastTopZ = st.Slices[len(st.Slices)-1].TopZ
				}
			} else {
				err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: specerr.ErrNonSliceStack}))
			}
		} else {
			err = append(err, specerr.NewAsset(path, i, r, &specerr.ResourcePropertyError{Index: j, Err: specerr.ErrMissingResource}))
		}
	}
	return err
}

func validateObjects(path string, res *go3mf.Resources, stacks map[resource]struct{}, err []error) (bool, []error) {
	var (
		ext         *SliceStackInfo
		mustRequire bool
	)
	for i, obj := range res.Objects {
		if ok := obj.ExtensionAttr.Get(&ext); !ok {
			continue
		}
		if ext.SliceStackID == 0 {
			err = append(err, specerr.NewObject(path, i, &specerr.MissingFieldError{Name: attrSliceRefID}))
			continue
		}
		if r, ok := res.FindAsset(ext.SliceStackID); ok {
			if r, ok := r.(*SliceStackResource); ok {
				if obj.ObjectType == go3mf.ObjectTypeModel || obj.ObjectType == go3mf.ObjectTypeSolidSupport {
					stacks[resource{path, r.ID}] = struct{}{}
				}
			} else {
				err = append(err, specerr.NewObject(path, i, specerr.ErrNonSliceStack))
			}
		} else {
			err = append(err, specerr.NewObject(path, i, specerr.ErrMissingResource))
		}
		for j, c := range obj.Components {
			if c.HasTransform() && !validTransform(c.Transform) {
				err = append(err, specerr.NewObject(path, i, &specerr.ComponentError{Index: j, Err: specerr.ErrSliceInvalidTranform}))
			}
		}
		if ext.SliceResolution == ResolutionLow {
			mustRequire = true
		}
	}
	return mustRequire, err
}

func validateChilds(model *go3mf.Model, stacks map[resource]struct{}, err []error) (bool, []error) {
	s := make([]string, 0, len(model.Childs))
	for path := range model.Childs {
		s = append(s, path)
	}
	sort.Strings(s)
	var mustRequire bool
	for _, path := range s {
		c := model.Childs[path]
		mustRequire, err = validateObjects(path, &c.Resources, stacks, err)
		err = validateAssets(path, model, &c.Resources, stacks, err)
	}
	return mustRequire, err
}
