package mesh

import (
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func Test_newbeamLattice(t *testing.T) {
	tests := []struct {
		name string
		want *beamLattice
	}{
		{"new", &beamLattice{CapMode: CapModeSphere, DefaultRadius: 1.0, MinLength: 0.0001}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newbeamLattice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newbeamLattice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_ClearBeamLattice(t *testing.T) {
	b := new(beamLattice)
	b.beams = append(b.beams, Beam{})
	b.beamSets = append(b.beamSets, BeamSet{})
	tests := []struct {
		name string
		b    *beamLattice
	}{
		{"base", b},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.ClearBeamLattice()
			if len(tt.b.beams) != 0 || len(tt.b.beamSets) != 0 {
				t.Error("beamLattice.ClearBeamLattice() have not cleared all the arrays")
			}
		})
	}
}

func Test_beamLattice_BeamCount(t *testing.T) {
	tests := []struct {
		name string
		b    *beamLattice
		want uint32
	}{
		{"zero", new(beamLattice), 0},
		{"one", &beamLattice{beams: make([]Beam, 2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.BeamCount(); got != tt.want {
				t.Errorf("beamLattice.BeamCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_Beam(t *testing.T) {
	b := new(beamLattice)
	b.beams = append(b.beams, Beam{})
	b.beams = append(b.beams, Beam{})
	type args struct {
		index uint32
	}
	tests := []struct {
		name string
		b    *beamLattice
		args args
		want *Beam
	}{
		{"zero", b, args{0}, &b.beams[0]},
		{"one", b, args{1}, &b.beams[1]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Beam(tt.args.index); got != tt.want {
				t.Errorf("beamLattice.Beam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_AddBeamSet(t *testing.T) {
	tests := []struct {
		name string
		b    *beamLattice
		want *BeamSet
	}{
		{"base", new(beamLattice), new(BeamSet)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AddBeamSet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("beamLattice.AddBeamSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_BeamSet(t *testing.T) {
	b := new(beamLattice)
	b.beamSets = append(b.beamSets, BeamSet{})
	b.beamSets = append(b.beamSets, BeamSet{})
	type args struct {
		index uint32
	}
	tests := []struct {
		name string
		b    *beamLattice
		args args
		want *BeamSet
	}{
		{"zero", b, args{0}, &b.beamSets[0]},
		{"one", b, args{1}, &b.beamSets[1]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.BeamSet(tt.args.index); got != tt.want {
				t.Errorf("beamLattice.BeamSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_AddBeam(t *testing.T) {
	type args struct {
		node1    uint32
		node2    uint32
		radius1  float64
		radius2  float64
		capMode1 BeamCapMode
		capMode2 BeamCapMode
	}
	tests := []struct {
		name      string
		b         *beamLattice
		args      args
		want      *Beam
		wantErr   bool
		wantPanic bool
	}{
		{"max", &beamLattice{maxBeamCount: 1, beams: []Beam{{}}}, args{1, 2, 1.0, 2.0, CapModeHemisphere, CapModeSphere}, nil, false, true},
		{"node1", new(beamLattice), args{1, 1, 1.0, 2.0, CapModeHemisphere, CapModeSphere}, nil, true, false},
		{"node2", new(beamLattice), args{2, 2, 1.0, 2.0, CapModeHemisphere, CapModeSphere}, nil, true, false},
		{"base", &beamLattice{beams: []Beam{{}}}, args{0, 1, 1.0, 2.0, CapModeHemisphere, CapModeSphere}, &Beam{
			NodeIndices: [2]uint32{0, 1},
			Index:       1,
			Radius:      [2]float64{1.0, 2.0},
			CapMode:     [2]BeamCapMode{CapModeHemisphere, CapModeSphere},
		}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); tt.wantPanic && r == nil {
					t.Error("beamLattice.AddBeam() want panic")
				}
			}()
			got, err := tt.b.AddBeam(tt.args.node1, tt.args.node2, tt.args.radius1, tt.args.radius2, tt.args.capMode1, tt.args.capMode2)
			if (err != nil) != tt.wantErr {
				t.Errorf("beamLattice.AddBeam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("beamLattice.AddBeam() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		{"max", &beamLattice{maxBeamCount: 1, beams: []Beam{{}, {}}}, args{0}, false},
		{"eq", &beamLattice{beams: []Beam{{NodeIndices: [2]uint32{1, 1}}}}, args{0}, false},
		{"high1", &beamLattice{beams: []Beam{{NodeIndices: [2]uint32{2, 1}}}}, args{2}, false},
		{"high2", &beamLattice{beams: []Beam{{NodeIndices: [2]uint32{1, 2}}}}, args{2}, false},
		{"good", &beamLattice{beams: []Beam{{NodeIndices: [2]uint32{1, 2}}}}, args{3}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.checkSanity(tt.args.nodeCount); got != tt.want {
				t.Errorf("beamLattice.checkSanity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_beamLattice_merge(t *testing.T) {
	type args struct {
		newNodes []uint32
	}
	tests := []struct {
		name    string
		b       *beamLattice
		args    args
		wantErr bool
		times   uint32
	}{
		{"err", &beamLattice{beams: []Beam{{}}}, args{[]uint32{0, 0}}, true, 1},
		{"zero", new(beamLattice), args{make([]uint32, 0)}, false, 0},
		{"merged", new(beamLattice), args{[]uint32{0, 1}}, false, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMesh := new(MockMergeableMesh)
			mockMesh.On("BeamCount").Return(tt.times)
			beam := &Beam{NodeIndices: [2]uint32{0, 1}, Radius: [2]float64{1.0, 2.0}, CapMode: [2]BeamCapMode{CapModeButt, CapModeHemisphere}}
			mockMesh.On("Beam", mock.Anything).Return(beam)
			if err := tt.b.merge(mockMesh, tt.args.newNodes); (err != nil) != tt.wantErr {
				t.Errorf("beamLattice.merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			emptyBeam := Beam{}
			if !tt.wantErr && len(tt.b.beams) > 0 && tt.b.beams[0] != emptyBeam {
				for i := 0; i < len(tt.b.beams); i++ {
					want := *beam
					want.Index = uint32(i)
					if got := tt.b.Beam(uint32(i)); *got != want {
						t.Errorf("beamLattice.merge() = %v, want %v", got, want)
						return
					}
				}
			}
		})
	}
}
