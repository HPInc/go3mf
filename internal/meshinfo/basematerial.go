package meshinfo

// BaseMaterial informs about a base material.
type BaseMaterial struct {
	MaterialGroupID uint32 // Identifier of the group.
	MaterialIndex   uint32 // Index of the base material used in the group.
}

func (b *BaseMaterial) Invalidate() {
	b.MaterialGroupID = 0
	b.MaterialIndex = 0
}

func (b *BaseMaterial) Copy(from interface{}) {
	other, ok := from.(*BaseMaterial)
	if !ok {
		return
	}
	b.MaterialGroupID = other.MaterialGroupID
	b.MaterialIndex = other.MaterialIndex
}

func (b *BaseMaterial) HasData() bool {
	return b.MaterialGroupID != 0
}

func (b *BaseMaterial) Permute(index1, index2, index3 uint32) {
	// nothing to permute
}

func (b *BaseMaterial) Merge(other interface{}) {
	// nothing to merge
}
