package meshinfo

// GenericMeshInfo is used as base struct for more specific classes.
type GenericMeshInfo struct {
	Container
	internalID uint64
}

// NewGenericMeshInfo creates a new GenericMeshInfo.
func NewGenericMeshInfo(container Container) *GenericMeshInfo {
	return &GenericMeshInfo{
		Container:  container,
		internalID: 0,
	}
}

func (b *GenericMeshInfo) clone(currentFaceCount uint32) MeshInfo {
	return NewGenericMeshInfo(b.Container.clone(currentFaceCount))
}

// FaceHasData checks if the specific face has any associated data.
func (b *GenericMeshInfo) FaceHasData(faceIndex uint32) bool {
	data, err := b.GetFaceData(faceIndex)
	if err != nil {
		return false
	}
	return data.HasData()
}

// Clear removes all the information stored in the container.
func (b *GenericMeshInfo) Clear() {
	count := int(b.GetCurrentFaceCount())
	for i := 0; i < count; i++ {
		b.resetFaceInformation(uint32(i))
	}
}

// resetFaceInformation clears the data of an specific face.
func (b *GenericMeshInfo) resetFaceInformation(faceIndex uint32) {
	data, err := b.GetFaceData(faceIndex)
	if err != nil {
		return
	}
	data.Invalidate()
}

// cloneFaceInfosFrom clones the data from another face.
func (b *GenericMeshInfo) cloneFaceInfosFrom(faceIndex uint32, otherInfo FaceQuerier, otherFaceIndex uint32) {
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

// permuteNodeInformation swap the data of the target mesh.
func (b *GenericMeshInfo) permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	data, err := b.GetFaceData(faceIndex)
	if err != nil {
		return
	}
	data.Permute(nodeIndex1, nodeIndex2, nodeIndex3)
}

// setInternalID sets an ID for the whole mesh information.
func (b *GenericMeshInfo) setInternalID(internalID uint64) {
	b.internalID = internalID
}

// getInternalID gets the internal ID of the mesh information.
func (b *GenericMeshInfo) getInternalID() uint64 {
	return b.internalID
}
