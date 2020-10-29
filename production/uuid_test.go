package production

import (
	"crypto/rand"
	"io"
	"testing"
)

func TestNewUUID(t *testing.T) {
	m := make(map[UUID]bool)
	for x := 1; x < 32; x++ {
		s := *NewUUID()
		if m[s] {
			t.Errorf("New returned duplicated UUID %s", s)
		}
		m[s] = true
		_, err := ParseUUID(string(s))
		if err != nil {
			t.Errorf("New.String() returned %q which does not decode", s)
			continue
		}
	}
}

func TestSetRand(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
	}{
		{"nil", args{nil}},
		{"rand", args{rand.Reader}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetRand(tt.args.r)
		})
	}
}
