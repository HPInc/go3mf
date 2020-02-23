package production

import (
	"testing"

	"github.com/qmuntal/go3mf"
)

func TestPathObject(t *testing.T) {
	type args struct {
		o            *go3mf.Object
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"emptyattr", args{&go3mf.Object{}, "/other.model"}, "/other.model"},
		{"emptypath", args{&go3mf.Object{ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{}}}, "/other.model"}, "/other.model"},
		{"emptyattr", args{&go3mf.Object{ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{Path: "/3dmodel.model"}}}, "/other.model"}, "/3dmodel.model"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PathObject(tt.args.o, tt.args.defaultValue); got != tt.want {
				t.Errorf("PathObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathItem(t *testing.T) {
	type args struct {
		o            *go3mf.Item
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"emptyattr", args{&go3mf.Item{}, "/other.model"}, "/other.model"},
		{"emptypath", args{&go3mf.Item{ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{}}}, "/other.model"}, "/other.model"},
		{"emptyattr", args{&go3mf.Item{ExtensionAttr: go3mf.ExtensionAttr{&PathUUID{Path: "/3dmodel.model"}}}, "/other.model"}, "/3dmodel.model"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PathItem(tt.args.o, tt.args.defaultValue); got != tt.want {
				t.Errorf("PathItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
