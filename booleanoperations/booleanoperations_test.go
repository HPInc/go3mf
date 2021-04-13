package booleanoperations

import (
	"testing"

	"github.com/qmuntal/go3mf/spec"
)

var _ spec.MarshalerAttr = new(BooleanOperationAttr)

func TestComponentAttr_Components(t *testing.T) {
	tests := []struct {
		name string
		p    *BooleanOperationAttr
		want string
	}{
		{"empty", new(BooleanOperationAttr), ""},
		{"association", &BooleanOperationAttr{Association: Association_logical}, "logical"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Association.String(); got != tt.want {
				t.Errorf("AssociationAttr.association() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationAttr_Components(t *testing.T) {
	tests := []struct {
		name string
		p    *BooleanOperationAttr
		want string
	}{
		{"empty", new(BooleanOperationAttr), ""},
		{"association", &BooleanOperationAttr{Operation: BooleanOperation_union}, "union"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Operation.String(); got != tt.want {
				t.Errorf("ComponentAttr.ObjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
