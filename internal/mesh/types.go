//go:generate mockgen -destination types_mock_test.go -package mesh -self_package github.com/qmuntal/go3mf/internal/mesh github.com/qmuntal/go3mf/internal/mesh MergeableMesh

package mesh

import "github.com/qmuntal/go3mf/internal/meshinfo"

// mergeableNodes defines a structure that can be merged with another node structure.
type mergeableNodes interface {
	// NodeCount returns the number of nodes in the mesh.
	NodeCount() uint32
	// Node retrieve the node with the target index.
	Node(index uint32) *Node
}

// mergeableFaces defines a structure that can be merged with another face structure.
type mergeableFaces interface {
	// FaceCount returns the number of faces in the mesh.
	FaceCount() uint32
	// Face retrieve the face with the target index.
	Face(index uint32) *Face
	// InformationHandler returns the information handler of the mesh. Can be nil.
	InformationHandler() *meshinfo.Handler
}

// mergeableBeams defines a structure that can be merged with another beam lattice.
type mergeableBeams interface {
	// BeamCount returns the number of beams in the mesh.
	BeamCount() uint32
	// Beam retrieve the beam with the target index.
	Beam(index uint32) *Beam
}

// MergeableMesh defines a structure that can be merged with another mesh.
type MergeableMesh interface {
	mergeableNodes
	mergeableFaces
	mergeableBeams
}
