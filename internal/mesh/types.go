package mesh

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Node defines a node of a mesh.
type Node struct {
	Index    int32      // Index of the node inside the mesh.
	Position mgl32.Vec3 // Coordinates of the node.
}

// Face defines a triangle of a mesh.
type Face struct {
	Index       int32    // Index of the face inside the mesh.
	NodeIndices [3]int32 // Coordinates of the three nodes that defines the mesh.
}

// BeamSet defines a set of beams.
type BeamSet struct {
	Refs       []uint32 // References to all the beams in the set.
	Name       string   // Name of the set.
	Identifier string   // Identifier of the set.
}

// Beam defines a single beam.
type Beam struct {
	Index       int32    // Index of the beam.
	NodeIndices [2]int32 // Indices of the two nodes that defines the beam.
	Radius      float64  // radius of the beam.
	CapMode     int32    // Capping mode.
}

// SliceNode defines a node of an slice.
type SliceNode struct {
	Index    int32      // Index of the slice.
	Position mgl32.Vec2 // Coordinates of the node.
}

// A BeamLatticeCapMode is an enumerable for the different capping modes.
type BeamLatticeCapMode int

const (
	// CapModeSphere when the capping is an sphere.
	CapModeSphere BeamLatticeCapMode = iota
	// CapModeHemisphere when the capping is an hemisphere.
	CapModeHemisphere
	// CapModeButt when the capping is an butt.
	CapModeButt
)
