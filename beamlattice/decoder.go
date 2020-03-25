package beamlattice

import (
	"encoding/xml"
	"strconv"

	"github.com/qmuntal/go3mf"
)

func (e Extension) NewNodeDecoder(parentNode interface{}, nodeName string) go3mf.NodeDecoder {
	if nodeName == attrBeamLattice {
		return &beamLatticeDecoder{mesh: parentNode.(*go3mf.Mesh)}
	}
	return nil
}

func (e Extension) DecodeAttribute(_ *go3mf.Scanner, _ interface{}, _ xml.Attr) {}

type beamLatticeDecoder struct {
	baseDecoder
	mesh *go3mf.Mesh
}

func (d *beamLatticeDecoder) Start(attrs []xml.Attr) {
	beamLattice := new(BeamLattice)
	d.mesh.Any = append(d.mesh.Any, beamLattice)
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrRadius:
			val, err := strconv.ParseFloat(a.Value, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			beamLattice.DefaultRadius = float32(val)
		case attrMinLength, attrPrecision: // lib3mf legacy
			val, err := strconv.ParseFloat(a.Value, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			beamLattice.MinLength = float32(val)
		case attrClippingMode, attrClipping: // lib3mf legacy
			var ok bool
			beamLattice.ClipMode, ok = newClipMode(a.Value)
			if !ok {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
			}
		case attrClippingMesh:
			val, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
			}
			beamLattice.ClippingMeshID = uint32(val)
		case attrRepresentationMesh:
			val, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
			}
			beamLattice.RepresentationMeshID = uint32(val)
		case attrCap:
			var ok bool
			beamLattice.CapMode, ok = newCapMode(a.Value)
			if !ok {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
			}
		}
	}
}

func (d *beamLatticeDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionSpace {
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

func (d *beamsDecoder) Start(_ []xml.Attr) {
	d.beamDecoder.mesh = d.mesh
}

func (d *beamsDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionSpace && name.Local == attrBeam {
		child = &d.beamDecoder
	}
	return
}

type beamDecoder struct {
	baseDecoder
	mesh *go3mf.Mesh
}

func (d *beamDecoder) Start(attrs []xml.Attr) {
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
			val, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			beam.NodeIndices[0] = uint32(val)
		case attrV2:
			val, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			beam.NodeIndices[1] = uint32(val)
		case attrR1:
			val, err := strconv.ParseFloat(a.Value, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
			}
			beam.Radius[0] = float32(val)
		case attrR2:
			val, err := strconv.ParseFloat(a.Value, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, false)
			}
			beam.Radius[1] = float32(val)
		case attrCap1:
			var ok bool
			beam.CapMode[0], ok = newCapMode(a.Value)
			if ok {
				hasCap1 = true
			}
		case attrCap2:
			var ok bool
			beam.CapMode[1], ok = newCapMode(a.Value)
			if ok {
				hasCap2 = true
			}
		}
	}
	if beam.Radius[0] == 0 {
		beam.Radius[0] = beamLattice.DefaultRadius
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
}

type beamSetsDecoder struct {
	baseDecoder
	mesh *go3mf.Mesh
}

func (d *beamSetsDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionSpace && name.Local == attrBeamSet {
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

func (d *beamSetDecoder) Start(attrs []xml.Attr) {
	d.beamRefDecoder.beamSet = &d.beamSet
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrName:
			d.beamSet.Name = a.Value
		case attrIdentifier:
			d.beamSet.Identifier = a.Value
		}
	}
}

func (d *beamSetDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionSpace && name.Local == attrRef {
		child = &d.beamRefDecoder
	}
	return
}

type beamRefDecoder struct {
	baseDecoder
	beamSet *BeamSet
}

func (d *beamRefDecoder) Start(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrIndex {
			val, err := strconv.ParseUint(a.Value, 10, 32)
			if err != nil {
				d.Scanner.InvalidAttr(a.Name.Local, a.Value, true)
			}
			d.beamSet.Refs = append(d.beamSet.Refs, uint32(val))
		}
	}
}

type baseDecoder struct {
	Scanner *go3mf.Scanner
}

func (d *baseDecoder) Start([]xml.Attr)                 {}
func (d *baseDecoder) Text([]byte)                      {}
func (d *baseDecoder) Child(xml.Name) go3mf.NodeDecoder { return nil }
func (d *baseDecoder) End()                             {}
func (d *baseDecoder) SetScanner(s *go3mf.Scanner)      { d.Scanner = s }
