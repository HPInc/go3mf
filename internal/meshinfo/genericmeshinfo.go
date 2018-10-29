package meshinfo

// genericMeshInfo is used as base struct for more specific classes.
type genericMeshInfo struct {
	Container
	internalID uint64
}

// NewGenericMeshInfo creates a new genericMeshInfo.
func NewGenericMeshInfo(container Container) MeshInfo {
	return &genericMeshInfo{
		Container:  container,
		internalID: 0,
	}
}

func (b *genericMeshInfo) Clone(currentFaceCount uint32) MeshInfo {
	return NewGenericMeshInfo(b.Container.Clone(currentFaceCount))
}

func (b *genericMeshInfo) FaceHasData(faceIndex uint32) bool {
	data, err := b.GetFaceData(faceIndex)
	if err != nil {
		return false
	}
	return data.HasData()
}

func (b *genericMeshInfo) Clear() {
	count := int(b.GetCurrentFaceCount())
	for i := 0; i < count; i++ {
		b.resetFaceInformation(uint32(i))
	}
}

func (b *genericMeshInfo) resetFaceInformation(faceIndex uint32) {
	data, err := b.GetFaceData(faceIndex)
	if err != nil {
		return
	}
	data.Invalidate()
}

func (b *genericMeshInfo) cloneFaceInfosFrom(faceIndex uint32, otherInfo MeshInfo, otherFaceIndex uint32) {
	targetData, err := b.GetFaceData(faceIndex)
	if err != nil {
		return
	}
	sourceData, err := otherInfo.GetFaceData(otherFaceIndex)
	if err != nil {
		return
	}
	targetData.Copy(sourceData)
}

func (b *genericMeshInfo) permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	data, err := b.GetFaceData(faceIndex)
	if err != nil {
		return
	}
	data.Permute(nodeIndex1, nodeIndex2, nodeIndex3)
}

func (b *genericMeshInfo) mergeInformationFrom(info MeshInfo) {
	// nothing to merge
}

func (b *genericMeshInfo) setInternalID(internalID uint64) {
	b.internalID = internalID
}

func (b *genericMeshInfo) getInternalID() uint64 {
	return b.internalID
}
