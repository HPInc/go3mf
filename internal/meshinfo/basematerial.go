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

// baseMaterialsMeshInfo specializes the baseMeshInfo struct to "base materials".
type baseMaterialsMeshInfo struct {
	baseMeshInfo
}

// newbaseMaterialsMeshInfo creates a new Base materials mesh information struct.
func newbaseMaterialsMeshInfo(container Container) *baseMaterialsMeshInfo {
	container.Clear()
	return &baseMaterialsMeshInfo{*newbaseMeshInfo(container)}
}

// GetType returns the type of information stored in this instance.
func (p *baseMaterialsMeshInfo) GetType() InformationType {
	return InfoBaseMaterials
}

// FaceHasData checks if the specific face has any associated data.
func (p *baseMaterialsMeshInfo) FaceHasData(faceIndex uint32) bool {
	data, err := p.GetFaceData(faceIndex)
	if err == nil {
		return data.(*BaseMaterial).MaterialGroupID != 0
	}
	return false
}

// Clone creates a deep copy of this instance.
func (p *baseMaterialsMeshInfo) Clone(currentFaceCount uint32) MeshInfo {
	return newbaseMaterialsMeshInfo(p.baseMeshInfo.Container.Clone(currentFaceCount))
}

// cloneFaceInfosFrom clones the data from another face.
func (p *baseMaterialsMeshInfo) cloneFaceInfosFrom(faceIndex uint32, otherInfo MeshInfo, otherFaceIndex uint32) {
	targetData, err := p.GetFaceData(faceIndex)
	if err != nil {
		return
	}
	sourceData, err := otherInfo.GetFaceData(otherFaceIndex)
	if err != nil {
		return
	}
	targetData.(*BaseMaterial).MaterialGroupID = sourceData.(*BaseMaterial).MaterialGroupID
	targetData.(*BaseMaterial).MaterialIndex = sourceData.(*BaseMaterial).MaterialIndex
}

//permuteNodeInformation does nothing.
func (p *baseMaterialsMeshInfo) permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	// nothing to merge
}

// mergeInformationFrom does nothing.
func (p *baseMaterialsMeshInfo) mergeInformationFrom(info MeshInfo) {
	// nothing to merge
}
