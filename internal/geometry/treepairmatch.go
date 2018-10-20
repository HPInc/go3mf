package geometry

type pairEntry struct {
	a, b int32
}

// TreePairMatch implements a n-log-n tree class which is able to identify
// duplicate pairs of numbers in a given data set.
type TreePairMatch struct {
	entries map[pairEntry]int32
}

// NewTreePairMatch creates a new TreePairMatch
func NewTreePairMatch() *TreePairMatch {
	return &TreePairMatch{map[pairEntry]int32{}}
}

// AddMatch adds a match to the set.
// If the match exists it is overriden.
func (t *TreePairMatch) AddMatch(data1, data2, param int32) {
	t.entries[newPairEntry(data1, data2)] = param
}

// CheckMatch check if a match is in the set.
func (t *TreePairMatch) CheckMatch(data1, data2 int32) (val int32, ok bool) {
	val, ok = t.entries[newPairEntry(data1, data2)]
	return
}

// DeleteMatch deletes a match from the set.
// If match doesn't exist it bevavhe as a no-op
func (t *TreePairMatch) DeleteMatch(data1, data2 int32) {
	delete(t.entries, newPairEntry(data1, data2))
}

func newPairEntry(data1, data2 int32) pairEntry {
	entry := pairEntry{}
	if data1 < data2 {
		entry.a = data1
		entry.b = data2
	} else {
		entry.a = data2
		entry.b = data1
	}
	return entry
}
