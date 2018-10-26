package meshinfo

// baseMeshInfo is used as base struct for more specific classes.
type baseMeshInfo struct {
	Container
	Invalidator
	internalID uint64
}

// newbaseMeshInfo creates a new baseMeshInfo.
func newbaseMeshInfo(container Container, invalidator Invalidator) *baseMeshInfo {
	return &baseMeshInfo{
		Container:   container,
		Invalidator: invalidator,
		internalID:  0,
	}
}

// ResetFaceInformation clears the data of an specific face.
func (b *baseMeshInfo) ResetFaceInformation(faceIndex uint32) {
	data, err := b.GetFaceData(faceIndex)
	if err == nil {
		b.Invalidator.Invalidate(data)
	}
}

// Clear resets the informations of all the faces.
func (b *baseMeshInfo) Clear() {
	count := int(b.GetCurrentFaceCount())
	for i := 0; i < count; i++ {
		b.ResetFaceInformation(uint32(i))
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
