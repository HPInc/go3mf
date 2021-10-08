// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package beamlattice

import (
	"encoding/xml"
	"strconv"

	"github.com/hpinc/go3mf"
	specerr "github.com/hpinc/go3mf/errors"
	"github.com/hpinc/go3mf/spec"
)

func (Spec) NewAttrGroup(xml.Name) spec.AttrGroup {
	return nil
}

func (Spec) NewElementDecoder(parent interface{}, name string) spec.ElementDecoder {
	if name == attrBeamLattice {
		return &beamLatticeDecoder{mesh: parent.(*go3mf.Mesh)}
	}
	return nil
}

type beamLatticeDecoder struct {
	baseDecoder
	mesh        *go3mf.Mesh
	beamLattice *BeamLattice
}

func (d *beamLatticeDecoder) Start(attrs []spec.XMLAttr) error {
	var errs error
	d.beamLattice = new(BeamLattice)
	d.mesh.Any = append(d.mesh.Any, d.beamLattice)
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrRadius:
			val, err := strconv.ParseFloat(string(a.Value), 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.beamLattice.Radius = float32(val)
		case attrMinLength, attrPrecision: // lib3mf legacy
			val, err := strconv.ParseFloat(string(a.Value), 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.beamLattice.MinLength = float32(val)
		case attrClippingMode, attrClipping: // lib3mf legacy
			var ok bool
			d.beamLattice.ClipMode, ok = newClipMode(string(a.Value))
			if !ok {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
			}
		case attrClippingMesh:
			val, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
			}
			d.beamLattice.ClippingMeshID = uint32(val)
		case attrRepresentationMesh:
			val, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
			}
			d.beamLattice.RepresentationMeshID = uint32(val)
		case attrCap:
			var ok bool
			d.beamLattice.CapMode, ok = newCapMode(string(a.Value))
			if !ok {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
			}
		}
	}
	if errs != nil {
		return specerr.Wrap(errs, d.beamLattice)
	}
	return nil
}

func (d *beamLatticeDecoder) WrapError(err error) error {
	return specerr.Wrap(err, d.beamLattice)
}

func (d *beamLatticeDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace {
		if name.Local == attrBeams {
			child = &beamsDecoder{beamLattice: d.beamLattice}
		} else if name.Local == attrBeamSets {
			child = &beamSetsDecoder{beamLattice: d.beamLattice}
		}
	}
	return
}

type beamsDecoder struct {
	baseDecoder
	beamLattice *BeamLattice
	beamDecoder beamDecoder
}

func (d *beamsDecoder) Start(_ []spec.XMLAttr) error {
	d.beamDecoder.beamLattice = d.beamLattice
	return nil
}

func (d *beamsDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrBeam {
		child = &d.beamDecoder
	}
	return
}

func (d *beamsDecoder) WrapError(err error) error {
	return specerr.Wrap(err, &d.beamLattice.Beams)
}

type beamDecoder struct {
	baseDecoder
	beamLattice *BeamLattice
}

func (d *beamDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		beam             Beam
		hasCap1, hasCap2 bool
		errs             error
	)
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrV1:
			val, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			beam.Indices[0] = uint32(val)
		case attrV2:
			val, err := strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			beam.Indices[1] = uint32(val)
		case attrR1:
			val, err := strconv.ParseFloat(string(a.Value), 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
			}
			beam.Radius[0] = float32(val)
		case attrR2:
			val, err := strconv.ParseFloat(string(a.Value), 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, false))
			}
			beam.Radius[1] = float32(val)
		case attrCap1:
			var ok bool
			beam.CapMode[0], ok = newCapMode(string(a.Value))
			if ok {
				hasCap1 = true
			}
		case attrCap2:
			var ok bool
			beam.CapMode[1], ok = newCapMode(string(a.Value))
			if ok {
				hasCap2 = true
			}
		}
	}
	if beam.Radius[0] == 0 {
		beam.Radius[0] = d.beamLattice.Radius
	}
	if beam.Radius[1] == 0 {
		beam.Radius[1] = beam.Radius[0]
	}
	if !hasCap1 {
		beam.CapMode[0] = d.beamLattice.CapMode
	}
	if !hasCap2 {
		beam.CapMode[1] = d.beamLattice.CapMode
	}
	d.beamLattice.Beams.Beam = append(d.beamLattice.Beams.Beam, beam)
	if errs != nil {
		return specerr.WrapIndex(errs, beam, len(d.beamLattice.Beams.Beam)-1)
	}
	return nil
}

type beamSetsDecoder struct {
	baseDecoder
	beamLattice *BeamLattice
}

func (d *beamSetsDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrBeamSet {
		child = &beamSetDecoder{beamLattice: d.beamLattice}
	}
	return
}

func (d *beamSetsDecoder) WrapError(err error) error {
	return specerr.Wrap(err, &d.beamLattice.BeamSets)
}

type beamSetDecoder struct {
	baseDecoder
	beamLattice    *BeamLattice
	beamSet        BeamSet
	beamRefDecoder beamRefDecoder
}

func (d *beamSetDecoder) End() {
	d.beamLattice.BeamSets.BeamSet = append(d.beamLattice.BeamSets.BeamSet, d.beamSet)
}

func (d *beamSetDecoder) Start(attrs []spec.XMLAttr) error {
	d.beamRefDecoder.beamSet = &d.beamSet
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrName:
			d.beamSet.Name = string(a.Value)
		case attrIdentifier:
			d.beamSet.Identifier = string(a.Value)
		}
	}
	return nil
}

func (d *beamSetDecoder) WrapError(err error) error {
	return specerr.WrapIndex(err, &d.beamSet, len(d.beamLattice.BeamSets.BeamSet))
}

func (d *beamSetDecoder) Child(name xml.Name) (child spec.ElementDecoder) {
	if name.Space == Namespace && name.Local == attrRef {
		child = &d.beamRefDecoder
	}
	return
}

type beamRefDecoder struct {
	baseDecoder
	beamSet *BeamSet
}

func (d *beamRefDecoder) Start(attrs []spec.XMLAttr) error {
	var (
		val  uint64
		errs error
	)
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrIndex {
			var err error
			val, err = strconv.ParseUint(string(a.Value), 10, 32)
			if err != nil {
				errs = specerr.Append(errs, specerr.NewParseAttrError(a.Name.Local, true))
			}
			break
		}
	}
	d.beamSet.Refs = append(d.beamSet.Refs, uint32(val))
	if errs != nil {
		return specerr.WrapIndex(errs, uint32(0), len(d.beamSet.Refs)-1)
	}
	return nil
}

type baseDecoder struct {
}

func (d *baseDecoder) Start([]spec.XMLAttr) error { return nil }
func (d *baseDecoder) End()                       {}
