package meshinfo

// BaseMaterial informs about a base material.
type BaseMaterial struct {
	MaterialGroupID uint32 // Identifier of the group.
	MaterialIndex   uint32 // Index of the base material used in the group.
}

// NewBaseMaterial creates a new base material.
func NewBaseMaterial(materialGroupID, materialIndex uint32) *BaseMaterial {
	return &BaseMaterial{materialGroupID, materialIndex}
}

type baseMaterialInvalidator struct {
}

func (p baseMaterialInvalidator) Invalidate(data FaceData) {
	if node, ok := data.(*BaseMaterial); ok {
		node.MaterialGroupID = 0
		node.MaterialIndex = 0
	}
}

// BaseMaterialsMeshInfo specializes the baseMeshInfo struct to "base materials".
type BaseMaterialsMeshInfo struct {
	baseMeshInfo
}

// NewBaseMaterialsMeshInfo creates a new Base materials mesh information struct.
func NewBaseMaterialsMeshInfo(container Container) *BaseMaterialsMeshInfo {
	container.Clear()
	return &BaseMaterialsMeshInfo{*newBaseMeshInfo(container, baseMaterialInvalidator{})}
}

// GetType returns the type of information stored in this instance.
func (p *BaseMaterialsMeshInfo) GetType() InformationType {
	return InfoBaseMaterials
}

// FaceHasData checks if the specific face has any associated data.
func (p *BaseMaterialsMeshInfo) FaceHasData(faceIndex uint32) bool {
	data, err := p.GetFaceData(faceIndex)
	if err == nil {
		return data.(*BaseMaterial).MaterialGroupID != 0
	}
	return false
}

// Clone creates a deep copy of this instance.
func (p *BaseMaterialsMeshInfo) Clone() MeshInfo {
	return NewBaseMaterialsMeshInfo(p.baseMeshInfo.Container.Clone())
}

// cloneFaceInfosFrom clones the data from another face.
func (p *BaseMaterialsMeshInfo) cloneFaceInfosFrom(faceIndex uint32, otherInfo MeshInfo, otherFaceIndex uint32) {
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
func (p *BaseMaterialsMeshInfo) permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	// nothing to merge
}

// mergeInformationFrom does nothing.
func (p *BaseMaterialsMeshInfo) mergeInformationFrom(info MeshInfo) {
	// nothing to merge
}
