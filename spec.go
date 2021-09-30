// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package go3mf

import (
	"github.com/hpinc/go3mf/spec"
)

type objectPather interface {
	ObjectPath() string
}

// UnknownAsset wraps a spec.UnknownTokens to fulfill
// the Asset interface.
type UnknownAsset struct {
	spec.UnknownTokens
	id uint32
}

func (u UnknownAsset) Identify() uint32 {
	return u.id
}
