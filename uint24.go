// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package go3mf

// Uint24 represents an unsigned integer of 24 bits.
type Uint24 [3]byte

// ToUint24 returns the uint24 representation of the number.
func ToUint24(number uint32) Uint24 {
	return Uint24{byte(number), byte(number >> 8), byte(number >> 16)}
}

// ToUint32 returns the uint32 representation of v.
func (v Uint24) ToUint32() uint32 {
	return uint32(v[0]) + (uint32(v[1]) << 8) + (uint32(v[2]) << 16)
}
