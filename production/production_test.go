package production

import "testing"

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
