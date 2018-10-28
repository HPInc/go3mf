package meshinfo

// NodeColor informs about the color of a node.
type NodeColor struct {
	Colors [3]Color // Colors of every vertex in a node.
}

func (n *NodeColor) Invalidate() {
	n.Colors[0] = 0
	n.Colors[1] = 0
	n.Colors[2] = 0
}

func (n *NodeColor) Copy(from interface{}) {
	other, ok := from.(*NodeColor)
	if !ok {
		return
	}
	n.Colors[0], n.Colors[1], n.Colors[2] = other.Colors[0], other.Colors[1], other.Colors[2]
}

func (n *NodeColor) HasData() bool {
	return (n.Colors[0] != 0) || (n.Colors[1] != 0) || (n.Colors[2] != 0)
}

func (n *NodeColor) Permute(index1, index2, index3 uint32) {
	if (index1 > 2) || (index2 > 2) || (index3 > 2) {
		return
	}
	n.Colors[0], n.Colors[1], n.Colors[2] = n.Colors[index1], n.Colors[index2], n.Colors[index3]
}

func (n *NodeColor) Merge(other interface{}) {
	// nothing to merge
}
