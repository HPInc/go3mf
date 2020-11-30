package beamlattice

import (
	"strconv"

	"github.com/qmuntal/go3mf"
	specerr "github.com/qmuntal/go3mf/errors"
	"github.com/qmuntal/go3mf/spec/xml"
)

func (e Spec) NewMeshElementDecoder(mesh *go3mf.Mesh, nodeName string) xml.NodeDecoder {
	if nodeName == attrBeamLattice {
		return &beamLatticeDecoder{mesh: mesh}
	}
	return nil
}

func (e Spec) DecodeAttribute(_ interface{}, _ xml.Attr) error { return nil }

type beamLatticeDecoder struct {
	baseDecoder
	mesh *go3mf.Mesh
}

func (d *beamLatticeDecoder) Start(attrs []xml.Attr) (err error) {
	beamLattice := new(BeamLattice)
	d.mesh.Any = append(d.mesh.Any, beamLattice)
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrRadius:
			val, err1 := strconv.ParseFloat(string(a.Value), 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			beamLattice.Radius = float32(val)
		case attrMinLength, attrPrecision: // lib3mf legacy
			val, err1 := strconv.ParseFloat(string(a.Value), 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			beamLattice.MinLength = float32(val)
		case attrClippingMode, attrClipping: // lib3mf legacy
			var ok bool
			beamLattice.ClipMode, ok = newClipMode(string(a.Value))
			if !ok {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, false))
			}
		case attrClippingMesh:
			val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, false))
			}
			beamLattice.ClippingMeshID = uint32(val)
		case attrRepresentationMesh:
			val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, false))
			}
			beamLattice.RepresentationMeshID = uint32(val)
		case attrCap:
			var ok bool
			beamLattice.CapMode, ok = newCapMode(string(a.Value))
			if !ok {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, false))
			}
		}
	}
	return
}

func (d *beamLatticeDecoder) Child(name xml.Name) (child xml.NodeDecoder) {
	if name.Space == Namespace {
		if name.Local == attrBeams {
			child = &beamsDecoder{mesh: d.mesh}
		} else if name.Local == attrBeamSets {
			child = &beamSetsDecoder{mesh: d.mesh}
		}
	}
	return
}

type beamsDecoder struct {
	baseDecoder
	mesh        *go3mf.Mesh
	beamDecoder beamDecoder
}

func (d *beamsDecoder) Start(_ []xml.Attr) error {
	d.beamDecoder.mesh = d.mesh
	return nil
}

func (d *beamsDecoder) Child(name xml.Name) (child xml.NodeDecoder) {
	if name.Space == Namespace && name.Local == attrBeam {
		child = &d.beamDecoder
	}
	return
}

type beamDecoder struct {
	baseDecoder
	mesh *go3mf.Mesh
}

func (d *beamDecoder) Start(attrs []xml.Attr) (err error) {
	var beam Beam
	var (
		hasCap1, hasCap2 bool
	)
	var beamLattice *BeamLattice
	d.mesh.Any.Get(&beamLattice)
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrV1:
			val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			beam.Indices[0] = uint32(val)
		case attrV2:
			val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			beam.Indices[1] = uint32(val)
		case attrR1:
			val, err1 := strconv.ParseFloat(string(a.Value), 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, false))
			}
			beam.Radius[0] = float32(val)
		case attrR2:
			val, err1 := strconv.ParseFloat(string(a.Value), 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, false))
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
		beam.Radius[0] = beamLattice.Radius
	}
	if beam.Radius[1] == 0 {
		beam.Radius[1] = beam.Radius[0]
	}
	if !hasCap1 {
		beam.CapMode[0] = beamLattice.CapMode
	}
	if !hasCap2 {
		beam.CapMode[1] = beamLattice.CapMode
	}
	beamLattice.Beams = append(beamLattice.Beams, beam)
	return
}

type beamSetsDecoder struct {
	baseDecoder
	mesh *go3mf.Mesh
}

func (d *beamSetsDecoder) Child(name xml.Name) (child xml.NodeDecoder) {
	if name.Space == Namespace && name.Local == attrBeamSet {
		child = &beamSetDecoder{mesh: d.mesh}
	}
	return
}

type beamSetDecoder struct {
	baseDecoder
	mesh           *go3mf.Mesh
	beamSet        BeamSet
	beamRefDecoder beamRefDecoder
}

func (d *beamSetDecoder) End() {
	var beamLattice *BeamLattice
	d.mesh.Any.Get(&beamLattice)
	beamLattice.BeamSets = append(beamLattice.BeamSets, d.beamSet)
}

func (d *beamSetDecoder) Start(attrs []xml.Attr) error {
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

func (d *beamSetDecoder) Child(name xml.Name) (child xml.NodeDecoder) {
	if name.Space == Namespace && name.Local == attrRef {
		child = &d.beamRefDecoder
	}
	return
}

type beamRefDecoder struct {
	baseDecoder
	beamSet *BeamSet
}

func (d *beamRefDecoder) Start(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrIndex {
			val, err1 := strconv.ParseUint(string(a.Value), 10, 32)
			if err1 != nil {
				err = specerr.Append(err, specerr.NewParseAttrError(a.Name.Local, true))
			}
			d.beamSet.Refs = append(d.beamSet.Refs, uint32(val))
		}
	}
	return
}

type baseDecoder struct {
}

func (d *baseDecoder) Start([]xml.Attr) error { return nil }
func (d *baseDecoder) End()                   {}
