package meshinfo

import (
	"fmt"
	reflect "reflect"
	"testing"
)

func TestHandlerOverflowError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *HandlerOverflowError
		want string
	}{
		{"error", new(HandlerOverflowError), "handler overflow"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("HandlerOverflowError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidInfoTypeError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *InvalidInfoTypeError
		want string
	}{
		{"error", &InvalidInfoTypeError{nil}, fmt.Sprintf("mesh information type '%v' is not supported", reflect.TypeOf(nil))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("InvalidInfoTypeError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFaceCountMissmatchError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *FaceCountMissmatchError
		want string
	}{
		{"error", &FaceCountMissmatchError{1, 2}, fmt.Sprintf("mesh information face count (%d) does not match with mesh face count (%d)", 1, 2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("FaceCountMissmatchError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFaceDataIndexError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *FaceDataIndexError
		want string
	}{
		{"error", &FaceDataIndexError{1, 2}, fmt.Sprintf("could not access face data (%d > %d)", 2, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("FaceDataIndexError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
