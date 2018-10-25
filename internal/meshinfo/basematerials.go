package meshinfo

type baseMaterialInvalidator struct {
}

func (p baseMaterialInvalidator) Invalidate(data FaceData) {
	if node, ok := data.(*BaseMaterial); ok {
		node.MaterialGroupID = 0
		node.MaterialIndex = 0
	}
}

// BaseMaterialsInformation specializes the baseMeshInformation struct to "base materials".
type BaseMaterialsInformation struct {
	baseMeshInformation
}

// NewBaseMaterialsInformation creates a new Base materials mesh information struct.
func NewBaseMaterialsInformation(container MeshInformationContainer) *BaseMaterialsInformation {
	container.Clear()
	return &BaseMaterialsInformation{*newBaseMeshInformation(container, baseMaterialInvalidator{})}
}

// GetType returns the type of information stored in this instance.
func (p *BaseMaterialsInformation) GetType() InformationType {
	return InfoBaseMaterials
}

// FaceHasData checks if the specific face has any associated data.
func (p *BaseMaterialsInformation) FaceHasData(faceIndex uint32) bool {
	data, err := p.GetFaceData(faceIndex)
	if err == nil {
		return data.(*BaseMaterial).MaterialGroupID != 0
	}
	return false
}

// Clone creates a deep copy of this instance.
func (p *BaseMaterialsInformation) Clone() MeshInformation {
	return NewBaseMaterialsInformation(p.clone())
}

// cloneFaceInfosFrom clones the data from another face.
func (p *BaseMaterialsInformation) cloneFaceInfosFrom(faceIndex uint32, otherInfo MeshInformation, otherFaceIndex uint32) {
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

//permuteNodeInformation swap the data of the target mesh.
func (p *BaseMaterialsInformation) permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	// nothing to merge
}

// mergeInformationFrom merges the information of the input mesh with the current information.
func (p *BaseMaterialsInformation) mergeInformationFrom(info MeshInformation) {
	// nothing to merge
}
