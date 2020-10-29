package uuid

import (
	"crypto/rand"
	"io"
	"testing"
)

func TestNew(t *testing.T) {
	m := make(map[string]bool)
	for x := 1; x < 32; x++ {
		s := New()
		if m[s] {
			t.Errorf("New returned duplicated UUID %s", s)
		}
		m[s] = true
		err := Validate(string(s))
		if err != nil {
			t.Errorf("New.String() returned %q which does not decode", s)
			continue
		}
	}
}

type test struct {
	in     string
	isuuid bool
}

var tests = []test{
	{"f47ac10b-58cc-0372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-1372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-2372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-3372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-5372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-6372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-7372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-8372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-9372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-a372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-b372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-c372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-d372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-e372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-f372-8567-0e02b2c3d479", true},

	{"urn:uuid:f47ac10b-58cc-4372-0567-0e02b2c3d479", true},
	{"URN:UUID:f47ac10b-58cc-4372-0567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-0567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-1567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-2567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-3567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-4567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-5567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-6567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-7567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-8567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-9567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-a567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-b567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-c567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-d567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-e567-0e02b2c3d479", true},
	{"f47ac10b-58cc-4372-f567-0e02b2c3d479", true},

	{"UR1:UUID:f47ac10b-58cc-4372-0567-0e02b2c3d479", false},
	{"f47ac10b158cc-5372-a567-0e02b2c3d479", false},
	{"f47ac10b-58cc25372-a567-0e02b2c3d479", false},
	{"f47ac10b-58cc-53723a567-0e02b2c3d479", false},
	{"f47ac10b-58cc-5372-a56740e02b2c3d479", false},
	{"f47ac10b-58cc-5372-a567-0e02-2c3d479", false},
	{"g47ac10b-58cc-4372-a567-0e02b2c3d479", false},

	{"{f47ac10b-58cc-0372-8567-0e02b2c3d479}", true},
	{"{f47ac10b-58cc-0372-8567-0e02b2c3d479", false},
	{"f47ac10b-58cc-0372-8567-0e02b2c3d479}", false},

	{"f47ac10b58cc037285670e02b2c3d479", true},
	{"f47ac10b58cc037285670e02b2c3d47Z", false},
	{"f47ac10b58cc037285670e02b2c3d4790", false},
	{"f47ac10b58cc037285670e02b2c3d47", false},
}

func testTest(t *testing.T, in string, tt test) {
	err := Validate(in)
	if ok := (err == nil); ok != tt.isuuid {
		t.Errorf("Validate(%s) got %v expected %v", in, ok, tt.isuuid)
	}
}

func TestValidate(t *testing.T) {
	for _, tt := range tests {
		testTest(t, tt.in, tt)
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
