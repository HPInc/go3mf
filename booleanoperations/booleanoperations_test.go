package booleanoperations

import (
	"testing"

	"github.com/qmuntal/go3mf/spec"
)

var _ spec.MarshalerAttr = new(AssociationAttr)
var _ spec.MarshalerAttr = new(OperationAttr)

func TestComponentAttr_Components(t *testing.T) {
	tests := []struct {
		name string
		p    *AssociationAttr
		want string
	}{
		{"empty", new(AssociationAttr), ""},
		{"association", &AssociationAttr{association: Association_logical}, "logical"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.association.String(); got != tt.want {
				t.Errorf("AssociationAttr.association() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationAttr_Components(t *testing.T) {
	tests := []struct {
		name string
		p    *OperationAttr
		want string
	}{
		{"empty", new(OperationAttr), ""},
		{"association", &OperationAttr{operation: BooleanOperation_union}, "union"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.operation.String(); got != tt.want {
				t.Errorf("ComponentAttr.ObjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
