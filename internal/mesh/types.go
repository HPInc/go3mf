package mesh

import (
	"github.com/go-gl/mathgl/mgl32"
)

// MaxNodeCount is the maximum number of nodes allowed.
const MaxNodeCount = 2147483646

// MaxEdgeCount is the maximum number of edges allowed.
const MaxEdgeCount = 2147483646

// MaxFaceCount is the maximum number of faces allowed.
const MaxFaceCount = 2147483646

// MaxBeamCount is the maximum number of beams allowed.
const MaxBeamCount = 2147483646

// MaxCoordinate is the maximum value of a coordinate.
const MaxCoordinate = 1000000000.0

// Node defines a node of a mesh.
type Node struct {
	Index    uint32     // Index of the node inside the mesh.
	Position mgl32.Vec3 // Coordinates of the node.
}

// Face defines a triangle of a mesh.
type Face struct {
	Index       uint32    // Index of the face inside the mesh.
	NodeIndices [3]uint32 // Coordinates of the three nodes that defines the mesh.
}

// BeamSet defines a set of beams.
type BeamSet struct {
	Refs       []uint32 // References to all the beams in the set.
	Name       string   // Name of the set.
	Identifier string   // Identifier of the set.
}

// Beam defines a single beam.
type Beam struct {
	Index       uint32         // Index of the beam.
	NodeIndices [2]uint32      // Indices of the two nodes that defines the beam.
	Radius      [2]float64     // radius of both ends of the beam.
	CapMode     [2]BeamCapMode // Capping mode.
}

// SliceNode defines a node of an slice.
type SliceNode struct {
	Index    uint32     // Index of the slice.
	Position mgl32.Vec2 // Coordinates of the node.
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
