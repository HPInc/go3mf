// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package materials

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/hpinc/go3mf"
	specerr "github.com/hpinc/go3mf/errors"
	"github.com/hpinc/go3mf/spec"
)

func (Spec) NewAttr3MF(string) spec.Attr3MF {
	return nil
}

func (Spec) NewElementDecoder(parent interface{}, name string) (child spec.ElementDecoder) {
	switch name {
	case attrColorGroup:
		child = &colorGroupDecoder{resources: parent.(*go3mf.Resources)}
	case attrTexture2DGroup:
		child = &tex2DGroupDecoder{resources: parent.(*go3mf.Resources)}
	case attrTexture2D:
		child = &texture2DDecoder{resources: parent.(*go3mf.Resources)}
	case attrCompositematerials:
		child = &compositeMaterialsDecoder{resources: parent.(*go3mf.Resources)}
	case attrMultiProps:
		child = &multiPropertiesDecoder{resources: parent.(*go3mf.Resources)}
	}
	return
}

type colorGroupDecoder struct {
	baseDecoder
	resources    *go3mf.Resources
	resource     ColorGroup
	colorDecoder colorDecoder
}

func (d *colorGroupDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *colorGroupDecoder) Wrap(err error) error {
	return specerr.WrapIndex(err, &d.resource, len(d.resources.Assets))
}

func (d *colorGroupDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrColor {
		child = &d.colorDecoder
	}
	return
}

func (d *colorGroupDecoder) Start(attrs []spec.XMLAttr) (errs error) {
	d.colorDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrID {
			id, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
			break
		}
	}
	return
}

type colorDecoder struct {
	baseDecoder
	resource *ColorGroup
}

func (d *colorDecoder) Start(attrs []spec.XMLAttr) error {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrColor {
			c, err := spec.ParseRGBA(string(a.Value))
			if err != nil {
				err = specerr.NewParseAttrError(a.Name.Local, true)
			}
			d.resource.Colors = append(d.resource.Colors, c)
			if err != nil {
				return specerr.WrapIndex(err, c, len(d.resource.Colors)-1)
			}
			break
		}
	}
	return nil
}

type tex2DCoordDecoder struct {
	baseDecoder
	resource *Texture2DGroup
}

func (d *tex2DCoordDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		text TextureCoord
		errs error
	)
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		val, err := strconv.ParseFloat(string(a.Value), 32)
		if err != nil {
			errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
		}
		switch a.Name.Local {
		case attrU:
			text[0] = float32(val)
		case attrV:
			text[1] = float32(val)
		}
	}
	d.resource.Coords = append(d.resource.Coords, text)
	if errs != nil {
		return specerr.WrapIndex(errs, text, len(d.resource.Coords)-1)
	}
	return nil
}

type tex2DGroupDecoder struct {
	baseDecoder
	resources         *go3mf.Resources
	resource          Texture2DGroup
	tex2DCoordDecoder tex2DCoordDecoder
}

func (d *tex2DGroupDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *tex2DGroupDecoder) Wrap(err error) error {
	return specerr.WrapIndex(err, &d.resource, len(d.resources.Assets))
}

func (d *tex2DGroupDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrTex2DCoord {
		child = &d.tex2DCoordDecoder
	}
	return
}

func (d *tex2DGroupDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	d.tex2DCoordDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
		case attrTexID:
			val, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.TextureID = uint32(val)
		}
	}
	if errs != nil {
		return specerr.WrapIndex(errs, &d.resource, len(d.resources.Assets))
	}
	return nil
}

type texture2DDecoder struct {
	baseDecoder
	resources *go3mf.Resources
	resource  Texture2D
}

