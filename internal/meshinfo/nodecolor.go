package meshinfo

// NodeColor informs about the color of a node.
// A node have three colors, one for each vertex of its assoicated face.
type NodeColor struct {
	Colors [3]Color // Colors of every vertex in a node.
}

// Invalidate sets to zero all the properties.
func (n *NodeColor) Invalidate() {
	n.Colors[0] = 0x00000000
	n.Colors[1] = 0x00000000
	n.Colors[2] = 0x00000000
}

// Copy copy the properties of another node color.
func (n *NodeColor) Copy(from interface{}) {
	other, ok := from.(*NodeColor)
	if !ok {
		return
	}
	n.Colors[0], n.Colors[1], n.Colors[2] = other.Colors[0], other.Colors[1], other.Colors[2]
}

// HasData returns true if the any of the colors is different from zero.
func (n *NodeColor) HasData() bool {
	return (n.Colors[0] != 0) || (n.Colors[1] != 0) || (n.Colors[2] != 0)
}

// Permute swap the colors using the given indexes. Do nothing if any of the indexes is bigger than 2.
func (n *NodeColor) Permute(index1, index2, index3 uint32) {
	if (index1 > 2) || (index2 > 2) || (index3 > 2) {
		return
	}
	n.Colors[0], n.Colors[1], n.Colors[2] = n.Colors[index1], n.Colors[index2], n.Colors[index3]
}

// Merge is not necessary for a node color.
func (n *NodeColor) Merge(other interface{}) {
	// nothing to merge
}
