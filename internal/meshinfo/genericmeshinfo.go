package meshinfo

// genericMeshInfo is used as base struct for more specific classes.
type genericMeshInfo struct {
	Container
	internalID uint64
	infoType   InformationType
}

// newgenericMeshInfo creates a new genericMeshInfo.
func newgenericMeshInfo(container Container, infoType InformationType) *genericMeshInfo {
	return &genericMeshInfo{
		Container:  container,
		internalID: 0,
		infoType:   infoType,
	}
}

// Clone creates a deep copy of this instance.
func (b *genericMeshInfo) Clone(currentFaceCount uint32) MeshInfo {
	return newgenericMeshInfo(b.Container.Clone(currentFaceCount), b.infoType)
}

// GetType returns the type of information stored in this instance.
func (b *genericMeshInfo) InfoType() InformationType {
	return b.infoType
}

// FaceHasData checks if the specific face has any associated data.
func (b *genericMeshInfo) FaceHasData(faceIndex uint32) bool {
	data, err := b.GetFaceData(faceIndex)
	if err != nil {
		return false
	}
	return data.HasData()
}

// Clear resets the informations of all the faces.
func (b *genericMeshInfo) Clear() {
	count := int(b.GetCurrentFaceCount())
	for i := 0; i < count; i++ {
		b.resetFaceInformation(uint32(i))
	}
}

// resetFaceInformation clears the data of an specific face.
func (b *genericMeshInfo) resetFaceInformation(faceIndex uint32) {
	data, err := b.GetFaceData(faceIndex)
	if err != nil {
		return
	}
	data.Invalidate()
}

// cloneFaceInfosFrom clones the data from another face.
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

//permuteNodeInformation swaps the data.
func (b *genericMeshInfo) permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	data, err := b.GetFaceData(faceIndex)
	if err != nil {
		return
	}
	data.Permute(nodeIndex1, nodeIndex2, nodeIndex3)
}

// mergeInformationFrom does nothing.
func (b *genericMeshInfo) mergeInformationFrom(info MeshInfo) {
	// nothing to merge
}

// setInternalID sets an ID for the whole mesh information.
func (b *genericMeshInfo) setInternalID(internalID uint64) {
	b.internalID = internalID
}

// getInternalId gets the internal ID of the mesh information.
func (b *genericMeshInfo) getInternalID() uint64 {
	return b.internalID
}
