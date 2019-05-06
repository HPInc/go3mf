package geo

// BeamSet defines a set of beams.
type BeamSet struct {
	Refs       []uint32
	Name       string
	Identifier string
}

// A CapMode is an enumerable for the different capping modes.
type CapMode int

const (
	// CapModeSphere when the capping is an sphere.
	CapModeSphere CapMode = iota
	// CapModeHemisphere when the capping is an hemisphere.
	CapModeHemisphere
	// CapModeButt when the capping is an butt.
	CapModeButt
)

func (b CapMode) String() string {
	return map[CapMode]string{
		CapModeSphere:     "sphere",
		CapModeHemisphere: "hemisphere",
		CapModeButt:       "butt",
	}[b]
}

// Beam defines a single beam.
type Beam struct {
	NodeIndices [2]uint32  // Indices of the two nodes that defines the beam.
	Radius      [2]float64 // Radius of both ends of the beam.
	CapMode     [2]CapMode // Capping mode.
}

// beamLattice defines a beam lattice structure.
type beamLattice struct {
	Beams                    []Beam
	BeamSets                 []BeamSet
	MinLength, DefaultRadius float64
	CapMode                  CapMode
}

func (b *beamLattice) checkSanity(nodeCount uint32) bool {
	for _, beam := range b.Beams {
		i0, i1 := beam.NodeIndices[0], beam.NodeIndices[1]
		if i0 == i1 {
			return false
		}
		if i0 >= nodeCount || i1 >= nodeCount {
			return false
		}
	}
	return true
}
