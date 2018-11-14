package go3mf

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestNewAttachement(t *testing.T) {
	type args struct {
		stream  io.Reader
		relType string
		uri     string
	}
	tests := []struct {
		name string
		args args
		want *Attachement
	}{
		{"base", args{new(bytes.Buffer), "a", "b"}, &Attachement{
			Stream:           new(bytes.Buffer),
			RelationshipType: "a",
			uri:              "b",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAttachement(tt.args.stream, tt.args.relType, tt.args.uri); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAttachement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachement_URI(t *testing.T) {
	tests := []struct {
		name string
		a    *Attachement
		want string
	}{
		{"base", NewAttachement(nil, "a", "b"), "b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.URI(); got != tt.want {
				t.Errorf("Attachement.URI() = %v, want %v", got, tt.want)
			}
		})
	}
}
