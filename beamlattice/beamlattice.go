package beamlattice

// ExtensionSpace is the canonical name of this extension.
const ExtensionSpace = "http://schemas.microsoft.com/3dmanufacturing/beamlattice/2017/02"

type Extension struct {
	LocalName string
}

func (e Extension) Space() string       { return ExtensionSpace }
func (e Extension) Required() bool      { return true }
func (e *Extension) SetRequired(r bool) {}
func (e *Extension) SetLocal(l string)  { e.LocalName = l }

func (e Extension) Local() string {
	if e.LocalName != "" {
		return e.LocalName
	}
	return "b"
}

// ClipMode defines the clipping modes for the beam lattices.
type ClipMode uint8

// Supported clip modes.
const (
	ClipNone ClipMode = iota
	ClipInside
	ClipOutside
)

func newClipMode(s string) (c ClipMode, ok bool) {
	c, ok = map[string]ClipMode{
		"none":    ClipNone,
		"inside":  ClipInside,
		"outside": ClipOutside,
	}[s]
	return
}

func (c ClipMode) String() string {
	return map[ClipMode]string{
		ClipNone:    "none",
		ClipInside:  "inside",
		ClipOutside: "outside",
	}[c]
}

// A CapMode is an enumerable for the different capping modes.
type CapMode uint8

// Supported cap modes.
const (
	CapModeSphere CapMode = iota
	CapModeHemisphere
	CapModeButt
)

func newCapMode(s string) (t CapMode, ok bool) {
	t, ok = map[string]CapMode{
		"sphere":     CapModeSphere,
		"hemisphere": CapModeHemisphere,
		"butt":       CapModeButt,
	}[s]
	return
}

func (b CapMode) String() string {
	return map[CapMode]string{
		CapModeSphere:     "sphere",
		CapModeHemisphere: "hemisphere",
		CapModeButt:       "butt",
	}[b]
}

// BeamLattice defines the Model Mesh BeamLattice Attributes class and is part of the BeamLattice extension to 3MF.
type BeamLattice struct {
	ClipMode                 ClipMode
	ClippingMeshID           uint32
	RepresentationMeshID     uint32
	Beams                    []Beam
	BeamSets                 []BeamSet
	MinLength, DefaultRadius float32
	CapMode                  CapMode
}

// BeamSet defines a set of beams.
type BeamSet struct {
	Refs       []uint32
	Name       string
	Identifier string
}

// Beam defines a single beam.
type Beam struct {
	NodeIndices [2]uint32  // Indices of the two nodes that defines the beam.
	Radius      [2]float32 // Radius of both ends of the beam.
	CapMode     [2]CapMode // Capping mode.
}

const (
	attrBeamLattice        = "beamlattice"
	attrRadius             = "radius"
	attrMinLength          = "minlength"
	attrPrecision          = "precision"
	attrClippingMode       = "clippingmode"
	attrClipping           = "clipping"
	attrClippingMesh       = "clippingmesh"
	attrRepresentationMesh = "representationmesh"
	attrCap                = "cap"
	attrBeams              = "beams"
	attrBeam               = "beam"
	attrBeamSets           = "beamsets"
	attrBeamSet            = "beamset"
	attrR1                 = "r1"
	attrR2                 = "r2"
	attrCap1               = "cap1"
	attrCap2               = "cap2"
	attrV1                 = "v1"
	attrV2                 = "v2"
	attrV3                 = "v3"
	attrName               = "name"
	attrIdentifier         = "identifier"
	attrRef                = "ref"
	attrIndex              = "index"
)
