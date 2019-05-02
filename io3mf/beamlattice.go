package io3mf

import (
	"encoding/xml"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/mesh"
)

type beamLatticeDecoder struct {
	emptyDecoder
	resource *go3mf.MeshResource
}

func (d *beamLatticeDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	var hasRadius, hasMinLength bool
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrRadius:
			d.resource.Mesh.DefaultRadius, ok = d.file.parser.ParseFloat64Required(attrRadius, a.Value)
			hasRadius = true
		case attrMinLength, attrPrecision: // lib3mf legacy
			d.resource.Mesh.MinLength, ok = d.file.parser.ParseFloat64Required(a.Name.Local, a.Value)
			hasMinLength = true
		case attrClippingMode, attrClipping: // lib3mf legacy
			d.resource.BeamLatticeAttributes.ClipMode, _ = newClipMode(a.Value)
		case attrClippingMesh:
			d.resource.BeamLatticeAttributes.ClippingMeshID = d.file.parser.ParseUint32Optional(attrClippingMesh, a.Value)
		case attrRepresentationMesh:
			d.resource.BeamLatticeAttributes.RepresentationMeshID = d.file.parser.ParseUint32Optional(attrRepresentationMesh, a.Value)
		case attrCap:
			d.resource.Mesh.CapMode, _ = newCapMode(a.Value)
		}
		if !ok {
			return false
		}
	}
	if !hasRadius {
		ok = d.file.parser.MissingAttr(attrRadius)
	}
	if !hasMinLength {
		ok = d.file.parser.MissingAttr(attrMinLength)
	}
	return ok
}

func (d *beamLatticeDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsBeamLatticeSpec {
		if name.Local == attrBeams {
			child = &beamsDecoder{mesh: d.resource.Mesh}
		} else if name.Local == attrBeamSets {
			child = &beamSetsDecoder{mesh: d.resource.Mesh}
		}
	}
	return
}

type beamsDecoder struct {
	emptyDecoder
	mesh        *mesh.Mesh
	beamDecoder beamDecoder
}

func (d *beamsDecoder) Open() {
	d.beamDecoder.mesh = d.mesh
}

func (d *beamsDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsBeamLatticeSpec && name.Local == attrBeam {
		child = &d.beamDecoder
	}
	return
}

type beamDecoder struct {
	emptyDecoder
	mesh *mesh.Mesh
}

func (d *beamDecoder) Attributes(attrs []xml.Attr) bool {
	beam := mesh.Beam{}
	var (
		hasV1, hasV2, hasCap1, hasCap2 bool
	)
	ok := true
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrV1:
			beam.NodeIndices[0], ok = d.file.parser.ParseUint32Required(attrV1, a.Value)
			hasV1 = true
		case attrV2:
			beam.NodeIndices[1], ok = d.file.parser.ParseUint32Required(attrV2, a.Value)
			hasV2 = true
		case attrR1:
			beam.Radius[0] = d.file.parser.ParseFloat64Optional(attrR1, a.Value)
		case attrR2:
			beam.Radius[1] = d.file.parser.ParseFloat64Optional(attrR2, a.Value)
		case attrCap1:
			beam.CapMode[0], _ = newCapMode(a.Value)
			hasCap1 = true
		case attrCap2:
			beam.CapMode[1], _ = newCapMode(a.Value)
			hasCap2 = true
		}
		if !ok {
			return false
		}
	}
	if !hasV1 {
		ok = d.file.parser.MissingAttr(attrV1)
	}
	if !hasV2 {
		ok = d.file.parser.MissingAttr(attrV2)
	}
	if ok {
		if beam.Radius[0] == 0 {
			beam.Radius[0] = d.mesh.DefaultRadius
		}
		if beam.Radius[1] == 0 {
			beam.Radius[1] = beam.Radius[0]
		}
		if !hasCap1 {
			beam.CapMode[0] = d.mesh.CapMode
		}
		if !hasCap2 {
			beam.CapMode[1] = d.mesh.CapMode
		}
		d.mesh.Beams = append(d.mesh.Beams, beam)
	}
	return ok
}

type beamSetsDecoder struct {
	emptyDecoder
	mesh *mesh.Mesh
}

func (d *beamSetsDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsBeamLatticeSpec && name.Local == attrBeamSet {
		child = &beamSetDecoder{mesh: d.mesh}
	}
	return
}

type beamSetDecoder struct {
	emptyDecoder
	mesh           *mesh.Mesh
	beamSet        mesh.BeamSet
	beamRefDecoder beamRefDecoder
}

func (d *beamSetDecoder) Open() {
	d.beamRefDecoder.beamSet = &d.beamSet
}

func (d *beamSetDecoder) Close() bool {
	d.mesh.BeamSets = append(d.mesh.BeamSets, d.beamSet)
	return true
}

func (d *beamSetDecoder) Attributes(attrs []xml.Attr) bool {
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
	return true
}

func (d *beamSetDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsBeamLatticeSpec && name.Local == attrRef {
		child = &d.beamRefDecoder
	}
	return
}

type beamRefDecoder struct {
	emptyDecoder
	beamSet *mesh.BeamSet
}

func (d *beamRefDecoder) Attributes(attrs []xml.Attr) bool {
	for _, a := range attrs {
		if a.Name.Space == "" && a.Name.Local == attrIndex {
			index, ok := d.file.parser.ParseUint32Required(attrIndex, a.Value)
			if ok {
				d.beamSet.Refs = append(d.beamSet.Refs, uint32(index))
				return true
			}
			break
		}
	}
	return d.file.parser.MissingAttr(attrIndex)
}
