package meshinfo

// BaseMaterial informs about a base material.
type BaseMaterial struct {
	GroupID uint32 // Identifier of the group.
	Index   uint32 // Index of the base material used in the group.
}

// Invalidate sets to zero all the properties.
func (b *BaseMaterial) Invalidate() {
	b.GroupID = 0
	b.Index = 0
}

// Copy copy the properties of another base material.
func (b *BaseMaterial) Copy(from FaceData) {
	other, ok := from.(*BaseMaterial)
	if !ok {
		return
	}
	b.GroupID = other.GroupID
	b.Index = other.Index
}

// HasData returns true if the group id is different from zero.
func (b BaseMaterial) HasData() bool {
	return b.GroupID != 0
}

// Permute is not necessary for a base material.
func (b *BaseMaterial) Permute(index1, index2, index3 uint32) {
	// nothing to permute
}

// Merge is not necessary for a base material.
func (b *BaseMaterial) Merge(other FaceData) {
	// nothing to merge
}

type baseMaterialContainer struct {
	dataBlocks []BaseMaterial
}

func newbaseMaterialContainer(currentFaceCount uint32) *baseMaterialContainer {
	return &baseMaterialContainer{
		dataBlocks: make([]BaseMaterial, currentFaceCount),
	}
}

func (m *baseMaterialContainer) clone(currentFaceCount uint32) Container {
	return newbaseMaterialContainer(currentFaceCount)
}

func (m *baseMaterialContainer) InfoType() DataType {
	return BaseMaterialType
}

func (m *baseMaterialContainer) AddFaceData(newFaceCount uint32) FaceData {
	m.dataBlocks = append(m.dataBlocks, BaseMaterial{})
	if len(m.dataBlocks) != int(newFaceCount) {
		panic(errFaceCountMissmatch)
	}
	return &m.dataBlocks[newFaceCount-1]
}

func (m *baseMaterialContainer) FaceData(faceIndex uint32) FaceData {
	return &m.dataBlocks[faceIndex]
}

func (m *baseMaterialContainer) FaceCount() uint32 {
	return uint32(len(m.dataBlocks))
}

func (m *baseMaterialContainer) Clear() {
	m.dataBlocks = m.dataBlocks[:0]
}
