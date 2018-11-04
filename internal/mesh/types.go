//go:generate mockgen -destination types_mock_test.go -package mesh -self_package github.com/qmuntal/go3mf/internal/mesh github.com/qmuntal/go3mf/internal/mesh MergeableNodes,MergeableFaces,MergeableBeams,MergeableMesh

package mesh

import "github.com/qmuntal/go3mf/internal/meshinfo"

// MergeableNodes defines a structure that can be merged with another node structure.
type MergeableNodes interface {
	// NodeCount returns the number of nodes in the mesh.
	NodeCount() uint32
	// Node retrieve the node with the target index.
	Node(index uint32) *Node
}

// MergeableFaces defines a structure that can be merged with another face structure.
type MergeableFaces interface {
	// FaceCount returns the number of faces in the mesh.
	FaceCount() uint32
	// Face retrieve the face with the target index.
	Face(index uint32) *Face
	// InformationHandler returns the information handler of the mesh. Can be nil.
	InformationHandler() *meshinfo.Handler
}

// MergeableBeams defines a structure that can be merged with another beam lattice.
type MergeableBeams interface {
	// BeamCount returns the number of beams in the mesh.
	BeamCount() uint32
	// Beam retrieve the beam with the target index.
	Beam(index uint32) *Beam
}

type MergeableMesh interface {
	MergeableNodes
	MergeableFaces
	MergeableBeams
}
