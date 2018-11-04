package geometry

import "fmt"

// InvalidUnitsError happens when an invalid unit is specified.
type InvalidUnitsError struct {
	units float32
}

func (e *InvalidUnitsError) Error() string {
	return fmt.Sprintf("the specified units (%.6f) are out of range (min: %.6f and max: %.6f)", e.units, VectorMinUnits, VectorMaxUnits)
}

// UnitsNotSettedError happens when the units cannot be setted because the dictionary already has some entries.
type UnitsNotSettedError struct {
}

func (e *UnitsNotSettedError) Error() string {
	return "the specified units could not be set because vector dictionary already has some entries"
}
