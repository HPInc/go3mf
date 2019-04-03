package io3mf

import (
	"encoding/xml"
	"errors"
	"image/color"
	"strconv"

	go3mf "github.com/qmuntal/go3mf"
)

type resourceDecoder struct {
	r             *Reader
	path          string
	progressCount int
}

func (d *resourceDecoder) Open() error {
	if !d.r.progress.progress(0.2, StageReadResources) {
		return ErrUserAborted
	}
	d.r.progress.pushLevel(0.2, 0.9)
	return nil
}

func (d *resourceDecoder) Close() error {
	d.r.progress.popLevel()
	return nil
}

func (d *resourceDecoder) Attributes(attrs []xml.Attr) error {
	return nil
}

func (d *resourceDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec {
		switch name.Local {
		case attrObject:
			d.progressCount++
			child = &objectDecoder{r: d.r, progressCount: d.progressCount, resource: go3mf.ObjectResource{ModelPath: d.path}}
		case attrBaseMaterials:
			child = &baseMaterialsDecoder{r: d.r, resource: go3mf.BaseMaterialsResource{ModelPath: d.path}}
		}
	} else if name.Space == nsMaterialSpec {
		switch name.Local {
		case attrColorGroup:
			child = &colorGroupDecoder{r: d.r, resource: go3mf.ColorGroupResource{ModelPath: d.path}}
		case attrTexture2DGroup:
			child = &tex2DGroupDecoder{r: d.r, resource: go3mf.Texture2DGroupResource{ModelPath: d.path}}
		case attrTexture2D:
			child = &texture2DDecoder{r: d.r, resource: go3mf.Texture2DResource{ModelPath: d.path}}
		case attrComposite:
			d.r.addWarning(&ReadError{InvalidOptionalValue, "go3mf: composite materials extension not supported"})
		}
	} else if name.Space == nsSliceSpec && name.Local == attrSliceStack {
		d.progressCount++
		child = &sliceStackDecoder{r: d.r, progressCount: d.progressCount, resource: go3mf.SliceStackResource{ModelPath: d.path}}
	}
	return
}

type baseMaterialsDecoder struct {
	r                   *Reader
	resource            go3mf.BaseMaterialsResource
	baseMaterialDecoder baseMaterialDecoder
}

func (d *baseMaterialsDecoder) Open() error {
	d.baseMaterialDecoder.r = d.r
	d.baseMaterialDecoder.resource = &d.resource
	return nil
}

func (d *baseMaterialsDecoder) Close() error {
	if d.resource.ID == 0 {
		return errors.New("go3mf: missing base materials resource id attribute")
	}
	d.r.addResource(&d.resource)
	return nil
}

func (d *baseMaterialsDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrBase {
		child = &d.baseMaterialDecoder
	}
	return
}

func (d *baseMaterialsDecoder) Attributes(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Space != "" || a.Name.Local != attrID {
			continue
		}
		if d.resource.ID == 0 {
			d.resource.ID, err = strconv.ParseUint(a.Value, 10, 64)
			if err != nil {
				err = errors.New("go3mf: base materials id is not valid")
			}
		} else {
			err = errors.New("go3mf: duplicated base materials id attribute")
		}
		if err != nil {
			break
		}
	}
	return
}

type baseMaterialDecoder struct {
	r        *Reader
	resource *go3mf.BaseMaterialsResource
}

func (d *baseMaterialDecoder) Open() error                                        { return nil }
func (d *baseMaterialDecoder) Close() error                                       { return nil }
func (d *baseMaterialDecoder) Child(name xml.Name) (child nodeDecoder) { return }

func (d *baseMaterialDecoder) Attributes(attrs []xml.Attr) (err error) {
	var name string
	var withColor bool
	baseColor := color.RGBA{}

	for _, a := range attrs {
		switch a.Name.Local {
		case attrName:
			name = a.Value
		case attrBaseMaterialColor:
			baseColor, err = strToSRGB(a.Value)
			if err != nil {
				return
			}
			withColor = true
		}
	}
	if name == "" || !withColor {
		return errors.New("go3mf: missing base material attributes")
	}
	d.resource.Materials = append(d.resource.Materials, go3mf.BaseMaterial{Name: name, Color: baseColor})
	return
}
