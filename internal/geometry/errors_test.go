package geometry

import (
	"fmt"
	"testing"
)

func TestInvalidUnitsError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *InvalidUnitsError
		want string
	}{
		{"base", &InvalidUnitsError{10.3}, fmt.Sprintf("the specified units (%.6f) are out of range (min: %.6f and max: %.6f)", 10.3, VectorMinUnits, VectorMaxUnits)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("InvalidUnitsError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnitsNotSettedError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *UnitsNotSettedError
		want string
	}{
		{"base", new(UnitsNotSettedError), "the specified units could not be set because vector dictionary already has some entries"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("UnitsNotSettedError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
