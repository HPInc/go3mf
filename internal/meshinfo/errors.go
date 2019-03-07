package meshinfo

import (
	"fmt"
)

// FaceCountMissmatchError defines an error that happens when a mesh info operation is done with a not matching current face number.
type FaceCountMissmatchError struct {
	current, new uint32
}

func (e *FaceCountMissmatchError) Error() string {
	return fmt.Sprintf("mesh information face count (%d) does not match with mesh face count (%d)", e.current, e.new)
}
