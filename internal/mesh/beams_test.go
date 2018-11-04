package mesh

import (
	"reflect"
	"testing"
)

func Test_newbeamLattice(t *testing.T) {
	tests := []struct {
		name string
		want *beamLattice
	}{
		{"new", &beamLattice{capMode: CapModeSphere, radius: 1.0, minLength: 0.0001}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newbeamLattice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newbeamLattice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_BeamLatticeMinLength(t *testing.T) {
	tests := []struct {
		name string
		b    *beamLattice
		want float64
	}{
		{"new", new(beamLattice), 0.0},
		{"base", &beamLattice{minLength: 2.0}, 2.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.BeamLatticeMinLength(); got != tt.want {
				t.Errorf("beamLattice.BeamLatticeMinLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_DefaultBeamLatticeRadius(t *testing.T) {
	tests := []struct {
		name string
		b    *beamLattice
		want float64
	}{
		{"new", new(beamLattice), 0.0},
		{"base", &beamLattice{radius: 2.0}, 2.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.DefaultBeamLatticeRadius(); got != tt.want {
				t.Errorf("beamLattice.DefaultBeamLatticeRadius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_BeamLatticeCapMode(t *testing.T) {
	tests := []struct {
		name string
		b    *beamLattice
		want BeamCapMode
	}{
		{"new", new(beamLattice), CapModeSphere},
		{"base", &beamLattice{capMode: CapModeHemisphere}, CapModeHemisphere},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.BeamLatticeCapMode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("beamLattice.BeamLatticeCapMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_SetBeamLatticeMinLength(t *testing.T) {
	type args struct {
		val float64
	}
	tests := []struct {
		name string
		b    *beamLattice
		args args
		want float64
	}{
		{"base", new(beamLattice), args{2.0}, 2.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetBeamLatticeMinLength(tt.args.val)
			if got := tt.b.minLength; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("beamLattice.SetBeamLatticeMinLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_SetDefaultBeamRadius(t *testing.T) {
	type args struct {
		val float64
	}
	tests := []struct {
		name string
		b    *beamLattice
		args args
		want float64
	}{
		{"base", new(beamLattice), args{2.0}, 2.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetDefaultBeamRadius(tt.args.val)
			if got := tt.b.radius; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("beamLattice.SetDefaultBeamRadius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_SetBeamLatticeCapMode(t *testing.T) {
	type args struct {
		val BeamCapMode
	}
	tests := []struct {
		name string
		b    *beamLattice
		args args
		want BeamCapMode
	}{
		{"base", new(beamLattice), args{CapModeHemisphere}, CapModeHemisphere},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetBeamLatticeCapMode(tt.args.val)
			if got := tt.b.capMode; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("beamLattice.SetBeamLatticeCapMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_ClearBeamLattice(t *testing.T) {
	b := new(beamLattice)
	b.Beams = append(b.Beams, new(Beam))
	b.BeamSets = append(b.BeamSets, new(BeamSet))
	tests := []struct {
		name string
		b    *beamLattice
	}{
		{"base", b},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.ClearBeamLattice()
			if len(tt.b.Beams) != 0 || len(tt.b.BeamSets) != 0 {
				t.Error("beamLattice.ClearBeamLattice() have not cleared all the arrays")
			}
		})
	}
}
