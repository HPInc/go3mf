package internal

import "strconv"

type Scanner interface {
	InvalidAttr(string,string)
	InvalidAttrOptional(string,string)
}

// ParseUint32 parses s as a uint32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func ParseUint32(p Scanner, attr string, s string) uint32 {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		p.InvalidAttr(attr, s)
		return 0
	}
	return uint32(n)
}

// ParseUint32Optional parses s as a uint32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func ParseUint32Optional(p Scanner, attr string, s string) uint32 {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		p.InvalidAttrOptional(attr, s)
	}
	return uint32(n)
}

// ParseFloat32 parses s as a float32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func ParseFloat32(p Scanner, attr string, s string) float32 {
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		p.InvalidAttr(attr, s)
		return 0
	}
	return float32(n)
}

// ParseFloat32Optional parses s as a float32.
// If it cannot be parsed a ParsePropertyError is added to the warnings.
func ParseFloat32Optional(p Scanner, attr string, s string) float32 {
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		p.InvalidAttrOptional(attr, s)
	}
	return float32(n)
}
