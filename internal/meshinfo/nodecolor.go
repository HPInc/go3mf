package meshinfo

// NodeColor informs about the color of a node.
type NodeColor struct {
	Colors [3]Color // Colors of every vertex in a node.
}

// NewNodeColor creates a new node color form an RGB color.
func NewNodeColor(r, g, b Color) *NodeColor {
	return &NodeColor{[3]Color{r, g, b}}
}

type nodeColorInvalidator struct {
}

func (p nodeColorInvalidator) Invalidate(data FaceData) {
	if node, ok := data.(*NodeColor); ok {
		node.Colors[0] = 0x00000000
		node.Colors[1] = 0x00000000
		node.Colors[2] = 0x00000000
	}
}

// nodeColorsMeshInfo specializes the baseMeshInfo struct to "colors defined per node".
// It implements functions to interpolate and reconstruct colors while the mesh topology is changing.
type nodeColorsMeshInfo struct {
	baseMeshInfo
}

// newnodeColorsMeshInfo creates a new Node colors mesh information struct.
func newnodeColorsMeshInfo(container Container) *nodeColorsMeshInfo {
	container.Clear()
	return &nodeColorsMeshInfo{*newbaseMeshInfo(container, nodeColorInvalidator{})}
}

// GetType returns the type of information stored in this instance.
func (p *nodeColorsMeshInfo) GetType() InformationType {
	return InfoNodeColors
}

// FaceHasData checks if the specific face has any associated data.
func (p *nodeColorsMeshInfo) FaceHasData(faceIndex uint32) bool {
	data, err := p.GetFaceData(faceIndex)
	if err == nil {
		node := data.(*NodeColor)
		return (node.Colors[0] != 0) || (node.Colors[1] != 0) || (node.Colors[2] != 0)
	}
	return false
}

// Clone creates a deep copy of this instance.
func (p *nodeColorsMeshInfo) Clone() MeshInfo {
	return newnodeColorsMeshInfo(p.baseMeshInfo.Container.Clone())
}

// cloneFaceInfosFrom clones the data from another face.
func (p *nodeColorsMeshInfo) cloneFaceInfosFrom(faceIndex uint32, otherInfo MeshInfo, otherFaceIndex uint32) {
	targetData, err := p.GetFaceData(faceIndex)
	if err != nil {
		return
	}
	sourceData, err := otherInfo.GetFaceData(otherFaceIndex)
	if err != nil {
		return
	}
	node1, node2 := targetData.(*NodeColor), sourceData.(*NodeColor)
	node1.Colors[0], node1.Colors[1], node1.Colors[2] = node2.Colors[0], node2.Colors[1], node2.Colors[2]
}

//permuteNodeInformation swaps the colors.
func (p *nodeColorsMeshInfo) permuteNodeInformation(faceIndex, nodeIndex1, nodeIndex2, nodeIndex3 uint32) {
	data, err := p.GetFaceData(faceIndex)
	if err == nil && (nodeIndex1 < 3) && (nodeIndex2 < 3) && (nodeIndex3 < 3) {
		node := data.(*NodeColor)
		node.Colors[0], node.Colors[1], node.Colors[2] = node.Colors[nodeIndex1], node.Colors[nodeIndex2], node.Colors[nodeIndex3]
	}
}

// mergeInformationFrom does nothing.
func (p *nodeColorsMeshInfo) mergeInformationFrom(info MeshInfo) {
	// nothing to merge
}
