package mesh

import (
	"reflect"
	"testing"
)

func TestNewBeamLattice(t *testing.T) {
	tests := []struct {
		name string
		want *BeamLattice
	}{
		{"new", &BeamLattice{capMode: CapModeSphere, radius: 1.0, minLength: 0.0001}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBeamLattice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBeamLattice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeamLattice_MinLength(t *testing.T) {
	tests := []struct {
		name string
		b    *BeamLattice
		want float64
	}{
		{"new", new(BeamLattice), 0.0},
		{"base", &BeamLattice{minLength: 2.0}, 2.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.MinLength(); got != tt.want {
				t.Errorf("BeamLattice.MinLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeamLattice_Radius(t *testing.T) {
	tests := []struct {
		name string
		b    *BeamLattice
		want float64
	}{
		{"new", new(BeamLattice), 0.0},
		{"base", &BeamLattice{radius: 2.0}, 2.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Radius(); got != tt.want {
				t.Errorf("BeamLattice.Radius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeamLattice_CapMode(t *testing.T) {
	tests := []struct {
		name string
		b    *BeamLattice
		want BeamLatticeCapMode
	}{
		{"new", new(BeamLattice), CapModeSphere},
		{"base", &BeamLattice{capMode: CapModeHemisphere}, CapModeHemisphere},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.CapMode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BeamLattice.CapMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeamLattice_SetMinLength(t *testing.T) {
	type args struct {
		val float64
	}
	tests := []struct {
		name string
		b    *BeamLattice
		args args
		want float64
	}{
		{"base", new(BeamLattice), args{2.0}, 2.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetMinLength(tt.args.val)
			if got := tt.b.minLength; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BeamLattice.SetMinLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeamLattice_SetRadius(t *testing.T) {
	type args struct {
		val float64
	}
	tests := []struct {
		name string
		b    *BeamLattice
		args args
		want float64
	}{
		{"base", new(BeamLattice), args{2.0}, 2.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetRadius(tt.args.val)
			if got := tt.b.radius; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BeamLattice.SetRadius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeamLattice_SetCapMode(t *testing.T) {
	type args struct {
		val BeamLatticeCapMode
	}
	tests := []struct {
		name string
		b    *BeamLattice
		args args
		want BeamLatticeCapMode
	}{
		{"base", new(BeamLattice), args{CapModeHemisphere}, CapModeHemisphere},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetCapMode(tt.args.val)
			if got := tt.b.capMode; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BeamLattice.SetCapMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeamLattice_Clear(t *testing.T) {
	b := new(BeamLattice)
	b.Beams = append(b.Beams, Beam{})
	b.BeamSets = append(b.BeamSets, BeamSet{})
	tests := []struct {
		name string
		b    *BeamLattice
	}{
		{"base", b},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Clear()
			if len(tt.b.Beams) != 0 || len(tt.b.BeamSets) != 0 {
				t.Error("BeamLattice.Clear() have not cleared all the arrays")
			}
		})
	}
}
