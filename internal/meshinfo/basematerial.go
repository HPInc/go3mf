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
func (b *BaseMaterial) HasData() bool {
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
	faceCount  uint32
	dataBlocks []*BaseMaterial
}

func newbaseMaterialContainer(currentFaceCount uint32) *baseMaterialContainer {
	m := &baseMaterialContainer{
		dataBlocks: make([]*BaseMaterial, 0, int(currentFaceCount)),
	}
	for i := 1; i <= int(currentFaceCount); i++ {
		m.AddFaceData(uint32(i))
	}
	return m
}

func (m *baseMaterialContainer) clone(currentFaceCount uint32) Container {
	return newbaseMaterialContainer(currentFaceCount)
}

func (m *baseMaterialContainer) InfoType() DataType {
	return BaseMaterialType
}

func (m *baseMaterialContainer) AddFaceData(newFaceCount uint32) FaceData {
	faceData := new(BaseMaterial)
	m.dataBlocks = append(m.dataBlocks, faceData)
	m.faceCount++
	if m.faceCount != newFaceCount {
		panic(&FaceCountMissmatchError{m.faceCount, newFaceCount})
	}
	return faceData
}

func (m *baseMaterialContainer) FaceData(faceIndex uint32) FaceData {
	return m.dataBlocks[int(faceIndex)]
}

func (m *baseMaterialContainer) FaceCount() uint32 {
	return m.faceCount
}

func (m *baseMaterialContainer) Clear() {
	m.dataBlocks = m.dataBlocks[:0]
	m.faceCount = 0
}
