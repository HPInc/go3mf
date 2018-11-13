package mesh

type comparer struct {
}

// CompareGeometry compares the geometry of two meshes to check if they are equal.
func (c comparer) CompareGeometry(m1, m2 *Mesh) bool {
	if !c.fastCheck(m1, m2) {
		return false
	}
	return c.compareNodes(m1, m2) && c.compareFaces(m1, m2) && c.compareBeams(m1, m2)
}

func (c comparer) fastCheck(m1, m2 *Mesh) bool {
	if m1 == nil || m2 == nil {
		return false
	}
	if m1 == m2 {
		return true
	}
	return len(m1.nodes) == len(m2.nodes) && len(m1.faces) == len(m2.faces) && len(m1.beams) == len(m2.beams)
}

func (c comparer) compareNodes(m1, m2 *Mesh) bool {
	for i := 0; i < len(m1.nodes); i++ {
		if !m1.nodes[i].Position.ApproxEqualThreshold(m2.nodes[i].Position, 0.0001) {
			return false
		}
	}
	return true
}

func (c comparer) compareFaces(m1, m2 *Mesh) bool {
	for i := 0; i < len(m1.faces); i++ {
		indices := m1.faces[i].NodeIndices
		other := m2.faces[i].NodeIndices
		if indices[0] != other[0] || indices[1] != other[1] || indices[2] != other[2] {
			return false
		}
	}
	return true
}

func (c comparer) compareBeams(m1, m2 *Mesh) bool {
	for i := 0; i < len(m1.beams); i++ {
		indices := m1.beams[i].NodeIndices
		other := m2.beams[i].NodeIndices
		if indices[0] != other[0] || indices[1] != other[1] {
			return false
		}
	}
	return true
}
