package mesh

// MaxBeamCount is the maximum number of beams allowed.
const MaxBeamCount = 2147483646

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
	beams             []*Beam
	beamSets          []*BeamSet
	minLength, radius float64
	capMode           BeamCapMode
	maxBeamCount      uint32 // If 0 MaxBeamCount will be used.
}

// newbeamLattice creates a new beamLattice with default values.
func newbeamLattice() *beamLattice {
	return &beamLattice{
		capMode:   CapModeSphere,
		radius:    1.0,
		minLength: 0.0001,
	}
}

// BeamLatticeMinLength gets the minium length for a beam in this lattice.
func (b *beamLattice) BeamLatticeMinLength() float64 {
	return b.minLength
}

// DefaultBeamLatticeRadius gets the default radius of a beam in this lattice.
func (b *beamLattice) DefaultBeamLatticeRadius() float64 {
	return b.radius
}

// BeamLatticeCapMode gets the default capping mode of a beam in this lattice.
func (b *beamLattice) BeamLatticeCapMode() BeamCapMode {
	return b.capMode
}

// SetBeamLatticeMinLength sets the minimum length of a beam in this lattice.
func (b *beamLattice) SetBeamLatticeMinLength(val float64) {
	b.minLength = val
}

// SetDefaultBeamRadius sets the default radius of a beam in this lattice.
func (b *beamLattice) SetDefaultBeamRadius(val float64) {
	b.radius = val
}

// SetBeamLatticeCapMode sets the default capping mode of a beam in this lattice.
func (b *beamLattice) SetBeamLatticeCapMode(val BeamCapMode) {
	b.capMode = val
}

// ClearBeamLattice resets the value of Beams and BeamSets.
func (b *beamLattice) ClearBeamLattice() {
	b.beams = make([]*Beam, 0)
	b.beamSets = make([]*BeamSet, 0)
}

// BeamCount returns the number of beams in the mesh.
func (b *beamLattice) BeamCount() uint32 {
	return uint32(len(b.beams))
}

// Beam retrieve the beam with the target index.
func (b *beamLattice) Beam(index uint32) *Beam {
	return b.beams[uint32(index)]
}

// AddBeamSet adds a new beam set to the mesh.
func (b *beamLattice) AddBeamSet() *BeamSet {
	set := new(BeamSet)
	b.beamSets = append(b.beamSets, set)
	return set
}

// BeamSet retrieve the beam set with the target index.
func (b *beamLattice) BeamSet(index uint32) *BeamSet {
	return b.beamSets[int(index)]
}

// AddBeam adds a beam to the mesh with the desried parameters.
func (b *beamLattice) AddBeam(node1, node2 *Node, radius1, radius2 float64, capMode1, capMode2 BeamCapMode) (*Beam, error) {
	if node1 == node2 {
		return nil, new(DuplicatedNodeError)
	}

	beamCount := b.BeamCount()
	if beamCount >= b.getMaxBeamCount() {
		panic(new(MaxBeamError))
	}

	beam := &Beam{
		Index:       beamCount,
		NodeIndices: [2]uint32{node1.Index, node2.Index},
		Radius:      [2]float64{radius1, radius2},
		CapMode:     [2]BeamCapMode{capMode1, capMode2},
	}

	b.beams = append(b.beams, beam)
	return beam, nil
}

func (b *beamLattice) checkSanity(nodeCount uint32) bool {
	beamCount := b.BeamCount()
	if beamCount > b.getMaxBeamCount() {
		return false
	}
	for i := 0; i < int(beamCount); i++ {
		beam := b.Beam(uint32(i))
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

func (b *beamLattice) merge(other mergeableBeams, newNodes []*Node) error {
	beamCount := other.BeamCount()
	if beamCount == 0 {
		return nil
	}
	for i := 0; i < int(beamCount); i++ {
		beam := other.Beam(uint32(i))
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
		return MaxBeamCount
	}
	return b.maxBeamCount
}
