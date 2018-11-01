package mesh

// BeamLattice defines a beam lattice structure.
type BeamLattice struct {
	Beams             []*Beam
	BeamSets          []*BeamSet
	minLength, radius float64
	capMode           BeamCapMode
}

// NewBeamLattice creates a new BeamLattice with default values.
func NewBeamLattice() *BeamLattice {
	return &BeamLattice{
		capMode:   CapModeSphere,
		radius:    1.0,
		minLength: 0.0001,
	}
}

// MinLength gets the minium length for a beam in this lattice.
func (b *BeamLattice) MinLength() float64 {
	return b.minLength
}

// Radius gets the default radius of a beam in this lattice.
func (b *BeamLattice) Radius() float64 {
	return b.radius
}

// CapMode gets the default capping mode of a beam in this lattice.
func (b *BeamLattice) CapMode() BeamCapMode {
	return b.capMode
}

// SetMinLength sets the minimum length of a beam in this lattice.
func (b *BeamLattice) SetMinLength(val float64) {
	b.minLength = val
}

// SetRadius sets the default radius of a beam in this lattice.
func (b *BeamLattice) SetRadius(val float64) {
	b.radius = val
}

// SetCapMode sets the default capping mode of a beam in this lattice.
func (b *BeamLattice) SetCapMode(val BeamCapMode) {
	b.capMode = val
}

// Clear resets the value of Beams and BeamSets.
func (b *BeamLattice) Clear() {
	b.Beams = make([]*Beam, 0)
	b.BeamSets = make([]*BeamSet, 0)
}
