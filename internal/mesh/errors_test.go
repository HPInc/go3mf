package mesh

import (
	"fmt"
	"testing"
)

func TestDuplicatedNodeError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *DuplicatedNodeError
		want string
	}{
		{"base", new(DuplicatedNodeError), "an Edge with two identical nodes has been tried to add to a mesh"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("DuplicatedNodeError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxFaceError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *MaxFaceError
		want string
	}{
		{"base", new(MaxFaceError), fmt.Sprintf("a Face has been tried to add to a mesh with too many faces (%d)", MaxFaceCount)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("MaxFaceError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxNodeError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *MaxNodeError
		want string
	}{
		{"base", new(MaxNodeError), fmt.Sprintf("a Node has been tried to add to a mesh with too many nodes (%d)", MaxNodeCount)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("MaxNodeError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxBeamError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *MaxBeamError
		want string
	}{
		{"base", new(MaxBeamError), fmt.Sprintf("a Beam has been tried to add to a mesh with too many beams (%d)", MaxBeamCount)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("MaxBeamError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
