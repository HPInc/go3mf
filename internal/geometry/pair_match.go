package geometry

type pairEntry struct {
	a, b uint32
}

// PairMatch implements a n-log-n tree class which is able to identify
// duplicate pairs of numbers in a given data set.
type PairMatch struct {
	entries map[pairEntry]uint32
}

// NewPairMatch creates a new PairMatch
func NewPairMatch() *PairMatch {
	return &PairMatch{map[pairEntry]uint32{}}
}

// AddMatch adds a match to the set.
// If the match exists it is overridden.
func (t *PairMatch) AddMatch(data1, data2, param uint32) {
	t.entries[newPairEntry(data1, data2)] = param
}

// CheckMatch check if a match is in the set.
func (t *PairMatch) CheckMatch(data1, data2 uint32) (val uint32, ok bool) {
	val, ok = t.entries[newPairEntry(data1, data2)]
	return
}

// DeleteMatch deletes a match from the set.
// If match doesn't exist it bevavhe as a no-op
func (t *PairMatch) DeleteMatch(data1, data2 uint32) {
	delete(t.entries, newPairEntry(data1, data2))
}

func newPairEntry(data1, data2 uint32) pairEntry {
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
