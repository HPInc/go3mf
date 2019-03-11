package meshinfo

import (
	"image/color"
)

// NodeColor informs about the color of a node.
// A node have three colors, one for each vertex of its associated face.
type NodeColor struct {
	Colors [3]color.RGBA // Colors of every vertex in a node.
}

// Invalidate sets to zero all the properties.
func (n *NodeColor) Invalidate() {
	n.Colors[0] = color.RGBA{}
	n.Colors[1] = color.RGBA{}
	n.Colors[2] = color.RGBA{}
}

// Copy copy the properties of another node color.
func (n *NodeColor) Copy(from FaceData) {
	other, ok := from.(*NodeColor)
	if !ok {
		return
	}
	n.Colors[0], n.Colors[1], n.Colors[2] = other.Colors[0], other.Colors[1], other.Colors[2]
}

// HasData returns true if the any of the colors is different from zero.
func (n *NodeColor) HasData() bool {
	return (n.Colors[0] != color.RGBA{}) || (n.Colors[1] != color.RGBA{}) || (n.Colors[2] != color.RGBA{})
}

// Permute swap the colors using the given indexes. Do nothing if any of the indexes is bigger than 2.
func (n *NodeColor) Permute(index1, index2, index3 uint32) {
	if (index1 > 2) || (index2 > 2) || (index3 > 2) {
		return
	}
	n.Colors[0], n.Colors[1], n.Colors[2] = n.Colors[index1], n.Colors[index2], n.Colors[index3]
}

// Merge is not necessary for a node color.
func (n *NodeColor) Merge(other FaceData) {
	// nothing to merge
}

type nodeColorContainer struct {
	dataBlocks []*NodeColor
}

func newnodeColorContainer(currentFaceCount uint32) *nodeColorContainer {
	m := &nodeColorContainer{
		dataBlocks: make([]*NodeColor, 0, int(currentFaceCount)),
	}
	for i := uint32(1); i <= currentFaceCount; i++ {
		m.AddFaceData(i)
	}
	return m
}

func (m *nodeColorContainer) clone(currentFaceCount uint32) Container {
	return newnodeColorContainer(currentFaceCount)
}

func (m *nodeColorContainer) InfoType() DataType {
	return NodeColorType
}

func (m *nodeColorContainer) AddFaceData(newFaceCount uint32) FaceData {
	faceData := new(NodeColor)
	m.dataBlocks = append(m.dataBlocks, faceData)
	if len(m.dataBlocks) != int(newFaceCount) {
		panic(errFaceCountMissmatch)
	}
	return faceData
}

func (m *nodeColorContainer) FaceData(faceIndex uint32) FaceData {
	return m.dataBlocks[int(faceIndex)]
}

func (m *nodeColorContainer) FaceCount() uint32 {
	return uint32(len(m.dataBlocks))
}

func (m *nodeColorContainer) Clear() {
	m.dataBlocks = m.dataBlocks[:0]
}
