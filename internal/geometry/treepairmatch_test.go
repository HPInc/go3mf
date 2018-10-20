package geometry

import (
	"reflect"
	"testing"
)

func TestNewTreePairMatch(t *testing.T) {
	tests := []struct {
		name string
		want *TreePairMatch
	}{
		{"new", &TreePairMatch{map[pairEntry]int32{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTreePairMatch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTreePairMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTreePairMatch_AddMatch(t *testing.T) {
	p := NewTreePairMatch()
	type args struct {
		data1 int32
		data2 int32
		param int32
	}
	tests := []struct {
		name string
		t    *TreePairMatch
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
				t.Error("TreePairMatch.AddMatch() haven't added the match")
			}
			if got != tt.args.param {
				t.Errorf("NewTreePairMatch() = %v, want %v", got, tt.args.param)
			}
		})
	}
}

func TestTreePairMatch_DeleteMatch(t *testing.T) {
	p := NewTreePairMatch()
	p.AddMatch(1, 2, 5)
	type args struct {
		data1 int32
		data2 int32
	}
	tests := []struct {
		name string
		t    *TreePairMatch
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
				t.Error("TreePairMatch.DeleteMatch() haven't deleted the match")
			}
		})
	}
}

func Test_newPairEntry(t *testing.T) {
	type args struct {
		data1 int32
		data2 int32
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
