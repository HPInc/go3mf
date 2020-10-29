package beamlattice

import (
	"syscall/js"

	"github.com/qmuntal/go3mf"
)

var (
	jsNS             = "BEAMLATTICE"
	arrayConstructor = js.Global().Get("Array")
	jsSpec           = go3mf.RegisterClass(jsNS, "X")
	jsBeamLattice    = go3mf.RegisterClass("BeamLattice", "X", jsNS)
	jsBeam           = go3mf.RegisterClass("Beam", "X", jsNS)
	jsBeamSet        = go3mf.RegisterClass("BeamSet", "X", jsNS)
)

// JSValue returns a JavaScript value associated with the object.
func (b Beam) JSValue() js.Value {
	v := jsBeam.New()
	v.Set(attrV1, b.Indices[0])
	v.Set(attrV2, b.Indices[1])
	if b.Radius[0] != 0 {
		v.Set(attrR1, b.Radius[0])
		v.Set(attrR2, b.Radius[1])
	} else {
		v.Set(attrR1, js.Undefined())
		v.Set(attrR2, js.Undefined())
	}
	v.Set(attrCap1, b.CapMode[0].String())
	v.Set(attrCap1, b.CapMode[1].String())
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (bs BeamSet) JSValue() js.Value {
	v := jsBeamSet.New()
	setString(v, attrName, bs.Name)
	setString(v, attrIdentifier, bs.Identifier)
	sr := arrayConstructor.New(len(bs.Refs))
	for i, b := range bs.Refs {
		sr.SetIndex(i, b)
	}
	v.Set(attrRef, sr)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (bl *BeamLattice) JSValue() js.Value {
	v := jsBeamLattice.New()
	v.Set("minLength", bl.MinLength)
	v.Set(attrRadius, bl.Radius)
	v.Set("clippingMode", bl.ClipMode.String())
	if bl.ClippingMeshID != 0 {
		v.Set("clippingMesh", bl.ClippingMeshID)
	} else {
		v.Set("clippingMesh", js.Undefined())
	}
	if bl.RepresentationMeshID != 0 {
		v.Set("representationMesh", bl.RepresentationMeshID)
	} else {
		v.Set("representationMesh", js.Undefined())
	}
	v.Set(attrCap, bl.CapMode.String())

	sb := arrayConstructor.New(len(bl.Beams))
	for i, b := range bl.Beams {
		sb.SetIndex(i, b)
	}
	v.Set(attrBeams, sb)

	sbs := arrayConstructor.New(len(bl.BeamSets))
	for i, bs := range bl.BeamSets {
		sbs.SetIndex(i, bs)
	}
	v.Set(attrBeamSets, sbs)
	return v
}

func setString(v js.Value, name, value string) {
	if value == "" {
		v.Set(name, js.Undefined())
	} else {
		v.Set(name, value)
	}
}
