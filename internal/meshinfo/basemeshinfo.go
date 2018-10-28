package meshinfo

// baseMeshInfo is used as base struct for more specific classes.
type baseMeshInfo struct {
	Container
	internalID uint64
}

// newbaseMeshInfo creates a new baseMeshInfo.
func newbaseMeshInfo(container Container) *baseMeshInfo {
	return &baseMeshInfo{
		Container:  container,
		internalID: 0,
	}
}

// resetFaceInformation clears the data of an specific face.
func (b *baseMeshInfo) resetFaceInformation(faceIndex uint32) {
	data, err := b.GetFaceData(faceIndex)
	if err == nil {
		data.Invalidate()
	}
}

// Clear resets the informations of all the faces.
func (b *baseMeshInfo) Clear() {
	count := int(b.GetCurrentFaceCount())
	for i := 0; i < count; i++ {
		b.resetFaceInformation(uint32(i))
	}
}

// setInternalID sets an ID for the whole mesh information.
func (b *baseMeshInfo) setInternalID(internalID uint64) {
	b.internalID = internalID
}

// getInternalId gets the internal ID of the mesh information.
func (b *baseMeshInfo) getInternalID() uint64 {
	return b.internalID
}
