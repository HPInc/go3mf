package meshinfo

import (
	"fmt"
	"testing"
)

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
