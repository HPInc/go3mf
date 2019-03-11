package mesh

const maxBeamCount = 2147483646

// BeamSet defines a set of beams.
type BeamSet struct {
	Refs       []uint32 // References to all the beams in the set.
	Name       string   // Name of the set.
	Identifier string   // Identifier of the set.
}

// A BeamCapMode is an enumerable for the different capping modes.
type BeamCapMode int

const (
	// CapModeSphere when the capping is an sphere.
	CapModeSphere BeamCapMode = iota
	// CapModeHemisphere when the capping is an hemisphere.
	CapModeHemisphere
	// CapModeButt when the capping is an butt.
	CapModeButt
)

// Beam defines a single beam.
type Beam struct {
	Index       uint32         // Index of the beam.
	NodeIndices [2]uint32      // Indices of the two nodes that defines the beam.
	Radius      [2]float64     // Radius of both ends of the beam.
	CapMode     [2]BeamCapMode // Capping mode.
}

// beamLattice defines a beam lattice structure.
type beamLattice struct {
	Beams                    []Beam
	BeamSets                 []BeamSet
	MinLength, DefaultRadius float64
	CapMode                  BeamCapMode
	maxBeamCount             int
}

// newbeamLattice creates a new beamLattice with default values.
func newbeamLattice() *beamLattice {
	return &beamLattice{
		CapMode:       CapModeSphere,
		DefaultRadius: 1.0,
		MinLength:     0.0001,
	}
}

func (b *beamLattice) clearBeamLattice() {
	b.Beams = make([]Beam, 0)
	b.BeamSets = make([]BeamSet, 0)
}

func (b *beamLattice) checkSanity(nodeCount uint32) bool {
	if len(b.Beams) > b.getMaxBeamCount() {
		return false
	}
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

func (b *beamLattice) merge(other *beamLattice, newNodes []uint32) {
	if len(other.Beams) == 0 {
		return
	}
	for _, beam := range other.Beams {
		n1, n2 := newNodes[beam.NodeIndices[0]], newNodes[beam.NodeIndices[1]]
		b.Beams = append(b.Beams, Beam{
			Index:       uint32(len(b.Beams)),
			NodeIndices: [2]uint32{n1, n2},
			Radius:      [2]float64{beam.Radius[0], beam.Radius[1]},
			CapMode:     [2]BeamCapMode{beam.CapMode[0], beam.CapMode[1]},
		})
	}
	return
}

func (b *beamLattice) getMaxBeamCount() int {
	if b.maxBeamCount == 0 {
		return maxBeamCount
	}
	return b.maxBeamCount
}
