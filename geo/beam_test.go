package geo

import (
	"testing"
)

func Test_beamLattice_checkSanity(t *testing.T) {
	type args struct {
		nodeCount uint32
	}
	tests := []struct {
		name string
		b    *beamLattice
		args args
		want bool
	}{
		{"eq", &beamLattice{Beams: []Beam{{NodeIndices: [2]uint32{1, 1}}}}, args{0}, false},
		{"high1", &beamLattice{Beams: []Beam{{NodeIndices: [2]uint32{2, 1}}}}, args{2}, false},
		{"high2", &beamLattice{Beams: []Beam{{NodeIndices: [2]uint32{1, 2}}}}, args{2}, false},
		{"good", &beamLattice{Beams: []Beam{{NodeIndices: [2]uint32{1, 2}}}}, args{3}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.checkSanity(tt.args.nodeCount); got != tt.want {
				t.Errorf("beamLattice.checkSanity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapMode_String(t *testing.T) {
	tests := []struct {
		name string
		b    CapMode
	}{
		{"sphere", CapModeSphere},
		{"hemisphere", CapModeHemisphere},
		{"butt", CapModeButt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.name {
				t.Errorf("CapMode.String() = %v, want %v", got, tt.name)
			}
		})
	}
}
