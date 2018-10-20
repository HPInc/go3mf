package geometry

// PairMatch defines an interface which is able to identify duplicate pairs of numbers in a given data set.
type PairMatch interface {
	// AddMatch adds a match to the set.
	AddMatch(data1, data2, param int32)
	// CheckMatch check if a match is in the set.
	CheckMatch(data1, data2 int32) (val int32, ok bool)
	// DeleteMatch deletes a match from the set.
	DeleteMatch(data1, data2 int32)
}
