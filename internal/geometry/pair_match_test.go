package geometry

import (
	"reflect"
	"testing"
)

func TestNewPairMatch(t *testing.T) {
	tests := []struct {
		name string
		want *PairMatch
	}{
		{"new", &PairMatch{map[pairEntry]uint32{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPairMatch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPairMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPairMatch_AddMatch(t *testing.T) {
	p := NewPairMatch()
	type args struct {
		data1 uint32
		data2 uint32
		param uint32
	}
	tests := []struct {
		name string
		t    *PairMatch
		args args
	}{
		{"new", p, args{1, 1, 2}},
		{"old", p, args{1, 1, 4}},
		{"new2", p, args{2, 1, 5}},
		{"old2", p, args{2, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.AddMatch(tt.args.data1, tt.args.data2, tt.args.param)
			got, ok := p.CheckMatch(tt.args.data1, tt.args.data2)
			if !ok {
				t.Error("PairMatch.AddMatch() haven't added the match")
				return
			}
			if got != tt.args.param {
				t.Errorf("PairMatch.CheckMatch() = %v, want %v", got, tt.args.param)
			}
		})
	}
}

func TestPairMatch_DeleteMatch(t *testing.T) {
	p := NewPairMatch()
	p.AddMatch(1, 2, 5)
	type args struct {
		data1 uint32
		data2 uint32
	}
	tests := []struct {
		name string
		t    *PairMatch
		args args
	}{
		{"nil", p, args{2, 3}},
		{"old", p, args{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.DeleteMatch(tt.args.data1, tt.args.data2)
			_, ok := p.CheckMatch(tt.args.data1, tt.args.data2)
			if ok {
				t.Error("PairMatch.DeleteMatch() haven't deleted the match")
			}
		})
	}
}

func Test_newPairEntry(t *testing.T) {
	type args struct {
		data1 uint32
		data2 uint32
	}
	tests := []struct {
		name string
		args args
		want pairEntry
	}{
		{"d1=d2", args{1, 1}, pairEntry{1, 1}},
		{"d1>d2", args{2, 1}, pairEntry{1, 2}},
		{"d1<d2", args{1, 2}, pairEntry{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPairEntry(tt.args.data1, tt.args.data2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPairEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}
