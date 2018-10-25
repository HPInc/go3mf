package meshinfo

// baseMeshInformation is used as base struct for more specific classes.
type baseMeshInformation struct {
	MeshInformationContainer
	Invalidator
	internalID uint64
}

// newBaseMeshInformation creates a new baseMeshInformation.
func newBaseMeshInformation(container MeshInformationContainer, invalidator Invalidator) *baseMeshInformation {
	return &baseMeshInformation{
		MeshInformationContainer: container,
		Invalidator:              invalidator,
		internalID:               0,
	}
}

// ResetFaceInformation clears the data of an specific face.
func (b *baseMeshInformation) ResetFaceInformation(faceIndex uint32) {
	data, err := b.GetFaceData(faceIndex)
	if data != nil && err == nil {
		b.Invalidator.Invalidate(data)
	}
}

// Clear resets the informations of all the faces.
func (b *baseMeshInformation) Clear() {
	count := int(b.GetCurrentFaceCount())
	for i := 0; i < count; i++ {
		b.ResetFaceInformation(uint32(i))
	}
}

// setInternalID sets an ID for the whole mesh information.
func (b *baseMeshInformation) setInternalID(internalID uint64) {
	b.internalID = internalID
}

// getInternalId gets the internal ID of the mesh information.
func (b *baseMeshInformation) getInternalID() uint64 {
	return b.internalID
}
