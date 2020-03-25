package production

import (
	"testing"

	"github.com/qmuntal/go3mf"
)

var _ go3mf.SpecDecoder = new(Spec)
var _ go3mf.SpecValidator = new(Spec)
var _ go3mf.AttrMarshaler = new(UUID)
var _ go3mf.AttrMarshaler = new(PathUUID)

func TestPathUUID_ObjectPath(t *testing.T) {
	tests := []struct {
		name string
		p    *PathUUID
		want string
	}{
		{"empty", new(PathUUID), ""},
		{"path", &PathUUID{Path: "/a.model"}, "/a.model"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ObjectPath(); got != tt.want {
				t.Errorf("PathUUID.ObjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
