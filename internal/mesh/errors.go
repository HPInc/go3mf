package mesh

import (
	"fmt"
)

type DuplicatedNodeError struct{}

func (e *DuplicatedNodeError) Error() string {
	return "an Edge with two identical nodes has been tried to add to a mesh"
}

type MaxFaceError struct {
}

func (e *MaxFaceError) Error() string {
	return fmt.Sprintf("a Face has been tried to add to a mesh with too many faces (%d)", MaxFaceCount)
}

type MaxNodeError struct {
}

func (e *MaxNodeError) Error() string {
	return fmt.Sprintf("a Node has been tried to add to a mesh with too many nodes (%d)", MaxNodeCount)
}

type MaxBeamError struct {
}

func (e *MaxBeamError) Error() string {
	return fmt.Sprintf("a Beam has been tried to add to a mesh with too many beams (%d)", MaxBeamCount)
}
