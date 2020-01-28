package beamlattice

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
)

func init() {
	go3mf.RegisterExtensionDecoder(ExtensionName, &extensionDecoder{})
}

type extensionDecoder struct{}

func (d *extensionDecoder) NodeDecoder(parentNode interface{}, nodeName string) go3mf.NodeDecoder {
	if nodeName == attrBeamLattice {
		return &beamLatticeDecoder{mesh: parentNode.(*go3mf.Mesh)}
	}
	return nil
}

func (d *extensionDecoder) DecodeAttribute(_ *go3mf.Scanner, _ interface{}, _ xml.Attr) {
}

type beamLatticeDecoder struct {
	go3mf.BaseDecoder
	mesh *go3mf.Mesh
}

func (d *beamLatticeDecoder) Attributes(attrs []xml.Attr) {
	var hasRadius, hasMinLength bool
	beamLattice := ExtensionBeamLattice(d.mesh)
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrRadius:
			beamLattice.DefaultRadius = d.Scanner.ParseFloat32Required(attrRadius, a.Value)
			hasRadius = true
		case attrMinLength, attrPrecision: // lib3mf legacy
			beamLattice.MinLength = d.Scanner.ParseFloat32Required(a.Name.Local, a.Value)
			hasMinLength = true
		case attrClippingMode, attrClipping: // lib3mf legacy
			var ok bool
			beamLattice.ClipMode, ok = newClipMode(a.Value)
			if !ok {
				d.Scanner.InvalidOptionalAttr(a.Name.Local, a.Value)
			}
		case attrClippingMesh:
			beamLattice.ClippingMeshID = d.Scanner.ParseUint32Optional(attrClippingMesh, a.Value)
		case attrRepresentationMesh:
			beamLattice.RepresentationMeshID = d.Scanner.ParseUint32Optional(attrRepresentationMesh, a.Value)
		case attrCap:
			var ok bool
			beamLattice.CapMode, ok = newCapMode(a.Value)
			if !ok {
				d.Scanner.InvalidOptionalAttr(a.Name.Local, a.Value)
			}
		}
	}
	if !hasRadius {
		d.Scanner.MissingAttr(attrRadius)
	}
	if !hasMinLength {
		d.Scanner.MissingAttr(attrMinLength)
	}
}

func (d *beamLatticeDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName {
		if name.Local == attrBeams {
			child = &beamsDecoder{mesh: d.mesh}
		} else if name.Local == attrBeamSets {
			child = &beamSetsDecoder{mesh: d.mesh}
		}
	}
	return
}

type beamsDecoder struct {
	go3mf.BaseDecoder
	mesh        *go3mf.Mesh
	beamDecoder beamDecoder
}

func (d *beamsDecoder) Open() {
	d.beamDecoder.mesh = d.mesh
}

func (d *beamsDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrBeam {
		child = &d.beamDecoder
	}
	return
}

type beamDecoder struct {
	go3mf.BaseDecoder
	mesh *go3mf.Mesh
}

func (d *beamDecoder) Attributes(attrs []xml.Attr) {
	beam := Beam{}
	var (
		hasV1, hasV2, hasCap1, hasCap2 bool
	)
	beamLattice := ExtensionBeamLattice(d.mesh)
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrV1:
			beam.NodeIndices[0] = d.Scanner.ParseUint32Required(attrV1, a.Value)
			hasV1 = true
		case attrV2:
			beam.NodeIndices[1] = d.Scanner.ParseUint32Required(attrV2, a.Value)
			hasV2 = true
		case attrR1:
			beam.Radius[0] = d.Scanner.ParseFloat32Optional(attrR1, a.Value)
		case attrR2:
			beam.Radius[1] = d.Scanner.ParseFloat32Optional(attrR2, a.Value)
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
	if !hasV1 {
		d.Scanner.MissingAttr(attrV1)
	}
	if !hasV2 {
		d.Scanner.MissingAttr(attrV2)
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
	go3mf.BaseDecoder
	mesh *go3mf.Mesh
}

func (d *beamSetsDecoder) Child(name xml.Name) (child go3mf.NodeDecoder) {
	if name.Space == ExtensionName && name.Local == attrBeamSet {
		child = &beamSetDecoder{mesh: d.mesh}
	}
	return
}

type beamSetDecoder struct {
	go3mf.BaseDecoder
	mesh           *go3mf.Mesh
	beamSet        BeamSet
	beamRefDecoder beamRefDecoder
}

func (d *beamSetDecoder) Open() {
	d.beamRefDecoder.beamSet = &d.beamSet
}

func (d *beamSetDecoder) Close() {
	ext := ExtensionBeamLattice(d.mesh)
	ext.BeamSets = append(ext.BeamSets, d.beamSet)
}

func (d *beamSetDecoder) Attributes(attrs []xml.Attr) {
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
	if name.Space == ExtensionName && name.Local == attrRef {
		child = &d.beamRefDecoder
	}
	return
}

type beamRefDecoder struct {
	go3mf.BaseDecoder
	beamSet *BeamSet
}

func (d *beamRefDecoder) Attributes(attrs []xml.Attr) {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrIndex {
			index := d.Scanner.ParseUint32Required(attrIndex, a.Value)
			d.beamSet.Refs = append(d.beamSet.Refs, uint32(index))
			return
		}
	}
	d.Scanner.MissingAttr(attrIndex)
}
