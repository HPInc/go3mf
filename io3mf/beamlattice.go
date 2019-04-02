package io3mf

import (
	"encoding/xml"
	"strconv"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/mesh"
)

type beamLatticeDecoder struct {
	r        *Reader
	resource *go3mf.MeshResource
}

func (d *beamLatticeDecoder) Decode(x xml.TokenReader, attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsBeamLatticeSpec {
				if tp.Name.Local == attrBeams {
					db := beamsDecoder{r: d.r, mesh: d.resource.Mesh}
					db.Decode(x)
				} else if tp.Name.Local == attrBeamSets {
					db := beamSetDecoder{r: d.r, mesh: d.resource.Mesh}
					db.Decode(x, tp.Attr)
				}
			}
		case xml.EndElement:
			if tp.Name.Space == nsBeamLatticeSpec && tp.Name.Local == attrBeamLattice {
				return nil
			}
		}
	}
}

func (d *beamLatticeDecoder) parseAttr(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrRadius:
			d.resource.Mesh.DefaultRadius, err = strconv.ParseFloat(a.Value, 64)
		case attrPrecision: // lib3mf legacy
			fallthrough
		case attrMinLength:
			d.resource.Mesh.MinLength, err = strconv.ParseFloat(a.Value, 64)
		case attrClipping: // lib3mf legacy
			fallthrough
		case attrClippingMode:
			d.resource.BeamLatticeAttributes.ClipMode, _ = newClipMode(a.Value)
		case attrClippingMesh:
			var val uint64
			val, err = strconv.ParseUint(a.Value, 10, 32)
			d.resource.BeamLatticeAttributes.ClippingMeshID = uint32(val)
		case attrRepresentationMesh:
			var val uint64
			val, err = strconv.ParseUint(a.Value, 10, 32)
			d.resource.BeamLatticeAttributes.RepresentationMeshID = uint32(val)
		case attrCap:
			d.resource.Mesh.CapMode, _ = newCapMode(a.Value)
		}
		if err != nil {
			break
		}
	}
	return
}

type beamsDecoder struct {
	r    *Reader
	mesh *mesh.Mesh
}

func (d *beamsDecoder) Decode(x xml.TokenReader) error {
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsBeamLatticeSpec && tp.Name.Local == attrBeam {
				err = d.parseBeam(tp.Attr)
			}
		case xml.EndElement:
			if tp.Name.Space == nsBeamLatticeSpec && tp.Name.Local == attrBeams {
				return nil
			}
		}
	}
}

func (d *beamsDecoder) parseBeam(attrs []xml.Attr) (err error) {
	var (
		v1, v2           uint64
		r1, r2           float64
		cap1, cap2       mesh.CapMode
		hasCap1, hasCap2 bool
	)
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}

		switch a.Name.Local {
		case attrV1:
			v1, err = strconv.ParseUint(a.Value, 10, 32)
		case attrV2:
			v2, err = strconv.ParseUint(a.Value, 10, 32)
		case attrR1:
			r1, err = strconv.ParseFloat(a.Value, 64)
		case attrR2:
			r2, err = strconv.ParseFloat(a.Value, 64)
		case attrCap1:
			cap1, _ = newCapMode(a.Value)
			hasCap1 = true
		case attrCap2:
			cap2, _ = newCapMode(a.Value)
			hasCap2 = true
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		return
	}
	if r1 == 0 {
		r1 = d.mesh.DefaultRadius
	}
	if r2 == 0 {
		r2 = r1
	}
	if !hasCap1 {
		cap1 = d.mesh.CapMode
	}
	if !hasCap2 {
		cap2 = d.mesh.CapMode
	}
	d.mesh.Beams = append(d.mesh.Beams, mesh.Beam{
		NodeIndices: [2]uint32{uint32(v1), uint32(v2)},
		Radius:      [2]float64{r1, r2},
		CapMode:     [2]mesh.CapMode{cap1, cap2},
	})
	return
}

type beamSetDecoder struct {
	r    *Reader
	mesh *mesh.Mesh
}

func (d *beamSetDecoder) Decode(x xml.TokenReader, attrs []xml.Attr) error {
	beamset := mesh.BeamSet{}
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrName:
			beamset.Name = a.Value
		case attrIdentifier:
			beamset.Identifier = a.Value
		}
	}
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsBeamLatticeSpec && tp.Name.Local == attrRef {
				err = d.parseRef(&beamset, tp.Attr)
			}
		case xml.EndElement:
			if tp.Name.Space == nsBeamLatticeSpec && tp.Name.Local == attrBeamSets {
				return nil
			}
		}
	}
}

func (d *beamSetDecoder) parseRef(beamset *mesh.BeamSet, attrs []xml.Attr) (err error) {
	var index uint64
	for _, a := range attrs {
		if a.Name.Space == nsBeamLatticeSpec && a.Name.Local == attrIndex {
			index, err = strconv.ParseUint(a.Value, 10, 32)
			break
		}
	}
	if err == nil {
		beamset.Refs = append(beamset.Refs, uint32(index))
	}
	return
}