func (d *texture2DDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *texture2DDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
		case attrPath:
			d.resource.Path = string(a.Value)
		case attrContentType:
			d.resource.ContentType, _ = newTexture2DType(string(a.Value))
		case attrTileStyleU:
			d.resource.TileStyleU, _ = newTileStyle(string(a.Value))
		case attrTileStyleV:
			d.resource.TileStyleV, _ = newTileStyle(string(a.Value))
		case attrFilter:
			d.resource.Filter, _ = newTextureFilter(string(a.Value))
		}
	}
	if errs != nil {
		return specerr.WrapIndex(errs, &d.resource, len(d.resources.Assets))
	}
	return nil
}

type compositeMaterialsDecoder struct {
	baseDecoder
	resources        *go3mf.Resources
	resource         CompositeMaterials
	compositeDecoder compositeDecoder
}

func (d *compositeMaterialsDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *compositeMaterialsDecoder) Wrap(err error) error {
	return specerr.WrapIndex(err, &d.resource, len(d.resources.Assets))
}

func (d *compositeMaterialsDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrComposite {
		child = &d.compositeDecoder
	}
	return
}

func (d *compositeMaterialsDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	d.compositeDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
		case attrMatID:
			val, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.MaterialID = uint32(val)
		case attrMatIndices:
			for _, f := range strings.Fields(string(a.Value)) {
				val, err := strconv.ParseUint(f, 10, 32)
				if err != nil {
					errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
				}
				d.resource.Indices = append(d.resource.Indices, uint32(val))
			}
		}
	}
	if errs != nil {
		return specerr.WrapIndex(errs, &d.resource, len(d.resources.Assets))
	}
	return nil
}

type compositeDecoder struct {
	baseDecoder
	resource *CompositeMaterials
}

func (d *compositeDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		composite Composite
		errs      error
	)
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrValues {
			for _, f := range strings.Fields(string(a.Value)) {
				val, err := strconv.ParseFloat(f, 32)
				if err != nil {
					errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
				}
				composite.Values = append(composite.Values, float32(val))
			}
		}
	}
	d.resource.Composites = append(d.resource.Composites, composite)
	if errs != nil {
		return specerr.WrapIndex(errs, composite, len(d.resource.Composites)-1)
	}
	return nil
}

type multiPropertiesDecoder struct {
	baseDecoder
	resources    *go3mf.Resources
	resource     MultiProperties
	multiDecoder multiDecoder
}

func (d *multiPropertiesDecoder) End() {
	d.resources.Assets = append(d.resources.Assets, &d.resource)
}

func (d *multiPropertiesDecoder) Wrap(err error) error {
	return specerr.WrapIndex(err, &d.resource, len(d.resources.Assets))
}

func (d *multiPropertiesDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrMulti {
		child = &d.multiDecoder
	}
	return
}

func (d *multiPropertiesDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	d.multiDecoder.resource = &d.resource
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrID:
			id, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.resource.ID = uint32(id)
		case attrBlendMethods:
			for _, f := range strings.Fields(string(a.Value)) {
				val, _ := newBlendMethod(f)
				d.resource.BlendMethods = append(d.resource.BlendMethods, val)
			}
		case attrPIDs:
			for _, f := range strings.Fields(string(a.Value)) {
				val, err := strconv.ParseUint(f, 10, 32)
				if err != nil {
					errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
				}
				d.resource.PIDs = append(d.resource.PIDs, uint32(val))
			}
		}
	}
	if errs != nil {
		return specerr.WrapIndex(errs, &d.resource, len(d.resources.Assets))
	}
	return nil
}

type multiDecoder struct {
	baseDecoder
	resource *MultiProperties
}

func (d *multiDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		multi Multi
		errs  error
	)
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrPIndices {
			for _, f := range strings.Fields(string(a.Value)) {
				val, err := strconv.ParseUint(f, 10, 32)
				if err != nil {
					errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
				}
				multi.PIndices = append(multi.PIndices, uint32(val))
			}
		}
	}
	d.resource.Multis = append(d.resource.Multis, multi)
	if errs != nil {
		return specerr.WrapIndex(errs, &d.resource, len(d.resource.Multis)-1)
	}
	return nil
}

type baseDecoder struct {
}

func (d *baseDecoder) End() {}
