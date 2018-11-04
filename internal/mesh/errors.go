package mesh

import (
	"fmt"
)

// DuplicatedNodeError happens when an Edge with two identical nodes has been tried to add to a mesh.
type DuplicatedNodeError struct{}

func (e *DuplicatedNodeError) Error() string {
	return "an Edge with two identical nodes has been tried to add to a mesh"
}

// MaxFaceError happens when a Face has been tried to add to a mesh with too many faces.
type MaxFaceError struct {
}

func (e *MaxFaceError) Error() string {
	return fmt.Sprintf("a Face has been tried to add to a mesh with too many faces (%d)", MaxFaceCount)
}

// MaxNodeError happens when a Node has been tried to add to a mesh with too many nodes.
type MaxNodeError struct {
}

func (e *MaxNodeError) Error() string {
	return fmt.Sprintf("a Node has been tried to add to a mesh with too many nodes (%d)", MaxNodeCount)
}

// MaxBeamError happens when a Beam has been tried to add to a mesh with too many beams.
type MaxBeamError struct {
}

func (e *MaxBeamError) Error() string {
	return fmt.Sprintf("a Beam has been tried to add to a mesh with too many beams (%d)", MaxBeamCount)
}
