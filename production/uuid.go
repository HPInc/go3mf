package production

import (
	"io"

	"github.com/qmuntal/go3mf/uuid"
)

// SetRand sets the random number generator to r, which implements io.Reader.
// If r.Read returns an error when the package requests random data then
// a panic will be issued.
func SetRand(r io.Reader) {
	uuid.SetRand(r)
}

// NewUUID returns a Random (Version 4) UUID.
//
// The strength of the UUIDs is based on the strength of the crypto/rand
// package.
func NewUUID() *UUID {
	u := UUID(uuid.New())
	return &u
}

// ParseUUID decodes s into a UUID or returns an error. Both the standard UUID
// forms of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx and
// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx are decoded as well as the
// Microsoft encoding {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx} and the raw hex
// encoding: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.
func ParseUUID(s string) (UUID, error) {
	if err := uuid.Validate(s); err != nil {
		return UUID(""), err
	}
	return UUID(s), nil
}
