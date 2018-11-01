package geometry

import "fmt"

type InvalidUnitsError struct {
	units float32
}

func (e *InvalidUnitsError) Error() string {
	return fmt.Sprintf("the specified units (%.6f) are out of range (min: %.6f and max: %.6f)", e.units, VectorMinUnits, VectorMaxUnits)
}

type UnitsNotSettedError struct {
}

func (e *UnitsNotSettedError) Error() string {
	return "the specified units could not be set because vector dictionary already has some entries"
}
