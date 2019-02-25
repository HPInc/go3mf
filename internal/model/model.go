package go3mf

// Metadata item is an in memory representation of the 3MF metadata,
// and can be attached to any 3MF model node.
type Metadata struct {
	Name  string // Name of the metadata.
	Value string // Value of the metadata.
}
