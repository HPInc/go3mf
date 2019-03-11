package mesh

import "errors"

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
	beams                    []Beam
	beamSets                 []BeamSet
	MinLength, DefaultRadius float64
	CapMode                  BeamCapMode
	maxBeamCount             uint32
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
	b.beams = make([]Beam, 0)
	b.beamSets = make([]BeamSet, 0)
}

// BeamCount returns the number of beams in the mesh.
func (b *beamLattice) BeamCount() uint32 {
	return uint32(len(b.beams))
}

// Beam retrieve the beam with the target index.
func (b *beamLattice) Beam(index uint32) *Beam {
	return &b.beams[uint32(index)]
}

// AddBeamSet adds a new beam set to the mesh.
func (b *beamLattice) AddBeamSet() *BeamSet {
	b.beamSets = append(b.beamSets, BeamSet{})
	return &b.beamSets[len(b.beamSets)-1]
}

// BeamSet retrieve the beam set with the target index.
func (b *beamLattice) BeamSet(index uint32) *BeamSet {
	return &b.beamSets[int(index)]
}

// AddBeam adds a beam to the mesh with the desried parameters.
func (b *beamLattice) AddBeam(node1, node2 uint32, radius1, radius2 float64, capMode1, capMode2 BeamCapMode) (*Beam, error) {
	if node1 == node2 {
		return nil, errors.New("go3mf: a beam with two identical nodes has been tried to add to a mesh")
	}

	beamCount := b.BeamCount()
	if beamCount >= b.getMaxBeamCount() {
		panic(errors.New("go3mf: too many beams has been tried to add to a mesh"))
	}

	b.beams = append(b.beams, Beam{
		Index:       beamCount,
		NodeIndices: [2]uint32{node1, node2},
		Radius:      [2]float64{radius1, radius2},
		CapMode:     [2]BeamCapMode{capMode1, capMode2},
	})
	return &b.beams[len(b.beams)-1], nil
}

func (b *beamLattice) checkSanity(nodeCount uint32) bool {
	beamCount := b.BeamCount()
	if beamCount > b.getMaxBeamCount() {
		return false
	}
	for i := uint32(0); i < beamCount; i++ {
		beam := b.Beam(i)
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

func (b *beamLattice) merge(other mergeableBeams, newNodes []uint32) error {
	beamCount := other.BeamCount()
	if beamCount == 0 {
		return nil
	}
	for i := uint32(0); i < beamCount; i++ {
		beam := other.Beam(i)
		n1, n2 := newNodes[beam.NodeIndices[0]], newNodes[beam.NodeIndices[1]]
		_, err := b.AddBeam(n1, n2, beam.Radius[0], beam.Radius[1], beam.CapMode[0], beam.CapMode[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *beamLattice) getMaxBeamCount() uint32 {
	if b.maxBeamCount == 0 {
		return maxBeamCount
	}
	return b.maxBeamCount
}
